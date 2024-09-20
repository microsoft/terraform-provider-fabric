---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "fabric_spark_job_definition Data Source - terraform-provider-fabric"
subcategory: ""
description: |-
  Get a Fabric Spark Job Definition.
  Use this data source to fetch a Spark Job Definition https://learn.microsoft.com/fabric/data-engineering/spark-job-definition.
  -> This item supports Service Principal authentication.
---

# fabric_spark_job_definition (Data Source)

Get a Fabric Spark Job Definition.

Use this data source to fetch a [Spark Job Definition](https://learn.microsoft.com/fabric/data-engineering/spark-job-definition).

-> This item supports Service Principal authentication.

## Example Usage

```terraform
# Get item details by id
data "fabric_spark_job_definition" "example_by_id" {
  id           = "11111111-1111-1111-1111-111111111111"
  workspace_id = "00000000-0000-0000-0000-000000000000"
}

# Get item details by name
data "fabric_spark_job_definition" "example_by_name" {
  display_name = "test1"
  workspace_id = "00000000-0000-0000-0000-000000000000"
}

# Get item details with definition
# Examples uses `id` but `display_name` can be used as well
data "fabric_spark_job_definition" "example_definition" {
  id                = "11111111-1111-1111-1111-111111111111"
  workspace_id      = "00000000-0000-0000-0000-000000000000"
  output_definition = true
}

# Access the content of the definition with JSONPath expression
output "example_definition_content_jsonpath" {
  value = provider::fabric::content_decode(data.fabric_spark_job_definition.example_definition.definition["SparkJobDefinitionV1.json"].content, ".defaultLakehouseArtifactId")
}

# Access the content of the definition as JSON object
output "example_definition_content_object" {
  value = provider::fabric::content_decode(data.fabric_spark_job_definition.example_definition.definition["SparkJobDefinitionV1.json"].content).defaultLakehouseArtifactId
}

# This is an invalid data source
# Do not specify `id` and `display_name` in the same data source block
# data "fabric_spark_job_definition" "example" {
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

- `display_name` (String) The Spark Job Definition display name.
- `id` (String) The Spark Job Definition ID.
- `output_definition` (Boolean) Output definition parts as gzip base64 content? Default: `false`

!> Your terraform state file may grow a lot if you output definition content. Only use it when you must use data from the definition.

- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

### Read-Only

- `definition` (Attributes Map) Definition parts. Possible path keys: `SparkJobDefinitionV1.json`. (see [below for nested schema](#nestedatt--definition))
- `description` (String) The Spark Job Definition description.
- `format` (String) The Spark Job Definition format. Possible values: `SparkJobDefinitionV1`.

<a id="nestedatt--timeouts"></a>

### Nested Schema for `timeouts`

Optional:

- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).

<a id="nestedatt--definition"></a>

### Nested Schema for `definition`

Read-Only:

- `content` (String) Gzip base64 content of definition part.
Use [`provider::fabric::content_decode`](../functions/content_decode.md) function to decode content.