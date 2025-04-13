resource "fabric_workspace_managed_private_endpoint" "example" {
  workspace_id                    = "00000000-0000-0000-0000-000000000000"
  name                            = "example"
  target_private_link_resource_id = "/subscriptions/11111111-1111-1111-1111-111111111111/resourceGroups/RESOURCE_GROUP_NAME/providers/Microsoft.Storage/storageAccounts/RESOURCE_NAME"
  target_subresource_type         = "blob"
  request_message                 = "Request message to approve private endpoint"
}
