data "fabric_onelake_data_access_security" "example_by_name" {
  workspace_id = "00000000-0000-0000-0000-000000000000"
  item_id      = "11111111-1111-1111-1111-111111111111"
  role_name    = "example"
}

data "fabric_onelake_data_access_security" "example_by_id" {
  workspace_id = "00000000-0000-0000-0000-000000000000"
  item_id      = "11111111-1111-1111-1111-111111111111"
  id           = "22222222-2222-2222-2222-222222222222"
}
