# Example 1 - Item without definition
resource "fabric_deployment_plan" "example" {
  display_name = "example"
  workspace_id = "00000000-0000-0000-0000-000000000000"
}

# Example 2 - Item with definition bootstrapping only
resource "fabric_deployment_plan" "example_definition_bootstrap" {
  display_name              = "example"
  description               = "example with definition bootstrapping"
  workspace_id              = "00000000-0000-0000-0000-000000000000"
  definition_update_enabled = false
  format                    = "Default"
  definition = {
    "plan.yml" = {
      source = "${local.path}/plan.yml.tmpl"
      tokens = {
        "GROUP_NAME" = "Group1"
        "LOGICAL_ID" = "11111111-1111-1111-1111-111111111111"
      }
    }
  }
}

# Example 3 - Item with definition update when source or tokens changed
resource "fabric_deployment_plan" "example_definition_update" {
  display_name = "example"
  description  = "example with definition update when source or tokens changed"
  workspace_id = "00000000-0000-0000-0000-000000000000"
  format       = "Default"
  definition = {
    "plan.yml" = {
      source = "${local.path}/plan.yml.tmpl"
      tokens = {
        "GROUP_NAME" = "Group1"
        "LOGICAL_ID" = "11111111-1111-1111-1111-111111111111"
      }
    }
  }
}

# Example 4 - Item with custom tokens delimiter
resource "fabric_deployment_plan" "example_custom_delimiter" {
  display_name = "example"
  description  = "example with custom tokens delimiter"
  workspace_id = "00000000-0000-0000-0000-000000000000"
  format       = "Default"
  definition = {
    "plan.yml" = {
      source           = "${local.path}/plan.yml.tmpl"
      tokens_delimiter = "{{}}"
      tokens = {
        "GROUP_NAME" = "Group1"
        "LOGICAL_ID" = "11111111-1111-1111-1111-111111111111"
      }
    }
  }
}
