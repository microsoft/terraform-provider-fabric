---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "fabric_kql_dashboards Data Source - terraform-provider-fabric"
subcategory: ""
description: |-
  The KQL Dashboards data-source allows you to retrieve a list of Fabric KQL Dashboards https://learn.microsoft.com/fabric/real-time-intelligence/dashboard-real-time-create.
  -> This data-source supports Service Principal authentication.
---

# fabric_kql_dashboards (Data Source)

The KQL Dashboards data-source allows you to retrieve a list of Fabric [KQL Dashboards](https://learn.microsoft.com/fabric/real-time-intelligence/dashboard-real-time-create).

-> This data-source supports Service Principal authentication.

## Example Usage

```terraform
data "fabric_kql_dashboards" "example" {
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

- `values` (Attributes Set) The set of KQL Dashboards. (see [below for nested schema](#nestedatt--values))

<a id="nestedatt--timeouts"></a>

### Nested Schema for `timeouts`

Optional:

- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).

<a id="nestedatt--values"></a>

### Nested Schema for `values`

Read-Only:

- `description` (String) The KQL Dashboard description.
- `display_name` (String) The KQL Dashboard display name.
- `id` (String) The KQL Dashboard ID.
- `workspace_id` (String) The Workspace ID.
