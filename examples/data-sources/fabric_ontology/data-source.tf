# Get item details by name
data "fabric_ontology" "example_by_name" {
  display_name = "example"
  workspace_id = "00000000-0000-0000-0000-000000000000"
}

# Get item details by id
data "fabric_ontology" "example_by_id" {
  id           = "11111111-1111-1111-1111-111111111111"
  workspace_id = "00000000-0000-0000-0000-000000000000"
}
# Get item details with definition
# Examples uses `id` but `display_name` can be used as well
data "fabric_ontology" "example_definition" {
  id                = "11111111-1111-1111-1111-111111111111"
  workspace_id      = "00000000-0000-0000-0000-000000000000"
  format            = "TMDL"
  output_definition = true
}

# Decode TMDL as text, not as JSON.
output "example_database_tmdl" {
  value = provider::fabric::content_decode(data.fabric_ontology.example_definition.definition["database.tmdl"].content)
}

# This is an invalid data source
# Do not specify `id` and `display_name` in the same data source block
# data "fabric_ontology" "example" {
#   display_name = "example"
#   id           = "11111111-1111-1111-1111-111111111111"
#   workspace_id = "00000000-0000-0000-0000-000000000000"
# }
