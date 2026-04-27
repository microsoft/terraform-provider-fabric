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
      tokens = {
        "SCHEMA" = "2.1.0"
      }
    }
    "Files/Config/draft/stage_config.json" = {
      source = "${local.path}/stage_config.json.tmpl"
      tokens = {
        "SCHEMA" = "1.0.0"
      }
    }
  }
}


# Example 4 - Item with custom tokens delimiter
resource "fabric_data_agent" "example_custom_delimiter" {
  display_name = "example"
  description  = "example with custom tokens delimiter"
  workspace_id = "00000000-0000-0000-0000-000000000000"
  format       = "Default"
  definition = {
    "Files/Config/data_agent.json" = {
      source           = "${local.path}/data_agent.json.tmpl"
      tokens_delimiter = "{{}}"
      tokens = {
        "SCHEMA" = "2.1.0"
      }
    }
    "Files/Config/draft/stage_config.json" = {
      source           = "${local.path}/stage_config.json.tmpl"
      tokens_delimiter = "{{}}"
      tokens = {
        "SCHEMA" = "1.0.0"
      }
    }
  }
}

# Example 5 - Item with parameters processing mode
resource "fabric_data_agent" "example_parameters" {
  display_name = "example"
  description  = "example with parameters processing mode"
  workspace_id = "00000000-0000-0000-0000-000000000000"
  format       = "Default"
  definition = {
    "Files/Config/data_agent.json" = {
      source          = "${local.path}/data_agent.json.tmpl"
      processing_mode = "Parameters"
      parameters = [
        {
          type  = "TextReplace"
          find  = "SCHEMA"
          value = "2.1.0"
        }
      ]
    }
    "Files/Config/draft/stage_config.json" = {
      source          = "${local.path}/stage_config.json.tmpl"
      processing_mode = "Parameters"
      parameters = [
        {
          type  = "TextReplace"
          find  = "SCHEMA"
          value = "1.0.0"
        }
      ]
    }
  }
}

