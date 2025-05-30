---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "fabric_domain_workspace_assignments Resource - terraform-provider-fabric"
subcategory: ""
description: |-
  The Domain Workspace Assignments resource allows you to manage a Fabric Domain Workspace Assignments https://learn.microsoft.com/fabric/governance/domains.
  -> This resource supports Service Principal authentication.
  ~> This resource is in preview. To access it, you must explicitly enable the preview mode in the provider level configuration.
---

# fabric_domain_workspace_assignments (Resource)

The Domain Workspace Assignments resource allows you to manage a Fabric [Domain Workspace Assignments](https://learn.microsoft.com/fabric/governance/domains).

-> This resource supports Service Principal authentication.

~> This resource is in **preview**. To access it, you must explicitly enable the `preview` mode in the provider level configuration.

## Example Usage

```terraform
resource "fabric_workspace" "example" {
  display_name = "example"
  description  = "Example Workspace"
}

resource "fabric_domain" "example" {
  display_name = "example"
  description  = "Example Domain"
}

resource "fabric_domain_workspace_assignments" "example" {
  domain_id = fabric_domain.example.id
  workspace_ids = [
    fabric_workspace.example.id
  ]
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `domain_id` (String) <i style="color:red;font-weight: bold">(ForceNew)</i> The Domain ID.
- `workspace_ids` (Set of String) The set of Workspace IDs. Set must contain at least 1 elements.

### Optional

- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

<a id="nestedatt--timeouts"></a>

### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
