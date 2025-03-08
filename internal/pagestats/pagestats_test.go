package pagestats

import (
	"log/slog"
	"lt-app/internal/constants"
	"lt-app/internal/pagedata"
	"lt-app/internal/webfetch"
	"reflect"
	"testing"
)

var callCount = 0

// MockFetcher is a mock implementation of the IFetcher interface
type MockFetcher struct{}

func (m *MockFetcher) GetInaccessibleLinks(urls []string) []string {
	return []string{"https://example.com/inaccessible1"}
}

func (m *MockFetcher) Fetch(webPageurl string, RLogger *slog.Logger) (string, *webfetch.ErrorResponse) {
	return "", nil
}

type MockPageData struct {
	WebPageUrl string
	Headings   map[string]int
}

func (m *MockPageData) GetHtmlVersion() string {
	return "html5"
}

func (m *MockPageData) GetHeadings() map[string]int {
	return m.Headings
}

func (m *MockPageData) GetTitle() string {
	return "Example Title"
}

func (m *MockPageData) ContainsLoginForm() bool {
	return true
}

func (m *MockPageData) GetLinkStats() (*pagedata.Links, []string) {
	var inaccessibleLinks []string
	if callCount == 0 {
		callCount++
		inaccessibleLinks = []string{"https://example.com/inaccessible1"}
	} else {
		inaccessibleLinks = make([]string, constants.INACC_LINKS_MAX_CAP+1)
	}
	return &pagedata.Links{
		Internal: 2,
		External: 3,
		Total:    5,
	}, inaccessibleLinks
}

func TestPageStatsBuilder_BuildSuccess(t *testing.T) {
	// Create a mock PageData
	mockPageData := &MockPageData{
		Headings: map[string]int{"h1": 1, "h2": 2},
	}

	// Create a mock Fetcher
	mockFetcher := &MockFetcher{}

	logger := slog.Default()

	// Create a PageStatsBuilder
	psb := &PageStatsBuilder{}

	// Call the Build function
	pageStats, err := psb.Build(mockPageData, mockFetcher, logger)

	// Assert that there is no error
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Assert that the PageStats is not nil
	if pageStats == nil {
		t.Errorf("Expected PageStats, got nil")
	}

	// Assert that the PageStats fields are correct
	expectedStats := &WebPageStats{
		HTMLVersion:       "html5", // You might need to mock GetHtmlVersion as well
		Title:             "Example Title",
		Headings:          map[string]int{"h1": 1, "h2": 2},
		InternalLinks:     2,
		ExternalLinks:     3,
		TotalLinks:        5,
		InaccessibleLinks: 1,
		HasLoginForm:      true,
	}

	if !reflect.DeepEqual(pageStats, expectedStats) {
		t.Errorf("Expected stats %v, got %v", expectedStats, pageStats)
	}
}

func TestPageStatsBuilder_BuildFailure_TooManyLinks(t *testing.T) {
	// Create a mock PageData
	mockPageData := &MockPageData{
		WebPageUrl: "https://example.com",
		Headings:   map[string]int{"h1": 1, "h2": 2},
	}

	// Create a mock Fetcher
	mockFetcher := &MockFetcher{}

	// Create a logger
	logger := slog.Default()

	// Create a PageStatsBuilder
	psb := &PageStatsBuilder{}

	// Call the Build function
	_, err := psb.Build(mockPageData, mockFetcher, logger)

	// Expecting error for too many links
	if err == nil {
		t.Errorf("Expected error for too many links, got nil")
	}
}
