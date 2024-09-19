# PowerShell

# Login to Entra ID
# See https://learn.microsoft.com/powershell/module/microsoft.graph.entra/connect-entra
Connect-Entra -Scopes 'Application.ReadWrite.All', 'User.Read'

# Get the current context
# See https://learn.microsoft.com/powershell/module/microsoft.graph.entra/get-entracontext
$ctx = Get-EntraContext

# Get the current user
# See https://learn.microsoft.com/powershell/module/microsoft.graph.entra/get-entrauser
$currentUser = (Get-EntraUser -Filter "UserPrincipalName eq '$($ctx.Account)'" -Property Id)

# Create a new Entra Application
# See https://learn.microsoft.com/powershell/module/microsoft.graph.entra/new-entraapplication for more details.
$app = (New-EntraApplication -DisplayName 'Fabric Terraform Provider' -SigninAudience AzureADMyOrg)

# (optional, recommended) Add the current user as an owner of the application
# See https://learn.microsoft.com/powershell/module/microsoft.graph.entra/add-entraapplicationowner for more details.
Add-EntraApplicationOwner -ObjectId $app.Id -RefObjectId $currentUser.Id

# Create a new Entra Service Principal associated with the application
# See https://learn.microsoft.com/powershell/module/microsoft.graph.entra/new-entraserviceprincipal for more details.
$sp = (New-EntraServicePrincipal -AppId $app.AppId)

# (optional, recommended) Add the current user as an owner of the service principal
# See https://learn.microsoft.com/powershell/module/microsoft.graph.entra/add-entraserviceprincipalowner for more details.
Add-EntraServicePrincipalOwner -ObjectId $sp.Id -RefObjectId $currentUser.Id
