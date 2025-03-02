package handlers

import (
	"log/slog"
	appLogger "lt-app/internal/logger"
	"lt-app/internal/services"
	"lt-app/internal/utils"

	"github.com/gofiber/fiber/v2"
)

// RequestBody represents the expected structure of the request body
type RequestBody struct {
	WebPageUrl string `json:"webPageUrl"`
}

func GetHome(c *fiber.Ctx) error {
	return c.SendFile("./public/html/index.html")
}

func validateRequestBody(c *fiber.Ctx, RLogger *slog.Logger) (RequestBody, string) {
	// Parse the request body
	var body RequestBody
	if err := c.BodyParser(&body); err != nil {
		RLogger.Error("Failed to parse request body", slog.String("error", err.Error()))
		return body, "Invalid request body"
	}

	// Validate the URL
	if !utils.IsValidURL(body.WebPageUrl) {
		RLogger.Warn("Invalid URL format", slog.String("url", body.WebPageUrl))
		return body, "Invalid URL format"
	}

	return body, ""
}

func AnalyzeWebPage(c *fiber.Ctx) error {
	RLogger := appLogger.RLoggerBuilder(c)

	// Parse the request reqBody
	reqBody, errMsg := validateRequestBody(c, RLogger)

	if errMsg != "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": errMsg, "statusCode": "Invalid url format",
		})
	}

	stats, err := services.FetchWebPageStats(reqBody.WebPageUrl, RLogger)

	if err != nil {
		return c.JSON(err)
	}

	RLogger.Info("ExtractedInfo", "webPageUrl", reqBody.WebPageUrl, "Stats", stats)

	// Send JSON response
	return c.JSON(stats)
}
