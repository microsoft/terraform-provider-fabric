# Example 1 - Item without definition properties
resource "fabric_org_app" "example" {
  display_name = "example1"
  workspace_id = "00000000-0000-0000-0000-000000000000"
}

# Example 2 - Item with definition bootstrapping only
resource "fabric_org_app" "example_definition_bootstrap" {
  display_name              = "example2"
  description               = "example with definition bootstrapping"
  workspace_id              = "00000000-0000-0000-0000-000000000000"
  definition_update_enabled = false
  format                    = "JSON"
  definition = {
    "definition.json" = {
      source = "${local.path}/definition.json"
    }
  }
}

# Example 3 - Item with definition update when source or tokens changed
resource "fabric_org_app" "example_definition_update" {
  display_name = "example3"
  description  = "example with definition update when source or tokens changed"
  workspace_id = "00000000-0000-0000-0000-000000000000"
  format       = "JSON"
  definition = {
    "definition.json" = {
      source = "${local.path}/definition.json.tmpl"
      tokens = {
        "OverviewTitle" = "Updated overview"
        "OverviewBody"  = "Updated Org App content."
      }
    }
  }
}

# Example 4 - Item with custom tokens delimiter
resource "fabric_org_app" "example_custom_delimiter" {
  display_name = "example4"
  description  = "example with custom tokens delimiter"
  workspace_id = "00000000-0000-0000-0000-000000000000"
  format       = "JSON"
  definition = {
    "definition.json" = {
      source           = "${local.path}/definition-custom-delimiter.json.tmpl"
      tokens_delimiter = "<<>>"
      tokens = {
        "OverviewTitle" = "Updated overview"
        "OverviewBody"  = "Updated Org App content."
      }
    }
  }
}

# Example 5 - Item with parameters processing mode
resource "fabric_org_app" "example_parameters" {
  display_name = "example5"
  description  = "example with parameters processing mode"
  workspace_id = "00000000-0000-0000-0000-000000000000"
  format       = "JSON"
  definition = {
    "definition.json" = {
      source          = "${local.path}/definition.json"
      processing_mode = "parameters"
      parameters = [
        {
          type  = "JsonPathReplace"
          find  = "$.elements[0].header.title"
          value = "Updated overview"
        },
        {
          type  = "TextReplace"
          find  = "Welcome to the Org App."
          value = "Updated Org App content."
        }
      ]
    }
  }
}
