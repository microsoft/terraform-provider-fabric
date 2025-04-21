# Example of using id to fetch the fabric capacity details
data "fabric_gateway_member" "example_by_id" {
  gateway_id = "00000000-0000-0000-0000-000000000000"
  id         = "11111111-1111-1111-1111-111111111111"
}

# Example of using display_name to fetch the fabric capacity details
data "fabric_gateway_member" "example_by_name" {
  gateway_id   = "00000000-0000-0000-0000-000000000000"
  display_name = "example"
}

# This is an invalid data source
# Do not specify `id` and `display_name` in the same data source block
# data "fabric_gateway_member" "example" {
#   gateway_id = "00000000-0000-0000-0000-000000000000"
#   display_name = "example"
#   id = "11111111-1111-1111-1111-111111111111"
# }
