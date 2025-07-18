---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "fabric_gateway_role_assignment Resource - terraform-provider-fabric"
subcategory: ""
description: |-
  The Gateway Role Assignment resource allows you to manage a Fabric Gateway Role Assignment https://learn.microsoft.com/power-bi/guidance/powerbi-implementation-planning-data-gateways.
  -> This resource supports Service Principal authentication.
---

# fabric_gateway_role_assignment (Resource)

The Gateway Role Assignment resource allows you to manage a Fabric [Gateway Role Assignment](https://learn.microsoft.com/power-bi/guidance/powerbi-implementation-planning-data-gateways).

-> This resource supports Service Principal authentication.

## Example Usage

```terraform
resource "fabric_gateway_role_assignment" "example" {
  gateway_id = "00000000-0000-0000-0000-000000000000"
  principal = {
    id   = "11111111-1111-1111-1111-111111111111"
    type = "User"
  }
  role = "ConnectionCreatorWithResharing"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `gateway_id` (String) <i style="color:red;font-weight: bold">(ForceNew)</i> The Gateway ID.
- `principal` (Attributes) The principal. (see [below for nested schema](#nestedatt--principal))
- `role` (String) The gateway role of the principal. Value must be one of : `Admin`, `ConnectionCreator`, `ConnectionCreatorWithResharing`.

### Optional

- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

### Read-Only

- `id` (String) The Gateway Role Assignment ID.

<a id="nestedatt--principal"></a>

### Nested Schema for `principal`

Required:

- `id` (String) <i style="color:red;font-weight: bold">(ForceNew)</i> The principal ID.
- `type` (String) <i style="color:red;font-weight: bold">(ForceNew)</i> The type of the principal. Value must be one of : `Group`, `ServicePrincipal`, `ServicePrincipalProfile`, `User`.

<a id="nestedatt--timeouts"></a>

### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).

## Import

Import is supported using the following syntax:

The [`terraform import` command](https://developer.hashicorp.com/terraform/cli/commands/import) can be used, for example:

```shell
# terraform import fabric_gateway_role_assignment.example "<GatewayID>/<GatewayRoleAssignmentID>"
terraform import fabric_gateway_role_assignment.example "00000000-0000-0000-0000-000000000000/11111111-1111-1111-1111-111111111111"
```
