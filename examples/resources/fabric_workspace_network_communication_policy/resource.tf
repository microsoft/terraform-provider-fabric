resource "fabric_workspace_network_communication_policy" "example" {
  workspace_id = "00000000-0000-0000-0000-000000000000"
  outbound = {
    public_access_rules = {
      default_action = "Deny"
    }
  }
  inbound = {
    public_access_rules = {
      default_action = "Deny"
    }
  }
}
