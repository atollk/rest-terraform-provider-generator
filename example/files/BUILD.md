# Build Instructions

This Terraform provider is written in Go and requires Go 1.21 or later to build.

## Prerequisites

1. Install Go 1.21 or later from https://golang.org/dl/
2. Ensure `$GOPATH/bin` is in your PATH

## Building the Provider

### Step 1: Download Dependencies

Before building, you need to download all Go module dependencies:

```bash
cd terraform-provider-petstore
go mod download
go mod tidy
```

This will:
- Download all required dependencies
- Update `go.mod` and `go.sum` with the correct versions
- Verify module checksums

### Step 2: Build

Once dependencies are downloaded, build the provider:

```bash
# Using Make
make build

# Or manually
go build -o terraform-provider-petstore
```

### Step 3: Install Locally

To use the provider with Terraform, install it to your local plugin directory:

```bash
# Using Make
make install

# Or manually (Linux/macOS)
mkdir -p ~/.terraform.d/plugins/registry.terraform.io/yourusername/petstore/0.1.0/$(go env GOOS)_$(go env GOARCH)
cp terraform-provider-petstore ~/.terraform.d/plugins/registry.terraform.io/yourusername/petstore/0.1.0/$(go env GOOS)_$(go env GOARCH)/

# For Windows
mkdir -p %APPDATA%\terraform.d\plugins\registry.terraform.io\yourusername\petstore\0.1.0\windows_amd64
copy terraform-provider-petstore.exe %APPDATA%\terraform.d\plugins\registry.terraform.io\yourusername\petstore\0.1.0\windows_amd64\
```

## Verifying the Build

After installation, create a test Terraform configuration:

```hcl
terraform {
  required_providers {
    petstore = {
      source = "yourusername/petstore"
    }
  }
}

provider "petstore" {
  base_url = "https://petstore3.swagger.io/api/v3"
}
```

Then run:

```bash
terraform init
```

You should see output indicating the provider was successfully initialized.

## Development Build

For development, you can use a local build without installing:

```bash
# Build with debug info
go build -gcflags="all=-N -l" -o terraform-provider-petstore

# Run Terraform with the local provider
terraform plan
```

## Troubleshooting

### Missing Dependencies

If you get errors about missing dependencies:

```bash
go mod download
go mod tidy
go mod verify
```

### Version Conflicts

If you have version conflicts:

```bash
go clean -modcache
go mod download
```

### Building for Different Platforms

To build for a different operating system or architecture:

```bash
# For Linux AMD64
GOOS=linux GOARCH=amd64 go build -o terraform-provider-petstore

# For macOS ARM64 (M1/M2)
GOOS=darwin GOARCH=arm64 go build -o terraform-provider-petstore

# For Windows AMD64
GOOS=windows GOARCH=amd64 go build -o terraform-provider-petstore.exe
```

## CI/CD Integration

For automated builds, add this to your CI/CD pipeline:

```bash
#!/bin/bash
set -e

# Download dependencies
go mod download

# Run tests
go test ./...

# Build
go build -o terraform-provider-petstore

# Optional: Run acceptance tests
TF_ACC=1 go test ./... -v -timeout 120m
```

## Next Steps

After building:
1. Review the [Quick Start Guide](QUICKSTART.md)
2. Try the [example configurations](examples/)
3. Read the [README](README.md) for usage details
