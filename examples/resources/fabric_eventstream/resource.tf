# Example 1 - Item without definition
resource "fabric_eventstream" "example" {
  display_name = "example"
  workspace_id = "00000000-0000-0000-0000-000000000000"
  folder_id    = "11111111-1111-1111-1111-111111111111"
}

# Example 2 - Item with definition bootstrapping only
resource "fabric_eventstream" "example_definition_bootstrap" {
  display_name              = "example"
  description               = "example with definition bootstrapping"
  workspace_id              = "00000000-0000-0000-0000-000000000000"
  definition_update_enabled = false
  format                    = "Default"
  definition = {
    "eventstream.json" = {
      source = "${local.path}/eventstream.json.tmpl"
    }
  }
}

# Example 3 - Item with definition update when source or tokens changed
resource "fabric_eventstream" "example_definition_update" {
  display_name = "example"
  description  = "example with definition update when source or tokens changed"
  workspace_id = "00000000-0000-0000-0000-000000000000"
  format       = "Default"
  definition = {
    "eventstream.json" = {
      source = "${local.path}/eventstream.json.tmpl"
      tokens = {
        "LakehouseWorkspaceID" = "11111111-1111-1111-1111-111111111111"
        "LakehouseID"          = "22222222-2222-2222-2222-222222222222"
      }
    }
  }
}
