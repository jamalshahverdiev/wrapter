package cmd

import (
	"fmt"
	"wrapter/utils"

	"github.com/spf13/cobra"
)

// Lock command
var lockCmd = &cobra.Command{
	Use:   "lock",
	Short: "Set providers lock",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Setting providers lock...")
		if err := utils.LockProviders(cfg); err != nil {
			utils.LogErrorAndExit("Locking providers failed", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(lockCmd)
}
