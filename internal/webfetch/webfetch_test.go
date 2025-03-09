package webfetch

import (
	"net/http"
	"testing"
)

func TestBuildErrorResponse(t *testing.T) {
	tests := []struct {
		name           string
		statusCode     int
		message        string
		expectedStatus int
		expectedError  string
	}{
		{
			name:           "Custom Message",
			statusCode:     http.StatusBadRequest,
			message:        "Invalid input",
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid input",
		},
		{
			name:           "NotFound Default Message",
			statusCode:     http.StatusNotFound,
			message:        "",
			expectedStatus: http.StatusNotFound,
			expectedError:  "Page not found from the url provided",
		},
		{
			name:           "Forbidden Default Message",
			statusCode:     http.StatusForbidden,
			message:        "",
			expectedStatus: http.StatusForbidden,
			expectedError:  "Access denied to the page",
		},
		{
			name:           "Unauthorized Default Message",
			statusCode:     http.StatusUnauthorized,
			message:        "",
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "Unauthorized access to the page",
		},
		{
			name:           "Other Status Code No Message",
			statusCode:     http.StatusInternalServerError,
			message:        "",
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "",
		},
		{
			name:           "Other Status Code With Message",
			statusCode:     http.StatusInternalServerError,
			message:        "Internal server error",
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "Internal server error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errResp := BuildErrorResponse(tt.statusCode, tt.message)

			if errResp == nil {
				t.Fatalf("buildErrorResponse(%d, %q) returned nil", tt.statusCode, tt.message)
			}

			if errResp.StatusCode != tt.expectedStatus {
				t.Errorf("buildErrorResponse(%d, %q) StatusCode = %d, want %d", tt.statusCode, tt.message, errResp.StatusCode, tt.expectedStatus)
			}

			if errResp.Error != tt.expectedError {
				t.Errorf("buildErrorResponse(%d, %q) Error = %q, want %q", tt.statusCode, tt.message, errResp.Error, tt.expectedError)
			}
		})
	}
}
