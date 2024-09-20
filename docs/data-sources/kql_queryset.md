---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "fabric_kql_queryset Data Source - terraform-provider-fabric"
subcategory: ""
description: |-
  Get a Fabric KQL Queryset.
  Use this data source to fetch a KQL Queryset https://learn.microsoft.com/fabric/real-time-intelligence/kusto-query-set.
  -> This item supports Service Principal authentication.
---

# fabric_kql_queryset (Data Source)

Get a Fabric KQL Queryset.

Use this data source to fetch a [KQL Queryset](https://learn.microsoft.com/fabric/real-time-intelligence/kusto-query-set).

-> This item supports Service Principal authentication.

## Example Usage

```terraform
data "fabric_kql_queryset" "example_by_id" {
  id           = "11111111-1111-1111-1111-111111111111"
  workspace_id = "00000000-0000-0000-0000-000000000000"
}

data "fabric_kql_queryset" "example_by_name" {
  display_name = "example"
  workspace_id = "00000000-0000-0000-0000-000000000000"
}

# This is an invalid data source
# Do not specify `id` and `display_name` in the same data source block
# data "fabric_kql_queryset" "example" {
#   display_name = "example"
#   id           = "11111111-1111-1111-1111-111111111111"
#   workspace_id = "00000000-0000-0000-0000-000000000000"
# }
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `workspace_id` (String) The Workspace ID.

### Optional

- `display_name` (String) The KQL Queryset display name.
- `id` (String) The KQL Queryset ID.
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

### Read-Only

- `description` (String) The KQL Queryset description.

<a id="nestedatt--timeouts"></a>

### Nested Schema for `timeouts`

Optional:

- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).