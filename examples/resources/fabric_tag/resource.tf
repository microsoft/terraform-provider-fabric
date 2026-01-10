# Create example
resource "fabric_tag" "example" {
  tags = [
    {
      display_name = "example"
    },
    {
      display_name = "example2"
    }
  ]
}

#Update
resource "fabric_tag" "example_update" {
  id           = "00000000-0000-0000-0000-000000000000"
  display_name = "example_updated"
  scope = {
    type = "Tenant"
  }
}
