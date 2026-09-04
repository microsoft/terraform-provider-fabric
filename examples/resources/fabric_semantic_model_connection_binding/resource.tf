resource "fabric_semantic_model_connection_binding" "example" {
  workspace_id      = "00000000-0000-0000-0000-000000000000"
  semantic_model_id = "11111111-1111-1111-1111-111111111111"
  connectivity_type = "ShareableCloud"
  connection_id     = "22222222-2222-2222-2222-222222222222"
  connection_details = {
    path = "https://contoso.database.windows.net;sales"
    type = "Sql"
  }
}
