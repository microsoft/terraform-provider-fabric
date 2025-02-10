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
      Install-Module -Name $ModuleName -AllowClobber -Force -Scope CurrentUser -Repository PSGallery
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
    'DataPipeline' {
      $itemEndpoint = 'dataPipelines'
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

  $definitionRequired = @('Report', 'SemanticModel', 'MirroredDatabase')
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

    $result = (Invoke-FabricRest -Method 'POST' -Endpoint "workspaces/$WorkspaceId/$itemEndpoint" -Payload $payload).Response
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
    $result = (Invoke-FabricRest -Method 'POST' -Endpoint "workspaces/$WorkspaceId/assignToCapacity" -Payload $payload).Response
    $workspace.Response.capacityId = $CapacityId
  }

  return $workspace.Response
}

function Set-FabricWorkspaceRoleAssignment {
  param (
    [Parameter(Mandatory = $true)]
    [string]$WorkspaceId,

    [Parameter(Mandatory = $true)]
    [object]$SPN
  )

  $results = Invoke-FabricRest -Method 'GET' -Endpoint "workspaces/$WorkspaceId/roleAssignments"
  $result = $results.Response.value | Where-Object { $_.id -eq $SPN.Id }
  if (!$result) {
    Write-Log -Message "Assigning SPN to Workspace: $($SPN.DisplayName)" -Level 'WARN'
    $payload = @{
      principal = @{
        id   = $SPN.Id
        type = 'ServicePrincipal'
      }
      role      = 'Admin'
    }
    $result = (Invoke-FabricRest -Method 'POST' -Endpoint "workspaces/$WorkspaceId/roleAssignments" -Payload $payload).Response
  }
}

# Define an array of modules to install
$modules = @('Az.Accounts', 'Az.Resources', 'Az.Fabric', 'pwsh-dotenv', 'ADOPS')

# Loop through each module and install if not installed
foreach ($module in $modules) {
  Install-ModuleIfNotInstalled -ModuleName $module
  Import-ModuleIfNotImported -ModuleName $module
}

# Import the .env file into the environment variables
if (Test-Path -Path './wellknown.env') {
  Import-Dotenv -Path ./wellknown.env -AllowClobber
}

if (!$Env:FABRIC_TESTACC_WELLKNOWN_ENTRA_TENANT_ID -or !$Env:FABRIC_TESTACC_WELLKNOWN_AZURE_SUBSCRIPTION_ID -or !$Env:FABRIC_TESTACC_WELLKNOWN_FABRIC_CAPACITY_NAME -or !$Env:FABRIC_TESTACC_WELLKNOWN_AZDO_ORGANIZATION_NAME -or !$Env:FABRIC_TESTACC_WELLKNOWN_NAME_PREFIX) {
  Write-Log -Message 'FABRIC_TESTACC_WELLKNOWN_ENTRA_TENANT_ID, FABRIC_TESTACC_WELLKNOWN_AZURE_SUBSCRIPTION_ID, FABRIC_TESTACC_WELLKNOWN_FABRIC_CAPACITY_NAME, FABRIC_TESTACC_WELLKNOWN_AZDO_ORGANIZATION_NAME and FABRIC_TESTACC_WELLKNOWN_NAME_PREFIX are required environment variables.' -Level 'ERROR'
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

$SPN = $null
if ($Env:FABRIC_TESTACC_WELLKNOWN_SPN_NAME) {
  $SPN = Get-AzADServicePrincipal -DisplayName $Env:FABRIC_TESTACC_WELLKNOWN_SPN_NAME
}

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
  'Dashboard'             = 'dash'
  'Datamart'              = 'dm'
  'DataPipeline'          = 'dp'
  'Environment'           = 'env'
  'Eventhouse'            = 'eh'
  'Eventstream'           = 'es'
  'KQLDashboard'          = 'kqldash'
  'KQLDatabase'           = 'kqldb'
  'KQLQueryset'           = 'kqlqs'
  'Lakehouse'             = 'lh'
  'MirroredDatabase'      = 'mdb'
  'MirroredWarehouse'     = 'mwh'
  'MLExperiment'          = 'mle'
  'MLModel'               = 'mlm'
  'Notebook'              = 'nb'
  'PaginatedReport'       = 'prpt'
  'Reflex'                = 'rx'
  'Report'                = 'rpt'
  'SemanticModel'         = 'sm'
  'SparkJobDefinition'    = 'sjd'
  'SQLDatabase'           = 'sqldb'
  'SQLEndpoint'           = 'sqle'
  'Warehouse'             = 'wh'
  'WorkspaceDS'           = 'wsds'
  'WorkspaceRS'           = 'wsrs'
  'DomainParent'          = 'parent'
  'DomainChild'           = 'child'
  'EntraServicePrincipal' = 'sp'
  'EntraGroup'            = 'grp'
  'AzDOProject'           = 'proj'
}

$baseName = Get-BaseName
$Env:FABRIC_TESTACC_WELLKNOWN_NAME_BASE = $baseName

# Save env vars wellknown.env file
$envVarNames = @(
  'FABRIC_TESTACC_WELLKNOWN_ENTRA_TENANT_ID',
  'FABRIC_TESTACC_WELLKNOWN_AZURE_SUBSCRIPTION_ID',
  'FABRIC_TESTACC_WELLKNOWN_FABRIC_CAPACITY_NAME',
  'FABRIC_TESTACC_WELLKNOWN_AZDO_ORGANIZATION_NAME',
  'FABRIC_TESTACC_WELLKNOWN_NAME_PREFIX',
  'FABRIC_TESTACC_WELLKNOWN_NAME_SUFFIX',
  'FABRIC_TESTACC_WELLKNOWN_NAME_BASE',
  'FABRIC_TESTACC_WELLKNOWN_SPN_NAME'
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

# Create WorkspaceRS if not exists
$displayNameTemp = "${displayName}_$($itemNaming['WorkspaceRS'])"
$workspace = Set-FabricWorkspace -DisplayName $displayNameTemp -CapacityId $capacity.id

# Assign WorkspaceDS to Capacity if not already assigned or assigned to a different capacity
$workspace = Set-FabricWorkspaceCapacity -WorkspaceId $workspace.id -CapacityId $capacity.id

Write-Log -Message "WorkspaceRS - Name: $($workspace.displayName) / ID: $($workspace.id)"
$wellKnown['WorkspaceRS'] = @{
  id          = $workspace.id
  displayName = $workspace.displayName
  description = $workspace.description
}

# Assign SPN to WorkspaceRS if not already assigned
if ($SPN) {
  Set-FabricWorkspaceRoleAssignment -WorkspaceId $workspace.id -SPN $SPN
}

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

# Assign SPN to WorkspaceRS if not already assigned
if ($SPN) {
  Set-FabricWorkspaceRoleAssignment -WorkspaceId $workspace.id -SPN $SPN
}

# Define an array of item types to create
$itemTypes = @('DataPipeline', 'Environment', 'Eventhouse', 'Eventstream', 'KQLDashboard', 'KQLQueryset', 'Lakehouse', 'MLExperiment', 'MLModel', 'Notebook', 'Reflex', 'SparkJobDefinition', 'Warehouse')

# Loop through each item type and create if not exists
foreach ($itemType in $itemTypes) {

  $displayNameTemp = "${displayName}_$($itemNaming[$itemType])"
  $item = Set-FabricItem -DisplayName $displayNameTemp -WorkspaceId $workspace.id -Type $itemType
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
$kqlDatabase = Set-FabricItem -DisplayName $displayNameTemp -WorkspaceId $workspace.id -Type 'KQLDatabase' -CreationPayload $creationPayload
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
      payload     = Get-DefinitionPartBase64 -Path 'internal/testhelp/fixtures/mirrored_database/mirroring.json'
      payloadType = 'InlineBase64'
    }
  )
}
$mirroredDatabase = Set-FabricItem -DisplayName $displayNameTemp -WorkspaceId $workspace.id -Type 'MirroredDatabase' -Definition $definition
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
$semanticModel = Set-FabricItem -DisplayName $displayNameTemp -WorkspaceId $workspace.id -Type 'SemanticModel' -Definition $definition
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
      path        = 'StaticResources/SharedResources/BaseThemes/CY24SU06.json'
      payload     = Get-DefinitionPartBase64 -Path 'internal/testhelp/fixtures/report_pbir_legacy/StaticResources/SharedResources/BaseThemes/CY24SU06.json'
      payloadType = 'InlineBase64'
    }
  )
}
$report = Set-FabricItem -DisplayName $displayNameTemp -WorkspaceId $workspace.id -Type 'Report' -Definition $definition
$wellKnown['Report'] = @{
  id          = $report.id
  displayName = $report.displayName
  description = $report.description
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

$results = Invoke-FabricRest -Method 'GET' -Endpoint "workspaces/$($workspace.id)/lakehouses/$($wellKnown['Lakehouse']['id'])/tables"
$result = $results.Response.data | Where-Object { $_.name -eq 'publicholidays' }
if (!$result) {
  Write-Log -Message "!!! Please go to the Lakehouse and manually run 'Start with sample data' to populate the data !!!" -Level 'ERROR' -Stop $false
  Write-Log -Message "Lakehouse: https://app.fabric.microsoft.com/groups/$($workspace.id)/lakehouses/$($wellKnown['Lakehouse']['id'])" -Level 'WARN'
}
$wellKnown['Lakehouse']['tableName'] = 'publicholidays'

$displayNameTemp = "${displayName}_$($itemNaming['Dashboard'])"
$results = Invoke-FabricRest -Method 'GET' -Endpoint "workspaces/$($workspace.id)/dashboards"
$result = $results.Response.value | Where-Object { $_.displayName -eq $displayNameTemp }
if (!$result) {
  Write-Log -Message "!!! Please create a Dashboard manually (with Display Name: ${displayNameTemp}), and update details in the well-known file !!!" -Level 'ERROR' -Stop $false
  Write-Log -Message "Workspace: https://app.fabric.microsoft.com/groups/$($workspace.id)" -Level 'WARN'
}
$wellKnown['Dashboard'] = @{
  id          = if ($result) { $result.id } else { '00000000-0000-0000-0000-000000000000' }
  displayName = if ($result) { $result.displayName } else { $displayNameTemp }
  description = if ($result) { $result.description } else { '' }
}

$displayNameTemp = "${displayName}_$($itemNaming['Datamart'])"
$results = Invoke-FabricRest -Method 'GET' -Endpoint "workspaces/$($workspace.id)/datamarts"
$result = $results.Response.value | Where-Object { $_.displayName -eq $displayNameTemp }
if (!$result) {
  Write-Log -Message "!!! Please create a Datamart manually (with Display Name: ${displayNameTemp}), and update details in the well-known file !!!" -Level 'ERROR' -Stop $false
  Write-Log -Message "Workspace: https://app.fabric.microsoft.com/groups/$($workspace.id)" -Level 'WARN'
}
$wellKnown['Datamart'] = @{
  id          = if ($result) { $result.id } else { '00000000-0000-0000-0000-000000000000' }
  displayName = if ($result) { $result.displayName } else { $displayNameTemp }
  description = if ($result) { $result.description } else { '' }
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
  type = 'group'
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

if ($SPN) {
  $body = @{
    originId = $SPN.Id
  }
  $bodyJson = $body | ConvertTo-Json
  $azdoSPN = Invoke-ADOPSRestMethod -Uri "https://vssps.dev.azure.com/$($azdoContext.Organization)/_apis/graph/serviceprincipals?api-version=7.2-preview.1" -Method Post -Body $bodyJson
  $result = Set-ADOPSGitPermission -ProjectId $azdoProject.id -RepositoryId $azdoRepo.id -Descriptor $azdoSPN.descriptor -Allow 'GenericContribute', 'PullRequestContribute', 'CreateBranch', 'CreateTag', 'GenericRead'
}

# Save wellknown.json file
$wellKnownJson = $wellKnown | ConvertTo-Json
$wellKnownJson
$wellKnownJson | Set-Content -Path './internal/testhelp/fixtures/.wellknown.json' -Force -NoNewline -Encoding utf8
