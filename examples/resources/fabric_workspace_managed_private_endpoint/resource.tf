resource "fabric_workspace_managed_private_endpoint" "example" {
  workspace_id                    = "3f18b478-6b93-4977-b116-e507d9e64a3f"
  name                            = "testprvendpoint1"
  target_private_link_resource_id = "/subscriptions/f4cea851-983c-45cb-954d-9fce8328d90c/resourceGroups/rg-fabric-tf-tests/providers/Microsoft.Storage/storageAccounts/testaccstlegxrfvazo"
  target_subresource_type         = "blob"
  request_message                 = "Request message to approve private endpoint"
}
