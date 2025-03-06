package utils

import (
	"fmt"
	appLogger "lt-app/internal/applogger"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
)

const REQUEST_TIMEOUT_SECONDS = 3

func IsValidURL(webPageUrl string) bool {
	re := regexp.MustCompile(`^https?:\/\/([\da-z\.-]+)\.([a-z\.]{2,6})([\/\w \.-]*)*\/?(\?[^\s]*)?(#[^\s]*)?$`)
	return re.MatchString(webPageUrl)
}

// Helper function to determine if a link is internal
func IsInternalLink(href string) bool {
	// Regex pattern to match internal links
	re := regexp.MustCompile(`^\/[^\/].*`)
	return re.MatchString(href)
}

// Function to extract the Doctype from the HTML source string
func ExtractDoctypeFromHtmlSource(htmlSource string) string {
	// re := regexp.MustCompile(`(?i)<!DOCTYPE\s+([^>]+)>`)
	re := regexp.MustCompile(`(?i)<!DOCTYPE\s+([^\s>]+)`)
	matches := re.FindStringSubmatch(htmlSource)
	if len(matches) > 1 {
		return strings.ToLower(matches[1])
	}
	return "unknown"
}

/*
Function to check the accessibility of a link using the HEAD method with Resty.
This function has not been used in the codebase. Done() is not called on the WaitGroup
when context timeout occurs. This can lead to a deadlock in the application. Not using
this function in the codebase.
*/
func CheckLinkAccessibilityWithResty(url string, wg *sync.WaitGroup, inaccessibleLinksChan chan<- string) {
	defer wg.Done()

	client := resty.New().SetTimeout(REQUEST_TIMEOUT_SECONDS * time.Second)

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
		return
	}
}

// Function to check the accessibility of a link using the HEAD method
func CheckLinkAccessibility(url string, wg *sync.WaitGroup, inaccessibleLinksChan chan<- string) {
	defer wg.Done()

	client := &http.Client{
		Timeout: REQUEST_TIMEOUT_SECONDS * time.Second, // Set timeout to 2 seconds
	}

	req, err := http.NewRequest("HEAD", url, nil)

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
		return
	}
}

func GetOriginFromURL(urlStr string) (string, error) {
	if !IsValidURL(urlStr) {
		return "", fmt.Errorf("invalid URL: %s", urlStr)
	}

	// Due to the previous check, the URL is valid. So no err
	parsedURL, _ := url.Parse(urlStr)

	origin := fmt.Sprintf("%s://%s", parsedURL.Scheme, parsedURL.Host)
	return origin, nil
}
