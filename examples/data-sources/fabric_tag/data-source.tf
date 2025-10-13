data "fabric_tag" "example_by_id" {
  id = "11111111-1111-1111-1111-111111111111"
}

data "fabric_tag" "example_by_name" {
  display_name = "example"
}

# This is an invalid data source
# Do not specify `id` and `display_name` in the same data source block
# data "fabric_tag" "example" {
#   display_name = "example"
#   id           = "11111111-1111-1111-1111-111111111111"
# }
