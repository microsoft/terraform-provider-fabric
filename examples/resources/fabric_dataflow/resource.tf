# Item without definition
resource "fabric_dataflow" "example" {
  display_name = "example"
  workspace_id = "00000000-0000-0000-0000-000000000000"
}

# Dataflow bootstrapping only
resource "fabric_dataflow" "example_bootstrap" {
  display_name              = "example"
  workspace_id              = "00000000-0000-0000-0000-000000000000"
  definition_update_enabled = false
  format                    = "Default"
  definition = {
    "queryMetadata.json" = {
      source = "${local.path}/queryMetadata.json.tmpl"
    }
    "mashup.pq" = {
      source = "${local.path}/mashup.pq.tmpl"
    }
  }
}

# Dataflow with definition update when source or tokens changed
resource "fabric_dataflow" "example_update" {
  display_name = "example with update"
  workspace_id = "00000000-0000-0000-0000-000000000000"
  format       = "Default"
  definition = {
    "queryMetadata.json" = {
      source = "${local.path}/queryMetadata.json.tmpl"
      tokens = {
        "MyValue1" = "my value 1"
        "MyValue2" = "my value 2"
      }
    }
    "mashup.pq" = {
      source = "${local.path}/mashup.pq.tmpl"
    }
  }
}
