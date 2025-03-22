# Example 1 - item with definition bootstrapping only
resource "mounted_data_factory" "example_definition_bootstrap" {
  display_name              = "example"
  description               = "example with definition bootstrapping"
  workspace_id              = "00000000-0000-0000-0000-000000000000"
  format                    = "Default"
  definition_update_enabled = false
  definition = {
    "mountedDataFactory-content.json" = {
      source = "${local.path}/mountedDataFactory-content.json.tmpl"
      tokens = {
        "DataFactoryResourceId" = "/subscriptions/<subscriptionId>/resourceGroups/<resourceGroup>/providers/Microsoft.DataFactory/factories/<factoryName>"
      }
    }
  }
}

# Example 2 - item with definition update when source or tokens changed
resource "mounted_data_factory" "example_definition_update" {
  display_name = "example"
  description  = "example with definition update when source or tokens changed"
  workspace_id = "00000000-0000-0000-0000-000000000000"
  format       = "Default"
  definition = {
    "mountedDataFactory-content.json" = {
      source = "${local.path}/mountedDataFactory-content.json.tmpl"
      tokens = {
        "DataFactoryResourceId" = "/subscriptions/<subscriptionId>/resourceGroups/<resourceGroup>/providers/Microsoft.DataFactory/factories/<factoryName>"
      }
    }
  }
}
