# Get item details
data "fabric_report" "example" {
  id           = "11111111-1111-1111-1111-111111111111"
  workspace_id = "00000000-0000-0000-0000-000000000000"
}

# Get item details with definition
data "fabric_report" "example_definition" {
  id                = "11111111-1111-1111-1111-111111111111"
  workspace_id      = "00000000-0000-0000-0000-000000000000"
  format            = "PBIR-Legacy"
  output_definition = true
}

# Access the content of the definition with JSONPath expression
output "example_definition_pbir_jsonpath" {
  value = provider::fabric::content_decode(data.fabric_report.example_definition.definition["definition.pbir"].content, ".datasetReference.byConnection.connectionString")
}
# Access the content of the definition as JSON object
output "example_definition_pbir_object" {
  value = provider::fabric::content_decode(data.fabric_report.example_definition.definition["definition.pbir"].content).datasetReference.byConnection.connectionString
}

output "example_definition_report_object" {
  value = provider::fabric::content_decode(data.fabric_report.example_definition.definition["report.json"].content)
}
