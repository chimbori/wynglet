package embedfs

import (
	"bytes"
	"context"
	"embed"
	"io"
	"net/http"
	"path/filepath"
)

//go:embed static
var staticFiles embed.FS

//go:embed "default-template.html"
var DefaultTemplate string

func ServeStaticFS(mux *http.ServeMux) {
	mux.Handle("GET /static/", maxAgeHandler(http.FileServer(http.FS(staticFiles))))
	mux.HandleFunc("GET /favicon.ico", func(w http.ResponseWriter, req *http.Request) {
		http.ServeFileFS(w, req, staticFiles, "static/favicon.svg")
	})
}

// maxAgeHandler wraps an HTTP handler to set cache control headers based on file extension.
// CSS files are cached for 1 day; all other static files for 1 year.
func maxAgeHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ext := filepath.Ext(req.URL.Path)
		if ext == ".css" {
			w.Header().Set("Cache-Control", "max-age=600") // 10 min
		} else {
			w.Header().Set("Cache-Control", "max-age=31536000, immutable") // 1 year
		}
		h.ServeHTTP(w, req)
	})
}

// [IconComponent] implements interface templ.Component, and can be used to embed
// SVG icons in Templ templates directly.
type IconComponent []byte

func (ic IconComponent) Render(ctx context.Context, w io.Writer) error {
	_, err := w.Write(bytes.ReplaceAll(ic, []byte("\n"), nil))
	return err
}

//go:embed static/bluesky.svg
var BlueSkyIcon IconComponent

//go:embed static/reddit.svg
var RedditIcon IconComponent

//go:embed static/linkedin.svg
var LinkedInIcon IconComponent

//go:embed static/hacker-news.svg
var HackerNewsIcon IconComponent
