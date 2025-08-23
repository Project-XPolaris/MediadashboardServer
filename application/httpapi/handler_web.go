package httpapi

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

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

	// SPA fallback: serve index.html for non-API routes when direct access
	e.Router.Prefix("/", func(c *haruka.Context) {
		path := c.Request.URL.Path
		// bypass API
		if strings.HasPrefix(path, "/api") {
			return
		}
		// only handle GET/HEAD
		if c.Request.Method != http.MethodGet && c.Request.Method != http.MethodHead {
			return
		}

		// try serve static file if exists, else fallback to index.html
		cleaned := filepath.Clean(path)
		if strings.HasPrefix(cleaned, "..") {
			cleaned = "/"
		}
		fsPath := filepath.Join(".", "dist", cleaned)
		if info, err := os.Stat(fsPath); err == nil && !info.IsDir() {
			http.ServeFile(c.Writer, c.Request, fsPath)
			return
		}
		http.ServeFile(c.Writer, c.Request, filepath.Join(".", "dist", "index.html"))
	})
}
