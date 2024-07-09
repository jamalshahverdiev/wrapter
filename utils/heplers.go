package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"wrapter/common"
	"wrapter/config"

	"github.com/AlecAivazis/survey/v2"
)

// ConstructStateKey constructs the state key for the Terraform backend configuration
func ConstructStateKey(project, currentDir string, cfg *config.Config) (string, error) {
	// Clean the current directory path to ensure it's in a standard format
	cleanDir := filepath.Clean(currentDir)

	// Extract the environment from the path
	environment := ExtractEnvironmentFromPath(cleanDir)
	if environment == "" {
		return "", fmt.Errorf("unable to determine environment from path: %s", cleanDir)
	}

	fmt.Printf("Extracted environment: %s\n", environment) // Debugging statement

	// Dynamically get the account ID from the profiles using the environment string
	accountID, err := getFieldValueByEnvironment(cfg.Profiles, environment)
	if err != nil {
		return "", err
	}

	// Get the region from the default_regions based on the account ID
	region, exists := cfg.DefaultRegions[accountID]
	if !exists {
		return "", fmt.Errorf("region not found for account ID: %s", accountID)
	}

	// Construct the expected prefix (account ID, environment, region)
	expectedPrefix := filepath.Join(accountID, environment, region)

	fmt.Printf("Expected prefix: %s\n", expectedPrefix) // Debugging statement

	// Find the relative service path after the expected prefix
	startIdx := strings.Index(cleanDir, expectedPrefix)
	if startIdx == -1 {
		return "", fmt.Errorf("directory path does not match expected prefix: %s", expectedPrefix)
	}
	// The service path is everything after the expected prefix
	relativeServicePath := strings.TrimPrefix(cleanDir[startIdx+len(expectedPrefix):], string(filepath.Separator))

	// Construct the final state key
	stateKey := filepath.Join(project, expectedPrefix, relativeServicePath, "service.tfstate")

	fmt.Printf("Constructed state key: %s\n", stateKey) // Debugging statement

	return stateKey, nil
}

// getFieldValueByEnvironment uses reflection to get the value of the field in the Profiles struct based on the environment string
func getFieldValueByEnvironment(profiles interface{}, environment string) (string, error) {
	// Define a mapping from lowercase environment names to struct field names
	envToFieldName := map[string]string{
		"dev":    "Dev",
		"stable": "Stable",
		"prod":   "Prod",
		"mgmt":   "Mgmt",
	}

	fieldName, exists := envToFieldName[strings.ToLower(environment)]
	if !exists {
		return "", fmt.Errorf("no profile found for environment: %s", environment)
	}

	// fmt.Printf("Looking for environment: %s (converted to: %s)\n", environment, fieldName) // Debugging statement

	r := reflect.ValueOf(profiles)
	f := reflect.Indirect(r).FieldByName(fieldName)
	if !f.IsValid() {
		return "", fmt.Errorf("no profile found for environment: %s", environment)
	}
	return f.String(), nil
}

// ExtractAccountIDAndRegionFromPath extracts the account ID and region from the directory path
func ExtractAccountIDAndRegionFromPath(path string) (string, string, error) {
	// Split the path into segments
	segments := strings.Split(filepath.Clean(path), string(filepath.Separator))

	if len(segments) < 6 {
		return "", "", fmt.Errorf("insufficient directory segments to determine account ID and region")
	}

	// Assume the account ID is the 5th segment and region is the 6th segment
	accountID := segments[4]
	region := segments[5]

	return accountID, region, nil
}

// ExtractEnvironmentFromPath extracts the environment from the directory path
func ExtractEnvironmentFromPath(path string) string {
	// Split the path into segments
	segments := strings.Split(filepath.Clean(path), string(filepath.Separator))

	// Assuming the structure /home/freshly/gitlab/diftech/scripts/111111111/prod/us-east-1/dwh/airbyte
	// We take the segment just before the region (e.g., 'prod' in this case)
	// Adjust the index based on your actual directory depth
	for i, segment := range segments {
		if isRegion(segment) {
			if i > 0 {
				return segments[i-1]
			}
		}
	}
	return "" // Return an empty string if the environment can't be determined
}

// Helper function to identify if a string is a region code
func isRegion(segment string) bool {
	regions := []string{
		"us-east-1", "us-west-2", "eu-central-1", // Add other regions as needed
	}
	for _, region := range regions {
		if segment == region {
			return true
		}
	}
	return false
}

// ListDirs lists all directories containing Terraform files, excluding .terraform directories and their subdirectories
func ListDirs() ([]string, error) {
	var dirs []string
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Exclude any directory that is part of a `.terraform` path
		if info.IsDir() && strings.Contains(path, ".terraform") {
			return filepath.SkipDir
		}

		// Add the directory if it passes the above checks
		if info.IsDir() {
			dirs = append(dirs, path)
		}
		return nil
	})
	return dirs, err
}

// Helper functions for prompts using the survey library
// promptForEnvironments prompts the user to select environments
func promptForEnvironments() ([]string, error) {
	options := []string{"dev", "stable", "prod"}
	var selected []string
	prompt := &survey.MultiSelect{
		Message: "Choose environments:",
		Options: options,
	}
	err := survey.AskOne(prompt, &selected)
	return selected, err
}

// promptForServiceType prompts the user to choose between new or custom service
func promptForServiceType() (string, error) {
	options := []string{"new", "custom"}
	var selected string
	prompt := &survey.Select{
		Message: "New or custom service?:",
		Options: options,
	}
	err := survey.AskOne(prompt, &selected)
	return selected, err
}

// promptForComponents prompts the user to select the required components
func promptForComponents(services []string) ([]string, error) {
	var selected []string
	prompt := &survey.MultiSelect{
		Message: "What components do you need?:",
		Options: services,
	}
	err := survey.AskOne(prompt, &selected)
	return selected, err
}

// promptForText prompts the user for a single line of text input
func promptForText(message string) (string, error) {
	var input string
	prompt := &survey.Input{
		Message: message,
	}
	err := survey.AskOne(prompt, &input)
	return input, err
}

// createServiceFiles creates the directory structure and required files for the service
func createServiceFiles(cfg *config.Config, environment, teamName, serviceName string, components []string) error {
	accountID, region, err := getAccountAndRegion(cfg, environment)
	if err != nil {
		return fmt.Errorf("could not get account and region: %w", err)
	}

	// Determine the root directory for git
	gitRoot, err := common.FindGitRoot()
	if err != nil {
		return fmt.Errorf("could not find git root: %w", err)
	}

	// Create the target directory structure
	targetDir := filepath.Join(gitRoot, accountID, environment, region, teamName, serviceName)
	if err := CreateTargetDir(targetDir); err != nil {
		return fmt.Errorf("could not create target directory %s: %w", targetDir, err)
	}

	// Generate and write locals.tf
	if err := generateLocalsTF(targetDir, environment, teamName, serviceName, accountID, region); err != nil {
		return fmt.Errorf("could not create locals.tf: %w", err)
	}

	// Generate and write main.tf
	if err := generateMainTF(cfg, targetDir, components); err != nil {
		return fmt.Errorf("could not create main.tf: %w", err)
	}

	// Copy static provider.tf
	if err := copyStaticFile("provider.tf", targetDir); err != nil {
		return fmt.Errorf("could not copy provider.tf: %w", err)
	}

	// Copy static variables.tf
	if err := copyStaticFile("variables.tf", targetDir); err != nil {
		return fmt.Errorf("could not copy variables.tf: %w", err)
	}

	return nil
}

// getAccountAndRegion retrieves the account ID and region for a given environment
func getAccountAndRegion(cfg *config.Config, environment string) (string, string, error) {
	var accountID string
	switch environment {
	case "dev", "stable":
		accountID = cfg.Profiles.Dev
	case "prod":
		accountID = cfg.Profiles.Prod
	case "mgmt":
		accountID = cfg.Profiles.Mgmt
	default:
		return "", "", fmt.Errorf("unknown environment: %s", environment)
	}

	region, exists := cfg.DefaultRegions[accountID]
	if !exists {
		return "", "", fmt.Errorf("region not found for account ID: %s", accountID)
	}

	return accountID, region, nil
}

// generateLocalsTF generates the locals.tf content and writes it to the target directory
func generateLocalsTF(targetDir, environment, teamName, serviceName, accountID, region string) error {
	localsContent := fmt.Sprintf(`locals {
  env                  = "%s"
  region               = "%s"
  team                 = "%s"
  service_name         = "%s"
  account_id           = "%s"
  ad_ou_dn             = "OU=company,DC=example,DC=local"
  scope                = "global"
  category             = "security"
  path                 = "iac"
  infra_path           = "infra"
  pg_path              = "postgresql"
  mongo_path           = "mongodb"
  keycloak_path        = "keycloak"
  vault_mount_path     = "vss1"
  services_secret_path = "services"
  secret               = "secret"
  callback_url         = "http://localhost:8080/callback"
}`, environment, region, teamName, serviceName, accountID)

	return WriteFile(filepath.Join(targetDir, "locals.tf"), localsContent)
}

// generateMainTF generates the main.tf content and writes it to the target directory
func generateMainTF(cfg *config.Config, targetDir string, components []string) error {
	// Map of component flags
	componentFlags := map[string]bool{
		"postgres_enabled": false,
		"mongodb_enabled":  false,
		"keycloak_enabled": false,
		"ad_enabled":       false,
	}

	// Set flags based on selected components
	for _, component := range components {
		switch component {
		case "PostgreSQL":
			componentFlags["postgres_enabled"] = true
		case "Mongo":
			componentFlags["mongodb_enabled"] = true
		case "Keycloak":
			componentFlags["keycloak_enabled"] = true
		case "AD":
			componentFlags["ad_enabled"] = true
		}
	}

	// Use module_git_url and module_version from the configuration
	moduleSource := fmt.Sprintf("%s?ref=%s", cfg.CommonService.ModuleGitURL, cfg.CommonService.ModuleVersion)

	// Build the main.tf content
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`module "common_modules" {
  source               = "%s"
  env                  = local.env
  team                 = local.team
  service_name         = local.service_name
  scope                = local.scope
  category             = local.category
  ad_ou_dn             = local.ad_ou_dn
  callback_url         = local.callback_url
  vault_mount_path     = local.vault_mount_path
  region               = local.region
  infra_path           = local.infra_path
  account_id           = local.account_id
  path                 = local.path
  pg_path              = local.pg_path
  mongo_path           = local.mongo_path
  keycloak_path        = local.keycloak_path
  services_secret_path = local.services_secret_path
  secret               = local.secret
  SERVICES_TOKEN       = var.SERVICES_TOKEN
`, moduleSource))

	// Conditionally add the component flags if they are enabled
	if componentFlags["postgres_enabled"] {
		sb.WriteString(fmt.Sprintf("  postgres_enabled     = %t\n", componentFlags["postgres_enabled"]))
	}
	if componentFlags["mongodb_enabled"] {
		sb.WriteString(fmt.Sprintf("  mongodb_enabled      = %t\n", componentFlags["mongodb_enabled"]))
	}
	if componentFlags["keycloak_enabled"] {
		sb.WriteString(fmt.Sprintf("  keycloak_enabled     = %t\n", componentFlags["keycloak_enabled"]))
	}
	if componentFlags["ad_enabled"] {
		sb.WriteString(fmt.Sprintf("  ad_enabled           = %t\n", componentFlags["ad_enabled"]))
	}

	// Close the module block
	sb.WriteString("}\n")

	// Convert the string builder content to a string and write it to the file
	mainContent := sb.String()
	return WriteFile(filepath.Join(targetDir, "main.tf"), mainContent)
}

// generateCustomSettingsTF generates the settings.tf content and writes it to the custom service directory
func generateCustomSettingsTF(targetDir, accountID string) error {
	settingsContent := fmt.Sprintf(`locals {
  service_name         = data.terraform_remote_state.wrapter.outputs.service_outputs.settings.name
  region               = data.terraform_remote_state.wrapter.outputs.service_outputs.settings.region
  team                 = data.terraform_remote_state.wrapter.outputs.service_outputs.settings.team
  env                  = data.terraform_remote_state.wrapter.outputs.service_outputs.settings.env
  account_id           = "%s"
  ad_ou_dn             = "OU=company,DC=example,DC=local"
  scope                = "global"
  category             = "security"
  path                 = "iac"
  infra_path           = "infra"
  pg_path              = "postgresql"
  mongo_path           = "mongodb"
  keycloak_path        = "keycloak"
  vault_mount_path     = "vss1"
  services_secret_path = "services"
  secret               = "secret"
  callback_url         = "http://localhost:8080/callback"
}

variable "MINIO_ACCESS_KEY" {
  type = string
}
variable "MINIO_SECRET_KEY" {
  type = string
}`, accountID)

	return WriteFile(filepath.Join(targetDir, "settings.tf"), settingsContent)
}

// generateTFStateTF generates the tfstate.tf content and writes it to the custom service directory
func generateTFStateTF(targetDir string, cfg *config.Config, environment string) error {
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

	// Get the region dynamically based on the current environment and configuration
	accountID, err := getFieldValueByEnvironment(cfg.Profiles, environment)
	if err != nil {
		return err
	}
	region := cfg.DefaultRegions[accountID]

	// Fetch the endpoint value from the configuration
	endpoint := cfg.Environments.Endpoint

	tfstateContent := fmt.Sprintf(`data "terraform_remote_state" "wrapter" {
  backend = "s3"
  config = {
    endpoint                    = "%s"
    bucket                      = "%s-tfstates"
    key                         = "%s"
    region                      = "%s"
    access_key                  = var.MINIO_ACCESS_KEY
    secret_key                  = var.MINIO_SECRET_KEY
    skip_credentials_validation = true
    skip_metadata_api_check     = true
    skip_requesting_account_id  = true
  }
}`, endpoint, cfg.Tofu.Project, stateKey, region)

	return WriteFile(filepath.Join(targetDir, "tfstate.tf"), tfstateContent)
}

// copyStaticFile writes the embedded content of a static file to the target directory
func copyStaticFile(filename, targetDir string) error {
	var content string

	switch filename {
	case "provider.tf":
		content = providerTfContent
	case "variables.tf":
		content = variablesTfContent
	default:
		return fmt.Errorf("unsupported static file: %s", filename)
	}

	targetPath := filepath.Join(targetDir, filename)

	return WriteFile(targetPath, content)
}

// createCustomServiceFiles creates the directory structure and required files for the custom service
func createCustomServiceFiles(cfg *config.Config, environment, teamName, serviceName string) error {
	accountID, region, err := getAccountAndRegion(cfg, environment)
	if err != nil {
		return fmt.Errorf("could not get account and region: %w", err)
	}

	// Determine the root directory for git
	gitRoot, err := common.FindGitRoot()
	if err != nil {
		return fmt.Errorf("could not find git root: %w", err)
	}

	// Create the target directory structure
	teamDir := filepath.Join(gitRoot, accountID, environment, region, teamName)
	serviceDir := filepath.Join(teamDir, serviceName)
	customServiceDir := serviceDir + "-custom"

	// Check if team directory exists
	if _, err := os.Stat(teamDir); os.IsNotExist(err) {
		fmt.Println("Entered team name doesn't exist. To create a new service please execute 'wrapter create' and then select 'new'.")
		return nil
	}

	// Check if service directory exists
	if _, err := os.Stat(serviceDir); os.IsNotExist(err) {
		fmt.Println("Entered service name doesn't exist. To create a new service please execute 'wrapter create' and then select 'new'.")
		return nil
	}

	// Create the custom service directory
	if err := CreateTargetDir(customServiceDir); err != nil {
		return fmt.Errorf("could not create custom service directory %s: %w", customServiceDir, err)
	}

	// Generate and write settings.tf
	if err := generateCustomSettingsTF(customServiceDir, accountID); err != nil {
		return fmt.Errorf("could not create settings.tf: %w", err)
	}

	// Generate and write tfstate.tf
	if err := generateTFStateTF(customServiceDir, cfg, environment); err != nil {
		return fmt.Errorf("could not create tfstate.tf: %w", err)
	}

	// Copy static provider.tf
	if err := copyStaticFile("provider.tf", customServiceDir); err != nil {
		return fmt.Errorf("could not copy provider.tf: %w", err)
	}

	return nil
}