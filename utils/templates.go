package utils

// Embedded content for provider.tf
const providerTfContent = `terraform {
  required_version = ">= 1.0.0"
  backend "s3" {
    skip_credentials_validation = true
    skip_metadata_api_check     = true
  }
}
`

// Embedded content for variables.tf
const variablesTfContent = `variable "SERVICES_TOKEN" {
  type = string
}
  
output "outputs" {
  value = module.common_modules
}
`
