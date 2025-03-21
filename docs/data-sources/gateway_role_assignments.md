---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "fabric_gateway_role_assignments Data Source - terraform-provider-fabric"
subcategory: ""
description: |-
  List Fabric Gateway Role Assignments.
  Use this data source to list [Gateway Role Assignments].
  -> This item supports Service Principal authentication.
  ~> This data-source is in preview. To access it, you must explicitly enable the preview mode in the provider level configuration.
---

# fabric_gateway_role_assignments (Data Source)

List Fabric Gateway Role Assignments.

Use this data source to list [Gateway Role Assignments].

-> This item supports Service Principal authentication.

~> This data-source is in **preview**. To access it, you must explicitly enable the `preview` mode in the provider level configuration.

## Example Usage

```terraform
data "fabric_gateway_role_assignments" "example" {
  gateway_id = "00000000-0000-0000-0000-000000000000"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `gateway_id` (String) The Gateway ID.

### Optional

- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

### Read-Only

- `values` (Attributes List) The list of Gateway Role Assignments. (see [below for nested schema](#nestedatt--values))

<a id="nestedatt--timeouts"></a>

### Nested Schema for `timeouts`

Optional:

- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).

<a id="nestedatt--values"></a>

### Nested Schema for `values`

Read-Only:

- `id` (String) The Gateway Role Assignment ID.
- `principal` (Attributes) The principal. (see [below for nested schema](#nestedatt--values--principal))
- `role` (String) The gateway role of the principal. Possible values: `Admin`, `ConnectionCreator`, `ConnectionCreatorWithResharing`.

<a id="nestedatt--values--principal"></a>

### Nested Schema for `values.principal`

Read-Only:

- `id` (String) The principal ID.
- `type` (String) The principal type. Possible values: `Group`, `ServicePrincipal`, `ServicePrincipalProfile`, `User`.
