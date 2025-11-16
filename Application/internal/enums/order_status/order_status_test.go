package order_status_test

import (
	"loopit/internal/enums/order_status"
	"testing"
)

func TestStatusString(t *testing.T) {
	tests := []struct {
		status   order_status.Status
		expected string
	}{
		{order_status.InUse, "In Use"},
		{order_status.ReturnRequested, "Return Requested"},
		{order_status.Returned, "Returned"},
		{order_status.Status(-1), "unknown"},
		{order_status.Status(99), "unknown"},
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
		expected    order_status.Status
		expectError bool
	}{
		{"In Use", order_status.InUse, false},
		{"Return Requested", order_status.ReturnRequested, false},
		{"Returned", order_status.Returned, false},
		{"Invalid", order_status.InUse, true},
	}

	for _, tt := range tests {
		status, err := order_status.ParseStatus(tt.input)
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
