package utils

import (
	"context"
	"fmt"
	appLogger "lt-app/internal/logger"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
)

func IsValidURL(webPageUrl string) bool {
	re := regexp.MustCompile(`^https?:\/\/([\da-z\.-]+)\.([a-z\.]{2,6})([\/\w \.-]*)*\/?(\?[^\s]*)?(#[^\s]*)?$`)
	return re.MatchString(webPageUrl)
}

// Helper function to determine if a link is internal
func IsInternalLink(href string) bool {
	// Simplified check for internal links
	return len(href) > 1 && (href[0] == '/')
}

// Function to extract the Doctype from the HTML source string
func ExtractDoctypeFromHtmlSource(htmlSource string) string {
	re := regexp.MustCompile(`(?i)<!DOCTYPE\s+([^>]+)>`)
	matches := re.FindStringSubmatch(htmlSource)
	if len(matches) > 1 {
		return strings.ToUpper(matches[1])
	}
	return "unknown"
}

/*
Function to check the accessibility of a link using the HEAD method with Resty.
This function has not been used in the codebase.
*/
func CheckLinkAccessibilityWithResty(url string, wg *sync.WaitGroup, inaccessibleLinksChan chan<- string) {
	defer wg.Done()

	client := resty.New().SetTimeout(5 * time.Second)

	resp, err := client.R().Head(url)

	if err != nil {
		appLogger.Logger.Info("Failed to check link accessibility", "url", url, "error", err)
		inaccessibleLinksChan <- url
		return
	}

	defer resp.RawBody().Close()

	if resp.StatusCode() != http.StatusOK {
		appLogger.Logger.Info("Inaccessible link", "url", url, "statusCode", resp.StatusCode())
		inaccessibleLinksChan <- url
	}
}

// Function to check the accessibility of a link using the HEAD method
func CheckLinkAccessibility(url string, wg *sync.WaitGroup, inaccessibleLinksChan chan<- string) {
	defer wg.Done()

	client := &http.Client{
		Timeout: 5 * time.Second, // Set timeout to 5 seconds
	}

	// Create a context with a timeout of 5 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "HEAD", url, nil)
	if err != nil {
		appLogger.Logger.Info("Failed to create request", "url", url, "error", err)
		inaccessibleLinksChan <- url
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		appLogger.Logger.Info("Failed to check link accessibility", "url", url, "error", err)
		inaccessibleLinksChan <- url
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		appLogger.Logger.Info("Inaccessible link", "url", url, "statusCode", resp.StatusCode)
		inaccessibleLinksChan <- url
	}
}

func GetOriginFromURL(urlStr string) (string, error) {
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return "", err
	}
	origin := fmt.Sprintf("%s://%s", parsedURL.Scheme, parsedURL.Host)
	return origin, nil
}
