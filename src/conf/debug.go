package conf

import (
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

// EnableDebugMode activates debug mode for the given domain.
// It will remain active for 1 hour from the current time.
func EnableDebugMode(domain string) {
	debugTracker.mu.Lock()
	defer debugTracker.mu.Unlock()

	debugTracker.modes[domain] = time.Now()
}

// DisableDebugMode deactivates debug mode for the given domain.
func DisableDebugMode(domain string) {
	debugTracker.mu.Lock()
	defer debugTracker.mu.Unlock()

	delete(debugTracker.modes, domain)
}

// CleanupExpiredDebugModes removes expired debug mode entries.
// This is called periodically to prevent the map from growing indefinitely.
func CleanupExpiredDebugModes() {
	debugTracker.mu.Lock()
	defer debugTracker.mu.Unlock()

	now := time.Now()
	for domain, activatedAt := range debugTracker.modes {
		if now.Sub(activatedAt) > debugModeDuration {
			delete(debugTracker.modes, domain)
		}
	}
}
