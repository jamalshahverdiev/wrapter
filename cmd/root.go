package cmd

import (
	"wrapter/config"
	"wrapter/utils"

	"github.com/spf13/cobra"
)

var cfg *config.Config

// Root command
var rootCmd = &cobra.Command{
	Use:   "wrapter",
	Short: "Wrapter - A Terraform wrapper in Go",
	Long:  `Wrapter is a CLI tool to manage Terraform codes for the Microservices requirements.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	var err error
	cfg, err = config.LoadConfig("invoke.yaml")
	if err != nil {
		utils.LogErrorAndExit("Failed to load config", err)
	}
}
