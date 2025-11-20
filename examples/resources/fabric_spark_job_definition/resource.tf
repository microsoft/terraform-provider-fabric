# Example 1 - Item without definition
resource "fabric_spark_job_definition" "example" {
  display_name = "example1"
  workspace_id = "00000000-0000-0000-0000-000000000000"
}

# Example 2 - Item with definition bootstrapping only
resource "fabric_spark_job_definition" "example_definition_bootstrap" {
  display_name              = "example2"
  description               = "example with definition bootstrapping"
  workspace_id              = "00000000-0000-0000-0000-000000000000"
  definition_update_enabled = false
  format                    = "SparkJobDefinitionV1"
  definition = {
    "SparkJobDefinitionV1.json" = {
      source = "${local.path}/SparkJobDefinitionV1.json.tmpl"
    }
  }
}

# Example 3 - Item with definition update when source or tokens changed
resource "fabric_spark_job_definition" "example_definition_update" {
  display_name = "example3"
  description  = "example with definition update when source or tokens changed"
  workspace_id = "00000000-0000-0000-0000-000000000000"
  format       = "SparkJobDefinitionV1"
  definition = {
    "SparkJobDefinitionV1.json" = {
      source = "${local.path}/SparkJobDefinitionV1.json.tmpl"
      tokens = {
        "DefaultLakehouseID"     = "11111111-1111-1111-1111-111111111111"
        "AdditionalLakehouseIDs" = "\"22222222-2222-2222-2222-222222222222\",\"33333333-3333-3333-3333-333333333333\""
      }
    }
  }
}

# Example 4 - Item with custom tokens delimiter
resource "fabric_spark_job_definition" "example_custom_delimiter" {
  display_name = "example4"
  description  = "example with custom tokens delimiter"
  workspace_id = "00000000-0000-0000-0000-000000000000"
  format       = "SparkJobDefinitionV1"
  definition = {
    "SparkJobDefinitionV1.json" = {
      source           = "${local.path}/SparkJobDefinitionV1.json.tmpl"
      tokens_delimiter = "##"
      tokens = {
        "DefaultLakehouseID"     = "11111111-1111-1111-1111-111111111111"
        "AdditionalLakehouseIDs" = "\"22222222-2222-2222-2222-222222222222\",\"33333333-3333-3333-3333-333333333333\""
      }
    }
  }
}

# Example 5 - Item with parameters processing mode
resource "fabric_spark_job_definition" "example_parameters" {
  display_name = "example5"
  description  = "example with parameters processing mode"
  workspace_id = "00000000-0000-0000-0000-000000000000"
  format       = "SparkJobDefinitionV1"
  definition = {
    "SparkJobDefinitionV1.json" = {
      source          = "${local.path}/SparkJobDefinitionV1.json.tmpl"
      processing_mode = "parameters"
      parameters = [
        {
          type  = "JsonPathReplace"
          find  = "$.defaultLakehouseArtifactId"
          value = "44444444-4444-4444-4444-444444444444"
        },
        {
          type  = "TextReplace"
          find  = "OldLakehouseID"
          value = "NewLakehouseID"
        }
      ]
    }
  }
}
