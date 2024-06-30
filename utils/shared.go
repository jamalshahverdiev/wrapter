package utils

import (
	"log"
	"os"
)

// WarnValidation returns the environment variable for validation warning
func WarnValidation() bool {
	return os.Getenv("ALLOW_FAIL_VALIDATION") != ""
}

// LogErrorAndExit logs an error message and exits the application
func LogErrorAndExit(msg string, err error) {
	log.Fatalf("%s: %v", msg, err)
}

// CreateTargetDir ensures that the target directory exists
func CreateTargetDir(path string) error {
	return os.MkdirAll(path, os.ModePerm)
}

// WriteFile writes content to a file at the specified path
func WriteFile(filePath, content string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	return err
}
