resource "fabric_folder" "example_workspace_root_folder" {
  display_name = "example"
  workspace_id = "00000000-0000-0000-0000-000000000000"
}

resource "fabric_folder" "example_subfolder" {
  display_name     = "example"
  workspace_id     = "00000000-0000-0000-0000-000000000000"
  parent_folder_id = "11111111-1111-1111-1111-111111111111"
}
#changing the parent_folder_id will move the folder to a different parent folder
# resource "fabric_folder" "example_subfolder" {
#   display_name     = "example"
#   workspace_id     = "00000000-0000-0000-0000-000000000000"
#   parent_folder_id = "00000000-0000-0000-0000-000000000000"
# }
