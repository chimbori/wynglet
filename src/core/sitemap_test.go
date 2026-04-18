package core

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestFetchSitemap(t *testing.T) {
	t.Parallel()

	childOne := `<?xml version="1.0" encoding="UTF-8"?><urlset><url><loc>https://example.com/one</loc></url><url><loc>https://example.com/two</loc></url></urlset>`
	childTwo := `<?xml version="1.0" encoding="UTF-8"?><urlset><url><loc>https://example.com/three</loc></url></urlset>`

	var server *httptest.Server
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/simple.xml":
			fmt.Fprint(w, `<?xml version="1.0" encoding="UTF-8"?><urlset><url><loc>https://example.com/a</loc></url><url><loc>https://example.com/b</loc></url><url><loc>https://example.com/c</loc></url></urlset>`)
		case "/duplicates.xml":
			fmt.Fprint(w, `<?xml version="1.0" encoding="UTF-8"?><urlset><url><loc>https://example.com/a</loc></url><url><loc>https://example.com/a</loc></url><url><loc>https://example.com/b</loc></url></urlset>`)
		case "/capped.xml":
			var builder strings.Builder
			builder.WriteString(`<?xml version="1.0" encoding="UTF-8"?><urlset>`)
			for index := range 20 {
				builder.WriteString(fmt.Sprintf(`<url><loc>https://example.com/%d</loc></url>`, index))
			}
			builder.WriteString(`</urlset>`)
			fmt.Fprint(w, builder.String())
		case "/empty.xml":
			fmt.Fprint(w, `<?xml version="1.0" encoding="UTF-8"?><urlset></urlset>`)
		case "/invalid.xml":
			fmt.Fprint(w, `<garbage>`)
		case "/mixed.xml":
			fmt.Fprint(w, `<?xml version="1.0" encoding="UTF-8"?><urlset><url><loc>http://example.com/insecure</loc></url><url><loc>https://example.com/secure</loc></url></urlset>`)
		case "/namespace.xml":
			fmt.Fprint(w, `<?xml version="1.0" encoding="UTF-8"?><urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9"><url><loc>https://example.com/namespaced</loc></url></urlset>`)
		case "/child-one.xml":
			fmt.Fprint(w, childOne)
		case "/child-two.xml":
			fmt.Fprint(w, childTwo)
		case "/index.xml":
			fmt.Fprintf(w, `<?xml version="1.0" encoding="UTF-8"?><sitemapindex><sitemap><loc>%s/child-one.xml</loc></sitemap><sitemap><loc>%s/child-two.xml</loc></sitemap></sitemapindex>`, server.URL, server.URL)
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	testCases := []struct {
		name      string
		path      string
		maxURLs   int
		wantURLs  []string
		wantError string
	}{
		{name: "simple urlset", path: "/simple.xml", maxURLs: 10, wantURLs: []string{"https://example.com/a", "https://example.com/b", "https://example.com/c"}},
		{name: "sitemap index", path: "/index.xml", maxURLs: 10, wantURLs: []string{"https://example.com/one", "https://example.com/two", "https://example.com/three"}},
		{name: "deduplicates urls", path: "/duplicates.xml", maxURLs: 10, wantURLs: []string{"https://example.com/a", "https://example.com/b"}},
		{name: "caps results", path: "/capped.xml", maxURLs: 5, wantURLs: []string{"https://example.com/0", "https://example.com/1", "https://example.com/2", "https://example.com/3", "https://example.com/4"}},
		{name: "empty sitemap", path: "/empty.xml", maxURLs: 10, wantURLs: []string{}},
		{name: "malformed xml", path: "/invalid.xml", maxURLs: 10, wantError: "XML syntax error"},
		{name: "filters non https", path: "/mixed.xml", maxURLs: 10, wantURLs: []string{"https://example.com/secure"}},
		{name: "supports namespace", path: "/namespace.xml", maxURLs: 10, wantURLs: []string{"https://example.com/namespaced"}},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			urls, err := fetchSitemapWithClient(ctx, server.Client(), server.URL+testCase.path, testCase.maxURLs)
			if testCase.wantError != "" {
				if err == nil || !strings.Contains(err.Error(), testCase.wantError) {
					t.Fatalf("expected error containing %q, got %v", testCase.wantError, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("FetchSitemap returned error: %v", err)
			}

			if len(urls) != len(testCase.wantURLs) {
				t.Fatalf("len(urls) = %d, want %d; urls=%v", len(urls), len(testCase.wantURLs), urls)
			}

			for index, wantURL := range testCase.wantURLs {
				if urls[index] != wantURL {
					t.Fatalf("urls[%d] = %q, want %q", index, urls[index], wantURL)
				}
			}
		})
	}
}
