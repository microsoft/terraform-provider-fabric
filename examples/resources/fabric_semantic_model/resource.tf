# Semantic Model bootstrapping only
resource "fabric_semantic_model" "example_bootstrap" {
  display_name              = "example"
  workspace_id              = "00000000-0000-0000-0000-000000000000"
  folder_id                 = "11111111-1111-1111-1111-111111111111"
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
