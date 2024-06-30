package cmd

import (
	"fmt"
	"wrapter/utils"

	"github.com/spf13/cobra"
)

// Doc command
var docCmd = &cobra.Command{
	Use:   "doc",
	Short: "Generate documentation",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Generating documentation...")
		if err := utils.GenerateDocs(cfg); err != nil {
			utils.LogErrorAndExit("Documentation generation failed", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(docCmd)
}
