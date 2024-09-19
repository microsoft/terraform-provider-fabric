# See https://learn.microsoft.com/entra/identity/managed-identities-azure-resources/how-to-configure-managed-identities?pivots=qs-configure-powershell-windows-vm#system-assigned-managed-identity for more details.

# Assign the system-assigned managed identity to the VM
$vm = Get-AzVM -ResourceGroupName "<RESOURCE GROUP NAME>" -Name "<VM NAME>"
Update-AzVM -ResourceGroupName "<RESOURCE GROUP NAME>" -VM $vm -IdentityType SystemAssigned

# Assign Contributor role for the system-assigned managed identity to the Fabric Capacity
# See https://learn.microsoft.com/azure/role-based-access-control/role-assignments-powershell for more details.
New-AzRoleAssignment -ObjectId "<PRINCIPAL ID>" -RoleDefinitionName Contributor -Scope "/subscriptions/<SUBSCRIPTION ID>/resourceGroups/<RESROURCE GROUP NAME>/providers/Microsoft.Fabric/capacities/<FABRIC CAPACITY NAME>"
