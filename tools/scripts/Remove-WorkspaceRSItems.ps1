#
# Remove-WorkspaceRSItems.ps1
#
# This script deletes all items in the "WorkspaceRS" workspace except for:
# 1. The "LakehouseRS" item (preserved by ID)
#
# It reads the workspace and lakehouse IDs from the .wellknown.json file.
#
# Usage:
#   .\Remove-WorkspaceRSItems.ps1 [-WellKnownPath <path>] [-Force] [-DryRun]
#
# Parameters:
#   -WellKnownPath: Path to the .wellknown.json file (defaults to ../../internal/testhelp/fixtures/.wellknown.json)
#   -Force: Skip confirmation prompt and proceed with deletion automatically
#   -DryRun: Show what would be deleted without actually deleting anything
#
# Examples:
#   .\Remove-WorkspaceRSItems.ps1
#   .\Remove-WorkspaceRSItems.ps1 -WellKnownPath "C:\path\to\.wellknown.json"
#   .\Remove-WorkspaceRSItems.ps1 -Force
#   .\Remove-WorkspaceRSItems.ps1 -DryRun
#   .\Remove-WorkspaceRSItems.ps1 -DryRun -WellKnownPath "custom-path.json"
#

param(
  [Parameter(Mandatory = $false)]
  [string]$WellKnownPath,

  [Parameter(Mandatory = $false)]
  [switch]$Force,

  [Parameter(Mandatory = $false)]
  [switch]$DryRun
)

function Write-Log {
  param (
    [Parameter(Mandatory = $true)]
    [string]$Message,

    [Parameter(Mandatory = $false)]
    [string]$Level = 'INFO',

    [Parameter(Mandatory = $false)]
    [bool]$Stop = $true
  )

  $color = switch ($Level) {
    'INFO' { 'Green' }
    'WARN' { 'Yellow' }
    'ERROR' { 'Red' }
    'DEBUG' { 'DarkMagenta' }
    default { 'Green' }
  }

  $prefix = switch ($Level) {
    'INFO' { '*' }
    'WARN' { '!' }
    'ERROR' { 'X' }
    'DEBUG' { 'D' }
    default { '*' }
  }

  Write-Host -ForegroundColor $color "[$prefix] $Message"

  if ($Stop -and $Level -eq 'ERROR') {
    exit 1
  }
}

function Install-ModuleIfNotInstalled {
  param (
    [Parameter(Mandatory = $true)]
    [string]$ModuleName
  )

  if (-not (Get-Module -Name $ModuleName -ListAvailable)) {
    try {
      Write-Log -Message "Installing module: $ModuleName" -Level 'DEBUG'
      Install-Module -Name $ModuleName -AllowClobber -Force -Scope CurrentUser -Repository PSGallery -Confirm:$false -SkipPublisherCheck -AcceptLicense
    }
    catch {
      Write-Error $_.Exception.Message
      Write-Log -Message "Unable to install module: $ModuleName" -Level 'ERROR'
    }
  }
}

function Import-ModuleIfNotImported {
  param (
    [Parameter(Mandatory = $true)]
    [string]$ModuleName
  )

  if (-not (Get-Module -Name $ModuleName)) {
    try {
      Write-Log -Message "Importing module: $ModuleName" -Level 'DEBUG'
      Import-Module -Name $ModuleName
    }
    catch {
      Write-Error $_.Exception.Message
      Write-Log -Message "Unable to import module: $ModuleName" -Level 'ERROR'
    }
  }
}


function Invoke-FabricRest {
  param (
    [Parameter(Mandatory = $false)]
    [string]$Method = 'GET',

    [Parameter(Mandatory = $true)]
    [string]$Endpoint,

    [Parameter(Mandatory = $false)]
    [object]$Payload,

    [Parameter(Mandatory = $false)]
    [int]$RetryCount = 3,

    [Parameter(Mandatory = $false)]
    [int]$RetryDelaySeconds = 30
  )

  try {
    # Retrieve the Fabric access token
    try {
      $secureAccessToken = (Get-AzAccessToken -WarningAction SilentlyContinue -AsSecureString -ResourceUrl 'https://api.fabric.microsoft.com').Token
    }
    catch {
      Write-Log -Message "Failed to retrieve access token." -Level 'ERROR'
    }

    $uri = "https://api.fabric.microsoft.com/v1/$Endpoint"
    $attempt = 0
    $response = $null
    $responseHeaders = $null
    $statusCode = $null

    while ($attempt -lt $RetryCount) {
      try {
        if ($Payload) {
          $body = $Payload | ConvertTo-Json -Depth 10 -Compress
          $response = Invoke-RestMethod -Authentication Bearer -Token $secureAccessToken -Uri $uri -Method $Method -ContentType 'application/json' -Body $body -ResponseHeadersVariable responseHeaders -StatusCodeVariable statusCode
        }
        else {
          $response = Invoke-RestMethod -Authentication Bearer -Token $secureAccessToken -Uri $uri -Method $Method -ResponseHeadersVariable responseHeaders -StatusCodeVariable statusCode
        }

        break
      }
      catch {
        $statusCode = $_.Exception.Response.StatusCode.value__

        if ($statusCode -eq 429) {
          $retryAfter = $_.Exception.Response.Headers.RetryAfter.Delta.TotalSeconds
          if ($retryAfter -gt 0) {
            Write-Log -Message "Rate limited. Retrying after $retryAfter seconds..." -Level 'WARN' -Stop $false
            Start-Sleep -Seconds $retryAfter
          }
          else {
            Write-Log -Message "Rate limited. Retrying after $RetryDelaySeconds seconds..." -Level 'WARN' -Stop $false
            Start-Sleep -Seconds $RetryDelaySeconds
          }
        }
        else {
          Write-Log -Message "Failed to invoke Fabric REST API: $($_.Exception.Message)" -Level 'ERROR' -Stop $false
          Write-Log -Message "Response: $($_.Exception.Response | ConvertTo-Json -Depth 10)" -Level 'DEBUG' -Stop $false
          throw $_
        }

        $attempt++
      }
    }

    if ($attempt -ge $RetryCount) {
      Write-Log -Message "Failed to invoke Fabric REST API after $RetryCount attempts." -Level 'ERROR'
    }

    return @{
      Response   = $response
      Headers    = $responseHeaders
      StatusCode = $statusCode
    }
  }
  catch {
    Write-Log -Message "Failed to invoke Fabric REST API: $($_.Exception.Message)" -Level 'ERROR'
  }
}

function Remove-WorkspaceRSItems {
  param (
    [Parameter(Mandatory = $true)]
    [string]$WellKnownJsonPath,

    [Parameter(Mandatory = $false)]
    [switch]$Force,

    [Parameter(Mandatory = $false)]
    [switch]$DryRun
  )

  # Check if the .wellknown.json file exists
  if (-not (Test-Path -Path $WellKnownJsonPath)) {
    Write-Log -Message "The .wellknown.json file was not found at path: $WellKnownJsonPath" -Level 'ERROR'
    return
  }

  # Read and parse the .wellknown.json file
  try {
    $wellKnownContent = Get-Content -Path $WellKnownJsonPath -Raw | ConvertFrom-Json
  }
  catch {
    Write-Log -Message "Failed to parse .wellknown.json file: $($_.Exception.Message)" -Level 'ERROR'
    return
  }

  # Extract WorkspaceRS and LakehouseRS IDs
  if (-not $wellKnownContent.WorkspaceRS -or -not $wellKnownContent.WorkspaceRS.id) {
    Write-Log -Message "WorkspaceRS ID not found in .wellknown.json" -Level 'ERROR'
    return
  }

  if (-not $wellKnownContent.LakehouseRS -or -not $wellKnownContent.LakehouseRS.id) {
    Write-Log -Message "LakehouseRS ID not found in .wellknown.json" -Level 'ERROR'
    return
  }

  $workspaceId = $wellKnownContent.WorkspaceRS.id
  $lakehouseRSId = $wellKnownContent.LakehouseRS.id

  Write-Log -Message "WorkspaceRS ID: $workspaceId" -Level 'INFO'
  Write-Log -Message "LakehouseRS ID to preserve: $lakehouseRSId" -Level 'INFO'

  # Get all items in the WorkspaceRS
  Write-Log -Message "Fetching all items in WorkspaceRS..." -Level 'INFO'
  try {
    $itemsResponse = Invoke-FabricRest -Method 'GET' -Endpoint "workspaces/$workspaceId/items"
    $items = $itemsResponse.Response.value
  }
  catch {
    Write-Log -Message "Failed to fetch items from WorkspaceRS: $($_.Exception.Message)" -Level 'ERROR'
    return
  }

  if (-not $items -or $items.Count -eq 0) {
    Write-Log -Message "No items found in WorkspaceRS" -Level 'INFO'
    return
  }

  Write-Log -Message "Found $($items.Count) items in WorkspaceRS" -Level 'INFO'

  # Filter out items to preserve:
  # 1. The LakehouseRS item by ID (exact match)
  $itemsToDelete = $items | Where-Object {
    $_.id -ne $lakehouseRSId -and
    $_.type -ne 'SQLEndpoint' -and
    $_.type -ne 'KQLDatabase'
  }

  Write-Log -Message "Preserving the LakehouseRS ID: $($lakehouseRSId)" -Level 'INFO' -Stop $false

  if ($itemsToDelete.Count -eq 0) {
    Write-Log -Message "No items to delete (only preserved items found)" -Level 'INFO'
    return
  }

  Write-Log -Message "Items to delete: $($itemsToDelete.Count)" -Level 'INFO'

  if ($DryRun) {
    Write-Log -Message "DRY RUN MODE - No items will actually be deleted" -Level 'WARN' -Stop $false
    Write-Log -Message "Items that would be deleted:" -Level 'INFO' -Stop $false
    foreach ($item in $itemsToDelete) {
      Write-Log -Message "  - $($item.displayName) (Type: $($item.type), ID: $($item.id))" -Level 'INFO' -Stop $false
    }
    Write-Log -Message "=== DRY RUN SUMMARY ===" -Level 'INFO'
    Write-Log -Message "Total items that would be deleted: $($itemsToDelete.Count)" -Level 'INFO'
    Write-Log -Message "DRY RUN completed - no actual deletions performed" -Level 'INFO'
    return
  }

  # Confirm deletion with user (unless Force is specified)
  if (-not $Force) {
    Write-Log -Message "About to delete $($itemsToDelete.Count) items from WorkspaceRS:" -Level 'WARN' -Stop $false
    foreach ($item in $itemsToDelete) {
      Write-Log -Message "  - $($item.displayName) (Type: $($item.type), ID: $($item.id))" -Level 'INFO' -Stop $false
    }

    $confirmation = Read-Host "Do you want to proceed with deletion? (y/N)"
    if ($confirmation -ne 'y' -and $confirmation -ne 'Y') {
      Write-Log -Message "Deletion cancelled by user" -Level 'INFO'
      return
    }
  }
  else {
    Write-Log -Message "Force mode enabled - proceeding with deletion of $($itemsToDelete.Count) items" -Level 'WARN' -Stop $false
  }

  # Delete items one by one
  $totalToDelete = $itemsToDelete.Count
  $deletedCount = 0
  $failedCount = 0

  Write-Log -Message "Starting deletion of $totalToDelete items..." -Level 'INFO'

  foreach ($item in $itemsToDelete) {
    Write-Log -Message "Deleting item: $($item.displayName) (Type: $($item.type), ID: $($item.id))" -Level 'INFO'

    try {
      $deleteResponse = Invoke-FabricRest -Method 'DELETE' -Endpoint "workspaces/$workspaceId/items/$($item.id)"

      if ($deleteResponse.StatusCode -eq 200 -or $deleteResponse.StatusCode -eq 204) {
        Write-Log -Message "Successfully deleted item: $($item.displayName)" -Level 'INFO'
        $deletedCount++
      }
      else {
        Write-Log -Message "Failed to delete item: $($item.displayName) - Status Code: $($deleteResponse.StatusCode)" -Level 'ERROR' -Stop $false
        $failedCount++
      }
    }
    catch {
      Write-Log -Message "Error deleting item: $($item.displayName) - $($_.Exception.Message)" -Level 'ERROR' -Stop $false
      $failedCount++
    }

    # Add a small delay between deletions to avoid rate limiting
    Start-Sleep -Milliseconds 500
  }

  # Summary
  Write-Log -Message "=== DELETION SUMMARY ===" -Level 'INFO'
  Write-Log -Message "Total items to delete: $totalToDelete" -Level 'INFO'
  Write-Log -Message "Successfully deleted: $deletedCount" -Level 'INFO'
  Write-Log -Message "Failed to delete: $failedCount" -Level 'INFO'

  if ($failedCount -eq 0) {
    Write-Log -Message "All items deleted successfully!" -Level 'INFO'
  }
  elseif ($deletedCount -eq 0) {
    Write-Log -Message "No items were deleted due to errors." -Level 'ERROR' -Stop $false
  }
  else {
    Write-Log -Message "Deletion completed with some failures." -Level 'WARN' -Stop $false
  }
}

# Main execution
# Define an array of modules to install
$modules = @('Az.Accounts', 'Az.Fabric', 'pwsh-dotenv')

# Loop through each module and install if not installed
foreach ($module in $modules) {
  Install-ModuleIfNotInstalled -ModuleName $module
  Import-ModuleIfNotImported -ModuleName $module
}


# Import the .env file into the environment variables
if (Test-Path -Path './wellknown.env') {
  Import-Dotenv -Path ./wellknown.env -AllowClobber
}

if (
  !$Env:FABRIC_TESTACC_WELLKNOWN_ENTRA_TENANT_ID -or
  !$Env:FABRIC_TESTACC_WELLKNOWN_AZURE_SUBSCRIPTION_ID
) {
  Write-Log -Message @'
  FABRIC_TESTACC_WELLKNOWN_ENTRA_TENANT_ID,
  FABRIC_TESTACC_WELLKNOWN_AZURE_SUBSCRIPTION_ID,
  are required environment variables.
'@ `
    -Level 'ERROR'
  return
}


# Check if already logged in to Azure, if not then login
$azContext = Get-AzContext
if (!$azContext -or $azContext.Tenant.Id -ne $Env:FABRIC_TESTACC_WELLKNOWN_ENTRA_TENANT_ID -or $azContext.Subscription.Id -ne $Env:FABRIC_TESTACC_WELLKNOWN_AZURE_SUBSCRIPTION_ID) {
  Write-Log -Message 'Logging in to Azure.' -Level 'DEBUG'
  Connect-AzAccount -Tenant $Env:FABRIC_TESTACC_WELLKNOWN_ENTRA_TENANT_ID -SubscriptionId $Env:FABRIC_TESTACC_WELLKNOWN_AZURE_SUBSCRIPTION_ID -UseDeviceAuthentication
  $azContext = Get-AzContext
  Write-Log -Message 'Logged in successfully' -Level 'DEBUG'
  # Disconnect-AzAccount
}

if (-not $WellKnownPath) {
  $WellKnownPath = Join-Path $PSScriptRoot "../../internal/testhelp/fixtures/.wellknown.json"
}

Write-Log -Message "Starting WorkspaceRS cleanup process..." -Level 'INFO'
Write-Log -Message "Using .wellknown.json file: $WellKnownPath" -Level 'INFO'

Remove-WorkspaceRSItems -WellKnownJsonPath $WellKnownPath -Force:$Force -DryRun:$DryRun

Write-Log -Message "WorkspaceRS cleanup process completed" -Level 'INFO'
