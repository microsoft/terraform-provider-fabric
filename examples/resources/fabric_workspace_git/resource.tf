# Example of Azure DevOps integration
resource "fabric_workspace_git" "azdo" {
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

# Example of GitHub integration
resource "fabric_workspace_git" "github" {
  workspace_id            = "00000000-0000-0000-0000-000000000000"
  initialization_strategy = "PreferWorkspace"
  git_provider_details = {
    git_provider_type = "GitHub"
    owner_name        = "ExampleOwner"
    repository_name   = "ExampleRepo"
    branch_name       = "ExampleBranch"
    directory_name    = "/ExampleDirectory"
  }
  git_credentials = {
    connection_id = "11111111-1111-1111-1111-111111111111"
  }
}
