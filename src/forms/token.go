// Package forms provides HTTP API handlers for form submission and CSRF token generation.
// This file implements cryptographically secure CSRF token generation with replay prevention.
package forms

import (
	"crypto/rand"
	"encoding/hex"
	"net"
	"net/http"
	"sync"
	"time"
)

const (
	tokenLength = 32
	tokenTTL    = 5 * time.Minute
)

// TokenResponse is the JSON response body for the token endpoint.
type TokenResponse struct {
	Token     string `json:"token"`
	ExpiresAt string `json:"expires_at"`
}

// tokenMetadata stores metadata for a form submission token.
type tokenMetadata struct {
	FormID    string
	CreatedAt time.Time
	Used      bool
}

// TokenCache manages in-memory token storage with automatic expiration.
type TokenCache struct {
	mu     sync.RWMutex
	tokens map[string]tokenMetadata
}

// NewTokenCache creates a new TokenCache and starts a background goroutine
// to periodically remove expired tokens.
func NewTokenCache() *TokenCache {
	tc := &TokenCache{
		tokens: make(map[string]tokenMetadata),
	}

	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			tc.removeExpired()
		}
	}()

	return tc
}

// Set stores a token with its associated metadata.
func (tc *TokenCache) Set(token string, meta tokenMetadata) {
	tc.mu.Lock()
	defer tc.mu.Unlock()
	tc.tokens[token] = meta
}

// Get retrieves a token’s metadata if it exists.
func (tc *TokenCache) Get(token string) (tokenMetadata, bool) {
	tc.mu.RLock()
	defer tc.mu.RUnlock()
	meta, found := tc.tokens[token]
	return meta, found
}

// MarkUsed marks a token as used to prevent replay attacks.
func (tc *TokenCache) MarkUsed(token string) {
	tc.mu.Lock()
	defer tc.mu.Unlock()
	if meta, found := tc.tokens[token]; found {
		meta.Used = true
		tc.tokens[token] = meta
	}
}

// removeExpired deletes all tokens that have exceeded the token TTL.
func (tc *TokenCache) removeExpired() {
	tc.mu.Lock()
	defer tc.mu.Unlock()
	now := time.Now()
	for token, meta := range tc.tokens {
		if now.Sub(meta.CreatedAt) > tokenTTL {
			delete(tc.tokens, token)
		}
	}
}

// generateToken creates a cryptographically secure random token.
func generateToken() (string, error) {
	bytes := make([]byte, tokenLength)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// getClientIP extracts the client’s IP address from the request.
func getClientIP(r *http.Request) string {
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	return ip
}

var tokenCache = NewTokenCache()
