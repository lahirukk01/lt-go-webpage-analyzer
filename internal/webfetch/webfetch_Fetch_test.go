package webfetch

import (
	"fmt"
	"log/slog"
	"lt-app/internal/utils"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var rootDir = utils.GetProjectRoot()

func startServer(statusCode int, contentStr string) *httptest.Server {
	// Start a local HTTP server to serve the mock HTML content in a goroutine
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(statusCode)
		w.Header().Set("Content-Type", "text/html")
		fmt.Println("StatusCode", statusCode)
		if statusCode == http.StatusOK {
			fmt.Fprint(w, contentStr)
		}
	}))
	// defer server.Close()
	return server
}

func TestFetchWebPageSourceContent(t *testing.T) {
	// Load the mock HTML content from file
	htmlContent, err := os.ReadFile(filepath.Join(rootDir, "mocks/scrape.html"))

	if err != nil {
		t.Fatalf("Failed to read mock HTML file: %v", err)
	}

	htmlContentStr := string(htmlContent)

	tests := []struct {
		name          string
		expectedBody  string
		statusCode    int
		expectedError *ErrorResponse
		containsError bool
	}{
		{
			name:          "Successful Fetch",
			expectedBody:  htmlContentStr,
			statusCode:    http.StatusOK,
			expectedError: nil,
			containsError: false,
		},
		{
			name:          "Status Not OK",
			expectedBody:  "",
			statusCode:    http.StatusBadRequest,
			expectedError: &ErrorResponse{StatusCode: http.StatusBadRequest, Error: ""},
			containsError: false,
		},
		{
			name:          "Invalid Url Domain",
			expectedBody:  "",
			statusCode:    0,
			expectedError: &ErrorResponse{StatusCode: http.StatusBadRequest, Error: "Domain of the url seems to be invalid."},
			containsError: true,
		},
	}

	fetcher := &HTTPFetcher{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a logger
			logger := slog.Default()

			var webPageUrl string

			if tt.statusCode == 0 {
				webPageUrl = "https://invalidurl"
			} else {
				// Start a local HTTP server to serve the mock HTML content
				server := startServer(tt.statusCode, htmlContentStr)
				defer server.Close()
				webPageUrl = server.URL
			}

			// Call the fetchWebPageSourceContent function in a goroutine
			bodyStr, errResp := fetcher.Fetch(webPageUrl, logger)

			// Check the error
			if tt.expectedError != nil {
				if errResp == nil {
					t.Errorf("Expected error, got nil")
				} else {
					if errResp.StatusCode != tt.expectedError.StatusCode {
						t.Errorf("Expected status code %d, got %d", tt.expectedError.StatusCode, errResp.StatusCode)
					}
					if tt.containsError {
						if !strings.Contains(errResp.Error, tt.expectedError.Error) {
							t.Errorf("Expected error message to contain %q, got %q", tt.expectedError.Error, errResp.Error)
						}
					} else {
						if errResp.Error != tt.expectedError.Error {
							t.Errorf("Expected error message %q, got %q", tt.expectedError.Error, errResp.Error)
						}
					}
				}
			} else {
				if errResp != nil {
					t.Errorf("Expected no error, got %v", errResp)
				}
			}

			// Check the body
			if bodyStr != tt.expectedBody {
				t.Errorf("Expected body %q, got %q", tt.expectedBody, bodyStr)
			}
		})
	}
}
