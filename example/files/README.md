# Terraform Provider for Petstore API

This Terraform provider allows you to manage resources in the Swagger Petstore API (OpenAPI 3.0 specification).

## Features

- **Resources**: Create, read, update, and delete pets, users, and store orders
- **Data Sources**: Query existing pets, users, and orders
- **Type-safe**: Built with the Terraform Plugin Framework for reliability

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.21 (for development)

## Building the Provider

1. Clone the repository:
```bash
git clone https://github.com/yourusername/terraform-provider-petstore
cd terraform-provider-petstore
```

2. Download dependencies:
```bash
go mod download
go mod tidy
```

3. Build the provider:
```bash
go build -o terraform-provider-petstore
```

Or use the Makefile (which handles dependencies automatically):
```bash
make build
```

4. Install the provider locally:
```bash
mkdir -p ~/.terraform.d/plugins/registry.terraform.io/yourusername/petstore/0.1.0/linux_amd64
cp terraform-provider-petstore ~/.terraform.d/plugins/registry.terraform.io/yourusername/petstore/0.1.0/linux_amd64/
```

## Using the Provider

### Provider Configuration

```hcl
terraform {
  required_providers {
    petstore = {
      source = "yourusername/petstore"
      version = "~> 0.1"
    }
  }
}

provider "petstore" {
  base_url = "https://petstore3.swagger.io/api/v3"
  api_key  = "your-api-key-here"
}
```

You can also use environment variables:
- `PETSTORE_BASE_URL`: The base URL for the API
- `PETSTORE_API_KEY`: The API key for authentication

### Resources

#### Pet Resource

```hcl
resource "petstore_pet" "my_dog" {
  name = "Buddy"
  photo_urls = [
    "https://example.com/photo1.jpg"
  ]
  status = "available"

  category {
    name = "Dogs"
  }

  tags {
    name = "friendly"
  }
}
```

**Attributes:**
- `name` (Required): Pet name
- `photo_urls` (Required): List of photo URLs
- `status` (Optional): Pet status (available, pending, sold)
- `category` (Optional): Pet category with `id` and `name`
- `tags` (Optional): List of tags with `id` and `name`
- `id` (Computed): Pet ID assigned by the API

#### User Resource

```hcl
resource "petstore_user" "john" {
  username    = "johndoe"
  first_name  = "John"
  last_name   = "Doe"
  email       = "john@example.com"
  password    = "secret123"
  phone       = "+1234567890"
  user_status = 1
}
```

**Attributes:**
- `username` (Required): Username (cannot be changed after creation)
- `first_name` (Optional): First name
- `last_name` (Optional): Last name
- `email` (Optional): Email address
- `password` (Optional, Sensitive): Password
- `phone` (Optional): Phone number
- `user_status` (Optional): User status
- `id` (Computed): User ID assigned by the API

#### Order Resource

```hcl
resource "petstore_order" "order1" {
  pet_id    = petstore_pet.my_dog.id
  quantity  = 1
  ship_date = "2024-12-25T10:00:00Z"
  status    = "placed"
  complete  = false
}
```

**Attributes:**
- `pet_id` (Required): ID of the pet being ordered
- `quantity` (Optional): Quantity ordered
- `ship_date` (Optional): Ship date in RFC3339 format
- `status` (Optional): Order status (placed, approved, delivered)
- `complete` (Optional): Whether order is complete
- `id` (Computed): Order ID assigned by the API

**Note:** The Petstore API doesn't support updating orders, so any changes will require replacing the resource.

### Data Sources

#### Pet Data Source

```hcl
data "petstore_pet" "existing_pet" {
  id = 12345
}

output "pet_name" {
  value = data.petstore_pet.existing_pet.name
}
```

#### User Data Source

```hcl
data "petstore_user" "existing_user" {
  username = "johndoe"
}

output "user_email" {
  value = data.petstore_user.existing_user.email
}
```

#### Order Data Source

```hcl
data "petstore_order" "existing_order" {
  id = 54321
}

output "order_status" {
  value = data.petstore_order.existing_order.status
}
```

## Importing Resources

You can import existing resources:

```bash
# Import a pet by ID
terraform import petstore_pet.my_dog 12345

# Import a user by username
terraform import petstore_user.john johndoe

# Import an order by ID
terraform import petstore_order.order1 54321
```

## Development

### Running Tests

```bash
go test ./...
```

### Generating Documentation

```bash
go generate ./...
```

This will generate documentation in the `docs/` directory.

## License

Apache 2.0

## Contributing

Contributions are welcome! Please open an issue or submit a pull request.
