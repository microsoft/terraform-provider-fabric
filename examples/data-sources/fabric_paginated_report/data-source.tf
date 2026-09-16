# Get item details by id
data "fabric_paginated_report" "example_by_id" {
  id           = "11111111-1111-1111-1111-111111111111"
  workspace_id = "00000000-0000-0000-0000-000000000000"
}

# Get item details by name
data "fabric_paginated_report" "example_by_name" {
  display_name = "example"
  workspace_id = "00000000-0000-0000-0000-000000000000"
}

# Get item details with its RDL definition
# This example uses `id`, but `display_name` can be used instead.
data "fabric_paginated_report" "example_definition" {
  id                = "11111111-1111-1111-1111-111111111111"
  workspace_id      = "00000000-0000-0000-0000-000000000000"
  format            = "Default"
  output_definition = true
}

# Decode the RDL definition as an XML string
output "example_definition_content" {
  value = provider::fabric::content_decode(data.fabric_paginated_report.example_definition.definition["example.rdl"].content)
}

# This is an invalid data source.
# Do not specify `id` and `display_name` in the same data source block.
# data "fabric_paginated_report" "example" {
#   display_name = "example"
#   id           = "11111111-1111-1111-1111-111111111111"
#   workspace_id = "00000000-0000-0000-0000-000000000000"
# }
