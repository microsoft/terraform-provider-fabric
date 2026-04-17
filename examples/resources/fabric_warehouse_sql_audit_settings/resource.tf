resource "fabric_warehouse_sql_audit_settings" "example" {
  workspace_id = "00000000-0000-0000-0000-000000000000"
  item_id      = "11111111-1111-1111-1111-111111111111"
  state          = "Enabled"
  retention_days = 10
  audit_actions_and_groups = [
    "SUCCESSFUL_DATABASE_AUTHENTICATION_GROUP",
    "FAILED_DATABASE_AUTHENTICATION_GROUP",
    "BATCH_COMPLETED_GROUP",
  ]
}
