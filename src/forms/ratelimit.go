// Package forms provides HTTP API handlers for form submission and CSRF token generation.
// This file implements per-form-per-IP rate limiting with configurable request windows.
package forms

import (
	"fmt"
	"sync"
	"time"

	"wynglet.chimbori.dev/conf"
)

// RateLimiter tracks form submission requests per IP to enforce rate limits.
type RateLimiter struct {
	mu    sync.RWMutex
	stats map[string][]time.Time
}

// NewRateLimiter creates a new RateLimiter and starts a background goroutine
// to periodically clean up stale rate limit records.
func NewRateLimiter() *RateLimiter {
	rl := &RateLimiter{
		stats: make(map[string][]time.Time),
	}

	go func() {
		ticker := time.NewTicker(10 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			rl.cleanup()
		}
	}()

	return rl
}

// Check returns true if the request from the given IP on the given form
// is within the rate limit, false otherwise. It records the current request time.
func (rl *RateLimiter) Check(formID, ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	key := fmt.Sprintf("%s:%s", formID, ip)
	now := time.Now()
	oneHourAgo := now.Add(-1 * time.Hour)

	if timestamps, found := rl.stats[key]; found {
		valid := []time.Time{}
		for _, ts := range timestamps {
			if ts.After(oneHourAgo) {
				valid = append(valid, ts)
			}
		}
		rl.stats[key] = valid
	}

	if len(rl.stats[key]) >= conf.Config.Forms.RateLimit.PerIPHour {
		return false
	}

	rl.stats[key] = append(rl.stats[key], now)
	return true
}

// cleanup removes rate limit records older than one hour.
func (rl *RateLimiter) cleanup() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	oneHourAgo := time.Now().Add(-1 * time.Hour)
	for key, timestamps := range rl.stats {
		valid := []time.Time{}
		for _, ts := range timestamps {
			if ts.After(oneHourAgo) {
				valid = append(valid, ts)
			}
		}
		if len(valid) == 0 {
			delete(rl.stats, key)
		} else {
			rl.stats[key] = valid
		}
	}
}

// rateLimiter is the global instance used for rate limiting form submissions.
var rateLimiter = NewRateLimiter()
