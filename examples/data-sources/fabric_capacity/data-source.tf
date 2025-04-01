# Example of using id to fetch the fabric capacity details
data "fabric_capacity" "example_by_id" {
  id = "00000000-0000-0000-0000-000000000000"
}

# Example of using display_name to fetch the fabric capacity details
data "fabric_capacity" "example_by_name" {
  display_name = "example"
}

# This is an invalid data source
# Do not specify `id` and `display_name` in the same data source block
# data "fabric_capacity" "example" {
#   display_name = "example"
#   id = "00000000-0000-0000-0000-000000000000"
# }

# It's recommended to use `lifecycle` with `postcondition` block to handle the state of the capacity.
data "fabric_capacity" "example" {
  display_name = "example"

  lifecycle {
    postcondition {
      condition     = self.state == "Active"
      error_message = "Fabric Capacity is not in Active state. Please check the Fabric Capacity status."
    }
  }
}
