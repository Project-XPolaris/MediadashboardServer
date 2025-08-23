package httpapi

import (
	"mediadashboard/config"
	"mediadashboard/module"
	"net/http"
	"strings"

	"github.com/allentom/haruka"
	"github.com/allentom/haruka/middleware"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
)

var Logger = log.New().WithFields(log.Fields{
	"scope": "Application",
})

func GetEngine() *haruka.Engine {
	e := haruka.NewEngine()
	e.UseCors(cors.AllowAll())
	module.Auth.AuthMiddleware.OnError = func(ctx *haruka.Context, err error) {
		AbortError(ctx, err, http.StatusForbidden)
		ctx.Abort()
	}
	module.Auth.AuthMiddleware.RequestFilter = func(c *haruka.Context) bool {
		NoAuthPath := []string{
			"/api/oauth/youauth",
			"/api/oauth/youauth/password",
			"/api/oauth/youplus",
			"/api/info",
			"/api/state",
		}
		for _, path := range NoAuthPath {
			if c.Pattern == path {
				return false
			}
		}
		// Special case for video link: token is in URL path. Extract and set Authorization header so auth middleware validates it.
		if vp := getVideoPrefix(); vp != "" && strings.HasPrefix(c.Request.URL.Path, vp+"/link/") {
			rest := strings.TrimPrefix(c.Request.URL.Path, vp+"/link/")
			segs := strings.Split(rest, "/")
			if len(segs) >= 3 {
				token := segs[2]
				if token != "" {
					c.Request.Header.Set("Authorization", "Bearer "+token)
				}
			}
		}
		return true
	}
	e.UseMiddleware(module.Auth.AuthMiddleware)
	e.UseMiddleware(middleware.NewLoggerMiddleware())
	e.UseMiddleware(middleware.NewPaginationMiddleware("page", "pageSize", 1, 20))
	// e.Router.GET("/oauth/youauth", generateAccessCodeWithYouAuthHandler)
	e.Router.POST("/api/oauth/youauth/password", generateAccessCodeWithYouAuthPasswordHandler)
	e.Router.GET("/api/oauth/check", checkTokenHandler)
	// e.Router.GET("/user/auth", youPlusTokenHandler)
	// e.Router.GET("/info", serviceInfoHandler)
	// photo proxy
	e.Router.GET("/api/ping", func(c *haruka.Context) {
		c.JSON(haruka.JSON{
			"message": "pong",
		})
	})
	RegisterServiceProxyHandler(e)
	e.Router.GET("/api/service/list", GetServiceListProxyHandler)
	e.Router.GET("/api/service/{name}", GetServiceDetailProxyHandler)
	// Add UI handler
	ServeUIHandler(e)

	InitErrorHandler()
	return e
}

// getVideoPrefix returns the configured prefix for the video service
func getVideoPrefix() string {
	for _, p := range config.Instance.ServiceProxy {
		if p != nil && p.Name == "video" {
			return p.Prefix
		}
	}
	return ""
}
