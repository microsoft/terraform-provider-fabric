# Example 1 - Item with definition bootstrapping only
resource "fabric_digital_twin_builder_flow" "example_definition_bootstrap" {
  display_name              = "example"
  description               = "example with definition bootstrapping"
  workspace_id              = "00000000-0000-0000-0000-000000000000"
  definition_update_enabled = false
  format                    = "Default"
  definition = {
    "definition.json" = {
      source = "${local.path}/definition.json.tmpl"
    }
  }
}

# Example 2 - Item with definition update when source or tokens changed
resource "fabric_digital_twin_builder_flow" "example_definition_update" {
  display_name = "example"
  description  = "example with definition update when source or tokens changed"
  workspace_id = "00000000-0000-0000-0000-000000000000"
  format       = "Default"
  definition = {
    "definition.json" = {
      source = "${local.path}/definition.json.tmpl"
      tokens = {
        "DIGITAL_TWIN_BUILDER_ID" = "11111111-1111-1111-1111-111111111111"
      }
    }
  }
}

#Example 3 - Item with creation payload
resource "fabric_digital_twin_builder_flow" "example_creation_payload" {
  display_name = "example"
  description  = "example with creation payload"
  workspace_id = "00000000-0000-0000-0000-000000000000"
  configuration = {
    digital_twin_builder_item_reference = {
      workspace_id   = "00000000-0000-0000-0000-000000000000",
      reference_type = "ById",
      item_id        = "11111111-1111-1111-1111-111111111111",
    }
  }
}
