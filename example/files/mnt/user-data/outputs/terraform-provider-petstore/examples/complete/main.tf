terraform {
  required_providers {
    petstore = {
      source = "yourusername/petstore"
    }
  }
}

provider "petstore" {
  base_url = "https://petstore3.swagger.io/api/v3"
  api_key  = "your-api-key-here"
}

# Create a user
resource "petstore_user" "customer" {
  username    = "customer123"
  first_name  = "Jane"
  last_name   = "Smith"
  email       = "jane.smith@example.com"
  password    = "securepassword"
  phone       = "+1987654321"
  user_status = 1
}

# Create multiple pets
resource "petstore_pet" "dog" {
  name = "Max"
  photo_urls = [
    "https://example.com/max1.jpg",
    "https://example.com/max2.jpg"
  ]
  status = "available"

  category {
    name = "Dogs"
  }

  tags {
    name = "energetic"
  }

  tags {
    name = "labrador"
  }
}

resource "petstore_pet" "cat" {
  name = "Luna"
  photo_urls = [
    "https://example.com/luna.jpg"
  ]
  status = "available"

  category {
    name = "Cats"
  }

  tags {
    name = "calm"
  }

  tags {
    name = "persian"
  }
}

# Create orders for the pets
resource "petstore_order" "dog_order" {
  pet_id    = petstore_pet.dog.id
  quantity  = 1
  ship_date = "2024-12-30T14:00:00Z"
  status    = "placed"
  complete  = false
}

resource "petstore_order" "cat_order" {
  pet_id    = petstore_pet.cat.id
  quantity  = 1
  ship_date = "2024-12-31T10:00:00Z"
  status    = "approved"
  complete  = false
}

# Use data sources to read information
data "petstore_user" "customer_data" {
  username = petstore_user.customer.username

  depends_on = [petstore_user.customer]
}

data "petstore_pet" "dog_data" {
  id = petstore_pet.dog.id

  depends_on = [petstore_pet.dog]
}

# Outputs
output "customer_email" {
  description = "Customer email address"
  value       = data.petstore_user.customer_data.email
}

output "dog_name" {
  description = "Dog's name"
  value       = data.petstore_pet.dog_data.name
}

output "dog_order_id" {
  description = "Dog order ID"
  value       = petstore_order.dog_order.id
}

output "cat_order_id" {
  description = "Cat order ID"
  value       = petstore_order.cat_order.id
}

output "all_pet_ids" {
  description = "All created pet IDs"
  value = {
    dog = petstore_pet.dog.id
    cat = petstore_pet.cat.id
  }
}
