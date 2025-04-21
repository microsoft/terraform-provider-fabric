data "fabric_connection" "example_by_id" {
  id = "9ad6978e-4a52-452e-91f8-6edd60e27e89"
}

# data "fabric_connection" "example_by_name" {
#   display_name = "azdo-board-aaaa"
# }

# This is an invalid data source
# Do not specify `id` and `display_name` in the same data source block
# data "fabric_connection" "example" {
#   display_name = "example"
#   id = "00000000-0000-0000-0000-000000000000"
# }
