resource "fabric_spark_environment_settings" "example" {
  workspace_id       = "00000000-0000-0000-0000-000000000000"
  environment_id     = "11111111-1111-1111-1111-111111111111"
  publication_status = "Published"

  driver_cores  = 4
  driver_memory = "28g"

  executor_cores  = 4
  executor_memory = "28g"

  runtime_version = "1.2"

  dynamic_executor_allocation = {
    /*
			your settings here
			*/
  }

  pool = {
    /*
			your settings here
			*/
  }

  spark_properties = {
    /*
			your settings here
			*/
  }
}
