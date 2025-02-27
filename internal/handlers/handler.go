package handlers

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
)

func GetHome(c *fiber.Ctx) error {
	return c.SendFile("./public/html/index.html")
}

func AnalyzeWebsite(c *fiber.Ctx) error {
	// Print request body
	fmt.Println(string(c.Body()))
	time.Sleep(5 * time.Second)

	// Create a response object
	response := map[string]string{
		"message": "Success",
	}

	// Send JSON response
	return c.JSON(response)
}
