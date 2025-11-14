# Example 1 - Item with definition bootstrapping only
resource "fabric_mounted_data_factory" "example_definition_bootstrap" {
  display_name              = "example1"
  description               = "example with definition bootstrapping"
  workspace_id              = "00000000-0000-0000-0000-000000000000"
  folder_id                 = "11111111-1111-1111-1111-111111111111"
  definition_update_enabled = false
  format                    = "Default"
  definition = {
    "mountedDataFactory-content.json" = {
      source = "${local.path}/mountedDataFactory-content.json.tmpl"
    }
  }
}

# Example 2 - Item with definition update when source or tokens changed
resource "fabric_mounted_data_factory" "example_definition_update" {
  display_name = "example2"
  description  = "example with definition update when source or tokens changed"
  workspace_id = "00000000-0000-0000-0000-000000000000"
  format       = "Default"
  definition = {
    "mountedDataFactory-content.json" = {
      source = "${local.path}/mountedDataFactory-content.json.tmpl"
      tokens = {
        "SUBSCRIPTION_ID"     = "00000000-0000-0000-0000-000000000000"
        "RESOURCE_GROUP_NAME" = "MyResourceGroupName"
        "FACTORY_NAME"        = "MyFactoryName"
      }
    }
  }
}
