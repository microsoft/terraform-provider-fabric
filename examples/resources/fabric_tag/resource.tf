# Create example with default scope
resource "fabric_tag" "example" {
  display_name = "example"
}

# Create example with explicit scope
resource "fabric_tag" "example_scope" {
  display_name = "example_scope"
  scope = {
    type = "Tenant"
  }
}

#Update
resource "fabric_tag" "example_update" {
  id           = "00000000-0000-0000-0000-000000000000"
  display_name = "example_updated"
  scope = {
    type = "Tenant"
  }
}
