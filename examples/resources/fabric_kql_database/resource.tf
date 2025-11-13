# Example 1 - Item with configuration, no definition - create a ReadWrite KQL database
resource "fabric_kql_database" "example1" {
  display_name = "example1"
  workspace_id = "00000000-0000-0000-0000-000000000000"

  configuration = {
    database_type = "ReadWrite"
    eventhouse_id = "11111111-1111-1111-1111-111111111111"
  }
}

# Example 2 - Item with configuration, no definition - create a Shortcut KQL database to source Azure Data Explorer cluster
resource "fabric_kql_database" "example2" {
  display_name = "example2"
  workspace_id = "00000000-0000-0000-0000-000000000000"

  configuration = {
    database_type        = "Shortcut"
    eventhouse_id        = "11111111-1111-1111-1111-111111111111"
    source_cluster_uri   = "https://clustername.westus.kusto.windows.net"
    source_database_name = "MyDatabase"
  }
}

# Example 3 - Item with configuration, no definition - create a Shortcut KQL database to source Azure Data Explorer cluster with invitation token
# Example below uses Write-only Arguments (https://developer.hashicorp.com/terraform/language/resources/ephemeral/write-only) and let you securely pass temporary values to Terraform's managed resources during an operation without persisting those values to state or plan files.
# Require Terraform 1.11 and later.
resource "fabric_kql_database" "example3" {
  display_name = "example3"
  workspace_id = "00000000-0000-0000-0000-000000000000"

  configuration = {
    database_type               = "Shortcut"
    eventhouse_id               = "11111111-1111-1111-1111-111111111111"
    invitation_token_wo         = "eyJ0...InvitationToken...iJKV"
    invitation_token_wo_version = 1
  }
}

# Example 4 - Item with configuration, no definition - create a Shortcut KQL database to source Azure Data Explorer cluster with invitation token
# Example below does NOT use Write-only Arguments and secret values are persisted to state and plan files. Recommended to use Write-only Arguments instead.
# Works on Terraform 1.10 and below.
resource "fabric_kql_database" "example4" {
  display_name = "example4"
  workspace_id = "00000000-0000-0000-0000-000000000000"

  configuration = {
    database_type    = "Shortcut"
    eventhouse_id    = "11111111-1111-1111-1111-111111111111"
    invitation_token = "eyJ0...InvitationToken...iJKV"
  }
}

# Example 5 - Item with configuration, no definition - create a Shortcut KQL database to source KQL database
resource "fabric_kql_database" "example5" {
  display_name = "example5"
  workspace_id = "00000000-0000-0000-0000-000000000000"

  configuration = {
    database_type        = "Shortcut"
    eventhouse_id        = "11111111-1111-1111-1111-111111111111"
    source_database_name = "MyDatabase"
  }
}


# Example 6 - Item with definition bootstrapping only
resource "fabric_kql_database" "example6" {
  display_name              = "example6"
  description               = "example with definition bootstrapping"
  workspace_id              = "00000000-0000-0000-0000-000000000000"
  format                    = "Default"
  definition_update_enabled = false # <-- Disable definition update
  definition = {
    "DatabaseProperties.json" = {
      source = "${local.path}/DatabaseProperties.json.tmpl"
    }
  }
}

# Example 7 - Item with definition update when source or tokens changed
resource "fabric_kql_database" "example7" {
  display_name = "example7"
  description  = "example with definition update when source or tokens changed"
  workspace_id = "00000000-0000-0000-0000-000000000000"
  format       = "Default"
  definition = {
    "DatabaseProperties.json" = {
      source = "${local.path}/DatabaseProperties.json.tmpl"
      tokens = {
        "MyKey" = "MyValue"
      }
    }
  }
}

# Example 8 - Item with custom tokens delimiter
resource "fabric_kql_database" "example8" {
  display_name = "example8"
  description  = "example with custom tokens delimiter"
  workspace_id = "00000000-0000-0000-0000-000000000000"
  format       = "Default"
  definition = {
    "DatabaseProperties.json" = {
      source           = "${local.path}/DatabaseProperties.json.tmpl"
      tokens_delimiter = "##"
      tokens = {
        "MyKey" = "MyValue"
      }
    }
  }
}

# Example 9 - Item with parameters processing mode
resource "fabric_kql_database" "example9" {
  display_name = "example9"
  description  = "example with parameters processing mode"
  workspace_id = "00000000-0000-0000-0000-000000000000"
  format       = "Default"
  definition = {
    "DatabaseProperties.json" = {
      source          = "${local.path}/DatabaseProperties.json.tmpl"
      processing_mode = "parameters"
      parameters = [
        {
          type  = "JsonPathReplace"
          find  = "$.properties.databaseType"
          value = "ReadWrite"
        },
        {
          type  = "TextReplace"
          find  = "OldValue"
          value = "NewValue"
        }
      ]
    }
  }
}
