resource "fabric_workspace" "example" {
  display_name = "example"
  description  = "Example Workspace"
}

resource "fabric_domain" "example" {
  display_name = "example"
  description  = "Example Domain"
}

resource "fabric_domain_workspace_assignments" "example" {
  domain_id = fabric_domain.example.id
  workspace_ids = [
    fabric_workspace.example.id
  ]
}
