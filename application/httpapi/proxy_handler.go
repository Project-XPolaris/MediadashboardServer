package httpapi

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"mediadashboard/config"
	"mediadashboard/plugins"

	"github.com/allentom/haruka"
)

func RegisterServiceProxyHandler(e *haruka.Engine) {
	for _, handler := range config.Instance.ServiceProxy {
		if handler.UseNacos {
			RegisterNacosProxy(e, handler)
			continue
		}
		RegisterProxy(e, handler)
	}
}

func RegisterProxy(e *haruka.Engine, p *config.ProxyConfig) {
	e.Router.Prefix(p.Prefix, func(c *haruka.Context) {
		target, err := url.Parse(p.Target)
		if err != nil {
			AbortError(c, err, http.StatusInternalServerError)
			return
		}

		// 替换请求路径
		c.Request.URL.Path = strings.Replace(c.Request.URL.Path, p.Prefix, "", 1)
		// 输出请求日志
		log.Println("proxy request", c.Request.URL.Path, "target", target.String())
		proxy := httputil.NewSingleHostReverseProxy(target)
		proxy.ServeHTTP(c.Writer, c.Request)
	})
}

// RegisterNacosProxy 在请求时通过 Nacos 选择健康实例并转发
func RegisterNacosProxy(e *haruka.Engine, p *config.ProxyConfig) {
	e.Router.Prefix(p.Prefix, func(c *haruka.Context) {
		targetStr := discoverFromNacos(p)
		if targetStr == "" {
			AbortError(c, fmt.Errorf("nacos discovery failed for %s", p.ServiceName), http.StatusBadGateway)
			return
		}
		target, err := url.Parse(targetStr)
		if err != nil {
			AbortError(c, err, http.StatusInternalServerError)
			return
		}
		c.Request.URL.Path = strings.Replace(c.Request.URL.Path, p.Prefix, "", 1)
		log.Println("proxy request", c.Request.URL.Path, "target", target.String())
		proxy := httputil.NewSingleHostReverseProxy(target)
		proxy.ServeHTTP(c.Writer, c.Request)
	})
}

func discoverFromNacos(p *config.ProxyConfig) string {
	// 统一复用 harukap 的插件，避免重复创建客户端
	if plugins.DefaultNacosPlugin == nil {
		return ""
	}
	group := p.Group
	if group == "" {
		group = "DEFAULT_GROUP"
	}
	ins, err := plugins.DefaultNacosPlugin.GetServiceInstance(p.ServiceName, group)
	if err != nil || ins == nil {
		log.Println("nacos select instance failed:", p.ServiceName, err)
		return ""
	}
	scheme := p.Scheme
	if scheme == "" {
		scheme = "http"
	}
	return fmt.Sprintf("%s://%s:%d", scheme, ins.Ip, ins.Port)
}

var GetServiceListProxyHandler = func(c *haruka.Context) {
	c.JSON(haruka.JSON{
		"success": true,
		"data":    config.Instance.ServiceProxy,
	})
}
