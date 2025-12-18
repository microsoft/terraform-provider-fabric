# Example 1 - Item without definition
resource "fabric_mirrored_azure_databricks_catalog" "example" {
  display_name = "example1"
  workspace_id = "00000000-0000-0000-0000-000000000000"
}

# Example 2 - Item with definition bootstrapping only
resource "fabric_mirrored_azure_databricks_catalog" "example_definition_bootstrap" {
  display_name              = "example2"
  description               = "example with definition bootstrapping"
  workspace_id              = "00000000-0000-0000-0000-000000000000"
  format                    = "Default"
  definition_update_enabled = false # <-- Disable definition update
  definition = {
    "definition.json" = {
      source = "${local.path}/definition.json.tmpl"
    }
  }
}

# Example 3 - Item with definition update when source or tokens changed
resource "fabric_mirrored_azure_databricks_catalog" "example_definition_update" {
  display_name = "example3"
  description  = "example with definition update when source or tokens changed"
  workspace_id = "00000000-0000-0000-0000-000000000000"
  format       = "Default"
  definition = {
    "definition.json" = {
      source = "${local.path}/definition.json.tmpl"
      tokens = {
        "CATALOG_NAME"                       = "MyCatalogName"
        "DATABRICKS_WORKSPACE_CONNECTION_ID" = "00000000-0000-0000-0000-000000000000"
      }
    }
  }
}

# Example 3a - Item with custom tokens delimiter
resource "fabric_mirrored_azure_databricks_catalog" "example_custom_delimiter" {
  display_name = "example3a"
  description  = "example with custom tokens delimiter"
  workspace_id = "00000000-0000-0000-0000-000000000000"
  format       = "Default"
  definition = {
    "definition.json" = {
      source           = "${local.path}/definition.json.tmpl"
      tokens_delimiter = "##"
      tokens = {
        "CATALOG_NAME"                       = "MyCatalogName"
        "DATABRICKS_WORKSPACE_CONNECTION_ID" = "00000000-0000-0000-0000-000000000000"
      }
    }
  }
}

# Example 3b - Item with parameters processing mode
resource "fabric_mirrored_azure_databricks_catalog" "example_parameters" {
  display_name = "example3b"
  description  = "example with parameters processing mode"
  workspace_id = "00000000-0000-0000-0000-000000000000"
  format       = "Default"
  definition = {
    "definition.json" = {
      source          = "${local.path}/definition.json.tmpl"
      processing_mode = "parameters"
      parameters = [
        {
          type  = "JsonPathReplace"
          find  = "$.catalogName"
          value = "UpdatedName"
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

# Example 4 - Item with configuration, no definition - configuration and definition cannot be used together at the same time
resource "fabric_mirrored_azure_databricks_catalog" "example_configuration" {
  display_name = "example4"
  description  = "example with configuration"
  workspace_id = "00000000-0000-0000-0000-000000000000"
  configuration = {
    catalog_name                       = "MyCatalogName",
    mirroring_mode                     = "Partial"
    databricks_workspace_connection_id = "00000000-0000-0000-0000-000000000000"
    storage_connection_id              = "11111111-1111-1111-1111-111111111111"
  }
}
