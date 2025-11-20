# Example 1 - Data Pipeline without definition
resource "fabric_data_pipeline" "example" {
  display_name = "example"
  workspace_id = "00000000-0000-0000-0000-000000000000"
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

# Example 4 - Data Pipeline with custom tokens delimiter
resource "fabric_data_pipeline" "example_custom_delimiter" {
  display_name = "example"
  description  = "example with custom tokens delimiter"
  workspace_id = "00000000-0000-0000-0000-000000000000"
  format       = "Default"
  definition = {
    "pipeline-content.json" = {
      source           = "${local.path}/pipeline-content.json"
      tokens_delimiter = "##"
      tokens = {
        "MyValue" = "World"
      }
    }
  }
}

# Example 5 - Data Pipeline with parameters processing mode
resource "fabric_data_pipeline" "example_parameters" {
  display_name = "example"
  description  = "example with parameters processing mode"
  workspace_id = "00000000-0000-0000-0000-000000000000"
  format       = "Default"
  definition = {
    "pipeline-content.json" = {
      source          = "${local.path}/pipeline-content.json"
      processing_mode = "parameters"
      parameters = [
        {
          type  = "JsonPathReplace"
          find  = "$.properties.activities[0].name"
          value = "UpdatedActivityName"
        },
        {
          type  = "TextReplace"
          find  = "OldValue"
          value = "NewValue"
        }
      ]
    }
  }
}
