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
