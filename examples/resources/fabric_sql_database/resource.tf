# Example 1 - Item without configuration or definition (default)
resource "fabric_sql_database" "example" {
  display_name = "example"
  workspace_id = "00000000-0000-0000-0000-000000000000"
}

# Example 2 - Item with configuration - create a new SQL database with custom settings
resource "fabric_sql_database" "example_new" {
  display_name = "example_new"
  workspace_id = "00000000-0000-0000-0000-000000000000"

  configuration = {
    creation_mode         = "New"
    backup_retention_days = 10
    collation             = "SQL_Latin1_General_CP1_CI_AS"
  }
}

# Example 3 - Item with configuration - restore a SQL database by ID
resource "fabric_sql_database" "example_restore_by_id" {
  display_name = "example_restore_by_id"
  workspace_id = "00000000-0000-0000-0000-000000000000"

  configuration = {
    creation_mode         = "Restore"
    restore_point_in_time = "2026-01-01T00:00:00Z"
    source_database_reference = {
      item_id        = "11111111-1111-1111-1111-111111111111"
      reference_type = "ById"
      workspace_id   = "00000000-0000-0000-0000-000000000000"
    }
  }
}

# Example 4 - Item with configuration - restore a SQL database by variable reference
resource "fabric_sql_database" "example_restore_by_variable" {
  display_name = "example_restore_by_variable"
  workspace_id = "00000000-0000-0000-0000-000000000000"

  configuration = {
    creation_mode         = "Restore"
    restore_point_in_time = "2026-01-01T00:00:00Z"
    source_database_reference = {
      reference_type     = "ByVariable"
      variable_reference = "$(/**/_VarLibrary_/_VarName_)"
    }
  }
}

# Example 5 - Item with definition only - deploy a SQL project
resource "fabric_sql_database" "example_sqlproj" {
  display_name = "example_sqlproj"
  workspace_id = "00000000-0000-0000-0000-000000000000"
  format       = "sqlproj"

  definition = {
    "definition.sqlproj" = {
      source = "${local.path}/definition.sqlproj.tmpl"
    }
    "dbo/Tables/TestTable.sql" = {
      source = "${local.path}/dbo/Tables/TestTable.sql.tmpl"
    }
  }
}
