package return_request_status_test

import (
	"loopit/internal/enums/return_request_status"
	"testing"
)

func TestStatusString(t *testing.T) {
	tests := []struct {
		status   return_request_status.Status
		expected string
	}{
		{return_request_status.Pending, "Pending"},
		{return_request_status.Approved, "Approved"},
		{return_request_status.Rejected, "Rejected"},
		{return_request_status.Status(-1), "unknown"},
		{return_request_status.Status(99), "unknown"},
	}

	for _, tt := range tests {
		result := tt.status.String()
		if result != tt.expected {
			t.Errorf("Status(%d).String() = %s; want %s", tt.status, result, tt.expected)
		}
	}
}

func TestParseStatus(t *testing.T) {
	tests := []struct {
		input       string
		expected    return_request_status.Status
		expectError bool
	}{
		{"Pending", return_request_status.Pending, false},
		{"Approved", return_request_status.Approved, false},
		{"Rejected", return_request_status.Rejected, false},
		{"Invalid", return_request_status.Pending, true},
	}

	for _, tt := range tests {
		status, err := return_request_status.ParseStatus(tt.input)
		if tt.expectError {
			if err == nil {
				t.Errorf("ParseStatus(%q) expected error, got nil", tt.input)
			}
		} else {
			if err != nil {
				t.Errorf("ParseStatus(%q) unexpected error: %v", tt.input, err)
			}
			if status != tt.expected {
				t.Errorf("ParseStatus(%q) = %v; want %v", tt.input, status, tt.expected)
			}
		}
	}
}
