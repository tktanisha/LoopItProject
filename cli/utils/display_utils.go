package utils

import (
	"fmt"
	"loopit/internal/config"
	"regexp"
	"strings"
)

// ShowBanner - Display the application welcome banner
func ShowBanner() {
	fmt.Println()
	fmt.Println(config.Green + strings.Repeat("═", 80) + config.Reset)
	fmt.Println(config.Green + "║" + centerText("", 78) + "║" + config.Reset)
	fmt.Println(config.Green + "║" + centerText(config.Bold+" WELCOME TO Loop IT CLI PROJECT", 78) + "║" + config.Reset)
	fmt.Println(config.Green + "║" + centerText("", 78) + "║" + config.Reset)
	fmt.Println(config.Green + strings.Repeat("═", 80) + config.Reset)
	fmt.Println()
}

// centerText - Helper to center text in given width
func centerText(text string, width int) string {
	padding := (width - len(stripAnsi(text))) / 2
	return strings.Repeat(" ", padding) + text + strings.Repeat(" ", width-padding-len(stripAnsi(text)))
}

// stripAnsi - Helper to strip ANSI escape codes for correct padding
func stripAnsi(str string) string {
	ansiEscape := `\x1b\[[0-9;]*m`
	re := regexp.MustCompile(ansiEscape)
	return re.ReplaceAllString(str, "")
}
