---
page_title: "Authenticating using a Service Principal and OpenID Connect (OIDC)"
subcategory: "Authentication"
description: |-

---

# Authenticating using a Service Principal and OpenID Connect (OIDC)

---

Federated identity credentials are a type of credential that allows workloads, such as GitHub Actions, workloads running on Kubernetes, or workloads running in compute platforms outside of Azure access Microsoft Entra protected resources without needing to manage secrets using workload identity federation.

## Setting up Entra Application and Service Principal

Follow [Creating an App Registration for the Service Principal context (SPN)](./auth_app_reg_spn.md) guide.

## Configure Entra App to trust a GitHub repository

### Using Azure Portal (GitHub)

1. In the [Microsoft Entra admin center](https://entra.microsoft.com), in **App registrations**, select your application.
1. Select **Certificates & secrets** > **Federated credentials** > **Add credential**.
1. In the **Federated credential scenario** drop-down box, select **GitHub Actions deploying Azure resources**.
1. Specify the **Organization** and **Repository** for your GitHub Actions workflow. For **Entity type**, select **Environment**, **Branch**, **Pull Request**, or **Tag** and specify the value. The values must exactly match the configuration in the GitHub workflow. For our example, let's select **Branch** and specify `main`.
1. Add a **Name** for the federated credential.
1. The **Issuer**, **Audiences**, and **Subject identifier** fields auto-populate based on the values you entered.
1. Click **Add** to configure the federated credential.

### Using Azure CLI (GitHub)

```shell
# Create application federated identity credential
az ad app federated-credential create --id "00000000-0000-0000-0000-000000000000" --parameters credential.json
```

Where the `credential.json` contains the following content:

```json
{
  "name": "branch-main",
  "issuer": "https://token.actions.githubusercontent.com",
  "subject": "repo:your-github-org/your-github-repo:refs:refs/heads/main",
  "description": "Deployments from the main branch",
  "audiences": [
    "api://AzureADTokenExchange"
  ]
}
```

See the [official documentation](https://learn.microsoft.com/cli/azure/ad/app/federated-credential?view=azure-cli-latest#az-ad-app-federated-credential-create) for more details.

## Configure Entra App to trust a Generic OIDC issuer

### Using Azure Portal (Generic)

1. In the [Microsoft Entra admin center](https://entra.microsoft.com), in **App registrations**, select your application.
1. Select **Certificates & secrets** > **Federated credentials** > **Add credential**.
1. In the **Federated credential scenario** drop-down box, select **Other issuer**.
1. Refer to the instructions from your OIDC provider for completing the form, before choosing a **Name** for the federated credential and clicking the **Add** button.

## Configuring Terraform to use the OIDC

Now that we have our federated credential for Entra App and ready to use, it's possible to configure Terraform in a few different ways.

### Environment Variables

```shell
# sh
export FABRIC_USE_OIDC=true
export FABRIC_TENANT_ID="00000000-0000-0000-0000-000000000000"
export FABRIC_CLIENT_ID="00000000-0000-0000-0000-000000000000"
```

```powershell
# PowerShell
$env:FABRIC_USE_OIDC = $true
$env:FABRIC_TENANT_ID = '00000000-0000-0000-0000-000000000000'
$env:FABRIC_CLIENT_ID = '00000000-0000-0000-0000-000000000000'
```

#### OIDC token

The provider will use the `FABRIC_OIDC_TOKEN` environment variable as an OIDC token. You can use this variable to specify the token provided by your OIDC provider. If your OIDC provider provides an ID token in a file, you can specify the path to this file with the `FABRIC_OIDC_TOKEN_FILE_PATH` environment variable.

#### GitHub Actions

When running in GitHub Actions, the provider will detect the `ACTIONS_ID_TOKEN_REQUEST_URL` and `ACTIONS_ID_TOKEN_REQUEST_TOKEN` environment variables set by the GitHub Actions runtime. You can also specify the `FABRIC_OIDC_REQUEST_TOKEN` and `FABRIC_OIDC_REQUEST_URL` environment variables.

For GitHub Actions workflows, you'll need to ensure the workflow has `write` permissions for the `id-token`.

```yaml
permissions:
  id-token: write
  contents: read
```

For more information about OIDC in GitHub Actions, see [official documentation](https://docs.github.com/actions/deployment/security-hardening-your-deployments/configuring-openid-connect-in-cloud-providers).

#### Azure DevOps Pipelines

When running in Azure DevOps Pipelines, the provider will detect the `SYSTEM_ACCESSTOKEN` environment variable set by the Azure DevOps runtime. You can also specify the `FABRIC_OIDC_REQUEST_TOKEN` environment variables.

```yaml
steps:
  # Bash example
  - bash: terraform apply -auto-approve
    env:
      FABRIC_OIDC_REQUEST_TOKEN: $(System.AccessToken) # or SYSTEM_ACCESSTOKEN: $(System.AccessToken)
      FABRIC_AZURE_DEVOPS_SERVICE_CONNECTION_ID: "your-service-connection-id"

  # PowerShell example
  - powershell: terraform apply -auto-approve
    env:
      FABRIC_OIDC_REQUEST_TOKEN: $(System.AccessToken) # or SYSTEM_ACCESSTOKEN: $(System.AccessToken)
      FABRIC_AZURE_DEVOPS_SERVICE_CONNECTION_ID: "your-service-connection-id"
```

For more information about OIDC in Azure DevOps Pipelines, see:

- [Create an Azure Resource Manager service connection that uses workload identity federation](https://learn.microsoft.com/azure/devops/pipelines/library/connect-to-azure?view=azure-devops#create-an-azure-resource-manager-service-connection-that-uses-workload-identity-federation)
- [System.AccessToken](https://learn.microsoft.com/azure/devops/pipelines/build/variables?view=azure-devops&tabs=yaml#systemaccesstoken).

### Provider Block

The following Terraform and Provider blocks can be specified, where `0.0.0-preview` is the version of the Fabric Provider that you'd like to use:

```terraform
# We strongly recommend using the required_providers block to set the Fabric Provider source and version being used
terraform {
  required_version = ">= 1.8, < 2.0"
  required_providers {
    fabric = {
      source  = "microsoft/fabric"
      version = "0.0.0" # Check for the latest version on the Terraform Registry
    }
  }
}

# Configure the Microsoft Fabric Provider
provider "fabric" {
  use_oidc = true
}
```
