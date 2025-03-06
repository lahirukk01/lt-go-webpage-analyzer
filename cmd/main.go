package main

import (
	"log"
	appLogger "lt-app/internal/applogger"
	"lt-app/internal/middleware"
	"lt-app/internal/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {
	appLogger.InitLogger()
	app := fiber.New()

	middleware.SetupMiddleware(app)

	routes.SetupRoutes(app)

	// Serve static files from the "public" directory
	app.Static("/static", "./public")

	// appLogger.Logger.Info("Server is running on port 3000")
	log.Fatal(app.Listen(":3000"))
}
