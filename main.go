package main

import (
	"log"
	"wrapter/cmd"
	"wrapter/utils"
)

func main() {
	// Verify that all required binaries are available
	if err := utils.VerifyRequirements(); err != nil {
		log.Fatalf("Verification failed: %v", err)
	}

	// Execute the command
	if err := cmd.Execute(); err != nil {
		log.Fatalf("error executing command: %v", err)
	}
}
