resource "fabric_workspace_outbound_gateway_rules" "example" {
  workspace_id   = "00000000-0000-0000-0000-000000000000"
  default_action = "Allow"
  allowed_gateways = [
    {
      id = "11111111-1111-1111-1111-111111111111"
    }
  ]
}
