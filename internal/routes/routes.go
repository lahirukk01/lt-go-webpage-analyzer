package routes

import (
	"lt-app/internal/handlers"
	appLogger "lt-app/internal/logger"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"
)

func buildApiRoutes(app *fiber.App) {
	api := app.Group("/api")

	api.Post("/analyze", handlers.AnalyzeWebPage)
}

func SetupRoutes(app *fiber.App) {
	appLogger.Logger.Info("Setting up routes")
	app.Get("/metrics", monitor.New())

	app.Get("/", handlers.GetHome)

	buildApiRoutes(app)
}
