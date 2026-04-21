package core

import (
	"log/slog"
	"net/http"
	"sync"
	"time"
)

// debugMode tracks which domains have debug mode enabled and when.
// Debug mode is not persistent and expires after 1 hour.
type debugModeTracker struct {
	mu    sync.RWMutex
	modes map[string]time.Time // domain -> activation time
}

var debugTracker = &debugModeTracker{
	modes: make(map[string]time.Time),
}

const debugModeDuration = 1 * time.Hour

// IsDebugModeActive checks if debug mode is active for the given domain.
// Returns true if debug mode is enabled and has not expired.
func IsDebugModeActive(domain string) bool {
	debugTracker.mu.RLock()
	defer debugTracker.mu.RUnlock()

	activatedAt, exists := debugTracker.modes[domain]
	if !exists {
		return false
	}

	if time.Since(activatedAt) > debugModeDuration {
		return false
	}

	return true
}

// GetDebugModeRemainingTime returns the remaining time for debug mode on the given domain.
// Returns 0 if debug mode is not active.
func GetDebugModeRemainingTime(domain string) time.Duration {
	debugTracker.mu.RLock()
	defer debugTracker.mu.RUnlock()

	activatedAt, exists := debugTracker.modes[domain]
	if !exists {
		return 0
	}

	elapsed := time.Since(activatedAt)
	if elapsed > debugModeDuration {
		return 0
	}

	return debugModeDuration - elapsed
}

// ToggleDebugMode enables or disables debug mode for the given domain.
// When enabled, it remains active for 1 hour from the current time.
// - domain: hostname whose debug mode state should be changed.
// - enable: whether debug mode should be turned on (true) or off (false).
// - req (optional): the triggering HTTP request used for logging context; it may be nil.
func ToggleDebugMode(domain string, enable bool, req *http.Request) {
	var msg string
	if enable {
		msg = "Debug Mode enabled"
	} else {
		msg = "Debug Mode disabled"
	}

	if req != nil {
		slog.Warn(msg,
			"method", req.Method,
			"path", req.URL.Path,
			"url", req.URL.String(),
			"status", http.StatusOK,
			"ip", ReadUserIP(req),
			"hostname", domain)
	} else {
		slog.Warn(msg, "hostname", domain)
	}

	debugTracker.mu.Lock()
	defer debugTracker.mu.Unlock()

	if enable {
		debugTracker.modes[domain] = time.Now()
	} else {
		delete(debugTracker.modes, domain)
	}
}

// CleanupExpiredDebugModes removes expired debug mode entries.
// This is called periodically to prevent the map from growing indefinitely.
func CleanupExpiredDebugModes() {
	debugTracker.mu.Lock()
	defer debugTracker.mu.Unlock()

	now := time.Now()
	for domain, activatedAt := range debugTracker.modes {
		if now.Sub(activatedAt) > debugModeDuration {
			slog.Warn("Debug Mode expired",
				"hostname", domain,
				"activated_at", activatedAt)
			delete(debugTracker.modes, domain)
		}
	}
}
