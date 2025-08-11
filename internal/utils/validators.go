package utils

import (
	"errors"
	"regexp"
	"strings"
)

// ValidateFullName checks if the full name is not empty and has at least two parts
func ValidateFullName(input string) error {
	if strings.TrimSpace(input) == "" {
		return errors.New("full name cannot be empty")
	}
	if len(strings.Fields(input)) < 2 {
		return errors.New("please enter at least first and last name")
	}
	return nil
}

// ValidateEmail uses a regex to check basic email format
func ValidateEmail(input string) error {
	if strings.TrimSpace(input) == "" {
		return errors.New("email cannot be empty")
	}
	regex := `^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`
	matched, _ := regexp.MatchString(regex, strings.ToLower(input))
	if !matched {
		return errors.New("invalid email format")
	}
	return nil
}

// ValidatePassword ensures minimum length and one number
func ValidatePassword(input string) error {
	if strings.TrimSpace(input) == "" {
		return errors.New("password cannot be empty")
	}
	if len(input) < 6 {
		return errors.New("password must be at least 6 characters")
	}
	hasDigit := regexp.MustCompile(`[0-9]`).MatchString
	if !hasDigit(input) {
		return errors.New("password must contain at least one digit")
	}
	hasSpecialChar := regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>]`).MatchString
	if !hasSpecialChar(input) {
		return errors.New("password must contain at least one special character")
	}
	hasUpperCase := regexp.MustCompile(`[A-Z]`).MatchString
	if !hasUpperCase(input) {
		return errors.New("password must contain at least one uppercase letter")
	}
	hasLowerCase := regexp.MustCompile(`[a-z]`).MatchString
	if !hasLowerCase(input) {
		return errors.New("password must contain at least one lowercase letter")
	}
	return nil
}

// ValidatePhoneNumber checks if it is exactly 10 digits
func ValidatePhoneNumber(input string) error {
	if strings.TrimSpace(input) == "" {
		return errors.New("phone number cannot be empty")
	}
	regex := `^[0-9]{10}$`
	matched, _ := regexp.MatchString(regex, input)
	if !matched {
		return errors.New("phone number must be 10 digits")
	}
	return nil
}

// ValidateAddress ensures it's not empty and has at least 10 characters
func ValidateAddress(input string) error {
	if strings.TrimSpace(input) == "" {
		return errors.New("address cannot be empty")
	}
	if len(strings.TrimSpace(input)) < 10 {
		return errors.New("address must be at least 10 characters")
	}
	return nil
}
