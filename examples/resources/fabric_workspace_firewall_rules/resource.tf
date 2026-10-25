# Firewall rules are only enforced when inbound public access is set to Deny.
resource "fabric_workspace_network_communication_policy" "example" {
  workspace_id = "00000000-0000-0000-0000-000000000000"
  inbound = {
    public_access_rules = {
      default_action = "Deny"
    }
  }
}

resource "fabric_workspace_firewall_rules" "example" {
  workspace_id = "00000000-0000-0000-0000-000000000000"

  rules = [
    {
      display_name = "corp-egress"
      value        = "203.0.113.0/24"
    },
    {
      display_name = "vpn-range"
      value        = "198.51.100.10-198.51.100.20"
    },
    {
      display_name = "single-host"
      value        = "192.0.2.42"
    },
  ]

  depends_on = [fabric_workspace_network_communication_policy.example]
}
