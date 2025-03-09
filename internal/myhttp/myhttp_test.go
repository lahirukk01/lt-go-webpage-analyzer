package myhttp

import (
	"net/http"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"

	"lt-app/internal/constants" // Assuming constants.REQUEST_TIMEOUT_SECONDS is defined
)

func TestNewRestyClient(t *testing.T) {
	client := NewRestyClient()

	if client == nil {
		t.Errorf("NewRestyClient returned nil")
	} else if client.client == nil {
		t.Errorf("client.client is nil")
	} else if client.client.GetClient().Timeout != constants.REQUEST_TIMEOUT_SECONDS*time.Second {
		t.Errorf("Timeout not set correctly")
	}
}

func TestRestyClient_Get(t *testing.T) {
	client := NewRestyClient()

	// Mocking resty's client behavior for testing
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "http://example.com/success",
		httpmock.NewStringResponder(http.StatusOK, "OK"))

	// Configure resty to use httpmock's transport
	client.client.SetTransport(httpmock.DefaultTransport)

	// Test successful response
	resp, err := client.Get("http://example.com/success")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if resp.StatusCode() != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode())
	}

	//Test error response.
	httpmock.RegisterResponder("GET", "http://example.com/error",
		httpmock.NewStringResponder(http.StatusNotFound, "Not Found"))

	resp, err = client.Get("http://example.com/error")

	if err != nil {
		t.Errorf("Did not expect an error, got %v", err)
	}

	if resp.StatusCode() != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, resp.StatusCode())
	}

	//Test a bad url.
	_, err = client.Get("bad url")
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}
