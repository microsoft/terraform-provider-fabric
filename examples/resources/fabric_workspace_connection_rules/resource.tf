resource "fabric_workspace_outbound_cloud_connection_rules" "example" {
  workspace_id   = "00000000-0000-0000-0000-000000000000"
  default_action = "Deny"
  rules = [
    {
      connection_type : "SQL",
      default_action : "Deny",
      allowed_endpoints = [
        {
          host_name_pattern = "*.microsoft.com"
        }
      ]
    },
    {
      connection_type : "LakeHouse",
      default_action : "Deny",
      allowed_workspaces = [
        {
          workspace_id = "11111111-1111-1111-1111-111111111111"
        }
      ]
    },
  ]
}
