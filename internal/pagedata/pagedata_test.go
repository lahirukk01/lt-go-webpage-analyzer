package pagedata

import (
	"log/slog"
	"lt-app/internal/utils"
	"os"
	"path/filepath"
	"testing"
)

var rootDir = utils.GetProjectRoot()

func getMockHtmlContent(filename string, t *testing.T) string {
	// Load the mock HTML content from file
	htmlContent, err := os.ReadFile(filepath.Join(rootDir, "mocks", filename))
	if err != nil {
		// Break the test if the mock HTML file cannot be read
		t.Fatalf("Failed to read mock HTML file: %v", err)
	}
	return string(htmlContent)
}

func TestPageDataBuilder(t *testing.T) {
	// Load the mock HTML content from file
	htmlContentStr := getMockHtmlContent("scrape.html", t)

	RLogger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	expectedPageUrl := "https://www.example.com/page1"
	expectedDoctype := "html"

	builder := &PageDataBuilder{}
	expectedOrigin, _ := utils.GetOriginFromURL(expectedPageUrl)

	pageData, _ := builder.Build(expectedPageUrl, htmlContentStr, RLogger)

	if pageData.WebPageUrl != expectedPageUrl {
		t.Errorf("Expected title %q, got %q", expectedPageUrl, pageData.WebPageUrl)
	}

	if pageData.WebPageOrigin != expectedOrigin {
		t.Errorf("Expected doctype %q, got %q", expectedOrigin, pageData.WebPageOrigin)
	}
	if pageData.DoctypeStr != expectedDoctype {
		t.Errorf("Expected doctype %q, got %q", expectedDoctype, pageData.DoctypeStr)
	}

	headings := pageData.GetHeadings()

	expectedHeadings := map[string]int{
		"h1": 1,
		"h2": 2,
		"h3": 3,
		"h4": 4,
		"h5": 5,
		"h6": 6,
	}

	for heading, count := range headings {
		if count != expectedHeadings[heading] {
			t.Errorf("Expected heading %q count %d, got %d", heading, expectedHeadings[heading], count)
		}
	}

	title := pageData.GetTitle()
	expectedTitle := "Webpage to Scrape"

	if title != expectedTitle {
		t.Errorf("Expected title %q, got %q", expectedTitle, title)
	}

	htmlVersion := pageData.GetHtmlVersion()
	expectedHtmlVersion := "html5"

	if htmlVersion != expectedHtmlVersion {
		t.Errorf("Expected HTML version %q, got %q", expectedHtmlVersion, htmlVersion)
	}

	hasLoginForm := pageData.ContainsLoginForm()

	if !hasLoginForm {
		t.Errorf("Expected page to contain login form")
	}

	linkStats, validLinks := pageData.GetLinkStats()

	expectedLinkStats := &Links{
		Internal: 2,
		External: 6,
	}

	if linkStats.Internal != expectedLinkStats.Internal {
		t.Errorf("Expected internal links %d, got %d", expectedLinkStats.Internal, linkStats.Internal)
	}

	if linkStats.External != expectedLinkStats.External {
		t.Errorf("Expected external links %d, got %d", expectedLinkStats.External, linkStats.External)
	}

	if len(validLinks) != 8 {
		t.Errorf("Expected 8 valid links, got %d", len(validLinks))
	}
}

// func TestPageDataBuilder_BuildFailure(t *testing.T) {
// 	invalidHtmlContentStr := getMockHtmlContent("invalid.html", t)
// 	expectedPageUrl := "https://www.example.com/page2"

// 	RLogger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

// 	builder := &PageDataBuilder{}

// 	pageData, err := builder.Build(expectedPageUrl, invalidHtmlContentStr, RLogger)

// 	fmt.Println("PageData", pageData)

// 	if err == nil {
// 		t.Errorf("Expected error, got nil")
// 	}
// }
