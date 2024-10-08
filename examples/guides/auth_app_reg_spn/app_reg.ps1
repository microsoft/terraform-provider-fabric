# PowerShell

# Set input variables
$fabricCapacityRgName = '<FABRIC CAPACITY RESOURCE GROUP NAME>' # Resource group where the Fabric Capacity is located
$fabricCapacityName = '<FABRIC CAPACITY NAME>'                  # Name of the existing Fabric Capacity

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
$app = (New-EntraApplication -DisplayName 'Microsoft Fabric Terraform Provider' -SigninAudience AzureADMyOrg)

# (optional, recommended) Add the current user as an owner of the application
# See https://learn.microsoft.com/powershell/module/microsoft.graph.entra/add-entraapplicationowner for more details.
Add-EntraApplicationOwner -ObjectId $app.Id -RefObjectId $currentUser.Id

# Create a new Entra Service Principal associated with the application
# See https://learn.microsoft.com/powershell/module/microsoft.graph.entra/new-entraserviceprincipal for more details.
$sp = (New-EntraServicePrincipal -AppId $app.AppId)

# (optional, recommended) Add the current user as an owner of the service principal
# See https://learn.microsoft.com/powershell/module/microsoft.graph.entra/add-entraserviceprincipalowner for more details.
Add-EntraServicePrincipalOwner -ObjectId $sp.Id -RefObjectId $currentUser.Id

# Install the Az.Fabric module
# https://www.powershellgallery.com/packages/Az.Fabric
Install-Module -Name Az.Fabric

# Get current admin members and add a new principal to the array
# See https://learn.microsoft.com/powershell/module/az.fabric/get-azfabriccapacity for more details.
$members = (Get-AzFabricCapacity -ResourceGroupName $fabricCapacityRgName -CapacityName $fabricCapacityName).AdministrationMember
$members += $sp.Id

# Update the Fabric Capacity with the new admin members
# See https://learn.microsoft.com/powershell/module/az.fabric/update-azfabriccapacity for more details.
Update-AzFabricCapacity -ResourceGroupName $fabricCapacityRgName -CapacityName $fabricCapacityName -AdministrationMember $members
