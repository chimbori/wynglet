package core

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

// ComputeKey computes a fully-qualified cache key for a URL.
// The cache key is filesystem-safe and includes the extension and optional thumbnail marker.
// Format: domain/path_path_path_file.HASH8.ext (or domain/path_path_path_file.t.HASH8.ext for thumbnails)
//
// Examples:
//
//	ComputeKey("https://example.com/foo/bar.html", "png", false) -> "example.com/foo_bar.abc12345.png"
//	ComputeKey("https://example.com/foo/bar.html", "webp", true) -> "example.com/foo_bar.t.abc12345.webp"
func ComputeKey(urlString string, ext string, isThumbnail bool) string {
	parsed, err := url.Parse(urlString)
	if err != nil {
		// Fallback: use hash-based key for invalid URLs
		hash := SHA256(urlString)
		return fmt.Sprintf("invalid/%s.%s", hash, ext)
	}

	domain := parsed.Hostname()
	if domain == "" {
		// Fallback: use hash-based key if no hostname
		hash := SHA256(urlString)
		return fmt.Sprintf("nohostname/%s.%s", hash, ext)
	}

	// Extract and sanitize the path component
	pathComponent := parsed.Path
	if pathComponent == "" || pathComponent == "/" {
		pathComponent = "index"
	} else {
		// Remove leading/trailing slashes
		pathComponent = strings.Trim(pathComponent, "/")
	}

	// Replace invalid filename characters with underscore
	pathComponent = sanitizePathComponent(pathComponent)

	// Calculate hash of domain+path (protocol-independent, first 8 chars)
	// This ensures https://example.com/foo and http://example.com/foo generate the same hash
	normalizedForHash := domain + parsed.Path
	if parsed.RawQuery != "" {
		normalizedForHash += "?" + parsed.RawQuery
	}
	fullHash := SHA256(normalizedForHash)
	hashPrefix := fullHash[:8]

	// Build filename
	var filename string
	if isThumbnail {
		filename = fmt.Sprintf("%s.t.%s.%s", pathComponent, hashPrefix, ext)
	} else {
		filename = fmt.Sprintf("%s.%s.%s", pathComponent, hashPrefix, ext)
	}

	return fmt.Sprintf("%s/%s", domain, filename)
}

// sanitizePathComponent replaces invalid filename characters with underscores.
func sanitizePathComponent(path string) string {
	// Replace invalid filename characters with underscore.
	// On most filesystems, these characters are problematic: / \ : * ? " < > |
	// We also replace spaces and other special characters for better readability.
	invalidChars := regexp.MustCompile(`[\s/\\:*?"<>|]`)
	return invalidChars.ReplaceAllString(path, "_")
}
