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
#!/bin/bash

# Set input variables
vmRgName="<VM RESOURCE GROUP NAME>"                          # Resource group where the VM is located
vmName="<VM NAME>"                                           # Name of the VM
fabricCapacityRgName="<FABRIC CAPACITY RESOURCE GROUP NAME>" # Resource group where the Fabric Capacity is located
fabricCapacityName="<FABRIC CAPACITY NAME>"                  # Name of the existing Fabric Capacity

# Install the Microsoft Fabric extension for Azure CLI
# See https://github.com/Azure/azure-cli-extensions/blob/main/src/microsoft-fabric/README.md for more details.
az extension add --name microsoft-fabric

# Assign the system-assigned managed identity to the VM
# See https://learn.microsoft.com/cli/azure/vm/identity#az-vm-identity-assign for more details.
identityPrincipalId=$(az vm identity assign --resource-group "${vmRgName}" --name "${vmName}" --identities "[system]" --output tsv --query systemAssignedIdentity)

# Get the Fabric Capacity
# See https://learn.microsoft.com/rest/api/microsoftfabric/fabric-capacities/get for more details.
fabricCapacity=$(az fabric capacity show --resource-group "${fabricCapacityRgName}" --capacity-name "${fabricCapacityName}" --output json)
fabricCapacityId=$(echo "${fabricCapacity}" | jq -r '.id')

# Assign Contributor role for the system-assigned managed identity to the Fabric Capacity
# See https://learn.microsoft.com/cli/azure/role/assignment#az-role-assignment-create for more details.
az role assignment create --assignee "${identityPrincipalId}" --role Contributor --scope "${fabricCapacityId}"

# Add a new principal to the the Fabric Capacity admin members
members=$(echo "${fabricCapacity}" | jq -c '.administration.members += ["'"${identityPrincipalId}"'"] | .administration')

# Update the Fabric Capacity with the new admin members
# See https://learn.microsoft.com/rest/api/microsoftfabric/fabric-capacities/update for more details.
az fabric capacity update --resource-group "${fabricCapacityRgName}" --capacity-name "${fabricCapacityName}" --administration "${members}"
```

```powershell
# PowerShell

# See https://learn.microsoft.com/entra/identity/managed-identities-azure-resources/how-to-configure-managed-identities?pivots=qs-configure-powershell-windows-vm#system-assigned-managed-identity for more details.

# Set input variables
$vmRgName = '<VM RESOURCE GROUP NAME>'                          # Resource group where the VM is located
$vmName = '<VM NAME>'                                           # Name of the VM
$fabricCapacityRgName = '<FABRIC CAPACITY RESOURCE GROUP NAME>' # Resource group where the Fabric Capacity is located
$fabricCapacityName = '<FABRIC CAPACITY NAME>'                  # Name of the existing Fabric Capacity

# Install the Az.Fabric module
# https://www.powershellgallery.com/packages/Az.Fabric
Install-Module -Name Az.Fabric

# Assign the system-assigned managed identity to the VM
# See https://learn.microsoft.com/powershell/module/az.compute/get-azvm for more details.
$vm = Get-AzVM -ResourceGroupName $vmRgName -Name $vmName
# See https://learn.microsoft.com/powershell/module/az.compute/update-azvm for more details.
Update-AzVM -ResourceGroupName $vmRgName -VM $vm -IdentityType SystemAssigned

# Get the system-assigned managed identity details
# See https://learn.microsoft.com/powershell/module/az.compute/get-azvm for more details.
$identityPrincipalId = (Get-AzVM -ResourceGroupName $vmRgName -Name $vmName).Identity.PrincipalId

# Get the Fabric Capacity
# See https://learn.microsoft.com/powershell/module/az.fabric/get-azfabriccapacity for more details.
$fabricCapacity = (Get-AzFabricCapacity -ResourceGroupName $fabricCapacityRgName -CapacityName $fabricCapacityName)

# Assign Contributor role for the system-assigned managed identity to the Fabric Capacity
# See https://learn.microsoft.com/azure/role-based-access-control/role-assignments-powershell for more details.
New-AzRoleAssignment -ObjectId $identityPrincipalId -RoleDefinitionName Contributor -Scope $fabricCapacity.Id

# Get current admin members and add a new principal to the array
$members = $fabricCapacity.AdministrationMember
$members += $identityPrincipalId

# Update the Fabric Capacity with the new admin members
# See https://learn.microsoft.com/powershell/module/az.fabric/update-azfabriccapacity for more details.
Update-AzFabricCapacity -ResourceGroupName $fabricCapacityRgName -CapacityName $fabricCapacityName -AdministrationMember $members
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
#!/bin/bash

# Set input variables
identityRgName="<IDENTITY RESOURCE GROUP NAME>"              # Resource group where the user-assigned managed identity will be created
identityName="<IDENTITY NAME>"                               # Name of the user-assigned managed identity
identityLocation="<IDENTITY LOCATION>"                       # Location where the user-assigned managed identity will be created
vmRgName="<VM RESOURCE GROUP NAME>"                          # Resource group where the VM is located
vmName="<VM NAME>"                                           # Name of the existing VM
fabricCapacityRgName="<FABRIC CAPACITY RESOURCE GROUP NAME>" # Resource group where the Fabric Capacity is located
fabricCapacityName="<FABRIC CAPACITY NAME>"                  # Name of the existing Fabric Capacity

# Install the Microsoft Fabric extension for Azure CLI
# See https://github.com/Azure/azure-cli-extensions/blob/main/src/microsoft-fabric/README.md for more details.
az extension add --name microsoft-fabric

# Create a user-assigned managed identity and get details
# See https://learn.microsoft.com/cli/azure/identity#az-identity-create for more details.
identity=$(az identity create --resource-group "${identityRgName}" --name "${identityName}" --location "${identityLocation}" --output json)
identityId=$(echo "${identity}" | jq -r '.id')
identityPrincipalId=$(echo "${identity}" | jq -r '.principalId')

# Assign the user-assigned managed identity to the VM
# See https://learn.microsoft.com/cli/azure/vm/identity#az-vm-identity-assign for more details.
az vm identity assign --resource-group "${vmRgName}" --name "${vmName}" --identities "${identityId}"

# Get the Fabric Capacity
# See https://learn.microsoft.com/rest/api/microsoftfabric/fabric-capacities/get for more details.
fabricCapacity=$(az fabric capacity show --resource-group "${fabricCapacityRgName}" --capacity-name "${fabricCapacityName}" --output json)
fabricCapacityId=$(echo "${fabricCapacity}" | jq -r '.id')

# Assign Contributor role for the user-assigned managed identity to the Fabric Capacity
# See https://learn.microsoft.com/cli/azure/role/assignment#az-role-assignment-create for more details.
az role assignment create --assignee "${identityPrincipalId}" --role Contributor --scope "${fabricCapacityId}"

# Add a new principal to the the Fabric Capacity admin members
members=$(echo "${fabricCapacity}" | jq -c '.administration.members += ["'"${identityPrincipalId}"'"] | .administration')

# Update the Fabric Capacity with the new admin members
# See https://learn.microsoft.com/rest/api/microsoftfabric/fabric-capacities/update for more details.
az fabric capacity update --resource-group "${fabricCapacityRgName}" --capacity-name "${fabricCapacityName}" --administration "${members}"
```

```powershell
# PowerShell

# See https://learn.microsoft.com/entra/identity/managed-identities-azure-resources/how-to-configure-managed-identities?pivots=qs-configure-powershell-windows-vm#user-assigned-managed-identity for more details.

# Set input variables
$identityRgName = '<IDENTITY RESOURCE GROUP NAME>'              # Resource group where the user-assigned managed identity will be created
$identityName = '<IDENTITY NAME>'                               # Name of the user-assigned managed identity
$identityLocation = '<IDENTITY LOCATION>'                       # Location where the user-assigned managed identity will be created
$vmRgName = '<VM RESOURCE GROUP NAME>'                          # Resource group where the VM is located
$vmName = '<VM NAME>'                                           # Name of the VM
$fabricCapacityRgName = '<FABRIC CAPACITY RESOURCE GROUP NAME>' # Resource group where the Fabric Capacity is located
$fabricCapacityName = '<FABRIC CAPACITY NAME>'                  # Name of the existing Fabric Capacity

# Install the Az.Fabric module
# https://www.powershellgallery.com/packages/Az.Fabric
Install-Module -Name Az.Fabric

# Create a user-assigned managed identity
# See https://learn.microsoft.com/powershell/module/az.managedserviceidentity/new-azuserassignedidentity for more details.
$identity = New-AzUserAssignedIdentity -ResourceGroupName $identityRgName -Name $identityName -Location $identityLocation

# Assign the user-assigned managed identity to the VM
# See https://learn.microsoft.com/powershell/module/az.compute/get-azvm for more details.
$vm = Get-AzVM -ResourceGroupName $vmRgName -Name $vmName
# See https://learn.microsoft.com/powershell/module/az.compute/update-azvm for more details.
Update-AzVM -ResourceGroupName $vmRgName -VM $vm -IdentityType UserAssigned -IdentityID $identity.Id

# Get the Fabric Capacity
# See https://learn.microsoft.com/powershell/module/az.fabric/get-azfabriccapacity for more details.
$fabricCapacity = (Get-AzFabricCapacity -ResourceGroupName $fabricCapacityRgName -CapacityName $fabricCapacityName)

# Assign Contributor role for the user-assigned managed identity to the Fabric Capacity
# See https://learn.microsoft.com/azure/role-based-access-control/role-assignments-powershell for more details.
New-AzRoleAssignment -ObjectId $identity.PrincipalId -RoleDefinitionName Contributor -Scope $fabricCapacity.Id

# Get current admin members and add a new principal to the array
$members = $fabricCapacity.AdministrationMember
$members += $identity.PrincipalId

# Update the Fabric Capacity with the new admin members
# See https://learn.microsoft.com/powershell/module/az.fabric/update-azfabriccapacity for more details.
Update-AzFabricCapacity -ResourceGroupName $fabricCapacityRgName -CapacityName $fabricCapacityName -AdministrationMember $members
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
