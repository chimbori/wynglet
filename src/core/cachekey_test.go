package core

import (
	"strings"
	"testing"
)

func TestComputeKey_Basic(t *testing.T) {
	key := ComputeKey("https://example.com/foo/bar.html", "png", false)

	if !strings.HasPrefix(key, "example.com/") {
		t.Errorf("Expected key to start with 'example.com/', got '%s'", key)
	}

	if !strings.Contains(key, "foo_bar") {
		t.Errorf("Expected key to contain 'foo_bar', got '%s'", key)
	}

	if !strings.HasSuffix(key, ".png") {
		t.Errorf("Expected key to end with .png, got '%s'", key)
	}
}

func TestComputeKey_RootPath(t *testing.T) {
	key := ComputeKey("https://example.com/", "png", false)

	if !strings.Contains(key, "index") {
		t.Errorf("Expected 'index' in key for root path, got '%s'", key)
	}
}

func TestComputeKey_WithExtension(t *testing.T) {
	key := ComputeKey("https://example.com/foo/bar.html", "webp", false)

	if !strings.HasSuffix(key, ".webp") {
		t.Errorf("Expected key to end with .webp, got '%s'", key)
	}
}

func TestComputeKey_Thumbnail(t *testing.T) {
	key := ComputeKey("https://example.com/foo/bar.html", "webp", true)

	if !strings.Contains(key, ".t.") {
		t.Errorf("Expected thumbnail marker '.t.' in key, got '%s'", key)
	}

	if !strings.HasSuffix(key, ".webp") {
		t.Errorf("Expected key to end with .webp, got '%s'", key)
	}
}

func TestComputeKey_SanitizeSpecialChars(t *testing.T) {
	key := ComputeKey("https://example.com/path?query=value&foo=bar", "png", false)

	// Special characters should be replaced with underscores
	if strings.Contains(key, "?") || strings.Contains(key, "&") {
		t.Errorf("Special characters not sanitized in key: %s", key)
	}
}

func TestComputeKey_InvalidURL(t *testing.T) {
	key := ComputeKey("not a valid url", "png", false)

	// Should fallback gracefully
	if !strings.HasPrefix(key, "nohostname/") {
		t.Errorf("Expected fallback key to start with 'nohostname/', got '%s'", key)
	}
}

func TestComputeKey_Consistency(t *testing.T) {
	url := "https://example.com/foo/bar.html"
	key1 := ComputeKey(url, "png", false)
	key2 := ComputeKey(url, "png", false)

	if key1 != key2 {
		t.Errorf("ComputeKey not consistent: %s vs %s", key1, key2)
	}
}

func TestComputeKey_DifferentExtensions(t *testing.T) {
	url := "https://example.com/foo/bar.html"
	key1 := ComputeKey(url, "png", false)
	key2 := ComputeKey(url, "webp", false)

	if key1 == key2 {
		t.Errorf("Different extensions should produce different keys")
	}

	if !strings.HasSuffix(key1, ".png") || !strings.HasSuffix(key2, ".webp") {
		t.Errorf("Keys don't have correct extensions")
	}
}

func TestComputeKey_DifferentThumbnailFlags(t *testing.T) {
	url := "https://example.com/foo/bar.html"
	key1 := ComputeKey(url, "webp", false)
	key2 := ComputeKey(url, "webp", true)

	if key1 == key2 {
		t.Errorf("Thumbnail flag should produce different keys")
	}

	if !strings.Contains(key2, ".t.") {
		t.Errorf("Thumbnail key should contain '.t.' marker")
	}
}

func TestSanitizePathComponent(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"foo/bar", "foo_bar"},
		{"foo bar", "foo_bar"},
		{"foo:bar", "foo_bar"},
		{"foo*bar", "foo_bar"},
		{"foo?bar", "foo_bar"},
		{"foo\"bar", "foo_bar"},
		{"foo<bar", "foo_bar"},
		{"foo>bar", "foo_bar"},
		{"foo|bar", "foo_bar"},
		{"foo\\bar", "foo_bar"},
		{"normal.html", "normal.html"},
	}

	for _, test := range tests {
		result := sanitizePathComponent(test.input)
		if result != test.expected {
			t.Errorf("sanitizePathComponent('%s'): expected '%s', got '%s'", test.input, test.expected, result)
		}
	}
}

func TestComputeKey_DifferentURLs(t *testing.T) {
	key1 := ComputeKey("https://example.com/foo", "png", false)
	key2 := ComputeKey("https://example.com/bar", "png", false)

	if key1 == key2 {
		t.Errorf("Different URLs should produce different keys")
	}
}

func TestComputeKey_ProtocolIndependent(t *testing.T) {
	// Same domain and path, different protocols should produce identical cache keys
	key1 := ComputeKey("https://example.com/foo/bar.html", "png", false)
	key2 := ComputeKey("http://example.com/foo/bar.html", "png", false)

	if key1 != key2 {
		t.Errorf("Different protocols should produce same cache key: %s vs %s", key1, key2)
	}
}

func TestComputeKey_ProtocolIndependent_WithQuery(t *testing.T) {
	// With query string
	key1 := ComputeKey("https://example.com/page?id=123", "png", false)
	key2 := ComputeKey("http://example.com/page?id=123", "png", false)

	if key1 != key2 {
		t.Errorf("Different protocols with query should produce same cache key: %s vs %s", key1, key2)
	}
}

func TestComputeKey_DifferentQueryStrings(t *testing.T) {
	// Different query strings should produce different cache keys
	key1 := ComputeKey("https://example.com/page?id=123", "png", false)
	key2 := ComputeKey("https://example.com/page?id=456", "png", false)

	if key1 == key2 {
		t.Errorf("Different query strings should produce different cache keys")
	}
}
