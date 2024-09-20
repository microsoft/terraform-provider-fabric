---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "fabric_sql_endpoints Data Source - terraform-provider-fabric"
subcategory: ""
description: |-
  List a Fabric SQL Endpoints.
  Use this data source to list SQL Endpoints https://learn.microsoft.com/power-bi/transform-model/sqlendpoints/sqlendpoints-overview.
  -> This item does not support Service Principal. Please use a User context authentication.
---

# fabric_sql_endpoints (Data Source)

List a Fabric SQL Endpoints.

Use this data source to list [SQL Endpoints](https://learn.microsoft.com/power-bi/transform-model/sqlendpoints/sqlendpoints-overview).

-> This item does not support Service Principal. Please use a User context authentication.

## Example Usage

```terraform
data "fabric_sql_endpoints" "example" {
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

- `values` (Attributes List) The list of SQL Endpoints. (see [below for nested schema](#nestedatt--values))

<a id="nestedatt--timeouts"></a>

### Nested Schema for `timeouts`

Optional:

- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).

<a id="nestedatt--values"></a>

### Nested Schema for `values`

Read-Only:

- `description` (String) The SQL Endpoint description.
- `display_name` (String) The SQL Endpoint display name.
- `id` (String) The SQL Endpoint ID.
- `workspace_id` (String) The Workspace ID.