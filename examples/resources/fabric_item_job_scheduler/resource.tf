resource "fabric_item_job_scheduler" "cron_configuration_example" {
  workspace_id = "00000000-0000-0000-0000-000000000000"
  item_id      = "11111111-1111-1111-1111-111111111111"
  job_type     = "MyJobType"
  enabled      = true #or false
  configuration = {
    start_date_time    = "YYYY-MM-DDTHH:mm:ssZ"
    end_date_time      = "YYYY-MM-DDTHH:mm:ssZ"
    local_time_zone_id = "YourTimezoneId"
    type               = "Cron"
    interval           = 10
  }
}

resource "fabric_item_job_scheduler" "daily_configuration_example" {
  workspace_id = "00000000-0000-0000-0000-000000000000"
  item_id      = "11111111-1111-1111-1111-111111111111"
  job_type     = "MyJobType"
  enabled      = true #or false
  configuration = {
    start_date_time    = "YYYY-MM-DDTHH:mm:ssZ"
    end_date_time      = "YYYY-MM-DDTHH:mm:ssZ"
    local_time_zone_id = "YourTimezoneId"
    type               = "Daily"
    times              = ["HH:mm:ss"]
  }
}

resource "fabric_item_job_scheduler" "weekly_configuration_example" {
  workspace_id = "00000000-0000-0000-0000-000000000000"
  item_id      = "11111111-1111-1111-1111-111111111111"
  job_type     = "MyJobType"
  enabled      = true #or false
  configuration = {
    start_date_time    = "YYYY-MM-DDTHH:mm:ssZ"
    end_date_time      = "YYYY-MM-DDTHH:mm:ssZ"
    local_time_zone_id = "YourTimezoneId"
    type               = "Weekly"
    times              = ["HH:mm:ss"]
    weekdays           = ["Monday"]
  }
}
