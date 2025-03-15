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

## Terraform Logging

The Microsoft Fabric Terraform Provider outputs logs that you can enable by setting the `TF_LOG` environment variable to `DEBUG` or any other log level that Terraform supports.

By default, logs are sent to `stderr`. To send logs to a file, set the `TF_LOG_PATH` environment variable to the target file path.

For example, you can run the following command to enable logging at the debug level, and to output logs in monochrome format to a file named `tf.log` relative to the current working directory, while the `terraform apply` command runs:

```sh
# sh
TF_LOG=DEBUG TF_LOG_PATH=tf.log terraform apply
```

```powershell
# PowerShell
$env:TF_LOG="DEBUG"
$env:TF_LOG_PATH="tf.log"
terraform apply
```

For more information about Terraform logging, see [Debugging Terraform](https://developer.hashicorp.com/terraform/internals/debugging).

## Fabric API logging

Low-level logging is possible which will handle API calls. This type of logging can be very useful for debugging issues related to API interactions. By setting the logging level to `TRACE`, you can capture detailed information about the API calls made by Terraform. This includes basic request and response details, which can help in diagnosing problems or understanding the behavior of the API.

To enable low-level logging for API calls, you need to setup environment variables `TF_LOG` and `FABRIC_SDK_GO_LOGGING` with `TRACE` value.

```sh
# sh
TF_LOG=TRACE FABRIC_SDK_GO_LOGGING=TRACE terraform apply
```

```powershell
# PowerShell
$env:TF_LOG="TRACE"
$env:FABRIC_SDK_GO_LOGGING="TRACE"
terraform apply
```

## Advanced Fabric API logging

The provider supports additional logging options for deeper troubleshooting of API interactions:

### Include Request and Response Bodies

By default, request and response bodies are excluded from logs.

~> Enabling this can lead to disclosure of sensitive information, use with care.

For detailed debugging, you can include these bodies in your logs using the `FABRIC_SDK_GO_LOGGING_INCLUDE_BODY` environment variable:

```sh
# sh
TF_LOG=TRACE FABRIC_SDK_GO_LOGGING=TRACE FABRIC_SDK_GO_LOGGING_INCLUDE_BODY=true terraform apply
```

```powershell
# PowerShell
$env:TF_LOG="TRACE"
$env:FABRIC_SDK_GO_LOGGING="TRACE"
$env:FABRIC_SDK_GO_LOGGING_INCLUDE_BODY="true"
terraform apply
```

#### Expose HTTP Headers values

By default, headers are redacted in logs.

~> Enabling this can lead to disclosure of sensitive information, use with care.

You can specify which headers to include in logs using the `FABRIC_SDK_GO_LOGGING_ALLOWED_HEADERS` environment variable:

```sh
# sh
TF_LOG=TRACE FABRIC_SDK_GO_LOGGING=TRACE FABRIC_SDK_GO_LOGGING_ALLOWED_HEADERS="requestid;x-ms-operation-id;x-ms-public-api-error-code" terraform apply
```

```powershell
# PowerShell
$env:TF_LOG="TRACE"
$env:FABRIC_SDK_GO_LOGGING="TRACE"
$env:FABRIC_SDK_GO_LOGGING_ALLOWED_HEADERS="requestid;x-ms-operation-id;x-ms-public-api-error-code"
terraform apply
```

Headers must be specified in a semicolon-separated list and must contain valid header names (alphanumeric characters and hyphens only).

#### Combining Body and Header Logging

For the most comprehensive troubleshooting of API calls, you can combine both body and header logging options:

```sh
# sh
TF_LOG=TRACE FABRIC_SDK_GO_LOGGING=TRACE FABRIC_SDK_GO_LOGGING_INCLUDE_BODY=true FABRIC_SDK_GO_LOGGING_ALLOWED_HEADERS="requestid;x-ms-operation-id;content-type" terraform apply
```

```powershell
$env:TF_LOG="TRACE"
$env:FABRIC_SDK_GO_LOGGING="TRACE"
$env:FABRIC_SDK_GO_LOGGING_INCLUDE_BODY="true"
$env:FABRIC_SDK_GO_LOGGING_ALLOWED_HEADERS="requestid;x-ms-operation-id;content-type"
terraform apply
```

~> Combining these options produces very verbose logs that include sensitive information. Use this configuration only in a secure environment and make sure to redact any sensitive data before sharing logs.

## FAQ

### I am getting error `The feature is not available`

- Check if your SPN, MSI or User that is used for Provider authentication is added to Fabric `Capacity administrators`.
- Check if your Fabric Capacity is not in the `paused` state.
- Majority of Fabric Items require to have Fabric Capacity assigned to the workspace. If you manage Workspace using [`fabric_workspace`](../resources/workspace.md) resource, ensure that you have assigned Fabric Capacity to the Workspace.

### I am getting error `Unable to find Capacity...`

- Check if your SPN, MSI or User that is used for Provider authentication is added to Fabric `Capacity administrators`.
- Check if your Fabric Capacity is not in the `paused` state.

### I am getting error `Failed to create workspace identity`

- If you manage Workspace using [`fabric_workspace`](../resources/workspace.md) resource, ensure that you have assigned Fabric Capacity to the Workspace.
- Check if your Fabric Capacity assigned to the Workspace is not in the `paused` state.

### I am getting error `Workspace name already exists`

- Ensure that you have provided the unique name for the Workspace that does not exist in the Fabric yet.

### Operations take too long to complete or timeout

You can observe some Terraform operations take time to complete with the messages like `Still creating...`, `Still reading...`, etc. or end up with a timeout error. This can happen due to various reasons such as network latency or [Fabric API throttling](https://learn.microsoft.com/rest/api/fabric/articles/throttling).

- Try to increase the global timeout for the operations by setting the [`timeout`](../index.md#timeout) attribute in the Provider block, or you can set the timeout for the specific Resource or Data-Source using the `timeouts` attribute.
- Change [Terraform parallelism](https://developer.hashicorp.com/terraform/internals/graph#walking-the-graph) to lower number than default (10x).
