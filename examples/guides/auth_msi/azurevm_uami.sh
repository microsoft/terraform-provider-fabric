# Create a user-assigned managed identity
# See https://learn.microsoft.com/cli/azure/identity#az-identity-create for more details.
az identity create --resource-group "<RESOURCE GROUP NAME>" --name "<IDENTITY NAME>"

# Assign the user-assigned managed identity to the VM
# See https://learn.microsoft.com/cli/azure/vm/identity#az-vm-identity-assign for more details.
az vm identity assign --resource-group "<RESOURCE GROUP NAME>" --name "<VM NAME>" --identities "<IDENTITY ID>"

# Assign Contributor role for the user-assigned managed identity to the Fabric Capacity
# See https://learn.microsoft.com/cli/azure/role/assignment#az-role-assignment-create for more details.
az role assignment create --assignee "<PRINCIPAL ID>" --role Contributor --scope "/subscriptions/<SUBSCRIPTION ID>/resourceGroups/<RESOURCE GROUP NAME>/providers/Microsoft.Fabric/capacities/<FABRIC CAPACITY NAME>"
