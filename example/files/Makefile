.PHONY: build install test clean fmt docs deps init

# Initialize and download dependencies
init:
	go mod download
	go mod tidy

# Build the provider (downloads deps first)
build: init
	go build -o terraform-provider-petstore

# Install the provider locally for testing
install: build
	mkdir -p ~/.terraform.d/plugins/registry.terraform.io/yourusername/petstore/0.1.0/$$(go env GOOS)_$$(go env GOARCH)
	cp terraform-provider-petstore ~/.terraform.d/plugins/registry.terraform.io/yourusername/petstore/0.1.0/$$(go env GOOS)_$$(go env GOARCH)/

# Run tests
test:
	go test ./... -v

# Clean build artifacts
clean:
	rm -f terraform-provider-petstore
	go clean

# Format code
fmt:
	go fmt ./...
	terraform fmt -recursive ./examples/

# Generate documentation
docs:
	go generate ./...

# Download dependencies only
deps:
	go mod download
	go mod tidy
	go mod verify
