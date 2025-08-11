package utils

import (
	"fmt"
	"strconv"
	"strings"
)

func IntConversion(input string) int {

	// Optional: Trim spaces
	input = strings.TrimSpace(input)

	// Try to convert input to int
	number, err := strconv.Atoi(input)
	if err != nil {
		fmt.Println("Invalid number:", err)
		return 0
	}

	return number
}

func FloatConversion(input string) float64 {

	// Optional: Trim spaces
	input = strings.TrimSpace(input)

	// Try to convert input to float64
	number, err := strconv.ParseFloat(input, 64)
	if err != nil {
		fmt.Println("Invalid number:", err)
		return 0.0
	}

	return number
}
