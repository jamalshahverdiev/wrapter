package utils

import (
	"fmt"
	"os"
	"os/exec"
	"wrapter/config"
)

// InitializeBackend initializes the Terraform backend
func InitializeBackend(cfg *config.Config) error {
	endpoint := cfg.Environments.Endpoint

	// Determine the environment based on the current directory
	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}

	// Construct the state key dynamically using the project name, current directory, and the configuration
	stateKey, err := ConstructStateKey(cfg.Tofu.Project, currentDir, cfg)
	if err != nil {
		return err
	}

	// fmt.Println("State key:", stateKey)
	// Get the region dynamically based on the current environment and configuration
	environment := ExtractEnvironmentFromPath(currentDir)
	accountID, err := getFieldValueByEnvironment(cfg.Profiles, environment)
	if err != nil {
		return err
	}
	region := cfg.DefaultRegions[accountID]
	fmt.Println("Region:", region)

	// Run the `tofu init` command only once
	command := exec.Command("tofu", "init",
		"-backend-config=endpoint="+endpoint,
		"-backend-config=bucket="+cfg.Tofu.Project+"-tfstates",
		"-backend-config=region="+region, // Use the dynamically fetched region
		"-backend-config=key="+stateKey,
		"-backend-config=access_key="+os.Getenv("MINIO_ACCESS_KEY"),
		"-backend-config=secret_key="+os.Getenv("MINIO_SECRET_KEY"),
		"-reconfigure",
	)

	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	command.Env = append(os.Environ(), "TF_CLI_CONFIG_FILE="+cfg.TerraformCliConfigPath)

	if err := command.Run(); err != nil {
		return err
	}

	return nil
}

// RunLinter runs the linter for Terraform configurations
func RunLinter(cfg *config.Config) error {
	dirs, err := ListDirs()
	if err != nil {
		return err
	}

	for _, dir := range dirs {
		println("Running TFlint for:", dir)
		command := exec.Command("sh", "-c", "tflint && tofu fmt -diff -check=true || (echo 'Run tofu fmt before commit!!!'; exit 1)")
		command.Dir = dir
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr
		command.Env = append(os.Environ(), "TF_CLI_CONFIG_FILE="+cfg.TerraformCliConfigPath)

		if err := command.Run(); err != nil {
			return err
		}
	}

	return nil
}

// GenerateDocs generates documentation for Terraform configurations
func GenerateDocs(cfg *config.Config) error {
	dirs, err := ListDirs()
	if err != nil {
		return err
	}

	for _, dir := range dirs {
		command := exec.Command("terraform-docs", "markdown", "table", "--output-file", "README.md", "--output-mode", "inject", dir)
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr
		command.Env = append(os.Environ(), "TF_CLI_CONFIG_FILE="+cfg.TerraformCliConfigPath)

		if err := command.Run(); err != nil {
			return err
		}
	}

	return nil
}

// FormatCode formats the Terraform code
func FormatCode(cfg *config.Config) error {
	dirs, err := ListDirs()
	if err != nil {
		return err
	}

	for _, dir := range dirs {
		println("Running tofu fmt in the", dir)
		command := exec.Command("tofu", "fmt")
		command.Dir = dir
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr
		command.Env = append(os.Environ(), "TF_CLI_CONFIG_FILE="+cfg.TerraformCliConfigPath)

		if err := command.Run(); err != nil {
			return err
		}
	}

	return nil
}

// LockProviders locks Terraform providers
func LockProviders(cfg *config.Config) error {
	dirs, err := ListDirs()
	if err != nil {
		return err
	}

	for _, dir := range dirs {
		println("Running tofu lock in the", dir)
		command := exec.Command("tofu", "providers", "lock",
			"-platform=linux_amd64",
			"-platform=darwin_amd64",
			"-platform=darwin_arm64",
		)
		command.Dir = dir
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr
		command.Env = append(os.Environ(), "TF_CLI_CONFIG_FILE="+cfg.TerraformCliConfigPath)

		if err := command.Run(); err != nil {
			return err
		}
	}

	return nil
}

// ValidateConfiguration validates the Terraform configuration
func ValidateConfiguration(cfg *config.Config) error {
	dirs, err := ListDirs()
	if err != nil {
		return err
	}

	for _, dir := range dirs {
		println("Running tofu validate in the", dir)
		command := exec.Command("sh", "-c",
			"tofu init -input=false -backend=false && tofu validate && tflint",
		)
		command.Dir = dir
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr
		command.Env = append(os.Environ(), "TF_CLI_CONFIG_FILE="+cfg.TerraformCliConfigPath)

		if err := command.Run(); err != nil {
			return err
		}
	}

	return nil
}

// Plan generates a Terraform plan
func Plan(cfg *config.Config) error {
	// Get the current working directory
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("could not get current directory: %w", err)
	}

	// Construct the state key dynamically using the project name, current directory, and the configuration
	stateKey, err := ConstructStateKey(cfg.Tofu.Project, currentDir, cfg)
	if err != nil {
		return fmt.Errorf("could not construct state key: %w", err)
	}

	// Extract the environment from the current directory
	environment := ExtractEnvironmentFromPath(currentDir)
	if environment == "" {
		return fmt.Errorf("could not determine environment from path: %s", currentDir)
	}

	// Get the account ID for the current environment from the profiles
	accountID, err := getFieldValueByEnvironment(cfg.Profiles, environment)
	if err != nil {
		return fmt.Errorf("could not get account ID for environment %s: %w", environment, err)
	}

	// Get the region for the account ID from the default regions map
	region, exists := cfg.DefaultRegions[accountID]
	if !exists {
		return fmt.Errorf("region not found for account ID: %s", accountID)
	}

	// Build the backend configuration command
	commandStr := fmt.Sprintf(
		"tofu init -backend-config=endpoint=%s -backend-config=bucket=%s-tfstates -backend-config=region=%s -backend-config=key=%s -backend-config=access_key=%s -backend-config=secret_key=%s -reconfigure && tofu plan -out tfplan.bin && tofu show --json tfplan.bin | jq > tfplan.json",
		cfg.Environments.Endpoint,
		cfg.Tofu.Project,
		region,
		stateKey,
		os.Getenv("MINIO_ACCESS_KEY"),
		os.Getenv("MINIO_SECRET_KEY"),
	)

	// Execute the command
	command := exec.Command("sh", "-c", commandStr)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	command.Env = append(os.Environ(), "TF_CLI_CONFIG_FILE="+cfg.TerraformCliConfigPath)

	if err := command.Run(); err != nil {
		return fmt.Errorf("plan generation failed: %w", err)
	}

	return nil
}

// BootstrapService bootstraps the service
func BootstrapService(cfg *config.Config) error {
	// Prompt for user inputs using the survey library
	environments, err := promptForEnvironments()
	if err != nil {
		return err
	}

	serviceType, err := promptForServiceType()
	if err != nil {
		return err
	}

	if serviceType != "new" {
		fmt.Println("Custom service handling is not implemented yet.")
		return nil
	}

	components, err := promptForComponents(cfg.CommonService.Services)
	if err != nil {
		return err
	}

	teamName, err := promptForText("Provide your team name and enter:")
	if err != nil {
		return err
	}

	serviceName, err := promptForText("Provide the service name and enter:")
	if err != nil {
		return err
	}

	// Create directories and files for each selected environment
	for _, env := range environments {
		if err := createServiceFiles(cfg, env, teamName, serviceName, components); err != nil {
			return fmt.Errorf("failed to create service files for environment %s: %w", env, err)
		}
	}

	return nil
}
