# Example 1 - Item with definition bootstrapping only
resource "fabric_mirrored_database" "example_definition_bootstrap" {
  display_name              = "example2"
  description               = "example with definition bootstrapping"
  workspace_id              = "00000000-0000-0000-0000-000000000000"
  definition_update_enabled = false
  format                    = "Default"
  definition = {
    "mirroring.json" = {
      source = "${local.path}/mirroring.json.tmpl"
    }
  }
}

# Example 2 - Item with definition update when source or tokens changed
resource "fabric_mirrored_database" "example_definition_update" {
  display_name = "example3"
  description  = "example with definition update when source or tokens changed"
  workspace_id = "00000000-0000-0000-0000-000000000000"
  format       = "Default"
  definition = {
    "mirroring.json" = {
      source = "${local.path}/mirroring.json.tmpl"
      tokens = {
        "DEFAULT_SCHEMA" = "my_schema"
      }
    }
  }
}

# Example 3 - Item with custom tokens delimiter
resource "fabric_mirrored_database" "example_custom_delimiter" {
  display_name = "example4"
  description  = "example with custom tokens delimiter"
  workspace_id = "00000000-0000-0000-0000-000000000000"
  format       = "Default"
  definition = {
    "mirroring.json" = {
      source           = "${local.path}/mirroring.json.tmpl"
      tokens_delimiter = "##"
      tokens = {
        "DEFAULT_SCHEMA" = "my_schema"
      }
    }
  }
}

# Example 4 - Item with parameters processing mode
resource "fabric_mirrored_database" "example_parameters" {
  display_name = "example5"
  description  = "example with parameters processing mode"
  workspace_id = "00000000-0000-0000-0000-000000000000"
  format       = "Default"
  definition = {
    "mirroring.json" = {
      source          = "${local.path}/mirroring.json.tmpl"
      processing_mode = "parameters"
      parameters = [
        {
          type  = "JsonPathReplace"
          find  = "$.properties.defaultSchema"
          value = "updated_schema"
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
