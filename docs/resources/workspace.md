---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "fabric_workspace Resource - terraform-provider-fabric"
subcategory: ""
description: |-
  The Workspace resource allows you to manage a Fabric Workspace https://learn.microsoft.com/fabric/get-started/workspaces.
  -> This resource supports Service Principal authentication.
---

# fabric_workspace (Resource)

The Workspace resource allows you to manage a Fabric [Workspace](https://learn.microsoft.com/fabric/get-started/workspaces).

-> This resource supports Service Principal authentication.

## Example Usage

```terraform
# Simple Workspace
resource "fabric_workspace" "example1" {
  display_name = "example1"
  description  = "Example Workspace 1"
}

# Workspace with Capacity and Identity
data "fabric_capacity" "example" {
  display_name = "example"
}

resource "fabric_workspace" "example2" {
  display_name = "example2"
  description  = "Example Workspace 2"
  capacity_id  = data.fabric_capacity.example.id
  identity = {
    type = "SystemAssigned"
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `display_name` (String) The Workspace display name. String length must be at most 256.Value must not be one of : .

### Optional

- `capacity_id` (String) The ID of the Fabric Capacity to assign to the Workspace.
- `description` (String) The Workspace description. Value defaults to ``. String length must be at most 4000.
- `identity` (Attributes) A workspace identity (see [Workspace Identity](https://learn.microsoft.com/fabric/security/workspace-identity) for more information). (see [below for nested schema](#nestedatt--identity))
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

### Read-Only

- `capacity_assignment_progress` (String) A Workspace assignment to capacity progress status. Value must be one of : `Completed`, `Failed`, `InProgress`.
- `capacity_region` (String) The region of the capacity associated with this workspace. Value must be one of : .
- `id` (String) The Workspace ID.
- `onelake_endpoints` (Attributes) The OneLake API endpoints associated with this workspace. (see [below for nested schema](#nestedatt--onelake_endpoints))
- `type` (String) The Workspace type.

<a id="nestedatt--identity"></a>

### Nested Schema for `identity`

Required:

- `type` (String) The identity type. Value must be one of : `SystemAssigned`.

Read-Only:

- `application_id` (String) The application ID.
- `service_principal_id` (String) The service principal ID.

<a id="nestedatt--timeouts"></a>

### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).

<a id="nestedatt--onelake_endpoints"></a>

### Nested Schema for `onelake_endpoints`

Read-Only:

- `blob_endpoint` (String) The OneLake API endpoint available for Blob API operations.
- `dfs_endpoint` (String) The OneLake API endpoint available for Distributed File System (DFS) or ADLSgen2 filesystem API operations.

## Import

Import is supported using the following syntax:

The [`terraform import` command](https://developer.hashicorp.com/terraform/cli/commands/import) can be used, for example:

```shell
# terraform import fabric_workspace.example "<WorkspaceID>"
terraform import fabric_workspace.example "00000000-0000-0000-0000-000000000000"
```
