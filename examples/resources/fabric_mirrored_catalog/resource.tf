# Example 1 - Item with definition bootstrapping only
resource "fabric_mirrored_catalog" "example_definition_bootstrap" {
  display_name              = "example"
  description               = "example with definition bootstrapping"
  workspace_id              = "00000000-0000-0000-0000-000000000000"
  definition_update_enabled = false
  format                    = "Default"
  definition = {
    "MirroredCatalogDefinition.json" = {
      source = "${local.path}/MirroredCatalogDefinition.json.tmpl"
      tokens = {
        "CONNECTION_ID" = "11111111-1111-1111-1111-111111111111"
        "SCOPE"         = "default"
      }
    }
  }
}

# Example 2 - Item with definition update when source or tokens changed
resource "fabric_mirrored_catalog" "example_definition_update" {
  display_name = "example2"
  description  = "example with definition update when source or tokens changed"
  workspace_id = "00000000-0000-0000-0000-000000000000"
  format       = "Default"
  definition = {
    "MirroredCatalogDefinition.json" = {
      source = "${local.path}/MirroredCatalogDefinition.json.tmpl"
      tokens = {
        "CONNECTION_ID" = "11111111-1111-1111-1111-111111111111"
        "SCOPE"         = "default"
      }
    }
  }
}
