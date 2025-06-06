# Get data-source details
data "fabric_eventstream_source_connection" "example" {
  workspace_id   = "00000000-0000-0000-0000-000000000000"
  eventstream_id = "11111111-1111-1111-1111-111111111111"
  source_id      = "22222222-2222-2222-2222-222222222222"
}
