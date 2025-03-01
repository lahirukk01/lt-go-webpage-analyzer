package main

import (
	"log"
	"lt-app/internal/middleware"
	"lt-app/internal/routes"
	"lt-app/internal/utils"

	"github.com/gofiber/fiber/v2"
)

func main() {
	utils.InitLogger()
	app := fiber.New()

	middleware.SetupMiddleware(app)

	routes.SetupRoutes(app)

	// Serve static files from the "public" directory
	app.Static("/static", "./public")

	utils.Logger.Info("Server is running on port 3000")
	log.Fatal(app.Listen(":3000"))
}
