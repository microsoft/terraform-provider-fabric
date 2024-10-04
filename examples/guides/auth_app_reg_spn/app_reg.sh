#!/bin/bash

# Set input variables
fabricCapacityRgName="<FABRIC CAPACITY RESOURCE GROUP NAME>" # Resource group where the Fabric Capacity is located
fabricCapacityName="<FABRIC CAPACITY NAME>"                  # Name of the existing Fabric Capacity

# Login to Azure with Entra ID credentials
# See https://learn.microsoft.com/cli/azure/authenticate-azure-cli for more details.
az login

# Get the current subscription ID
# See https://learn.microsoft.com/cli/azure/account#az-account-show for more details.
subscriptionId=$(az account show --output tsv --query id)

# Get the current user ID
# See https://learn.microsoft.com/cli/azure/ad/signed-in-user#az-ad-signed-in-user-show for more details.
currentUserObjId=$(az ad signed-in-user show --output tsv --query id)

# Create a new Entra Application
# See https://learn.microsoft.com/cli/azure/ad/app#az-ad-app-create) for more details.
appObjId=$(az ad app create --display-name "Microsoft Fabric Terraform Provider" --sign-in-audience AzureADMyOrg --output tsv --query id)

# (optional, recommended) Add the current user as an owner of the application
# See https://learn.microsoft.com/cli/azure/ad/app/owner#az-ad-app-owner-add for more details.
az ad app owner add --id "${appObjId}" --owner-object-id "${currentUserObjId}"

# Create a new Entra Service Principal associated with the application
# see https://learn.microsoft.com/cli/azure/ad/sp#az-ad-sp-create for more details.
spObjId=$(az ad sp create --id "${appObjId}" --output tsv --query id)

# (optional, recommended) Add the current user as an owner of the service principal
# See https://learn.microsoft.com/cli/azure/reference-index#az-rest for more details.
az rest --method POST --url "https://graph.microsoft.com/v1.0/servicePrincipals/${spObjId}/owners/\$ref" --body "{\"@odata.id\": \"https://graph.microsoft.com/v1.0/users/${currentUserObjId}\"}"

# Assign Contributor role for the user-assigned managed identity to the Fabric Capacity
# See https://learn.microsoft.com/cli/azure/role/assignment#az-role-assignment-create for more details.
az role assignment create --assignee "${spObjId}" --role Contributor --scope "/subscriptions/${subscriptionId}/resourceGroups/${fabricCapacityRgName}/providers/Microsoft.Fabric/capacities/${fabricCapacityName}"

# Install the Microsoft Fabric extension for Azure CLI
# See https://github.com/Azure/azure-cli-extensions/blob/main/src/microsoft-fabric/README.md for more details.
az extension add --name microsoft-fabric

# Get current admin members and add a new principal to the array
# See https://learn.microsoft.com/rest/api/microsoftfabric/fabric-capacities/get for more details.
members=$(az fabric capacity show --resource-group "${fabricCapacityRgName}" --capacity-name "${fabricCapacityName}" --output json --query administration | jq --compact-output '.members += ["'"${spObjId}"'"]')

# Update the Fabric Capacity with the new admin members
# See https://learn.microsoft.com/rest/api/microsoftfabric/fabric-capacities/update for more details.
az fabric capacity update --resource-group "${fabricCapacityRgName}" --capacity-name "${fabricCapacityName}" --administration "${members}"
