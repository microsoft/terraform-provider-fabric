# Report bootstrapping only
resource "fabric_report" "example_bootstrap" {
  display_name              = "example"
  workspace_id              = "00000000-0000-0000-0000-000000000000"
  definition_update_enabled = false
  definition = {
    "report.json" = {
      source = "${local.path}/report.json"
    }
    "definition.pbir" = {
      source = "${local.path}/definition.pbir.tmpl"
      tokens = {
        "SemanticModelID" = "00000000-0000-0000-0000-000000000000"
      }
    }
    "StaticResources/SharedResources/BaseThemes/CY24SU06.json" = {
      source = "${local.path}/StaticResources/SharedResources/BaseThemes/CY24SU06.json"
    }
  }
}

# Report with update when source or tokens changed
resource "fabric_report" "example_update" {
  display_name = "example with update"
  workspace_id = "00000000-0000-0000-0000-000000000000"
  definition = {
    "report.json" = {
      source = "${local.path}/report.json"
    }
    "definition.pbir" = {
      source = "${local.path}/definition.pbir.tmpl"
      tokens = {
        "SemanticModelID" = "00000000-0000-0000-0000-000000000000"
      }
    }
    "StaticResources/SharedResources/BaseThemes/CY24SU06.json" = {
      source = "${local.path}/StaticResources/SharedResources/BaseThemes/CY24SU06.json"
    }
  }
}
