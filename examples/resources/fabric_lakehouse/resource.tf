# Simple Lakehouse resource
resource "fabric_lakehouse" "example1" {
  display_name = "example1"
  workspace_id = "00000000-0000-0000-0000-000000000000"
  folder_id    = "11111111-1111-1111-1111-111111111111"
}

# Lakehouse resource with enabled schemas
resource "fabric_lakehouse" "example2" {
  display_name = "example2"
  description  = "example2 with enabled schemas"
  workspace_id = "00000000-0000-0000-0000-000000000000"

  configuration = {
    enable_schemas = true
  }
}
