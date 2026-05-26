package core

import (
	"net"
	"net/http"
	"strings"
)

// GetOutboundIP returns the preferred outbound IP address of this machine.
// It establishes a UDP connection to 8.8.8.8 to determine the local IP.
func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return nil
	}
	defer conn.Close()
	return conn.LocalAddr().(*net.UDPAddr).IP
}

// ReadUserIP extracts the client IP address from an HTTP request.
// It checks X-Real-Ip, X-Forwarded-For headers, and falls back to RemoteAddr.
func ReadUserIP(r *http.Request) string {
	IPAddress := r.Header.Get("X-Real-Ip")
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Forwarded-For")
	}
	if IPAddress == "" {
		IPAddress = r.RemoteAddr
	}
	return IPAddress
}

// SecurityHeaders adds standard HTTP security headers to all responses.
func SecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Referrer-Policy", "no-referrer")
		w.Header().Set("Permissions-Policy", "geolocation=(), microphone=(), camera=()")
		w.Header().Set("Content-Security-Policy", strings.Join([]string{
			"default-src 'self'",
			"img-src 'self' data:",
			"script-src 'self'",
			"style-src 'self' 'unsafe-inline'",
			"connect-src 'self' https:",
			"object-src 'none'",
			"base-uri 'self'",
			"frame-ancestors 'none'",
		}, "; "))

		next.ServeHTTP(w, r)
	})
}

// NormalizeClientIP normalizes a client IP address string.
// It trims whitespace, extracts the first IP from X-Forwarded-For headers,
// and removes the port from the host if present.
func NormalizeClientIP(raw string) string {
	first := strings.TrimSpace(raw)
	if strings.Contains(first, ",") {
		first, _, _ = strings.Cut(first, ",")
		first = strings.TrimSpace(first)
	}
	if host, _, err := net.SplitHostPort(first); err == nil {
		first = host
	}
	return first
}

// SetCORSHeaders sets CORS headers for an authorized origin.
func SetCORSHeaders(w http.ResponseWriter, origin string) {
	w.Header().Set("Access-Control-Allow-Origin", origin)
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Max-Age", "3600")
}
