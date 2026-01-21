resource "fabric_item_job_scheduler" "cron_configuration_example" {
  workspace_id = "00000000-0000-0000-0000-000000000000"
  item_id      = "11111111-1111-1111-1111-111111111111"
  job_type     = "Execute"
  enabled      = true #or false
  configuration = {
    start_date_time = "2025-11-11T10:00:00Z"
    end_date_time   = "2025-11-12T10:00:00Z"
    type            = "Cron"
    interval        = 10
  }
}

resource "fabric_item_job_scheduler" "daily_configuration_example" {
  workspace_id = "00000000-0000-0000-0000-000000000000"
  item_id      = "11111111-1111-1111-1111-111111111111"
  job_type     = "Execute"
  enabled      = true #or false
  configuration = {
    start_date_time = "2025-11-11T10:00:00Z"
    end_date_time   = "2025-11-12T10:00:00Z"
    type            = "Daily"
    times           = ["10:00"]
  }
}

resource "fabric_item_job_scheduler" "weekly_configuration_example" {
  workspace_id = "00000000-0000-0000-0000-000000000000"
  item_id      = "11111111-1111-1111-1111-111111111111"
  job_type     = "Execute"
  enabled      = true #or false
  configuration = {
    start_date_time = "2025-11-11T10:00:00Z"
    end_date_time   = "2025-11-12T10:00:00Z"
    type            = "Weekly"
    times           = ["10:00"]
    weekdays        = ["Monday"]
  }
}

resource "fabric_item_job_scheduler" "monthly_configuration_day_of_month_example" {
  workspace_id = "00000000-0000-0000-0000-000000000000"
  item_id      = "11111111-1111-1111-1111-111111111111"
  job_type     = "Execute"
  enabled      = true #or false
  configuration = {
    start_date_time = "2025-11-11T10:00:00Z"
    end_date_time   = "2025-11-12T10:00:00Z"
    type            = "Monthly"
    times           = ["10:00"]
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
  job_type     = "Execute"
  enabled      = true #or false
  configuration = {
    start_date_time = "2025-11-11T10:00:00Z"
    end_date_time   = "2025-11-12T10:00:00Z"
    type            = "Monthly"
    times           = ["10:00"]
    recurrence      = 1
    occurrence = {
      occurrence_type = "OrdinalWeekday"
      week_index      = "First"
      weekday         = "Monday"
    }
  }
}
