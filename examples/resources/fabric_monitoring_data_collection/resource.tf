# The existing Monitoring item must already have a routing definition.
resource "fabric_monitoring_data_collection" "example" {
  workspace_id  = "00000000-0000-0000-0000-000000000000"
  monitoring_id = "11111111-1111-1111-1111-111111111111"
  enabled       = true
}

# Set enabled to false to stop collection without changing its definition.
# Removing this resource also disables collection, but does not delete the item.
