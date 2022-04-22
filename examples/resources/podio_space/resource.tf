resource "podio_space" "kanban_board" {
  org_id  = 1000
  name    = "Team Kanban"
  privacy = "closed"

  ignore_delete_errors = true
}
