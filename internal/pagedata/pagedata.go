package pagedata

import (
	"log/slog"
	"lt-app/internal/utils"
	"lt-app/internal/webfetch"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type PageData struct {
	Doc           *goquery.Document
	DoctypeStr    string
	WebPageUrl    string
	WebPageOrigin string
}

type Links struct {
	Internal int
	External int
	Total    int
}

type IPageDataBuilder interface {
	Build(webPageUrl string, bodyString string, RLogger *slog.Logger) (*PageData, *webfetch.ErrorResponse)
}

type PageDataBuilder struct{}

func (pdb *PageDataBuilder) Build(webPageUrl string, bodyString string, RLogger *slog.Logger) (*PageData, *webfetch.ErrorResponse) {
	/**
	Go query does not validate html. So no error is returned if the html is invalid.
	Ideally string should be validated for html before parsing it.
	*/
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(bodyString))

	// if err != nil {
	// 	RLogger.Error("Error parsing HTML", "HTML parsing error", err)
	// 	return nil, webfetch.BuildErrorResponse(http.StatusBadRequest, "Web page url contains invalid html content")
	// }

	docTypeStr := utils.ExtractDoctypeFromHtmlSource(bodyString)
	origin, _ := utils.GetOriginFromURL(webPageUrl)

	return &PageData{
		Doc:           doc,
		DoctypeStr:    docTypeStr,
		WebPageUrl:    webPageUrl,
		WebPageOrigin: origin,
	}, nil
}

func (pd *PageData) GetHeadings() map[string]int {
	headings := map[string]int{
		"h1": 0,
		"h2": 0,
		"h3": 0,
		"h4": 0,
		"h5": 0,
		"h6": 0,
	}

	// Example: Count the headings
	pd.Doc.Find("h1, h2, h3, h4, h5, h6").Each(func(i int, s *goquery.Selection) {
		tag := goquery.NodeName(s)
		headings[tag]++
	})

	return headings
}

func (pd *PageData) GetTitle() string {
	return pd.Doc.Find("title").Text()
}

// Function to check if the page contains a login form
func (pd *PageData) ContainsLoginForm() bool {
	doc := pd.Doc
	// Look for form elements with input fields for username and password
	hasUsername := doc.Find("input[type='text'], input[type='email'], input[name='username'], input[name='email']").Length() > 0
	hasPassword := doc.Find("input[type='password']").Length() > 0
	hasSubmit := doc.Find("input[type='submit'], button[type='submit']").Length() > 0

	return hasUsername && hasPassword && hasSubmit
}

func (pd *PageData) GetHtmlVersion() string {
	if pd.DoctypeStr == "HTML" {
		return "HTML5"
	}
	return pd.DoctypeStr
}

func (pd *PageData) GetLinkStats() (*Links, []string) {
	// Store external links in slice to check for accessibility
	var validLinks []string
	links := &Links{}

	// Count internal and external links
	pd.Doc.Find("a").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists && len(href) > 1 && href[0] != '#' && href[0] != '?' {
			if utils.IsInternalLink(href) {
				links.Internal++
				validLinks = append(validLinks, pd.WebPageOrigin+href)
			} else {
				links.External++
				validLinks = append(validLinks, href)
			}
		}
	})

	links.Total = links.Internal + links.External
	return links, validLinks
}
