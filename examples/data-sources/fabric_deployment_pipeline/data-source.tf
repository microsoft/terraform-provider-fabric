data "fabric_deployment_pipeline" "example_by_id" {
  id = "11111111-1111-1111-1111-111111111111"
}

data "fabric_deployment_pipeline" "example_by_name" {
  display_name = "example"
}
