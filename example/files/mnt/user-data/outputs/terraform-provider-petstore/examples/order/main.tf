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

# Create a pet first
resource "petstore_pet" "my_cat" {
  name = "Whiskers"
  photo_urls = [
    "https://example.com/cat.jpg"
  ]
  status = "available"

  category {
    name = "Cats"
  }
}

# Create an order for the pet
resource "petstore_order" "cat_order" {
  pet_id   = petstore_pet.my_cat.id
  quantity = 1
  ship_date = "2024-12-25T10:00:00Z"
  status   = "placed"
  complete = false
}

# Data source to read an order
data "petstore_order" "existing_order" {
  id = 54321
}

output "order_status" {
  value = data.petstore_order.existing_order.status
}
