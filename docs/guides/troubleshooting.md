---
page_title: "Troubleshooting"
subcategory: "Guides"
description: |-
  Troubleshooting Guide.
---

# Troubleshooting

The guide provides troubleshooting steps for common issues that you might encounter while using the Microsoft Fabric Terraform Provider.

---

-> For Terraform-specific support, see the latest Terraform topics on the [HashiCorp Discuss website](https://discuss.hashicorp.com/). For issues specific to the Microsoft Fabric Terraform Provider, see Issues in the [microsoft/terraform-provider-fabric](https://github.com/microsoft/terraform-provider-fabric) GitHub repository.

---

## Logging

The Microsoft Fabric Terraform Provider outputs logs that you can enable by setting the `TF_LOG` environment variable to `DEBUG` or any other log level that Terraform supports.

By default, logs are sent to `stderr`. To send logs to a file, set the `TF_LOG_PATH` environment variable to the target file path.

For example, you can run the following command to enable logging at the debug level, and to output logs in monochrome format to a file named `tf.log` relative to the current working directory, while the `terraform apply` command runs:

```bash
TF_LOG=DEBUG TF_LOG_PATH=tf.log terraform apply -no-color
```

For more information about Terraform logging, see [Debugging Terraform](https://developer.hashicorp.com/terraform/internals/debugging).

## FAQ

### I am getting error `The feature is not available`

- Check if your SPN, MSI or User that is used for Provider authentication is added to Fabric `Capacity administrators`.
- Check if your Fabric Capacity is not in the `paused` state.
- Majority of Fabric Items require to have Fabric Capacity assigned to the workspace. If you manage Workspace using [`fabric_workspace` resource](../resources/workspace), ensure that you have assigned Fabric Capacity to the Workspace.

### I am getting error `Unable to find Capacity...`

- Check if your SPN, MSI or User that is used for is added to Fabric `Capacity administrators`.
- Check if your Fabric Capacity is not in the `paused` state.

### I am getting error `Failed to create workspace identity`

- If you manage Workspace using [`fabric_workspace` resource](../resources/workspace), ensure that you have assigned Fabric Capacity to the Workspace.
- Check if your Fabric Capacity assigned to the Workspace is not in the `paused` state.

### I am getting error `Workspace name already exists`

- Ensure that you have provided the unique name for the Workspace.
