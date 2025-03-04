package appLogger

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

func TestInitLogger(t *testing.T) {
	// Backup the original logger and replace it with a temporary one
	originalLogger := Logger
	defer func() {
		Logger = originalLogger
	}()

	// Call InitLogger
	InitLogger()

	// Check if the logger is initialized
	if Logger == nil {
		t.Error("Logger is not initialized")
	}

	// Check if the logger level is correct
	if Logger.Enabled(context.TODO(), slog.LevelInfo) != true {
		t.Error("Logger level is not set to Info")
	}
}

func TestLoggerMethods(t *testing.T) {
	// Backup the original logger and replace it with a temporary one
	originalLogger := Logger
	defer func() {
		Logger = originalLogger
	}()

	// Create a buffer to capture the logger output
	var buf bytes.Buffer
	handler := slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug})
	Logger = slog.New(handler)

	// Test Debug method
	Logger.Debug("Debug message", "key", "value")
	if !strings.Contains(buf.String(), "level=DEBUG") {
		t.Error("Debug message not logged")
	}
	buf.Reset()

	// Test Info method
	Logger.Info("Info message", "key", "value")
	if !strings.Contains(buf.String(), "level=INFO") {
		t.Error("Info message not logged")
	}
	buf.Reset()

	// Test Warn method
	Logger.Warn("Warn message", "key", "value")
	if !strings.Contains(buf.String(), "level=WARN") {
		t.Error("Warn message not logged")
	}
	buf.Reset()

	// Test Error method
	Logger.Error("Error message", "key", "value")
	if !strings.Contains(buf.String(), "level=ERROR") {
		t.Error("Error message not logged")
	}
	buf.Reset()
}

// captureOutput captures os.Stdout and returns it as a string.
func captureOutput(f func()) string {
	old := os.Stdout // keep backup of the real stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	err := w.Close()
	if err != nil {
		panic(err)
	}
	os.Stdout = old // restoring the real stdout
	var buf bytes.Buffer
	_, err = buf.ReadFrom(r)
	if err != nil {
		panic(err)
	}
	return buf.String()
}

func TestRLoggerBuilder(t *testing.T) {
	// Backup the original logger and replace it with a temporary one
	// originalLogger := Logger
	// defer func() {
	// 	Logger = originalLogger
	// }()

	// // Create a buffer to capture the logger output
	// var buf bytes.Buffer
	// handler := slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug})
	// logger = slog.New(handler)

	// Create a new Fiber app
	app := fiber.New()
	app.Use(requestid.New())

	// Define a test route
	app.Get("/test", func(c *fiber.Ctx) error {
		// Set a request ID in the context
		c.Locals("requestid", "test-request-id")

		old := os.Stdout // keep backup of the real stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		// Build the logger using RLoggerBuilder
		logger := RLoggerBuilder(c)

		// Check if the logger is initialized
		if logger == nil {
			t.Error("RLoggerBuilder returned nil logger")
			return nil
		}

		// Log a message
		logger.Info("Test message")

		err := w.Close()
		if err != nil {
			panic(err)
		}
		os.Stdout = old // restoring the real stdout
		var buf bytes.Buffer
		_, err = buf.ReadFrom(r)
		if err != nil {
			panic(err)
		}
		output := buf.String()
		fmt.Println("Output", output)

		// Check if the request ID is present in the log output
		if !strings.Contains(output, "test-request-id") {
			t.Errorf("Request ID not found in log output: %s", output)
		}

		return c.SendString("OK")
	})

	// Perform a test request
	req := httptest.NewRequest("GET", "/test", nil)
	resp, _ := app.Test(req)

	// Check the response status code
	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("Expected status code %d, got %d", fiber.StatusOK, resp.StatusCode)
	}
}
