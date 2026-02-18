resource "fabric_external_data_share" "example" {
  workspace_id = "00000000-0000-0000-0000-000000000000"
  item_id      = "11111111-1111-1111-1111-111111111111"
  paths        = ["Files/Sales/Contoso_Sales_2023"]
  recipient = {
    "user_principal_name" = "example@example.com"
  }
}

resource "fabric_external_data_share" "example_with_tenant" {
  workspace_id = "00000000-0000-0000-0000-000000000000"
  item_id      = "11111111-1111-1111-1111-111111111111"
  paths        = ["Files/Sales/Contoso_Sales_2023"]
  recipient = {
    "user_principal_name" = "example@example.com"
    "tenant_id"           = "22222222-2222-2222-2222-222222222222"
  }
}
