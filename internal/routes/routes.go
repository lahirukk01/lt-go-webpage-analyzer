package routes

import (
	"lt-app/internal/handlers"
	"lt-app/internal/utils"

	"github.com/gofiber/fiber/v2"
)

func buildApiRoutes(app *fiber.App) {
	api := app.Group("/api")

	api.Post("/analyze", handlers.AnalyzeWebPage)
}

func SetupRoutes(app *fiber.App) {
	utils.Logger.Info("Setting up routes")
	app.Get("/", handlers.GetHome)

	buildApiRoutes(app)
}
