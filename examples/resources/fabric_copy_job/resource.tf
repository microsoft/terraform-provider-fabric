# Example 1 - Item without definition
resource "fabric_copy_job" "example_definition" {
  display_name = "example"
  description  = "example without definition"
  workspace_id = "00000000-0000-0000-0000-000000000000"
}

# Example 2 - Item with definition bootstrapping only
resource "fabric_copy_job" "example_definition_bootstrap" {
  display_name              = "example"
  description               = "example with definition bootstrapping"
  workspace_id              = "00000000-0000-0000-0000-000000000000"
  definition_update_enabled = false
  format                    = "Default"
  definition = {
    "copyjob-content.json" = {
      source = "${local.path}/copyjob-content.json.tmpl"
    }
  }
}

# Example 3 - Item with definition update when source or tokens changed
resource "fabric_copy_job" "example_definition_update" {
  display_name = "example"
  description  = "example with definition update when source or tokens changed"
  workspace_id = "00000000-0000-0000-0000-000000000000"
  format       = "Default"
  definition = {
    "copyjob-content.json" = {
      source = "${local.path}/copyjob-content.json.tmpl"
      tokens = {
        "SOURCE_WORKSPACE_ID"      = "Source Workspace ID"
        "SOURCE_ARTIFACT_ID"       = "Source Artifact ID"
        "DESTINATION_WORKSPACE_ID" = "Destination Workspace ID"
        "DESTINATION_ARTIFACT_ID"  = "Destination Artifact ID"
      }
    }
  }
}

# Example 4 - Item with custom tokens delimiter
resource "fabric_copy_job" "example_custom_delimiter" {
  display_name = "example"
  description  = "example with custom tokens delimiter"
  workspace_id = "00000000-0000-0000-0000-000000000000"
  format       = "Default"
  definition = {
    "copyjob-content.json" = {
      source           = "${local.path}/copyjob-content.json.tmpl"
      tokens_delimiter = "##"
      tokens = {
        "SOURCE_WORKSPACE_ID"      = "Source Workspace ID"
        "SOURCE_ARTIFACT_ID"       = "Source Artifact ID"
        "DESTINATION_WORKSPACE_ID" = "Destination Workspace ID"
        "DESTINATION_ARTIFACT_ID"  = "Destination Artifact ID"
      }
    }
  }
}

# Example 5 - Item with parameters processing mode
resource "fabric_copy_job" "example_parameters" {
  display_name = "example"
  description  = "example with parameters processing mode"
  workspace_id = "00000000-0000-0000-0000-000000000000"
  format       = "Default"
  definition = {
    "copyjob-content.json" = {
      source          = "${local.path}/copyjob-content.json.tmpl"
      processing_mode = "parameters"
      parameters = [
        {
          type  = "JsonPathReplace"
          find  = "$.properties.source.workspaceId"
          value = "00000000-0000-0000-0000-000000000001"
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
