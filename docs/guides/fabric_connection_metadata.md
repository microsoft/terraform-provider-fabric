# Fabric Connection Supported Metadata

## `AdlsGen2CosmosStructuredStream` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `ServicePrincipal`
- Supported Connection Encryption Types: `NotEncrypted`

### `AdlsGen2CosmosStructuredStream.Contents` Creation Method

| Name  | Type   | Required | Allowed Values |
| ----- | ------ | -------- | -------------- |
| `url` | `Text` | `True`   | N/A            |

## `AmazonRdsForSqlServer` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `Basic`
- Supported Connection Encryption Types: `NotEncrypted`, `Encrypted`

### `AmazonRdsForSqlServer.Database` Creation Method

| Name       | Type   | Required | Allowed Values |
| ---------- | ------ | -------- | -------------- |
| `server`   | `Text` | `True`   | N/A            |
| `database` | `Text` | `True`   | N/A            |

## `AmazonRedshift` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `Basic`
- Supported Connection Encryption Types: `NotEncrypted`, `Encrypted`

### `AmazonRedshift.Database` Creation Method

| Name           | Type     | Required | Allowed Values |
| -------------- | -------- | -------- | -------------- |
| `server`       | `Text`   | `True`   | N/A            |
| `database`     | `Text`   | `True`   | N/A            |
| `ProviderName` | `Text`   | `False`  | N/A            |
| `BatchSize`    | `Number` | `False`  | N/A            |

## `AmazonS3` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `Basic`
- Supported Connection Encryption Types: `NotEncrypted`

### `AmazonS3.Storage` Creation Method

| Name  | Type   | Required | Allowed Values |
| ----- | ------ | -------- | -------------- |
| `url` | `Text` | `False`  | N/A            |

## `AmazonS3Compatible` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `Basic`
- Supported Connection Encryption Types: `NotEncrypted`

### `AmazonS3Compatible.Storage` Creation Method

| Name  | Type   | Required | Allowed Values |
| ----- | ------ | -------- | -------------- |
| `url` | `Text` | `True`   | N/A            |

## `AnalysisServices` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `Basic`
- Supported Connection Encryption Types: `NotEncrypted`

### `AnalysisServices` Creation Method

| Name       | Type   | Required | Allowed Values |
| ---------- | ------ | -------- | -------------- |
| `server`   | `Text` | `True`   | N/A            |
| `database` | `Text` | `False`  | N/A            |

## `Anaplan` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `Basic`
- Supported Connection Encryption Types: `NotEncrypted`

### `Anaplan.Contents` Creation Method

| Name      | Type   | Required | Allowed Values |
| --------- | ------ | -------- | -------------- |
| `apiUrl`  | `Text` | `True`   | N/A            |
| `authUrl` | `Text` | `True`   | N/A            |

## `ApacheHive` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `Basic`
- Supported Connection Encryption Types: `NotEncrypted`, `Encrypted`

### `ApacheHiveLLAP.Database` Creation Method

| Name              | Type     | Required | Allowed Values |
| ----------------- | -------- | -------- | -------------- |
| `server`          | `Text`   | `True`   | N/A            |
| `database`        | `Text`   | `True`   | N/A            |
| `thriftTransport` | `Number` | `True`   | `1`, `2`       |

## `Applixure` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `Basic`
- Supported Connection Encryption Types: `NotEncrypted`

### `Applixure.EnvironmentName` Creation Method

| Name | Type   | Required | Allowed Values |
| ---- | ------ | -------- | -------------- |
| `id` | `Text` | `True`   | N/A            |

### `Applixure.GetFullObjects` Creation Method

| Name         | Type   | Required | Allowed Values |
| ------------ | ------ | -------- | -------------- |
| `id`         | `Text` | `True`   | N/A            |
| `objectType` | `Text` | `True`   | N/A            |

### `Applixure.GetScoreHistory` Creation Method

| Name       | Type     | Required | Allowed Values |
| ---------- | -------- | -------- | -------------- |
| `id`       | `Text`   | `True`   | N/A            |
| `days`     | `Number` | `True`   | N/A            |
| `interval` | `Number` | `True`   | N/A            |

## `AtScale` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `Basic`
- Supported Connection Encryption Types: `NotEncrypted`

### `AtScale.Cubes` Creation Method

| Name                | Type       | Required | Allowed Values |
| ------------------- | ---------- | -------- | -------------- |
| `server`            | `Text`     | `True`   | N/A            |
| `ConnectionTimeout` | `Duration` | `False`  | N/A            |
| `CommandTimeout`    | `Duration` | `False`  | N/A            |

## `AutomationAnywhere` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `Basic`
- Supported Connection Encryption Types: `NotEncrypted`

### `AutomationAnywhere.Feed` Creation Method

| Name         | Type   | Required | Allowed Values                                      |
| ------------ | ------ | -------- | --------------------------------------------------- |
| `CRVersion`  | `Text` | `True`   | `10.x/11.x`, `Automation 360`, `11.3.5.1 Or Higher` |
| `CRHostName` | `Text` | `True`   | N/A                                                 |

## `AutomyDataAnalytics` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `Key`
- Supported Connection Encryption Types: `NotEncrypted`

## `AzureAISearch` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `Key`
- Supported Connection Encryption Types: `NotEncrypted`

### `AzureAISearch.Contents` Creation Method

| Name  | Type   | Required | Allowed Values |
| ----- | ------ | -------- | -------------- |
| `url` | `Text` | `True`   | N/A            |

## `AzureBatch` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `Key`
- Supported Connection Encryption Types: `NotEncrypted`

### `AzureBatch.Contents` Creation Method

| Name          | Type   | Required | Allowed Values |
| ------------- | ------ | -------- | -------------- |
| `accountName` | `Text` | `True`   | N/A            |
| `batchUrl`    | `Text` | `True`   | N/A            |
| `poolName`    | `Text` | `True`   | N/A            |

## `AzureBlobs` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `Anonymous`, `Key`, `SharedAccessSignature`, `ServicePrincipal`
- Supported Connection Encryption Types: `NotEncrypted`

### `AzureBlobs` Creation Method

| Name      | Type   | Required | Allowed Values |
| --------- | ------ | -------- | -------------- |
| `account` | `Text` | `True`   | N/A            |
| `domain`  | `Text` | `True`   | N/A            |

## `AzureCosmosDBForMongoDB` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `Basic`
- Supported Connection Encryption Types: `NotEncrypted`, `Encrypted`

### `AzureCosmosDBForMongoDB.Database` Creation Method

| Name            | Type   | Required | Allowed Values     |
| --------------- | ------ | -------- | ------------------ |
| `server`        | `Text` | `True`   | N/A                |
| `serverVersion` | `Text` | `True`   | `Above 3.2`, `3.2` |

## `AzureDatabaseForMySQL` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `Basic`
- Supported Connection Encryption Types: `NotEncrypted`, `Encrypted`

### `AzureDatabaseForMySQL.Database` Creation Method

| Name       | Type   | Required | Allowed Values |
| ---------- | ------ | -------- | -------------- |
| `server`   | `Text` | `True`   | N/A            |
| `database` | `Text` | `True`   | N/A            |

## `AzureDatabricksWorkspace` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `Key`, `ServicePrincipal`
- Supported Connection Encryption Types: `NotEncrypted`

### `AzureDatabricksWorkspace.Actions` Creation Method

| Name  | Type   | Required | Allowed Values |
| ----- | ------ | -------- | -------------- |
| `url` | `Text` | `True`   | N/A            |

## `AzureDataFactory` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `ServicePrincipal`
- Supported Connection Encryption Types: `NotEncrypted`

### `AzureDataFactory.Actions` Creation Method

| Name              | Type   | Required | Allowed Values |
| ----------------- | ------ | -------- | -------------- |
| `subscriptionId`  | `Text` | `True`   | N/A            |
| `resourceGroup`   | `Text` | `True`   | N/A            |
| `dataFactoryName` | `Text` | `True`   | N/A            |

## `AzureDataLakeStorage` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `Key`, `SharedAccessSignature`, `ServicePrincipal`, `WorkspaceIdentity`
- Supported Connection Encryption Types: `NotEncrypted`

### `AzureDataLakeStorage` Creation Method

| Name     | Type   | Required | Allowed Values |
| -------- | ------ | -------- | -------------- |
| `server` | `Text` | `True`   | N/A            |
| `path`   | `Text` | `True`   | N/A            |

## `AzureDataLakeStoreCosmosStructuredStream` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `ServicePrincipal`
- Supported Connection Encryption Types: `NotEncrypted`

### `AzureDataLakeStoreCosmosStructuredStream.Contents` Creation Method

| Name  | Type   | Required | Allowed Values |
| ----- | ------ | -------- | -------------- |
| `url` | `Text` | `True`   | N/A            |

## `AzureEnterprise` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `Key`
- Supported Connection Encryption Types: `NotEncrypted`

### `AzureEnterprise.Contents` Creation Method

| Name  | Type   | Required | Allowed Values |
| ----- | ------ | -------- | -------------- |
| `url` | `Text` | `True`   | N/A            |

### `AzureEnterprise.Tables` Creation Method

| Name  | Type   | Required | Allowed Values |
| ----- | ------ | -------- | -------------- |
| `url` | `Text` | `True`   | N/A            |

## `AzureFiles` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `Key`
- Supported Connection Encryption Types: `NotEncrypted`

### `AzureFiles.Contents` Creation Method

| Name       | Type   | Required | Allowed Values |
| ---------- | ------ | -------- | -------------- |
| `shareUrl` | `Text` | `True`   | N/A            |
| `snapshot` | `Text` | `False`  | N/A            |

## `AzureFunction` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `Key`, `Anonymous`
- Supported Connection Encryption Types: `NotEncrypted`

### `AzureFunction.Contents` Creation Method

| Name             | Type   | Required | Allowed Values |
| ---------------- | ------ | -------- | -------------- |
| `functionAppUrl` | `Text` | `True`   | N/A            |

## `AzureHDInsightCluster` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `Basic`
- Supported Connection Encryption Types: `NotEncrypted`

### `AzureHDInsightCluster.Actions` Creation Method

| Name                   | Type      | Required | Allowed Values |
| ---------------------- | --------- | -------- | -------------- |
| `hdiUrl`               | `Text`    | `True`   | N/A            |
| `entSecPackageEnabled` | `Boolean` | `False`  | N/A            |

## `AzureHDInsightOnDemandCluster` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `ServicePrincipal`
- Supported Connection Encryption Types: `NotEncrypted`

### `AzureHDInsightOnDemandCluster.Actions` Creation Method

| Name                | Type   | Required | Allowed Values |
| ------------------- | ------ | -------- | -------------- |
| `subscriptionId`    | `Text` | `True`   | N/A            |
| `resourceGroupName` | `Text` | `True`   | N/A            |

## `AzureHive` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `Basic`
- Supported Connection Encryption Types: `NotEncrypted`, `Encrypted`

### `AzureHiveLLAP.Database` Creation Method

| Name       | Type   | Required | Allowed Values |
| ---------- | ------ | -------- | -------------- |
| `server`   | `Text` | `True`   | N/A            |
| `database` | `Text` | `True`   | N/A            |

## `AzureKeyVault` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `ServicePrincipal`
- Supported Connection Encryption Types: `NotEncrypted`

### `AzureKeyVault.Actions` Creation Method

| Name          | Type   | Required | Allowed Values |
| ------------- | ------ | -------- | -------------- |
| `accountName` | `Text` | `True`   | N/A            |

## `AzureMachineLearning` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `ServicePrincipal`
- Supported Connection Encryption Types: `NotEncrypted`

### `AzureMachineLearning.Contents` Creation Method

| Name                | Type   | Required | Allowed Values |
| ------------------- | ------ | -------- | -------------- |
| `subscriptionId`    | `Text` | `True`   | N/A            |
| `resourceGroupName` | `Text` | `True`   | N/A            |
| `workspaceName`     | `Text` | `True`   | N/A            |

## `AzurePostgreSQL` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `Basic`
- Supported Connection Encryption Types: `NotEncrypted`, `Encrypted`

### `AzurePostgreSQL.Database` Creation Method

| Name       | Type   | Required | Allowed Values |
| ---------- | ------ | -------- | -------------- |
| `server`   | `Text` | `True`   | N/A            |
| `database` | `Text` | `True`   | N/A            |

## `AzureServiceBus` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `Basic`
- Supported Connection Encryption Types: `NotEncrypted`

### `AzureServiceBus.Contents` Creation Method

| Name       | Type   | Required | Allowed Values |
| ---------- | ------ | -------- | -------------- |
| `hostName` | `Text` | `True`   | N/A            |

## `AzureSqlMI` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `Basic`, `ServicePrincipal`
- Supported Connection Encryption Types: `NotEncrypted`, `Encrypted`

### `AzureSqlMI.Database` Creation Method

| Name       | Type   | Required | Allowed Values |
| ---------- | ------ | -------- | -------------- |
| `server`   | `Text` | `True`   | N/A            |
| `database` | `Text` | `True`   | N/A            |
| `Query`    | `Text` | `False`  | N/A            |

## `AzureSynapseWorkspace` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `ServicePrincipal`
- Supported Connection Encryption Types: `NotEncrypted`

### `AzureSynapseWorkspace.Actions` Creation Method

| Name            | Type   | Required | Allowed Values |
| --------------- | ------ | -------- | -------------- |
| `workspaceName` | `Text` | `True`   | N/A            |

## `AzureTables` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `Key`
- Supported Connection Encryption Types: `NotEncrypted`

### `AzureTables` Creation Method

| Name      | Type   | Required | Allowed Values |
| --------- | ------ | -------- | -------------- |
| `account` | `Text` | `True`   | N/A            |
| `domain`  | `Text` | `True`   | N/A            |

## `BI360` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `Key`
- Supported Connection Encryption Types: `NotEncrypted`

### `BI360.Contents` Creation Method

| Name  | Type   | Required | Allowed Values |
| ----- | ------ | -------- | -------------- |
| `Url` | `Text` | `True`   | N/A            |

## `BitSightSecurityRatings` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `Key`
- Supported Connection Encryption Types: `NotEncrypted`

### `BitSightSecurityRatings.Contents` Creation Method

| Name                      | Type      | Required | Allowed Values |
| ------------------------- | --------- | -------- | -------------- |
| `company_guid`            | `Text`    | `False`  | N/A            |
| `affects_rating_findings` | `Boolean` | `False`  | N/A            |

## `CCHTagetik` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `Basic`
- Supported Connection Encryption Types: `NotEncrypted`

### `CCHTagetik.Contents` Creation Method

| Name       | Type   | Required | Allowed Values |
| ---------- | ------ | -------- | -------------- |
| `URL`      | `Text` | `True`   | N/A            |
| `Database` | `Text` | `True`   | N/A            |
| `AW`       | `Text` | `False`  | N/A            |
| `Dataset`  | `Text` | `False`  | N/A            |

## `Celonis` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `Key`
- Supported Connection Encryption Types: `NotEncrypted`

### `Celonis.Navigation` Creation Method

| Name  | Type   | Required | Allowed Values |
| ----- | ------ | -------- | -------------- |
| `URL` | `Text` | `True`   | N/A            |

### `Celonis.KnowledgeModels` Creation Method

| Name  | Type   | Required | Allowed Values |
| ----- | ------ | -------- | -------------- |
| `URL` | `Text` | `True`   | N/A            |

## `CloudBluePSA` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `Key`
- Supported Connection Encryption Types: `NotEncrypted`

### `CloudBluePSA.Feed` Creation Method

| Name     | Type   | Required | Allowed Values |
| -------- | ------ | -------- | -------------- |
| `url`    | `Text` | `True`   | N/A            |
| `filter` | `Text` | `True`   | N/A            |

## `CloudScope` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `Key`
- Supported Connection Encryption Types: `NotEncrypted`

## `CloudScopeInstagram` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `Key`
- Supported Connection Encryption Types: `NotEncrypted`

## `Cognite` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `Key`
- Supported Connection Encryption Types: `NotEncrypted`

### `Cognite.Contents` Creation Method

| Name          | Type   | Required | Allowed Values |
| ------------- | ------ | -------- | -------------- |
| `project`     | `Text` | `True`   | N/A            |
| `environment` | `Text` | `False`  | N/A            |

## `CommonDataService` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `ServicePrincipal`
- Supported Connection Encryption Types: `NotEncrypted`, `Encrypted`

### `CommonDataService.Database` Creation Method

| Name                         | Type      | Required | Allowed Values |
| ---------------------------- | --------- | -------- | -------------- |
| `server`                     | `Text`    | `False`  | N/A            |
| `CreateNavigationProperties` | `Boolean` | `False`  | N/A            |

## `comScore` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `Basic`
- Supported Connection Encryption Types: `NotEncrypted`

### `comScore.GetReport` Creation Method

| Name            | Type   | Required | Allowed Values |
| --------------- | ------ | -------- | -------------- |
| `datacenter`    | `Text` | `True`   | N/A            |
| `client`        | `Text` | `True`   | N/A            |
| `itemId`        | `Text` | `True`   | N/A            |
| `site`          | `Text` | `True`   | N/A            |
| `startDate`     | `Date` | `False`  | N/A            |
| `endDate`       | `Date` | `False`  | N/A            |
| `SegmentId`     | `Text` | `False`  | N/A            |
| `VisitFilterId` | `Text` | `False`  | N/A            |
| `EventFilterId` | `Text` | `False`  | N/A            |
| `fullUrlString` | `Text` | `False`  | N/A            |

### `comScore.ReportItems` Creation Method

| Name            | Type   | Required | Allowed Values |
| --------------- | ------ | -------- | -------------- |
| `datacenter`    | `Text` | `True`   | N/A            |
| `client`        | `Text` | `True`   | N/A            |
| `itemId`        | `Text` | `True`   | N/A            |
| `site`          | `Text` | `True`   | N/A            |
| `startDate`     | `Date` | `False`  | N/A            |
| `endDate`       | `Date` | `False`  | N/A            |
| `SegmentId`     | `Text` | `False`  | N/A            |
| `VisitFilterId` | `Text` | `False`  | N/A            |
| `EventFilterId` | `Text` | `False`  | N/A            |
| `fullUrlString` | `Text` | `False`  | N/A            |

### `comScore.NavTable` Creation Method

| Name         | Type   | Required | Allowed Values |
| ------------ | ------ | -------- | -------------- |
| `datacenter` | `Text` | `True`   | N/A            |
| `client`     | `Text` | `True`   | N/A            |
| `startDate`  | `Date` | `False`  | N/A            |
| `endDate`    | `Date` | `False`  | N/A            |

## `ConfluentCloud` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `Basic`
- Supported Connection Encryption Types: `NotEncrypted`

### `ConfluentCloud.Contents` Creation Method

| Name     | Type   | Required | Allowed Values |
| -------- | ------ | -------- | -------------- |
| `server` | `Text` | `True`   | N/A            |

## `CosmosDB` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `Key`
- Supported Connection Encryption Types: `NotEncrypted`

### `CosmosDB.Contents` Creation Method

| Name                                        | Type   | Required | Allowed Values |
| ------------------------------------------- | ------ | -------- | -------------- |
| `host`                                      | `Text` | `True`   | N/A            |
| `NUMBER_OF_RETRIES`                         | `Text` | `False`  | N/A            |
| `ENABLE_AVERAGE_FUNCTION_PASSDOWN`          | `Text` | `False`  | N/A            |
| `ENABLE_SORT_PASSDOWN_FOR_MULTIPLE_COLUMNS` | `Text` | `False`  | N/A            |

## `Databricks` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `Basic`, `Key`
- Supported Connection Encryption Types: `NotEncrypted`

### `Databricks.Catalogs` Creation Method

| Name                            | Type   | Required | Allowed Values        |
| ------------------------------- | ------ | -------- | --------------------- |
| `host`                          | `Text` | `True`   | N/A                   |
| `httpPath`                      | `Text` | `True`   | N/A                   |
| `Catalog`                       | `Text` | `False`  | N/A                   |
| `Database`                      | `Text` | `False`  | N/A                   |
| `EnableAutomaticProxyDiscovery` | `Text` | `False`  | `enabled`, `disabled` |

### `Databricks.Contents` Creation Method

| Name                            | Type   | Required | Allowed Values        |
| ------------------------------- | ------ | -------- | --------------------- |
| `host`                          | `Text` | `True`   | N/A                   |
| `httpPath`                      | `Text` | `True`   | N/A                   |
| `Catalog`                       | `Text` | `False`  | N/A                   |
| `Database`                      | `Text` | `False`  | N/A                   |
| `EnableAutomaticProxyDiscovery` | `Text` | `False`  | `enabled`, `disabled` |

### `Databricks.Query` Creation Method

| Name       | Type   | Required | Allowed Values |
| ---------- | ------ | -------- | -------------- |
| `host`     | `Text` | `True`   | N/A            |
| `httpPath` | `Text` | `True`   | N/A            |

## `DatabricksMultiCloud` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `Basic`, `Key`
- Supported Connection Encryption Types: `NotEncrypted`

### `DatabricksMultiCloud.Catalogs` Creation Method

| Name                            | Type   | Required | Allowed Values        |
| ------------------------------- | ------ | -------- | --------------------- |
| `host`                          | `Text` | `True`   | N/A                   |
| `httpPath`                      | `Text` | `True`   | N/A                   |
| `Catalog`                       | `Text` | `False`  | N/A                   |
| `Database`                      | `Text` | `False`  | N/A                   |
| `EnableAutomaticProxyDiscovery` | `Text` | `False`  | `enabled`, `disabled` |

### `DatabricksMultiCloud.Query` Creation Method

| Name       | Type   | Required | Allowed Values |
| ---------- | ------ | -------- | -------------- |
| `host`     | `Text` | `True`   | N/A            |
| `httpPath` | `Text` | `True`   | N/A            |

## `DataLake` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `ServicePrincipal`
- Supported Connection Encryption Types: `NotEncrypted`

### `DataLake.Contents` Creation Method

| Name       | Type     | Required | Allowed Values |
| ---------- | -------- | -------- | -------------- |
| `url`      | `Text`   | `True`   | N/A            |
| `PageSize` | `Number` | `False`  | N/A            |

### `DataLake.Files` Creation Method

| Name  | Type   | Required | Allowed Values |
| ----- | ------ | -------- | -------------- |
| `url` | `Text` | `True`   | N/A            |

## `DataLakeAnalytics` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `ServicePrincipal`
- Supported Connection Encryption Types: `NotEncrypted`

### `DataLakeAnalytics.Account` Creation Method

| Name                | Type   | Required | Allowed Values |
| ------------------- | ------ | -------- | -------------- |
| `accountName`       | `Text` | `True`   | N/A            |
| `subscriptionId`    | `Text` | `True`   | N/A            |
| `resourceGroupName` | `Text` | `False`  | N/A            |

## `DataPipelineCosmosDb` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `Key`
- Supported Connection Encryption Types: `NotEncrypted`

### `DataPipelineCosmosDb.Contents` Creation Method

| Name              | Type   | Required | Allowed Values |
| ----------------- | ------ | -------- | -------------- |
| `AccountEndpoint` | `Text` | `True`   | N/A            |
| `Database`        | `Text` | `True`   | N/A            |

## `DataWorld` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `Key`
- Supported Connection Encryption Types: `NotEncrypted`

### `DataWorld.Dataset` Creation Method

| Name    | Type   | Required | Allowed Values |
| ------- | ------ | -------- | -------------- |
| `owner` | `Text` | `True`   | N/A            |
| `id`    | `Text` | `True`   | N/A            |
| `query` | `Text` | `False`  | N/A            |

## `DeltaSharing` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `Key`
- Supported Connection Encryption Types: `NotEncrypted`

### `DeltaSharing.Contents` Creation Method

| Name           | Type     | Required | Allowed Values |
| -------------- | -------- | -------- | -------------- |
| `host`         | `Text`   | `True`   | N/A            |
| `rowLimitHint` | `Number` | `False`  | N/A            |

## `DocumentDB` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `Key`
- Supported Connection Encryption Types: `NotEncrypted`

### `DocumentDB.Contents` Creation Method

| Name         | Type   | Required | Allowed Values |
| ------------ | ------ | -------- | -------------- |
| `url`        | `Text` | `True`   | N/A            |
| `database`   | `Text` | `False`  | N/A            |
| `collection` | `Text` | `False`  | N/A            |

## `Dremio` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `Basic`
- Supported Connection Encryption Types: `NotEncrypted`

### `Dremio.Databases` Creation Method

| Name           | Type   | Required | Allowed Values |
| -------------- | ------ | -------- | -------------- |
| `server`       | `Text` | `True`   | N/A            |
| `engine`       | `Text` | `False`  | N/A            |
| `routingTag`   | `Text` | `False`  | N/A            |
| `routingQueue` | `Text` | `False`  | N/A            |

### `Dremio.DatabasesV300` Creation Method

| Name           | Type   | Required | Allowed Values                       |
| -------------- | ------ | -------- | ------------------------------------ |
| `server`       | `Text` | `True`   | N/A                                  |
| `encryption`   | `Text` | `True`   | `Enabled`, `Disabled`, `Enabled-PEM` |
| `engine`       | `Text` | `False`  | N/A                                  |
| `routingTag`   | `Text` | `False`  | N/A                                  |
| `routingQueue` | `Text` | `False`  | N/A                                  |

### `Dremio.DatabasesV370` Creation Method

| Name           | Type   | Required | Allowed Values                       |
| -------------- | ------ | -------- | ------------------------------------ |
| `server`       | `Text` | `True`   | N/A                                  |
| `encryption`   | `Text` | `True`   | `Enabled`, `Disabled`, `Enabled-PEM` |
| `engine`       | `Text` | `False`  | N/A                                  |
| `routingTag`   | `Text` | `False`  | N/A                                  |
| `routingQueue` | `Text` | `False`  | N/A                                  |

## `DremioCloud` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `Key`
- Supported Connection Encryption Types: `NotEncrypted`

### `DremioCloud.Databases` Creation Method

| Name        | Type   | Required | Allowed Values |
| ----------- | ------ | -------- | -------------- |
| `projectId` | `Text` | `True`   | N/A            |
| `engine`    | `Text` | `False`  | N/A            |

### `DremioCloud.DatabasesByServer` Creation Method

| Name        | Type   | Required | Allowed Values |
| ----------- | ------ | -------- | -------------- |
| `server`    | `Text` | `True`   | N/A            |
| `projectId` | `Text` | `False`  | N/A            |
| `engine`    | `Text` | `False`  | N/A            |

### `DremioCloud.DatabasesByServerV330` Creation Method

| Name           | Type   | Required | Allowed Values |
| -------------- | ------ | -------- | -------------- |
| `server`       | `Text` | `True`   | N/A            |
| `projectId`    | `Text` | `False`  | N/A            |
| `engine`       | `Text` | `False`  | N/A            |
| `routingTag`   | `Text` | `False`  | N/A            |
| `routingQueue` | `Text` | `False`  | N/A            |

### `DremioCloud.DatabasesByServerV360` Creation Method

| Name           | Type   | Required | Allowed Values |
| -------------- | ------ | -------- | -------------- |
| `server`       | `Text` | `True`   | N/A            |
| `projectId`    | `Text` | `False`  | N/A            |
| `engine`       | `Text` | `False`  | N/A            |
| `routingTag`   | `Text` | `False`  | N/A            |
| `routingQueue` | `Text` | `False`  | N/A            |
| `encryption`   | `Text` | `False`  | `Enabled-PEM`  |

### `DremioCloud.DatabasesByServerV370` Creation Method

| Name           | Type   | Required | Allowed Values |
| -------------- | ------ | -------- | -------------- |
| `server`       | `Text` | `True`   | N/A            |
| `projectId`    | `Text` | `False`  | N/A            |
| `engine`       | `Text` | `False`  | N/A            |
| `routingTag`   | `Text` | `False`  | N/A            |
| `routingQueue` | `Text` | `False`  | N/A            |
| `encryption`   | `Text` | `False`  | `Enabled-PEM`  |

## `Dynamics 365 Business Central (on-premises)` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `Basic`
- Supported Connection Encryption Types: `NotEncrypted`

### `Dynamics365BusinessCentralOnPremises.Contents` Creation Method

| Name      | Type   | Required | Allowed Values |
| --------- | ------ | -------- | -------------- |
| `url`     | `Text` | `True`   | N/A            |
| `company` | `Text` | `False`  | N/A            |

## `Dynamics NAV` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `Basic`
- Supported Connection Encryption Types: `NotEncrypted`

### `DynamicsNav.Contents` Creation Method

| Name      | Type   | Required | Allowed Values |
| --------- | ------ | -------- | -------------- |
| `url`     | `Text` | `True`   | N/A            |
| `company` | `Text` | `False`  | N/A            |

## `Dynamics365` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `ServicePrincipal`
- Supported Connection Encryption Types: `NotEncrypted`

### `Dynamics365.Contents` Creation Method

| Name     | Type   | Required | Allowed Values |
| -------- | ------ | -------- | -------------- |
| `server` | `Text` | `True`   | N/A            |

## `DynamicsAX` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `ServicePrincipal`
- Supported Connection Encryption Types: `NotEncrypted`

### `DynamicsAX.Contents` Creation Method

| Name  | Type   | Required | Allowed Values |
| ----- | ------ | -------- | -------------- |
| `url` | `Text` | `True`   | N/A            |

## `DynamicsCrm` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `ServicePrincipal`
- Supported Connection Encryption Types: `NotEncrypted`

### `DynamicsCrm.Contents` Creation Method

| Name     | Type   | Required | Allowed Values |
| -------- | ------ | -------- | -------------- |
| `server` | `Text` | `True`   | N/A            |

## `EduFrame` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `Key`
- Supported Connection Encryption Types: `NotEncrypted`

### `EduFrame.Contents` Creation Method

| Name         | Type   | Required | Allowed Values |
| ------------ | ------ | -------- | -------------- |
| `domainSlug` | `Text` | `True`   | N/A            |

## `ElasticSearch` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `Basic`, `Anonymous`
- Supported Connection Encryption Types: `NotEncrypted`

### `ElasticSearch.Database` Creation Method

| Name            | Type   | Required | Allowed Values |
| --------------- | ------ | -------- | -------------- |
| `url`           | `Text` | `True`   | N/A            |
| `keyColumnName` | `Text` | `False`  | N/A            |

## `EmigoDataSourceConnector` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `Basic`
- Supported Connection Encryption Types: `NotEncrypted`

### `Emigo.Contents` Creation Method

| Name                   | Type   | Required | Allowed Values                                                     |
| ---------------------- | ------ | -------- | ------------------------------------------------------------------ |
| `DataRestrictionType`  | `Text` | `False`  | `Not set`, `Days`, `Weeks`, `Months`, `Quarters`, `Years`          |
| `DataRestrictionValue` | `Text` | `False`  | N/A                                                                |
| `DataRestrictionMode`  | `Text` | `False`  | `Default`, `Exact`                                                 |
| `AuthorizationMode`    | `Text` | `False`  | `Default`, `EmigoObszary`, `EmigoHierarchia`, `CustomRestrictions` |

## `EQuIS` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `Key`, `Basic`
- Supported Connection Encryption Types: `NotEncrypted`

### `EQuIS.Contents` Creation Method

| Name      | Type   | Required | Allowed Values |
| --------- | ------ | -------- | -------------- |
| `baseUri` | `Text` | `True`   | N/A            |

## `EventHub` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `Basic`
- Supported Connection Encryption Types: `NotEncrypted`

### `EventHub.Contents` Creation Method

| Name         | Type   | Required | Allowed Values |
| ------------ | ------ | -------- | -------------- |
| `endpoint`   | `Text` | `True`   | N/A            |
| `entityPath` | `Text` | `True`   | N/A            |

## `FactSetAnalytics` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `Basic`
- Supported Connection Encryption Types: `NotEncrypted`

## `FactSetRMS` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `Basic`
- Supported Connection Encryption Types: `NotEncrypted`

## `Fhir` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `Anonymous`
- Supported Connection Encryption Types: `NotEncrypted`

### `Fhir.Contents` Creation Method

| Name          | Type   | Required | Allowed Values |
| ------------- | ------ | -------- | -------------- |
| `url`         | `Text` | `True`   | N/A            |
| `searchQuery` | `Text` | `False`  | N/A            |

## `FTP` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `Basic`, `Anonymous`
- Supported Connection Encryption Types: `NotEncrypted`

### `FTP.Contents` Creation Method

| Name     | Type   | Required | Allowed Values |
| -------- | ------ | -------- | -------------- |
| `server` | `Text` | `True`   | N/A            |

## `GitHubSourceControl` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `Key`
- Supported Connection Encryption Types: `NotEncrypted`

### `GitHubSourceControl.Contents` Creation Method

| Name  | Type   | Required | Allowed Values |
| ----- | ------ | -------- | -------------- |
| `url` | `Text` | `False`  | N/A            |

## `GoogleAds` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `Basic`
- Supported Connection Encryption Types: `NotEncrypted`

### `GoogleAds.Contents` Creation Method

| Name                  | Type   | Required | Allowed Values |
| --------------------- | ------ | -------- | -------------- |
| `googleAdsApiVersion` | `Text` | `True`   | N/A            |
| `clientCustomerID`    | `Text` | `True`   | N/A            |
| `loginCustomerID`     | `Text` | `False`  | N/A            |

## `GoogleBigQuery` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `Basic`
- Supported Connection Encryption Types: `NotEncrypted`, `Encrypted`

### `GoogleBigQuery.Database` Creation Method

| Name                | Type       | Required | Allowed Values |
| ------------------- | ---------- | -------- | -------------- |
| `BillingProject`    | `Text`     | `False`  | N/A            |
| `UseStorageApi`     | `Boolean`  | `False`  | N/A            |
| `ConnectionTimeout` | `Duration` | `False`  | N/A            |
| `CommandTimeout`    | `Duration` | `False`  | N/A            |
| `ProjectId`         | `Text`     | `False`  | N/A            |

## `GoogleCloudStorage` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `Basic`
- Supported Connection Encryption Types: `NotEncrypted`

### `GoogleCloudStorage.Storage` Creation Method

| Name  | Type   | Required | Allowed Values |
| ----- | ------ | -------- | -------------- |
| `url` | `Text` | `False`  | N/A            |

## `GooglePubSub` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `Basic`
- Supported Connection Encryption Types: `NotEncrypted`

### `GooglePubSub.Contents` Creation Method

| Name        | Type   | Required | Allowed Values |
| ----------- | ------ | -------- | -------------- |
| `projectId` | `Text` | `True`   | N/A            |

## `HttpServer` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `Basic`, `Anonymous`
- Supported Connection Encryption Types: `NotEncrypted`

### `HttpServer.Contents` Creation Method

| Name  | Type   | Required | Allowed Values |
| ----- | ------ | -------- | -------------- |
| `url` | `Text` | `True`   | N/A            |

## `Impala` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `Anonymous`, `Basic`
- Supported Connection Encryption Types: `NotEncrypted`, `Encrypted`

### `Impala.Database` Creation Method

| Name                | Type       | Required | Allowed Values |
| ------------------- | ---------- | -------- | -------------- |
| `server`            | `Text`     | `True`   | N/A            |
| `ConnectionTimeout` | `Duration` | `False`  | N/A            |
| `CommandTimeout`    | `Duration` | `False`  | N/A            |

## `InformationGrid` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `Key`
- Supported Connection Encryption Types: `NotEncrypted`

### `InformationGrid.Contents` Creation Method

| Name     | Type   | Required | Allowed Values |
| -------- | ------ | -------- | -------------- |
| `server` | `Text` | `True`   | N/A            |

## `Insightly` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `Key`
- Supported Connection Encryption Types: `NotEncrypted`

### `Insightly.PagedTable` Creation Method

| Name  | Type   | Required | Allowed Values |
| ----- | ------ | -------- | -------------- |
| `url` | `Text` | `True`   | N/A            |

## `IoTHub` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `Basic`
- Supported Connection Encryption Types: `NotEncrypted`

### `IoTHub.Contents` Creation Method

| Name         | Type   | Required | Allowed Values |
| ------------ | ------ | -------- | -------------- |
| `entityPath` | `Text` | `True`   | N/A            |

## `JamfPro` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `Basic`, `Key`
- Supported Connection Encryption Types: `NotEncrypted`

### `JamfPro.Contents` Creation Method

| Name      | Type   | Required | Allowed Values |
| --------- | ------ | -------- | -------------- |
| `jamfUrl` | `Text` | `True`   | N/A            |

## `JDIConnector` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `Basic`
- Supported Connection Encryption Types: `NotEncrypted`

### `JDIConnector.Contents` Creation Method

| Name      | Type   | Required | Allowed Values |
| --------- | ------ | -------- | -------------- |
| `DataUrl` | `Text` | `True`   | N/A            |

## `Kinesis` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `Basic`
- Supported Connection Encryption Types: `NotEncrypted`

### `Kinesis.Contents` Creation Method

| Name         | Type   | Required | Allowed Values |
| ------------ | ------ | -------- | -------------- |
| `streamName` | `Text` | `True`   | N/A            |

## `LinkedInLearning` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `Basic`
- Supported Connection Encryption Types: `NotEncrypted`

### `LinkedInLearning.Contents` Creation Method

| Name         | Type       | Required | Allowed Values |
| ------------ | ---------- | -------- | -------------- |
| `start_date` | `DateTime` | `False`  | N/A            |
| `end_date`   | `DateTime` | `False`  | N/A            |

## `Mandrill` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `Key`
- Supported Connection Encryption Types: `NotEncrypted`

### `Mandrill.Contents` Creation Method

| Name  | Type   | Required | Allowed Values |
| ----- | ------ | -------- | -------------- |
| `url` | `Text` | `True`   | N/A            |

## `MariaDBForPipeline` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `Basic`
- Supported Connection Encryption Types: `NotEncrypted`, `Encrypted`

### `MariaDBForPipeline.Database` Creation Method

| Name       | Type   | Required | Allowed Values |
| ---------- | ------ | -------- | -------------- |
| `server`   | `Text` | `True`   | N/A            |
| `database` | `Text` | `True`   | N/A            |

## `Microsoft365` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `ServicePrincipal`
- Supported Connection Encryption Types: `NotEncrypted`

## `MicroStrategyDataset` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `Basic`
- Supported Connection Encryption Types: `NotEncrypted`

### `MicroStrategyDataset.Contents` Creation Method

| Name         | Type     | Required | Allowed Values     |
| ------------ | -------- | -------- | ------------------ |
| `libraryUrl` | `Text`   | `True`   | N/A                |
| `authMode`   | `Text`   | `False`  | `Standard`, `LDAP` |
| `limit`      | `Number` | `False`  | N/A                |
| `timeout`    | `Number` | `False`  | N/A                |

### `MicroStrategyDataset.TestConnection` Creation Method

| Name         | Type   | Required | Allowed Values |
| ------------ | ------ | -------- | -------------- |
| `libraryUrl` | `Text` | `True`   | N/A            |

## `MongoDBAtlasForPipeline` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `Basic`
- Supported Connection Encryption Types: `NotEncrypted`

### `MongoDBAtlasForPipeline.Database` Creation Method

| Name           | Type   | Required | Allowed Values |
| -------------- | ------ | -------- | -------------- |
| `server`       | `Text` | `True`   | N/A            |
| `cluster`      | `Text` | `False`  | N/A            |
| `randomString` | `Text` | `False`  | N/A            |

## `MongoDBForPipeline` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `Basic`
- Supported Connection Encryption Types: `NotEncrypted`

### `MongoDBForPipeline.Contents` Creation Method

| Name     | Type   | Required | Allowed Values |
| -------- | ------ | -------- | -------------- |
| `server` | `Text` | `True`   | N/A            |

## `MySql` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `Basic`
- Supported Connection Encryption Types: `NotEncrypted`, `Encrypted`

### `MySql` Creation Method

| Name       | Type   | Required | Allowed Values |
| ---------- | ------ | -------- | -------------- |
| `server`   | `Text` | `True`   | N/A            |
| `database` | `Text` | `True`   | N/A            |

## `Netezza` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `Basic`
- Supported Connection Encryption Types: `NotEncrypted`, `Encrypted`

### `Netezza.Database` Creation Method

| Name                         | Type       | Required | Allowed Values |
| ---------------------------- | ---------- | -------- | -------------- |
| `server`                     | `Text`     | `True`   | N/A            |
| `database`                   | `Text`     | `True`   | N/A            |
| `ConnectionTimeout`          | `Duration` | `False`  | N/A            |
| `CommandTimeout`             | `Duration` | `False`  | N/A            |
| `NormalizeDatabaseName`      | `Boolean`  | `False`  | N/A            |
| `HierarchicalNavigation`     | `Boolean`  | `False`  | N/A            |
| `CreateNavigationProperties` | `Boolean`  | `False`  | N/A            |

## `OData` Type

- Support 'Skip Test Connection': `True`
- Supported Credential Types: `Anonymous`, `Basic`, `Key`
- Supported Connection Encryption Types: `NotEncrypted`

### `OData` Creation Method

| Name  | Type   | Required | Allowed Values |
| ----- | ------ | -------- | -------------- |
| `url` | `Text` | `True`   | N/A            |

## `OracleCloudStorage` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `Basic`
- Supported Connection Encryption Types: `NotEncrypted`

### `OracleCloudStorage.Contents` Creation Method

| Name          | Type   | Required | Allowed Values |
| ------------- | ------ | -------- | -------------- |
| `APIEndpoint` | `Text` | `True`   | N/A            |

## `Paxata` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `Key`
- Supported Connection Encryption Types: `NotEncrypted`

### `Paxata.Contents` Creation Method

| Name  | Type   | Required | Allowed Values |
| ----- | ------ | -------- | -------------- |
| `url` | `Text` | `True`   | N/A            |

## `Plantronics` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `Anonymous`
- Supported Connection Encryption Types: `NotEncrypted`

### `Plantronics.Feed` Creation Method

| Name              | Type     | Required | Allowed Values |
| ----------------- | -------- | -------- | -------------- |
| `URL`             | `Text`   | `True`   | N/A            |
| `Tenant`          | `Text`   | `True`   | N/A            |
| `URL1`            | `Text`   | `False`  | N/A            |
| `URL2`            | `Text`   | `False`  | N/A            |
| `Parameters`      | `Text`   | `False`  | N/A            |
| `ElementsPerPage` | `Number` | `False`  | N/A            |

### `Plantronics.Test` Creation Method

| Name     | Type   | Required | Allowed Values |
| -------- | ------ | -------- | -------------- |
| `URL`    | `Text` | `True`   | N/A            |
| `Tenant` | `Text` | `True`   | N/A            |

## `PlanviewEnterprise` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `Basic`
- Supported Connection Encryption Types: `NotEncrypted`

### `PlanviewEnterprise.Feed` Creation Method

| Name       | Type   | Required | Allowed Values |
| ---------- | ------ | -------- | -------------- |
| `url`      | `Text` | `True`   | N/A            |
| `database` | `Text` | `True`   | N/A            |

### `PlanviewEnterprise.CallQueryService` Creation Method

| Name             | Type   | Required | Allowed Values |
| ---------------- | ------ | -------- | -------------- |
| `url`            | `Text` | `True`   | N/A            |
| `database`       | `Text` | `True`   | N/A            |
| `sqlQueryString` | `Text` | `True`   | N/A            |

## `PostgreSQL` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `Basic`
- Supported Connection Encryption Types: `NotEncrypted`, `Encrypted`

### `PostgreSql` Creation Method

| Name       | Type   | Required | Allowed Values |
| ---------- | ------ | -------- | -------------- |
| `server`   | `Text` | `True`   | N/A            |
| `database` | `Text` | `True`   | N/A            |

## `PowerBIDatasets` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `ServicePrincipal`
- Supported Connection Encryption Types: `NotEncrypted`

## `PowerGP` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `Basic`
- Supported Connection Encryption Types: `NotEncrypted`

### `PowerGP.GetData` Creation Method

| Name          | Type   | Required | Allowed Values |
| ------------- | ------ | -------- | -------------- |
| `baseURL`     | `Text` | `True`   | N/A            |
| `company`     | `Text` | `True`   | N/A            |
| `RelativeURL` | `Text` | `True`   | N/A            |

## `Prevedere` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `Key`
- Supported Connection Encryption Types: `NotEncrypted`

### `Prevedere.Contents` Creation Method

| Name  | Type   | Required | Allowed Values |
| ----- | ------ | -------- | -------------- |
| `url` | `Text` | `True`   | N/A            |

## `ProductioneerMExt` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `Basic`
- Supported Connection Encryption Types: `NotEncrypted`

### `ProductioneerExt.Feed` Creation Method

| Name          | Type   | Required | Allowed Values |
| ------------- | ------ | -------- | -------------- |
| `CompanyName` | `Text` | `True`   | N/A            |
| `endpoint`    | `Text` | `True`   | N/A            |
| `StartDate`   | `Text` | `True`   | N/A            |

### `ProductioneerMExt.Contents` Creation Method

| Name          | Type   | Required | Allowed Values |
| ------------- | ------ | -------- | -------------- |
| `CompanyName` | `Text` | `True`   | N/A            |
| `endpoint`    | `Text` | `True`   | N/A            |
| `StartDate`   | `Text` | `True`   | N/A            |

## `ProjectIntelligence` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `Key`
- Supported Connection Encryption Types: `NotEncrypted`

### `ProjectIntelligence.Service` Creation Method

| Name                     | Type   | Required | Allowed Values |
| ------------------------ | ------ | -------- | -------------- |
| `ShareAdvanceWebAddress` | `Text` | `True`   | N/A            |
| `Dimension`              | `Text` | `False`  | N/A            |
| `DataType`               | `Text` | `False`  | N/A            |

## `QuestionPro` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `Key`
- Supported Connection Encryption Types: `NotEncrypted`

### `QuestionPro.Contents` Creation Method

| Name        | Type   | Required | Allowed Values |
| ----------- | ------ | -------- | -------------- |
| `access_id` | `Text` | `True`   | N/A            |

## `QuickBase` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `Key`
- Supported Connection Encryption Types: `NotEncrypted`

### `QuickBase.Contents` Creation Method

| Name  | Type   | Required | Allowed Values |
| ----- | ------ | -------- | -------------- |
| `url` | `Text` | `True`   | N/A            |

## `RestService` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `Basic`, `Anonymous`, `ServicePrincipal`
- Supported Connection Encryption Types: `NotEncrypted`

### `RestService.Contents` Creation Method

| Name       | Type   | Required | Allowed Values |
| ---------- | ------ | -------- | -------------- |
| `baseUrl`  | `Text` | `True`   | N/A            |
| `audience` | `Text` | `False`  | N/A            |

## `RiskAssurancePlatform` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `Basic`
- Supported Connection Encryption Types: `NotEncrypted`

### `RiskAssurance.GetData` Creation Method

| Name          | Type   | Required | Allowed Values |
| ------------- | ------ | -------- | -------------- |
| `baseURL`     | `Text` | `True`   | N/A            |
| `RelativeURL` | `Text` | `True`   | N/A            |
| `customer`    | `Text` | `True`   | N/A            |

## `Roamler` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `Key`
- Supported Connection Encryption Types: `NotEncrypted`

## `SalesforceServiceCloud` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `Basic`
- Supported Connection Encryption Types: `NotEncrypted`

### `SalesforceServiceCloud.Contents` Creation Method

| Name             | Type   | Required | Allowed Values |
| ---------------- | ------ | -------- | -------------- |
| `environmentURL` | `Text` | `True`   | N/A            |

## `Samsara` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `Key`
- Supported Connection Encryption Types: `NotEncrypted`

### `Samsara.Records` Creation Method

| Name         | Type   | Required | Allowed Values |
| ------------ | ------ | -------- | -------------- |
| `Region`     | `Text` | `True`   | `US`, `EU`     |
| `RangeStart` | `Text` | `True`   | N/A            |
| `RangeEnd`   | `Text` | `False`  | N/A            |

## `SDMX` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `Anonymous`
- Supported Connection Encryption Types: `NotEncrypted`

### `SDMX.Contents` Creation Method

| Name       | Type   | Required | Allowed Values                                                 |
| ---------- | ------ | -------- | -------------------------------------------------------------- |
| `url`      | `Text` | `True`   | N/A                                                            |
| `Option`   | `Text` | `True`   | `Show codes and labels`, `Show codes only`, `Show labels only` |
| `Language` | `Text` | `False`  | N/A                                                            |

## `ServiceNow` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `Basic`
- Supported Connection Encryption Types: `NotEncrypted`

### `ServiceNow.Data` Creation Method

| Name       | Type   | Required | Allowed Values |
| ---------- | ------ | -------- | -------------- |
| `instance` | `Text` | `True`   | N/A            |

### `ServiceNow.Simple` Creation Method

| Name       | Type   | Required | Allowed Values |
| ---------- | ------ | -------- | -------------- |
| `instance` | `Text` | `True`   | N/A            |
| `endpoint` | `Text` | `True`   | N/A            |

### `ServiceNow.Contents` Creation Method

| Name             | Type     | Required | Allowed Values |
| ---------------- | -------- | -------- | -------------- |
| `instance`       | `Text`   | `True`   | N/A            |
| `endpoint`       | `Text`   | `True`   | N/A            |
| `textArguments`  | `Text`   | `False`  | N/A            |
| `maxRecordCount` | `Number` | `False`  | N/A            |

### `ServiceNow.PerformanceAnalytics` Creation Method

| Name            | Type   | Required | Allowed Values |
| --------------- | ------ | -------- | -------------- |
| `instance`      | `Text` | `True`   | N/A            |
| `textArguments` | `Text` | `False`  | N/A            |

### `ServiceNow.Aggregate` Creation Method

| Name            | Type   | Required | Allowed Values |
| --------------- | ------ | -------- | -------------- |
| `instance`      | `Text` | `True`   | N/A            |
| `endpoint`      | `Text` | `True`   | N/A            |
| `textArguments` | `Text` | `False`  | N/A            |

## `SFTP` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `Basic`
- Supported Connection Encryption Types: `NotEncrypted`

### `SFTP.Contents` Creation Method

| Name          | Type   | Required | Allowed Values |
| ------------- | ------ | -------- | -------------- |
| `server`      | `Text` | `True`   | N/A            |
| `fingerprint` | `Text` | `False`  | N/A            |

## `SharePoint` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `Anonymous`, `ServicePrincipal`
- Supported Connection Encryption Types: `NotEncrypted`

### `SharePointList` Creation Method

| Name                | Type   | Required | Allowed Values |
| ------------------- | ------ | -------- | -------------- |
| `sharePointSiteUrl` | `Text` | `True`   | N/A            |

## `ShortcutsBI` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `Key`
- Supported Connection Encryption Types: `NotEncrypted`

## `Siteimprove` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `Basic`
- Supported Connection Encryption Types: `NotEncrypted`

## `Snowflake` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `Basic`
- Supported Connection Encryption Types: `NotEncrypted`

### `Snowflake.Databases` Creation Method

| Name                         | Type      | Required | Allowed Values |
| ---------------------------- | --------- | -------- | -------------- |
| `server`                     | `Text`    | `True`   | N/A            |
| `warehouse`                  | `Text`    | `True`   | N/A            |
| `Role`                       | `Text`    | `False`  | N/A            |
| `CreateNavigationProperties` | `Boolean` | `False`  | N/A            |
| `ConnectionTimeout`          | `Number`  | `False`  | N/A            |
| `CommandTimeout`             | `Number`  | `False`  | N/A            |

## `SolarWindsServiceDesk` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `Key`
- Supported Connection Encryption Types: `NotEncrypted`

### `SolarWindsServiceDesk.ContentsV110` Creation Method

| Name         | Type       | Required | Allowed Values |
| ------------ | ---------- | -------- | -------------- |
| `RangeStart` | `DateTime` | `False`  | N/A            |
| `RangeEnd`   | `DateTime` | `False`  | N/A            |

## `Spark` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `Basic`
- Supported Connection Encryption Types: `NotEncrypted`

### `Spark.Tables` Creation Method

| Name                     | Type      | Required | Allowed Values |
| ------------------------ | --------- | -------- | -------------- |
| `server`                 | `Text`    | `True`   | N/A            |
| `protocol`               | `Number`  | `True`   | `0`, `1`, `2`  |
| `BatchSize`              | `Number`  | `False`  | N/A            |
| `HierarchicalNavigation` | `Boolean` | `False`  | N/A            |

### `AzureSpark.Tables` Creation Method

| Name                     | Type      | Required | Allowed Values |
| ------------------------ | --------- | -------- | -------------- |
| `server`                 | `Text`    | `True`   | N/A            |
| `BatchSize`              | `Number`  | `False`  | N/A            |
| `HierarchicalNavigation` | `Boolean` | `False`  | N/A            |

### `ApacheSpark.Tables` Creation Method

| Name                     | Type      | Required | Allowed Values |
| ------------------------ | --------- | -------- | -------------- |
| `server`                 | `Text`    | `True`   | N/A            |
| `protocol`               | `Number`  | `True`   | `0`, `2`       |
| `BatchSize`              | `Number`  | `False`  | N/A            |
| `HierarchicalNavigation` | `Boolean` | `False`  | N/A            |

## `SparkPost` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `Key`
- Supported Connection Encryption Types: `NotEncrypted`

### `SparkPost.GetList` Creation Method

| Name   | Type   | Required | Allowed Values |
| ------ | ------ | -------- | -------------- |
| `Path` | `Text` | `True`   | N/A            |

### `SparkPost.GetTable` Creation Method

| Name               | Type     | Required | Allowed Values |
| ------------------ | -------- | -------- | -------------- |
| `DaysToAggregate`  | `Number` | `True`   | N/A            |
| `MetricColumns`    | `Text`   | `True`   | N/A            |
| `NonMetricColumns` | `Text`   | `True`   | N/A            |
| `Path`             | `Text`   | `True`   | N/A            |

### `SparkPost.NavTable` Creation Method

| Name              | Type     | Required | Allowed Values |
| ----------------- | -------- | -------- | -------------- |
| `DaysToAggregate` | `Number` | `True`   | N/A            |

## `SQL` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `Basic`, `ServicePrincipal`
- Supported Connection Encryption Types: `NotEncrypted`, `Encrypted`

### `Sql` Creation Method

| Name       | Type   | Required | Allowed Values |
| ---------- | ------ | -------- | -------------- |
| `server`   | `Text` | `True`   | N/A            |
| `database` | `Text` | `False`  | N/A            |

## `Stripe` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `Key`
- Supported Connection Encryption Types: `NotEncrypted`

### `Stripe.Contents` Creation Method

| Name        | Type     | Required | Allowed Values |
| ----------- | -------- | -------- | -------------- |
| `method`    | `Text`   | `True`   | N/A            |
| `pageLimit` | `Number` | `False`  | N/A            |

## `SurveyMonkey` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `Key`
- Supported Connection Encryption Types: `NotEncrypted`

## `SweetIQ` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `Key`
- Supported Connection Encryption Types: `NotEncrypted`

### `SweetIQ.Contents` Creation Method

| Name                | Type   | Required | Allowed Values |
| ------------------- | ------ | -------- | -------------- |
| `clientId`          | `Text` | `False`  | N/A            |
| `path`              | `Text` | `False`  | N/A            |
| `optionalParameter` | `Text` | `False`  | N/A            |

## `Synapse` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `ServicePrincipal`
- Supported Connection Encryption Types: `NotEncrypted`

## `Tenforce` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `Basic`
- Supported Connection Encryption Types: `NotEncrypted`

### `Tenforce.Contents` Creation Method

| Name             | Type   | Required | Allowed Values              |
| ---------------- | ------ | -------- | --------------------------- |
| `ApplicationUrl` | `Text` | `True`   | N/A                         |
| `ListId`         | `Text` | `True`   | N/A                         |
| `DataType`       | `Text` | `True`   | `Do not include`, `Include` |

## `Timelog` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `Anonymous`
- Supported Connection Encryption Types: `NotEncrypted`

### `Timelog.Tables` Creation Method

| Name                   | Type   | Required | Allowed Values |
| ---------------------- | ------ | -------- | -------------- |
| `SiteCode`             | `Text` | `True`   | N/A            |
| `ApiID`                | `Text` | `True`   | N/A            |
| `ApiPassword`          | `Text` | `True`   | N/A            |
| `URLAccountName`       | `Text` | `True`   | N/A            |
| `DefaultNumberAccount` | `Text` | `True`   | N/A            |
| `CurrentMonthOnly`     | `Text` | `True`   | N/A            |

## `Troux` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `Basic`
- Supported Connection Encryption Types: `NotEncrypted`

### `Troux.Feed` Creation Method

| Name  | Type   | Required | Allowed Values |
| ----- | ------ | -------- | -------------- |
| `url` | `Text` | `True`   | N/A            |

### `Troux.CustomFeed` Creation Method

| Name    | Type   | Required | Allowed Values |
| ------- | ------ | -------- | -------------- |
| `url`   | `Text` | `True`   | N/A            |
| `query` | `Text` | `True`   | N/A            |

### `Troux.TestConnection` Creation Method

| Name  | Type   | Required | Allowed Values |
| ----- | ------ | -------- | -------------- |
| `url` | `Text` | `True`   | N/A            |

## `Usercube` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `Basic`
- Supported Connection Encryption Types: `NotEncrypted`

### `Usercube.Universes` Creation Method

| Name        | Type   | Required | Allowed Values |
| ----------- | ------ | -------- | -------------- |
| `serverUrl` | `Text` | `True`   | N/A            |

## `Vena` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `Basic`
- Supported Connection Encryption Types: `NotEncrypted`

### `Vena.Contents` Creation Method

| Name         | Type   | Required | Allowed Values                                                                                                    |
| ------------ | ------ | -------- | ----------------------------------------------------------------------------------------------------------------- |
| `source`     | `Text` | `True`   | `https://ca3.vena.io`, `https://us3.vena.io`, `https://us2.vena.io`, `https://us1.vena.io`, `https://eu1.vena.io` |
| `modelQuery` | `Text` | `False`  | N/A                                                                                                               |
| `apiVersion` | `Text` | `False`  | `v1`, `v2`                                                                                                        |

## `Vertica` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `Basic`
- Supported Connection Encryption Types: `NotEncrypted`

### `Vertica.Database` Creation Method

| Name       | Type   | Required | Allowed Values |
| ---------- | ------ | -------- | -------------- |
| `server`   | `Text` | `True`   | N/A            |
| `database` | `Text` | `True`   | N/A            |

## `Visual Studio Team Services` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `Basic`
- Supported Connection Encryption Types: `NotEncrypted`

### `VSTS.Contents` Creation Method

| Name  | Type   | Required | Allowed Values |
| ----- | ------ | -------- | -------------- |
| `url` | `Text` | `True`   | N/A            |

## `VSTS` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `Basic`, `Key`
- Supported Connection Encryption Types: `NotEncrypted`

### `VSTS.AccountContents` Creation Method

| Name  | Type   | Required | Allowed Values |
| ----- | ------ | -------- | -------------- |
| `url` | `Text` | `True`   | N/A            |

### `VSTS.Feed` Creation Method

| Name  | Type   | Required | Allowed Values |
| ----- | ------ | -------- | -------------- |
| `url` | `Text` | `True`   | N/A            |

### `VSTS.AnalyticsViews` Creation Method

| Name      | Type   | Required | Allowed Values |
| --------- | ------ | -------- | -------------- |
| `url`     | `Text` | `True`   | N/A            |
| `project` | `Text` | `True`   | N/A            |

## `Web` Type

- Support 'Skip Test Connection': `True`
- Supported Credential Types: `Anonymous`, `Basic`, `ServicePrincipal`
- Supported Connection Encryption Types: `NotEncrypted`

### `Web` Creation Method

| Name  | Type   | Required | Allowed Values |
| ----- | ------ | -------- | -------------- |
| `url` | `Text` | `True`   | N/A            |

## `WebForPipeline` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `Basic`, `Anonymous`, `ServicePrincipal`
- Supported Connection Encryption Types: `NotEncrypted`

### `WebForPipeline.Contents` Creation Method

| Name       | Type   | Required | Allowed Values |
| ---------- | ------ | -------- | -------------- |
| `baseUrl`  | `Text` | `True`   | N/A            |
| `audience` | `Text` | `False`  | N/A            |

## `Webtrends` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `Basic`
- Supported Connection Encryption Types: `NotEncrypted`

### `Webtrends.KeyMetrics` Creation Method

| Name        | Type   | Required | Allowed Values |
| ----------- | ------ | -------- | -------------- |
| `ProfileId` | `Text` | `True`   | N/A            |
| `startDate` | `Date` | `False`  | N/A            |
| `endDate`   | `Date` | `False`  | N/A            |

### `Webtrends.ReportContents` Creation Method

| Name        | Type   | Required | Allowed Values |
| ----------- | ------ | -------- | -------------- |
| `ProfileId` | `Text` | `True`   | N/A            |
| `ReportId`  | `Text` | `True`   | N/A            |
| `startDate` | `Date` | `False`  | N/A            |
| `endDate`   | `Date` | `False`  | N/A            |

### `Webtrends.Tables` Creation Method

| Name        | Type   | Required | Allowed Values |
| ----------- | ------ | -------- | -------------- |
| `ProfileId` | `Text` | `True`   | N/A            |
| `startDate` | `Date` | `False`  | N/A            |
| `endDate`   | `Date` | `False`  | N/A            |

## `WebtrendsAnalytics` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `Basic`
- Supported Connection Encryption Types: `NotEncrypted`

### `WebtrendsAnalytics.Tables` Creation Method

| Name          | Type   | Required | Allowed Values                   |
| ------------- | ------ | -------- | -------------------------------- |
| `ProfileId`   | `Text` | `True`   | N/A                              |
| `Period`      | `Text` | `True`   | `Custom Date`, `Report Period`   |
| `reportType`  | `Text` | `True`   | `Summary`, `Trend`, `Individual` |
| `startDate`   | `Date` | `False`  | N/A                              |
| `endDate`     | `Date` | `False`  | N/A                              |
| `startPeriod` | `Text` | `False`  | N/A                              |
| `endPeriod`   | `Text` | `False`  | N/A                              |

## `WorkforceDimensions` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `Basic`
- Supported Connection Encryption Types: `NotEncrypted`

### `WorkforceDimensions.Contents` Creation Method

| Name                        | Type   | Required | Allowed Values                                                                                                                                                                                                                                            |
| --------------------------- | ------ | -------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `configurationServer`       | `Text` | `True`   | N/A                                                                                                                                                                                                                                                       |
| `workForceDimensionsServer` | `Text` | `True`   | N/A                                                                                                                                                                                                                                                       |
| `symbolicPeriod`            | `Text` | `True`   | `Date Range (start and end dates are required)`, `Previous Pay Period`, `Current Pay Period`, `Next Pay Period`, `Today`, `Yesterday`, `Yesterday, Today, Tomorrow`, `Yesterday Plus 6 Days`, `Last 30 Days`, `Last 90 Days`, `Last Week`, `Current Week` |
| `startDate`                 | `Date` | `False`  | N/A                                                                                                                                                                                                                                                       |
| `endDate`                   | `Date` | `False`  | N/A                                                                                                                                                                                                                                                       |

## `WtsParadigm` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `Basic`
- Supported Connection Encryption Types: `NotEncrypted`

## `Zucchetti` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `Basic`
- Supported Connection Encryption Types: `NotEncrypted`

### `Zucchetti.Contents` Creation Method

| Name          | Type   | Required | Allowed Values |
| ------------- | ------ | -------- | -------------- |
| `Url`         | `Text` | `True`   | N/A            |
| `Environment` | `Text` | `True`   | N/A            |

## `Zuora` Type

- Support 'Skip Test Connection': `False`
- Supported Credential Types: `Basic`
- Supported Connection Encryption Types: `NotEncrypted`

### `Zuora.Export` Creation Method

| Name          | Type   | Required | Allowed Values |
| ------------- | ------ | -------- | -------------- |
| `QueryString` | `Text` | `True`   | N/A            |
| `APIRoot`     | `Text` | `True`   | N/A            |
| `SchemaName`  | `Text` | `False`  | N/A            |
| `APIVersion`  | `Text` | `False`  | N/A            |
