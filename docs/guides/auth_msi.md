---
page_title: "Authenticating using a Managed Identity (MSI)"
subcategory: "Authentication"
description: |-

---

# Authenticating using a Managed Identity (MSI)

---

## Setting up Microsoft Fabric to allow Service Principals

Follow [Configure Microsoft Fabric to allow Service Principals (SPN) and Managed Identities (MSI)](./auth_app_reg_spn.md#configure-microsoft-fabric-to-allow-service-principals-spn-and-managed-identities-msi) guide.

## Using a System-Assigned Managed Identity

### Configuring a Virtual Machine to use a System-Assigned Managed Identity

```shell
# Assign the system-assigned managed identity to the VM
# See https://learn.microsoft.com/cli/azure/vm/identity#az-vm-identity-assign for more details.
az vm identity assign --resource-group "<RESOURCE GROUP NAME>" --name "<VM NAME>" --identities system

# Assign Contributor role for the system-assigned managed identity to the Fabric Capacity
# See https://learn.microsoft.com/cli/azure/role/assignment#az-role-assignment-create for more details.
az role assignment create --assignee "<PRINCIPAL ID>" --role Contributor --scope "/subscriptions/<SUBSCRIPTION ID>/resourceGroups/<RESOURCE GROUP NAME>/providers/Microsoft.Fabric/capacities/<FABRIC CAPACITY NAME>"

# Get current admin members and add a new principal to the array
# See https://learn.microsoft.com/rest/api/microsoftfabric/fabric-capacities/get for more details.
members=$(az rest --method get --uri "/subscriptions/<SUBSCRIPTION ID>/resourceGroups/<RESOURCE GROUP NAME>/providers/Microsoft.Fabric/capacities/<FABRIC CAPACITY NAME>?api-version=2023-11-01" --output json --query properties.administration.members | jq '. += ["<PRINCIPAL ID>"]')

# Update the Fabric Capacity with the new admin members
# See https://learn.microsoft.com/rest/api/microsoftfabric/fabric-capacities/update for more details.
az rest --method patch --uri "/subscriptions/<SUBSCRIPTION ID>/resourceGroups/<RESOURCE GROUP NAME>/providers/Microsoft.Fabric/capacities/<FABRIC CAPACITY NAME>?api-version=2023-11-01" --body "{\"properties\":{\"administration\":{\"members\":${members}}}}"
```

```powershell
# See https://learn.microsoft.com/entra/identity/managed-identities-azure-resources/how-to-configure-managed-identities?pivots=qs-configure-powershell-windows-vm#system-assigned-managed-identity for more details.

# Assign the system-assigned managed identity to the VM
# See https://learn.microsoft.com/powershell/module/az.compute/get-azvm for more details.
$vm = Get-AzVM -ResourceGroupName "<RESOURCE GROUP NAME>" -Name "<VM NAME>"
# See https://learn.microsoft.com/powershell/module/az.compute/update-azvm for more details.
Update-AzVM -ResourceGroupName "<RESOURCE GROUP NAME>" -VM $vm -IdentityType SystemAssigned

# Assign Contributor role for the system-assigned managed identity to the Fabric Capacity
# See https://learn.microsoft.com/azure/role-based-access-control/role-assignments-powershell for more details.
New-AzRoleAssignment -ObjectId "<PRINCIPAL ID>" -RoleDefinitionName Contributor -Scope "/subscriptions/<SUBSCRIPTION ID>/resourceGroups/<RESROURCE GROUP NAME>/providers/Microsoft.Fabric/capacities/<FABRIC CAPACITY NAME>"

# Get current admin members and add a new principal to the array
# See https://learn.microsoft.com/powershell/module/az.accounts/invoke-azrestmethod for more details.
$members = ((Invoke-AzRestMethod -Method GET -Path "/subscriptions/<SUBSCRIPTION ID>/resourceGroups/<RESROURCE GROUP NAME>/providers/Microsoft.Fabric/capacities/<FABRIC CAPACITY NAME>?api-version=2023-11-01").Content | ConvertFrom-Json).properties.administration.members += "<PRINCIPAL ID>"

# Update the Fabric Capacity with the new admin members
$payload = @{
 properties = @{
  administration = @{
   members = $members
  }
 }
} | ConvertTo-Json -Depth 10
Invoke-AzRestMethod -Method PATCH -Path "/subscriptions/<SUBSCRIPTION ID>/resourceGroups/<RESROURCE GROUP NAME>/providers/Microsoft.Fabric/capacities/<FABRIC CAPACITY NAME>?api-version=2023-11-01" -Payload $payload
```

### Configuring Terraform to use the System-Assigned Managed Identity

At this point we assume that managed identity is configured on the resource (e.g. virtual machine) being used, that permissions have been granted, and that you are running Terraform on that resource.

Terraform can be configured to use managed identity for authentication in one of two ways: using Environment Variables or by defining the fields within the Provider Block.

#### Environment Variables setup for the System-Assigned Managed Identity

You can configure Terraform to use Managed Identity by setting the Environment Variable `FABRIC_USE_MSI` to `true`; as shown below:

```shell
# sh
export FABRIC_USE_MSI=true
export FABRIC_TENANT_ID="00000000-0000-0000-0000-000000000000"
```

```powershell
# PowerShell
$env:FABRIC_USE_MSI = $true
$env:FABRIC_TENANT_ID = '00000000-0000-0000-0000-000000000000'
```

#### Provider Block setup for the System-Assigned Managed Identity

The following Terraform and Provider blocks can be specified, where `0.0.0` is the version of the Fabric Provider that you'd like to use:

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
  use_msi   = true
  tenant_id = "00000000-0000-0000-0000-000000000000"
}
```

## Using a User-Assigned Managed Identity

### Configuring a Virtual Machine to use a User-Assigned Managed Identity

```shell
# Create a user-assigned managed identity
# See https://learn.microsoft.com/cli/azure/identity#az-identity-create for more details.
az identity create --resource-group "<RESOURCE GROUP NAME>" --name "<IDENTITY NAME>"

# Assign the user-assigned managed identity to the VM
# See https://learn.microsoft.com/cli/azure/vm/identity#az-vm-identity-assign for more details.
az vm identity assign --resource-group "<RESOURCE GROUP NAME>" --name "<VM NAME>" --identities "<IDENTITY ID>"

# Assign Contributor role for the user-assigned managed identity to the Fabric Capacity
# See https://learn.microsoft.com/cli/azure/role/assignment#az-role-assignment-create for more details.
az role assignment create --assignee "<PRINCIPAL ID>" --role Contributor --scope "/subscriptions/<SUBSCRIPTION ID>/resourceGroups/<RESOURCE GROUP NAME>/providers/Microsoft.Fabric/capacities/<FABRIC CAPACITY NAME>"

# Get current admin members and add a new principal to the array
# See https://learn.microsoft.com/rest/api/microsoftfabric/fabric-capacities/get for more details.
members=$(az rest --method get --uri "/subscriptions/<SUBSCRIPTION ID>/resourceGroups/<RESOURCE GROUP NAME>/providers/Microsoft.Fabric/capacities/<FABRIC CAPACITY NAME>?api-version=2023-11-01" --output json --query properties.administration.members | jq '. += ["<PRINCIPAL ID>"]')

# Update the Fabric Capacity with the new admin members
# See https://learn.microsoft.com/rest/api/microsoftfabric/fabric-capacities/update for more details.
az rest --method patch --uri "/subscriptions/<SUBSCRIPTION ID>/resourceGroups/<RESOURCE GROUP NAME>/providers/Microsoft.Fabric/capacities/<FABRIC CAPACITY NAME>?api-version=2023-11-01" --body "{\"properties\":{\"administration\":{\"members\":${members}}}}"
```

```powershell
# See https://learn.microsoft.com/entra/identity/managed-identities-azure-resources/how-to-configure-managed-identities?pivots=qs-configure-powershell-windows-vm#user-assigned-managed-identity for more details.

# Create a user-assigned managed identity
# See https://learn.microsoft.com/powershell/module/az.managedserviceidentity/new-azuserassignedidentity for more details.
New-AzUserAssignedIdentity -ResourceGroupName "<RESROURCE GROUP NAME>" -Name "<USER ASSIGNED IDENTITY NAME>"

# Assign the user-assigned managed identity to the VM
# See https://learn.microsoft.com/powershell/module/az.compute/get-azvm for more details.
$vm = Get-AzVM -ResourceGroupName "<RESROURCE GROUP NAME>" -Name "<VM NAME>"
# See https://learn.microsoft.com/powershell/module/az.compute/update-azvm for more details.
Update-AzVM -ResourceGroupName "<RESROURCE GROUP NAME>" -VM $vm -IdentityType UserAssigned -IdentityID "/subscriptions/<SUBSCRIPTION ID>/resourcegroups/<RESROURCE GROUP NAME>/providers/Microsoft.ManagedIdentity/userAssignedIdentities/<USER ASSIGNED IDENTITY NAME>"

# Assign Contributor role for the user-assigned managed identity to the Fabric Capacity
# See https://learn.microsoft.com/azure/role-based-access-control/role-assignments-powershell for more details.
New-AzRoleAssignment -ObjectId "<PRINCIPAL ID>" -RoleDefinitionName Contributor -Scope "/subscriptions/<SUBSCRIPTION ID>/resourceGroups/<RESROURCE GROUP NAME>/providers/Microsoft.Fabric/capacities/<FABRIC CAPACITY NAME>"

# Get current admin members and add a new principal to the array
# See https://learn.microsoft.com/powershell/module/az.accounts/invoke-azrestmethod for more details.
$members = ((Invoke-AzRestMethod -Method GET -Path "/subscriptions/<SUBSCRIPTION ID>/resourceGroups/<RESROURCE GROUP NAME>/providers/Microsoft.Fabric/capacities/<FABRIC CAPACITY NAME>?api-version=2023-11-01").Content | ConvertFrom-Json).properties.administration.members += "<PRINCIPAL ID>"

# Update the Fabric Capacity with the new admin members
$payload = @{
 properties = @{
  administration = @{
   members = $members
  }
 }
} | ConvertTo-Json -Depth 10
Invoke-AzRestMethod -Method PATCH -Path "/subscriptions/<SUBSCRIPTION ID>/resourceGroups/<RESROURCE GROUP NAME>/providers/Microsoft.Fabric/capacities/<FABRIC CAPACITY NAME>?api-version=2023-11-01" -Payload $payload
```

### Configuring Terraform to use the User-Assigned Managed Identity

At this point we assume that managed identity is configured on the resource (e.g. virtual machine) being used, that permissions have been granted, and that you are running Terraform on that resource.

Terraform can be configured to use managed identity for authentication in one of two ways: using Environment Variables or by defining the fields within the Provider block.

#### Environment Variables setup for the User-Assigned Managed Identity

You can configure Terraform to use Managed Identity by setting the Environment Variable `FABRIC_USE_MSI` to `true`; as shown below:

```shell
# sh
export FABRIC_USE_MSI=true
export FABRIC_TENANT_ID="00000000-0000-0000-0000-000000000000"
export FABRIC_CLIENT_ID="00000000-0000-0000-0000-000000000000"
```

```powershell
# PowerShell
$env:FABRIC_USE_MSI = $true
$env:FABRIC_TENANT_ID = '00000000-0000-0000-0000-000000000000'
$env:FABRIC_CLIENT_ID = '00000000-0000-0000-0000-000000000000'
```

#### Provider Block setup for the User-Assigned Managed Identity

The following Terraform and Provider blocks can be specified, where `0.0.0` is the version of the Fabric Provider that you'd like to use:

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
  use_msi   = true
  tenant_id = "00000000-0000-0000-0000-000000000000"
  client_id = "00000000-0000-0000-0000-000000000000"
}
```
