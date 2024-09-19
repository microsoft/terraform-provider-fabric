resource "fabric_workspace_git" "example" {
  workspace_id            = "00000000-0000-0000-0000-000000000000"
  initialization_strategy = "PreferWorkspace"
  git_provider_details = {
    git_provider_type = "AzureDevOps"
    organization_name = "MyExampleOrg"
    project_name      = "MyExampleProject"
    repository_name   = "ExampleRepo"
    branch_name       = "ExampleBranch"
    directory_name    = "/ExampleDirectory"
  }
}
