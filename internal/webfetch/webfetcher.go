package webfetch

import (
	"io"
	"log/slog"
	"lt-app/internal/constants"
	"lt-app/internal/utils"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
)

// ErrorResponse represents an error response with a status code and message
type ErrorResponse struct {
	StatusCode int    `json:"statusCode"`
	Error      string `json:"error"`
}

type FetchPageSourceResult struct {
	BodyBytes          []byte
	FetchErrorResponse *ErrorResponse
}

type IFetcher interface {
	Fetch(webPageurl string, RLogger *slog.Logger) (string, *ErrorResponse)
	GetInaccessibleLinks(urls []string) []string
}

type HTTPFetcher struct{}

func (f *HTTPFetcher) Fetch(webPageurl string, RLogger *slog.Logger) (string, *ErrorResponse) {
	var wg sync.WaitGroup
	fetchResult := make(chan FetchPageSourceResult)

	wg.Add(1)
	go fetchPageSource(webPageurl, &wg, fetchResult, RLogger)

	go func() {
		wg.Wait()
		close(fetchResult)
	}()

	result := <-fetchResult

	if result.FetchErrorResponse != nil {
		return "", result.FetchErrorResponse
	}

	return string(result.BodyBytes), nil
}

func fetchPageSource(webPageurl string, wg *sync.WaitGroup, fetchResult chan<- FetchPageSourceResult, RLogger *slog.Logger) {
	defer wg.Done()

	client := resty.New().SetTimeout(5 * time.Second) // Set a timeout of 5 seconds
	client.SetDoNotParseResponse(true)                // Do not parse the response body

	resp, err := client.R().
		Get(webPageurl)

	if err != nil {
		RLogger.Error("Request Error:", "error", err)

		var errorResponse *ErrorResponse

		if strings.Contains(err.Error(), "no such host") {
			errorResponse = BuildErrorResponse(http.StatusBadRequest, "Domain of the url seems to be invalid.")
		} else {
			errorResponse = BuildErrorResponse(http.StatusInternalServerError, "Something went wrong")
		}

		fetchResult <- FetchPageSourceResult{nil, errorResponse}
		return
	}

	defer resp.RawBody().Close()

	RLogger.Info("Web page fetch response Status", "url", webPageurl, slog.Int("StatusCode", resp.StatusCode()), slog.String("Status", resp.Status()))

	if resp.StatusCode() >= http.StatusBadRequest {
		fetchResult <- FetchPageSourceResult{nil, BuildErrorResponse(resp.StatusCode(), "")}
		return
	}

	bodyBytes, err := io.ReadAll(resp.RawBody())

	if err != nil {
		RLogger.Error("Error reading response body", "error", err)
		fetchResult <- FetchPageSourceResult{nil, BuildErrorResponse(http.StatusInternalServerError, "Something went wrong")}
		return
	}

	fetchResult <- FetchPageSourceResult{bodyBytes, nil}
}

func BuildErrorResponse(statusCode int, message string) *ErrorResponse {
	var errorMessage = message

	if message == "" {
		switch statusCode {
		case http.StatusNotFound:
			errorMessage = "Page not found from the url provided"
		case http.StatusForbidden:
			errorMessage = "Access denied to the page"
		case http.StatusUnauthorized:
			errorMessage = "Unauthorized access to the page"
		}
	}

	return &ErrorResponse{
		StatusCode: statusCode,
		Error:      errorMessage,
	}
}

func (f *HTTPFetcher) GetInaccessibleLinks(urls []string) []string {
	// Count the number of inaccessible links
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, constants.CONCURRENT_GOROUTINE_LIMIT) // Limit the number of concurrent requests

	// Reduce reallocation by setting the capacity of the slice
	var inaccessibleLinks = make([]string, 0, len(urls))
	var mu sync.Mutex

	wg.Add(len(urls))

	for _, link := range urls {
		semaphore <- struct{}{} // Acquire a semaphore

		go func(link string) {
			defer func() {
				<-semaphore // Release the semaphore
				wg.Done()
			}()
			if !utils.CheckLinkAccessibilityWithResty(link) {
				mu.Lock()
				inaccessibleLinks = append(inaccessibleLinks, link)
				mu.Unlock()
			}
		}(link)
	}

	wg.Wait()

	return inaccessibleLinks
}
