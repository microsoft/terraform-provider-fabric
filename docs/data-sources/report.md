---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "fabric_report Data Source - terraform-provider-fabric"
subcategory: ""
description: |-
  The Report data-source allows you to retrieve details about a Fabric Report https://learn.microsoft.com/power-bi/developer/projects/projects-report.
  -> This data-source supports Service Principal authentication.
---

# fabric_report (Data Source)

The Report data-source allows you to retrieve details about a Fabric [Report](https://learn.microsoft.com/power-bi/developer/projects/projects-report).

-> This data-source supports Service Principal authentication.

## Example Usage

```terraform
# Get item details
data "fabric_report" "example" {
  id           = "11111111-1111-1111-1111-111111111111"
  workspace_id = "00000000-0000-0000-0000-000000000000"
}

# Get item details with definition
data "fabric_report" "example_definition" {
  id                = "11111111-1111-1111-1111-111111111111"
  workspace_id      = "00000000-0000-0000-0000-000000000000"
  format            = "PBIR-Legacy"
  output_definition = true
}

# Access the content of the definition with JSONPath expression
output "example_definition_pbir_jsonpath" {
  value = provider::fabric::content_decode(data.fabric_report.example_definition.definition["definition.pbir"].content, ".datasetReference.byConnection.connectionString")
}
# Access the content of the definition as JSON object
output "example_definition_pbir_object" {
  value = provider::fabric::content_decode(data.fabric_report.example_definition.definition["definition.pbir"].content).datasetReference.byConnection.connectionString
}

output "example_definition_report_object" {
  value = provider::fabric::content_decode(data.fabric_report.example_definition.definition["report.json"].content)
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `id` (String) The Report ID.
- `workspace_id` (String) The Workspace ID.

### Optional

- `format` (String) The Report format. Possible values: `PBIR`, `PBIR-Legacy`
- `output_definition` (Boolean) Output definition parts as gzip base64 content? Default: `false`

!> Your terraform state file may grow a lot if you output definition content. Only use it when you must use data from the definition.

- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

### Read-Only

- `definition` (Attributes Map) Definition parts. Possible path keys: **PBIR** format: `StaticResources/RegisteredResources/*`, `StaticResources/SharedResources/*`, `definition.pbir`, `definition/bookmarks/*.json`, `definition/pages/*.json`, `definition/report.json`, `definition/version.json` **PBIR-Legacy** format: `StaticResources/RegisteredResources/*`, `StaticResources/SharedResources/*`, `definition.pbir`, `report.json` (see [below for nested schema](#nestedatt--definition))
- `description` (String) The Report description.
- `display_name` (String) The Report display name.

<a id="nestedatt--timeouts"></a>

### Nested Schema for `timeouts`

Optional:

- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).

<a id="nestedatt--definition"></a>

### Nested Schema for `definition`

Read-Only:

- `content` (String) Gzip base64 content of definition part.
Use [`provider::fabric::content_decode`](../functions/content_decode.md) function to decode content.
