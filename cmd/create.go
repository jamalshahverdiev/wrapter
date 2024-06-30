package cmd

import (
	"fmt"
	"wrapter/utils"

	"github.com/spf13/cobra"
)

// Create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Bootstrap the service",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Bootstrapping the service...")
		if err := utils.BootstrapService(cfg); err != nil {
			utils.LogErrorAndExit("Service bootstrap failed", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
}
