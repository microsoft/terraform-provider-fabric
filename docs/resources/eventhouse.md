---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "fabric_eventhouse Resource - terraform-provider-fabric"
subcategory: ""
description: |-
  This resource manages a Fabric Eventhouse.
  See Eventhouse https://learn.microsoft.com/fabric/real-time-intelligence/eventhouse for more information.
  -> This item supports Service Principal authentication.
---

# fabric_eventhouse (Resource)

This resource manages a Fabric Eventhouse.

See [Eventhouse](https://learn.microsoft.com/fabric/real-time-intelligence/eventhouse) for more information.

-> This item supports Service Principal authentication.

## Example Usage

```terraform
resource "fabric_eventhouse" "example" {
  display_name = "example"
  workspace_id = "00000000-0000-0000-0000-000000000000"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `display_name` (String) The Eventhouse display name.
- `workspace_id` (String) The Workspace ID.

### Optional

- `description` (String) The Eventhouse description.
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

### Read-Only

- `id` (String) The Eventhouse ID.
- `properties` (Attributes) The Eventhouse properties. (see [below for nested schema](#nestedatt--properties))

<a id="nestedatt--timeouts"></a>

### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).

<a id="nestedatt--properties"></a>

### Nested Schema for `properties`

Read-Only:

- `database_ids` (List of String) The IDs list of KQL Databases.
- `ingestion_service_uri` (String) Ingestion service URI.
- `query_service_uri` (String) Query service URI.

## Import

Import is supported using the following syntax:

```shell
# terraform import fabric_eventhouse.example "<WorkspaceID>/<EventhouseID>"
terraform import fabric_eventhouse.example "00000000-0000-0000-0000-000000000000/11111111-1111-1111-1111-111111111111"
```
