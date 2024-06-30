package utils

import (
	"fmt"
	"os"
	"os/exec"
)

// CheckBinary verifies if a binary exists in the system's PATH and is non-empty
func CheckBinary(binaryName string) error {
	path, err := exec.LookPath(binaryName)
	if err != nil {
		return fmt.Errorf("binary %s not found in PATH", binaryName)
	}

	fileInfo, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("unable to stat the binary %s: %w", binaryName, err)
	}

	if fileInfo.Size() == 0 {
		return fmt.Errorf("binary %s is found but it is an empty file", binaryName)
	}

	return nil
}

// VerifyRequirements checks for all required binaries
func VerifyRequirements() error {
	requiredBinaries := []string{"tofu", "jq", "terraform-docs"}

	for _, binary := range requiredBinaries {
		if err := CheckBinary(binary); err != nil {
			return fmt.Errorf("requirement check failed: %w", err)
		}
	}

	return nil
}
