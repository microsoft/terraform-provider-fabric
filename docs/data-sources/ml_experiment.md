---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "fabric_ml_experiment Data Source - terraform-provider-fabric"
subcategory: ""
description: |-
  The ML Experiment data-source allows you to retrieve details about a Fabric ML Experiment https://learn.microsoft.com/fabric/data-science/machine-learning-experiment.
  -> This data-source does not support Service Principal. Please use a User context authentication.
  ~> This data-source is in preview. To access it, you must explicitly enable the preview mode in the provider level configuration.
---

# fabric_ml_experiment (Data Source)

The ML Experiment data-source allows you to retrieve details about a Fabric [ML Experiment](https://learn.microsoft.com/fabric/data-science/machine-learning-experiment).

-> This data-source does not support Service Principal. Please use a User context authentication.

~> This data-source is in **preview**. To access it, you must explicitly enable the `preview` mode in the provider level configuration.

## Example Usage

```terraform
data "fabric_ml_experiment" "example_by_id" {
  id           = "11111111-1111-1111-1111-111111111111"
  workspace_id = "00000000-0000-0000-0000-000000000000"
}

data "fabric_ml_experiment" "example_by_name" {
  display_name = "example"
  workspace_id = "00000000-0000-0000-0000-000000000000"
}

# This is an invalid data source
# Do not specify `id` and `display_name` in the same data source block
# data "fabric_ml_experiment" "example" {
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

- `display_name` (String) The ML Experiment display name.
- `id` (String) The ML Experiment ID.
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

### Read-Only

- `description` (String) The ML Experiment description.

<a id="nestedatt--timeouts"></a>

### Nested Schema for `timeouts`

Optional:

- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
