// Package dashboard provides HTTP handlers and Templ templates for the dashboard UI,
// including forms management, viewing submissions, and exporting data.
package dashboard

import (
	"encoding/csv"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/lmittmann/tint"
	"wynglet.chimbori.dev/conf"
	"wynglet.chimbori.dev/db"
)

// formsPageHandler serves the Forms Dashboard page
// GET /dashboard/forms
func formsPageHandler(w http.ResponseWriter, r *http.Request) {
	FormsPageTempl().Render(r.Context(), w)
}

// formsListHandler returns HTML list of forms
// GET /dashboard/forms/list
func formsListHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	queries := db.New(db.Pool)

	// Fetch distinct form_ids and their counts from the database
	forms, err := queries.ListForms(ctx)
	if err != nil {
		slog.Error("Failed to fetch forms", tint.Err(err),
			"method", r.Method,
			"path", r.URL.Path,
			"url", r.URL.String(),
			"hostname", r.Host,
			"user-agent", r.UserAgent(),
			"ip", r.RemoteAddr,
			"status", http.StatusInternalServerError)
		http.Error(w, "Failed to load forms", http.StatusInternalServerError)
		return
	}

	FormsListTempl(forms).Render(ctx, w)
}

// formDetailHandler shows submissions for a specific form
// GET /dashboard/forms/{form_id}
func formDetailHandler(w http.ResponseWriter, r *http.Request) {
	formID := r.PathValue("form_id")
	if formID == "" {
		slog.Warn("Missing form_id parameter",
			"method", r.Method,
			"path", r.URL.Path,
			"url", r.URL.String(),
			"hostname", r.Host,
			"user-agent", r.UserAgent(),
			"ip", r.RemoteAddr,
			"status", http.StatusBadRequest)
		http.Error(w, "Missing form_id", http.StatusBadRequest)
		return
	}

	page := 1
	if p := r.URL.Query().Get("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	offset := int64((page - 1) * conf.Config.Dashboard.Pagination.Limit)

	ctx := r.Context()
	queries := db.New(db.Pool)

	// Count total submissions
	total, err := queries.CountFormSubmissions(ctx, formID)
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to count submissions (form_id: %s)", formID), tint.Err(err),
			"method", r.Method,
			"path", r.URL.Path,
			"url", r.URL.String(),
			"hostname", r.Host,
			"user-agent", r.UserAgent(),
			"ip", r.RemoteAddr,
			"status", http.StatusInternalServerError)
		http.Error(w, "Failed to load submissions", http.StatusInternalServerError)
		return
	}

	// Fetch submissions
	submissions, err := queries.ListFormSubmissions(ctx, db.ListFormSubmissionsParams{
		FormID: formID,
		Limit:  int32(conf.Config.Dashboard.Pagination.Limit),
		Offset: int32(offset),
	})
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to fetch submissions (form_id: %s, page: %d)", formID, page), tint.Err(err),
			"method", r.Method,
			"path", r.URL.Path,
			"url", r.URL.String(),
			"hostname", r.Host,
			"user-agent", r.UserAgent(),
			"ip", r.RemoteAddr,
			"status", http.StatusInternalServerError)
		http.Error(w, "Failed to load submissions", http.StatusInternalServerError)
		return
	}

	FormSubmissionsTempl(formID, submissions, page, total, conf.Config.Dashboard.Pagination.Limit).Render(ctx, w)
}

// formExportCSVHandler exports submissions as CSV
// GET /dashboard/forms/{form_id}/export.csv
func formExportCSVHandler(w http.ResponseWriter, r *http.Request) {
	formID := r.PathValue("form_id")
	if formID == "" {
		slog.Warn("Missing form_id parameter",
			"method", r.Method,
			"path", r.URL.Path,
			"url", r.URL.String(),
			"hostname", r.Host,
			"user-agent", r.UserAgent(),
			"ip", r.RemoteAddr,
			"status", http.StatusBadRequest)
		http.Error(w, "Missing form_id", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	// Fetch all submissions for export
	submissions, err := db.New(db.Pool).ExportFormSubmissions(ctx, formID)
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to fetch submissions for export (form_id: %s)", formID), tint.Err(err),
			"method", r.Method,
			"path", r.URL.Path,
			"url", r.URL.String(),
			"hostname", r.Host,
			"user-agent", r.UserAgent(),
			"ip", r.RemoteAddr,
			"status", http.StatusInternalServerError)
		http.Error(w, "Failed to export", http.StatusInternalServerError)
		return
	}

	// Generate CSV
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=form_%s_%s.csv", formID, time.Now().Format("2006-01-02")))

	writer := csv.NewWriter(w)
	defer writer.Flush()

	// Write header
	writer.Write([]string{"Submitted", "Domain", "IP", "Spam", "Email Sent", "Fields"})

	// Write rows
	for _, submission := range submissions {
		spam := "false"
		if submission.IsSpam {
			spam = "true"
		}
		emailSent := ""
		if submission.EmailSentAt != nil {
			emailSent = submission.EmailSentAt.Format(time.RFC3339)
		}
		writer.Write([]string{
			submission.SubmittedAt.Format(time.RFC3339),
			submission.Domain,
			submission.IpAddress,
			spam,
			emailSent,
			submission.FormData,
		})
	}
}

// formDeleteHandler deletes a specific submission and returns the updated submissions table
// DELETE /dashboard/forms/{form_id}/{id}
func formDeleteHandler(w http.ResponseWriter, r *http.Request) {
	formID := r.PathValue("form_id")
	idStr := r.PathValue("id")

	if formID == "" || idStr == "" {
		slog.Warn("Missing parameters",
			"method", r.Method,
			"path", r.URL.Path,
			"url", r.URL.String(),
			"hostname", r.Host,
			"user-agent", r.UserAgent(),
			"ip", r.RemoteAddr,
			"status", http.StatusBadRequest)
		http.Error(w, "Missing parameters", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		slog.Warn("Invalid submission ID",
			"method", r.Method,
			"path", r.URL.Path,
			"url", r.URL.String(),
			"hostname", r.Host,
			"user-agent", r.UserAgent(),
			"ip", r.RemoteAddr,
			"status", http.StatusBadRequest)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	// Delete the submission
	err = db.New(db.Pool).DeleteFormSubmission(ctx, id)
	if err != nil && err != pgx.ErrNoRows {
		slog.Error(fmt.Sprintf("Failed to delete submission (form_id: %s, submission_id: %d)", formID, id), tint.Err(err),
			"method", r.Method,
			"path", r.URL.Path,
			"url", r.URL.String(),
			"hostname", r.Host,
			"user-agent", r.UserAgent(),
			"ip", r.RemoteAddr,
			"status", http.StatusInternalServerError)
		http.Error(w, "Failed to delete", http.StatusInternalServerError)
		return
	}

	// Get the page parameter from query string (default to 1)
	page := 1
	if p := r.URL.Query().Get("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	offset := int64((page - 1) * conf.Config.Dashboard.Pagination.Limit)

	// Fetch updated submissions for re-render
	queries := db.New(db.Pool)
	submissions, err := queries.ListFormSubmissions(ctx, db.ListFormSubmissionsParams{
		FormID: formID,
		Limit:  int32(conf.Config.Dashboard.Pagination.Limit),
		Offset: int32(offset),
	})
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to fetch submissions after delete (form_id: %s)", formID), tint.Err(err),
			"method", r.Method,
			"path", r.URL.Path,
			"url", r.URL.String(),
			"hostname", r.Host,
			"user-agent", r.UserAgent(),
			"ip", r.RemoteAddr,
			"status", http.StatusInternalServerError)
		http.Error(w, "Failed to refresh submissions", http.StatusInternalServerError)
		return
	}

	// Set HTMX response headers
	w.Header().Set("HX-Retarget", "#submissions-table")
	w.Header().Set("HX-Reswap", "outerHTML")

	// Render and return the updated submissions table
	FormSubmissionsTableTempl(formID, submissions).Render(ctx, w)
}
