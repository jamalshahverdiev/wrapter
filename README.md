# Wrapter - A Terraform Wrapper in Go

Wrapter is a CLI tool to manage Terraform codes for the Microservices requirements. It simplifies and automates various Terraform tasks, including initialization, validation, formatting, and more.

## Features

- **Initialization**: Initialize the Terraform backend.
- **Documentation**: Generate documentation for Terraform configurations.
- **Validation**: Validate the Terraform configuration.
- **Linting**: Run the linter for Terraform configurations.
- **Formatting**: Format the Terraform code.
- **Lock Providers**: Set providers lock.
- **Bootstrap Service**: Bootstrap new or custom services.
- **Plan Generation**: Generate a Terraform plan.

## Installation

To install Wrapter, you need to have Go installed on your machine. Then, you can use the following command:

```bash
go get github.com/jamalshahverdiev/wrapter
```

### Compile and move binary file to the PATH dir. We can use [!GVM](https://github.com/moovweb/gvm) to manage Go versions.

```bash
gvm use go1.22.3
```