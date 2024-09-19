# See https://learn.microsoft.com/cli/azure/ad/app/credential#az-ad-app-credential-reset for more details.
az ad app credential reset --id "00000000-0000-0000-0000-000000000000" --append --cert "@~/client.crt"
