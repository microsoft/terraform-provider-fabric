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
    "mirroringAzureDatabricksCatalog.json" = {
      source = "${local.path}/mirroringAzureDatabricksCatalog.json.tmpl"
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
    "mirroringAzureDatabricksCatalog.json" = {
      source = "${local.path}/mirroringAzureDatabricksCatalog.json.tmpl"
      tokens = {
        "MyKey" = "MyValue"
      }
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
    mirroring_mode                     = "Partial/Full"
    databricks_workspace_connection_id = "00000000-0000-0000-0000-000000000000"
    storage_connection_id              = "11111111-1111-1111-1111-111111111111"
  }
}
