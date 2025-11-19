resource "fabric_warehouse" "example" {
  display_name = "warehouse_example"
  workspace_id = "11111111-1111-1111-1111-111111111111"
}

# Warehouse resource with enabled schemas
resource "fabric_warehouse" "example2" {
  display_name = "warehouse_example2"
  description  = "warehouse_example2 with collation_type"
  workspace_id = "00000000-0000-0000-0000-000000000000"

  configuration = {
    collation_type = "Latin1_General_100_BIN2_UTF8"
  }
}
