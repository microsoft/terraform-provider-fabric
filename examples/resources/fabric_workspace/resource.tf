# Simple Workspace
resource "fabric_workspace" "example1" {
  display_name = "example1"
  description  = "Example Workspace 1"
}

# Workspace with Capacity and Identity
data "fabric_capacity" "example" {
  display_name = "example"
}

resource "fabric_workspace" "example2" {
  display_name = "example2"
  description  = "Example Workspace 2"
  capacity_id  = data.fabric_capacity.example.id
  identity = {
    type = "SystemAssigned"
  }
}

# Workspace with skip_capacity_state_validation
# Use this when the caller does not have permissions to list capacities
resource "fabric_workspace" "example3" {
  display_name                   = "example3"
  description                    = "Example Workspace 3"
  capacity_id                    = "00000000-0000-0000-0000-000000000000"
  skip_capacity_state_validation = true
}
