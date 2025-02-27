package main

import (
	"fmt"
	"lt-app/internal/routes"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

func main() {
	app := fiber.New()

	// Apply the rate limiter middleware
	app.Use(limiter.New(limiter.Config{
		Max:        10,               // Maximum number of requests per time window
		Expiration: 30 * time.Second, // Time window 30 seconds
	}))

	routes.SetupRoutes(app)

	// Serve static files from the "public" directory
	app.Static("/static", "./public")

	fmt.Println("Server is running on port 3000")
	app.Listen(":3000")
}
