package core

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const sitemapMaxBytes int64 = 50 * 1024 * 1024

var sitemapHTTPClient = &http.Client{Timeout: 30 * time.Second}

type xmlURLSet struct {
	URLs []xmlURL `xml:"url"`
}

type xmlURL struct {
	Loc string `xml:"loc"`
}

type xmlSitemapIndex struct {
	Sitemaps []xmlSitemap `xml:"sitemap"`
}

type xmlSitemap struct {
	Loc string `xml:"loc"`
}

// FetchSitemap downloads a sitemap or sitemap index, returns deduplicated HTTPS URLs,
// and limits the result set to maxURLs.
func FetchSitemap(ctx context.Context, sitemapURL string, maxURLs int) ([]string, error) {
	return fetchSitemapWithClient(ctx, sitemapHTTPClient, sitemapURL, maxURLs)
}

func fetchSitemapWithClient(ctx context.Context, client *http.Client, sitemapURL string, maxURLs int) ([]string, error) {
	if maxURLs <= 0 {
		return nil, fmt.Errorf("maxURLs must be greater than zero")
	}

	seen := make(map[string]struct{}, maxURLs)
	urls, err := fetchSitemapRecursive(ctx, client, sitemapURL, maxURLs, false, seen)
	if err != nil {
		return nil, err
	}
	return urls, nil
}

func fetchSitemapRecursive(ctx context.Context, client *http.Client, sitemapURL string, maxURLs int, nested bool, seen map[string]struct{}) ([]string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, sitemapURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; Butterfly/1.0; +https://butterfly.chimbori.dev)")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, sitemapMaxBytes))
	if err != nil {
		return nil, err
	}

	var index xmlSitemapIndex
	if err := xml.Unmarshal(body, &index); err == nil && len(index.Sitemaps) > 0 {
		if nested {
			return nil, fmt.Errorf("nested sitemap index recursion is not supported")
		}

		var urls []string
		for _, sitemap := range index.Sitemaps {
			childURL := strings.TrimSpace(sitemap.Loc)
			if childURL == "" {
				continue
			}

			childURLs, err := fetchSitemapRecursive(ctx, client, childURL, maxURLs-len(seen), true, seen)
			if err != nil {
				return nil, err
			}
			urls = append(urls, childURLs...)
			if len(seen) >= maxURLs {
				break
			}
		}
		return urls, nil
	}

	var urlset xmlURLSet
	if err := xml.Unmarshal(body, &urlset); err != nil {
		return nil, err
	}

	urls := make([]string, 0, len(urlset.URLs))
	for _, item := range urlset.URLs {
		if len(seen) >= maxURLs {
			break
		}

		rawURL := strings.TrimSpace(item.Loc)
		if rawURL == "" {
			continue
		}

		parsed, err := url.Parse(rawURL)
		if err != nil || !strings.EqualFold(parsed.Scheme, "https") {
			continue
		}

		normalized := parsed.String()
		if _, ok := seen[normalized]; ok {
			continue
		}
		seen[normalized] = struct{}{}
		urls = append(urls, normalized)
	}

	return urls, nil
}
