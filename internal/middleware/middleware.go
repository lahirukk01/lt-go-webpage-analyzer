package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

func SetupMiddleware(app *fiber.App) {
	// Apply the rate limiter middleware
	app.Use(limiter.New(limiter.Config{
		Max:        10,               // Maximum number of requests per time window
		Expiration: 30 * time.Second, // Time window 30 seconds
	}))

	app.Use(pprof.New())
	app.Use(requestid.New())
	app.Use(logger.New(logger.Config{
		TimeFormat: time.RFC3339,
		Format:     "${time} | ${status} | ${latency} | ${ip} | reqid: ${locals:requestid} | ${method} | ${path} | ${error}\n",
	}))
}
