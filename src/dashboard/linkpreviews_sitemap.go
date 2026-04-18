package dashboard

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"butterfly.chimbori.dev/db"
	"butterfly.chimbori.dev/linkpreviews"
	"butterfly.chimbori.dev/validation"
	"github.com/lmittmann/tint"
)

func sitemapPageHandler(w http.ResponseWriter, req *http.Request) {
	SitemapImportPageTempl(linkpreviews.GetSitemapImportStatus(), "").Render(req.Context(), w)
}

func startSitemapImportHandler(w http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		slog.Error("failed to parse sitemap import form", tint.Err(err),
			"method", req.Method,
			"path", req.URL.Path,
			"url", req.URL.String(),
			"status", http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		SitemapImportPageTempl(linkpreviews.GetSitemapImportStatus(), "Invalid form submission").Render(req.Context(), w)
		return
	}

	sitemapURL := strings.TrimSpace(req.FormValue("sitemap_url"))
	if sitemapURL == "" {
		err := fmt.Errorf("missing sitemap URL")
		slog.Error("missing sitemap URL", tint.Err(err),
			"method", req.Method,
			"path", req.URL.Path,
			"url", req.URL.String(),
			"status", http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		SitemapImportPageTempl(linkpreviews.GetSitemapImportStatus(), "Sitemap URL is required").Render(req.Context(), w)
		return
	}

	canonicalURL, err := validation.Canonicalize(sitemapURL)
	if err != nil || canonicalURL.Scheme != "https" {
		if err == nil {
			err = fmt.Errorf("unsupported URL scheme: %s", canonicalURL.Scheme)
		}
		slog.Error("invalid sitemap URL", tint.Err(err),
			"method", req.Method,
			"path", req.URL.Path,
			"url", req.URL.String(),
			"sitemap_url", sitemapURL,
			"status", http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		SitemapImportPageTempl(linkpreviews.GetSitemapImportStatus(), "Sitemap URL must be a valid HTTPS URL").Render(req.Context(), w)
		return
	}

	queries := db.New(db.Pool)
	if _, _, err := validation.ValidateUrl(req.Context(), queries, canonicalURL.String()); err != nil {
		slog.Error("sitemap URL is not authorized", tint.Err(err),
			"method", req.Method,
			"path", req.URL.Path,
			"url", req.URL.String(),
			"sitemap_url", canonicalURL.String(),
			"status", http.StatusUnauthorized)
		w.WriteHeader(http.StatusUnauthorized)
		SitemapImportPageTempl(linkpreviews.GetSitemapImportStatus(), err.Error()).Render(req.Context(), w)
		return
	}

	if err := linkpreviews.StartSitemapImport(req.Context(), canonicalURL.String()); err != nil {
		slog.Error("failed to start sitemap import", tint.Err(err),
			"method", req.Method,
			"path", req.URL.Path,
			"url", req.URL.String(),
			"sitemap_url", canonicalURL.String(),
			"status", http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		SitemapImportPageTempl(linkpreviews.GetSitemapImportStatus(), err.Error()).Render(req.Context(), w)
		return
	}

	http.Redirect(w, req, "/dashboard/link-previews/sitemap", http.StatusSeeOther)
}

func sitemapStatusHandler(w http.ResponseWriter, req *http.Request) {
	SitemapImportStatusTempl(linkpreviews.GetSitemapImportStatus()).Render(req.Context(), w)
}

func cancelSitemapImportHandler(w http.ResponseWriter, req *http.Request) {
	linkpreviews.CancelSitemapImport()
	SitemapImportStatusTempl(linkpreviews.GetSitemapImportStatus()).Render(req.Context(), w)
}

func sitemapImportProgress(status *linkpreviews.SitemapImportStatus) int {
	if status == nil || status.TotalURLs == 0 {
		return 0
	}
	processed := status.Completed + status.Skipped + status.Failed
	return min(100, processed*100/status.TotalURLs)
}

func sitemapImportSummary(status *linkpreviews.SitemapImportStatus) string {
	if status == nil {
		return "No sitemap import has run yet"
	}
	if status.InProgress {
		return fmt.Sprintf("Processing %d URLs from %s", status.TotalURLs, status.SitemapURL)
	}
	return fmt.Sprintf("Last import: %d completed, %d skipped, %d failed", status.Completed, status.Skipped, status.Failed)
}
