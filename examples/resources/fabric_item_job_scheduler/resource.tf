resource "fabric_item_job_scheduler" "cron_configuration_example" {
  workspace_id = "00000000-0000-0000-0000-000000000000"
  item_id      = "11111111-1111-1111-1111-111111111111"
  job_type     = "MyJobType"
  enabled      = true #or false
  configuration = {
    start_date_time = "YYYY-MM-DDTHH:mm:ssZ"
    end_date_time   = "YYYY-MM-DDTHH:mm:ssZ"
    type            = "Cron"
    interval        = 10
  }
}

resource "fabric_item_job_scheduler" "daily_configuration_example" {
  workspace_id = "00000000-0000-0000-0000-000000000000"
  item_id      = "11111111-1111-1111-1111-111111111111"
  job_type     = "MyJobType"
  enabled      = true #or false
  configuration = {
    start_date_time = "YYYY-MM-DDTHH:mm:ssZ"
    end_date_time   = "YYYY-MM-DDTHH:mm:ssZ"
    type            = "Daily"
    times           = ["HH:mm:ss"]
  }
}

resource "fabric_item_job_scheduler" "weekly_configuration_example" {
  workspace_id = "00000000-0000-0000-0000-000000000000"
  item_id      = "11111111-1111-1111-1111-111111111111"
  job_type     = "MyJobType"
  enabled      = true #or false
  configuration = {
    start_date_time = "YYYY-MM-DDTHH:mm:ssZ"
    end_date_time   = "YYYY-MM-DDTHH:mm:ssZ"
    type            = "Weekly"
    times           = ["HH:mm:ss"]
    weekdays        = ["Monday"]
  }
}

resource "fabric_item_job_scheduler" "monthly_configuration_day_of_month_example" {
  workspace_id = "00000000-0000-0000-0000-000000000000"
  item_id      = "11111111-1111-1111-1111-111111111111"
  job_type     = "MyJobType"
  enabled      = true #or false
  configuration = {
    start_date_time = "YYYY-MM-DDTHH:mm:ssZ"
    end_date_time   = "YYYY-MM-DDTHH:mm:ssZ"
    type            = "Monthly"
    times           = ["HH:mm:ss"]
    weekdays        = ["Monday"]
    recurrence      = 1
    occurrence = {
      occurrence_type = "DayOfMonth"
      day_of_month    = 10
    }
  }
}

resource "fabric_item_job_scheduler" "monthly_configuration_ordinal_weekday_example" {
  workspace_id = "00000000-0000-0000-0000-000000000000"
  item_id      = "11111111-1111-1111-1111-111111111111"
  job_type     = "MyJobType"
  enabled      = true #or false
  configuration = {
    start_date_time = "YYYY-MM-DDTHH:mm:ssZ"
    end_date_time   = "YYYY-MM-DDTHH:mm:ssZ"
    type            = "Monthly"
    times           = ["HH:mm:ss"]
    weekdays        = ["Monday"]
    recurrence      = 1
    occurrence = {
      occurrence_type = "DayOfMonth"
      week_index      = "First"
      weekday         = "Monday"
    }
  }
}
