resource "fabric_spark_custom_pool" "example" {
  workspace_id = "00000000-0000-0000-0000-000000000000"
  name         = "example"
  node_family  = "MemoryOptimized"
  node_size    = "Small"
  type         = "Workspace"

  auto_scale = {
    enabled        = true
    min_node_count = 1
    max_node_count = 3
  }

  dynamic_executor_allocation = {
    enabled       = true
    min_executors = 1
    max_executors = 2
  }
}
