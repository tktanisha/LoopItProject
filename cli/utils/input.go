package utils

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
)

// Input prompts for simple input (no validation)
func Input(label string) string {
	var result string
	prompt := &survey.Input{
		Message: label,
	}
	err := survey.AskOne(prompt, &result)
	if err != nil {
		fmt.Println("Prompt failed:", err)
		return ""
	}
	return result
}

// InputWithValidation prompts for input with a validator
func InputWithValidation(label string, validate func(string) error) string {
	var result string
	prompt := &survey.Input{
		Message: label,
	}

	wrappedValidator := func(ans interface{}) error {
		str, ok := ans.(string)
		if !ok {
			return fmt.Errorf("input must be a string")
		}
		return validate(str)
	}

	err := survey.AskOne(prompt, &result, survey.WithValidator(wrappedValidator))
	if err != nil {
		fmt.Println("Prompt failed:", err)
		return ""
	}
	return result
}

// InputPassword prompts for password input with masking and validation
func InputPassword(label string, validate func(string) error) string {
	var result string
	prompt := &survey.Password{
		Message: label,
	}
	wrappedValidator := func(ans interface{}) error {
		str, ok := ans.(string)
		if !ok {
			return fmt.Errorf("input must be a string")
		}
		return validate(str)
	}

	err := survey.AskOne(prompt, &result, survey.WithValidator(wrappedValidator))
	if err != nil {
		fmt.Println("Prompt failed:", err)
		return ""
	}
	return result
}

// SelectFromList prompts for selection from a list
func SelectFromList(label string, items []string) (int, string) {
	var result string
	prompt := &survey.Select{
		Message: label,
		Options: items,
	}
	err := survey.AskOne(prompt, &result)
	if err != nil {
		fmt.Println("Prompt failed:", err)
		return -1, ""
	}

	for i, item := range items {
		if item == result {
			return i, result
		}
	}
	return -1, result
}
