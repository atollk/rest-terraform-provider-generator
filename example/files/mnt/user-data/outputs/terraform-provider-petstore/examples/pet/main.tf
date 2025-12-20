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

# Create a pet
resource "petstore_pet" "my_dog" {
  name = "Buddy"
  photo_urls = [
    "https://example.com/photo1.jpg",
    "https://example.com/photo2.jpg"
  ]
  status = "available"

  category {
    name = "Dogs"
  }

  tags {
    name = "friendly"
  }

  tags {
    name = "golden-retriever"
  }
}

# Data source to read a pet
data "petstore_pet" "existing_pet" {
  id = 12345
}

output "pet_name" {
  value = data.petstore_pet.existing_pet.name
}
