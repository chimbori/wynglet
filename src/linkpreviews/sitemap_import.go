package linkpreviews

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"sync/atomic"

	"butterfly.chimbori.dev/conf"
	"butterfly.chimbori.dev/core"
	"butterfly.chimbori.dev/db"
	"butterfly.chimbori.dev/embedfs"
	"butterfly.chimbori.dev/validation"
	"github.com/lmittmann/tint"
)

const (
	sitemapImportUserAgent = "sitemap-import"
	sitemapImportMaxErrors = 50
)

type SitemapImportStatus struct {
	SitemapURL string
	TotalURLs  int
	Completed  int
	Skipped    int
	Failed     int
	InProgress bool
	Errors     []string

	mu         sync.Mutex
	cancelFunc context.CancelFunc
}

var sitemapImportState struct {
	mu     sync.Mutex
	status *SitemapImportStatus
}

func GetSitemapImportStatus() *SitemapImportStatus {
	sitemapImportState.mu.Lock()
	defer sitemapImportState.mu.Unlock()
	return cloneSitemapImportStatus(sitemapImportState.status)
}

func StartSitemapImport(parent context.Context, sitemapURL string) error {
	if Cache == nil {
		return fmt.Errorf("link preview cache is disabled")
	}

	sitemapImportState.mu.Lock()
	if sitemapImportState.status != nil && sitemapImportState.status.InProgress {
		sitemapImportState.mu.Unlock()
		return fmt.Errorf("a sitemap import is already in progress")
	}
	sitemapImportState.mu.Unlock()

	queries := db.New(db.Pool)
	urls, err := core.FetchSitemap(parent, sitemapURL, conf.Config.LinkPreviews.Sitemap.MaxURLs)
	if err != nil {
		return err
	}

	status := &SitemapImportStatus{
		SitemapURL: sitemapURL,
		TotalURLs:  len(urls),
		InProgress: true,
	}

	ctx, cancel := context.WithCancel(context.Background())
	status.cancelFunc = cancel

	sitemapImportState.mu.Lock()
	sitemapImportState.status = status
	sitemapImportState.mu.Unlock()

	if len(urls) == 0 {
		finishSitemapImport(status)
		return nil
	}

	existingURLs, err := queries.GetExistingLinkPreviewURLs(parent, urls)
	if err != nil {
		appendSitemapImportError(status, err.Error())
		finishSitemapImport(status)
		cancel()
		return err
	}

	existing := make(map[string]struct{}, len(existingURLs))
	for _, existingURL := range existingURLs {
		existing[existingURL] = struct{}{}
	}

	concurrency := max(conf.Config.LinkPreviews.Sitemap.ConcurrentURLs, 1)

	go runSitemapImport(ctx, status, urls, existing, concurrency)
	return nil
}

func CancelSitemapImport() bool {
	sitemapImportState.mu.Lock()
	defer sitemapImportState.mu.Unlock()
	if sitemapImportState.status == nil || !sitemapImportState.status.InProgress || sitemapImportState.status.cancelFunc == nil {
		return false
	}
	sitemapImportState.status.cancelFunc()
	return true
}

func runSitemapImport(ctx context.Context, status *SitemapImportStatus, urls []string, existing map[string]struct{}, concurrency int) {
	defer finishSitemapImport(status)

	jobs := make(chan string)
	var waitGroup sync.WaitGroup
	var skipped int64
	var completed int64
	var failed int64

	for range concurrency {
		waitGroup.Go(func() {
			queries := db.New(db.Pool)
			for url := range jobs {
				select {
				case <-ctx.Done():
					return
				default:
				}

				if _, ok := existing[url]; ok {
					atomic.AddInt64(&skipped, 1)
					updateSitemapImportCounts(status, int(atomic.LoadInt64(&completed)), int(atomic.LoadInt64(&skipped)), int(atomic.LoadInt64(&failed)))
					continue
				}

				validatedURL, _, err := validation.ValidateUrl(ctx, queries, url)
				if err != nil {
					atomic.AddInt64(&failed, 1)
					appendSitemapImportError(status, fmt.Sprintf("%s: %v", url, err))
					updateSitemapImportCounts(status, int(atomic.LoadInt64(&completed)), int(atomic.LoadInt64(&skipped)), int(atomic.LoadInt64(&failed)))
					continue
				}

				cached, err := Cache.Find(validatedURL)
				if err != nil {
					atomic.AddInt64(&failed, 1)
					appendSitemapImportError(status, fmt.Sprintf("%s: %v", validatedURL, err))
					updateSitemapImportCounts(status, int(atomic.LoadInt64(&completed)), int(atomic.LoadInt64(&skipped)), int(atomic.LoadInt64(&failed)))
					continue
				}

				if cached != nil {
					recordLinkPreviewCreated(validatedURL, sitemapImportUserAgent)
					atomic.AddInt64(&completed, 1)
					updateSitemapImportCounts(status, int(atomic.LoadInt64(&completed)), int(atomic.LoadInt64(&skipped)), int(atomic.LoadInt64(&failed)))
					continue
				}

				if err := generateAndCacheSitemapPreview(ctx, validatedURL); err != nil {
					atomic.AddInt64(&failed, 1)
					appendSitemapImportError(status, fmt.Sprintf("%s: %v", validatedURL, err))
					updateSitemapImportCounts(status, int(atomic.LoadInt64(&completed)), int(atomic.LoadInt64(&skipped)), int(atomic.LoadInt64(&failed)))
					continue
				}

				recordLinkPreviewCreated(validatedURL, sitemapImportUserAgent)
				atomic.AddInt64(&completed, 1)
				updateSitemapImportCounts(status, int(atomic.LoadInt64(&completed)), int(atomic.LoadInt64(&skipped)), int(atomic.LoadInt64(&failed)))
			}
		})
	}

	for _, url := range urls {
		select {
		case <-ctx.Done():
			close(jobs)
			waitGroup.Wait()
			if !errors.Is(ctx.Err(), context.Canceled) {
				appendSitemapImportError(status, ctx.Err().Error())
			}
			return
		case jobs <- url:
		}
	}

	close(jobs)
	waitGroup.Wait()

	if err := ctx.Err(); err != nil && !errors.Is(err, context.Canceled) {
		appendSitemapImportError(status, err.Error())
	}
}

func generateAndCacheSitemapPreview(parent context.Context, url string) error {
	ctx, cancel := context.WithTimeout(parent, conf.Config.LinkPreviews.Screenshot.Timeout)
	defer cancel()

	screenshot, err := core.TakeScreenshot(ctx, url, "#link-preview")
	if err != nil {
		if !errors.Is(err, core.ErrMissingSelector) {
			return err
		}

		title, description, fetchErr := core.FetchTitleAndDescription(ctx, url)
		if fetchErr != nil {
			return fetchErr
		}

		screenshot, err = core.TakeScreenshotWithTemplate(ctx, embedfs.DefaultTemplate, url, "#link-preview", title, description)
		if err != nil {
			return err
		}
	}

	compressed, err := core.CompressPNG(screenshot)
	if err != nil {
		slog.Error("PNG compression failed during sitemap import", tint.Err(err), "url", url)
		compressed = screenshot
	}

	return Cache.Write(url, compressed)
}

func updateSitemapImportCounts(status *SitemapImportStatus, completed, skipped, failed int) {
	status.mu.Lock()
	defer status.mu.Unlock()
	status.Completed = completed
	status.Skipped = skipped
	status.Failed = failed
}

func appendSitemapImportError(status *SitemapImportStatus, err string) {
	status.mu.Lock()
	defer status.mu.Unlock()
	status.Errors = append(status.Errors, err)
	if len(status.Errors) > sitemapImportMaxErrors {
		status.Errors = status.Errors[len(status.Errors)-sitemapImportMaxErrors:]
	}
}

func finishSitemapImport(status *SitemapImportStatus) {
	status.mu.Lock()
	status.InProgress = false
	status.cancelFunc = nil
	status.mu.Unlock()
	clearRunningSitemapImport(status)
}

func clearRunningSitemapImport(status *SitemapImportStatus) {
	sitemapImportState.mu.Lock()
	defer sitemapImportState.mu.Unlock()
	if sitemapImportState.status == status {
		sitemapImportState.status = status
	}
}

func cloneSitemapImportStatus(status *SitemapImportStatus) *SitemapImportStatus {
	if status == nil {
		return nil
	}

	status.mu.Lock()
	defer status.mu.Unlock()
	clone := &SitemapImportStatus{
		SitemapURL: status.SitemapURL,
		TotalURLs:  status.TotalURLs,
		Completed:  status.Completed,
		Skipped:    status.Skipped,
		Failed:     status.Failed,
		InProgress: status.InProgress,
		Errors:     append([]string(nil), status.Errors...),
	}
	return clone
}
