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

          $retryDelaySeconds = $RetryDelaySeconds
          if ($retryAfter) {
            $retryDelaySeconds = $retryAfter
          }

          Write-Log -Message "Throttled. Waiting for $retryDelaySeconds seconds before retrying..." -Level 'DEBUG'
          Start-Sleep -Seconds $retryDelaySeconds

          $attempt++
        }
        else {
          throw $_
        }
      }
    }

    if ($attempt -ge $RetryCount) {
      Write-Log -Message "Maximum retry attempts reached. Request failed." -Level 'ERROR'
    }

    if ($statusCode -eq 200 -or $statusCode -eq 201) {
      return [PSCustomObject]@{
        Response = $response
        Headers  = $responseHeaders
      }
    }

    if ($statusCode -eq 202 -and $responseHeaders.Location -and $responseHeaders['x-ms-operation-id']) {
      $operationId = [string]$responseHeaders['x-ms-operation-id']
      Write-Log -Message "Long Running Operation initiated. Operation ID: $operationId" -Level 'DEBUG'
      $result = Get-LroResult -OperationId $operationId

      return [PSCustomObject]@{
        Response = $result.Response
        Headers  = $result.Headers
      }
    }
  }
  catch {
    Write-Log -Message $_.Exception.Message -Level 'ERROR'
  }
}

function Get-LroResult {
  param (
    [Parameter(Mandatory = $true)]
    [string]$OperationId
  )

  $operationStatus = $null
  while ($operationStatus -ne 'Succeeded') {
    $result = Invoke-FabricRest -Method 'GET' -Endpoint "operations/$OperationId"

    $operationStatus = $result.Response.status

    if ($operationStatus -eq 'Failed') {
      Write-Log -Message "Operation failed. Status: $operationStatus" -Level 'ERROR'
    }

    if ($operationStatus -ne "Succeeded") {
      $retryAfter = [int]$result.Headers['Retry-After'][0]
      Start-Sleep -Seconds $retryAfter
    }
  }

  return Invoke-FabricRest -Method 'GET' -Endpoint "operations/$OperationId/result"
}

function Set-FabricItem {
  param (
    [Parameter(Mandatory = $true)]
    [string]$DisplayName,

    [Parameter(Mandatory = $true)]
    [string]$WorkspaceId,

    [Parameter(Mandatory = $true)]
    [string]$Type,

    [Parameter(Mandatory = $false)]
    [object]$CreationPayload,

    [Parameter(Mandatory = $false)]
    [object]$Definition
  )

  switch ($Type) {
    'ApacheAirflowJob' {
      $itemEndpoint = 'apacheAirflowJobs'
    }
    'CopyJob' {
      $itemEndpoint = 'copyJobs'
    }
    'Dataflow' {
      $itemEndpoint = 'dataflows'
    }
    'DataPipeline' {
      $itemEndpoint = 'dataPipelines'
    }
    'DigitalTwinBuilder' {
      $itemEndpoint = 'digitalTwinBuilders'
    }
    'Environment' {
      $itemEndpoint = 'environments'
    }
    'Eventhouse' {
      $itemEndpoint = 'eventhouses'
    }
    'Eventstream' {
      $itemEndpoint = 'eventstreams'
    }
    'GraphQLApi' {
      $itemEndpoint = 'GraphQLApis'
    }
    'KQLDashboard' {
      $itemEndpoint = 'kqlDashboards'
    }
    'KQLDatabase' {
      $itemEndpoint = 'kqlDatabases'
    }
    'KQLQueryset' {
      $itemEndpoint = 'kqlQuerysets'
    }
    'Lakehouse' {
      $itemEndpoint = 'lakehouses'
    }
    'MirroredDatabase' {
      $itemEndpoint = 'mirroredDatabases'
    }
    'MLExperiment' {
      $itemEndpoint = 'mlExperiments'
    }
    'MLModel' {
      $itemEndpoint = 'mlModels'
    }
    'MountedDataFactory' {
      $itemEndpoint = 'mountedDataFactories'
    }
    'Notebook' {
      $itemEndpoint = 'notebooks'
    }
    'Reflex' {
      $itemEndpoint = 'reflexes'
    }
    'Report' {
      $itemEndpoint = 'reports'
    }
    'SemanticModel' {
      $itemEndpoint = 'semanticModels'
    }
    'SparkJobDefinition' {
      $itemEndpoint = 'sparkJobDefinitions'
    }
    'SQLDatabase' {
      $itemEndpoint = 'sqlDatabases'
    }
    'Warehouse' {
      $itemEndpoint = 'warehouses'
    }
    default {
      $itemEndpoint = 'items'
    }
  }

  If ($CreationPayload -and $Definition) {
    Write-Log -Message 'Only one of CreationPayload or Definition is allowed at time.' -Level 'ERROR'
  }

  $definitionRequired = @('ApacheAirflowJob', 'Report', 'SemanticModel', 'MirroredDatabase', 'MountedDataFactory', 'Eventstream')
  if ($Type -in $definitionRequired -and !$Definition) {
    Write-Log -Message "Definition is required for Type: $Type" -Level 'ERROR'
  }

  $results = Invoke-FabricRest -Method 'GET' -Endpoint "workspaces/$WorkspaceId/$itemEndpoint"
  $result = $results.Response.value | Where-Object { $_.displayName -eq $DisplayName }
  if (!$result) {
    Write-Log -Message "Creating ${Type}: $DisplayName" -Level 'WARN'
    $payload = @{
      displayName = $DisplayName
      description = $DisplayName
    }

    if ($itemEndpoint -eq 'items') {
      $payload['type'] = $Type
    }

    if ($CreationPayload) {
      $payload['creationPayload'] = $CreationPayload
    }

    if ($Definition) {
      $payload['definition'] = $Definition
    }

    $resultPost = (Invoke-FabricRest -Method 'POST' -Endpoint "workspaces/$WorkspaceId/$itemEndpoint" -Payload $payload)
    $result = $resultPost.Response
  }

  Write-Log -Message "${Type} - Name: $($result.displayName) / ID: $($result.id)"

  return $result
}

function Get-DefinitionPartBase64 {
  param (
    [Parameter(Mandatory = $true)]
    [string]$Path,

    [Parameter(Mandatory = $false)]
    [object]$Values
  )

  if (-not (Test-Path -Path $Path)) {
    Write-Log -Message "File not found: $Path" -Level 'ERROR'
  }

  $content = (Get-Content -Path $Path -Raw).Trim().ToString()

  if ($Values) {
    foreach ($value in $Values) {
      $content = $content.Replace($value.key, $value.value)
    }
  }

  return [Convert]::ToBase64String([System.Text.Encoding]::UTF8.GetBytes($content))
}

function Set-DeploymentPipeline {
  param(
    [Parameter(Mandatory = $true)]
    [string]$DisplayName
  )
  $results = Invoke-FabricRest -Method 'GET' -Endpoint "deploymentPipelines"
  $result = $results.Response.value | Where-Object { $_.displayName -eq $DisplayName }
  if (!$result) {
    Write-Log -Message "Creating Deployment Pipeline: $DisplayName" -Level 'WARN'
    $payload = @{
      displayName = $DisplayName
      description = $DisplayName
      stages      = @(
        @{
          displayName = "Development"
          description = "Development stage description"
          isPublic    = $false
        },
        @{
          displayName = "Test"
          description = "Test stage description"
          isPublic    = $false
        },
        @{
          displayName = "Production"
          description = "Production stage description"
          isPublic    = $true
        }
      )
    }
    $result = (Invoke-FabricRest -Method 'POST' -Endpoint "deploymentPipelines" -Payload $payload).Response
  }
  else {
    $result = Invoke-FabricRest -Method 'GET' -Endpoint "deploymentPipelines/$($result.id)"
    $result = $result.Response
  }
  Write-Log -Message "Deployment Pipeline - Name: $($result.displayName) / ID: $($result.id)"

  return $result
}

function Set-FabricDomain {
  param (
    [Parameter(Mandatory = $true)]
    [string]$DisplayName,

    [Parameter(Mandatory = $false)]
    [string]$ParentDomainId
  )

  $results = Invoke-FabricRest -Method 'GET' -Endpoint "admin/domains"
  $result = $results.Response.domains | Where-Object { $_.displayName -eq $DisplayName }
  if (!$result) {
    Write-Log -Message "Creating Domain: $DisplayName" -Level 'WARN'
    $payload = @{
      displayName = $DisplayName
      description = $DisplayName
    }

    if ($ParentDomainId) {
      $payload['parentDomainId'] = $ParentDomainId
    }

    $result = (Invoke-FabricRest -Method 'POST' -Endpoint "admin/domains" -Payload $payload).Response
  }

  if ($ParentDomainId) {
    Write-Log -Message "Child Domain - Name: $($result.displayName) / ID: $($result.id)"
  }
  else {
    Write-Log -Message "Parent Domain - Name: $($result.displayName) / ID: $($result.id)"
  }

  return $result
}

function Get-BaseName {
  param (
    [Parameter(Mandatory = $false)]
    [int]$Length = 10
  )

  $base = $Env:FABRIC_TESTACC_WELLKNOWN_NAME_BASE

  if (!$base) {
    $base = -join ((65..90) + (97..122) | Get-Random -Count $Length | ForEach-Object { [char]$_ })
  }

  return $base
}

function Get-DisplayName {
  param (
    [Parameter(Mandatory = $true)]
    [string]$Base,

    [Parameter(Mandatory = $false)]
    [string]$Prefix = $Env:FABRIC_TESTACC_WELLKNOWN_NAME_PREFIX,

    [Parameter(Mandatory = $false)]
    [string]$Suffix = $Env:FABRIC_TESTACC_WELLKNOWN_NAME_SUFFIX,

    [Parameter(Mandatory = $false)]
    [string]$Separator = '_'
  )

  $result = $Base

  # add prefix and suffix
  if ($Prefix) {
    $result = "${Prefix}${Separator}${result}"
  }

  if ($Suffix) {
    $result = "${result}${Separator}${Suffix}"
  }

  return $result
}

function Set-FabricWorkspace {
  param (
    [Parameter(Mandatory = $true)]
    [string]$DisplayName,

    [Parameter(Mandatory = $true)]
    [string]$CapacityId
  )

  $workspaces = Invoke-FabricRest -Method 'GET' -Endpoint 'workspaces'
  $workspace = $workspaces.Response.value | Where-Object { $_.displayName -eq $DisplayName }
  if (!$workspace) {
    Write-Log -Message "Creating Workspace: $DisplayName" -Level 'WARN'
    $payload = @{
      displayName = $DisplayName
      description = $DisplayName
      capacityId  = $CapacityId
    }
    $workspace = (Invoke-FabricRest -Method 'POST' -Endpoint 'workspaces' -Payload $payload).Response
  }

  return $workspace
}

function Set-FabricWorkspaceCapacity {
  param (
    [Parameter(Mandatory = $true)]
    [string]$WorkspaceId,

    [Parameter(Mandatory = $true)]
    [string]$CapacityId
  )

  $workspace = Invoke-FabricRest -Method 'GET' -Endpoint "workspaces/$WorkspaceId"
  if ($workspace.Response.capacityId -ne $CapacityId) {
    Write-Log -Message "Assigning Workspace to Capacity ID: $CapacityId" -Level 'WARN'
    $payload = @{
      capacityId = $CapacityId
    }
    _ = (Invoke-FabricRest -Method 'POST' -Endpoint "workspaces/$WorkspaceId/assignToCapacity" -Payload $payload).Response
    $workspace.Response.capacityId = $CapacityId
  }

  return $workspace.Response
}

function Set-FabricWorkspaceRoleAssignment {
  param (
    [Parameter(Mandatory = $true)]
    [string]$WorkspaceId,

    [Parameter(Mandatory = $true)]
    [object]$SG
  )

  $results = Invoke-FabricRest -Method 'GET' -Endpoint "workspaces/$WorkspaceId/roleAssignments"
  $result = $results.Response.value | Where-Object { $_.id -eq $SG.Id }
  if (!$result) {
    Write-Log -Message "Assigning SG to Workspace: $($SG.DisplayName)" -Level 'WARN'
    $payload = @{
      principal = @{
        id   = $SG.Id
        type = 'Group'
      }
      role      = 'Admin'
    }
    $result = (Invoke-FabricRest -Method 'POST' -Endpoint "workspaces/$WorkspaceId/roleAssignments" -Payload $payload).Response
  }
}

function Set-FabricGatewayVirtualNetwork {
  [CmdletBinding()]
  param(
    [Parameter(Mandatory = $true)]
    [string]$DisplayName,

    [Parameter(Mandatory = $true)]
    [string]$CapacityId,

    # Inactivity time (in minutes) before the gateway goes to auto-sleep.
    # Allowed values: 30, 60, 90, 120, 150, 240, 360, 480, 720, 1440.
    [Parameter(Mandatory = $true)]
    [ValidateSet(30, 60, 90, 120, 150, 240, 360, 480, 720, 1440)]
    [int]$InactivityMinutesBeforeSleep,

    # Number of member gateways (between 1 and 7).
    [Parameter(Mandatory = $true)]
    [ValidateRange(1, 7)]
    [int]$NumberOfMemberGateways,

    # Azure virtual network details:
    [Parameter(Mandatory = $true)]
    [string]$SubscriptionId,

    [Parameter(Mandatory = $true)]
    [string]$ResourceGroupName,

    [Parameter(Mandatory = $true)]
    [string]$VirtualNetworkName,

    [Parameter(Mandatory = $true)]
    [string]$SubnetName
  )

  # Attempt to check for an existing gateway with the same display name.
  $existingGateways = Invoke-FabricRest -Method 'GET' -Endpoint "gateways"
  $result = $existingGateways.Response.value | Where-Object { $_.displayName -eq $DisplayName }
  if (!$result) {
    # Construct the payload for creating a Virtual Network gateway.
    # Refer to the API documentation for details on the request format and the Virtual Network Azure Resource.
    $payload = @{
      type                         = "VirtualNetwork"
      displayName                  = $DisplayName
      capacityId                   = $CapacityId
      inactivityMinutesBeforeSleep = $InactivityMinutesBeforeSleep
      numberOfMemberGateways       = $NumberOfMemberGateways
      virtualNetworkAzureResource  = @{
        subscriptionId     = $SubscriptionId
        resourceGroupName  = $ResourceGroupName
        virtualNetworkName = $VirtualNetworkName
        subnetName         = $SubnetName
      }
    }

    Write-Log -Message "Creating Virtual Network Gateway: $DisplayName" -Level 'WARN'
    $newGateway = Invoke-FabricRest -Method 'POST' -Endpoint "gateways" -Payload $payload
    $result = $newGateway.Response
  }

  Write-Log -Message "Gateway Virtual Network - Name: $($result.displayName) / ID: $($result.id)"
  return $result
}

function Set-FabricGatewayRoleAssignment {
  param (
    [Parameter(Mandatory = $true)]
    [string]$GatewayId,

    [Parameter(Mandatory = $true)]
    [string]$PrincipalId,

    [Parameter(Mandatory = $true)]
    [ValidateSet('User', 'Group', 'ServicePrincipal')]
    [string]$PrincipalType,

    [Parameter(Mandatory = $true)]
    [ValidateSet('Admin', 'ConnectionCreator', 'ConnectionCreatorWithResharing')]
    [string]$Role
  )

  $results = Invoke-FabricRest -Method 'GET' -Endpoint "gateways/$GatewayId/roleAssignments"
  $result = $results.Response.value | Where-Object { $_.id -eq $PrincipalId }
  if (!$result) {
    Write-Log -Message "Assigning Principal ($PrincipalType / $PrincipalId) to Gateway: $($GatewayId)" -Level 'WARN'
    $payload = @{
      principal = @{
        id   = $PrincipalId
        type = $PrincipalType
      }
      role      = $Role
    }
    $result = (Invoke-FabricRest -Method 'POST' -Endpoint "gateways/$GatewayId/roleAssignments" -Payload $payload).Response
  }
}

function Set-AzureVirtualNetwork {
  param(
    [Parameter(Mandatory = $true)]
    [string]$ResourceGroupName,

    [Parameter(Mandatory = $true)]
    [string]$VNetName,

    [Parameter(Mandatory = $true)]
    [string]$Location,

    [Parameter(Mandatory = $true)]
    [string[]]$AddressPrefixes,

    [Parameter(Mandatory = $true)]
    [string]$SubnetName,

    [Parameter(Mandatory = $true)]
    [string[]]$SubnetAddressPrefixes,

    [Parameter(Mandatory = $true)]
    [object]$SG
  )

  # Attempt to get the existing Virtual Network
  try {
    $vnet = Get-AzVirtualNetwork -Name $VNetName -ResourceGroupName $ResourceGroupName -ErrorAction Stop
  }
  catch {
    # VNet does not exist, so create it
    Write-Log -Message "Creating Azure VNet: $VNetName in Resource Group: $ResourceGroupName" -Level 'WARN'
    $subnetConfig = New-AzVirtualNetworkSubnetConfig `
      -Name $SubnetName `
      -AddressPrefix $SubnetAddressPrefixes `

    $subnetConfig = Add-AzDelegation `
      -Name 'PowerPlatformVnetAccess' `
      -ServiceName 'Microsoft.PowerPlatform/vnetaccesslinks' `
      -Subnet $subnetConfig

    $vnet = New-AzVirtualNetwork `
      -Name $VNetName `
      -ResourceGroupName $ResourceGroupName `
      -Location $Location `
      -AddressPrefix $AddressPrefixes `
      -Subnet $subnetConfig

    # Commit creation
    $vnet = $vnet | Set-AzVirtualNetwork
    Write-Log -Message "Created Azure VNet: $VNetName" -Level 'INFO'
  }

  # If the VNet already exists, check for the subnet
  $subnet = $vnet.Subnets | Where-Object { $_.Name -eq $SubnetName }
  if (-not $subnet) {
    # Subnet does not exist; add one with the delegation
    Write-Log -Message "Adding subnet '$SubnetName' with delegation 'Microsoft.PowerPlatform/vnetaccesslinks' to VNet '$VNetName'." -Level 'WARN'
    $subnetConfig = New-AzVirtualNetworkSubnetConfig `
      -Name $SubnetName `
      -AddressPrefix $SubnetAddressPrefixes `

    $subnetConfig = Add-AzDelegation `
      -Name 'PowerPlatformVnetAccess' `
      -ServiceName 'Microsoft.PowerPlatform/vnetaccesslinks' `
      -Subnet $subnetConfig

    $vnet = $vnet | Set-AzVirtualNetwork
  }
  else {
    # Subnet exists; ensure it has the correct delegation
    $existingDelegation = $subnet.Delegations | Where-Object { $_.ServiceName -eq 'Microsoft.PowerPlatform/vnetaccesslinks' }
    if (-not $existingDelegation) {
      Write-Log -Message "Subnet '$SubnetName' found but missing delegation to 'Microsoft.PowerPlatform/vnetaccesslinks'. Adding Microsoft.PowerPlatform/vnetaccesslinks delegation..." -Level 'WARN'

      $subnetConfig = Add-AzDelegation `
        -Name 'PowerPlatformVnetAccess' `
        -ServiceName 'Microsoft.PowerPlatform/vnetaccesslinks' `
        -Subnet $subnet

      $vnet = $vnet | Set-AzVirtualNetwork
      Write-Log -Message "Added missing delegation to subnet '$SubnetName'." -Level 'INFO'
    }
  }
  Write-Log -Message "Az Virtual Network - Name: $($vnet.Name)"

  $userPrincipalName = $azContext.Account.Id
  $principal = Get-AzADUser -UserPrincipalName $userPrincipalName

  # Check if the principal already has the Network Contributor role on the VNet, if not then assign it.
  $existingAssignment = Get-AzRoleAssignment -Scope $vnet.Id -ObjectId $principal.Id -ErrorAction SilentlyContinue | Where-Object {
    $_.RoleDefinitionName -eq "Network Contributor"
  }
  Write-Log "Assigning Network Contributor role to the principal on the virtual network $($VNetName)"
  if (!$existingAssignment) {
    New-AzRoleAssignment -ObjectId $principal.Id -RoleDefinitionName "Network Contributor" -Scope $vnet.Id
  }

  # Check if the spns SG already has the Network Contributor role on the VNet, if not then assign it.
  $existingAssignment = Get-AzRoleAssignment -Scope $vnet.Id -ObjectId $SG.Id -ErrorAction SilentlyContinue | Where-Object {
    $_.RoleDefinitionName -eq "Network Contributor"
  }
  Write-Log "Assigning Network Contributor role to the spns security group $($SG.DisplayName) on the virtual network $($VNetName)"
  if (!$existingAssignment) {
    New-AzRoleAssignment -ObjectId $SG.Id -RoleDefinitionName "Network Contributor" -Scope $vnet.Id
  }

  return $vnet
}

function Set-AzureDataFactory {
  param(
    [Parameter(Mandatory = $true)]
    [string]$ResourceGroupName,

    [Parameter(Mandatory = $true)]
    [string]$DataFactoryName,

    [Parameter(Mandatory = $true)]
    [string]$Location,

    [Parameter(Mandatory = $true)]
    [object]$SG
  )

  # Register the Microsoft.DataFactory resource provider
  $df = Get-AzResourceProvider -ProviderNamespace "Microsoft.DataFactory"
  if ($df.RegistrationState -ne 'Registered') {
    Write-Log -Message "Registering Microsoft.DataFactory resource provider" -Level 'WARN'
    Register-AzResourceProvider -ProviderNamespace "Microsoft.DataFactory"
  }
  else {
    Write-Log -Message "Microsoft.DataFactory resource provider already registered" -Level 'INFO'
  }

  # Attempt to get the existing Data Factory
  try {
    $dataFactory = Get-AzDataFactoryV2 -ResourceGroupName $ResourceGroupName -Name $DataFactoryName -ErrorAction Stop
  }
  catch {
    # Data Factory does not exist, so create it
    Write-Log -Message "Creating Data Factory: $DataFactoryName in Resource Group: $ResourceGroupName" -Level 'WARN'
    $dataFactory = Set-AzDataFactoryV2 -ResourceGroupName $ResourceGroupName -Name $DataFactoryName -Location $Location
    Write-Log -Message "Created Data Factory: $DataFactoryName" -Level 'INFO'
  }

  Write-Log -Message "Az Data Factory - Name: $($dataFactory.DataFactoryName)"

  $userPrincipalName = $azContext.Account.Id
  $principal = Get-AzADUser -UserPrincipalName $userPrincipalName

  # Check if the principal already has the Data Factory Contributor role on the Data Factory, if not then assign it.
  $existingAssignment = Get-AzRoleAssignment -Scope $dataFactory.DataFactoryId -ObjectId $principal.Id -ErrorAction SilentlyContinue | Where-Object {
    $_.RoleDefinitionName -eq "Data Factory Contributor"
  }
  Write-Log "Assigning Data Factory Contributor role to the principal on the Data Factory $($dataFactory.DataFactoryName)"
  if (!$existingAssignment) {
    New-AzRoleAssignment -ObjectId $principal.Id -RoleDefinitionName "Data Factory Contributor" -Scope $dataFactory.DataFactoryId
  }

  # Check if the spns SG already has the Data Factory Contributor role on the Data Factory, if not then assign it.
  $existingAssignment = Get-AzRoleAssignment -Scope $dataFactory.DataFactoryId -ObjectId $SG.Id -ErrorAction SilentlyContinue | Where-Object {
    $_.RoleDefinitionName -eq "Data Factory Contributor"
  }
  Write-Log "Assigning Data Factory Contributor role to the spns security group $($SG.DisplayName) on the Data Factory $($dataFactory.DataFactoryName)"
  if (!$existingAssignment) {
    New-AzRoleAssignment -ObjectId $SG.Id -RoleDefinitionName "Data Factory Contributor" -Scope $dataFactory.DataFactoryId
  }

  return $dataFactory
}

function Set-Shortcut {
  param (
    [Parameter(Mandatory = $true)]
    [string]$WorkspaceId,
    # The unique identifier of the data item (e.g., lakehouse ID)
    [Parameter(Mandatory = $true)]
    [string]$ItemId,
    # OneLake data source payload
    [Parameter(Mandatory = $true)]
    [object]$Payload
  )

  # Attempt to get the existing shortcut
  $results = Invoke-FabricRest -Method 'GET' -Endpoint "workspaces/$WorkspaceId/items/$ItemId/shortcuts"
  $result = $results.Response.value | Where-Object { $_.name -eq $Payload.name -and ($_.path.TrimStart('/') -eq $Payload.path.TrimStart('/')) }

  if (!$result) {
    # Shortcut does not exist, so create it
    Write-Log -Message "Creating Shortcut: $($Payload.name)" -Level 'INFO'

    $result = (Invoke-FabricRest -Method 'POST' -Endpoint "workspaces/$WorkspaceId/items/$ItemId/shortcuts" -Payload $Payload).Response
  }
  $result.path = $result.path.TrimStart('/')
  Write-Log -Message "Shortcut - Name: $($result.name) / Path: $($result.path)"

  return $result
}

function Set-FabricFolder {
  param (
    [Parameter(Mandatory = $true)]
    [string]$WorkspaceId,

    [Parameter(Mandatory = $true)]
    [string]$DisplayName,

    [Parameter(Mandatory = $false)]
    [string]$ParentFolderId
  )
  # Attempt to get the existing folder
  $results = Invoke-FabricRest -Method 'GET' -Endpoint "workspaces/$WorkspaceId/folders"

  if (!$ParentFolderId) {
    # Looking for a root folder - filter for folders that don't have parentFolderId property
    $result = $results.Response.value | Where-Object {
      $_.displayName -eq $DisplayName -and
      (-not (Get-Member -InputObject $_ -Name "parentFolderId" -MemberType Properties))
    }
    Write-Log -Message "$result" -Level 'INFO'
  } else {
    # Looking for a folder with a specific parentFolderId
    $result = $results.Response.value | Where-Object {
      $_.displayName -eq $DisplayName -and
      $_.parentFolderId -eq $ParentFolderId
    }
  }

  if (!$result) {
    # Folder does not exist, so create it
    Write-Log -Message "Creating Folder: $DisplayName" -Level 'INFO'

    $payload = @{
      displayName    = $DisplayName
      parentFolderId = $ParentFolderId
    }

    $result = (Invoke-FabricRest -Method 'POST' -Endpoint "workspaces/$WorkspaceId/folders" -Payload $payload).Response
  }
  Write-Log -Message "Folder - Name: $($result.displayName) / ParentFolderId: $($result.parentFolderId)"

  return $result
}


# Define an array of modules to install
$modules = @('Az.Accounts', 'Az.Resources', 'Az.Storage', 'Az.Fabric', 'pwsh-dotenv', 'ADOPS', 'Az.Network', 'Az.DataFactory')

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
  !$Env:FABRIC_TESTACC_WELLKNOWN_AZURE_SUBSCRIPTION_ID -or
  !$Env:FABRIC_TESTACC_WELLKNOWN_FABRIC_CAPACITY_NAME -or
  !$Env:FABRIC_TESTACC_WELLKNOWN_AZDO_ORGANIZATION_NAME -or
  !$Env:FABRIC_TESTACC_WELLKNOWN_NAME_PREFIX -or
  !$Env:FABRIC_TESTACC_WELLKNOWN_AZURE_LOCATION -or
  !$Env:FABRIC_TESTACC_WELLKNOWN_AZURE_SPNS_SG_NAME
) {
  Write-Log -Message @'
  FABRIC_TESTACC_WELLKNOWN_ENTRA_TENANT_ID,
  FABRIC_TESTACC_WELLKNOWN_AZURE_SUBSCRIPTION_ID,
  FABRIC_TESTACC_WELLKNOWN_FABRIC_CAPACITY_NAME,
  FABRIC_TESTACC_WELLKNOWN_AZDO_ORGANIZATION_NAME,
  FABRIC_TESTACC_WELLKNOWN_NAME_PREFIX,
  FABRIC_TESTACC_WELLKNOWN_AZURE_LOCATION,
  and FABRIC_TESTACC_WELLKNOWN_AZURE_SPNS_SG_NAME
  are required environment variables.
'@ `
    -Level 'ERROR'
}

# Check if already logged in to Azure, if not then login
$azContext = Get-AzContext
if (!$azContext -or $azContext.Tenant.Id -ne $Env:FABRIC_TESTACC_WELLKNOWN_ENTRA_TENANT_ID -or $azContext.Subscription.Id -ne $Env:FABRIC_TESTACC_WELLKNOWN_AZURE_SUBSCRIPTION_ID) {
  Write-Log -Message 'Logging in to Azure.' -Level 'DEBUG'
  Connect-AzAccount -Tenant $Env:FABRIC_TESTACC_WELLKNOWN_ENTRA_TENANT_ID -SubscriptionId $Env:FABRIC_TESTACC_WELLKNOWN_AZURE_SUBSCRIPTION_ID -UseDeviceAuthentication
  $azContext = Get-AzContext
  # Disconnect-AzAccount
}
# $currentUser = Get-AzADUser -SignedIn

# Logged in to Azure DevOps
Write-Log -Message 'Logging in to Azure DevOps.' -Level 'DEBUG'
$secureAccessToken = (Get-AzAccessToken -WarningAction SilentlyContinue -AsSecureString -ResourceUrl '499b84ac-1321-427f-aa17-267ca6975798').Token
$unsecureAccessToken = $secureAccessToken | ConvertFrom-SecureString -AsPlainText
$azdoContext = Connect-ADOPS -TenantId $azContext.Tenant.Id -Organization $Env:FABRIC_TESTACC_WELLKNOWN_AZDO_ORGANIZATION_NAME -OAuthToken $unsecureAccessToken

$SPNS_SG = Get-AzADGroup -DisplayName $Env:FABRIC_TESTACC_WELLKNOWN_AZURE_SPNS_SG_NAME

$wellKnown = @{}

# Get Fabric Capacity ID
$capacities = Invoke-FabricRest -Method 'GET' -Endpoint 'capacities'
$capacity = $capacities.Response.value | Where-Object { $_.displayName -eq $Env:FABRIC_TESTACC_WELLKNOWN_FABRIC_CAPACITY_NAME }
if (!$capacity) {
  Write-Log -Message "Fabric Capacity: $($Env:FABRIC_TESTACC_WELLKNOWN_FABRIC_CAPACITY_NAME)"
}
Write-Log -Message "Fabric Capacity - Name: $($Env:FABRIC_TESTACC_WELLKNOWN_FABRIC_CAPACITY_NAME) / ID: $($capacity.id)"
$wellKnown['Capacity'] = @{
  id          = $capacity.id
  displayName = $capacity.displayName
  sku         = $capacity.sku
}

$itemNaming = @{
  'ApacheAirflowJob'       = 'aaj'
  'AzureDataFactory'       = 'adf'
  'CopyJob'                = 'cj'
  'Dashboard'              = 'dash'
  'Dataflow'               = 'df'
  'Datamart'               = 'dm'
  'DataPipeline'           = 'dp'
  'DeploymentPipeline'     = 'deployp'
  'Environment'            = 'env'
  'Eventhouse'             = 'eh'
  'Eventstream'            = 'es'
  'Folder'                 = 'fld'
  'GraphQLApi'             = 'gql'
  'KQLDashboard'           = 'kqldash'
  'KQLDatabase'            = 'kqldb'
  'KQLQueryset'            = 'kqlqs'
  'Lakehouse'              = 'lh'
  'MirroredDatabase'       = 'mdb'
  'MirroredWarehouse'      = 'mwh'
  'MLExperiment'           = 'mle'
  'MLModel'                = 'mlm'
  'MountedDataFactory'     = 'mdf'
  'Notebook'               = 'nb'
  'Shortcut'               = 'srt'
  'PaginatedReport'        = 'prpt'
  'Reflex'                 = 'rx'
  'Report'                 = 'rpt'
  'SemanticModel'          = 'sm'
  'SparkJobDefinition'     = 'sjd'
  'SQLDatabase'            = 'sqldb'
  'SQLEndpoint'            = 'sqle'
  'Warehouse'              = 'wh'
  'WorkspaceDS'            = 'wsds'
  'WorkspaceRS'            = 'wsrs'
  'WorkspaceMPE'           = 'wsmpe'
  'DomainParent'           = 'parent'
  'DomainChild'            = 'child'
  'EntraServicePrincipal'  = 'sp'
  'EntraGroup'             = 'grp'
  'AzDOProject'            = 'proj'
  'VirtualNetwork01'       = 'vnet01'
  'VirtualNetwork02'       = 'vnet02'
  'VirtualNetworkSubnet'   = 'subnet'
  'GatewayVirtualNetwork'  = 'gvnet'
  'ManagedPrivateEndpoint' = 'mpe'
  'StorageAccount'         = 'st'
  'ResourceGroup'          = 'rg'
  'FabricCapacity'         = 'fc'
}

$baseName = Get-BaseName
$Env:FABRIC_TESTACC_WELLKNOWN_NAME_BASE = $baseName

# Save env vars wellknown.env file
$envVarNames = @(
  'FABRIC_TESTACC_WELLKNOWN_ENTRA_TENANT_ID',
  'FABRIC_TESTACC_WELLKNOWN_AZURE_SUBSCRIPTION_ID',
  'FABRIC_TESTACC_WELLKNOWN_AZURE_LOCATION',
  'FABRIC_TESTACC_WELLKNOWN_AZURE_SPNS_SG_NAME',
  'FABRIC_TESTACC_WELLKNOWN_FABRIC_CAPACITY_NAME',
  'FABRIC_TESTACC_WELLKNOWN_AZDO_ORGANIZATION_NAME',
  'FABRIC_TESTACC_WELLKNOWN_NAME_PREFIX',
  'FABRIC_TESTACC_WELLKNOWN_NAME_SUFFIX',
  'FABRIC_TESTACC_WELLKNOWN_NAME_BASE',
  'FABRIC_TESTACC_WELLKNOWN_SPN_NAME',
  'FABRIC_TESTACC_WELLKNOWN_GITHUB_CONNECTION_ID'
)

$envVars = $envVarNames | ForEach-Object {
  $envVarName = $_
  if (Test-Path "Env:${envVarName}") {
    $value = (Get-ChildItem "Env:${envVarName}").Value
    "$envVarName=`"$value`""
  }
}

$envVars -join "`n" | Set-Content -Path './wellknown.env' -Force -NoNewline -Encoding utf8

$displayName = Get-DisplayName -Base $baseName

# Create WorkspaceMPE if not exists
$displayNameTemp = "${displayName}_$($itemNaming['WorkspaceMPE'])"
$workspace = Set-FabricWorkspace -DisplayName $displayNameTemp -CapacityId $capacity.id

# Assign WorkspaceMPE to Capacity if not already assigned or assigned to a different capacity
$workspace = Set-FabricWorkspaceCapacity -WorkspaceId $workspace.id -CapacityId $capacity.id

Write-Log -Message "WorkspaceMPE - Name: $($workspace.displayName) / ID: $($workspace.id)"
$wellKnown['WorkspaceMPE'] = @{
  id          = $workspace.id
  displayName = $workspace.displayName
  description = $workspace.description
}
# Assign SPN to WorkspaceMPE if not already assigned
Set-FabricWorkspaceRoleAssignment -WorkspaceId $workspace.id -SG $SPNS_SG

# Create WorkspaceRS if not exists
$displayNameTemp = "${displayName}_$($itemNaming['WorkspaceRS'])"
$workspace = Set-FabricWorkspace -DisplayName $displayNameTemp -CapacityId $capacity.id

# Assign WorkspaceRS to Capacity if not already assigned or assigned to a different capacity
$workspace = Set-FabricWorkspaceCapacity -WorkspaceId $workspace.id -CapacityId $capacity.id

Write-Log -Message "WorkspaceRS - Name: $($workspace.displayName) / ID: $($workspace.id)"
$wellKnown['WorkspaceRS'] = @{
  id          = $workspace.id
  displayName = $workspace.displayName
  description = $workspace.description
}

# Assign SPN to WorkspaceRS if not already assigned
Set-FabricWorkspaceRoleAssignment -WorkspaceId $workspace.id -SG $SPNS_SG

# Create WorkspaceDS if not exists
$displayNameTemp = "${displayName}_$($itemNaming['WorkspaceDS'])"
$workspace = Set-FabricWorkspace -DisplayName $displayNameTemp -CapacityId $capacity.id

# Assign WorkspaceDS to Capacity if not already assigned or assigned to a different capacity
$workspace = Set-FabricWorkspaceCapacity -WorkspaceId $workspace.id -CapacityId $capacity.id

Write-Log -Message "WorkspaceDS - Name: $($workspace.displayName) / ID: $($workspace.id)"
$wellKnown['WorkspaceDS'] = @{
  id          = $workspace.id
  displayName = $workspace.displayName
  description = $workspace.description
}

# Assign SPN to WorkspaceDS if not already assigned
Set-FabricWorkspaceRoleAssignment -WorkspaceId $workspace.id -SG $SPNS_SG

# Define an array of item types to create
$itemTypes = @('CopyJob', 'Dataflow', 'DataPipeline', 'DigitalTwinBuilder', 'Environment', 'Eventhouse', 'GraphQLApi', 'KQLDashboard', 'KQLQueryset', 'Lakehouse', 'MLExperiment', 'MLModel', 'Notebook', 'Reflex', 'SparkJobDefinition', 'SQLDatabase', 'Warehouse')

# Loop through each item type and create if not exists
foreach ($itemType in $itemTypes) {

  $displayNameTemp = "${displayName}_$($itemNaming[$itemType])"
  $item = Set-FabricItem -DisplayName $displayNameTemp -WorkspaceId $wellKnown['WorkspaceDS'].id -Type $itemType
  $wellKnown[$itemType] = @{
    id          = $item.id
    displayName = $item.displayName
    description = $item.description
  }
}

# Create KQLDatabase if not exists
$displayNameTemp = "${displayName}_$($itemNaming['KQLDatabase'])"
$creationPayload = @{
  databaseType           = 'ReadWrite'
  parentEventhouseItemId = $wellKnown['Eventhouse'].id
}
$kqlDatabase = Set-FabricItem -DisplayName $displayNameTemp -WorkspaceId $wellKnown['WorkspaceDS'].id -Type 'KQLDatabase' -CreationPayload $creationPayload
$wellKnown['KQLDatabase'] = @{
  id          = $kqlDatabase.id
  displayName = $kqlDatabase.displayName
  description = $kqlDatabase.description
}

# Create MirroredDatabase if not exists
$displayNameTemp = "${displayName}_$($itemNaming['MirroredDatabase'])"
$definition = @{
  parts = @(
    @{
      path        = "mirroring.json"
      payload     = Get-DefinitionPartBase64 -Path 'internal/testhelp/fixtures/mirrored_database/mirroring.json.tmpl' -Values @(@{ key = '{{ .DEFAULT_SCHEMA }}'; value = 'dbo' })
      payloadType = 'InlineBase64'
    }
  )
}
$mirroredDatabase = Set-FabricItem -DisplayName $displayNameTemp -WorkspaceId $wellKnown['WorkspaceDS'].id -Type 'MirroredDatabase' -Definition $definition
$wellKnown['MirroredDatabase'] = @{
  id          = $mirroredDatabase.id
  displayName = $mirroredDatabase.displayName
  description = $mirroredDatabase.description
}

# Create SemanticModel if not exists
$displayNameTemp = "${displayName}_$($itemNaming['SemanticModel'])"
$definition = @{
  parts = @(
    @{
      path        = 'definition.pbism'
      payload     = Get-DefinitionPartBase64 -Path 'internal/testhelp/fixtures/semantic_model_tmsl/definition.pbism'
      payloadType = 'InlineBase64'
    }
    @{
      path        = 'model.bim'
      payload     = Get-DefinitionPartBase64 -Path 'internal/testhelp/fixtures/semantic_model_tmsl/model.bim.tmpl' -Values @(@{ key = '{{ .ColumnName }}'; value = 'ColumnTest1' })
      payloadType = 'InlineBase64'
    }
  )
}

$semanticModel = Set-FabricItem -DisplayName $displayNameTemp -WorkspaceId $wellKnown['WorkspaceDS'].id -Type 'SemanticModel' -Definition $definition
$wellKnown['SemanticModel'] = @{
  id          = $semanticModel.id
  displayName = $semanticModel.displayName
  description = $semanticModel.description
}

# Create Report if not exists
$displayNameTemp = "${displayName}_$($itemNaming['Report'])"
$definition = @{
  parts = @(
    @{
      path        = 'definition.pbir'
      payload     = Get-DefinitionPartBase64 -Path 'internal/testhelp/fixtures/report_pbir_legacy/definition.pbir.tmpl' -Values @(@{ key = '{{ .SemanticModelID }}'; value = $semanticModel.id })
      payloadType = 'InlineBase64'
    },
    @{
      path        = 'report.json'
      payload     = Get-DefinitionPartBase64 -Path 'internal/testhelp/fixtures/report_pbir_legacy/report.json'
      payloadType = 'InlineBase64'
    },
    @{
      path        = 'StaticResources/SharedResources/BaseThemes/CY24SU10.json'
      payload     = Get-DefinitionPartBase64 -Path 'internal/testhelp/fixtures/report_pbir_legacy/StaticResources/SharedResources/BaseThemes/CY24SU10.json'
      payloadType = 'InlineBase64'
    }
    @{
      path        = 'StaticResources/RegisteredResources/fabric_48_color10148978481469717.png'
      payload     = Get-DefinitionPartBase64 -Path 'internal/testhelp/fixtures/report_pbir_legacy/StaticResources/RegisteredResources/fabric_48_color10148978481469717.png'
      payloadType = 'InlineBase64'
    }
  )
}
$report = Set-FabricItem -DisplayName $displayNameTemp -WorkspaceId $wellKnown['WorkspaceDS'].id -Type 'Report' -Definition $definition
$wellKnown['Report'] = @{
  id          = $report.id
  displayName = $report.displayName
  description = $report.description
}

# Create Deployment Pipeline if not exists
$displayNameTemp = "${displayName}_$($itemNaming['DeploymentPipeline'])"
$deploymentPipeline = Set-DeploymentPipeline -DisplayName $displayNameTemp
$wellKnown['DeploymentPipeline'] = @{
  id          = $deploymentPipeline.id
  displayName = $deploymentPipeline.displayName
  description = $deploymentPipeline.description
  stages      = $deploymentPipeline.stages
}

# Create Eventstream if not exists
$displayNameTemp = "${displayName}_$($itemNaming['Eventstream'])"
$definition = @{
  parts = @(
    @{
      path        = "eventstream.json"
      payload     = Get-DefinitionPartBase64 -Path 'internal/testhelp/fixtures/eventstream/eventstream.json.tmpl' -Values @(
        @{ key = '{{ .LakehouseID }}'; value = $wellKnown['Lakehouse'].id },
        @{ key = '{{ .LakehouseWorkspaceID }}'; value = $wellKnown['WorkspaceDS'].id }
      )
      payloadType = 'InlineBase64'
    }
  )
}
$eventstream = Set-FabricItem -DisplayName $displayNameTemp -WorkspaceId $wellKnown['WorkspaceDS'].id -Type 'Eventstream' -Definition $definition
$wellKnown['Eventstream'] = @{
  id          = $eventstream.id
  displayName = $eventstream.displayName
  description = $eventstream.description
}

# Set Eventstream source connection
$eventstreamTopology = (Invoke-FabricRest -Method 'GET' -Endpoint "workspaces/$($wellKnown['WorkspaceDS'].id)/eventstreams/$($eventstream.id)/topology").Response
$eventstreamSource = $eventstreamTopology.sources | Where-Object { $_.type -eq 'CustomEndpoint' } | Select-Object -First 1
$eventstreamSourceId = $eventstreamSource.id

$eventstreamConnection = (Invoke-FabricRest -Method 'GET' -Endpoint "workspaces/$($wellKnown['WorkspaceDS'].id)/eventstreams/$($eventstream.id)/sources/$($eventstreamSourceId)/connection").Response
$wellKnown['Eventstream']['sourceConnection'] = @{
  sourceId                = $eventstreamSourceId
  eventHubName            = $eventstreamConnection.eventHubName
  fullyQualifiedNamespace = $eventstreamConnection.fullyQualifiedNamespace
}



# Create Parent Domain if not exists
$displayNameTemp = "${displayName}_$($itemNaming['DomainParent'])"
$parentDomain = Set-FabricDomain -DisplayName $displayNameTemp
$wellKnown['DomainParent'] = @{
  id          = $parentDomain.id
  displayName = $parentDomain.displayName
  description = $parentDomain.description
}

# Create Child Domain if not exists
$displayNameTemp = "${displayName}_$($itemNaming['DomainChild'])"
$childDomain = Set-FabricDomain -DisplayName $displayNameTemp -ParentDomainId $parentDomain.id
$wellKnown['DomainChild'] = @{
  id          = $childDomain.id
  displayName = $childDomain.displayName
  description = $childDomain.description
}

$IS_LAKEHOUSE_POPULATED = $false
$results = Invoke-FabricRest -Method 'GET' -Endpoint "workspaces/$($wellKnown['WorkspaceDS'].id)/lakehouses/$($wellKnown['Lakehouse']['id'])/tables"
$result = $results.Response.data | Where-Object { $_.name -eq 'publicholidays' }
if (!$result) {
  Write-Log -Message "!!! Please go to the Lakehouse and manually run 'Start with sample data' -> 'Public holidays' to populate the data !!!" -Level 'ERROR' -Stop $false
  Write-Log -Message "Lakehouse: https://app.fabric.microsoft.com/groups/$($wellKnown['WorkspaceDS'].id)/lakehouses/$($wellKnown['Lakehouse']['id'])" -Level 'WARN'
}
else {
  $IS_LAKEHOUSE_POPULATED = $true
}
$wellKnown['Lakehouse']['tableName'] = 'publicholidays'

$displayNameTemp = "${displayName}_$($itemNaming['Dashboard'])"
$results = Invoke-FabricRest -Method 'GET' -Endpoint "workspaces/$($wellKnown['WorkspaceDS'].id)/dashboards"
$result = $results.Response.value | Where-Object { $_.displayName -eq $displayNameTemp }
if (!$result) {
  Write-Log -Message "!!! Please create a Dashboard manually (with Display Name: ${displayNameTemp}), and update details in the well-known file !!!" -Level 'ERROR' -Stop $false
  Write-Log -Message "Workspace: https://app.fabric.microsoft.com/groups/$($wellKnown['WorkspaceDS'].id)" -Level 'WARN'
}
$wellKnown['Dashboard'] = @{
  id          = if ($result) { $result.id } else { '00000000-0000-0000-0000-000000000000' }
  displayName = if ($result) { $result.displayName } else { $displayNameTemp }
  description = if ($result) { $result.description } else { '' }
}

$displayNameTemp = "${displayName}_$($itemNaming['Datamart'])"
$results = Invoke-FabricRest -Method 'GET' -Endpoint "workspaces/$($wellKnown['WorkspaceDS'].id)/datamarts"
$result = $results.Response.value | Where-Object { $_.displayName -eq $displayNameTemp }
if (!$result) {
  Write-Log -Message "!!! Please create a Datamart manually (with Display Name: ${displayNameTemp}), and update details in the well-known file !!!" -Level 'ERROR' -Stop $false
  Write-Log -Message "Workspace: https://app.fabric.microsoft.com/groups/$($wellKnown['WorkspaceDS'].id)" -Level 'WARN'
}
$wellKnown['Datamart'] = @{
  id          = if ($result) { $result.id } else { '00000000-0000-0000-0000-000000000000' }
  displayName = if ($result) { $result.displayName } else { $displayNameTemp }
  description = if ($result) { $result.description } else { '' }
}

# Create Resource Group if not exists
$displayNameTemp = "$($itemNaming['ResourceGroup'])-${displayName}"
$resourceGroup = Get-AzResourceGroup -Name $displayNameTemp
if (!$resourceGroup) {
  Write-Log -Message "Creating Resource Group: $displayNameTemp" -Level 'WARN'
  $resourceGroup = New-AzResourceGroup -Name $displayNameTemp -Location $Env:FABRIC_TESTACC_WELLKNOWN_AZURE_LOCATION
}
Write-Log -Message "Resource Group - Name: $($resourceGroup.ResourceGroupName) / ID: $($resourceGroup.ResourceId)"
$wellKnown['ResourceGroup'] = @{
  id       = $resourceGroup.ResourceId
  name     = $resourceGroup.ResourceGroupName
  location = $resourceGroup.Location
}

# Set Azure Context
$wellKnown['Azure'] = @{
  subscriptionId = $Env:FABRIC_TESTACC_WELLKNOWN_AZURE_SUBSCRIPTION_ID
  location       = $wellKnown['ResourceGroup'].location
}

# Create Storage Account if not exists
$displayNameTemp = "$($itemNaming['StorageAccount'])${displayName}".ToLower() -replace '[^a-z0-9]', ''
$storageAccount = Get-AzStorageAccount -ResourceGroupName $wellKnown['ResourceGroup'].name -Name $displayNameTemp -ErrorAction SilentlyContinue
if (!$storageAccount) {
  Write-Log -Message "Creating Storage Account: $displayNameTemp" -Level 'WARN'
  $storageAccount = New-AzStorageAccount -ResourceGroupName $wellKnown['ResourceGroup'].name -Name $displayNameTemp -SkuName Standard_LRS -Kind StorageV2 -Location $wellKnown['ResourceGroup'].location
}
Write-Log -Message "Storage Account - Name: $($storageAccount.StorageAccountName) / ID: $($storageAccount.Id)"
$wellKnown['StorageAccount'] = @{
  id   = $storageAccount.Id
  name = $storageAccount.StorageAccountName
}

# Create Managed Private Endpoint if not exists
$displayNameTemp = "${displayName}_$($itemNaming['ManagedPrivateEndpoint'])"
$managedPrivateEndpoints = Invoke-FabricRest -Method 'GET' -Endpoint "workspaces/$($wellKnown['WorkspaceMPE'].id)/managedPrivateEndpoints"
$managedPrivateEndpoint = $managedPrivateEndpoints.Response.value | Where-Object { $_.name -eq $displayNameTemp }
if (!$managedPrivateEndpoint) {
  Write-Log -Message "Creating Managed Private Endpoint: $displayNameTemp" -Level 'WARN'
  $payload = @{
    name                        = $displayNameTemp
    targetPrivateLinkResourceId = $wellKnown['StorageAccount']['id']
    targetSubresourceType       = 'blob'
    requestMessage              = $displayNameTemp
  }
  $managedPrivateEndpoint = (Invoke-FabricRest -Method 'POST' -Endpoint "workspaces/$($wellKnown['WorkspaceMPE'].id)/managedPrivateEndpoints" -Payload $payload).Response
}
Write-Log -Message "Managed Private Endpoint - Name: $($managedPrivateEndpoint.name) / ID: $($managedPrivateEndpoint.id)"
$wellKnown['ManagedPrivateEndpoint'] = @{
  id   = $managedPrivateEndpoint.id
  name = $managedPrivateEndpoint.name
}

# Create SP if not exists
$displayNameTemp = "${displayName}_$($itemNaming['EntraServicePrincipal'])"
$entraSp = Get-AzADServicePrincipal -DisplayName $displayNameTemp
if (!$entraSp) {
  Write-Log -Message "Creating Service Principal: $displayNameTemp" -Level 'WARN'
  $entraApp = New-AzADApplication -DisplayName $displayNameTemp
  $entraSp = New-AzADServicePrincipal -ApplicationId $entraApp.AppId
}
Write-Log -Message "Service Principal - Name: $($entraSp.DisplayName) / ID: $($entraSp.id)"
$wellKnown['Principal'] = @{
  id    = $entraSp.Id
  type  = 'ServicePrincipal'
  name  = $entraSp.DisplayName
  appId = $entraSp.AppId
}

# Create Group if not exists
$displayNameTemp = "${displayName}_$($itemNaming['EntraGroup'])"
$entraGroup = Get-AzADGroup -DisplayName $displayNameTemp
if (!$entraGroup) {
  Write-Log -Message "Creating Group: $displayNameTemp" -Level 'WARN'
  $entraGroup = New-AzADGroup -DisplayName $displayNameTemp -MailNickname $displayNameTemp -SecurityEnabled
  # New-AzADGroupOwner -GroupId $entraGroup.Id -OwnerId $currentUser.Id
}
Write-Log -Message "Group - Name: $($entraGroup.DisplayName) / ID: $($entraGroup.Id)"
$wellKnown['Group'] = @{
  type = 'Group'
  id   = $entraGroup.Id
  name = $entraGroup.DisplayName
}

# Create AzDO Project if not exists
$displayNameTemp = "${displayName}_$($itemNaming['AzDOProject'])"
$azdoProject = Get-ADOPSProject -Name $displayNameTemp
if (!$azdoProject) {
  Write-Log -Message "Creating AzDO Project: $displayNameTemp" -Level 'WARN'
  $azdoProject = New-ADOPSProject -Name $displayNameTemp -Visibility Private -Wait
}
Write-Log -Message "AzDO Project - Name: $($azdoProject.name) / ID: $($azdoProject.id)"

# Create AzDO Repository if not exists
$azdoRepo = Get-ADOPSRepository -Project $azdoProject.name -Repository 'test'
if (!$azdoRepo) {
  Write-Log -Message "Creating AzDO Repository: test" -Level 'WARN'
  $azdoRepo = New-ADOPSRepository -Project $azdoProject.name -Name 'test'
  Initialize-ADOPSRepository -RepositoryId $azdoRepo.id | Out-Null
}
Write-Log -Message "AzDO Repository - Name: $($azdoRepo.name) / ID: $($azdoRepo.id)"
$wellKnown['AzDO'] = @{
  organizationName = $azdoContext.Organization
  projectId        = $azdoProject.id
  projectName      = $azdoProject.name
  repositoryId     = $azdoRepo.id
  repositoryName   = $azdoRepo.name
}

$body = @{
  originId = $SPNS_SG.Id
}
$bodyJson = $body | ConvertTo-Json
$azdoSG = Invoke-ADOPSRestMethod -Uri "https://vssps.dev.azure.com/$($azdoContext.Organization)/_apis/graph/groups?api-version=7.2-preview.1" -Method Post -Body $bodyJson
$result = Set-ADOPSGitPermission -ProjectId $azdoProject.id -RepositoryId $azdoRepo.id -Descriptor $azdoSG.descriptor -Allow 'GenericContribute', 'PullRequestContribute', 'CreateBranch', 'CreateTag', 'GenericRead'

# Set GitHub
if (!$Env:FABRIC_TESTACC_WELLKNOWN_GITHUB_CONNECTION_ID) {
  Write-Log -Message "!!! Please go to the Connections and manually add 'GitHub - Source control' connection !!!" -Level 'ERROR' -Stop $false
  Write-Log -Message "Connections: https://app.fabric.microsoft.com/groups/me/gateways" -Level 'ERROR' -Stop $false
  Write-Log -Message "and set FABRIC_TESTACC_WELLKNOWN_GITHUB_CONNECTION_ID" -Level 'ERROR' -Stop $true
}

$results = Invoke-FabricRest -Method 'GET' -Endpoint "connections/$Env:FABRIC_TESTACC_WELLKNOWN_GITHUB_CONNECTION_ID"
Write-Log -Message "GitHub - Name: $($results.Response.displayName) / ID: $($results.Response.id) / Path: $($results.Response.connectionDetails.path)"
$segments = $results.Response.connectionDetails.path.TrimEnd('/') -split '/'
$wellKnown['GitHub'] = @{
  connectionId   = $Env:FABRIC_TESTACC_WELLKNOWN_GITHUB_CONNECTION_ID
  ownerName      = $segments[3]
  repositoryName = $segments[4]
}

# Register the Microsoft.PowerPlatform resource provider
$pp = Get-AzResourceProvider -ProviderNamespace 'Microsoft.PowerPlatform'
if ($pp.RegistrationState -ne 'Registered') {
  Write-Log -Message 'Registering Microsoft.PowerPlatform resource provider' -Level 'WARN'
  Register-AzResourceProvider -ProviderNamespace 'Microsoft.PowerPlatform'
}
else {
  Write-Log -Message 'Microsoft.PowerPlatform resource provider already registered' -Level 'INFO'
}

# Create Azure Virtual Network 1 if not exists
$vnetName = "${displayName}_$($itemNaming['VirtualNetwork01'])"
$addrRange = '10.10.0.0/16'
$subName = "${displayName}_$($itemNaming['VirtualNetworkSubnet'])"
$subRange = '10.10.1.0/24'

$vnet = Set-AzureVirtualNetwork `
  -ResourceGroupName $wellKnown['ResourceGroup'].name `
  -VNetName $vnetName `
  -Location $wellKnown['ResourceGroup'].location `
  -AddressPrefixes $addrRange `
  -SubnetName $subName `
  -SubnetAddressPrefixes $subRange `
  -SG $SPNS_SG

$wellKnown['VirtualNetwork01'] = @{
  name              = $vnet.Name
  resourceGroupName = $wellKnown['ResourceGroup'].name
  subnetName        = $subName
  subscriptionId    = $wellKnown['Azure'].subscriptionId
}

# Create Azure Virtual Network 2 if not exists
$vnetName = "${displayName}_$($itemNaming['VirtualNetwork02'])"
$addrRange = '10.10.0.0/16'
$subName = "${displayName}_$($itemNaming['VirtualNetworkSubnet'])"
$subRange = '10.10.1.0/24'

$vnet = Set-AzureVirtualNetwork `
  -ResourceGroupName $wellKnown['ResourceGroup'].name `
  -VNetName $vnetName `
  -Location $wellKnown['ResourceGroup'].location `
  -AddressPrefixes $addrRange `
  -SubnetName $subName `
  -SubnetAddressPrefixes $subRange `
  -SG $SPNS_SG

$wellKnown['VirtualNetwork02'] = @{
  name              = $vnet.Name
  resourceGroupName = $wellKnown['ResourceGroup'].name
  subnetName        = $subName
  subscriptionId    = $wellKnown['Azure'].subscriptionId
}

# Create Fabric Gateway Virtual Network if not exists
$displayNameTemp = "${displayName}_$($itemNaming['GatewayVirtualNetwork'])"
$inactivityMinutesBeforeSleep = 30
$numberOfMemberGateways = 1

$gateway = Set-FabricGatewayVirtualNetwork `
  -DisplayName $displayNameTemp `
  -CapacityId $capacity.id `
  -InactivityMinutesBeforeSleep $inactivityMinutesBeforeSleep `
  -NumberOfMemberGateways $numberOfMemberGateways `
  -SubscriptionId $wellKnown['Azure'].subscriptionId `
  -ResourceGroupName $wellKnown['ResourceGroup'].name `
  -VirtualNetworkName $wellKnown['VirtualNetwork01'].name `
  -SubnetName $wellKnown['VirtualNetwork01'].subnetName

$wellKnown['GatewayVirtualNetwork'] = @{
  id          = $gateway.id
  displayName = $gateway.displayName
  type        = $gateway.type
}

Set-FabricGatewayRoleAssignment -GatewayId $gateway.id -PrincipalId $SPNS_SG.Id -PrincipalType 'Group' -Role 'Admin'
Set-FabricGatewayRoleAssignment -GatewayId $gateway.id -PrincipalId $wellKnown['Principal'].id -PrincipalType $wellKnown['Principal'].type -Role 'ConnectionCreator'

# Create the Azure Data Factory if not exists
$displayNameTemp = "$Env:FABRIC_TESTACC_WELLKNOWN_NAME_PREFIX-$Env:FABRIC_TESTACC_WELLKNOWN_NAME_BASE-$($itemNaming['AzureDataFactory'])"

$dataFactory = Set-AzureDataFactory `
  -ResourceGroupName $wellKnown['ResourceGroup'].name `
  -DataFactoryName $displayNameTemp `
  -Location $wellKnown['ResourceGroup'].location `
  -SG $SPNS_SG

$wellKnown['AzureDataFactory'] = @{
  name              = $displayNameTemp
  resourceGroupName = $wellKnown['ResourceGroup'].name
  location          = $wellKnown['ResourceGroup'].location
  subscriptionId    = $wellKnown['Azure'].subscriptionId
}

# Create the Mounted Data Factory if not exists
$displayNameTemp = "${displayName}_$($itemNaming['MountedDataFactory'])"
$definition = @{
  parts = @(
    @{
      path        = 'mountedDataFactory-content.json'
      payload     = Get-DefinitionPartBase64 -Path 'internal/testhelp/fixtures/mounted_data_factory/mountedDataFactory-content.json.tmpl' -Values @(
        @{ key = '{{ .SUBSCRIPTION_ID }}'; value = $wellKnown['Azure'].subscriptionId },
        @{ key = '{{ .RESOURCE_GROUP_NAME }}'; value = $wellKnown['ResourceGroup'].name },
        @{ key = '{{ .FACTORY_NAME }}'; value = $dataFactory.DataFactoryName }
      )
      payloadType = 'InlineBase64'
    }
  )
}

$mountedDataFactory = Set-FabricItem -DisplayName $displayNameTemp -WorkspaceId $wellKnown['WorkspaceDS'].id -Type 'MountedDataFactory' -Definition $definition

$wellKnown['MountedDataFactory'] = @{
  id          = $mountedDataFactory.id
  displayName = $mountedDataFactory.displayName
  description = $mountedDataFactory.description
}

# Create the Apache Airflow Job if not exists
$displayNameTemp = "${displayName}_$($itemNaming['ApacheAirflowJob'])"
$definition = @{
  parts = @(
    @{
      path        = 'apacheAirflowJob-content.json'
      payload     = Get-DefinitionPartBase64 -Path 'internal/testhelp/fixtures/apache_airflow_job/apacheairflowjob-content.json.tmpl'
      payloadType = 'InlineBase64'
    }
  )
}

$apacheAirflowJob = Set-FabricItem -DisplayName $displayNameTemp -WorkspaceId $wellKnown['WorkspaceDS'].id -Type 'ApacheAirflowJob' -Definition $definition

$wellKnown['ApacheAirflowJob'] = @{
  id          = $apacheAirflowJob.id
  displayName = $apacheAirflowJob.displayName
  description = $apacheAirflowJob.description
}

$displayNameTemp = "${displayName}_$($itemNaming['Shortcut'])"
if ($IS_LAKEHOUSE_POPULATED -eq $false) {
  Write-Log -Message "Lakehouse is not populated. Skipping shortcut creation." -Level 'ERROR' -Stop:$false

  $wellKnown['Shortcut'] = @{
    shortcutName = $displayNameTemp
    shortcutPath = ''
    workspaceId  = $wellKnown['WorkspaceDS'].id
    lakehouseId  = $wellKnown['Lakehouse'].id
  }
}
else {
  $TABLES_PATH = "Tables"
  $shortcutPayload = @{
    path   = $TABLES_PATH
    name   = $displayNameTemp
    target = @{
      onelake = @{
        workspaceId = $wellKnown['WorkspaceDS'].id
        itemId      = $wellKnown['Lakehouse'].id
        path        = $TABLES_PATH + "/" + $wellKnown['Lakehouse'].tableName
      }
    }
  }

  $shortcut = Set-Shortcut `
    -WorkspaceId $wellKnown['WorkspaceDS'].id `
    -ItemId $wellKnown['Lakehouse'].id `
    -Payload $shortcutPayload

  $wellKnown['Shortcut'] = @{
    shortcutName = $shortcut.name
    shortcutPath = $shortcut.path
    workspaceId  = $wellKnown['WorkspaceDS'].id
    lakehouseId  = $wellKnown['Lakehouse'].id
  }
}

# Create the Folder if not exists
$displayNameTemp = "${displayName}_$($itemNaming['Folder'])"

$folder = Set-FabricFolder `
  -WorkspaceId $wellKnown['WorkspaceDS'].id `
  -DisplayName $displayNameTemp

$wellKnown['Folder'] = @{
  id             = $folder.id
  displayName    = $folder.displayName
  parentFolderId = $folder.parentFolderId
}

# Save wellknown.json file
$wellKnownJson = $wellKnown | ConvertTo-Json -Depth 10
$wellKnownJson
$wellKnownJson | Set-Content -Path './internal/testhelp/fixtures/.wellknown.json' -Force -NoNewline -Encoding utf8
