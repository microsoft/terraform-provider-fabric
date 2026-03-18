resource "fabric_workspace_git_outbound_policy" "example" {
  workspace_id   = "00000000-0000-0000-0000-000000000000"
  default_action = "Deny"
}
