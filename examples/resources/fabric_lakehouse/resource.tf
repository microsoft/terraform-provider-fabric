# Simple Lakehouse resource
resource "fabric_lakehouse" "example1" {
  display_name = "example1"
  workspace_id = "00000000-0000-0000-0000-000000000000"
}

# Lakehouse resource with enabled schemas
resource "fabric_lakehouse" "example2" {
  display_name = "example2"
  description  = "example2 with enabled schemas"
  workspace_id = "00000000-0000-0000-0000-000000000000"

  configuration = {
    enable_schemas = true
  }
}


# Item with definition bootstrapping only
resource "fabric_lakehouse" "example_definition_bootstrap" {
  display_name              = "example2"
  description               = "example with definition bootstrapping"
  workspace_id              = "00000000-0000-0000-0000-000000000000"
  format                    = "Default"
  definition_update_enabled = false # <-- Disable definition update
  definition = {
    "lakehouse.metadata.json" = {
      source = "${local.path}/lakehouse.metadata.json.tmpl"
    }
  }
}

# Item with definition
resource "fabric_lakehouse" "example_definition_update" {
  display_name = "example3"
  description  = "example with definition"
  workspace_id = "00000000-0000-0000-0000-000000000000"
  format       = "Default"
  definition = {
    "lakehouse.metadata.json" = {
      source = "${local.path}/lakehouse.metadata.json.tmpl"
    }
    "shortcuts.metadata.json" = {
      source = "${local.path}/shortcuts.metadata.json.tmpl"
    }
    "data-access-roles.json" = {
      source = "${local.path}/data-access-roles.json.tmpl"
    }
    "alm.settings.json" = {
      source = "${local.path}/alm.settings.json.tmpl"
    }
  }
}

# Item with custom tokens delimiter
resource "fabric_lakehouse" "example_custom_delimiter" {
  display_name = "example3a"
  description  = "example with custom tokens delimiter"
  workspace_id = "00000000-0000-0000-0000-000000000000"
  format       = "Default"
  definition = {
    "lakehouse.metadata.json" = {
      source           = "${local.path}/lakehouse.metadata.json.tmpl"
      tokens_delimiter = "##"
      tokens = {
        "MyKey" = "MyValue"
      }
    }
  }
}

# Example - Item with parameters processing mode
resource "fabric_lakehouse" "example_parameters" {
  display_name = "example3a"
  description  = "example with parameters processing mode"
  workspace_id = "00000000-0000-0000-0000-000000000000"
  format       = "Default"
  definition = {
    "lakehouse.metadata.json" = {
      source          = "${local.path}/lakehouse.metadata.json.tmpl"
      processing_mode = "Parameters"
      parameters = [
        {
          type  = "TextReplace"
          find  = "DefaultSchema"
          value = "dbo"
        }
      ]
    },
    "shortcuts.metadata.json" = {
      source          = "${local.path}/shortcuts.metadata.json.tmpl"
      processing_mode = "Parameters"
      parameters = [
        {
          type  = "JsonPathReplace"
          find  = "$.target.oneLake.itemID"
          value = "11111111-1111-1111-1111-111111111111"
        }
      ]
    }
  }
}
