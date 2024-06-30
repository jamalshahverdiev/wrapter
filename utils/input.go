package utils

import (
	"fmt"
	"strings"
)

// GetUserInput prompts the user for a single string input
func GetUserInput(prompt string) (string, error) {
	// Simulating user input logic; replace with your actual input mechanism
	fmt.Print(prompt)
	var input string
	fmt.Scanln(&input) // Replace with actual input capture mechanism
	return input, nil
}

// Confirm prompts the user for a yes/no confirmation
func Confirm(prompt string) bool {
	input, _ := GetUserInput(fmt.Sprintf("%s [y/N]: ", prompt))
	return strings.ToLower(input) == "y"
}
