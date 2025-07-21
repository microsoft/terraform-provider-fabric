data "fabric_folders" "example_recursively" {
  workspace_id = "00000000-0000-0000-0000-000000000000"
  #recursive     = true this is the default value if not provided
}

data "fabric_folders" "example_non_recursively" {
  workspace_id = "00000000-0000-0000-0000-000000000000"
  recursive    = false
}

data "fabric_folders" "example_with_root_folder_recursively" {
  workspace_id   = "00000000-0000-0000-0000-000000000000"
  root_folder_id = "11111111-1111-1111-1111-111111111111"
  #recursive     = true this is the default value if not provided
}

data "fabric_folders" "example_with_root_folder_non_recursively" {
  workspace_id   = "00000000-0000-0000-0000-000000000000"
  root_folder_id = "11111111-1111-1111-1111-111111111111"
  recursive      = false
}
