data "fabric_dataflow" "example_by_id" {
  id           = "11111111-1111-1111-1111-111111111111"
  workspace_id = "00000000-0000-0000-0000-000000000000"
}

data "fabric_dataflow" "example_by_name" {
  display_name = "example"
  workspace_id = "00000000-0000-0000-0000-000000000000"
}

# Get item details with definition
# Examples uses `id` but `display_name` can be used as well
data "fabric_dataflow" "example_definition" {
  id                = "11111111-1111-1111-1111-111111111111"
  workspace_id      = "00000000-0000-0000-0000-000000000000"
  format            = "Default"
  output_definition = true
}

# Access the content of the definition with JSONPath expression
output "example_definition_content_jsonpath" {
  value = provider::fabric::content_decode(data.fabric_dataflow.example_definition.definition["queryMetadata.json"].content, ".")
}

# Access the content of the definition as JSON object
output "example_definition_content_json_object" {
  value = provider::fabric::content_decode(data.fabric_dataflow.example_definition.definition["queryMetadata.json"].content)
}
output "example_definition_pq_object" {
  value = provider::fabric::content_decode(data.fabric_dataflow.example_definition.definition["mashup.pq"].content)
}

# This is an invalid data source
# Do not specify `id` and `display_name` in the same data source block
# data "fabric_dataflow" "example" {
#   display_name = "dp_example"
#   id           = "11111111-1111-1111-1111-111111111111"
#   workspace_id = "00000000-0000-0000-0000-000000000000"
# }
