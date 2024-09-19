# (optional, useful in the multi-tenant scenarios) Disable the new login experience
# See https://learn.microsoft.com/powershell/azure/authenticate-interactive#disable-the-new-login-experience for more details.
Update-AzConfig -LoginExperienceV2 Off

# (optional, Windows only, useful in the multi-tenant scenarios) Disable WAM on Windows
# See https://learn.microsoft.com/powershell/azure/authenticate-interactive#web-account-manager-wam for more details.
Update-AzConfig -EnableLoginByWam $false

# Login to Azure with Entra ID credentials
# See https://learn.microsoft.com/powershell/module/az.accounts/connect-azaccount for more details.
Connect-AzAccount -Tenant '00000000-0000-0000-0000-000000000000' -AuthScope 'api://fabric_terraform_provider'

# Set the FABRIC_TOKEN environment variable to the access token
# See https://learn.microsoft.com/powershell/module/az.accounts/get-azaccesstoken for more details.
$env:FABRIC_TOKEN = (Get-AzAccessToken -ResourceUrl 'https://api.fabric.microsoft.com').Token
