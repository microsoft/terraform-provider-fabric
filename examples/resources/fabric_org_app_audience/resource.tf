# Notes:
# - The `definition` is required. The `parentAppId` in the definition must reference an existing Org App in the same workspace.
# - The `parentAppId` cannot be changed after creation. Use `lifecycle.replace_triggered_by` to recreate the Org App Audience when the parent Org App is replaced.
# - Deleting the parent Org App also deletes all of its Org App Audiences.
# - An Org App Audience can't be moved to another folder on its own; move the parent Org App instead.

# Parent Org App
resource "fabric_org_app" "example" {
  display_name = "example"
  workspace_id = "00000000-0000-0000-0000-000000000000"
}

# Example 1 - Item with definition bootstrapping only
resource "fabric_org_app_audience" "example" {
  display_name              = "example1"
  description               = "example with definition bootstrapping"
  workspace_id              = "00000000-0000-0000-0000-000000000000"
  definition_update_enabled = false
  format                    = "Default"
  definition = {
    "definition.json" = {
      source = "${local.path}/definition.json.tmpl"
      tokens = {
        "ParentAppId"              = fabric_org_app.example.id
        "HasAccessToHiddenContent" = "false"
        "TabOrder"                 = "1"
      }
    }
  }

  lifecycle {
    replace_triggered_by = [fabric_org_app.example.id]
  }
}

# Example 2 - Item with definition update when source or tokens changed
resource "fabric_org_app_audience" "example_definition_update" {
  display_name = "example2"
  description  = "example with definition update when source or tokens changed"
  workspace_id = "00000000-0000-0000-0000-000000000000"
  format       = "Default"
  definition = {
    "definition.json" = {
      source = "${local.path}/definition.json.tmpl"
      tokens = {
        "ParentAppId"              = fabric_org_app.example.id
        "HasAccessToHiddenContent" = "true"
        "TabOrder"                 = "2"
      }
    }
  }

  lifecycle {
    replace_triggered_by = [fabric_org_app.example.id]
  }
}
