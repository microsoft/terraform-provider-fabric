// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package connection

import (
	"context"

	supertypes "github.com/FrangipaneTeam/terraform-plugin-framework-supertypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
)

type baseDataSourceConnectionModel struct {
	connectionModel
	ConnectionDetails supertypes.SingleNestedObjectValueOf[connectionDetailsModel] `tfsdk:"connection_details"`
	CredentialDetails supertypes.SingleNestedObjectValueOf[credentialDetailsModel] `tfsdk:"credential_details"`
}

type baseResourceConnectionModel struct {
	connectionModel
	ConnectionDetails supertypes.SingleNestedObjectValueOf[rsConnectionDetailsModel] `tfsdk:"connection_details"`
	CredentialDetails supertypes.SingleNestedObjectValueOf[rsCredentialDetailsModel] `tfsdk:"credential_details"`
}

func (to *baseResourceConnectionModel) setConnectionDetails(ctx context.Context, from fabcore.Connection) diag.Diagnostics {
	connectionDetails := supertypes.NewSingleNestedObjectValueOfNull[rsConnectionDetailsModel](ctx)

	connectionDetails = to.ConnectionDetails

	if from.ConnectionDetails != nil {
		connectionDetailsModel := &rsConnectionDetailsModel{}
		connectionDetailsModel.set(*from.ConnectionDetails)

		diags := connectionDetails.Set(ctx, connectionDetailsModel)
		if diags.HasError() {
			return diags
		}
	}

	to.ConnectionDetails = connectionDetails

	return nil
}

func (to *baseResourceConnectionModel) setCredentialDetails(ctx context.Context, from fabcore.Connection) diag.Diagnostics {
	credentialDetails := supertypes.NewSingleNestedObjectValueOfNull[rsCredentialDetailsModel](ctx)

	credentialDetails = to.CredentialDetails

	if from.CredentialDetails != nil {
		credentialDetailsModel := &rsCredentialDetailsModel{}
		credentialDetailsModel.set(*from.CredentialDetails)

		diags := credentialDetails.Set(ctx, credentialDetailsModel)
		if diags.HasError() {
			return diags
		}

		to.CredentialDetails = credentialDetails
	}

	return nil
}

type connectionModel struct {
	ID               customtypes.UUID `tfsdk:"id"`
	DisplayName      types.String     `tfsdk:"display_name"`
	GatewayID        customtypes.UUID `tfsdk:"gateway_id"`
	ConnectivityType types.String     `tfsdk:"connectivity_type"`
	PrivacyLevel     types.String     `tfsdk:"privacy_level"`
}

func (to *connectionModel) set(from fabcore.Connection) {
	to.ID = customtypes.NewUUIDPointerValue(from.ID)
	to.DisplayName = types.StringPointerValue(from.DisplayName)
	to.GatewayID = customtypes.NewUUIDPointerValue(from.GatewayID)
	to.ConnectivityType = types.StringPointerValue((*string)(from.ConnectivityType))
	to.PrivacyLevel = types.StringPointerValue((*string)(from.PrivacyLevel))
}

func (to *baseDataSourceConnectionModel) setConnectionDetails(ctx context.Context, from *fabcore.ListConnectionDetails) diag.Diagnostics {
	connectionDetails := supertypes.NewSingleNestedObjectValueOfNull[connectionDetailsModel](ctx)

	if from != nil {
		connectionDetailsModel := &connectionDetailsModel{}
		connectionDetailsModel.set(*from)

		diags := connectionDetails.Set(ctx, connectionDetailsModel)
		if diags.HasError() {
			return diags
		}
	}

	to.ConnectionDetails = connectionDetails

	return nil
}

func (to *baseDataSourceConnectionModel) setCredentialDetails(ctx context.Context, from *fabcore.ListCredentialDetails) diag.Diagnostics {
	credentialDetails := supertypes.NewSingleNestedObjectValueOfNull[credentialDetailsModel](ctx)

	if from != nil {
		credentialDetailsModel := &credentialDetailsModel{}
		credentialDetailsModel.set(*from)

		diags := credentialDetails.Set(ctx, credentialDetailsModel)
		if diags.HasError() {
			return diags
		}
	}

	to.CredentialDetails = credentialDetails

	return nil
}

type connectionDetailsModel struct {
	Path types.String `tfsdk:"path"`
	Type types.String `tfsdk:"type"`
}

func (to *connectionDetailsModel) set(from fabcore.ListConnectionDetails) {
	to.Path = types.StringPointerValue(from.Path)
	to.Type = types.StringPointerValue(from.Type)
}

type rsConnectionDetailsModel struct {
	Path           types.String                  `tfsdk:"path"` // computed
	Type           types.String                  `tfsdk:"type"`
	CreationMethod types.String                  `tfsdk:"creation_method"`
	Parameters     supertypes.MapValueOf[string] `tfsdk:"parameters"`
}

func (m rsConnectionDetailsModel) getParameters(ctx context.Context) (map[string]string, diag.Diagnostics) {
	if !m.Parameters.IsNull() && !m.Parameters.IsUnknown() {
		return m.Parameters.Get(ctx)
	}

	return nil, nil
}

func (to *rsConnectionDetailsModel) set(from fabcore.ListConnectionDetails) {
	to.Path = types.StringPointerValue(from.Path)
	to.Type = types.StringPointerValue(from.Type)
}

type credentialDetailsModel struct {
	ConnectionEncryption types.String `tfsdk:"connection_encryption"`
	CredentialType       types.String `tfsdk:"credential_type"`
	SingleSignOnType     types.String `tfsdk:"single_sign_on_type"`
	SkipTestConnection   types.Bool   `tfsdk:"skip_test_connection"`
}

func (to *credentialDetailsModel) set(from fabcore.ListCredentialDetails) {
	to.CredentialType = types.StringPointerValue((*string)(from.CredentialType))
	to.ConnectionEncryption = types.StringPointerValue((*string)(from.ConnectionEncryption))
	to.SingleSignOnType = types.StringPointerValue((*string)(from.SingleSignOnType))
	to.SkipTestConnection = types.BoolPointerValue(from.SkipTestConnection)
}

type rsCredentialDetailsModel struct {
	ConnectionEncryption types.String `tfsdk:"connection_encryption"`
	SingleSignOnType     types.String `tfsdk:"single_sign_on_type"`
	SkipTestConnection   types.Bool   `tfsdk:"skip_test_connection"`
	CredentialType       types.String `tfsdk:"credential_type"`

	// AnonymousCredentials                   supertypes.SingleNestedObjectValueOf[anonymousCredentialsModel]                   `tfsdk:"anonymous_credentials"`
	BasicCredentials                 supertypes.SingleNestedObjectValueOf[credentialsBasicModel]                 `tfsdk:"basic_credentials"`
	KeyCredentials                   supertypes.SingleNestedObjectValueOf[credentialsKeyModel]                   `tfsdk:"key_credentials"`
	ServicePrincipalCredentials      supertypes.SingleNestedObjectValueOf[credentialsServicePrincipalModel]      `tfsdk:"service_principal_credentials"`
	SharedAccessSignatureCredentials supertypes.SingleNestedObjectValueOf[credentialsSharedAccessSignatureModel] `tfsdk:"shared_access_signature_credentials"`
	WindowsCredentials               supertypes.SingleNestedObjectValueOf[credentialsWindowsModel]               `tfsdk:"windows_credentials"`
	// WindowsWithoutImpersonationCredentials supertypes.SingleNestedObjectValueOf[credentialsWindowsWithoutImpersonationModel] `tfsdk:"windows_without_impersonation_credentials"`
	// WorkspaceIdentityCredentials           supertypes.SingleNestedObjectValueOf[credentialsWorkspaceIdentityModel]           `tfsdk:"workspace_identity_credentials"`
}

func (to *rsCredentialDetailsModel) set(from fabcore.ListCredentialDetails) {
	to.ConnectionEncryption = types.StringPointerValue((*string)(from.ConnectionEncryption))
	to.SingleSignOnType = types.StringPointerValue((*string)(from.SingleSignOnType))
	to.SkipTestConnection = types.BoolPointerValue(from.SkipTestConnection)
	to.CredentialType = types.StringPointerValue((*string)(from.CredentialType))
}

type credentialsBasicModel struct {
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
}

type credentialsKeyModel struct {
	Key types.String `tfsdk:"key"`
}

type credentialsServicePrincipalModel struct {
	TenantID     types.String `tfsdk:"tenant_id"`
	ClientID     types.String `tfsdk:"client_id"`
	ClientSecret types.String `tfsdk:"client_secret"`
}

type credentialsSharedAccessSignatureModel struct {
	Token types.String `tfsdk:"token"`
}

type credentialsWindowsModel struct {
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
}

func possibleConnectionTypeValues() []string {
	return []string{
		"SQL",
		"AnalysisServices",
		"SharePoint",
		"Web",
		"OData",
		"MySql",
		"PostgreSQL",
		"AzureTables",
		"AzureBlobs",
		"GoogleAnalytics",
		"Salesforce",
		"AdobeAnalytics",
		"AzureDataLakeStorage",
		"Exchange",
		"AIInsights",
		"AIFunctions",
		"AzureMLFunctions",
		"appFigures",
		"GoogleBigQuery",
		"GoogleBigQueryAad",
		"PowerBI",
		"Cds",
		"DataLake",
		"DataWorld",
		"DocumentDB",
		"Dynamics365BusinessCentral",
		"Dynamics 365 Business Central (on-premises)",
		"Dynamics NAV",
		"github",
		"Impala",
		"AzureHive",
		"ApacheHive",
		"Kusto",
		"AzureDataExplorer",
		"LinkedIn",
		"MailChimp",
		"Netezza",
		"PlanviewEnterprise",
		"Projectplace",
		"QuickBooksOnline",
		"AmazonRedshift",
		"Smartsheet",
		"Snowflake",
		"Spark",
		"SparkPost",
		"Stripe",
		"SweetIQ",
		"Troux",
		"twilio",
		"Visual Studio Team Services",
		"VSTS",
		"Webtrends",
		"Vertica",
		"zendesk",
		"Acterys",
		"Anaplan",
		"Asana",
		"AutodeskConstructionCloud",
		"AutomationAnywhere",
		"AutomyDataAnalytics",
		"BI360",
		"BitSightSecurityRatings",
		"Bloomberg",
		"BQL",
		"BQECore",
		"BuildingConnected",
		"CCHTagetik",
		"CDataConnectCloud",
		"Celonis",
		"Cherwell",
		"CloudBluePSA",
		"Cognite",
		"Databricks",
		"DatabricksMultiCloud",
		"DeltaSharing",
		"Dremio",
		"DremioCloud",
		"EduFrame",
		"EmigoDataSourceConnector",
		"EntersoftBusinessSuite",
		"EQuIS",
		"eWayCRM",
		"FactSetAnalytics",
		"FactSetRMS",
		"Funnel",
		"HexagonSmartApi",
		"IndustrialAppStore",
		"InformationGrid",
		"inwink",
		"JamfPro",
		"Kognitwin",
		"kxkdbinsightsenterprise",
		"LEAP",
		"LinkedInLearning",
		"MicroStrategyDataset",
		"OneStream",
		"Paxata",
		"PlanviewOKR",
		"PlanviewProjectplace",
		"Profisee",
		"QuickBase",
		"Roamler",
		"Samsara",
		"SDMX",
		"ShortcutsBI",
		"Siteimprove",
		"SmartsheetGlobal",
		"SoftOneBI",
		"SolarWindsServiceDesk",
		"Spigit",
		"SumTotal",
		"Supermetrics",
		"SurveyMonkey",
		"Tenforce",
		"Usercube",
		"Vena",
		"VesselInsight",
		"WebtrendsAnalytics",
		"Windsor",
		"Witivio",
		"WorkforceDimensions",
		"Wrike",
		"ZendeskData",
		"ZohoCreator",
		"Zucchetti",
		"AtScale",
		"AzureCostManagement",
		"AzureResourceGraph",
		"AzureTrino",
		"CommonDataService",
		"CosmosDB",
		"CustomerInsights",
		"Fhir",
		"GoogleSheets",
		"Intune",
		"Lakehouse",
		"MicrosoftGraphSecurity",
		"PowerBIDatamarts",
		"PowerPlatformDataflows",
		"ProductInsights",
		"Synapse",
		"TeamsAnalytics",
		"Warehouse",
		"VivaInsights",
		"VivaInsightsApi",
		"WorkplaceAnalytics",
		"AdlsGen2CosmosStructuredStream",
		"AmazonRdsForSqlServer",
		"AmazonS3",
		"AmazonS3Compatible",
		"AzureAISearch",
		"AzureBatch",
		"AzureCosmosDBForMongoDB",
		"AzureDatabaseForMySQL",
		"AzureDatabricksWorkspace",
		"AzureDataFactory",
		"AzureDataLakeStoreCosmosStructuredStream",
		"AzureFiles",
		"AzureFunction",
		"AzureHDInsightCluster",
		"AzureHDInsightOnDemandCluster",
		"AzureKeyVault",
		"AzureMachineLearning",
		"AzurePostgreSQL",
		"AzureServiceBus",
		"AzureSqlMI",
		"AzureSynapseWorkspace",
		"ConfluentCloud",
		"DataLakeAnalytics",
		"DataPipelineCosmosDb",
		"Dynamics365",
		"DynamicsAX",
		"DynamicsCrm",
		"EventHub",
		"FabricDataPipelines",
		"FTP",
		"GoogleCloudStorage",
		"GooglePubSub",
		"HttpServer",
		"IoTHub",
		"Kinesis",
		"MariaDBForPipeline",
		"MicrosoftOutlook",
		"MicrosoftTeams",
		"MongoDBAtlasForPipeline",
		"MongoDBForPipeline",
		"Notebook",
		"OracleCloudStorage",
		"PowerBIDatasets",
		"RestService",
		"SalesforceServiceCloud",
		"ServiceNow",
		"SFTP",
		"SparkJobDefinition",
		"UserDataFunctions",
		"WebForPipeline",
		"AzureDevOpsSourceControl",
		"GitHubSourceControl",
		"GoogleAds",
		"Microsoft365",
		"accelo",
		"Adobe Analytics",
		"AlpineMetrics",
		"Applixure",
		"AriaConnector",
		"AriaConnector.External",
		"AzureEnterprise",
		"Basemethod",
		"CloudScope",
		"CloudScopeInstagram",
		"Commercial Partner Analytics Connector",
		"comScore",
		"datawiz",
		"InfinityConnector",
		"InLooxNow",
		"insightCentr",
		"Insightly",
		"IntelliBoard",
		"JDIConnector",
		"KaizalaAttendanceReports",
		"KaizalaReports",
		"KaizalaSurveyReports",
		"Mandrill",
		"AdMaD",
		"myob_ar",
		"Office365Adoption",
		"Office365Adoption AAD",
		"Office365Mon Reporting Data",
		"Office365Mon2",
		"Plantronics",
		"PowerGP",
		"Prevedere",
		"primavera",
		"ProductioneerMExt",
		"ProjectIntelligence",
		"QuestionPro",
		"Quosal",
		"Radian6",
		"RiskAssurancePlatform",
		"scoop",
		"ScopevisioPowerBICon",
		"SentryOne",
		"SocialEngagement",
		"SpotlightCloudReports",
		"Timelog",
		"WtsParadigm",
		"Xero",
		"Ziosk",
		"Zuora",
		"AdminInsights",
		"AutoPremium",
		"CapacityMetricsCES",
		"Goals",
		"MetricsCES",
		"MetricsDataConnector",
		"MicrosoftCallQuality",
		"UsageMetricsDataConnector",
		"UsageMetricsCES",
		"ElasticSearch",
		"FabricSql",
	}
}
