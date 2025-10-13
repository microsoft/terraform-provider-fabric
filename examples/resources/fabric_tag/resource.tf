#Example with default scope (Tenant)
resource "fabric_tag" "example" {
  create_tags_request = [
    {
      display_name = "example"
    },
    {
      display_name = "example2"
    }
  ]
}
