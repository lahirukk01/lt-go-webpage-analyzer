package applogger

import (
	"log/slog"
	"os"

	"github.com/gofiber/fiber/v2"
)

var Logger *slog.Logger

func createLogger() *slog.Logger {
	opts := &slog.HandlerOptions{Level: slog.LevelDebug} // Adjust level as needed
	handler := slog.NewJSONHandler(os.Stdout, opts)      // or slog.NewJSONHandler
	return slog.New(handler)
}

func InitLogger() {
	Logger = createLogger()
}

func RLoggerBuilder(c *fiber.Ctx) *slog.Logger {
	requestID := c.Locals("requestid").(string)
	return createLogger().With("request_id", requestID)
}
