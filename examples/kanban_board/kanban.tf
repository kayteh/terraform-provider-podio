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

data "podio_organization" "my_org" {
  url_label = var.org_slug
}

resource "podio_space" "kanban" {
  name   = "My Cool Kanban Workspace"
  org_id = data.podio_organization.my_org.org_id
}

output "space_url" {
  value = podio_space.kanban.url
}
