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
