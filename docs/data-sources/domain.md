---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "fabric_domain Data Source - terraform-provider-fabric"
subcategory: ""
description: |-
  Get a Fabric Domain.
  Use this data source to get Domain https://learn.microsoft.com/fabric/governance/domains.
  -> This item supports Service Principal authentication.
  ~> This data-source is in preview. To access it, you must explicitly enable the preview mode in the provider level configuration.
---

# fabric_domain (Data Source)

Get a Fabric Domain.

Use this data source to get [Domain](https://learn.microsoft.com/fabric/governance/domains).

-> This item supports Service Principal authentication.

~> This data-source is in **preview**. To access it, you must explicitly enable the `preview` mode in the provider level configuration.

## Example Usage

```terraform
data "fabric_domain" "example" {
  id = "00000000-0000-0000-0000-000000000000"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `id` (String) The Domain ID.

### Optional

- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

### Read-Only

- `contributors_scope` (String) The Domain contributors scope. Possible values: `AdminsOnly`, `AllTenant`, `SpecificUsersAndGroups`.
- `description` (String) The Domain description.
- `display_name` (String) The Domain display name.
- `parent_domain_id` (String) The Domain parent ID.

<a id="nestedatt--timeouts"></a>

### Nested Schema for `timeouts`

Optional:

- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
