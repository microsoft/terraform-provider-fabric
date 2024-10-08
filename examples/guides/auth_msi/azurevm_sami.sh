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
