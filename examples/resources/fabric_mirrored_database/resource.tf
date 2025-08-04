# Example 1 - Item without definition
resource "fabric_mirrored_database" "example" {
  display_name = "example1"
  workspace_id = "00000000-0000-0000-0000-000000000000"
}

# Example 2 - Item with definition bootstrapping only
resource "fabric_mirrored_database" "example_definition_bootstrap" {
  display_name              = "example2"
  description               = "example with definition bootstrapping"
  workspace_id              = "00000000-0000-0000-0000-000000000000"
  definition_update_enabled = false
  format                    = "Default"
  definition = {
    "mirroring.json" = {
      source = "${local.path}/mirroring.json.tmpl"
    }
  }
}

# Example 3 - Item with definition update when source or tokens changed
resource "fabric_mirrored_database" "example_definition_update" {
  display_name = "example3"
  description  = "example with definition update when source or tokens changed"
  workspace_id = "00000000-0000-0000-0000-000000000000"
  format       = "Default"
  definition = {
    "mirroring.json" = {
      source = "${local.path}/mirroring.json.tmpl"
      tokens = {
        "DEFAULT_SCHEMA" = "my_schema"
      }
    }
  }
}
