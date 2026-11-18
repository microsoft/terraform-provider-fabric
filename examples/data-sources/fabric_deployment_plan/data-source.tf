# Get item details by name
data "fabric_deployment_plan" "example_by_name" {
  display_name = "example"
  workspace_id = "00000000-0000-0000-0000-000000000000"
}

# Get item details by id
data "fabric_deployment_plan" "example_by_id" {
  id           = "11111111-1111-1111-1111-111111111111"
  workspace_id = "00000000-0000-0000-0000-000000000000"
}

# Get item details with definition
# Examples uses `id` but `display_name` can be used as well
data "fabric_deployment_plan" "example_definition" {
  id                = "11111111-1111-1111-1111-111111111111"
  workspace_id      = "00000000-0000-0000-0000-000000000000"
  format            = "Default"
  output_definition = true
}

# Access the content of the definition (plan.yml is YAML, so content_decode returns a string)
output "example_definition_content" {
  value = provider::fabric::content_decode(data.fabric_deployment_plan.example_definition.definition["plan.yml"].content)
}

# Access the content of the definition as an object
output "example_definition_content_object" {
  value = yamldecode(provider::fabric::content_decode(data.fabric_deployment_plan.example_definition.definition["plan.yml"].content)).groups
}

# This is an invalid data source
# Do not specify `id` and `display_name` in the same data source block
# data "fabric_deployment_plan" "example" {
#   display_name = "example"
#   id           = "11111111-1111-1111-1111-111111111111"
#   workspace_id = "00000000-0000-0000-0000-000000000000"
# }
