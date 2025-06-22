package httpapi

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"mediadashboard/config"

	"github.com/allentom/haruka"
)

func RegisterServiceProxyHandler(e *haruka.Engine) {
	for _, handler := range config.Instance.ServiceProxy {
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

var GetServiceListProxyHandler = func(c *haruka.Context) {
	c.JSON(haruka.JSON{
		"success": true,
		"data":    config.Instance.ServiceProxy,
	})
}
