package routes

import (
	"fmt"
	"lt-app/internal/handlers"

	"github.com/gofiber/fiber/v2"
)

func buildApiRoutes(app *fiber.App) {
	api := app.Group("/api")

	api.Post("/analyze", handlers.AnalyzeWebsite)
}

func SetupRoutes(app *fiber.App) {
	fmt.Println("Setting up routes")
	app.Get("/", handlers.GetHome)

	buildApiRoutes(app)
}
