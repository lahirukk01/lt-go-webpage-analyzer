package services

import (
	"log/slog"
	"lt-app/internal/pagedata"
	"lt-app/internal/pagestats"
	"lt-app/internal/webfetch"
)

func FetchWebPageStats(webPageUrl string, RLogger *slog.Logger) (*pagestats.WebPageStats, *webfetch.ErrorResponse) {
	webfetcher := &webfetch.HTTPFetcher{}
	bodyString, fetchError := webfetcher.Fetch(webPageUrl, RLogger)

	if fetchError != nil {
		RLogger.Error("Error loading HTTP response body.", "url", webPageUrl, "error", fetchError)
		return nil, fetchError
	}

	pgBuilder := &pagedata.PageDataBuilder{}
	pageData, pdBuildErr := pgBuilder.Build(webPageUrl, bodyString, RLogger)

	if pdBuildErr != nil {
		RLogger.Error("Error building PageData", "error", pdBuildErr)
		return nil, pdBuildErr
	}

	// Create an instance of WebPageStats
	psBuilder := &pagestats.PageStatsBuilder{}
	stats, statBuildErr := psBuilder.Build(pageData, webfetcher, RLogger)

	if statBuildErr != nil {
		RLogger.Error("Error building WebPageStats", "error", statBuildErr)
		return nil, statBuildErr
	}

	return stats, nil
}
