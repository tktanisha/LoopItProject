package utils

import (
	"errors"
	"loopit/internal/models"
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

func ValidateDescription(input string) error {
	if strings.TrimSpace(input) == "" {
		return errors.New("description cannot be empty")
	}
	return nil
}

func ValidateOrderID(input int) error {
	if input <= 0 {
		return errors.New("order_id must be a positive integer")
	}
	return nil
}

func ValidateFeedbackText(input string) error {
	if strings.TrimSpace(input) == "" {
		return errors.New("feedback_text cannot be empty")
	}
	return nil
}

func ValidateRating(input int) error {
	if input < 1 || input > 5 {
		return errors.New("rating must be between 1 and 5")
	}
	return nil
}

func ValidateProduct(p *models.Product) error {
	if p.CategoryID <= 0 {
		return errors.New("category_id must be a positive integer")
	}
	if strings.TrimSpace(p.Name) == "" {
		return errors.New("name cannot be empty")
	}
	if len(p.Name) > 100 {
		return errors.New("name cannot exceed 100 characters")
	}
	if len(p.Description) > 500 {
		return errors.New("description cannot exceed 500 characters")
	}
	if p.Duration <= 0 {
		return errors.New("duration must be greater than 0")
	}
	return nil
}

// ValidateSociety validates the Society creation request
func ValidateSociety(name, location, pincode string) error {
	if strings.TrimSpace(name) == "" {
		return errors.New("society name cannot be empty")
	}
	if len(name) < 3 {
		return errors.New("society name must be at least 3 characters")
	}
	if strings.TrimSpace(location) == "" {
		return errors.New("location cannot be empty")
	}
	if len(location) < 5 {
		return errors.New("location must be at least 5 characters")
	}
	// Inline pincode validation (instead of separate ValidatePincode)
	if strings.TrimSpace(pincode) == "" {
		return errors.New("pincode cannot be empty")
	}
	regex := `^[0-9]{6}$`
	matched, _ := regexp.MatchString(regex, pincode)
	if !matched {
		return errors.New("pincode must be 6 digits")
	}
	return nil
}

// ValidateCategory validates category creation request
func ValidateCategory(name string, price, security float64) error {
	if strings.TrimSpace(name) == "" {
		return errors.New("category name cannot be empty")
	}
	if len(name) < 3 {
		return errors.New("category name must be at least 3 characters")
	}
	if price <= 0 {
		return errors.New("price must be a positive number")
	}
	if security < 0 {
		return errors.New("security cannot be negative")
	}
	return nil
}
