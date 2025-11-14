resource "fabric_warehouse_snapshot" "example" {
  display_name = "warehouse_example"
  workspace_id = "00000000-0000-0000-0000-000000000000"
  folder_id    = "11111111-1111-1111-1111-111111111111"
  configuration = {
    parent_warehouse_id = "11111111-1111-1111-1111-111111111111"
    #snapshot_date_time if not provided the current date and time will be taken
  }
}

