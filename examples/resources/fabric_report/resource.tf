resource "fabric_report" "example_pbir" {
  display_name = "test"
  workspace_id = "00000000-0000-0000-0000-000000000000"
  format       = "PBIR"
  definition = {
    "definition/report.json" = {
      source = "${local.path}/definition/report.json"
    }
    "definition/version.json" = {
      source = "${local.path}/definition/version.json"
    }
    "definition.pbir" = {
      source = "${local.path}/definition.pbir"
    }
    "definition/pages/pages.json" = {
      source = "${local.path}/definition/pages/pages.json"
    }
    "definition/pages/f0275333137c0ea79df2/page.json" = {
      source = "${local.path}/definition/pages/f0275333137c0ea79df2/page.json"
    }
  }
}

resource "fabric_report" "example_pbir_with_visuals" {
  display_name = "test"
  workspace_id = "00000000-0000-0000-0000-000000000000"
  format       = "PBIR"
  definition = {
    "definition/report.json" = {
      source = "${local.path}/definition/report.json"
    }
    "definition/version.json" = {
      source = "${local.path}/definition/version.json"
    }
    "definition.pbir" = {
      source = "${local.path}/definition.pbir"
    }
    "definition/pages/pages.json" = {
      source = "${local.path}/definition/pages/pages.json"
    }
    "definition/pages/f0275333137c0ea79df2/page.json" = {
      source = "${local.path}/definition/pages/f0275333137c0ea79df2/page.json"
    }
    "definition/pages/f0275333137c0ea79df2/visuals/a3c8f5e1b79d42f0c6a1/visual.json" = {
      source = "${local.path}/definition/pages/f0275333137c0ea79df2/visuals/a3c8f5e1b79d42f0c6a1/visual.json"
    }
    "StaticResources/SharedResources/BaseThemes/CY23SU10.json" = {
      source = "${local.path}/StaticResources/SharedResources/BaseThemes/CY23SU10.json"
    }
  }
}

# Report bootstrapping only
resource "fabric_report" "example_bootstrap" {
  display_name              = "example"
  workspace_id              = "00000000-0000-0000-0000-000000000000"
  definition_update_enabled = false
  format                    = "PBIR-Legacy"
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
    "StaticResources/SharedResources/BaseThemes/CY24SU10.json" = {
      source = "${local.path}/StaticResources/SharedResources/BaseThemes/CY24SU10.json"
    }
    "StaticResources/RegisteredResources/fabric_48_color10148978481469717.png" = {
      source = "${local.path}/StaticResources/RegisteredResources/fabric_48_color10148978481469717.png"
    }
  }
}

# Report with update when source or tokens changed
resource "fabric_report" "example_update" {
  display_name = "example with update"
  workspace_id = "00000000-0000-0000-0000-000000000000"
  format       = "PBIR-Legacy"
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
    "StaticResources/SharedResources/BaseThemes/CY24SU10.json" = {
      source = "${local.path}/StaticResources/SharedResources/BaseThemes/CY24SU10.json"
    }
  }
}
