resource "fabric_item_job" "example" {
  workspace_id = "00000000-0000-0000-0000-000000000000"
  item_id      = "11111111-1111-1111-1111-111111111111"
  job_type     = "Execute"
}

# Trigger a Data Pipeline on-demand with parameters
resource "fabric_item_job" "with_parameters_example" {
  workspace_id = "00000000-0000-0000-0000-000000000000"
  item_id      = "11111111-1111-1111-1111-111111111111"
  job_type     = "Execute"
  execution_data = jsonencode({
    parameters = {
      param_name = "param_value"
    }
  })
}
