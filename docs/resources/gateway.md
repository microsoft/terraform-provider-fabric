---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "fabric_gateway Resource - terraform-provider-fabric"
subcategory: ""
description: |-
  This resource manages a Fabric Gateway.
  See Gateway https://learn.microsoft.com/power-bi/guidance/powerbi-implementation-planning-data-gateways for more information.
  -> This item supports Service Principal authentication.
  ~> This resource is in preview. To access it, you must explicitly enable the preview mode in the provider level configuration.
---

# fabric_gateway (Resource)

This resource manages a Fabric Gateway.

See [Gateway](https://learn.microsoft.com/power-bi/guidance/powerbi-implementation-planning-data-gateways) for more information.

-> This item supports Service Principal authentication.

~> This resource is in **preview**. To access it, you must explicitly enable the `preview` mode in the provider level configuration.

## Example Usage

```terraform
# Virtual Network Gateway
resource "fabric_gateway" "example" {
  type                            = "VirtualNetwork"
  display_name                    = "example"
  inactivity_minutes_before_sleep = 30
  number_of_member_gateways       = 1
  virtual_network_azure_resource = {
    resource_group_name  = "example resource group"
    virtual_network_name = "example virtual network"
    subnet_name          = "example subnet"
    subscription_id      = "00000000-0000-0000-0000-000000000000"
  }
  capacity_id = "11111111-1111-1111-1111-111111111111"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `type` (String) The Gateway type. Accepted values: `VirtualNetwork`

### Optional

- `capacity_id` (String) The Gateway capacity ID.
- `display_name` (String) The Gateway display name.
- `inactivity_minutes_before_sleep` (Number) The Gateway inactivity minutes before sleep.
- `number_of_member_gateways` (Number) The Gateway number of member gateways.
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))
- `virtual_network_azure_resource` (Attributes) The Gateway virtual network Azure resource. (see [below for nested schema](#nestedatt--virtual_network_azure_resource))

### Read-Only

- `id` (String) The Gateway ID.

<a id="nestedatt--timeouts"></a>

### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).

<a id="nestedatt--virtual_network_azure_resource"></a>

### Nested Schema for `virtual_network_azure_resource`

Required:

- `resource_group_name` (String) The resource group name.
- `subnet_name` (String) The subnet name.
- `subscription_id` (String) The subscription ID.
- `virtual_network_name` (String) The virtual network name.

## Import

Import is supported using the following syntax:

```shell
# terraform import fabric_gateway.example "<GatewayID>"
terraform import fabric_gateway.example "00000000-0000-0000-0000-000000000000"
```
