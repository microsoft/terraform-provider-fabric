---
page_title: "Creating an App Registration for the User context"
subcategory: "Authentication"
description: |-

---

# Creating an App Registration for the User context

---

## Create an App Registration

You can follow this [guide](https://learn.microsoft.com/entra/identity-platform/quickstart-register-app#register-an-application) to create an app registration.

## Set API Permissions

In the **API permissions** menu of your App Registration, add required API permissions to use with the Fabric Terraform provider:

-> Depends on your solution you may add various set of API permissions. The below one are just an example.

- Power BI Service
  - App.Read.All
  - Capacity.ReadWrite.All
  - ...
  - ...
  - Workspace.Read.Al
  - Workspace.ReadWrite.All

- Microsoft Graph
  - User.Read
  - User.ReadBasic.All

## Expose API

In the **Expose an API** menu of your App Registration, you need to define your application ID URI:

- Application ID URI: `api://<client_id>`, for example:

```text
api://fabric_terraform_provider
```

- Add required scope in the `Scopes defined by this API` section:

1. Scope name: `default`
1. Who can consent: `Admins and users`
1. Admin consent display name: `Fabric Terraform Provider`
1. Admin consent description: `Allows connection to backend services for Fabric Terraform Provider`
1. User consent display name: `Fabric Terraform Provider`
1. User consent description: `Allows connection to backend services for Fabric Terraform Provider`
1. State: `Enabled`

- You will finally need to pre-authorize Azure CLI/Azure PowerShell and Power BI to access your exposed API permissions by adding Azure CLI/Azure PowerShell and Power BI 1st party Microsoft applications. In the `Authorized client applications` section add:
  - for Azure CLI: `04b07795-8ddb-461a-bbee-02f9e1bf7b46`
  - for Azure PowerShell: `1950a258-227b-4e31-a9cf-717495945fc2`
  - for Power BI: `00000009-0000-0000-c000-000000000000` and `871c010f-5e61-4fb1-83ac-98610a7e9110`

Read more about first-party Microsoft applications ine the [Application IDs of commonly used Microsoft applications](https://learn.microsoft.com/troubleshoot/azure/entra/entra-id/governance/verify-first-party-apps-sign-in#application-ids-of-commonly-used-microsoft-applications) article.

## Usage with Azure CLI

After above steps you should be able to authenticate using [Azure CLI](https://learn.microsoft.com/cli/azure/):

```shell
# (optional, useful in the multi-tenant scenarios) Disable the new login experience
# See https://learn.microsoft.com/cli/azure/authenticate-azure-cli-interactively#sign-in-with-a-different-tenant for more details.
az config set core.login_experience_v2=off

# (optional, Windows only, useful in the multi-tenant scenarios) Disable WAM on Windows
# See https://learn.microsoft.com/cli/azure/authenticate-azure-cli-interactively#sign-in-with-web-account-manager-wam-on-windows for more details.
az config set core.enable_broker_on_windows=false

# Login to Azure with Entra ID credentials
# See https://learn.microsoft.com/cli/azure/authenticate-azure-cli for more details.
az login --allow-no-subscriptions --tenant 00000000-0000-0000-0000-000000000000 --scope api://fabric_terraform_provider/default
```

The following Fabric Provider block can be specified to use the Azure CLI:

```terraform
provider "fabric" {
  use_cli = true
}
```

## Usage with Azure PowerShell

After above steps you should be able to authenticate using [Azure PowerShell](https://learn.microsoft.com/powershell/azure/):

```powershell
# (optional, useful in the multi-tenant scenarios) Disable the new login experience
# See https://learn.microsoft.com/powershell/azure/authenticate-interactive#disable-the-new-login-experience for more details.
Update-AzConfig -LoginExperienceV2 Off

# (optional, Windows only, useful in the multi-tenant scenarios) Disable WAM on Windows
# See https://learn.microsoft.com/powershell/azure/authenticate-interactive#web-account-manager-wam for more details.
Update-AzConfig -EnableLoginByWam $false

# Login to Azure with Entra ID credentials
# See https://learn.microsoft.com/powershell/module/az.accounts/connect-azaccount for more details.
Connect-AzAccount -Tenant '00000000-0000-0000-0000-000000000000' -AuthScope 'api://fabric_terraform_provider'

# Set the FABRIC_TOKEN environment variable to the access token
# See https://learn.microsoft.com/powershell/module/az.accounts/get-azaccesstoken for more details.
$env:FABRIC_TOKEN = (Get-AzAccessToken -ResourceUrl 'https://api.fabric.microsoft.com').Token
```

The following Fabric Provider block can be specified to use the Azure PowerShell:

```terraform
provider "fabric" {}
```
