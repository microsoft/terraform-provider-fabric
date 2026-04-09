# Example 1 - Data Agent without definition
resource "fabric_data_agent" "example" {
  display_name = "example"
  workspace_id = "00000000-0000-0000-0000-000000000000"
}

# Example 2 - Data Agent with definition bootstrapping only
resource "fabric_data_agent" "example_definition_bootstrap" {
  display_name              = "example"
  description               = "example with definition bootstrapping"
  workspace_id              = "00000000-0000-0000-0000-000000000000"
  definition_update_enabled = false
  definition = {
    "Files/Config/data_agent.json" = {
      source = "${local.path}/data_agent.json.tmpl"
    }
    "Files/Config/draft/stage_config.json" = {
      source = "${local.path}/stage_config.json.tmpl"
    }
  }
}

# Example 3 - Data Agent with definition update when source changed
resource "fabric_data_agent" "example_definition_update" {
  display_name = "example"
  description  = "example with definition update when source changed"
  workspace_id = "00000000-0000-0000-0000-000000000000"
  definition = {
    "Files/Config/data_agent.json" = {
      source = "${local.path}/data_agent.json.tmpl"
    }
    "Files/Config/draft/stage_config.json" = {
      source = "${local.path}/stage_config.json.tmpl"
    }
  }
}
