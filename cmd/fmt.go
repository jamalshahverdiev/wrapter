package cmd

import (
	"fmt"
	"wrapter/utils"

	"github.com/spf13/cobra"
)

// Fmt command
var fmtCmd = &cobra.Command{
	Use:   "fmt",
	Short: "Format the Terraform code",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Formatting Terraform code...")
		if err := utils.FormatCode(cfg); err != nil {
			utils.LogErrorAndExit("Formatting failed", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(fmtCmd)
}
