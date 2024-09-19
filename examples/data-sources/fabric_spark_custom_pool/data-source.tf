data "fabric_spark_custom_pool" "example_by_id" {
  id           = "11111111-1111-1111-1111-111111111111"
  workspace_id = "00000000-0000-0000-0000-000000000000"
}

data "fabric_spark_custom_pool" "example_by_name" {
  name         = "example"
  workspace_id = "00000000-0000-0000-0000-000000000000"
}

# This is an invalid data source
# Do not specify id and name in the same data source block
# data "fabric_spark_custom_pool" "example" {
#   name = "example"
#   id = "11111111-1111-1111-1111-111111111111"
#   workspace_id = "00000000-0000-0000-0000-000000000000"
# }
