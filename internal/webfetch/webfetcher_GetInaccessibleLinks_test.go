package webfetch

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	appLogger "lt-app/internal/applogger"
)

func startServer_accessTest() *httptest.Server {
	// Start a local HTTP server to serve the mock HTML content in a goroutine
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/page1":
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, "This is page 1")
		case "/page2":
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, "This is page 2")
		case "/page3":
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, "This is page 3")
		default:
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "404 Not Found")
		}
	}))
	return server
}

func TestGetInaccessibleLinks(t *testing.T) {
	// Spinnup server with 3 valid GET routes
	server := startServer_accessTest()
	defer server.Close()

	appLogger.InitLogger()

	// Define test URLs
	baseURL := server.URL
	validURL1 := baseURL + "/page1"
	validURL2 := baseURL + "/page2"
	validURL3 := baseURL + "/page3"
	invalidURL1 := baseURL + "/invalid1"
	invalidURL2 := baseURL + "/invalid2"

	// Create a list of URLs to test
	urls := []string{validURL1, validURL2, validURL3, invalidURL1, invalidURL2}

	// Call GetInaccessibleLinks
	fetcher := &HTTPFetcher{}
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
