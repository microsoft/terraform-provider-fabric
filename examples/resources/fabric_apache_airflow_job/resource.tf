# Example 1 - Item without definition
resource "fabric_apache_airflow_job" "example_definition" {
  display_name = "example"
  description  = "example without definition"
  workspace_id = "00000000-0000-0000-0000-000000000000"
  folder_id    = "11111111-1111-1111-1111-111111111111"
}

# Example 2 - Item with definition bootstrapping only
resource "fabric_apache_airflow_job" "example_definition_bootstrap" {
  display_name              = "example"
  description               = "example with definition bootstrapping"
  workspace_id              = "00000000-0000-0000-0000-000000000000"
  definition_update_enabled = false
  definition = {
    "apacheairflowjob-content.json" = {
      source = "${local.path}/apacheairflowjob-content.json.tmpl"
    }
  }
}

# Example 3 - Item with definition update when source or tokens changed
resource "fabric_apache_airflow_job" "example_definition_update" {
  display_name = "example"
  description  = "example with definition update when source or tokens changed"
  workspace_id = "00000000-0000-0000-0000-000000000000"
  definition = {
    "apacheairflowjob-content.json" = {
      source = "${local.path}/apacheairflowjob-content.json.tmpl"
      tokens = {}
    }
  }
}
