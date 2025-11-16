package utils_test

import (
	"encoding/json"
	"loopit/internal/utils"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestWriteErrorResponse(t *testing.T) {
	// Prepare a response recorder to capture the output
	rr := httptest.NewRecorder()
	statusCode := http.StatusBadRequest
	errMsg := "Invalid request"
	details := "Missing required fields"

	// Call the function under test
	utils.WriteErrorResponse(rr, statusCode, errMsg, details)

	// Verify HTTP status code
	if rr.Code != statusCode {
		t.Errorf("Expected status code %d, got %d", statusCode, rr.Code)
	}

	// Verify Content-Type header
	contentType := rr.Header().Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		t.Errorf("Expected Content-Type to contain application/json, got %s", contentType)
	}

	// Decode the JSON response
	var resp utils.ErrorResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("Failed to decode JSON response: %v", err)
	}

	// Validate response fields
	if resp.Status != false {
		t.Errorf("Expected Status to be false, got %v", resp.Status)
	}
	if resp.StatusCode != statusCode {
		t.Errorf("Expected StatusCode to be %d, got %d", statusCode, resp.StatusCode)
	}
	if resp.Error != errMsg {
		t.Errorf("Expected Error to be %q, got %q", errMsg, resp.Error)
	}
	if resp.Details != details {
		t.Errorf("Expected Details to be %q, got %q", details, resp.Details)
	}
}
