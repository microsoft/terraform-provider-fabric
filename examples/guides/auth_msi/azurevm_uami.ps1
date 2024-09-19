# See https://learn.microsoft.com/entra/identity/managed-identities-azure-resources/how-to-configure-managed-identities?pivots=qs-configure-powershell-windows-vm#user-assigned-managed-identity for more details.

# Create a user-assigned managed identity
New-AzUserAssignedIdentity -ResourceGroupName "<RESROURCE GROUP NAME>" -Name "<USER ASSIGNED IDENTITY NAME>"

# Assign the user-assigned managed identity to the VM
$vm = Get-AzVM -ResourceGroupName "<RESROURCE GROUP NAME>" -Name "<VM NAME>"
Update-AzVM -ResourceGroupName "<RESROURCE GROUP NAME>" -VM $vm -IdentityType UserAssigned -IdentityID "/subscriptions/<SUBSCRIPTION ID>/resourcegroups/<RESROURCE GROUP NAME>/providers/Microsoft.ManagedIdentity/userAssignedIdentities/<USER ASSIGNED IDENTITY NAME>"

# Assign Contributor role for the user-assigned managed identity to the Fabric Capacity
# See https://learn.microsoft.com/azure/role-based-access-control/role-assignments-powershell for more details.
New-AzRoleAssignment -ObjectId "<PRINCIPAL ID>" -RoleDefinitionName Contributor -Scope "/subscriptions/<SUBSCRIPTION ID>/resourceGroups/<RESROURCE GROUP NAME>/providers/Microsoft.Fabric/capacities/<FABRIC CAPACITY NAME>"
