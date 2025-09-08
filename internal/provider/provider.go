// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	azlog "github.com/Azure/azure-sdk-for-go/sdk/azcore/log"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoft/fabric-sdk-go/fabric"

	"github.com/microsoft/terraform-provider-fabric/internal/auth"
	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/functions"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
	pclient "github.com/microsoft/terraform-provider-fabric/internal/provider/client"
	pconfig "github.com/microsoft/terraform-provider-fabric/internal/provider/config"
	putils "github.com/microsoft/terraform-provider-fabric/internal/provider/utils"
	"github.com/microsoft/terraform-provider-fabric/internal/services/activator"
	"github.com/microsoft/terraform-provider-fabric/internal/services/apacheairflowjob"
	"github.com/microsoft/terraform-provider-fabric/internal/services/capacity"
	"github.com/microsoft/terraform-provider-fabric/internal/services/connection"
	"github.com/microsoft/terraform-provider-fabric/internal/services/copyjob"
	"github.com/microsoft/terraform-provider-fabric/internal/services/dashboard"
	"github.com/microsoft/terraform-provider-fabric/internal/services/dataflow"
	"github.com/microsoft/terraform-provider-fabric/internal/services/datamart"
	"github.com/microsoft/terraform-provider-fabric/internal/services/datapipeline"
	"github.com/microsoft/terraform-provider-fabric/internal/services/deploymentpipeline"
	"github.com/microsoft/terraform-provider-fabric/internal/services/deploymentpipelinera"
	"github.com/microsoft/terraform-provider-fabric/internal/services/digitaltwinbuilder"
	"github.com/microsoft/terraform-provider-fabric/internal/services/domain"
	"github.com/microsoft/terraform-provider-fabric/internal/services/domainra"
	"github.com/microsoft/terraform-provider-fabric/internal/services/domainwa"
	"github.com/microsoft/terraform-provider-fabric/internal/services/environment"
	"github.com/microsoft/terraform-provider-fabric/internal/services/eventhouse"
	"github.com/microsoft/terraform-provider-fabric/internal/services/eventstream"
	"github.com/microsoft/terraform-provider-fabric/internal/services/eventstreamsourceconnection"
	"github.com/microsoft/terraform-provider-fabric/internal/services/folder"
	"github.com/microsoft/terraform-provider-fabric/internal/services/gateway"
	"github.com/microsoft/terraform-provider-fabric/internal/services/gatewayra"
	"github.com/microsoft/terraform-provider-fabric/internal/services/graphqlapi"
	"github.com/microsoft/terraform-provider-fabric/internal/services/kqldashboard"
	"github.com/microsoft/terraform-provider-fabric/internal/services/kqldatabase"
	"github.com/microsoft/terraform-provider-fabric/internal/services/kqlqueryset"
	"github.com/microsoft/terraform-provider-fabric/internal/services/lakehouse"
	"github.com/microsoft/terraform-provider-fabric/internal/services/lakehousetable"
	"github.com/microsoft/terraform-provider-fabric/internal/services/mirroreddatabase"
	"github.com/microsoft/terraform-provider-fabric/internal/services/mirroredwarehouse"
	"github.com/microsoft/terraform-provider-fabric/internal/services/mlexperiment"
	"github.com/microsoft/terraform-provider-fabric/internal/services/mlmodel"
	"github.com/microsoft/terraform-provider-fabric/internal/services/mounteddatafactory"
	"github.com/microsoft/terraform-provider-fabric/internal/services/notebook"
	"github.com/microsoft/terraform-provider-fabric/internal/services/paginatedreport"
	"github.com/microsoft/terraform-provider-fabric/internal/services/report"
	"github.com/microsoft/terraform-provider-fabric/internal/services/semanticmodel"
	"github.com/microsoft/terraform-provider-fabric/internal/services/shortcut"
	"github.com/microsoft/terraform-provider-fabric/internal/services/sparkcustompool"
	"github.com/microsoft/terraform-provider-fabric/internal/services/sparkenvsettings"
	"github.com/microsoft/terraform-provider-fabric/internal/services/sparkjobdefinition"
	"github.com/microsoft/terraform-provider-fabric/internal/services/sparkwssettings"
	"github.com/microsoft/terraform-provider-fabric/internal/services/sqldatabase"
	"github.com/microsoft/terraform-provider-fabric/internal/services/sqlendpoint"
	"github.com/microsoft/terraform-provider-fabric/internal/services/variablelibrary"
	"github.com/microsoft/terraform-provider-fabric/internal/services/warehouse"
	"github.com/microsoft/terraform-provider-fabric/internal/services/warehousesnapshot"
	"github.com/microsoft/terraform-provider-fabric/internal/services/workspace"
	"github.com/microsoft/terraform-provider-fabric/internal/services/workspacegit"
	"github.com/microsoft/terraform-provider-fabric/internal/services/workspacempe"
	"github.com/microsoft/terraform-provider-fabric/internal/services/workspacera"
)

// Ensure FabricProvider satisfies various provider interfaces.
var (
	_ provider.Provider                       = (*FabricProvider)(nil)
	_ provider.ProviderWithFunctions          = (*FabricProvider)(nil)
	_ provider.ProviderWithEphemeralResources = (*FabricProvider)(nil)
	// _ provider.ProviderWithConfigValidators = (*FabricProvider)(nil)
	// _ provider.ProviderWithValidateConfig   = (*FabricProvider)(nil)
	// _ provider.ProviderWithMetaSchema = (*FabricProvider)(nil).
)

// FabricProvider defines the provider implementation.
type FabricProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version              string
	config               *pconfig.ProviderConfig
	createClient         func(ctx context.Context, cfg *pconfig.ProviderConfig) (*fabric.Client, error)
	affirmProviderConfig func(cfg *pconfig.ProviderConfig)
}

func New(version string) pclient.ProviderWithFabricClient {
	cfg := pconfig.ProviderConfig{}
	cfg.Auth = &auth.Config{}
	cfg.ProviderData = &pconfig.ProviderData{}
	cfg.Endpoint = pconfig.DefaultFabricEndpointURL
	cfg.Timeout, _ = time.ParseDuration(pconfig.DefaultTimeout)
	cfg.Version = version

	return &FabricProvider{
		config:               &cfg,
		version:              version,
		createClient:         createDefaultClient,
		affirmProviderConfig: nil,
	}
}

func NewFunc(version string) func() provider.Provider {
	return func() provider.Provider {
		return New(version)
	}
}

func createDefaultClient(ctx context.Context, cfg *pconfig.ProviderConfig) (*fabric.Client, error) {
	resp, err := auth.NewCredential(*cfg.Auth)
	if err != nil {
		tflog.Error(ctx, "Failed to initialize authentication", map[string]any{"error": err.Error()})

		return nil, err
	}

	tflog.Info(ctx, resp.Info)

	fabricClientOpt := &policy.ClientOptions{}

	// MaxRetries specifies the maximum number of attempts a failed operation will be retried before producing an error.
	// Not really an unlimited cap, but sufficiently large enough to be considered as such.
	fabricClientOpt.Retry.MaxRetries = math.MaxInt32

	// MaxRetryDelay specifies the maximum delay allowed before retrying an operation.
	// A value less than zero means there is no cap.
	fabricClientOpt.Retry.MaxRetryDelay = -1

	ctx, lvl, err := pclient.NewFabricSDKLoggerSubsystem(ctx)
	if err != nil {
		tflog.Error(ctx, "Failed to initialize Microsoft Fabric SDK logger subsystem", map[string]any{"error": err.Error()})

		return nil, err
	}

	if cls := os.Getenv(pclient.AzureSDKLoggingEnvVar); cls == pclient.AzureSDKLoggingAll && lvl != hclog.Off {
		logOptions, err := pclient.ConfigureLoggingOptions(ctx, lvl)
		if err != nil {
			tflog.Error(ctx, "Failed to configure logging options", map[string]any{"error": err.Error()})

			return nil, err
		}

		if logOptions != nil {
			fabricClientOpt.Logging = *logOptions
		}

		azlog.SetListener(func(ev azlog.Event, msg string) {
			tflog.SubsystemTrace(ctx, pclient.FabricSDKLoggerName, "SDK", map[string]any{
				"event":   ev,
				"message": msg,
			})
		})
	}

	perCallPolicies := make([]policy.Policy, 0)
	perCallPolicies = append(perCallPolicies, pclient.WithUserAgent(pclient.BuildUserAgent(cfg.TerraformVersion, fabric.Version, cfg.Version, cfg.PartnerID, cfg.DisableTerraformPartnerID)))
	fabricClientOpt.PerCallPolicies = perCallPolicies

	client, err := fabric.NewClient(resp.Cred, &cfg.Endpoint, fabricClientOpt)
	if err != nil {
		tflog.Error(ctx, "Failed to initialize Microsoft Fabric client", map[string]any{"error": err.Error()})

		return nil, err
	}

	return client, nil
}

func (p *FabricProvider) ConfigureCreateClient(createFabricClientFunction func(ctx context.Context, cfg *pconfig.ProviderConfig) (*fabric.Client, error)) {
	p.createClient = createFabricClientFunction
}

func (p *FabricProvider) ConfigureAffirmProviderConfig(affirmProviderConfigFunction func(cfg *pconfig.ProviderConfig)) {
	p.affirmProviderConfig = affirmProviderConfigFunction
}

func (p *FabricProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "fabric"
	resp.Version = p.version
}

func (p *FabricProvider) Schema(ctx context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	tflog.Debug(ctx, "Schema request received")

	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"timeout": schema.StringAttribute{
				MarkdownDescription: "Default timeout for all requests. It can be overridden at any Resource/Data-Source\n" +
					"   A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as `30s` or `2h45m`. Valid time units are \"s\" (seconds), \"m\" (minutes), \"h\" (hours)\n" +
					"   If not set, the default timeout is `" + pconfig.DefaultTimeout + "`.",
				Optional:   true,
				CustomType: timetypes.GoDurationType{},
			},

			// Fabric specific fields
			"endpoint": schema.StringAttribute{
				MarkdownDescription: "The Endpoint of the Microsoft Fabric API.",
				Optional:            true,
				CustomType:          customtypes.URLType{},
			},

			// Azure specific fields
			"environment": schema.StringAttribute{
				MarkdownDescription: "The cloud environment which should be used. Possible values are 'public', 'usgovernment' and 'china'. Defaults to 'public'",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive("public", "usgovernment", "china"),
				},
			},
			"tenant_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the Microsoft Entra ID tenant that Fabric API uses to authenticate with.",
				Optional:            true,
				CustomType:          customtypes.UUIDType{},
			},
			"auxiliary_tenant_ids": schema.SetAttribute{
				MarkdownDescription: "The Auxiliary Tenant IDs which should be used.",
				ElementType:         customtypes.UUIDType{},
				Optional:            true,
				Validators: []validator.Set{
					setvalidator.SizeAtMost(3),
				},
			},

			// Client ID specific fields
			"client_id": schema.StringAttribute{
				MarkdownDescription: "The Client ID of the app registration.",
				Optional:            true,
				CustomType:          customtypes.UUIDType{},
			},
			"client_id_file_path": schema.StringAttribute{
				MarkdownDescription: "The path to a file containing the Client ID which should be used.",
				Optional:            true,
			},

			// Client Secret specific fields
			"client_secret": schema.StringAttribute{
				MarkdownDescription: "The Client Secret of the app registration. For use when authenticating as a Service Principal using a Client Secret.",
				Optional:            true,
				Sensitive:           true,
			},
			"client_secret_file_path": schema.StringAttribute{
				MarkdownDescription: "The path to a file containing the Client Secret which should be used. For use when authenticating as a Service Principal using a Client Secret.",
				Optional:            true,
			},

			// Client Certificate specific fields
			"client_certificate": schema.StringAttribute{
				MarkdownDescription: "Base64 encoded PKCS#12 certificate bundle. For use when authenticating as a Service Principal using a Client Certificate.",
				Optional:            true,
				Sensitive:           true,
			},
			"client_certificate_file_path": schema.StringAttribute{
				MarkdownDescription: "The path to the Client Certificate associated with the Service Principal for use when authenticating as a Service Principal using a Client Certificate.",
				Optional:            true,
			},
			"client_certificate_password": schema.StringAttribute{
				MarkdownDescription: "The password associated with the Client Certificate. For use when authenticating as a Service Principal using a Client Certificate.",
				Optional:            true,
				Sensitive:           true,
			},

			// OIDC specific fields
			"use_oidc": schema.BoolAttribute{
				MarkdownDescription: "Allow OpenID Connect to be used for authentication.",
				Optional:            true,
			},
			"oidc_request_token": schema.StringAttribute{
				MarkdownDescription: "The bearer token for the request to the OIDC provider. For use when authenticating as a Service Principal using OpenID Connect.",
				Optional:            true,
				Sensitive:           true,
			},
			"oidc_request_url": schema.StringAttribute{
				MarkdownDescription: "The URL for the OIDC provider from which to request an ID token. For use when authenticating as a Service Principal using OpenID Connect.",
				Optional:            true,
			},
			"oidc_token": schema.StringAttribute{
				MarkdownDescription: "The OIDC token for use when authenticating as a Service Principal using OpenID Connect.",
				Optional:            true,
				Sensitive:           true,
			},
			"oidc_token_file_path": schema.StringAttribute{
				MarkdownDescription: "The path to a file containing an OIDC token for use when authenticating as a Service Principal using OpenID Connect.",
				Optional:            true,
			},
			"azure_devops_service_connection_id": schema.StringAttribute{
				MarkdownDescription: "The Azure DevOps Service Connection ID that uses Workload Identity Federation.",
				Optional:            true,
			},

			// Use Azure CLI for auth
			"use_cli": schema.BoolAttribute{
				MarkdownDescription: "Allow Azure CLI to be used for authentication.",
				Optional:            true,
			},

			// Use Azure Developer CLI for auth
			"use_dev_cli": schema.BoolAttribute{
				MarkdownDescription: "Allow Azure Developer CLI to be used for authentication.",
				Optional:            true,
			},

			// Use Managed Service Identity for auth
			"use_msi": schema.BoolAttribute{
				MarkdownDescription: "Allow Managed Service Identity (MSI) to be used for authentication.",
				Optional:            true,
			},

			// Use preview features
			"preview": schema.BoolAttribute{
				MarkdownDescription: "Enable preview mode to use preview features.",
				Optional:            true,
			},

			// Telemetry
			"partner_id": schema.StringAttribute{
				MarkdownDescription: "A GUID/UUID that is [registered](https://learn.microsoft.com/partner-center/marketplace-offers/azure-partner-customer-usage-attribution#register-guids-and-offers) with Microsoft to facilitate partner resource usage attribution.",
				Optional:            true,
				CustomType:          customtypes.UUIDType{},
			},
			"disable_terraform_partner_id": schema.BoolAttribute{
				MarkdownDescription: "Disable sending the Terraform Partner ID if a custom `partner_id` isn't specified, which allows Microsoft to better understand the usage of Terraform. The Partner ID does not give HashiCorp any direct access to usage information. This can also be sourced from the `FABRIC_DISABLE_TERRAFORM_PARTNER_ID` environment variable. Defaults to `false`.",
				Optional:            true,
			},
		},
	}
}

// func (p *FabricProvider) ValidateConfig(ctx context.Context, req provider.ValidateConfigRequest, resp *provider.ValidateConfigResponse) {
// 	var config pconfig.ProviderConfigModel

// 	if resp.Diagnostics.Append(req.Config.Get(ctx, &config)...); resp.Diagnostics.HasError() {
// 		return
// 	}
// }

// func (p *FabricProvider) ConfigValidators(ctx context.Context) []provider.ConfigValidator {
// 	return []provider.ConfigValidator{}
// }

func (p *FabricProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Debug(ctx, "Microsoft Fabric Terraform Provider configuration started")
	tflog.Info(ctx, "Initializing Terraform Provider for Microsoft Fabric")

	var config pconfig.ProviderConfigModel

	if resp.Diagnostics.Append(req.Config.Get(ctx, &config)...); resp.Diagnostics.HasError() {
		return
	}

	ctx = p.setConfig(ctx, &config, resp)
	tflog.Debug(ctx, "Setting configuration")

	if resp.Diagnostics.HasError() {
		return
	}

	p.mapConfig(ctx, &config, resp)
	tflog.Debug(ctx, "Mapping configuration")

	p.config.TerraformVersion = req.TerraformVersion

	if resp.Diagnostics.HasError() {
		return
	}

	if p.affirmProviderConfig != nil {
		p.affirmProviderConfig(p.config)
	}

	tflog.Debug(ctx, "Creating Microsoft Fabric client")

	client, err := p.createClient(ctx, p.config)
	if err != nil {
		tflog.Error(ctx, "Error configuring Microsoft Fabric client", map[string]any{"error": err})

		return
	}

	p.config.FabricClient = client

	tflog.Debug(ctx, "Assigning Microsoft Fabric client to DataSourceData")

	resp.DataSourceData = p.config.ProviderData

	tflog.Debug(ctx, "Assigning Microsoft Fabric client to ResourceData")

	resp.ResourceData = p.config.ProviderData

	tflog.Debug(ctx, "Assigning Microsoft Fabric client to EphemeralResourceData")

	resp.EphemeralResourceData = p.config.ProviderData

	tflog.Info(ctx, "Configured Microsoft Fabric client", map[string]any{"success": true})
}

func (p *FabricProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		apacheairflowjob.NewResourceApacheAirflowJob,
		copyjob.NewResourceCopyJob,
		dataflow.NewResourceDataflow,
		datapipeline.NewResourceDataPipeline,
		digitaltwinbuilder.NewResourceDigitalTwinBuilder,
		domain.NewResourceDomain,
		domainra.NewResourceDomainRoleAssignments,
		domainwa.NewResourceDomainWorkspaceAssignments,
		connection.NewResourceConnection,
		deploymentpipeline.NewResourceDeploymentPipeline,
		deploymentpipelinera.NewResourceDeploymentPipelineRoleAssignment,
		func() resource.Resource { return environment.NewResourceEnvironment(ctx) },
		func() resource.Resource { return eventhouse.NewResourceEventhouse(ctx) },
		eventstream.NewResourceEventstream,
		folder.NewResourceFolder,
		gateway.NewResourceGateway,
		gatewayra.NewResourceGatewayRoleAssignment,
		graphqlapi.NewResourceGraphQLApi,
		kqldashboard.NewResourceKQLDashboard,
		kqldatabase.NewResourceKQLDatabase,
		kqlqueryset.NewResourceKQLQueryset,
		func() resource.Resource { return lakehouse.NewResourceLakehouse(ctx) },
		func() resource.Resource { return mirroreddatabase.NewResourceMirroredDatabase(ctx) },
		mounteddatafactory.NewResourceMountedDataFactory,
		mlexperiment.NewResourceMLExperiment,
		mlmodel.NewResourceMLModel,
		shortcut.NewResourceShortcut,
		notebook.NewResourceNotebook,
		activator.NewResourceActivator,
		report.NewResourceReport,
		semanticmodel.NewResourceSemanticModel,
		sparkcustompool.NewResourceSparkCustomPool,
		sparkenvsettings.NewResourceSparkEnvironmentSettings,
		sparkwssettings.NewResourceSparkWorkspaceSettings,
		sparkjobdefinition.NewResourceSparkJobDefinition,
		sqldatabase.NewResourceSQLDatabase,
		variablelibrary.NewResourceVariableLibrary,
		warehouse.NewResourceWarehouse,
		warehousesnapshot.NewResourceWarehouseSnapshot,
		workspace.NewResourceWorkspace,
		workspacera.NewResourceWorkspaceRoleAssignment,
		workspacegit.NewResourceWorkspaceGit,
		workspacempe.NewResourceWorkspaceManagedPrivateEndpoint,
	}
}

func (p *FabricProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		apacheairflowjob.NewDataSourceApacheAirflowJob,
		apacheairflowjob.NewDataSourceApacheAirflowJobs,
		capacity.NewDataSourceCapacity,
		capacity.NewDataSourceCapacities,
		connection.NewDataSourceConnection,
		connection.NewDataSourceConnections,
		copyjob.NewDataSourceCopyJob,
		copyjob.NewDataSourceCopyJobs,
		dashboard.NewDataSourceDashboards,
		dataflow.NewDataSourceDataflow,
		dataflow.NewDataSourceDataflows,
		datapipeline.NewDataSourceDataPipeline,
		datapipeline.NewDataSourceDataPipelines,
		datamart.NewDataSourceDatamarts,
		digitaltwinbuilder.NewDataSourceDigitalTwinBuilder,
		digitaltwinbuilder.NewDataSourceDigitalTwinBuilders,
		deploymentpipeline.NewDataSourceDeploymentPipeline,
		deploymentpipeline.NewDataSourceDeploymentPipelines,
		deploymentpipelinera.NewDataSourceDeploymentPipelineRoleAssignments,
		domain.NewDataSourceDomain,
		domain.NewDataSourceDomains,
		domainwa.NewDataSourceDomainWorkspaceAssignments,
		func() datasource.DataSource { return environment.NewDataSourceEnvironment(ctx) },
		func() datasource.DataSource { return environment.NewDataSourceEnvironments(ctx) },
		func() datasource.DataSource { return eventhouse.NewDataSourceEventhouse(ctx) },
		func() datasource.DataSource { return eventhouse.NewDataSourceEventhouses(ctx) },
		eventstream.NewDataSourceEventstream,
		eventstream.NewDataSourceEventstreams,
		eventstreamsourceconnection.NewDataSourceEventstreamSourceConnection,
		folder.NewDataSourceFolder,
		folder.NewDataSourceFolders,
		gateway.NewDataSourceGateway,
		gateway.NewDataSourceGateways,
		gatewayra.NewDataSourceGatewayRoleAssignment,
		gatewayra.NewDataSourceGatewayRoleAssignments,
		graphqlapi.NewDataSourceGraphQLApi,
		graphqlapi.NewDataSourceGraphQLApis,
		kqldashboard.NewDataSourceKQLDashboard,
		kqldashboard.NewDataSourceKQLDashboards,
		kqldatabase.NewDataSourceKQLDatabase,
		kqldatabase.NewDataSourceKQLDatabases,
		kqlqueryset.NewDataSourceKQLQueryset,
		kqlqueryset.NewDataSourceKQLQuerysets,
		func() datasource.DataSource { return lakehouse.NewDataSourceLakehouse(ctx) },
		func() datasource.DataSource { return lakehouse.NewDataSourceLakehouses(ctx) },
		lakehousetable.NewDataSourceLakehouseTable,
		lakehousetable.NewDataSourceLakehouseTables,
		func() datasource.DataSource { return mirroreddatabase.NewDataSourceMirroredDatabase(ctx) },
		func() datasource.DataSource { return mirroreddatabase.NewDataSourceMirroredDatabases(ctx) },
		mirroredwarehouse.NewDataSourceMirroredWarehouses,
		mlexperiment.NewDataSourceMLExperiment,
		mlexperiment.NewDataSourceMLExperiments,
		mlmodel.NewDataSourceMLModel,
		mlmodel.NewDataSourceMLModels,
		mounteddatafactory.NewDataSourceMountedDataFactory,
		mounteddatafactory.NewDataSourceMountedDataFactories,
		notebook.NewDataSourceNotebook,
		notebook.NewDataSourceNotebooks,
		shortcut.NewDataSourceShortcut,
		shortcut.NewDataSourceShortcuts,
		paginatedreport.NewDataSourcePaginatedReports,
		activator.NewDataSourceActivator,
		activator.NewDataSourceActivators,
		report.NewDataSourceReport,
		report.NewDataSourceReports,
		semanticmodel.NewDataSourceSemanticModel,
		semanticmodel.NewDataSourceSemanticModels,
		sparkcustompool.NewDataSourceSparkCustomPool,
		sparkenvsettings.NewDataSourceSparkEnvironmentSettings,
		sparkwssettings.NewDataSourceSparkWorkspaceSettings,
		sparkjobdefinition.NewDataSourceSparkJobDefinition,
		sparkjobdefinition.NewDataSourceSparkJobDefinitions,
		sqldatabase.NewDataSourceSQLDatabase,
		sqldatabase.NewDataSourceSQLDatabases,
		sqlendpoint.NewDataSourceSQLEndpoints,
		variablelibrary.NewDataSourceVariableLibrary,
		variablelibrary.NewDataSourceVariableLibraries,
		warehouse.NewDataSourceWarehouse,
		warehouse.NewDataSourceWarehouses,
		warehousesnapshot.NewDataSourceWarehouseSnapshot,
		warehousesnapshot.NewDataSourceWarehouseSnapshots,
		workspace.NewDataSourceWorkspace,
		workspace.NewDataSourceWorkspaces,
		workspacera.NewDataSourceWorkspaceRoleAssignment,
		workspacera.NewDataSourceWorkspaceRoleAssignments,
		workspacegit.NewDataSourceWorkspaceGit,
		workspacempe.NewDataSourceWorkspaceManagedPrivateEndpoint,
		workspacempe.NewDataSourceWorkspaceManagedPrivateEndpoints,
	}
}

func (p *FabricProvider) EphemeralResources(_ context.Context) []func() ephemeral.EphemeralResource {
	return []func() ephemeral.EphemeralResource{
		eventstreamsourceconnection.NewEphemeralResourceEventstreamSourceConnection,
	}
}

func (p *FabricProvider) Functions(_ context.Context) []func() function.Function {
	return []func() function.Function{
		functions.NewFunctionContentDecode,
	}
}

func (p *FabricProvider) setConfig(ctx context.Context, config *pconfig.ProviderConfigModel, resp *provider.ConfigureResponse) context.Context {
	timeout, diags := config.Timeout.ToStringValue(ctx)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return ctx
	}

	config.Timeout = timetypes.NewGoDurationValueFromStringMust(putils.GetStringValue(timeout, pconfig.GetEnvVarsTimeout(), pconfig.DefaultTimeout).ValueString())
	if config.Timeout.IsUnknown() || config.Timeout.IsNull() {
		resp.Diagnostics.AddAttributeError(
			path.Root("timeout"),
			common.ErrorInvalidConfig,
			"The provider cannot create the Microsoft Fabric API client as there is an unknown configuration value for the timeout. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the "+utils.ConvertStringSlicesToString(pconfig.GetEnvVarsTimeout(), false, false)+" environment variables.",
		)

		return ctx
	}

	ctx = tflog.SetField(ctx, "timeout", config.Timeout)

	endpoint, diags := config.Endpoint.ToStringValue(ctx)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return ctx
	}

	config.Endpoint = customtypes.NewURLValue(putils.GetStringValue(endpoint, pconfig.GetEnvVarsEndpoint(), pconfig.DefaultFabricEndpointURL).ValueString())
	if config.Endpoint.IsUnknown() || config.Endpoint.IsNull() {
		resp.Diagnostics.AddAttributeError(
			path.Root("endpoint"),
			common.ErrorInvalidConfig,
			"The provider cannot create the Microsoft Fabric API client as there is an unknown configuration value for the Microsoft Fabric API endpoint. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the "+utils.ConvertStringSlicesToString(pconfig.GetEnvVarsEndpoint(), false, false)+" environment variables.",
		)

		return ctx
	}

	ctx = tflog.SetField(ctx, "endpoint", config.Endpoint)

	config.Environment = putils.GetStringValue(config.Environment, pconfig.GetEnvVarsEnvironment(), "public")
	environmentAllowedValues := map[string]bool{
		"public":       true,
		"usgovernment": true,
		"china":        true,
	}

	if _, ok := environmentAllowedValues[config.Environment.ValueString()]; !ok {
		resp.Diagnostics.AddAttributeError(
			path.Root("environment"),
			common.ErrorInvalidConfig,
			"The provider cannot create the Microsoft Fabric API client as there is an unknown configuration value for the Cloud environment. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the "+strings.Join(pconfig.GetEnvVarsEnvironment(), ",")+" environment variables.",
		)

		return ctx
	}

	ctx = tflog.SetField(ctx, "environment", config.Environment.ValueString())

	tenantID, diags := config.TenantID.ToStringValue(ctx)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return ctx
	}

	config.TenantID = customtypes.NewUUIDValue(putils.GetStringValue(tenantID, pconfig.GetEnvVarsTenantID(), "").ValueString())
	ctx = tflog.SetField(ctx, "tenant_id", config.TenantID.ValueString())

	config.AuxiliaryTenantIDs = putils.GetListStringValues(config.AuxiliaryTenantIDs, pconfig.GetEnvVarsAuxiliaryTenantIDs(), []string{})
	ctx = tflog.SetField(ctx, "auxiliary_tenant_ids", config.AuxiliaryTenantIDs.String())

	clientID, diags := config.ClientID.ToStringValue(ctx)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return ctx
	}

	config.ClientID = customtypes.NewUUIDValue(putils.GetStringValue(clientID, pconfig.GetEnvVarsClientID(), "").ValueString())
	ctx = tflog.SetField(ctx, "client_id", config.ClientID.ValueString())

	config.ClientIDFilePath = putils.GetStringValue(config.ClientIDFilePath, pconfig.GetEnvVarsClientIDFilePath(), "")
	ctx = tflog.SetField(ctx, "client_id_file_path", config.ClientIDFilePath.ValueString())

	config.ClientSecret = putils.GetStringValue(config.ClientSecret, pconfig.GetEnvVarsClientSecret(), "")
	ctx = tflog.SetField(ctx, "client_secret", config.ClientSecret.ValueString())
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "client_secret")

	config.ClientSecretFilePath = putils.GetStringValue(config.ClientSecretFilePath, pconfig.GetEnvVarsClientSecretFilePath(), "")
	ctx = tflog.SetField(ctx, "client_secret_file_path", config.ClientSecretFilePath.ValueString())

	config.ClientCertificate = putils.GetStringValue(config.ClientCertificate, pconfig.GetEnvVarsClientCertificate(), "")
	ctx = tflog.SetField(ctx, "client_certificate", config.ClientCertificate.ValueString())
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "client_certificate")

	config.ClientCertificateFilePath = putils.GetStringValue(config.ClientCertificateFilePath, pconfig.GetEnvVarsClientCertificateFilePath(), "")
	ctx = tflog.SetField(ctx, "client_certificate_file_path", config.ClientCertificateFilePath.ValueString())

	config.ClientCertificatePassword = putils.GetStringValue(config.ClientCertificatePassword, pconfig.GetEnvVarsClientCertificatePassword(), "")
	ctx = tflog.SetField(ctx, "client_certificate_password", config.ClientCertificatePassword.ValueString())
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "client_certificate_password")

	config.OIDCRequestURL = putils.GetStringValue(config.OIDCRequestURL, pconfig.GetEnvVarsOIDCRequestURL(), "")
	ctx = tflog.SetField(ctx, "oidc_request_url", config.OIDCRequestURL.ValueString())
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "oidc_request_url")

	config.OIDCRequestToken = putils.GetStringValue(config.OIDCRequestToken, pconfig.GetEnvVarsOIDCRequestToken(), "")
	ctx = tflog.SetField(ctx, "oidc_request_token", config.OIDCRequestToken.ValueString())
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "oidc_request_token")

	config.OIDCToken = putils.GetStringValue(config.OIDCToken, pconfig.GetEnvVarsOIDCToken(), "")
	ctx = tflog.SetField(ctx, "oidc_token", config.OIDCToken.ValueString())
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "oidc_token")

	config.OIDCTokenFilePath = putils.GetStringValue(config.OIDCTokenFilePath, pconfig.GetEnvVarsOIDCTokenFilePath(), "")
	ctx = tflog.SetField(ctx, "oidc_token_file_path", config.OIDCTokenFilePath.ValueString())

	config.AzureDevOpsServiceConnectionID = putils.GetStringValue(config.AzureDevOpsServiceConnectionID, pconfig.GetEnvVarsAzureDevOpsServiceConnectionID(), "")
	ctx = tflog.SetField(ctx, "azure_devops_service_connection_id", config.AzureDevOpsServiceConnectionID.ValueString())

	config.UseOIDC = putils.GetBoolValue(config.UseOIDC, pconfig.GetEnvVarsUseOIDC(), false)
	ctx = tflog.SetField(ctx, "use_oidc", config.UseOIDC.ValueBool())

	config.UseMSI = putils.GetBoolValue(config.UseMSI, pconfig.GetEnvVarsUseMSI(), false)
	ctx = tflog.SetField(ctx, "use_msi", config.UseMSI.ValueBool())

	config.UseDevCLI = putils.GetBoolValue(config.UseDevCLI, pconfig.GetEnvVarsUseDevCLI(), false)
	ctx = tflog.SetField(ctx, "use_dev_cli", config.UseDevCLI.ValueBool())

	config.UseCLI = putils.GetBoolValue(config.UseCLI, pconfig.GetEnvVarsUseCLI(), false)
	ctx = tflog.SetField(ctx, "use_cli", config.UseCLI.ValueBool())

	trueCount := 0
	if config.UseOIDC.ValueBool() {
		trueCount++
	}

	if config.UseMSI.ValueBool() {
		trueCount++
	}

	if config.UseDevCLI.ValueBool() {
		trueCount++
	}

	if config.UseCLI.ValueBool() {
		trueCount++
	}

	if trueCount > 1 {
		resp.Diagnostics.AddError(
			common.ErrorAttComboInvalid,
			"Only one of 'use_oidc', 'use_msi', 'use_dev_cli', and 'use_cli' can be true at a time.\n"+
				"Current configuration:\n"+
				"\tuse_cli: "+config.UseCLI.String()+"\n"+
				"\tuse_dev_cli: "+config.UseDevCLI.String()+"\n"+
				"\tuse_msi: "+config.UseMSI.String()+"\n"+
				"\tuse_oidc: "+config.UseOIDC.String(),
		)

		return ctx
	}

	config.Preview = putils.GetBoolValue(config.Preview, pconfig.GetEnvVarsPreview(), false)
	ctx = tflog.SetField(ctx, "preview", config.Preview.ValueBool())

	partnerID, diags := config.PartnerID.ToStringValue(ctx)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return ctx
	}

	config.PartnerID = customtypes.NewUUIDValue(putils.GetStringValue(partnerID, pconfig.GetEnvVarsPartnerID(), "").ValueString())
	ctx = tflog.SetField(ctx, "partner_id", config.PartnerID.ValueString())

	config.DisableTerraformPartnerID = putils.GetBoolValue(config.DisableTerraformPartnerID, pconfig.GetEnvVarsDisableTerraformPartnerID(), false)
	ctx = tflog.SetField(ctx, "disable_terraform_partner_id", config.DisableTerraformPartnerID.ValueBool())

	return ctx
}

func (p *FabricProvider) mapConfig(ctx context.Context, config *pconfig.ProviderConfigModel, resp *provider.ConfigureResponse) {
	// Map the provider configuration model the configuration provider configuration
	var err error

	timeout, diags := config.Timeout.ValueGoDuration()
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	p.config.Timeout = timeout

	p.config.Endpoint = config.Endpoint.ValueString()

	switch strings.ToLower(config.Environment.ValueString()) {
	case "usgovernment":
		p.config.Auth.Environment = cloud.AzureGovernment
	case "china":
		p.config.Auth.Environment = cloud.AzureChina
	default:
		p.config.Auth.Environment = cloud.AzurePublic
	}

	p.config.Auth.TenantID = config.TenantID.ValueString()

	var auxiliaryTenantIDs []string

	resp.Diagnostics.Append(config.AuxiliaryTenantIDs.ElementsAs(ctx, &auxiliaryTenantIDs, true)...)
	p.config.Auth.AuxiliaryTenantIDs = auxiliaryTenantIDs

	clientID, diags := config.ClientID.ToStringValue(ctx)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	p.config.Auth.ClientID, err = putils.GetValueOrFileValue("client_id", "client_id_file_path", clientID, config.ClientIDFilePath)
	if err != nil {
		resp.Diagnostics.AddError(common.ErrorInvalidValue, err.Error())

		return
	}

	p.config.Auth.ClientSecret, err = putils.GetValueOrFileValue("client_secret", "client_id_file_path", config.ClientSecret, config.ClientSecretFilePath)
	if err != nil {
		resp.Diagnostics.AddError(common.ErrorInvalidValue, err.Error())

		return
	}

	clientCertificateRaw, err := putils.GetCertOrFileCert("client_certificate", "client_certificate_file_path", config.ClientCertificate, config.ClientCertificateFilePath)
	if err != nil {
		resp.Diagnostics.AddError(common.ErrorInvalidValue, err.Error())

		return
	}

	if clientCertificateRaw != "" {
		cert, key, err := auth.ConvertBase64ToCert(clientCertificateRaw, config.ClientCertificatePassword.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(common.ErrorInvalidValue, err.Error())

			return
		}

		p.config.Auth.ClientCertificate = cert
		p.config.Auth.ClientCertificateKey = key
	}

	p.config.Auth.OIDC.RequestToken = config.OIDCRequestToken.ValueString()
	p.config.Auth.OIDC.RequestURL = config.OIDCRequestURL.ValueString()

	p.config.Auth.OIDC.Token, err = putils.GetValueOrFileValue("oidc_token", "oidc_token_file_path", config.OIDCToken, config.OIDCTokenFilePath)
	if err != nil {
		resp.Diagnostics.AddError(common.ErrorInvalidValue, err.Error())

		return
	}

	p.config.Auth.OIDC.AzureDevOpsServiceConnectionID = config.AzureDevOpsServiceConnectionID.ValueString()
	p.config.Auth.UseOIDC = config.UseOIDC.ValueBool()
	p.config.Auth.UseMSI = config.UseMSI.ValueBool()
	p.config.Auth.UseDevCLI = config.UseDevCLI.ValueBool()
	p.config.Auth.UseCLI = config.UseCLI.ValueBool()

	p.validateConfigAuthOIDC(resp)
	p.validateConfigAuthMSI(resp)
	p.validateConfigAuthCertificate(resp)
	p.validateConfigAuthSecret(resp)

	p.config.Preview = config.Preview.ValueBool()
	p.config.PartnerID = config.PartnerID.ValueString()
	p.config.DisableTerraformPartnerID = config.DisableTerraformPartnerID.ValueBool()
}

func (p *FabricProvider) validateConfigAuthOIDC(resp *provider.ConfigureResponse) {
	if p.config.Auth.UseOIDC {
		infoAuthType := " when using OpenID Connect (OIDC) authentication."

		if p.config.Auth.TenantID == "" {
			resp.Diagnostics.AddError(
				common.ErrorInvalidConfig,
				"'tenant_id' is required"+infoAuthType,
			)

			return
		}

		if p.config.Auth.ClientID == "" {
			resp.Diagnostics.AddError(
				common.ErrorInvalidConfig,
				"Client ID is required"+infoAuthType+"\n"+
					"Please provide a valid 'client_id' or 'client_id_file_path' in the provider configuration.",
			)

			return
		}

		if p.config.Auth.ClientSecret != "" {
			resp.Diagnostics.AddError(
				common.ErrorInvalidConfig,
				"'client_secret' is not accepted"+infoAuthType,
			)

			return
		}

		if p.config.Auth.OIDC.AzureDevOpsServiceConnectionID != "" && p.config.Auth.OIDC.RequestToken == "" {
			resp.Diagnostics.AddError(
				common.ErrorInvalidConfig,
				"'oidc_request_token' is required when 'azure_devops_service_connection_id' is provided"+infoAuthType,
			)
		}

		if p.config.Auth.OIDC.AzureDevOpsServiceConnectionID == "" && (p.config.Auth.OIDC.Token == "" && (p.config.Auth.OIDC.RequestURL == "" || p.config.Auth.OIDC.RequestToken == "")) {
			maskValue := (func(v string) string {
				if v != "" {
					return "***"
				}

				return ""
			})

			resp.Diagnostics.AddError(
				common.ErrorInvalidConfig,
				"either 'oidc_token' or both 'oidc_request_url' and 'oidc_request_token' must be provided"+infoAuthType+"\n"+
					"Current configuration:\n"+
					"\tuse_oidc: "+strconv.FormatBool(p.config.Auth.UseOIDC)+"\n"+
					"\toidc_token: "+maskValue(p.config.Auth.OIDC.Token)+"\n"+
					"\toidc_request_url: "+maskValue(p.config.Auth.OIDC.RequestURL)+"\n"+
					"\toidc_request_token: "+maskValue(p.config.Auth.OIDC.RequestToken),
			)

			return
		}
	}
}

// validateConfigAuthMSI validates the configuration for Managed Service Identity (MSI) authentication.
func (p *FabricProvider) validateConfigAuthMSI(resp *provider.ConfigureResponse) {
	if p.config.Auth.UseMSI {
		infoAuthType := " when using Managed Service Identity (MSI) authentication."

		if p.config.Auth.TenantID == "" {
			resp.Diagnostics.AddError(
				common.ErrorInvalidConfig,
				"'tenant_id' is required"+infoAuthType,
			)

			return
		}

		if p.config.Auth.ClientSecret != "" {
			resp.Diagnostics.AddError(
				common.ErrorInvalidConfig,
				"'client_secret' is not accepted"+infoAuthType,
			)

			return
		}
	}
}

func (p *FabricProvider) validateConfigAuthCertificate(resp *provider.ConfigureResponse) {
	if p.config.Auth.ClientCertificate != nil {
		infoAuthType := " when using Service Principal with Certificate authentication."

		if p.config.Auth.TenantID == "" {
			resp.Diagnostics.AddError(
				common.ErrorInvalidConfig,
				"'tenant_id' is required"+infoAuthType,
			)

			return
		}

		if p.config.Auth.ClientID == "" {
			resp.Diagnostics.AddError(
				common.ErrorInvalidConfig,
				"Client ID is required"+infoAuthType+"\n"+
					"Please provide a valid 'client_id' or 'client_id_file_path' in the provider configuration.",
			)

			return
		}

		if p.config.Auth.ClientSecret != "" {
			resp.Diagnostics.AddError(
				common.ErrorInvalidConfig,
				"'client_secret' is not accepted"+infoAuthType,
			)

			return
		}
	}
}

// validateConfigAuthSecret validates the configuration for Service Principal authentication with Secret.
func (p *FabricProvider) validateConfigAuthSecret(resp *provider.ConfigureResponse) {
	if p.config.Auth.ClientSecret != "" {
		infoAuthType := " when using Service Principal with Secret authentication."

		if p.config.Auth.TenantID == "" {
			resp.Diagnostics.AddError(
				common.ErrorInvalidConfig,
				"'tenant_id' is required"+infoAuthType,
			)

			return
		}

		if p.config.Auth.ClientID == "" {
			resp.Diagnostics.AddError(
				common.ErrorInvalidConfig,
				"Client ID is required"+infoAuthType+"\n"+
					"Please provide a valid 'client_id' or 'client_id_file_path' in the provider configuration.",
			)

			return
		}
	}
}
