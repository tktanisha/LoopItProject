package enums_test

import (
	"encoding/json"
	"loopit/internal/enums"
	"testing"
)

func TestRoleString(t *testing.T) {
	tests := []struct {
		role     enums.Role
		expected string
	}{
		{enums.RoleUser, "user"},
		{enums.RoleLender, "lender"},
		{enums.RoleAdmin, "admin"},
		{enums.Role(-1), "unknown"},
		{enums.Role(99), "unknown"},
	}

	for _, tt := range tests {
		result := tt.role.String()
		if result != tt.expected {
			t.Errorf("Role(%d).String() = %s; want %s", tt.role, result, tt.expected)
		}
	}
}

func TestParseRole(t *testing.T) {
	tests := []struct {
		input       string
		expected    enums.Role
		expectError bool
	}{
		{"user", enums.RoleUser, false},
		{"lender", enums.RoleLender, false},
		{"admin", enums.RoleAdmin, false},
		{"invalid", enums.RoleUser, true},
	}

	for _, tt := range tests {
		role, err := enums.ParseRole(tt.input)
		if tt.expectError {
			if err == nil {
				t.Errorf("ParseRole(%q) expected error, got nil", tt.input)
			}
		} else {
			if err != nil {
				t.Errorf("ParseRole(%q) unexpected error: %v", tt.input, err)
			}
			if role != tt.expected {
				t.Errorf("ParseRole(%q) = %v; want %v", tt.input, role, tt.expected)
			}
		}
	}
}

func TestMarshalJSON(t *testing.T) {
	tests := []struct {
		role     enums.Role
		expected string
	}{
		{enums.RoleUser, `"user"`},
		{enums.RoleLender, `"lender"`},
		{enums.RoleAdmin, `"admin"`},
		{enums.Role(99), `"unknown"`},
	}

	for _, tt := range tests {
		data, err := json.Marshal(tt.role)
		if err != nil {
			t.Errorf("json.Marshal(%v) unexpected error: %v", tt.role, err)
		}
		if string(data) != tt.expected {
			t.Errorf("json.Marshal(%v) = %s; want %s", tt.role, string(data), tt.expected)
		}
	}
}

func TestUnmarshalJSON(t *testing.T) {
	tests := []struct {
		jsonInput   string
		expected    enums.Role
		expectError bool
	}{
		{`"user"`, enums.RoleUser, false},
		{`"lender"`, enums.RoleLender, false},
		{`"admin"`, enums.RoleAdmin, false},
		{`"invalid"`, enums.RoleUser, true},
		{`123`, enums.RoleUser, true}, // invalid type
	}

	for _, tt := range tests {
		var r enums.Role
		err := json.Unmarshal([]byte(tt.jsonInput), &r)
		if tt.expectError {
			if err == nil {
				t.Errorf("json.Unmarshal(%s) expected error, got nil", tt.jsonInput)
			}
		} else {
			if err != nil {
				t.Errorf("json.Unmarshal(%s) unexpected error: %v", tt.jsonInput, err)
			}
			if r != tt.expected {
				t.Errorf("json.Unmarshal(%s) = %v; want %v", tt.jsonInput, r, tt.expected)
			}
		}
	}
}
