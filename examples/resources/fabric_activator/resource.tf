# Example 1 - Item without definition
resource "fabric_activator" "example" {
  display_name = "example"
  workspace_id = "00000000-0000-0000-0000-000000000000"
}

# Example 2 - Item with definition bootstrapping only
resource "fabric_activator" "example_definition_bootstrap" {
  display_name              = "example"
  description               = "example with definition bootstrapping"
  workspace_id              = "00000000-0000-0000-0000-000000000000"
  definition_update_enabled = false
  format                    = "Default"
  definition = {
    "ReflexEntities.json" = {
      source = "${local.path}/ReflexEntities.json"
    }
  }
}

# Example 3 - Item with definition update when source or tokens changed
resource "fabric_activator" "example_definition_update" {
  display_name = "example"
  description  = "example with definition update when source or tokens changed"
  workspace_id = "00000000-0000-0000-0000-000000000000"
  format       = "Default"
  definition = {
    "ReflexEntities.json" = {
      source = "${local.path}/ReflexEntities.json"
      tokens = {
        "MyValue1" = "my value 1"
        "MyValue2" = "my value 2"
      }
    }
  }
}

# Example 4 - Item with custom tokens delimiter
resource "fabric_activator" "example_custom_delimiter" {
  display_name = "example"
  description  = "example with custom tokens delimiter"
  workspace_id = "00000000-0000-0000-0000-000000000000"
  format       = "Default"
  definition = {
    "ReflexEntities.json" = {
      source           = "${local.path}/ReflexEntities.json"
      tokens_delimiter = "##"
      tokens = {
        "MyValue1" = "my value 1"
        "MyValue2" = "my value 2"
      }
    }
  }
}

# Example 5 - Item with parameters processing mode
resource "fabric_activator" "example_parameters" {
  display_name = "example"
  description  = "example with parameters processing mode"
  workspace_id = "00000000-0000-0000-0000-000000000000"
  format       = "Default"
  definition = {
    "ReflexEntities.json" = {
      source          = "${local.path}/ReflexEntities.json"
      processing_mode = "parameters"
      parameters = [
        {
          type  = "JsonPathReplace"
          find  = "$.properties.name"
          value = "UpdatedName"
        },
        {
          type  = "TextReplace"
          find  = "OldText"
          value = "NewText"
        }
      ]
    }
  }
}
