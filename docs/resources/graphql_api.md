---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "fabric_graphql_api Resource - terraform-provider-fabric"
subcategory: ""
description: |-
  The GraphQL API resource allows you to manage a Fabric GraphQL API https://learn.microsoft.com/fabric/data-engineering/api-graphql-overview.
  -> This resource supports Service Principal authentication.
---

# fabric_graphql_api (Resource)

The GraphQL API resource allows you to manage a Fabric [GraphQL API](https://learn.microsoft.com/fabric/data-engineering/api-graphql-overview).

-> This resource supports Service Principal authentication.

## Example Usage

```terraform
resource "fabric_graphql_api" "example" {
  display_name = "example"
  workspace_id = "00000000-0000-0000-0000-000000000000"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `display_name` (String) The GraphQL API display name.
- `workspace_id` (String) The Workspace ID.

### Optional

- `description` (String) The GraphQL API description.
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

### Read-Only

- `id` (String) The GraphQL API ID.

<a id="nestedatt--timeouts"></a>

### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).

## Import

Import is supported using the following syntax:

The [`terraform import` command](https://developer.hashicorp.com/terraform/cli/commands/import) can be used, for example:

```shell
# terraform import fabric_graphql_api.example "<WorkspaceID>/<GraphQLApiID>"
terraform import fabric_graphql_api.example "00000000-0000-0000-0000-000000000000/11111111-1111-1111-1111-111111111111"
```
