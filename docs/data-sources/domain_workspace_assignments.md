---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "fabric_domain_workspace_assignments Data Source - terraform-provider-fabric"
subcategory: ""
description: |-
  List a Fabric Domain Workspace Assignments.
  Use this data source to list Domain Workspace Assignments https://learn.microsoft.com/fabric/governance/domains.
  -> This item supports Service Principal authentication.
  ~> This data-source is in preview. To access it, you must explicitly enable the preview mode in the provider level configuration.
---

# fabric_domain_workspace_assignments (Data Source)

List a Fabric Domain Workspace Assignments.

Use this data source to list [Domain Workspace Assignments](https://learn.microsoft.com/fabric/governance/domains).

-> This item supports Service Principal authentication.

~> This data-source is in **preview**. To access it, you must explicitly enable the `preview` mode in the provider level configuration.

## Example Usage

```terraform
data "fabric_domain_workspace_assignments" "example" {
  domain_id = "00000000-0000-0000-0000-000000000000"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `domain_id` (String) The Domain ID.

### Optional

- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

### Read-Only

- `values` (Attributes List) The list of Domain Workspace Assignments. (see [below for nested schema](#nestedatt--values))

<a id="nestedatt--timeouts"></a>

### Nested Schema for `timeouts`

Optional:

- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).

<a id="nestedatt--values"></a>

### Nested Schema for `values`

Read-Only:

- `display_name` (String) The Workspace display name.
- `id` (String) The Workspace ID.
