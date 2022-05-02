terraform {
  required_providers {
    podio = {
      source  = "kayteh/podio"
      version = ">=1.0.0"
    }
  }
}

provider "podio" {
  client_id     = var.client_id
  client_secret = var.client_secret
  username      = var.username
  password      = var.password
}
