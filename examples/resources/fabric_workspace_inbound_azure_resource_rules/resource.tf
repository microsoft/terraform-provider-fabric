resource "fabric_workspace_inbound_azure_resource_rules" "example" {
  workspace_id = "00000000-0000-0000-0000-000000000000"
  rules = [
    {
      display_name = "Azure Data Factory - RESOURCE_NAME"
      resource_id  = "/subscriptions/11111111-1111-1111-1111-111111111111/resourceGroups/RESOURCE_GROUP_NAME/providers/Microsoft.DataFactory/factories/RESOURCE_NAME"
    },
    {
      display_name = "Azure SQL Server - RESOURCE_NAME"
      resource_id  = "/subscriptions/11111111-1111-1111-1111-111111111111/resourceGroups/RESOURCE_GROUP_NAME/providers/Microsoft.Sql/servers/RESOURCE_NAME"
    }
  ]
}
