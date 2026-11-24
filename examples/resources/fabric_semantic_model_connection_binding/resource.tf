resource "fabric_semantic_model_connection_binding" "example" {
  workspace_id      = "00000000-0000-0000-0000-000000000000"
  semantic_model_id = "00000000-0000-0000-0000-000000000000"

  connection_id     = "00000000-0000-0000-0000-000000000000"
  connectivity_type = "ShareableCloud"

  connection_details = {
    path = "server;db"
    type = "SQL"
  }
}
