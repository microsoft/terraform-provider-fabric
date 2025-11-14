# Example 1 - Data Pipeline without definition
resource "fabric_data_pipeline" "example" {
  display_name = "example"
  workspace_id = "00000000-0000-0000-0000-000000000000"
  folder_id    = "11111111-1111-1111-1111-111111111111"
}

# Example 2 - Data Pipeline with definition bootstrapping only
resource "fabric_data_pipeline" "example_definition_bootstrap" {
  display_name              = "example"
  description               = "example with definition bootstrapping"
  workspace_id              = "00000000-0000-0000-0000-000000000000"
  format                    = "Default"
  definition_update_enabled = false
  definition = {
    "pipeline-content.json" = {
      source = "${local.path}/pipeline-content.json"
      tokens = {
        "MyValue" = "World"
      }
    }
  }
}

# Example 3 - Data Pipeline with definition update when source or tokens changed
resource "fabric_data_pipeline" "example_definition_update" {
  display_name = "example"
  description  = "example with definition update when source or tokens changed"
  workspace_id = "00000000-0000-0000-0000-000000000000"
  format       = "Default"
  definition = {
    "pipeline-content.json" = {
      source = "${local.path}/pipeline-content.json"
      tokens = {
        "MyValue" = "World"
      }
    }
  }
}
