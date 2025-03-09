package handlers

import (
	"log/slog"
	appLogger "lt-app/internal/applogger"
	"lt-app/internal/services"
	"lt-app/internal/utils"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
)

// RequestBody represents the expected structure of the request body
type RequestBody struct {
	WebPageUrl string `json:"webPageUrl"`
}

type RequestBodyValidationErr struct {
	StatusCode int    `json:"statusCode"`
	Error      string `json:"error"`
}

var rootPath = utils.GetProjectRoot()

func GetHome(c *fiber.Ctx) error {
	publicPath := filepath.Join(rootPath, "public")
	return c.SendFile(filepath.Join(publicPath, "html/index.html"))
}

func validateRequestBody(c *fiber.Ctx, RLogger *slog.Logger) (RequestBody, *RequestBodyValidationErr) {
	// Parse the request body
	var body RequestBody
	if err := c.BodyParser(&body); err != nil {
		RLogger.Error("Failed to parse request body", slog.String("error", err.Error()))
		return body, &RequestBodyValidationErr{StatusCode: fiber.StatusBadRequest, Error: "Invalid request body"}
	}

	// Validate the URL
	if !utils.IsValidURL(body.WebPageUrl) {
		RLogger.Warn("Invalid URL format", slog.String("url", body.WebPageUrl))
		return body, &RequestBodyValidationErr{StatusCode: fiber.StatusBadRequest, Error: "Invalid url format"}
	}

	return body, nil
}

func AnalyzeWebPage(c *fiber.Ctx) error {
	RLogger := appLogger.RLoggerBuilder(c)

	// Parse the request reqBody
	reqBody, validationErr := validateRequestBody(c, RLogger)

	if validationErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(validationErr)
	}

	stats, err := services.FetchWebPageStats(reqBody.WebPageUrl, RLogger)

	if err != nil {
		return c.JSON(err)
	}

	RLogger.Info("ExtractedInfo", "webPageUrl", reqBody.WebPageUrl, "Stats", stats)

	// Send JSON response
	return c.JSON(stats)
}
