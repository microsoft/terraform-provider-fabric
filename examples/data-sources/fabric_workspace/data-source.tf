data "fabric_workspace" "example_by_id" {
  id = "00000000-0000-0000-0000-000000000000"
}

data "fabric_workspace" "example_by_name" {
  display_name = "example"
}

# Workspace with skip_capacity_state_validation
# Use this when the caller does not have permissions to list capacities
data "fabric_workspace" "example_skip_validation" {
  display_name                   = "example"
  skip_capacity_state_validation = true
}

# This is an invalid data source
# Do not specify `id` and `display_name` in the same data source block
# data "fabric_workspace" "example" {
#   display_name = "example"
#   id = "00000000-0000-0000-0000-000000000000"
# }
