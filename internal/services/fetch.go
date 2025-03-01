package services

import (
	"io"
	"log/slog"
	"lt-app/internal/utils"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-resty/resty/v2"
)

// ErrorResponse represents an error response with a status code and message
type ErrorResponse struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
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
	Doc        *goquery.Document
	DoctypeStr string
}

func buildErrorResponse(statusCode int, message string) *ErrorResponse {
	return &ErrorResponse{
		StatusCode: statusCode,
		Message:    message,
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

	return &PageData{
		Doc:        doc,
		DoctypeStr: docTypeStr,
	}, nil
}

func (pd *PageData) getHeadings() map[string]int {
	headings := make(map[string]int)

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

func (pd *PageData) setLinkStats(stats *WebPageStats) {
	// Count internal and external links
	pd.Doc.Find("a").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists {
			if utils.IsInternalLink(href) {
				stats.InternalLinks++
			} else {
				stats.ExternalLinks++
			}
		}
	})
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
	pageData.setLinkStats(stats)

	return stats, nil
}
