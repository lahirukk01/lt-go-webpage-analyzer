package pagestats

import (
	"fmt"
	"log/slog"
	"lt-app/internal/constants"
	"lt-app/internal/pagedata"
	"lt-app/internal/webfetch"
	"net/http"
)

type IPageStatsBuilder interface {
	Build(pageData pagedata.IPageData, RLogger *slog.Logger) (*WebPageStats, *webfetch.ErrorResponse)
}

type PageStatsBuilder struct{}

type WebPageStats struct {
	HTMLVersion       string         `json:"htmlVersion"`
	Title             string         `json:"title"`
	Headings          map[string]int `json:"headings"`
	InternalLinks     int            `json:"internalLinks"`
	ExternalLinks     int            `json:"externalLinks"`
	TotalLinks        int            `json:"totalLinks"`
	InaccessibleLinks int            `json:"inaccessibleLinks"`
	HasLoginForm      bool           `json:"hasLoginForm"`
}

func (psb *PageStatsBuilder) Build(pageData pagedata.IPageData, fetcher webfetch.IFetcher, RLogger *slog.Logger) (*WebPageStats, *webfetch.ErrorResponse) {
	stats := &WebPageStats{
		HTMLVersion:       pageData.GetHtmlVersion(),
		Title:             pageData.GetTitle(),
		Headings:          pageData.GetHeadings(),
		InaccessibleLinks: 0,
		HasLoginForm:      pageData.ContainsLoginForm(),
	}

	links, validLinks := pageData.GetLinkStats()
	stats.InternalLinks = links.Internal
	stats.ExternalLinks = links.External
	stats.TotalLinks = links.Total

	// Won't allow more than 300 links to be checked
	if len(validLinks) > constants.INACC_LINKS_MAX_CAP {
		errMessage := fmt.Sprintf("Too many links to check. Exceeded the %d limit", constants.INACC_LINKS_MAX_CAP)
		return nil, webfetch.BuildErrorResponse(http.StatusBadRequest, errMessage)
	}

	inaccessibleLinks := fetcher.GetInaccessibleLinks(validLinks)
	stats.InaccessibleLinks = len(inaccessibleLinks)
	RLogger.Info("InaccessibleLinks", "inaccessibleLinks", inaccessibleLinks, "inaccessibleLinkCount", stats.InaccessibleLinks)

	return stats, nil
}
