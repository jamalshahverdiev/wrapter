package config

import (
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"wrapter/common"
)

// Config represents the configuration structure
type Config struct {
	Tofu struct {
		Version string `yaml:"version"`
		Project string `yaml:"project"`
		Region  string `yaml:"region"`
	} `yaml:"tofu"`
	CommonService struct {
		BootstrapURL  string   `yaml:"bootstrap_url"`
		ModuleGitURL  string   `yaml:"module_git_url"`
		ModuleVersion string   `yaml:"module_version"`
		Services      []string `yaml:"services"`
	} `yaml:"common_service"`
	DefaultRegions map[string]string `yaml:"default_regions"`
	Profiles       struct {
		Dev    string `yaml:"dev"`
		Stable string `yaml:"stable"`
		Prod   string `yaml:"prod"`
		Mgmt   string `yaml:"mgmt"`
	} `yaml:"profiles"`
	Environments struct {
		Endpoint string             `yaml:"endpoint"`
		Dev      EnvironmentDetails `yaml:"dev"`
		Stable   EnvironmentDetails `yaml:"stable"`
		Prod     EnvironmentDetails `yaml:"prod"`
		Mgmt     EnvironmentDetails `yaml:"mgmt"`
	} `yaml:"environments"`

	TerraformCliConfigPath string `yaml:"-"` // Path to the terraform.tfrc file
}

// EnvironmentDetails captures details for each environment configuration
type EnvironmentDetails struct {
	AWS    string `yaml:"aws"`
	EKS    string `yaml:"eks"`
	Aurora string `yaml:"aurora"`
	Redis  struct {
		Host  string `yaml:"host"`
		Group string `yaml:"group"`
	} `yaml:"redis"`
	Atlas struct {
		Cluster string `yaml:"cluster"`
		Project string `yaml:"project"`
	} `yaml:"atlas"`
}

// LoadConfig searches for the invoke.yaml file from the git root directory and loads it
func LoadConfig(filename string) (*Config, error) {
	gitRoot, err := common.FindGitRoot()
	if err != nil {
		return nil, err
	}

	configPath := filepath.Join(gitRoot, filename)

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	// Set the path to the terraform.tfrc file found in the same directory as the .git folder
	config.TerraformCliConfigPath = filepath.Join(gitRoot, "terraform.tfrc")

	return &config, nil
}
