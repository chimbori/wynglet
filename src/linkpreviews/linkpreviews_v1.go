package linkpreviews

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"regexp"

	"github.com/lmittmann/tint"
	"wynglet.chimbori.dev/conf"
	"wynglet.chimbori.dev/core"
	"wynglet.chimbori.dev/db"
	"wynglet.chimbori.dev/embedfs"
	"wynglet.chimbori.dev/validation"
)

var Cache *core.DiskCache

var selectorRegex = regexp.MustCompile(`^[#.][a-zA-Z0-9_-]+$`)

func Init(mux *http.ServeMux) error {
	var err error
	Cache, err = core.NewDiskCache(
		filepath.Join(conf.Config.DataDir, "cache", "link-previews"),
		core.WithTTL(conf.Config.LinkPreviews.Cache.TTL),
		core.WithMaxSize(conf.Config.LinkPreviews.Cache.MaxSizeBytes),
	)
	if err != nil {
		return err
	}

	mux.HandleFunc("GET /link-previews/v1", linkPreviewsV1Handler)
	return nil
}

// Validates the URL, checks if it’s cached, generates screenshots, and serves them.
// GET /link-previews/v1
func linkPreviewsV1Handler(w http.ResponseWriter, req *http.Request) {
	slog.Debug("handleLinkPreview", "url", req.Method+" "+req.URL.String())

	reqUrl := req.URL.Query().Get("url")
	userAgent := req.Header.Get("User-Agent")
	canonicalUserAgent := core.GetCanonicalUserAgent(userAgent)
	queries := db.New(db.Pool)

	url, hostname, err := validation.ValidateUrl(req.Context(), queries, reqUrl)
	if err != nil {
		slog.Error("URL validation failed", tint.Err(err),
			"method", req.Method,
			"path", req.URL.Path,
			"url", reqUrl,
			"hostname", hostname,
			"user-agent", userAgent,
			"status", http.StatusUnauthorized)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	selector := req.URL.Query().Get("sel")
	if selector == "" {
		selector = "#link-preview"
	} else if !selectorRegex.MatchString(selector) {
		selectorErr := fmt.Errorf("invalid selector")
		slog.Error("selector validation failed", tint.Err(selectorErr),
			"method", req.Method,
			"path", req.URL.Path,
			"url", reqUrl,
			"hostname", hostname,
			"user-agent", userAgent,
			"status", http.StatusBadRequest)
		http.Error(w, "invalid selector", http.StatusBadRequest)
		return
	}

	var cached []byte

	cached, err = Cache.Find(core.ComputeKey(url, "png", false))
	if err != nil {
		err = fmt.Errorf("url: %s, %w", url, err)
		slog.Error("error during cache lookup", tint.Err(err),
			"method", req.Method,
			"path", req.URL.Path,
			"url", url,
			"hostname", hostname,
			"user-agent", userAgent,
			"status", http.StatusInternalServerError)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if cached != nil {
		slog.Info("cached screenshot served",
			"method", req.Method,
			"path", req.URL.Path,
			"url", url,
			"hostname", hostname,
			"user-agent", userAgent,
			"status", http.StatusOK)
		w.Header().Set("Content-Type", "image/png")
		w.Header().Set("Cache-Control", "max-age=31536000, immutable") // 1 year
		w.Write(cached)
		recordLinkPreviewAccessed(url, canonicalUserAgent)

	} else {
		ctx, cancel := context.WithTimeout(req.Context(), conf.Config.LinkPreviews.Screenshot.Timeout)
		defer cancel()
		screenshot, err := core.TakeScreenshot(ctx, url, selector)
		if err != nil {
			if !errors.Is(err, core.ErrMissingSelector) {
				err = fmt.Errorf("url: %s, %w", url, err)
				slog.Error("error taking screenshot", tint.Err(err),
					"method", req.Method,
					"path", req.URL.Path,
					"url", url,
					"hostname", hostname,
					"user-agent", userAgent,
					"status", http.StatusInternalServerError)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			slog.Info("attempting with default template",
				"method", req.Method,
				"path", req.URL.Path,
				"url", url,
				"hostname", hostname,
				"user-agent", userAgent,
				"status", http.StatusOK)
			title, description, fetchErr := core.FetchTitleAndDescription(ctx, url)
			if fetchErr != nil {
				err = fmt.Errorf("fetchTitleAndDescription failed: %w", fetchErr)
			} else {
				screenshot, err = core.TakeScreenshotWithTemplate(ctx, embedfs.DefaultTemplate, url, "#link-preview", title, description)
			}
			if err != nil {
				err = fmt.Errorf("url: %s, %w", url, err)
				slog.Error("error using default template", tint.Err(err),
					"method", req.Method,
					"path", req.URL.Path,
					"url", url,
					"hostname", hostname,
					"user-agent", userAgent,
					"status", http.StatusInternalServerError)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		// Serve the screenshot immediately after generation, without waiting for compression.
		slog.Info("new screenshot generated",
			"method", req.Method,
			"path", req.URL.Path,
			"url", url,
			"hostname", hostname,
			"user-agent", userAgent,
			"status", http.StatusOK)
		w.Header().Set("Content-Type", "image/png")
		w.Header().Set("Cache-Control", "max-age=31536000, immutable") // 1 year
		w.Write(screenshot)
		recordLinkPreviewCreated(url, canonicalUserAgent)

		// Compress and cache the generated screenshot, but without holding up the HTTP request
		go func() {
			dataToWrite := screenshot
			compressed, err := core.CompressPNG(screenshot)
			if err == nil {
				dataToWrite = compressed
				slog.Info("PNG compressed", "from", len(screenshot), "to", len(compressed), "%", (len(compressed) * 100 / len(screenshot)))
			} else {
				slog.Error("PNG compression failed", tint.Err(err), "url", url)
			}

			if err := Cache.Write(core.ComputeKey(url, "png", false), dataToWrite); err != nil {
				err = fmt.Errorf("error writing to cache: %s, %w", url, err)
				slog.Error("error writing to cache", tint.Err(err),
					"method", req.Method,
					"path", req.URL.Path,
					"url", url,
					"hostname", hostname,
					"user-agent", userAgent,
					"status", http.StatusInternalServerError)
			}
		}()
	}
}

// Record when a link preview is created (for the first time)
func recordLinkPreviewCreated(url string, canonicalUserAgent string) {
	queries := db.New(db.Pool)
	err := queries.RecordLinkPreviewCreated(context.Background(), db.RecordLinkPreviewCreatedParams{
		Url:                url,
		CanonicalUserAgent: &canonicalUserAgent,
	})
	if err != nil {
		slog.Error("failed to log link preview created", tint.Err(err))
	}
	// Don’t return an error to the caller; fulfill the request anyway.
}

// Record when a link preview is accessed from the cache
func recordLinkPreviewAccessed(url string, canonicalUserAgent string) {
	queries := db.New(db.Pool)
	rowsUpdated, err := queries.RecordLinkPreviewAccessed(context.Background(), db.RecordLinkPreviewAccessedParams{
		Url:                url,
		CanonicalUserAgent: &canonicalUserAgent,
	})
	if err != nil {
		slog.Error("failed to log link preview created", tint.Err(err))
	}
	if rowsUpdated == 0 { // If not already in the database, add it now.
		recordLinkPreviewCreated(url, canonicalUserAgent)
	}
	// Don’t return an error to the caller; fulfill the request anyway.
}

// DeleteCached removes a cached screenshot file from disk.
func DeleteCached(url string) error {
	err := Cache.Delete(url)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}
