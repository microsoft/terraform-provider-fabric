data "fabric_warehouse_snapshot" "example_by_id" {
  id           = "00000000-0000-0000-0000-000000000000"
  workspace_id = "11111111-1111-1111-1111-111111111111"
}

data "fabric_warehouse_snapshot" "example_by_name" {
  display_name = "example"
  workspace_id = "11111111-1111-1111-1111-111111111111"
}

# This is an invalid data source
# Do not specify `id` and `display_name` in the same data source block
# data "fabric_warehouse" "example" {
#   display_name = "example"
#   id           = "00000000-0000-0000-0000-000000000000"
#   workspace_id = "11111111-1111-1111-1111-111111111111"
# }
