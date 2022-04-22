provider "podio" {
}

variable "org_slug" {
  type         = string
  ddescription = "Human-readable Org name/slug. The `citrix` in https://podio.com/citrix/hello-world"
}

data "podio_organization" "my_org" {
  slug = var.org_slug
}

resource "podio_space" "kanban" {
  name   = "Team Kanban"
  org_id = data.podio_organization.id
}