package handlers

import (
	"io"
	appLogger "lt-app/internal/logger"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func setupTestApp() *fiber.App {

	// Initialize logger
	appLogger.InitLogger()

	// Create a new Fiber app
	app := fiber.New()

	return app
}

func TestGetHome(t *testing.T) {
	app := setupTestApp()

	// Define the GetHome route
	app.Get("/", GetHome)

	// Create a test request
	req := httptest.NewRequest("GET", "/", nil)

	// Perform the request
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Error testing / route: %v", err)
	}

	// Check the status code
	if resp.StatusCode != http.StatusOK {
		t.Errorf("/ route: expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Error reading response body: %v", err)
	}

	pageContent := string(body)
	// Assert that the response body contains the expected content
	if !strings.Contains(pageContent, "Welcome to the Web Page Analyzer") {
		t.Errorf("Unexpected response body: %s", pageContent)
	}
}
