package utils

import (
	appLogger "lt-app/internal/applogger"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

func TestIsValidURL(t *testing.T) {
	tests := []struct {
		url      string
		expected bool
	}{
		{"https://example.com", true},
		{"http://example.com", true},
		{"https://example.com/path", true},
		{"http://example.com/path?query=1", true},
		{"https://example.com/path#fragment", true},
		{"ftp://example.com", false},
		{"example.com", false},
		{"http://", false},
		{"", false},
		{"https://", false},
		{"https://example", false},
	}

	for _, test := range tests {
		t.Run(test.url, func(t *testing.T) {
			result := IsValidURL(test.url)
			if result != test.expected {
				t.Errorf("IsValidURL(%q) = %v; want %v", test.url, result, test.expected)
			}
		})
	}
}

func TestIsInternalLink(t *testing.T) {
	tests := []struct {
		href     string
		expected bool
	}{
		{"/internal-link", true},
		{"/another/internal/link", true},
		{"http://example.com/external-link", false},
		{"https://example.com/external-link", false},
		{"//example.com/external-link", false},
		{"internal-link", false},
		{"", false},
		{"/", false},
		{"#", false},
		{"?query=1", false},
	}

	for _, test := range tests {
		t.Run(test.href, func(t *testing.T) {
			result := IsInternalLink(test.href)
			if result != test.expected {
				t.Errorf("IsInternalLink(%q) = %v; want %v", test.href, result, test.expected)
			}
		})
	}
}

func TestExtractDoctypeFromHtmlSource(t *testing.T) {
	tests := []struct {
		htmlSource string
		expected   string
	}{
		{"<!DOCTYPE html>", "html"},
		{"<!DOCTYPE HTML>", "html"},
		{"<!DOCTYPE html PUBLIC \"-//W3C//DTD XHTML 1.0 Transitional//EN\" \"http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd\">", "html"},
		{"<!DOCTYPE HTML PUBLIC \"-//W3C//DTD XHTML 1.0 Transitional//EN\" \"http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd\">", "html"},
		{"<!DOCTYPE svg PUBLIC \"-//W3C//DTD SVG 1.1//EN\" \"http://www.w3.org/Graphics/SVG/1.1/DTD/svg11.dtd\">", "svg"},
		{"<!DOCTYPE unknown>", "unknown"},
		{"<html>", "unknown"},
		{"", "unknown"},
	}

	for _, test := range tests {
		t.Run(test.htmlSource, func(t *testing.T) {
			result := ExtractDoctypeFromHtmlSource(test.htmlSource)
			if result != test.expected {
				t.Errorf("ExtractDoctypeFromHtmlSource(%q) = %v; want %v", test.htmlSource, result, test.expected)
			}
		})
	}
}

type LinkAccessTest struct {
	name           string
	mockStatusCode int
	expectedResult bool
	mockURL        string
}

var tests = []LinkAccessTest{
	{"Accessible Link", http.StatusOK, false, ""},
	{"Inaccessible Link", http.StatusNotFound, true, ""},
	{"Server Error", http.StatusInternalServerError, true, ""},
	{"Invalid URL", 0, true, "://invalid-url"},
}

func TestCheckLinkAccessibilityWithResty(t *testing.T) {
	linkAccessTestRuner(t, CheckLinkAccessibilityWithResty)
}

func TestCheckLinkAccessibility(t *testing.T) {
	linkAccessTestRuner(t, CheckLinkAccessibility)
}

func linkAccessTestRuner(t *testing.T, fetchFunc func(url string, wg *sync.WaitGroup, inaccessibleLinksChan chan<- string)) {
	appLogger.InitLogger()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Determine the URL to use
			urlToTest := test.mockURL
			if urlToTest == "" {
				// Create a mock server
				mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(test.mockStatusCode)
				}))
				defer mockServer.Close()
				urlToTest = mockServer.URL
			}

			// Create a wait group
			var wg sync.WaitGroup
			wg.Add(1)

			// Create a channel to receive inaccessible links
			inaccessibleLinksChan := make(chan string, 1)
			// Call the function to test
			go fetchFunc(urlToTest, &wg, inaccessibleLinksChan)

			// Wait for the function to complete
			wg.Wait()
			close(inaccessibleLinksChan)

			// Check the result
			select {
			case link := <-inaccessibleLinksChan:
				if !test.expectedResult && link != "" {
					t.Errorf("Expected no inaccessible link, but got %v", link)
				}
			default:
				if test.expectedResult {
					t.Error("Expected an inaccessible link, but got none")
				}
			}
		})
	}
}

func TestGetOriginFromURL(t *testing.T) {
	tests := []struct {
		url      string
		expected string
		wantErr  bool
	}{
		{
			url:      "https://example.com/path",
			expected: "https://example.com",
			wantErr:  false,
		},
		{
			url:      "http://example.com",
			expected: "http://example.com",
			wantErr:  false,
		},
		{
			url:      "invalid-url",
			expected: "",
			wantErr:  true,
		},
		{
			url:      "",
			expected: "",
			wantErr:  true,
		},
		{
			url:      "%", // Malformed URL
			expected: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.url, func(t *testing.T) {
			origin, err := GetOriginFromURL(tt.url)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetOriginFromURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if origin != tt.expected {
				t.Errorf("GetOriginFromURL() got = %v, want %v", origin, tt.expected)
			}
		})
	}
}
