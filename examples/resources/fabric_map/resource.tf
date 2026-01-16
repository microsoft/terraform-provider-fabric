# Example 1 - Item without definition
resource "fabric_map" "example" {
  display_name = "example"
  workspace_id = "00000000-0000-0000-0000-000000000000"
}

# Example 2 - Item with definition bootstrapping only
resource "fabric_map" "example_definition_bootstrap" {
  display_name              = "example"
  description               = "example with definition bootstrapping"
  workspace_id              = "00000000-0000-0000-0000-000000000000"
  definition_update_enabled = false
  format                    = "Default"
  definition = {
    "map.json" = {
      source = "${local.path}/map.json.tmpl"
    }
  }
}

# Example 3 - Item with definition update when source or tokens changed
resource "fabric_map" "example_definition_update" {
  display_name = "example"
  description  = "example with definition update when source or tokens changed"
  workspace_id = "00000000-0000-0000-0000-000000000000"
  format       = "Default"
  definition = {
    "map.json" = {
      source = "${local.path}/map.json.tmpl"
      tokens = {
        "LAKEHOUSES" = "{\"workspaceId\": \"00000000-0000-0000-0000-000000000000\", \"artifactId\": \"11111111-1111-1111-1111-111111111111\"}"
      }
    }
  }
}

# Example 4 - Item with custom tokens delimiter
resource "fabric_map" "example_custom_delimiter" {
  display_name = "example"
  description  = "example with custom tokens delimiter"
  workspace_id = "00000000-0000-0000-0000-000000000000"
  format       = "Default"
  definition = {
    "map.json" = {
      source           = "${local.path}/map.json.tmpl"
      tokens_delimiter = "##"
      tokens = {
        "LAKEHOUSES" = "{\"workspaceId\": \"00000000-0000-0000-0000-000000000000\", \"artifactId\": \"11111111-1111-1111-1111-111111111111\"}"
      }
    }
  }
}

# Example 5 - Item with parameters processing mode
resource "fabric_map" "example_parameters" {
  display_name = "example"
  description  = "example with parameters processing mode"
  workspace_id = "00000000-0000-0000-0000-000000000000"
  format       = "Default"
  definition = {
    "map.json" = {
      source          = "${local.path}/map.json.tmpl"
      processing_mode = "parameters"
      parameters = [
        {
          type  = "JsonPathReplace"
          find  = "$.dataSources.lakehouses[0].artifactId"
          value = "11111111-1111-1111-1111-111111111111"
        },
        {
          type  = "TextReplace"
          find  = "OldLakehouseID"
          value = "NewLakehouseID"
        }
      ]
    }
  }
}
