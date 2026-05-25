// Package forms provides HTTP API handlers for form submission and CSRF token generation.
// Handlers enforce CORS validation, rate limiting, and replay prevention.
package forms

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/lmittmann/tint"
	"gopkg.in/yaml.v3"
	"wynglet.chimbori.dev/core"
	"wynglet.chimbori.dev/db"
	"wynglet.chimbori.dev/validation"
)

type FormSubmission struct {
	Token    string            `form:"_token"`
	FormID   string            `form:"_form_id"`
	Subject  string            `form:"_subject"`
	Redirect string            `form:"_redirect"`
	Honeypot string            `form:"_honeypot"`
	Fields   map[string]string `form:"-"`
}

func Init(mux *http.ServeMux) {
	mux.HandleFunc("GET /forms/v1/token", getTokenHandler)
	mux.HandleFunc("POST /forms/v1/submit", formSubmitHandler)
}

// Validates CORS origin, CSRF token, form data, and rate limits.
// Records form submissions to the database with spam detection via honeypot fields.
// POST /forms/v1/submit
func formSubmitHandler(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")
	clientIP := getClientIP(r)

	// Validate CORS origin early, but only if present.
	// Missing Origin header indicates a same-origin request (browsers only send Origin for CORS requests),
	// so we allow it through. However, origin="null" (file:// URLs) is rejected for security.
	if origin == "null" {
		slog.Warn(fmt.Sprintf("Form submission rejected: origin from file:// URL: origin=`%s`", origin),
			"method", r.Method,
			"path", r.URL.Path,
			"url", r.URL.String(),
			"hostname", r.Host,
			"user-agent", r.UserAgent(),
			"ip", clientIP,
			"status", http.StatusForbidden)
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	// If Origin header is present (not empty and not "null"), validate it.
	if origin != "" {
		parsedURL, err := url.Parse(origin)
		if err != nil || parsedURL.Host == "" {
			slog.Warn(fmt.Sprintf("Failed to parse origin: origin=`%s`", origin),
				"method", r.Method,
				"path", r.URL.Path,
				"url", r.URL.String(),
				"hostname", r.Host,
				"user-agent", r.UserAgent(),
				"ip", clientIP,
				"status", http.StatusForbidden,
				tint.Err(err))
			http.Error(w, "Unauthorized", http.StatusForbidden)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		queries := db.New(db.Pool)
		authorized, authErr := validation.IsAuthorized(ctx, queries, parsedURL)
		if authErr != nil {
			slog.Warn(fmt.Sprintf("Failed to check domain authorization: origin=`%s`", origin),
				"method", r.Method,
				"path", r.URL.Path,
				"url", r.URL.String(),
				"hostname", r.Host,
				"user-agent", r.UserAgent(),
				"ip", clientIP,
				"status", http.StatusForbidden,
				tint.Err(authErr))
			http.Error(w, "Unauthorized", http.StatusForbidden)
			return
		}

		if !authorized {
			slog.Warn(fmt.Sprintf("Origin not authorized: origin=%s", origin),
				"method", r.Method,
				"path", r.URL.Path,
				"url", r.URL.String(),
				"hostname", r.Host,
				"user-agent", r.UserAgent(),
				"ip", clientIP,
				"status", http.StatusForbidden)
			http.Error(w, "Unauthorized", http.StatusForbidden)
			return
		}

		core.SetCORSHeaders(w, origin)
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		slog.Error("Failed to parse form data", tint.Err(err),
			"method", r.Method,
			"path", r.URL.Path,
			"url", r.URL.String(),
			"hostname", r.Host,
			"user-agent", r.UserAgent(),
			"ip", clientIP,
			"status", http.StatusBadRequest)
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	token := r.FormValue("_token")
	formID := r.FormValue("_form_id")
	honeypotField := r.FormValue("_honeypot")
	subject := r.FormValue("_subject")
	redirectURL := r.FormValue("_redirect")

	if token == "" || formID == "" {
		slog.Warn(fmt.Sprintf("Form submission missing required fields: has_token=%v, has_form_id=%v", token != "", formID != ""),
			"method", r.Method,
			"path", r.URL.Path,
			"url", r.URL.String(),
			"hostname", r.Host,
			"user-agent", r.UserAgent(),
			"ip", clientIP,
			"status", http.StatusBadRequest)
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	meta, found := tokenCache.Get(token)
	if !found {
		slog.Warn(fmt.Sprintf("Token not found: token=%s", token[:8]+"..."),
			"method", r.Method,
			"path", r.URL.Path,
			"url", r.URL.String(),
			"hostname", r.Host,
			"user-agent", r.UserAgent(),
			"ip", clientIP,
			"status", http.StatusForbidden)
		http.Error(w, "Invalid or expired token", http.StatusForbidden)
		return
	}

	if meta.Used {
		slog.Warn(fmt.Sprintf("Token already used (replay attack): form_id=%s", formID),
			"method", r.Method,
			"path", r.URL.Path,
			"url", r.URL.String(),
			"hostname", r.Host,
			"user-agent", r.UserAgent(),
			"ip", clientIP,
			"status", http.StatusForbidden)
		http.Error(w, "Token already used", http.StatusForbidden)
		return
	}

	if meta.FormID != formID {
		slog.Warn(fmt.Sprintf("Token form_id mismatch: expected=%s, got=%s, form_id=%s", meta.FormID, formID, formID),
			"method", r.Method,
			"path", r.URL.Path,
			"url", r.URL.String(),
			"hostname", r.Host,
			"user-agent", r.UserAgent(),
			"ip", clientIP,
			"status", http.StatusForbidden)
		http.Error(w, "Invalid token for this form", http.StatusForbidden)
		return
	}

	tokenCache.MarkUsed(token)

	if !rateLimiter.Check(formID, clientIP) {
		slog.Warn(fmt.Sprintf("Rate limit exceeded: form_id=%s", formID),
			"method", r.Method,
			"path", r.URL.Path,
			"url", r.URL.String(),
			"hostname", r.Host,
			"user-agent", r.UserAgent(),
			"ip", clientIP,
			"status", http.StatusTooManyRequests)
		http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
		return
	}

	isSpam := false
	if honeypotField != "" {
		honeypotValue := r.FormValue(honeypotField)
		if honeypotValue != "" {
			isSpam = true
			slog.Info(fmt.Sprintf("Honeypot triggered: honeypot_field=%s, form_id=%s", honeypotField, formID),
				"method", r.Method,
				"path", r.URL.Path,
				"url", r.URL.String(),
				"hostname", r.Host,
				"user-agent", r.UserAgent(),
				"ip", clientIP)
		}
	}

	userFields := make(map[string]string)
	for key, values := range r.Form {
		if len(key) > 0 && key[0] != '_' && len(values) > 0 && len(strings.TrimSpace(values[0])) > 0 {
			userFields[key] = values[0]
		}
	}

	fieldsYAML, err := yaml.Marshal(userFields)
	if err != nil {
		slog.Error("Failed to marshal form fields", tint.Err(err),
			"method", r.Method,
			"path", r.URL.Path,
			"url", r.URL.String(),
			"hostname", r.Host,
			"user-agent", r.UserAgent(),
			"ip", clientIP,
			"status", http.StatusInternalServerError)
		http.Error(w, "Failed to process submission", http.StatusInternalServerError)
		return
	}

	_, err = db.New(db.Pool).CreateFormSubmission(r.Context(), db.CreateFormSubmissionParams{
		FormID:      formID,
		Domain:      origin,
		IpAddress:   clientIP,
		FormData:    string(fieldsYAML),
		IsSpam:      isSpam,
		EmailSentAt: nil,
	})
	if err != nil {
		slog.Error("Failed to store form submission", tint.Err(err),
			"method", r.Method,
			"path", r.URL.Path,
			"url", r.URL.String(),
			"hostname", r.Host,
			"user-agent", r.UserAgent(),
			"ip", clientIP,
			"status", http.StatusInternalServerError)
		http.Error(w, "Failed to process submission", http.StatusInternalServerError)
		return
	}

	slog.Info(fmt.Sprintf("Form submitted: spam=%v, subject=%s, form_id=%s", isSpam, subject, formID),
		"method", r.Method,
		"path", r.URL.Path,
		"url", r.URL.String(),
		"hostname", r.Host,
		"user-agent", r.UserAgent(),
		"ip", clientIP)

	if r.Header.Get("Accept") == "application/json" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	} else {
		if redirectURL != "" {
			http.Redirect(w, r, redirectURL, http.StatusSeeOther)
		} else {
			w.Header().Set("Content-Type", "text/plain")
			fmt.Fprintf(w, "Thank you for your submission.\n")
		}
	}
}

// Returns a new CSRF token for use in form submissions.
// GET /forms/v1/token
func getTokenHandler(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")
	clientIP := getClientIP(r)

	// Validate CORS origin early, but only if present.
	// Missing Origin header indicates a same-origin request (browsers only send Origin for CORS requests),
	// so we allow it through. However, origin="null" (file:// URLs) is rejected for security.
	if origin == "null" {
		slog.Warn(fmt.Sprintf("Token request rejected: origin from file:// URL: origin=`%s`", origin),
			"method", r.Method,
			"path", r.URL.Path,
			"url", r.URL.String(),
			"hostname", r.Host,
			"user-agent", r.UserAgent(),
			"ip", clientIP,
			"status", http.StatusForbidden)
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	// If Origin header is present (not empty and not "null"), validate it
	if origin != "" {
		parsedURL, err := url.Parse(origin)
		if err != nil || parsedURL.Host == "" {
			slog.Warn(fmt.Sprintf("Failed to parse origin: origin=`%s`", origin),
				"method", r.Method,
				"path", r.URL.Path,
				"url", r.URL.String(),
				"hostname", r.Host,
				"user-agent", r.UserAgent(),
				"ip", clientIP,
				"status", http.StatusForbidden,
				tint.Err(err))
			http.Error(w, "Unauthorized", http.StatusForbidden)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		queries := db.New(db.Pool)
		authorized, authErr := validation.IsAuthorized(ctx, queries, parsedURL)
		if authErr != nil {
			slog.Warn(fmt.Sprintf("Failed to check domain authorization: origin=`%s`", origin),
				"method", r.Method,
				"path", r.URL.Path,
				"url", r.URL.String(),
				"hostname", r.Host,
				"user-agent", r.UserAgent(),
				"ip", clientIP,
				"status", http.StatusForbidden,
				tint.Err(authErr))
			http.Error(w, "Unauthorized", http.StatusForbidden)
			return
		}

		if !authorized {
			slog.Warn(fmt.Sprintf("Origin not authorized: origin=`%s`", origin),
				"method", r.Method,
				"path", r.URL.Path,
				"url", r.URL.String(),
				"hostname", r.Host,
				"user-agent", r.UserAgent(),
				"ip", clientIP,
				"status", http.StatusForbidden)
			http.Error(w, "Unauthorized", http.StatusForbidden)
			return
		}

		core.SetCORSHeaders(w, origin)
	}

	formID := r.URL.Query().Get("form_id")

	if formID == "" {
		slog.Error("Missing form_id parameter",
			"method", r.Method,
			"path", r.URL.Path,
			"url", r.URL.String(),
			"hostname", r.Host,
			"user-agent", r.UserAgent(),
			"ip", clientIP,
			"status", http.StatusBadRequest)
		http.Error(w, "Missing form_id parameter", http.StatusBadRequest)
		return
	}

	if !rateLimiter.Check(formID+":token", clientIP) {
		slog.Error(fmt.Sprintf("Rate limit exceeded: form_id=%s", formID),
			"method", r.Method,
			"path", r.URL.Path,
			"url", r.URL.String(),
			"hostname", r.Host,
			"user-agent", r.UserAgent(),
			"ip", clientIP,
			"status", http.StatusTooManyRequests)
		http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
		return
	}

	token, err := generateToken()
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to generate token: form_id=%s", formID), tint.Err(err),
			"method", r.Method,
			"path", r.URL.Path,
			"url", r.URL.String(),
			"hostname", r.Host,
			"user-agent", r.UserAgent(),
			"ip", clientIP,
			"status", http.StatusInternalServerError)
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	tokenCache.Set(token, tokenMetadata{
		FormID:    formID,
		CreatedAt: time.Now(),
		Used:      false,
	})

	expiresAt := time.Now().Add(tokenTTL)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(TokenResponse{
		Token:     token,
		ExpiresAt: expiresAt.Format(time.RFC3339),
	})
}
