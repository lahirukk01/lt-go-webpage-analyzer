package pagestats

import (
	"fmt"
	"log/slog"
	"lt-app/internal/pagedata"
	"lt-app/internal/utils"
	"lt-app/internal/webfetch"
	"net/http"
	"sync"
)

const INACC_LINKS_MAX_CAP = 300
const CONCURRENT_GOROUTINE_LIMIT = 20

type WebPageStats struct {
	HTMLVersion       string         `json:"htmlVersion"`
	Title             string         `json:"title"`
	Headings          map[string]int `json:"headings"`
	InternalLinks     int            `json:"internalLinks"`
	ExternalLinks     int            `json:"externalLinks"`
	TotalLinks        int            `json:"totalLinks"`
	InaccessibleLinks int            `json:"inaccessibleLinks"`
	HasLoginForm      bool           `json:"hasLoginForm"`
}

func BuildWebPageStats(pageData *pagedata.PageData, RLogger *slog.Logger) (*WebPageStats, *webfetch.ErrorResponse) {
	stats := &WebPageStats{
		HTMLVersion:       pageData.DoctypeStr,
		Title:             pageData.GetTitle(),
		Headings:          pageData.GetHeadings(),
		InaccessibleLinks: 0,
		HasLoginForm:      pageData.ContainsLoginForm(),
	}

	links, validLinks := pageData.GetLinkStats()
	stats.InternalLinks = links.Internal
	stats.ExternalLinks = links.External
	stats.TotalLinks = links.Total

	// Won't allow more than 300 links to be checked
	if len(validLinks) > INACC_LINKS_MAX_CAP {
		return nil, webfetch.BuildErrorResponse(http.StatusBadRequest, fmt.Sprintf("Too many links to check. Exceeded the %d limit", INACC_LINKS_MAX_CAP))
	}

	inaccessibleLinks := getInaccessibleLinks(validLinks)
	stats.InaccessibleLinks = len(inaccessibleLinks)
	RLogger.Info("InaccessibleLinks", "inaccessibleLinks", inaccessibleLinks, "inaccessibleLinkCount", stats.InaccessibleLinks)

	return stats, nil
}

func getInaccessibleLinks(urls []string) []string {
	// Count the number of inaccessible links
	var wg sync.WaitGroup
	inaccessibleLinksChan := make(chan string, len(urls))
	semaphore := make(chan struct{}, CONCURRENT_GOROUTINE_LIMIT) // Limit the number of concurrent requests

	// Reduce reallocation by setting the capacity of the slice
	inaccessibleLinks := make([]string, 0, min(len(urls), INACC_LINKS_MAX_CAP))

	wg.Add(len(urls))

	for _, link := range urls {
		semaphore <- struct{}{} // Acquire a semaphore

		go func(link string) {
			defer func() {
				<-semaphore // Release the semaphore
			}()
			utils.CheckLinkAccessibility(link, &wg, inaccessibleLinksChan)
		}(link)
	}

	go func() {
		wg.Wait()
		close(inaccessibleLinksChan)
	}()

	for link := range inaccessibleLinksChan {
		inaccessibleLinks = append(inaccessibleLinks, link)
	}
	return inaccessibleLinks
}
