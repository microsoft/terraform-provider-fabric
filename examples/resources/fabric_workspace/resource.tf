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

# Workspace with Domain Assignment
data "fabric_domain" "example" {
  display_name = "example"
}

resource "fabric_workspace" "example3" {
  display_name = "example3"
  description  = "Example Workspace 3 with Domain"
  domain_id    = data.fabric_domain.example.id
}

# Workspace with both Capacity and Domain
resource "fabric_workspace" "example4" {
  display_name = "example4"
  description  = "Example Workspace 4 with Capacity and Domain"
  capacity_id  = data.fabric_capacity.example.id
  domain_id    = data.fabric_domain.example.id
}
