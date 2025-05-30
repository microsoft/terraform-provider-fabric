---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "fabric_workspace_git Data Source - terraform-provider-fabric"
subcategory: ""
description: |-
  The Workspace Git data-source allows you to retrieve details about a Fabric Workspace Git https://learn.microsoft.com/fabric/cicd/git-integration/intro-to-git-integration.
  -> This data-source does not support Service Principal. Please use a User context authentication.
  ~> This data-source is in preview. To access it, you must explicitly enable the preview mode in the provider level configuration.
---

# fabric_workspace_git (Data Source)

The Workspace Git data-source allows you to retrieve details about a Fabric [Workspace Git](https://learn.microsoft.com/fabric/cicd/git-integration/intro-to-git-integration).

-> This data-source does not support Service Principal. Please use a User context authentication.

~> This data-source is in **preview**. To access it, you must explicitly enable the `preview` mode in the provider level configuration.

## Example Usage

```terraform
data "fabric_workspace_git" "example" {
  workspace_id = "00000000-0000-0000-0000-000000000000"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `workspace_id` (String) The Workspace ID.

### Optional

- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

### Read-Only

- `git_connection_state` (String) The Git connection state. Value must be one of : `Connected`, `ConnectedAndInitialized`, `NotConnected`.
- `git_credentials` (Attributes) The Git credentials details. (see [below for nested schema](#nestedatt--git_credentials))
- `git_provider_details` (Attributes) The Git provider details. (see [below for nested schema](#nestedatt--git_provider_details))
- `git_sync_details` (Attributes) The Git sync details. (see [below for nested schema](#nestedatt--git_sync_details))
- `id` (String) The Workspace Git ID.

<a id="nestedatt--timeouts"></a>

### Nested Schema for `timeouts`

Optional:

- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).

<a id="nestedatt--git_credentials"></a>

### Nested Schema for `git_credentials`

Read-Only:

- `connection_id` (String) The connection ID.
- `source` (String) The Git credentials source. Value must be one of : `Automatic`, `ConfiguredConnection`, `None`.

<a id="nestedatt--git_provider_details"></a>

### Nested Schema for `git_provider_details`

Read-Only:

- `branch_name` (String) The branch name.
- `directory_name` (String) The directory name.
- `git_provider_type` (String) The git provider type. Value must be one of : `AzureDevOps`, `GitHub`.
- `organization_name` (String) The Azure DevOps organization name.
- `owner_name` (String) The GitHub owner name.
- `project_name` (String) The Azure DevOps project name.
- `repository_name` (String) The repository name.

<a id="nestedatt--git_sync_details"></a>

### Nested Schema for `git_sync_details`

Read-Only:

- `head` (String) The git head.
- `last_sync_time` (String) The last sync time.
