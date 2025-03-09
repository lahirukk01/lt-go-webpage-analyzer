package utils

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

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

func GetOriginFromURL(urlStr string) (string, error) {
	if !IsValidURL(urlStr) {
		return "", fmt.Errorf("invalid URL: %s", urlStr)
	}

	// Due to the previous check, the URL is valid. So no err
	parsedURL, _ := url.Parse(urlStr)

	origin := fmt.Sprintf("%s://%s", parsedURL.Scheme, parsedURL.Host)
	return origin, nil
}

/*
	Function to check the accessibility of a link using the GET method.

This function has not been used in the codebase. Improved the performance
by using the above function CheckLinkAccessibilityWithResty.
*/
// func CheckLinkAccessibility(url string, wg *sync.WaitGroup, inaccessibleLinksChan chan<- string) {
// 	defer wg.Done()

// 	client := &http.Client{
// 		Timeout: constants.REQUEST_TIMEOUT_SECONDS * time.Second, // Set timeout to 2 seconds
// 	}

// 	req, err := http.NewRequest("GET", url, nil)

// 	if err != nil {
// 		appLogger.Logger.Info("Failed to create request", "url", url, "error", err)
// 		inaccessibleLinksChan <- url
// 		return
// 	}

// 	resp, err := client.Do(req)
// 	if err != nil {
// 		appLogger.Logger.Info("Failed to check link accessibility", "url", url, "error", err)
// 		inaccessibleLinksChan <- url
// 		return
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusOK {
// 		appLogger.Logger.Info("Inaccessible link", "url", url, "statusCode", resp.StatusCode)
// 		inaccessibleLinksChan <- url
// 		return
// 	}
// }

/*
*This function has not been used in the code. However need a html tag validation
function as goquery does not validate the html tags.
*/
// func ValidateHTMLTags(htmlString string) error {
// 	reader := strings.NewReader(htmlString)
// 	tokenizer := html.NewTokenizer(reader)

// 	for {
// 		tokenType := tokenizer.Next()
// 		switch tokenType {
// 		case html.ErrorToken:
// 			err := tokenizer.Err()
// 			if err != nil {
// 				if err.Error() == "EOF" {
// 					return nil // End of file, valid
// 				}
// 				return err // Return other errors
// 			}
// 			return nil //EOF with no error
// 		case html.StartTagToken, html.EndTagToken, html.SelfClosingTagToken:
// 			// You can perform more advanced checks here, such as:
// 			// - Ensuring tags are properly nested.
// 			// - Checking for required attributes.
// 			// - Ensuring tag names are valid.
// 			// For basic well-formedness, the parsing itself is sufficient.
// 		}
// 	}
// }
