# Get item details
data "fabric_semantic_model" "example" {
  id           = "11111111-1111-1111-1111-111111111111"
  workspace_id = "00000000-0000-0000-0000-000000000000"
}

# Get item details with definition
data "fabric_semantic_model" "example_definition" {
  id                = "11111111-1111-1111-1111-111111111111"
  workspace_id      = "00000000-0000-0000-0000-000000000000"
  output_definition = true
}

# Access the content of the definition with JSONPath expression
output "example_definition_pbism_jsonpath" {
  value = provider::fabric::content_decode(data.fabric_semantic_model.example_definition.definition["definition.pbism"].content, ".version")
}

# Access the content of the definition as JSON object
output "example_definition_bim_object" {
  value = provider::fabric::content_decode(data.fabric_semantic_model.example_definition.definition["model.bim"].content).model.tables[0].partitions
}
