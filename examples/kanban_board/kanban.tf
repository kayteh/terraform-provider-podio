data "podio_organization" "my_org" {
  url_label = var.org_slug
}

resource "podio_space" "kanban" {
  name      = "Team Kanban"
  org_id    = data.podio_organization.my_org.org_id
  auto_join = false
}

resource "podio_app" "kanban" {
  name        = "Kanban"
  space_id    = podio_space.kanban.space_id
  description = "My cool team kanban app"
  usage       = "Backlog is open, do stuff, and mark as done!"
  item_name   = "Task"
  icon        = "22.png"
}

output "space_url" {
  value = podio_space.kanban.url
}
