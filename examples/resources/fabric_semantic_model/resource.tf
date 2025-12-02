# Semantic Model bootstrapping only
resource "fabric_semantic_model" "example_bootstrap" {
  display_name              = "example"
  workspace_id              = "00000000-0000-0000-0000-000000000000"
  definition_update_enabled = false
  format                    = "TMSL"
  definition = {
    "model.bim" = {
      source = "${local.path}/model.bim.tmpl"
    }
    "definition.pbism" = {
      source = "${local.path}/definition.pbism"
    }
  }
}

# Semantic Model with definition update when source or tokens changed
resource "fabric_semantic_model" "example_update" {
  display_name = "example with update"
  workspace_id = "00000000-0000-0000-0000-000000000000"
  format       = "TMSL"
  definition = {
    "model.bim" = {
      source = "${local.path}/model.bim.tmpl"
      tokens = {
        "ColumnName" = "Hello"
      }
    }
    "definition.pbism" = {
      source = "${local.path}/definition.pbism"
    }
  }
}

# Semantic Model with custom tokens delimiter
resource "fabric_semantic_model" "example_custom_delimiter" {
  display_name = "example with custom delimiter"
  workspace_id = "00000000-0000-0000-0000-000000000000"
  format       = "TMSL"
  definition = {
    "model.bim" = {
      source           = "${local.path}/model.bim.tmpl"
      tokens_delimiter = "##"
      tokens = {
        "ColumnName" = "Hello"
      }
    }
    "definition.pbism" = {
      source = "${local.path}/definition.pbism"
    }
  }
}

# Semantic Model with parameters processing mode
resource "fabric_semantic_model" "example_parameters" {
  display_name = "example with parameters"
  workspace_id = "00000000-0000-0000-0000-000000000000"
  format       = "TMSL"
  definition = {
    "model.bim" = {
      source          = "${local.path}/model.bim.tmpl"
      processing_mode = "parameters"
      parameters = [
        {
          type  = "JsonPathReplace"
          find  = "$.model.tables[0].columns[0].name"
          value = "UpdatedColumnName"
        },
        {
          type  = "TextReplace"
          find  = "OldValue"
          value = "NewValue"
        }
      ]
    }
    "definition.pbism" = {
      source = "${local.path}/definition.pbism"
    }
  }
}
