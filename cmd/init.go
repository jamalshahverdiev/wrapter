package cmd

import (
	"fmt"
	"wrapter/utils"

	"github.com/spf13/cobra"
)

// Init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the Terraform backend",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Initializing Terraform backend...")
		if err := utils.InitializeBackend(cfg); err != nil {
			utils.LogErrorAndExit("Initialization failed", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
