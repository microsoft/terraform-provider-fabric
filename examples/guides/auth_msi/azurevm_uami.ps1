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
