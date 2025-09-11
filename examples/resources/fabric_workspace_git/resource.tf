# Example of Azure DevOps integration with automatic credentials (no SPN support, only User identity is supported)
resource "fabric_workspace_git" "azdo_automatic" {
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
  git_credentials = {
    source = "Automatic"
  }
}

# Example of Azure DevOps integration with configured credentials
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
  git_credentials = {
    source        = "ConfiguredConnection"
    connection_id = "11111111-1111-1111-1111-111111111111"
  }
}

# Example of GitHub integration with configured credentials
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
    source        = "ConfiguredConnection"
    connection_id = "11111111-1111-1111-1111-111111111111"
  }
}
