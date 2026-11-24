resource "fabric_onelake_diagnostic_setting" "example" {
  workspace_id = "00000000-0000-0000-0000-000000000000"
  status       = "Enabled"

  destination = {
    lakehouse = {
      workspace_id = "00000000-0000-0000-0000-000000000000"
      item_id      = "00000000-0000-0000-0000-000000000000"
    }
  }
}
