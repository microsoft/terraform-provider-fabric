# Example 1 - Item without definition
resource "fabric_variable_library" "example" {
  display_name = "example"
  workspace_id = "00000000-0000-0000-0000-000000000000"
}

# Example 2 - Item with definition without valueSets
resource "fabric_variable_library" "example_definition_bootstrap" {
  display_name              = "example"
  description               = "example with definition bootstrapping"
  workspace_id              = "00000000-0000-0000-0000-000000000000"
  definition_update_enabled = false
  format                    = "Default"
  definition = {
    "settings.json" = {
      source = "${local.path}/settings.json"
    }
    "variables.json" = {
      source = "${local.path}/variables.json"
    }
  }
}

# Example 3 - Item with definition with valueSets
resource "fabric_variable_library" "example_definition_update" {
  display_name = "example"
  description  = "Item with definition with valueSets"
  workspace_id = "00000000-0000-0000-0000-000000000000"
  format       = "Default"
  definition = {
    "settings.json" = {
      source = "${local.path}/settings.json"
    }
    "variables.json" = {
      source = "${local.path}/variables.json"
    }
    "valueSets/valueSet1.json" = {
      source = "${local.path}/valueSets/valueSet1.json"
    }
  }
}

# Example 4 - Item with custom tokens delimiter
resource "fabric_variable_library" "example_custom_delimiter" {
  display_name = "example"
  description  = "example with custom tokens delimiter"
  workspace_id = "00000000-0000-0000-0000-000000000000"
  format       = "Default"
  definition = {
    "variablelibrary-content.json" = {
      source           = "${local.path}/variablelibrary-content.json"
      tokens_delimiter = "##"
      tokens = {
        "MyValue1" = "my value 1"
        "MyValue2" = "my value 2"
      }
    }
  }
}

# Example 5 - Item with parameters processing mode
resource "fabric_variable_library" "example_parameters" {
  display_name = "example"
  description  = "example with parameters processing mode"
  workspace_id = "00000000-0000-0000-0000-000000000000"
  format       = "Default"
  definition = {
    "variablelibrary-content.json" = {
      source          = "${local.path}/variablelibrary-content.json"
      processing_mode = "parameters"
      parameters = [
        {
          type  = "JsonPathReplace"
          find  = "$.variables[0].name"
          value = "UpdatedVariableName"
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
