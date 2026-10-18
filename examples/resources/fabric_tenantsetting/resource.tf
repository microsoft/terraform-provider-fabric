# Example 1 - Updating a Tenant Setting
resource "fabric_tenantsetting" "example" {
  setting_name     = "example"
  enabled          = false
  delete_behaviour = "NoChange"
}

# Example 2 - Updating a Tenant Setting with Security Groups
resource "fabric_tenantsetting" "example_with_security_groups" {
  setting_name     = "example"
  enabled          = true
  delete_behaviour = "NoChange"
  enabled_security_groups = [
    {
      graph_id = "00000000-0000-0000-0000-000000000000"
      name     = "example"
    }
  ]
}

# Example 3 - Updating a Tenant Setting with delete_behaviour (Setting will be set to disabled on delete)
resource "fabric_tenantsetting" "example_with_delete_behaviour" {
  setting_name     = "example"
  enabled          = true
  delete_behaviour = "Disable"
}

# Note: if delete_behaviour is not specified, it defaults to "NoChange"
