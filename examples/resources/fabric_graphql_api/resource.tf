# Example 1 - Item without definition
resource "fabric_graphql_api" "example" {
  display_name = "example"
  workspace_id = "00000000-0000-0000-0000-000000000000"
}

# Example 2 - Item with definition bootstrapping only
resource "fabric_graphql_api" "example_definition_bootstrap" {
  display_name              = "example"
  description               = "example with definition bootstrapping"
  workspace_id              = "00000000-0000-0000-0000-000000000000"
  definition_update_enabled = false
  definition = {
    "graphql-definition.json" = {
      source = "${local.path}/graphql-definition.json.tmpl"
    }
  }
}

# Example 3 - Item with definition update when source or tokens changed
resource "fabric_graphql_api" "example_definition_update" {
  display_name = "example"
  description  = "example with definition update when source or tokens changed"
  workspace_id = "00000000-0000-0000-0000-000000000000"
  definition = {
    "graphql-definition.json" = {
      source = "${local.path}/graphql-definition.json.tmpl"
      tokens = {
        "CONNECTION_ID" = "11111111-1111-1111-1111-111111111111"
        "TABLE_NAME"    = "my_table"
      }
    }
  }
}


# Example 4 - Item with custom tokens delimiter
resource "fabric_graphql_api" "example_custom_delimiter" {
  display_name = "example"
  description  = "example with custom tokens delimiter"
  workspace_id = "00000000-0000-0000-0000-000000000000"
  format       = "Default"
  definition = {
    "graphql-definition.json" = {
      source           = "${local.path}/graphql-definition.json.tmpl"
      tokens_delimiter = "##"
      tokens = {
        "CONNECTION_ID" = "11111111-1111-1111-1111-111111111111"
        "TABLE_NAME"    = "my_table"
      }
    }
  }
}

# Example 5 - Item with parameters processing mode
resource "fabric_graphql_api" "example_parameters" {
  display_name = "example"
  description  = "example with parameters processing mode"
  workspace_id = "00000000-0000-0000-0000-000000000000"
  format       = "Default"
  definition = {
    "graphql-definition.json" = {
      source          = "${local.path}/graphql-definition.json.tmpl"
      processing_mode = "parameters"
      parameters = [
        {
          type  = "JsonPathReplace"
          find  = "$.connectionId"
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
