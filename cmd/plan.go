package cmd

import (
	"fmt"
	"wrapter/utils"

	"github.com/spf13/cobra"
)

// Plan command
var planCmd = &cobra.Command{
	Use:   "plan",
	Short: "Generate a Terraform plan",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Generating Terraform plan...")
		if err := utils.Plan(cfg); err != nil {
			utils.LogErrorAndExit("Plan generation failed", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(planCmd)
}
