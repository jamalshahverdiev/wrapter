package cmd

import (
	"fmt"
	"wrapter/utils"

	"github.com/spf13/cobra"
)

// Validate command
var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate the Terraform configuration",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Validating Terraform configuration...")
		if err := utils.ValidateConfiguration(cfg); err != nil {
			utils.LogErrorAndExit("Validation failed", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(validateCmd)
}
