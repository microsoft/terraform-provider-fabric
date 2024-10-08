---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "fabric_lakehouse_tables Data Source - terraform-provider-fabric"
subcategory: ""
description: |-
  List a Fabric Lakehouse Tables.
  Use this data source to list Lakehouse Tables https://learn.microsoft.com/fabric/data-engineering/lakehouse-and-delta-tables.
  -> This item supports Service Principal authentication.
---

# fabric_lakehouse_tables (Data Source)

List a Fabric Lakehouse Tables.

Use this data source to list [Lakehouse Tables](https://learn.microsoft.com/fabric/data-engineering/lakehouse-and-delta-tables).

-> This item supports Service Principal authentication.

## Example Usage

```terraform
data "fabric_lakehouse_tables" "example" {
  lakehouse_id = "11111111-1111-1111-1111-111111111111"
  workspace_id = "00000000-0000-0000-0000-000000000000"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `lakehouse_id` (String) The Lakehouse ID.
- `workspace_id` (String) The Workspace ID.

### Optional

- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

### Read-Only

- `values` (Attributes List) The list of Lakehouse Tables. (see [below for nested schema](#nestedatt--values))

<a id="nestedatt--timeouts"></a>

### Nested Schema for `timeouts`

Optional:

- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).

<a id="nestedatt--values"></a>

### Nested Schema for `values`

Read-Only:

- `format` (String) The Format of the table.
- `location` (String) The Location of the table.
- `name` (String) The Name of the table.
- `type` (String) The Type of the table. Possible values: `External`, `Managed`.
