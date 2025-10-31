# Example 1 - Item without definition
resource "fabric_variable_library" "example" {
  display_name = "example"
  workspace_id = "00000000-0000-0000-0000-000000000000"
}

# Example 2 - Item with definition bootstrapping only
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
