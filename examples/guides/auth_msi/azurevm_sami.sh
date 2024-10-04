#!/bin/bash

# Set input variables
vmRgName="<VM RESOURCE GROUP NAME>"                          # Resource group where the VM is located
vmName="<VM NAME>"                                           # Name of the VM
fabricCapacityRgName="<FABRIC CAPACITY RESOURCE GROUP NAME>" # Resource group where the Fabric Capacity is located
fabricCapacityName="<FABRIC CAPACITY NAME>"                  # Name of the existing Fabric Capacity

# Login to Azure with Entra ID credentials
# See https://learn.microsoft.com/cli/azure/authenticate-azure-cli for more details.
az login

# Get the current subscription ID
# See https://learn.microsoft.com/cli/azure/account#az-account-show for more details.
subscriptionId=$(az account show --output tsv --query id)

# Assign the system-assigned managed identity to the VM
# See https://learn.microsoft.com/cli/azure/vm/identity#az-vm-identity-assign for more details.
identityPrincipalId=$(az vm identity assign --resource-group "${vmRgName}" --name "${vmName}" --identities "[system]" --output tsv --query systemAssignedIdentity)

# Assign Contributor role for the system-assigned managed identity to the Fabric Capacity
# See https://learn.microsoft.com/cli/azure/role/assignment#az-role-assignment-create for more details.
az role assignment create --assignee "${identityPrincipalId}" --role Contributor --scope "/subscriptions/${subscriptionId}/resourceGroups/${fabricCapacityRgName}/providers/Microsoft.Fabric/capacities/${fabricCapacityName}"

# Install the Microsoft Fabric extension for Azure CLI
# See https://github.com/Azure/azure-cli-extensions/blob/main/src/microsoft-fabric/README.md for more details.
az extension add --name microsoft-fabric

# Get current Fabric Capacity admin members and add a new principal to the array
# See https://learn.microsoft.com/rest/api/microsoftfabric/fabric-capacities/get for more details.
members=$(az fabric capacity show --resource-group "${fabricCapacityRgName}" --capacity-name "${fabricCapacityName}" --output json --query administration | jq --compact-output '.members += ["'"${identityPrincipalId}"'"]')

# Update the Fabric Capacity with the new admin members
# See https://learn.microsoft.com/rest/api/microsoftfabric/fabric-capacities/update for more details.
az fabric capacity update --resource-group "${fabricCapacityRgName}" --capacity-name "${fabricCapacityName}" --administration "${members}"
