package common

import (
	"os"
	"path/filepath"
)

// FindGitRoot finds the root directory of the Git repository and returns its path.
// It starts from the current directory and moves upwards.
func FindGitRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, ".git")); err == nil {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", err
		}
		dir = parent
	}
}
