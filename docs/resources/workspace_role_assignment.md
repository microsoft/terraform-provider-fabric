---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "fabric_workspace_role_assignment Resource - terraform-provider-fabric"
subcategory: ""
description: |-
  Manage a Workspace Role Assignment.
  See Roles in Workspaces https://learn.microsoft.com/fabric/get-started/roles-workspaces for more information.
  -> This item supports Service Principal authentication.
---

# fabric_workspace_role_assignment (Resource)

Manage a Workspace Role Assignment.

See [Roles in Workspaces](https://learn.microsoft.com/fabric/get-started/roles-workspaces) for more information.

-> This item supports Service Principal authentication.

## Example Usage

```terraform
resource "fabric_workspace_role_assignment" "example" {
  workspace_id   = "00000000-0000-0000-0000-000000000000"
  principal_id   = "11111111-1111-1111-1111-111111111111"
  principal_type = "User"
  role           = "Member"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `principal_id` (String) The Principal ID.
- `principal_type` (String) The type of the principal. Accepted values: `Group`, `ServicePrincipal`, `ServicePrincipalProfile`, `User`.
- `role` (String) The Workspace Role of the principal. Accepted values: `Admin`, `Contributor`, `Member`, `Viewer`.
- `workspace_id` (String) The Workspace ID.

### Optional

- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

### Read-Only

- `id` (String) The Workspace Role Assignment ID.

<a id="nestedatt--timeouts"></a>

### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).

## Import

Import is supported using the following syntax:

```shell
# terraform import fabric_workspace_role_assignment.example "<WorkspaceID>/<WorkspaceRoleAssignmentID>"
terraform import fabric_workspace_role_assignment.example "00000000-0000-0000-0000-000000000000/11111111-1111-1111-1111-111111111111"
```