# Example 1 - Paginated Report with required definition
resource "fabric_paginated_report" "example" {
  display_name = "example"
  workspace_id = "00000000-0000-0000-0000-000000000000"
  format       = "Default"
  definition = {
    "example.rdl" = {
      source = "${local.path}/report.rdl"
    }
  }
}

# Example 2 - Paginated Report with definition bootstrapping only
resource "fabric_paginated_report" "example_definition_bootstrap" {
  display_name              = "example_definition_bootstrap"
  description               = "example with definition bootstrapping"
  workspace_id              = "00000000-0000-0000-0000-000000000000"
  definition_update_enabled = false
  format                    = "Default"
  definition = {
    "example_definition_bootstrap.rdl" = {
      source = "${local.path}/report.rdl"
    }
  }
}

# Example 3 - Paginated Report with definition update when source or tokens change
resource "fabric_paginated_report" "example_definition_update" {
  display_name = "example_definition_update"
  description  = "example with definition update when source or tokens change"
  workspace_id = "00000000-0000-0000-0000-000000000000"
  format       = "Default"
  definition = {
    "example_definition_update.rdl" = {
      source = "${local.path}/report.rdl.tmpl"
      tokens = {
        "ReportTitle" = "Updated Paginated Report"
      }
    }
  }
}

# Example 4 - Paginated Report with custom tokens delimiter
resource "fabric_paginated_report" "example_custom_delimiter" {
  display_name = "example_custom_delimiter"
  description  = "example with custom tokens delimiter"
  workspace_id = "00000000-0000-0000-0000-000000000000"
  format       = "Default"
  definition = {
    "example_custom_delimiter.rdl" = {
      source           = "${local.path}/report_at.rdl.tmpl"
      tokens_delimiter = "@{}@"
      tokens = {
        "ReportTitle" = "Custom Delimiter Paginated Report"
      }
    }
  }
}

# Example 5 - Paginated Report with parameters processing mode
resource "fabric_paginated_report" "example_parameters" {
  display_name = "example_parameters"
  description  = "example with parameters processing mode"
  workspace_id = "00000000-0000-0000-0000-000000000000"
  format       = "Default"
  definition = {
    "example_parameters.rdl" = {
      source          = "${local.path}/report.rdl"
      processing_mode = "parameters"
      parameters = [
        {
          type  = "TextReplace"
          find  = "Paginated Report"
          value = "Parameters Paginated Report"
        }
      ]
    }
  }
}
