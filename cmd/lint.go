package cmd

import (
	"fmt"
	"wrapter/utils"

	"github.com/spf13/cobra"
)

// Lint command
var lintCmd = &cobra.Command{
	Use:   "lint",
	Short: "Run the linter",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running linter...")
		if err := utils.RunLinter(cfg); err != nil {
			utils.LogErrorAndExit("Linter failed", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(lintCmd)
}
