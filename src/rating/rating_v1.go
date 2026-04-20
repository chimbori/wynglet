// Package rating provides HTTP handlers for embedded rating widgets.
// It supports two UI types: thumbs (👍/👎) and stars (⭐ 1-5).
// Ratings are validated against an allowlist of domains stored in the database.
// Rate limiting is applied per IP, with an option to disable in debug mode.
package rating

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"butterfly.chimbori.dev/conf"
	"butterfly.chimbori.dev/core"
	"butterfly.chimbori.dev/db"
	"butterfly.chimbori.dev/validation"
	"github.com/lmittmann/tint"
)

//go:embed widget.min.css
var widgetCSS string

//go:embed success.min.css
var successCSS string

const (
	uiThumbs = "thumbs"
	uiStars  = "stars"
)

// ratingButton represents a single rating button in the UI.
// It contains the emoji display, numeric value, and accessibility text.
type ratingButton struct {
	Emoji   string // The emoji character to display
	Value   int16  // The numeric value submitted when clicked
	AltText string // Accessibility label for screen readers
}

// Init registers the rating service HTTP handlers on the provided mux.
// Routes:
//   - GET /rating/v1: Returns an embedded rating widget (iframe-safe)
//   - POST /rating/v1/rate: Processes a rating submission
func Init(mux *http.ServeMux) {
	mux.HandleFunc("GET /rating/v1", handleRatingWidget)
	mux.HandleFunc("POST /rating/v1/rate", handleRate)
}

// handleRatingWidget serves the interactive rating widget.
// It validates the URL against the domain allowlist, selects the appropriate UI (thumbs or stars),
// and renders the widget as a self-contained HTML page suitable for embedding in iframes.
// Query parameters:
//   - url: The URL to rate (required)
//   - ui: The UI type, either "thumbs" or "stars" (optional, defaults to "thumbs")
func handleRatingWidget(w http.ResponseWriter, req *http.Request) {
	reqURL := req.URL.Query().Get("url")
	ui := req.URL.Query().Get("ui")
	queries := db.New(db.Pool)

	url, hostname, err := validation.ValidateUrl(req.Context(), queries, reqURL)
	if err != nil {
		// Set CSP headers even for UI validation errors, otherwise clients won’t see it.
		if hostname != "" {
			setFrameHeadersForDebugMode(w, hostname)
		}
		slog.Error("URL validation failed", tint.Err(err),
			"method", req.Method,
			"path", req.URL.Path,
			"url", reqURL,
			"hostname", hostname,
			"status", http.StatusUnauthorized)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	ipAddress := core.NormalizeClientIP(core.ReadUserIP(req))

	// Determine if we should show rating buttons based on duplicate detection
	showButtons := true
	if !conf.Config.Debug && !conf.IsDebugModeActive(hostname) {
		// Check if user already has a recent rating for this URL
		exists, err := queries.HasRecentRatingByIPForURL(req.Context(), db.HasRecentRatingByIPForURLParams{
			Url:       url,
			IpAddress: ipAddress,
		})
		if err != nil {
			slog.Error("failed to check for duplicate rating", tint.Err(err),
				"method", req.Method,
				"path", req.URL.Path,
				"url", url,
				"hostname", hostname)
			// Continue anyway, showing empty widget on error
		}
		showButtons = !exists
	}

	var buttons []ratingButton
	css := ""
	if showButtons {
		buttons, err = ratingButtonsForUI(ui)
		if err != nil {
			slog.Error("invalid rating UI", tint.Err(err),
				"method", req.Method,
				"path", req.URL.Path,
				"url", url,
				"hostname", hostname,
				"status", http.StatusBadRequest)
			// Set CSP headers even for UI validation errors, otherwise clients won’t see it.
			setFrameHeadersForDebugMode(w, hostname)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		css = widgetCSS
	}

	w.Header().Set("Cache-Control", "no-store")
	setFrameHeadersForDebugMode(w, hostname)
	if err := RatingWidget(url, ui, buttons, css).Render(req.Context(), w); err != nil {
		slog.Error("failed to render rating widget", tint.Err(err),
			"method", req.Method,
			"path", req.URL.Path,
			"url", url,
			"hostname", hostname,
			"status", http.StatusInternalServerError)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// handleRate processes rating submissions from the widget.
// It validates the URL and rating, checks for duplicates per IP (unless debug mode is enabled),
// records the rating in the database, and returns a success confirmation page.
// Form parameters:
//   - url: The URL that was rated
//   - ui: The UI type used (thumbs or stars)
//   - rating: The numeric rating value submitted
func handleRate(w http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		slog.Error("failed to parse form", tint.Err(err),
			"method", req.Method,
			"path", req.URL.Path,
			"status", http.StatusBadRequest,
			"url", "",
			"hostname", "",
			"ip", "")
		http.Error(w, "invalid form payload", http.StatusBadRequest)
		return
	}

	reqURL := req.FormValue("url")
	ui := req.FormValue("ui")
	ratingStr := req.FormValue("rating")
	queries := db.New(db.Pool)

	url, hostname, err := validation.ValidateUrl(req.Context(), queries, reqURL)
	if err != nil {
		slog.Error("URL validation failed", tint.Err(err),
			"method", req.Method,
			"path", req.URL.Path,
			"url", reqURL,
			"hostname", hostname,
			"status", http.StatusUnauthorized)
		// Set CSP headers even for UI validation errors, otherwise clients won’t see it.
		if hostname != "" {
			setFrameHeadersForDebugMode(w, hostname)
		}
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	rating, err := parseRating(ui, ratingStr)
	if err != nil {
		slog.Error("invalid rating payload", tint.Err(err),
			"method", req.Method,
			"path", req.URL.Path,
			"url", url,
			"hostname", hostname,
			"status", http.StatusBadRequest)
		// Set CSP headers even for UI validation errors, otherwise clients won’t see it.
		setFrameHeadersForDebugMode(w, hostname)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ipAddress := core.NormalizeClientIP(core.ReadUserIP(req))
	if ipAddress == "" {
		slog.Error("unable to read client IP",
			"method", req.Method,
			"path", req.URL.Path,
			"url", url,
			"hostname", hostname,
			"status", http.StatusBadRequest)
		// Set CSP headers even for UI validation errors, otherwise clients won’t see it.
		setFrameHeadersForDebugMode(w, hostname)
		http.Error(w, "unable to read client IP", http.StatusBadRequest)
		return
	}

	duplicate, err := recordRating(req.Context(), url, ui, rating, ipAddress)
	if err != nil {
		slog.Error("failed to record rating", tint.Err(err),
			"method", req.Method,
			"path", req.URL.Path,
			"url", url,
			"hostname", hostname,
			"ip", ipAddress,
			"status", http.StatusInternalServerError)
		// Set CSP headers even for UI validation errors, otherwise clients won’t see it.
		setFrameHeadersForDebugMode(w, hostname)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if duplicate {
		slog.Warn("duplicate rating blocked",
			"method", req.Method,
			"path", req.URL.Path,
			"url", url,
			"hostname", hostname,
			"ip", ipAddress,
			"debug", conf.Config.Debug,
			"status", http.StatusTooManyRequests)
		// Set CSP headers even for UI validation errors, otherwise clients won’t see it.
		setFrameHeadersForDebugMode(w, hostname)
		http.Error(w, "rating already recorded for this URL", http.StatusTooManyRequests)
		return
	}

	slog.Info("rating recorded",
		"method", req.Method,
		"path", req.URL.Path,
		"url", url,
		"hostname", hostname,
		"ip", ipAddress,
		"rating", rating,
		"ui", ui,
		"status", http.StatusCreated)

	setFrameHeadersForDebugMode(w, hostname)
	w.WriteHeader(http.StatusCreated)
	if err := RatingSuccess(successCSS).Render(req.Context(), w); err != nil {
		slog.Error("failed to render rating success page", tint.Err(err),
			"method", req.Method,
			"path", req.URL.Path,
			"url", url,
			"hostname", hostname,
			"ip", ipAddress,
			"status", http.StatusInternalServerError)
	}
}

// setFrameHeadersForDebugMode sets CSP and X-Frame-Options headers based on debug mode status.
// When debug mode is active for the hostname, it allows unrestricted iframe embedding.
// Otherwise, it maintains restrictive security headers.
func setFrameHeadersForDebugMode(w http.ResponseWriter, hostname string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	cspParts := []string{
		"default-src 'self'",
		"img-src 'self' data:",
		"script-src 'self'",
		"style-src 'self' 'unsafe-inline'",
		"object-src 'none'",
		"base-uri 'self'",
	}
	if conf.IsDebugModeActive(hostname) {
		// In Debug mode, override restrictive middleware headers to allow iframe embedding.
		w.Header().Del("X-Frame-Options")
		cspParts = append(cspParts, "frame-ancestors *")
	} else {
		// Normal mode: keep restrictive headers set in [core.SecurityHeaders] middleware.
		w.Header().Set("X-Frame-Options", "DENY")
		cspParts = append(cspParts, "frame-ancestors 'https://"+hostname+"'")
	}
	w.Header().Set("Content-Security-Policy", strings.Join(cspParts, "; "))
}

// recordRating inserts a rating into the database with duplicate detection.
// It checks for recent ratings from the same IP, except in debug mode (conf.Config.Debug),
// duplicate detection is skipped to allow unlimited ratings.
// Returns (isDuplicate, error): isDuplicate is true if a recent rating was found from this IP.
func recordRating(ctx context.Context, url string, ui string, rating int16, ipAddress string) (bool, error) {
	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		return false, err
	}
	defer tx.Rollback(ctx)

	q := db.New(tx)
	// In debug mode (conf.Config.Debug set), skip duplicate check and allow unlimited ratings per IP.
	if !conf.Config.Debug {
		exists, err := q.HasRecentRatingByIPForURL(ctx, db.HasRecentRatingByIPForURLParams{
			Url:       url,
			IpAddress: ipAddress,
		})
		if err != nil {
			return false, err
		}
		if exists {
			return true, nil
		}
	}

	err = q.InsertRating(ctx, db.InsertRatingParams{
		Url:       url,
		Ui:        ui,
		Rating:    rating,
		IpAddress: ipAddress,
	})
	if err != nil {
		return false, err
	}

	if err := tx.Commit(ctx); err != nil {
		return false, err
	}
	return false, nil
}

// ratingButtonsForUI returns the button definitions for the specified UI type.
// Supports "thumbs" (👍/👎) and "stars" (⭐ 1-5 stars).
func ratingButtonsForUI(ui string) ([]ratingButton, error) {
	switch ui {
	case uiThumbs:
		return []ratingButton{
			{Emoji: "👍", Value: 1, AltText: "Thumbs Up"},
			{Emoji: "👎", Value: -1, AltText: "Thumbs Down"},
		}, nil
	case uiStars:
		return []ratingButton{
			{Emoji: "⭐", Value: 1, AltText: "1 star"},
			{Emoji: "⭐", Value: 2, AltText: "2 stars"},
			{Emoji: "⭐", Value: 3, AltText: "3 stars"},
			{Emoji: "⭐", Value: 4, AltText: "4 stars"},
			{Emoji: "⭐", Value: 5, AltText: "5 stars"},
		}, nil
	default:
		return nil, errors.New("invalid ui parameter, expected thumbs or stars")
	}
}

// parseRating validates that the provided rating value matches the selected UI type.
// For thumbs UI, accepts -1 (👎) or 1 (👍).
// For stars UI, accepts 1-5.
func parseRating(ui string, rating string) (int16, error) {
	buttons, err := ratingButtonsForUI(ui)
	if err != nil {
		return 0, err
	}

	var parsed int16
	if _, err := fmt.Sscanf(rating, "%d", &parsed); err != nil {
		return 0, errors.New("invalid rating parameter")
	}

	for _, b := range buttons {
		if b.Value == parsed {
			return parsed, nil
		}
	}
	return 0, errors.New("rating value does not match selected UI")
}
