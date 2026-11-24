resource "fabric_workspace_onelake_immutability_policy" "example" {
  workspace_id   = "00000000-0000-0000-0000-000000000000"
  scope          = "DiagnosticLogs"
  retention_days = 30
}
