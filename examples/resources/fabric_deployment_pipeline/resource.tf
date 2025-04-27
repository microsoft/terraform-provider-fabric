resource "fabric_deployment_pipeline" "example" {
  display_name = "Deployment Pipeline Example"
  description  = "Deployment Pipeline Example"
  id           = data.fabric_deployment_pipeline.example.id
  stages = [
    {
      display_name = "Stage 1",
      description  = "Stage 1",
      is_public    = true,
    },
    {
      display_name = "Stage 2",
      description  = "Stage 2",
      is_public    = false,
    }
  ]
}
