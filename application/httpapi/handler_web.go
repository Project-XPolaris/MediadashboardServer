package httpapi

import (
	"github.com/allentom/haruka"
)

// ServeUIHandler handles the UI static files and API proxy
func ServeUIHandler(e *haruka.Engine) {
	// Serve static files from the dist directory
	e.Router.Static("/", "./dist")

	// Setup API proxy to self
	// apiProxyConfig := &ProxyConfig{
	// 	Target: "/", // 代理到自身根路径
	// 	Prefix: "/api",
	// }
	// apiProxyConfig.Register(e)
}
