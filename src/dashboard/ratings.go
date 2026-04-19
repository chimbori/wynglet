package dashboard

import (
	"fmt"
	"log/slog"
	"math"
	"net/http"
	"net/url"
	"strconv"

	"butterfly.chimbori.dev/db"
	"butterfly.chimbori.dev/validation"
	"github.com/lmittmann/tint"
)

// GET /dashboard/ratings
func ratingsPageHandler(w http.ResponseWriter, req *http.Request) {
	slog.Debug("ratingsPageHandler", "url", req.Method+" "+req.URL.String())

	page := 1
	if pageStr := req.URL.Query().Get("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}
	RatingsPageTempl(page).Render(req.Context(), w)
}

// GET /dashboard/ratings/embed?url=...&ui=thumbs|stars
func ratingsEmbedBuilderHandler(w http.ResponseWriter, req *http.Request) {
	rawURL := req.URL.Query().Get("url")
	ui := req.URL.Query().Get("ui")
	if ui == "" {
		ui = "thumbs"
	}
	if ui != "thumbs" && ui != "stars" {
		RatingsIframeCodeErrorTempl("UI must be either thumbs or stars.").Render(req.Context(), w)
		return
	}

	if rawURL == "" {
		RatingsIframeCodeErrorTempl("Paste a URL to generate an embed iframe.").Render(req.Context(), w)
		return
	}

	queries := db.New(db.Pool)
	validatedURL, hostname, err := validation.ValidateUrl(req.Context(), queries, rawURL)
	if err != nil {
		slog.Warn("ratings embed validation failed", tint.Err(err),
			"method", req.Method,
			"path", req.URL.Path,
			"url", rawURL,
			"hostname", hostname,
			"status", http.StatusBadRequest)
		RatingsIframeCodeErrorTempl(err.Error()).Render(req.Context(), w)
		return
	}

	iframeCode := buildRatingsIframeEmbed(req, validatedURL, ui)
	RatingsIFrameCodeTempl(iframeCode).Render(req.Context(), w)
}

func buildRatingsIframeEmbed(req *http.Request, validatedURL string, ui string) string {
	params := url.Values{}
	params.Set("url", validatedURL)
	params.Set("ui", ui)

	scheme := req.Header.Get("X-Forwarded-Proto")
	if scheme == "" {
		if req.TLS != nil {
			scheme = "https"
		} else {
			scheme = "http"
		}
	}

	width, height := 0, 0
	switch ui {
	case "stars":
		width, height = 200, 50
	case "thumbs":
		width, height = 120, 50
	}

	iframeSrc := fmt.Sprintf("%s://%s/rating/v1?%s", scheme, req.Host, params.Encode())
	return fmt.Sprintf(
		`<iframe src="%s" width="%d" height="%d" style="border:0;overflow:hidden;" loading="lazy" referrerpolicy="no-referrer"></iframe>`,
		iframeSrc, width, height)
}

func formatThumbsSummary(thumbsUp int64, thumbsDown int64) string {
	total := thumbsUp + thumbsDown
	if total == 0 {
		return "No votes"
	}
	upPct := (float64(thumbsUp) * 100) / float64(total)
	downPct := (float64(thumbsDown) * 100) / float64(total)
	return fmt.Sprintf("👍 %.0f%% / 👎 %.0f%%", math.Round(upPct), math.Round(downPct))
}

func formatStarsSummary(avg float64, total int64) string {
	if total == 0 {
		return "No votes"
	}
	return fmt.Sprintf("%.2f/5 (%d)", avg, total)
}
