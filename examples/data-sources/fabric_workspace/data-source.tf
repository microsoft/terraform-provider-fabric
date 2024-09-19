data "fabric_workspace" "example_by_id" {
  id = "00000000-0000-0000-0000-000000000000"
}

data "fabric_workspace" "example_by_name" {
  display_name = "example"
}

# This is an invalid data source
# Do not specify `id` and `display_name` in the same data source block
# data "fabric_workspace" "example" {
#   display_name = "example"
#   id = "00000000-0000-0000-0000-000000000000"
# }
