package services

import (
	"io"
	"log/slog"
	"lt-app/internal/utils"
	"net/http"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-resty/resty/v2"
)

var MAX_LINKS_PROCESS = 170
var CONCURRENT_GOROUTINE_LIMIT = 20

// ErrorResponse represents an error response with a status code and message
type ErrorResponse struct {
	StatusCode int    `json:"statusCode"`
	Error      string `json:"error"`
}

type WebPageStats struct {
	HTMLVersion       string         `json:"htmlVersion"`
	Title             string         `json:"title"`
	Headings          map[string]int `json:"headings"`
	InternalLinks     int            `json:"internalLinks"`
	ExternalLinks     int            `json:"externalLinks"`
	InaccessibleLinks int            `json:"inaccessibleLinks"`
	HasLoginForm      bool           `json:"hasLoginForm"`
}

type PageData struct {
	Doc           *goquery.Document
	DoctypeStr    string
	WebPageUrl    string
	webPageOrigin string
}

func buildErrorResponse(statusCode int, message string) *ErrorResponse {
	return &ErrorResponse{
		StatusCode: statusCode,
		Error:      message,
	}
}

func fetchWebPageSourceContent(webPageurl string, RLogger *slog.Logger) (string, *ErrorResponse) {
	client := resty.New()
	client.SetDoNotParseResponse(true) // Do not parse the response body

	resp, err := client.R().
		Get(webPageurl)

	if err != nil {
		RLogger.Error("Request Error:", "error", err)
		return "", buildErrorResponse(http.StatusBadRequest, "Invalid web page URL")
	}

	defer resp.RawBody().Close()

	RLogger.Info("Response Status", slog.Int("StatusCode", resp.StatusCode()), slog.String("Status", resp.Status()))

	if resp.StatusCode() >= http.StatusBadRequest {
		return "", buildErrorResponse(resp.StatusCode(), resp.Status())
	}

	bodyBytes, err := io.ReadAll(resp.RawBody())

	if err != nil {
		RLogger.Error("Error reading response body", "error", err)
		return "", buildErrorResponse(http.StatusInternalServerError, "Something went wrong")
	}
	bodyString := string(bodyBytes)

	return bodyString, nil
}

func GetWebPageData(webPageurl string, RLogger *slog.Logger) (*PageData, *ErrorResponse) {
	bodyString, errResp := fetchWebPageSourceContent(webPageurl, RLogger)

	if errResp != nil {
		return nil, errResp
	}

	docTypeStr := utils.ExtractDoctypeFromHtmlSource(bodyString)

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(bodyString))

	if err != nil {
		RLogger.Error("Error parsing HTML", "HTML parsing error", err)

		return nil, buildErrorResponse(http.StatusInternalServerError, "Something went wrong")
	}

	origin, _ := utils.GetOriginFromURL(webPageurl)

	return &PageData{
		Doc:           doc,
		DoctypeStr:    docTypeStr,
		WebPageUrl:    webPageurl,
		webPageOrigin: origin,
	}, nil
}

func (pd *PageData) getHeadings() map[string]int {
	headings := map[string]int{
		"h1": 0,
		"h2": 0,
		"h3": 0,
		"h4": 0,
		"h5": 0,
		"h6": 0,
	}

	// Example: Count the headings
	pd.Doc.Find("h1, h2, h3, h4, h5, h6").Each(func(i int, s *goquery.Selection) {
		tag := goquery.NodeName(s)
		headings[tag]++
	})

	return headings
}

func (pd *PageData) getTitle() string {
	return pd.Doc.Find("title").Text()
}

// Function to check if the page contains a login form
func (pd *PageData) containsLoginForm() bool {
	doc := pd.Doc
	// Look for form elements with input fields for username and password
	hasUsername := doc.Find("input[type='text'], input[type='email'], input[name='username'], input[name='email']").Length() > 0
	hasPassword := doc.Find("input[type='password']").Length() > 0
	hasSubmit := doc.Find("input[type='submit'], button[type='submit']").Length() > 0

	return hasUsername && hasPassword && hasSubmit
}

func (pd *PageData) getHtmlVersion() string {
	if pd.DoctypeStr == "HTML" {
		return "HTML5"
	}
	return pd.DoctypeStr
}

func (pd *PageData) setLinkStats(stats *WebPageStats) []string {
	// Store external links in slice to check for accessibility
	var validLinks []string

	// Count internal and external links
	pd.Doc.Find("a").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists && len(href) > 1 && href[0] != '#' && href[0] != '?' {
			if utils.IsInternalLink(href) {
				stats.InternalLinks++
				validLinks = append(validLinks, pd.webPageOrigin+href)
			} else {
				stats.ExternalLinks++
				validLinks = append(validLinks, href)
			}
		}
	})
	return validLinks
}

func setInaccessibleLinksCount(urls []string, stats *WebPageStats, RLogger *slog.Logger) {
	// Count the number of inaccessible links
	var wg sync.WaitGroup
	inaccessibleLinksChan := make(chan string)
	semaphore := make(chan struct{}, CONCURRENT_GOROUTINE_LIMIT) // Limit the number of concurrent requests

	inaccessibleLinks := []string{}

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

	wg.Wait()
	close(inaccessibleLinksChan)

	for link := range inaccessibleLinksChan {
		stats.InaccessibleLinks++
		inaccessibleLinks = append(inaccessibleLinks, link)
	}

	RLogger.Info("InaccessibleLinks", "inaccessibleLinks", inaccessibleLinks, "inaccessibleLinkCount", stats.InaccessibleLinks)
}

func FetchWebPageStats(webPageUrl string, RLogger *slog.Logger) (*WebPageStats, *ErrorResponse) {
	pageData, err := GetWebPageData(webPageUrl, RLogger)

	if err != nil {
		RLogger.Error("Error loading HTTP response body.", "error", err)
		return nil, err
	}

	// Create an instance of WebPageStats
	stats := &WebPageStats{
		Title:        pageData.getTitle(),
		Headings:     pageData.getHeadings(),
		HTMLVersion:  pageData.getHtmlVersion(),
		HasLoginForm: pageData.containsLoginForm(),
	}

	// Count internal and external links
	validLinks := pageData.setLinkStats(stats)

	/* When exceeding the limit, process only the first 160 links.
	Otherwise channel keeps hanging and then the request
	*/
	validLinksCount := utils.Min(len(validLinks), MAX_LINKS_PROCESS)
	setInaccessibleLinksCount(validLinks[0:validLinksCount], stats, RLogger)

	return stats, nil
}
