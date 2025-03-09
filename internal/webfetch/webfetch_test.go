package webfetch

import (
	"errors"
	"lt-app/internal/applogger"
	"lt-app/internal/myhttp"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestBuildErrorResponse(t *testing.T) {
	tests := []struct {
		name           string
		statusCode     int
		message        string
		expectedStatus int
		expectedError  string
	}{
		{
			name:           "Custom Message",
			statusCode:     http.StatusBadRequest,
			message:        "Invalid input",
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid input",
		},
		{
			name:           "NotFound Default Message",
			statusCode:     http.StatusNotFound,
			message:        "",
			expectedStatus: http.StatusNotFound,
			expectedError:  "Page not found from the url provided",
		},
		{
			name:           "Forbidden Default Message",
			statusCode:     http.StatusForbidden,
			message:        "",
			expectedStatus: http.StatusForbidden,
			expectedError:  "Access denied to the page",
		},
		{
			name:           "Unauthorized Default Message",
			statusCode:     http.StatusUnauthorized,
			message:        "",
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "Unauthorized access to the page",
		},
		{
			name:           "Other Status Code No Message",
			statusCode:     http.StatusInternalServerError,
			message:        "",
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "",
		},
		{
			name:           "Other Status Code With Message",
			statusCode:     http.StatusInternalServerError,
			message:        "Internal server error",
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "Internal server error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errResp := BuildErrorResponse(tt.statusCode, tt.message)

			if errResp == nil {
				t.Fatalf("buildErrorResponse(%d, %q) returned nil", tt.statusCode, tt.message)
			}

			if errResp.StatusCode != tt.expectedStatus {
				t.Errorf("buildErrorResponse(%d, %q) StatusCode = %d, want %d", tt.statusCode, tt.message, errResp.StatusCode, tt.expectedStatus)
			}

			if errResp.Error != tt.expectedError {
				t.Errorf("buildErrorResponse(%d, %q) Error = %q, want %q", tt.statusCode, tt.message, errResp.Error, tt.expectedError)
			}
		})
	}
}

func TestFetchWebPageSourceContent(t *testing.T) {
	applogger.InitLogger()
	rclient := myhttp.NewRestyClient()
	// Initialize WebFetcher with a mock HTTP client
	fetcher := NewWebFetcher(rclient)

	// Activate httpmock
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	rclient.SetTransport(httpmock.DefaultTransport)

	// Test case 1: Successful link (200 OK)
	httpmock.RegisterResponder("GET", "http://example.com/success",
		httpmock.NewStringResponder(http.StatusOK, "OK"))

	content, err := fetcher.Fetch("http://example.com/success", applogger.Logger)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if content != "OK" {
		t.Errorf("Expected content 'OK', got %s", content)
	}

	// Test case 2: Inaccessible link (404 Not Found)
	httpmock.RegisterResponder("GET", "http://example.com/notfound",
		httpmock.NewStringResponder(http.StatusNotFound, "Not Found"))

	content, err = fetcher.Fetch("http://example.com/notfound", applogger.Logger)
	if err == nil {
		t.Errorf("Missing expected error")
	}
	if content != "" {
		t.Errorf("Expected empty content, got %s", content)
	}

	// Test case 3: Error during request
	httpmock.RegisterResponder("GET", "http://example.com/error",
		func(req *http.Request) (*http.Response, error) {
			return nil, errors.New("simulated error")
		})

	content, err = fetcher.Fetch("http://example.com/error", applogger.Logger)
	if err == nil {
		t.Errorf("Missing expected error")
	}
	if content != "" {
		t.Errorf("Expected empty content, got %s", content)
	}
}

func TestGetInaccessibleLinks(t *testing.T) {
	applogger.InitLogger()

	rclient := myhttp.NewRestyClient()
	// Initialize WebFetcher with a mock HTTP client
	fetcher := NewWebFetcher(rclient)

	// Activate httpmock
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// fetcher.httpClient.(*myhttp.RestyClient).SetTransport(httpmock.DefaultTransport)
	rclient.SetTransport(httpmock.DefaultTransport)

	// Define test URLs
	baseURL := "http://example.com"
	validURL1 := baseURL + "/page1"
	validURL2 := baseURL + "/page2"
	validURL3 := baseURL + "/page3"
	invalidURL1 := baseURL + "/invalid1"
	invalidURL2 := baseURL + "/invalid2"

	// Create a list of URLs to test
	urls := []string{validURL1, validURL2, validURL3, invalidURL1, invalidURL2}

	for i, url := range urls {
		status := http.StatusOK
		// Register a responder for each URL
		if i > 2 {
			status = http.StatusNotFound
		}

		httpmock.RegisterResponder("GET", url,
			httpmock.NewStringResponder(status, "OK"))
	}

	inaccessible := fetcher.GetInaccessibleLinks(urls)

	// Define expected inaccessible URLs
	expected := []string{invalidURL1, invalidURL2}

	// Check if the returned inaccessible URLs match the expected ones
	if !equal(inaccessible, expected) {
		t.Errorf("GetInaccessibleLinks(%v) = %v, want %v", urls, inaccessible, expected)
	}
}

// Helper function to compare two string slices
func equal(arr1, arr2 []string) bool {
	if len(arr1) != len(arr2) {
		return false
	}

	for _, link1 := range arr1 {
		found := false
		for _, link2 := range arr2 {
			if link1 == link2 {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	return true
}

func TestWebFetcher_checkLinkAccessibilityWithResty(t *testing.T) {
	applogger.InitLogger()
	rclient := myhttp.NewRestyClient()
	// Initialize WebFetcher with a mock HTTP client
	fetcher := NewWebFetcher(rclient)

	// Activate httpmock
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	rclient.SetTransport(httpmock.DefaultTransport)

	// Test case 1: Successful link (200 OK)
	httpmock.RegisterResponder("GET", "http://example.com/success",
		httpmock.NewStringResponder(http.StatusOK, "OK"))

	accessible := fetcher.checkLinkAccessibilityWithResty("http://example.com/success")
	assert.True(t, accessible, "Expected link to be accessible")

	// Test case 2: Inaccessible link (404 Not Found)
	httpmock.RegisterResponder("GET", "http://example.com/notfound",
		httpmock.NewStringResponder(http.StatusNotFound, "Not Found"))

	accessible = fetcher.checkLinkAccessibilityWithResty("http://example.com/notfound")
	assert.False(t, accessible, "Expected link to be inaccessible")

	// // Test case 3: Error during request
	httpmock.RegisterResponder("GET", "http://example.com/error",
		func(req *http.Request) (*http.Response, error) {
			return nil, errors.New("simulated error")
		})

	accessible = fetcher.checkLinkAccessibilityWithResty("http://example.com/error")
	assert.False(t, accessible, "Expected link to be inaccessible due to error")
}
