# Example 1 - basic item
resource "fabric_graphql_api" "example" {
  display_name = "example"
  workspace_id = "00000000-0000-0000-0000-000000000000"
}

# Example 2 - item with definition bootstrapping only
resource "fabric_graphql_api" "example_definition_bootstrap" {
  display_name              = "example"
  description               = "example with definition bootstrapping"
  workspace_id              = "00000000-0000-0000-0000-000000000000"
  format                    = "Default"
  definition_update_enabled = false
  definition = {
    "graphql-definition.json" = {
      source = "${local.path}/graphql-definition.json"
      tokens = {
        "MyKey" = "MyValue"
      }
    }
  }
}

# Example 3 - item with definition update when source or tokens changed
resource "fabric_graphql_api" "example_definition_update" {
  display_name = "example"
  description  = "example with definition update when source or tokens changed"
  workspace_id = "00000000-0000-0000-0000-000000000000"
  format       = "Default"
  definition = {
    "graphql-definition.json" = {
      source = "${local.path}/graphql-definition.json"
      tokens = {
        "MyKey" = "MyValue"
      }
    }
  }
}
