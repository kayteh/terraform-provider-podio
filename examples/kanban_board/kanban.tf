data "podio_organization" "my_org" {
  url_label = var.org_slug
}

resource "podio_space" "kanban" {
  name   = "Team Kanban"
  org_id = data.podio_organization.my_org.org_id
}

output "space_url" {
  value = podio_space.kanban.url
}
