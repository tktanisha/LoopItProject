package utils

import (
	"testing"
)

func TestValidateFullName(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		expectErr bool
	}{
		{"Empty name", "", true},
		{"Single name", "Tanisha", true},
		{"Valid full name", "Tanisha Khandelwal", false},
		{"Extra spaces", "   Tanisha   Khandelwal   ", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateFullName(tt.input)
			if (err != nil) != tt.expectErr {
				t.Errorf("got error %v, want error? %v", err, tt.expectErr)
			}
		})
	}
}

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		expectErr bool
	}{
		{"Empty email", "", true},
		{"Invalid format no @", "userexample.com", true},
		{"Invalid format no domain", "user@", true},
		{"Valid lowercase email", "user@example.com", false},
		{"Valid uppercase email", "USER@EXAMPLE.COM", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateEmail(tt.input)
			if (err != nil) != tt.expectErr {
				t.Errorf("got error %v, want error? %v", err, tt.expectErr)
			}
		})
	}
}

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		expectErr bool
	}{
		{"Empty password", "", true},
		{"Too short", "Ab1@", true},
		{"No digit", "Abcdef@", true},
		{"No special char", "Abcdef1", true},
		{"No uppercase", "abcdef1@", true},
		{"No lowercase", "ABCDEF1@", true},
		{"Valid password", "Abcdef1@", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePassword(tt.input)
			if (err != nil) != tt.expectErr {
				t.Errorf("got error %v, want error? %v", err, tt.expectErr)
			}
		})
	}
}

func TestValidatePhoneNumber(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		expectErr bool
	}{
		{"Empty phone", "", true},
		{"Too short", "12345", true},
		{"Too long", "123456789012", true},
		{"Non-digit chars", "12345abcde", true},
		{"Valid phone", "9876543210", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePhoneNumber(tt.input)
			if (err != nil) != tt.expectErr {
				t.Errorf("got error %v, want error? %v", err, tt.expectErr)
			}
		})
	}
}

func TestValidateAddress(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		expectErr bool
	}{
		{"Empty address", "", true},
		{"Too short", "abcde", true},
		{"Valid address", "123 Main Street", false},
		{"Address with spaces", "   45B, Central Park   ", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateAddress(tt.input)
			if (err != nil) != tt.expectErr {
				t.Errorf("got error %v, want error? %v", err, tt.expectErr)
			}
		})
	}
}
