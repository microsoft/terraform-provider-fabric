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
