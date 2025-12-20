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
resource "petstore_user" "john_doe" {
  username   = "johndoe"
  first_name = "John"
  last_name  = "Doe"
  email      = "john.doe@example.com"
  password   = "supersecret123"
  phone      = "+1234567890"
  user_status = 1
}

# Data source to read a user
data "petstore_user" "existing_user" {
  username = "janedoe"
}

output "user_email" {
  value = data.petstore_user.existing_user.email
}
