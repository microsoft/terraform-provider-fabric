---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "fabric_dataflows Data Source - terraform-provider-fabric"
subcategory: ""
description: |-
  The Dataflows data-source allows you to retrieve a list of Fabric Dataflows https://learn.microsoft.com/rest/api/fabric/articles/item-management/definitions/dataflow-definition.
  -> This data-source supports Service Principal authentication.
---

# fabric_dataflows (Data Source)

The Dataflows data-source allows you to retrieve a list of Fabric [Dataflows](https://learn.microsoft.com/rest/api/fabric/articles/item-management/definitions/dataflow-definition).

-> This data-source supports Service Principal authentication.

## Example Usage

```terraform
data "fabric_dataflows" "example" {
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

- `values` (Attributes Set) The set of Dataflows. (see [below for nested schema](#nestedatt--values))

<a id="nestedatt--timeouts"></a>

### Nested Schema for `timeouts`

Optional:

- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).

<a id="nestedatt--values"></a>

### Nested Schema for `values`

Read-Only:

- `description` (String) The Dataflow description.
- `display_name` (String) The Dataflow display name.
- `id` (String) The Dataflow ID.
- `workspace_id` (String) The Workspace ID.
