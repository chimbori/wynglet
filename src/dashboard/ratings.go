package dashboard

import (
	"fmt"
	"log/slog"
	"math"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/lmittmann/tint"
	"wynglet.chimbori.dev/conf"
	"wynglet.chimbori.dev/db"
)

// Renders the ratings dashboard page.
// GET /dashboard/ratings
func ratingsPageHandler(w http.ResponseWriter, req *http.Request) {
	slog.Debug("ratingsPageHandler", "url", req.Method+" "+req.URL.String())

	page := 1
	if pageStr := req.URL.Query().Get("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}
	days := parseRatingDays(req.URL.Query().Get("days"))
	RatingsPageTempl(page, days).Render(req.Context(), w)
}

// Returns paginated list of ratings with filtering and sorting.
// GET /dashboard/ratings/list?days=0&page=1&sortBy=rating&sortOrder=desc
func ratingsListHandler(w http.ResponseWriter, req *http.Request) {
	slog.Debug("ratingsListHandler", "url", req.Method+" "+req.URL.String())

	ctx := req.Context()
	queries := db.New(db.Pool)

	page := 1
	if pageStr := req.URL.Query().Get("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}
	days := parseRatingDays(req.URL.Query().Get("days"))
	sortBy := req.URL.Query().Get("sortBy")
	sortOrder := req.URL.Query().Get("sortOrder")

	totalCount, err := queries.CountRatingGroups(ctx, int32(days))
	if err != nil {
		slog.Error("failed to count ratings", tint.Err(err),
			"method", req.Method,
			"path", req.URL.Path,
			"url", req.URL.String(),
			"status", http.StatusInternalServerError)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	offset := int32((page - 1) * conf.Config.Dashboard.Pagination.Limit)
	ratings, err := queries.ListRatingsWithDistribution(ctx, db.ListRatingsWithDistributionParams{
		Days:             int32(days),
		PaginationLimit:  int32(conf.Config.Dashboard.Pagination.Limit),
		PaginationOffset: offset,
	})
	if err != nil {
		slog.Error("failed to list ratings", tint.Err(err),
			"method", req.Method,
			"path", req.URL.Path,
			"url", req.URL.String(),
			"status", http.StatusInternalServerError)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Client-side sorting
	sortRatings(ratings, sortBy, sortOrder)

	RatingsListTempl(ratings, page, days, totalCount, sortBy, sortOrder).Render(ctx, w)
}

// sortRatings sorts the ratings in-place based on the specified sort column and order.
func sortRatings(ratings []db.ListRatingsWithDistributionRow, sortBy, sortOrder string) {
	if sortBy == "" {
		sortBy = "rating" // Default sort by rating
	}
	if sortOrder == "" {
		sortOrder = "desc" // Default order
	}

	ascending := sortOrder == "asc"

	switch sortBy {
	case "url":
		sort.SliceStable(ratings, func(i, j int) bool {
			cmp := strings.Compare(ratings[i].Url, ratings[j].Url)
			if ascending {
				return cmp < 0
			}
			return cmp > 0
		})
	case "votes":
		sort.SliceStable(ratings, func(i, j int) bool {
			if ascending {
				return ratings[i].TotalRatings < ratings[j].TotalRatings
			}
			return ratings[i].TotalRatings > ratings[j].TotalRatings
		})
	case "rating":
		fallthrough
	default:
		sort.SliceStable(ratings, func(i, j int) bool {
			if ascending {
				return ratings[i].NormalizedScore < ratings[j].NormalizedScore
			}
			return ratings[i].NormalizedScore > ratings[j].NormalizedScore
		})
	}
}

// parseRatingDays parses the days query parameter for ratings time filtering.
// Returns 0 for "All" (default), or 1/7/28 for specific ranges.
func parseRatingDays(daysParam string) int {
	switch daysParam {
	case "1":
		return 1
	case "7":
		return 7
	case "28":
		return 28
	default:
		return 0
	}
}

func pctOf(n, total int64) float64 {
	if total == 0 {
		return 0
	}
	return math.Round(float64(n) * 100 / float64(total))
}

func formatAverageRating(r db.ListRatingsWithDistributionRow) string {
	if r.TotalRatings == 0 {
		return "—"
	}
	switch r.Ui {
	case "thumbs":
		return fmt.Sprintf("👍 %.0f%%", pctOf(r.ThumbsUp, r.TotalRatings))
	case "stars":
		return fmt.Sprintf("⭐ %.1f / 5", r.AverageStars)
	default:
		return "—"
	}
}
