---
page_title: "Creating an App Registration for the Service Principal context (SPN)"
subcategory: "Authentication"
description: |-

---

# Creating an App Registration for the Service Principal context (SPN)

---

## Create an App Registration

### Using Azure Portal

1. Sign in to the [Microsoft Entra admin center](https://entra.microsoft.com).
1. Browse to **Identity** > **Applications** > **App registrations** and select **New registration**.
1. Enter a display Name for your application.
1. Don't enter anything for **Redirect URI (optional)**

For more details and advanced scenarios, please follow this [guide](https://learn.microsoft.com/entra/identity-platform/quickstart-register-app#register-an-application).

### Using Azure CLI

Run the following commands to create an App Registration with Service Principal using [Azure CLI](https://learn.microsoft.com/cli/azure/):

```shell
#!/bin/bash

# Set input variables
fabricCapacityRgName="<FABRIC CAPACITY RESOURCE GROUP NAME>" # Resource group where the Fabric Capacity is located
fabricCapacityName="<FABRIC CAPACITY NAME>"                  # Name of the existing Fabric Capacity

# Install the Microsoft Fabric extension for Azure CLI
# See https://github.com/Azure/azure-cli-extensions/blob/main/src/microsoft-fabric/README.md for more details.
az extension add --name microsoft-fabric

# Get the current user ID
# See https://learn.microsoft.com/cli/azure/ad/signed-in-user#az-ad-signed-in-user-show for more details.
currentUserObjId=$(az ad signed-in-user show --output tsv --query id)

# Create a new Entra Application
# See https://learn.microsoft.com/cli/azure/ad/app#az-ad-app-create) for more details.
appObjId=$(az ad app create --display-name "Microsoft Fabric Terraform Provider" --sign-in-audience AzureADMyOrg --output tsv --query id)

# (optional, recommended) Add the current user as an owner of the application
# See https://learn.microsoft.com/cli/azure/ad/app/owner#az-ad-app-owner-add for more details.
az ad app owner add --id "${appObjId}" --owner-object-id "${currentUserObjId}"

# Create a new Entra Service Principal associated with the application
# see https://learn.microsoft.com/cli/azure/ad/sp#az-ad-sp-create for more details.
spObjId=$(az ad sp create --id "${appObjId}" --output tsv --query id)

# (optional, recommended) Add the current user as an owner of the service principal
# See https://learn.microsoft.com/cli/azure/reference-index#az-rest for more details.
az rest --method POST --url "https://graph.microsoft.com/v1.0/servicePrincipals/${spObjId}/owners/\$ref" --body "{\"@odata.id\": \"https://graph.microsoft.com/v1.0/users/${currentUserObjId}\"}"

# Get the Fabric Capacity
# See https://learn.microsoft.com/rest/api/microsoftfabric/fabric-capacities/get for more details.
fabricCapacity=$(az fabric capacity show --resource-group "${fabricCapacityRgName}" --capacity-name "${fabricCapacityName}" --output json)
fabricCapacityId=$(echo "${fabricCapacity}" | jq -r '.id')

# Assign Contributor role for the user-assigned managed identity to the Fabric Capacity
# See https://learn.microsoft.com/cli/azure/role/assignment#az-role-assignment-create for more details.
az role assignment create --assignee "${spObjId}" --role Contributor --scope "${fabricCapacityId}"

# Add a new principal to the the Fabric Capacity admin members
members=$(echo "${fabricCapacity}" | jq -c '.administration.members += ["'"${spObjId}"'"] | .administration')

# Update the Fabric Capacity with the new admin members
# See https://learn.microsoft.com/rest/api/microsoftfabric/fabric-capacities/update for more details.
az fabric capacity update --resource-group "${fabricCapacityRgName}" --capacity-name "${fabricCapacityName}" --administration "${members}"
```

### Using Entra PowerShell

Run the following commands to create an App Registration with Service Principal using [Entra PowerShell](https://learn.microsoft.com/powershell/entra-powershell/):

```powershell
# PowerShell

# Set input variables
$fabricCapacityRgName = '<FABRIC CAPACITY RESOURCE GROUP NAME>' # Resource group where the Fabric Capacity is located
$fabricCapacityName = '<FABRIC CAPACITY NAME>'                  # Name of the existing Fabric Capacity

# Login to Entra ID
# See https://learn.microsoft.com/powershell/module/microsoft.graph.entra/connect-entra
Connect-Entra -Scopes 'Application.ReadWrite.All', 'User.Read'

# Get the current context
# See https://learn.microsoft.com/powershell/module/microsoft.graph.entra/get-entracontext
$ctx = Get-EntraContext

# Get the current user
# See https://learn.microsoft.com/powershell/module/microsoft.graph.entra/get-entrauser
$currentUser = (Get-EntraUser -Filter "UserPrincipalName eq '$($ctx.Account)'" -Property Id)

# Create a new Entra Application
# See https://learn.microsoft.com/powershell/module/microsoft.graph.entra/new-entraapplication for more details.
$app = (New-EntraApplication -DisplayName 'Microsoft Fabric Terraform Provider' -SigninAudience AzureADMyOrg)

# (optional, recommended) Add the current user as an owner of the application
# See https://learn.microsoft.com/powershell/module/microsoft.graph.entra/add-entraapplicationowner for more details.
Add-EntraApplicationOwner -ObjectId $app.Id -RefObjectId $currentUser.Id

# Create a new Entra Service Principal associated with the application
# See https://learn.microsoft.com/powershell/module/microsoft.graph.entra/new-entraserviceprincipal for more details.
$sp = (New-EntraServicePrincipal -AppId $app.AppId)

# (optional, recommended) Add the current user as an owner of the service principal
# See https://learn.microsoft.com/powershell/module/microsoft.graph.entra/add-entraserviceprincipalowner for more details.
Add-EntraServicePrincipalOwner -ObjectId $sp.Id -RefObjectId $currentUser.Id

# Install the Az.Fabric module
# https://www.powershellgallery.com/packages/Az.Fabric
Install-Module -Name Az.Fabric

# Get current admin members and add a new principal to the array
# See https://learn.microsoft.com/powershell/module/az.fabric/get-azfabriccapacity for more details.
$members = (Get-AzFabricCapacity -ResourceGroupName $fabricCapacityRgName -CapacityName $fabricCapacityName).AdministrationMember
$members += $sp.Id

# Update the Fabric Capacity with the new admin members
# See https://learn.microsoft.com/powershell/module/az.fabric/update-azfabriccapacity for more details.
Update-AzFabricCapacity -ResourceGroupName $fabricCapacityRgName -CapacityName $fabricCapacityName -AdministrationMember $members
```

## Configure Microsoft Fabric to allow Service Principals (SPN) and Managed Identities (MSI)

1. Sign in to the [Microsoft Fabric admin portal](https://app.fabric.microsoft.com/admin-portal).
1. Browse to **Tenant settings** > **Developer settings** > [Service principals can use Fabric APIs](https://learn.microsoft.com/fabric/admin/service-admin-portal-developer#service-principals-can-use-fabric-apis) and check **Enable**.
1. Apply security restrictions to **The entire organization** or **Specific security groups**
