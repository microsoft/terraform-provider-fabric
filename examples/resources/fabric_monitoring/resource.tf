resource "fabric_monitoring" "example" {
  display_name = "WorkspaceMonitoring"
  description  = "Workspace diagnostic collection configuration"
  workspace_id = "00000000-0000-0000-0000-000000000000"
  format       = "Default"
  definition = {
    "MonitoringProperties.json" = {
      source = "${path.module}/MonitoringProperties.json"
    }
  }
}

# Defining routing does not enable collection.
resource "fabric_monitoring_data_collection" "example" {
  workspace_id  = fabric_monitoring.example.workspace_id
  monitoring_id = fabric_monitoring.example.id
  enabled       = true
}
