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

func TestPageDataBuilder_BuildSuccess(t *testing.T) {
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
