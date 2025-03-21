---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "fabric_workspace_role_assignments Data Source - terraform-provider-fabric"
subcategory: ""
description: |-
  List Fabric Workspace Role Assignments.
  Use this data source to list Workspace Role Assignments https://learn.microsoft.com/fabric/fundamentals/roles-workspaces.
  -> This item supports Service Principal authentication.
---

# fabric_workspace_role_assignments (Data Source)

List Fabric Workspace Role Assignments.

Use this data source to list [Workspace Role Assignments](https://learn.microsoft.com/fabric/fundamentals/roles-workspaces).

-> This item supports Service Principal authentication.

## Example Usage

```terraform
data "fabric_workspace_role_assignments" "example" {
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

- `values` (Attributes List) The list of Workspace Role Assignments. (see [below for nested schema](#nestedatt--values))

<a id="nestedatt--timeouts"></a>

### Nested Schema for `timeouts`

Optional:

- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).

<a id="nestedatt--values"></a>

### Nested Schema for `values`

Read-Only:

- `id` (String) The Workspace Role Assignment ID.
- `principal_details` (Attributes) The principal details. (see [below for nested schema](#nestedatt--values--principal_details))
- `principal_display_name` (String) The principal's display name.
- `principal_id` (String) The Principal ID.
- `principal_type` (String) The type of the principal. Possible values: `Group`, `ServicePrincipal`, `ServicePrincipalProfile`, `User`.
- `role` (String) The workspace role of the principal. Possible values: `Admin`, `Contributor`, `Member`, `Viewer`.

<a id="nestedatt--values--principal_details"></a>

### Nested Schema for `values.principal_details`

Read-Only:

- `app_id` (String) The service principal's Microsoft Entra App ID.
- `group_type` (String) The type of the group. Possible values: `DistributionList`, `SecurityGroup`, `Unknown`.
- `parent_principal_id` (String) The parent principal ID of Service Principal Profile.
- `user_principal_name` (String) The user principal name.
