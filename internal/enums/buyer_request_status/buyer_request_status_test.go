package buyer_request_status_test

import (
	"loopit/internal/enums/buyer_request_status"
	"testing"
)

func TestStatusString(t *testing.T) {
	tests := []struct {
		status   buyer_request_status.Status
		expected string
	}{
		{buyer_request_status.Pending, "Pending"},
		{buyer_request_status.Approved, "Approved"},
		{buyer_request_status.Rejected, "Rejected"},
		{buyer_request_status.Status(-1), "unknown"},
		{buyer_request_status.Status(99), "unknown"},
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
		expected    buyer_request_status.Status
		expectError bool
	}{
		{"Pending", buyer_request_status.Pending, false},
		{"Approved", buyer_request_status.Approved, false},
		{"Rejected", buyer_request_status.Rejected, false},
		{"Invalid", buyer_request_status.Pending, true},
	}

	for _, tt := range tests {
		status, err := buyer_request_status.ParseStatus(tt.input)
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
