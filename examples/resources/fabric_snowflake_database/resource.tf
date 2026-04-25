# Example 1 - Item without configuration or definition (empty database)
resource "fabric_snowflake_database" "example" {
  display_name = "example1"
  workspace_id = "00000000-0000-0000-0000-000000000000"
}

# Example 2 - Item with configuration (binds to an existing Snowflake connection)
resource "fabric_snowflake_database" "example_configuration" {
  display_name = "example2"
  description  = "example with configuration"
  workspace_id = "00000000-0000-0000-0000-000000000000"
  configuration = {
    connection_id           = "11111111-1111-1111-1111-111111111111"
    snowflake_database_name = "ExampleDatabase"
  }
}

# Example 3 - Item with definition bootstrapping only
resource "fabric_snowflake_database" "example_definition_bootstrap" {
  display_name              = "example3"
  description               = "example with definition bootstrapping"
  workspace_id              = "00000000-0000-0000-0000-000000000000"
  format                    = "Default"
  definition_update_enabled = false # <-- Disable definition update
  definition = {
    "SnowflakeDatabaseProperties.json" = {
      source = "${local.path}/SnowflakeDatabaseProperties.json.tmpl"
    }
  }
}

# Example 4 - Item with definition update when source or tokens changed
resource "fabric_snowflake_database" "example_definition_update" {
  display_name = "example4"
  description  = "example with definition update when source or tokens changed"
  workspace_id = "00000000-0000-0000-0000-000000000000"
  format       = "Default"
  definition = {
    "SnowflakeDatabaseProperties.json" = {
      source = "${local.path}/SnowflakeDatabaseProperties.json.tmpl"
      tokens = {
        "DATABASE_NAME" = "ExampleDatabase"
      }
    }
  }
}


# Example 5 - Item with custom tokens delimiter
resource "fabric_snowflake_database" "example_custom_delimiter" {
  display_name = "example"
  description  = "example with custom tokens delimiter"
  workspace_id = "00000000-0000-0000-0000-000000000000"
  format       = "Default"
  definition = {
    "SnowflakeDatabaseProperties.json" = {
      source           = "${local.path}/SnowflakeDatabaseProperties.json.tmpl"
      tokens_delimiter = "{{}}"
      tokens = {
        "DATABASE_NAME" = "ExampleDatabase"
      }
    }
  }
}

# Example 5 - Item with parameters processing mode
resource "fabric_snowflake_database" "example_parameters" {
  display_name = "example"
  description  = "example with parameters processing mode"
  workspace_id = "00000000-0000-0000-0000-000000000000"
  format       = "Default"
  definition = {
    "SnowflakeDatabaseProperties.json" = {
      source          = "${local.path}/SnowflakeDatabaseProperties.json.tmpl"
      processing_mode = "Parameters"
      parameters = [
        {
          type  = "TextReplace"
          find  = "OldDatabaseName"
          value = "NewDatabaseName"
        }
      ]
    }
  }
}
