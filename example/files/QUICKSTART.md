# Quick Start Guide

This guide will help you get started with the Petstore Terraform Provider.

## Prerequisites

- Terraform 1.0 or later
- Go 1.21 or later (for building from source)
- Access to a Petstore API instance

## Installation

### Option 1: Build from Source

1. Clone the repository:
```bash
git clone https://github.com/yourusername/terraform-provider-petstore
cd terraform-provider-petstore
```

2. Build and install:
```bash
make build
make install
```

### Option 2: Manual Installation

1. Build the provider:
```bash
go build -o terraform-provider-petstore
```

2. Create the plugin directory:
```bash
mkdir -p ~/.terraform.d/plugins/registry.terraform.io/yourusername/petstore/0.1.0/$(go env GOOS)_$(go env GOARCH)
```

3. Copy the binary:
```bash
cp terraform-provider-petstore ~/.terraform.d/plugins/registry.terraform.io/yourusername/petstore/0.1.0/$(go env GOOS)_$(go env GOARCH)/
```

## Basic Usage

### 1. Create a Terraform Configuration

Create a file named `main.tf`:

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
  api_key  = "special-key"
}

resource "petstore_pet" "my_first_pet" {
  name = "Fluffy"
  photo_urls = [
    "https://example.com/fluffy.jpg"
  ]
  status = "available"

  category {
    name = "Cats"
  }
}

output "pet_id" {
  value = petstore_pet.my_first_pet.id
}
```

### 2. Initialize Terraform

```bash
terraform init
```

### 3. Plan the Changes

```bash
terraform plan
```

### 4. Apply the Configuration

```bash
terraform apply
```

Type `yes` when prompted to confirm.

### 5. View the Output

After successful apply, you'll see the pet ID:

```
Outputs:

pet_id = 12345
```

## Common Operations

### Creating a User

```hcl
resource "petstore_user" "john" {
  username   = "johndoe"
  first_name = "John"
  last_name  = "Doe"
  email      = "john@example.com"
}
```

### Creating an Order

```hcl
resource "petstore_order" "order1" {
  pet_id   = petstore_pet.my_first_pet.id
  quantity = 1
  status   = "placed"
}
```

### Using Data Sources

```hcl
data "petstore_pet" "existing" {
  id = 12345
}

output "existing_pet_name" {
  value = data.petstore_pet.existing.name
}
```

### Importing Existing Resources

```bash
# Import a pet
terraform import petstore_pet.my_pet 12345

# Import a user
terraform import petstore_user.my_user johndoe

# Import an order
terraform import petstore_order.my_order 67890
```

## Environment Variables

Instead of hardcoding credentials, use environment variables:

```bash
export PETSTORE_BASE_URL="https://petstore3.swagger.io/api/v3"
export PETSTORE_API_KEY="your-api-key"
```

Then in your Terraform configuration:

```hcl
provider "petstore" {
  # base_url and api_key will be read from environment variables
}
```

## Cleaning Up

To destroy all created resources:

```bash
terraform destroy
```

## Next Steps

- Check out the [complete example](examples/complete/main.tf) for a full demonstration
- Read the [README](README.md) for detailed documentation
- See individual examples in the `examples/` directory:
  - [Pet example](examples/pet/main.tf)
  - [User example](examples/user/main.tf)
  - [Order example](examples/order/main.tf)

## Troubleshooting

### Provider not found

If you get an error about the provider not being found, ensure:
1. The provider is built and installed correctly
2. The provider name in `required_providers` matches the installation path
3. Run `terraform init` again

### API errors

If you get API errors:
1. Check that your `base_url` is correct
2. Verify your `api_key` is valid
3. Ensure the Petstore API is accessible

### State issues

To inspect the current state:
```bash
terraform state list
terraform state show petstore_pet.my_first_pet
```

To refresh the state from the API:
```bash
terraform refresh
```
