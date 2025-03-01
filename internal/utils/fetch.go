package utils

import (
	"regexp"
	"strings"
)

func IsValidURL(webPageUrl string) bool {
	re := regexp.MustCompile(`^https?:\/\/([\da-z\.-]+)\.([a-z\.]{2,6})([\/\w \.-]*)*\/?$`)
	return re.MatchString(webPageUrl)
}

// Helper function to determine if a link is internal
func IsInternalLink(href string) bool {
	// Simplified check for internal links
	return len(href) > 0 && href[0] == '/'
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
