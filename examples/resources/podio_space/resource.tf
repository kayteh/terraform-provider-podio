data "podio_organization" "my_org" {
    url_label = "my-org"
}

resource "podio_space" "kanban_board" {
  org_id  = data.podio_organization.my_org.org_id
  name    = "Team Kanban"
  privacy = "closed"
}
