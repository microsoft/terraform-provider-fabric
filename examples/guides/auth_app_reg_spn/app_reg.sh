# sh

# Login to Azure with Entra ID credentials
# See https://learn.microsoft.com/cli/azure/authenticate-azure-cli for more details.
az login --allow-no-subscriptions

# Get the current user ID
# See https://learn.microsoft.com/cli/azure/ad/signed-in-user#az-ad-signed-in-user-show for more details.
currentUserObjId=$(az ad signed-in-user show --output tsv --query id)

# Create a new Entra Application
# See https://learn.microsoft.com/cli/azure/ad/app#az-ad-app-create) for more details.
appObjId=$(az ad app create --display-name "Fabric Terraform Provider" --sign-in-audience AzureADMyOrg --output tsv --query id)

# (optional, recommended) Add the current user as an owner of the application
# See https://learn.microsoft.com/cli/azure/ad/app/owner#az-ad-app-owner-add for more details.
az ad app owner add --id "${appObjId}" --owner-object-id "${currentUserObjId}"

# Create a new Entra Service Principal associated with the application
# see https://learn.microsoft.com/cli/azure/ad/sp#az-ad-sp-create for more details.
spObjId=$(az ad sp create --id "${appObjId}" --output tsv --query id)

# (optional, recommended) Add the current user as an owner of the service principal
# See https://learn.microsoft.com/cli/azure/reference-index#az-rest for more details.
az rest --method POST --url "https://graph.microsoft.com/v1.0/servicePrincipals/${spObjId}/owners/\$ref" --body "{\"@odata.id\": \"https://graph.microsoft.com/v1.0/users/${currentUserObjId}\"}"
