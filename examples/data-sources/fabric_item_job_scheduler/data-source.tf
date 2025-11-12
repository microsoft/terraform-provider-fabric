data "fabric_item_job_scheduler" "example" {
  workspace_id = "00000000-0000-0000-0000-000000000000"
  item_id      = "11111111-1111-1111-1111-111111111111"
  scheduleId   = "22222222-2222-2222-2222-222222222222"
  jobType      = "Execute"
}
