package main

import (
	"fmt"
	"log"
	"lt-app/internal/applogger"
	"lt-app/internal/constants"
	"lt-app/internal/middleware"
	"lt-app/internal/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {
	applogger.InitLogger()
	app := fiber.New()

	middleware.SetupMiddleware(app)

	routes.SetupRoutes(app)

	// Serve static files from the "public" directory
	app.Static("/static", "./public")

	applogger.Logger.Info(fmt.Sprintf("Server is listening on port %d", constants.SERVER_PORT))
	log.Fatal(app.Listen(fmt.Sprintf(":%d", constants.SERVER_PORT)))
}
