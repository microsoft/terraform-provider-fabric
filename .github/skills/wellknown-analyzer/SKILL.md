---
name: wellknown-analyzer
description: Analyze dependencies and determine what infrastructure must be added to tools/scripts/Set-WellKnown.ps1 and internal/testhelp/wellknown.go so that acceptance tests for a new resource have all required pre-requisite resources available. USE FOR: analyzing well-known test fixture requirements for new resources, determining creation strategies, and generating code snippets for Set-WellKnown.ps1 and wellknown.go additions.
---

# Skill: Well-Known Analyzer

Determine what infrastructure must be added to `tools/scripts/Set-WellKnown.ps1` so that acceptance tests for a new resource have all required pre-requisite resources available.

## Background

Acceptance tests run against real Fabric APIs and need pre-existing infrastructure. The PowerShell script `tools/scripts/Set-WellKnown.ps1` creates all required resources and writes their IDs/metadata to a `.wellknown.json` fixture that tests consume at runtime.

When adding a new resource, you must determine:

1. Can the new item be created simply (no payload, no definition)?
2. Does it need a `CreationPayload` or `Definition`?
3. Does it depend on other Fabric items that must exist first?
4. Does it require Azure infrastructure (Resource Groups, VNets, Storage Accounts, etc.)?
5. Does it need Entra ID objects (service principals, groups)?
6. Does it need external service connections (AzDO, GitHub)?

## Step 1 — Analyze the New Resource's Dependencies

Read the SDK analysis (from `#skill:sdk-contract-navigator`) and Fabric API documentation to answer:

### Fabric Item Dependencies

| Question                                    | How to Check                                 | Example                                                                 |
| ------------------------------------------- | -------------------------------------------- | ----------------------------------------------------------------------- |
| Does Create require a `workspace_id`?       | Almost always yes                            | All workspace-scoped items                                              |
| Does Create require another item's ID?      | Check `CreationPayload` for reference fields | KQL Database needs `parentEventhouseItemId`                             |
| Does Create require a `Definition`?         | Check if API requires definition on create   | Report needs `definition.pbir`, Semantic Model needs `definition.pbism` |
| Does it have a parent item relationship?    | Check if item is scoped under another item   | DigitalTwinBuilderFlow requires a DigitalTwinBuilder                    |
| Does it need data populated after creation? | Check if tests rely on item content          | Lakehouse needs sample data loaded for shortcut tests                   |

### Azure Infrastructure Dependencies

| Question                              | How to Check                                     | Example                            |
| ------------------------------------- | ------------------------------------------------ | ---------------------------------- |
| Does it need an Azure Resource Group? | Check if the resource references Azure resources | Mounted Data Factory, VNet Gateway |
| Does it need a Storage Account?       | Check for blob/storage references                | Managed Private Endpoints          |
| Does it need a Virtual Network?       | Check for VNet/subnet references                 | Virtual Network Gateway            |
| Does it need an Azure Data Factory?   | Check for ADF references                         | Mounted Data Factory               |
| Does it need Azure RBAC assignments?  | Check if Azure roles must be pre-assigned        | Network Contributor on VNet        |

### Entra ID Dependencies

| Question                                 | How to Check                      | Example                    |
| ---------------------------------------- | --------------------------------- | -------------------------- |
| Does it need a Service Principal?        | Check for role assignment tests   | Role assignment resources  |
| Does it need an Entra Group?             | Check for group-based assignments | Workspace role assignments |
| Does it need specific app registrations? | Check for OAuth/auth requirements | Connection resources       |

### External Service Dependencies

| Question                           | How to Check                        | Example                                    |
| ---------------------------------- | ----------------------------------- | ------------------------------------------ |
| Does it need an AzDO project/repo? | Check for Git integration           | Workspace Git                              |
| Does it need a GitHub connection?  | Check for GitHub references         | Workspace Git                              |
| Does it need a Fabric Connection?  | Check for connection references     | Connection role assignments                |
| Does it need a Gateway?            | Check for gateway-scoped operations | Gateway role assignments, VNet connections |

## Step 2 — Check Current Well-Known Script

Read `tools/scripts/Set-WellKnown.ps1` to determine what already exists and what's missing.

### Check the `Set-FabricItem` switch block

Look at the `switch ($Type)` block in the `Set-FabricItem` function. If the new item type is missing, it needs an entry mapping to its REST API endpoint:

```powershell
'<NewItemType>' {
    $itemEndpoint = '<camelCasePluralEndpoint>'
}
```

The endpoint is the camelCase plural of the item type (e.g. `lakehouses`, `eventhouses`, `sqlDatabases`, `dataPipelines`).

### Check the `$itemNaming` hashtable

Every item/resource that gets created needs a short naming abbreviation in `$itemNaming`:

```powershell
'<NewItemType>' = '<2-5 char abbreviation>'
```

### Check existing creation patterns

Determine which creation pattern applies.

## Step 3 — Determine the Creation Strategy

### Strategy A: Simple Item (no payload, no definition)

If the Fabric API can create the item with just `displayName` and `description`, add it to the `$itemTypes` array:

```powershell
$itemTypes = @('ApacheAirflowJob', ..., '<NewItemType>', ..., 'Warehouse')
```

The loop handles creation automatically:

```powershell
foreach ($itemType in $itemTypes) {
    $displayNameTemp = "${displayName}_$($itemNaming[$itemType])"
    $item = Set-FabricItem -DisplayName $displayNameTemp -WorkspaceId $wellKnown['WorkspaceDS'].id -Type $itemType
    $wellKnown[$itemType] = @{
        id          = $item.id
        displayName = $item.displayName
        description = $item.description
    }
}
```

**Items currently using this pattern:** ApacheAirflowJob, CopyJob, Dataflow, DataPipeline, DigitalTwinBuilder, Environment, Eventhouse, GraphQLApi, KQLDashboard, KQLQueryset, Lakehouse, Map, MLExperiment, MLModel, Notebook, Reflex, SparkJobDefinition, SQLDatabase, VariableLibrary, Warehouse

### Strategy B: Item with CreationPayload

If the API requires a `creationPayload` on create, add a dedicated block **after** any dependencies have been created:

```powershell
# Create <NewItemType> if not exists
$displayNameTemp = "${displayName}_$($itemNaming['<NewItemType>'])"
$creationPayload = @{
    <requiredField> = <value>
}
$item = Set-FabricItem -DisplayName $displayNameTemp -WorkspaceId $wellKnown['WorkspaceDS'].id -Type '<NewItemType>' -CreationPayload $creationPayload
$wellKnown['<NewItemType>'] = @{
    id          = $item.id
    displayName = $item.displayName
    description = $item.description
}
```

**Current examples:**

- **KQLDatabase** — needs `databaseType` and `parentEventhouseItemId` (depends on Eventhouse)
- **DigitalTwinBuilderFlow** — needs `digitalTwinBuilderItemReference` (depends on DigitalTwinBuilder)
- **WarehouseSnapshot** — needs `parentWarehouseId` (depends on Warehouse)

### Strategy C: Item with Definition

If the API requires a `definition` on create (items that are definition-required):

```powershell
$displayNameTemp = "${displayName}_$($itemNaming['<NewItemType>'])"
$definition = @{
    parts = @(
        @{
            path        = '<definition-path>'
            payload     = Get-DefinitionPartBase64 -Path 'internal/testhelp/fixtures/<item_type>/<file>.tmpl' -Values @(
                @{ key = '{{ .PlaceholderVar }}'; value = $actualValue }
            )
            payloadType = 'InlineBase64'
        }
    )
}
$item = Set-FabricItem -DisplayName $displayNameTemp -WorkspaceId $wellKnown['WorkspaceDS'].id -Type '<NewItemType>' -Definition $definition
$wellKnown['<NewItemType>'] = @{
    id          = $item.id
    displayName = $item.displayName
    description = $item.description
}
```

**Current examples:**

- **MirroredDatabase** — needs `mirroring.json` definition
- **SemanticModel** — needs `definition.pbism` + `model.bim`
- **Report** — needs `definition.pbir` + `report.json` + static resources (depends on SemanticModel)
- **Eventstream** — needs `eventstream.json` definition (depends on Lakehouse)
- **MountedDataFactory** — needs `mountedDataFactory-content.json` (depends on Azure Data Factory)

This strategy also requires:

1. Create template fixture files in `internal/testhelp/fixtures/<item_type>/`
2. Use `Get-DefinitionPartBase64` to Base64-encode the definition content

### Strategy D: Non-Fabric-Item Infrastructure

For non-Fabric-item resources that need dedicated setup:

**Azure resources:**

```powershell
# Create Azure resource
$resource = Set-Azure<Resource> -ResourceGroupName $wellKnown['ResourceGroup'].name -Name $displayNameTemp ...
$wellKnown['<ResourceKey>'] = @{
    id   = $resource.Id
    name = $resource.Name
}
```

**Fabric connections:**

```powershell
$connection = Set-FabricConnection -DisplayName $displayNameTemp -ConnectivityType "<type>"
$wellKnown['<ConnectionKey>'] = @{
    id          = $connection.id
    displayName = $connection.displayName
}
```

**Role assignments:**

```powershell
Set-FabricGatewayRoleAssignment -GatewayId $gatewayId -PrincipalId $principalId -PrincipalType 'ServicePrincipal' -Role 'Admin'
```

**Domains:**

```powershell
$domain = Set-FabricDomain -DisplayName $displayNameTemp
$wellKnown['<DomainKey>'] = @{ id = $domain.id; displayName = $domain.displayName }
```

## Step 4 — Check for Ordering Requirements

Well-known resources are created sequentially. Ensure dependencies are created before dependents:

1. **Azure infrastructure first** — Resource Groups, VNets, Storage Accounts, Data Factories
2. **Workspaces second** — WorkspaceMPE, WorkspaceOAP, WorkspaceRS, WorkspaceDS
3. **Simple items third** — items in the `$itemTypes` array
4. **Dependent items fourth** — items that reference other items (KQLDatabase→Eventhouse, DigitalTwinBuilderFlow→DigitalTwinBuilder)
5. **Items with definitions fifth** — items requiring definition fixtures
6. **Non-item resources sixth** — Connections, Gateways, Domains, etc.
7. **Sub-resources last** — Role assignments, shortcuts, folders, schedulers

If the new resource depends on something not yet created, note it as a prerequisite.

## Step 5 — Produce the Recommendation

Output a structured recommendation:

### Summary

```
Resource: fabric_<name>
Category: Fabric Item / Non-Item
Creation Strategy: A (simple) / B (payload) / C (definition) / D (non-item infra)
```

### Dependencies Found

List all identified dependencies:

```
✅ Already exists: <dependency>
❌ Missing: <dependency> — needs to be added
```

### Required Changes to `Set-WellKnown.ps1`

1. **Switch block** — Add/confirm entry
2. **$itemNaming** — Add abbreviation
3. **Creation logic** — Add to array or add dedicated block
4. **$wellKnown output** — Ensure entry is written
5. **Fixture files** — List any template files needed in `internal/testhelp/fixtures/`
6. **Azure infra** — List any Azure resources needed
7. **Ordering** — Where in the script the new code should be placed

### Code Snippets

Provide ready-to-paste PowerShell code for each required change.

## Go Access Pattern (`wellknown.go`)

Well-known data is loaded from `internal/testhelp/fixtures/.wellknown.json` (or the `FABRIC_TESTACC_WELLKNOWN` env var) and accessed in tests via:

```go
entity := testhelp.WellKnown()["Lakehouse"].(map[string]any)
entityID := entity["id"].(string)
```

### Key Workspaces

| Key            | Purpose                               |
| -------------- | ------------------------------------- |
| `WorkspaceRS`  | Resource tests (create/update/delete) |
| `WorkspaceDS`  | Data source tests (read-only)         |
| `WorkspaceMPE` | Managed private endpoint tests        |
| `WorkspaceOAP` | Outbound access policy tests          |

## Reference

- Well-known setup script: `tools/scripts/Set-WellKnown.ps1`
- Well-known Go accessor: `internal/testhelp/wellknown.go`
- Test fixture directory: `internal/testhelp/fixtures/`
- Well-known JSON output: `internal/testhelp/fixtures/.wellknown.json`
