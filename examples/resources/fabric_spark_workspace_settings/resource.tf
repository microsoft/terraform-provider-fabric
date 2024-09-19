# Example 1
resource "fabric_spark_workspace_settings" "example" {
  workspace_id = "00000000-0000-0000-0000-000000000000"

  automatic_log = {
    /*
			your settings here
			*/
  }

  high_concurrency = {
    /*
			your settings here
			*/
  }

  environment = {
    /*
			your settings here
			*/
  }

  pool = {
    /*
			your settings here
			*/
  }
}

# Example 2
resource "fabric_spark_workspace_settings" "example2" {
  workspace_id = "00000000-0000-0000-0000-000000000000"

  automatic_log = {
    enabled = true
  }

  high_concurrency = {
    notebook_interactive_run_enabled = false
  }

  environment = {
    name            = "MyExampleEnvironment"
    runtime_version = "1.3"
  }

  pool = {
    default_pool = {
      name = "MyExampleCustomPool"
      type = "Workspace"
    }
    starter_pool = {
      max_executors  = 3
      max_node_count = 1
    }
    customize_compute_enabled = true
  }
}
