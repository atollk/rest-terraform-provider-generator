# Terraform API Provider Generator

⚠️ This project is in an unstable development phase.

A code generation tool that automatically creates Terraform providers from OpenAPI/Swagger specifications.
This eliminates the need to manually write boilerplate Terraform provider code for REST APIs.

Compared to the [Hashicorp OpenAPI generator](https://developer.hashicorp.com/terraform/plugin/code-generation/openapi-generator), this tool aims to provide a complete end-to-end experience, preventing the user from having to write any Go code themselves.

## Overview

The Terraform API Provider Generator takes two inputs:
1. An OpenAPI/Swagger specification (JSON format) that describes your REST API
2. A generation specification (YAML format) that defines how to map API resources to Terraform resources

It then generates a complete, working Terraform provider with:
- Resource implementations (Create, Read, Update, Delete operations)
- Data source implementations
- Proper schema definitions based on OpenAPI models
- HTTP client code for API interactions

## Requirements

- [Go](https://golang.org/doc/install) >= 1.21
- OpenAPI 3.0 or 3.1 specification file (JSON format)
- Generation specification file (YAML format)

## Installation

Clone the repository:

```bash
git clone https://github.com/yourusername/terraform-api-provider-generator
cd terraform-api-provider-generator
```

Download dependencies:

```bash
go mod download
```

## Usage

### 1. Prepare Your OpenAPI Specification

Ensure you have an OpenAPI 3.x specification file in JSON format. Example:

```json
{
  "openapi": "3.0.0",
  "info": {
    "title": "My API",
    "version": "1.0.0"
  },
  "paths": {
    "/pet": {
      "post": { ... },
      "get": { ... }
    }
  }
}
```

### 2. Create a Generation Specification

Create a YAML file (`genspec.yaml`) that defines how to generate your Terraform provider:

```yaml
$schema: "./internal/provider_spec/rest_api_provider_schema.json"

global_defaults:
  uri: "api.example.com"
  id_attribute: "id"
  create_method: POST
  destroy_method: DELETE
  update_method: PUT
  read_method: GET

resources:
  pet:
    path: "/pet"
    force_recreate: true
    update:
      path: "/pet"
    read:
      path: "/pet"
    destroy:
      path: "/pet"

  user:
    path: "/user"
    id_attribute: "username"
    force_new:
      - "username"
```

Look at `internal/provider_spec/rest_api_provider_schema.json` for detailed information about all options.

#### Generation Specification Options

### 3. Generate the Provider Code

Run the generator (modify paths in `internal/main.go` or create your own entry point):

```bash
go run internal/main.go
```

## Contributing

Contributions are welcome! Please feel free to submit issues or pull requests.
