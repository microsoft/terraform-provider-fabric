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
