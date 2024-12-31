// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package connection

import (
	"context"
	"fmt"

	azto "github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	supertypes "github.com/FrangipaneTeam/terraform-plugin-framework-supertypes"
	superstringvalidator "github.com/FrangipaneTeam/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
	pconfig "github.com/microsoft/terraform-provider-fabric/internal/provider/config"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.ResourceWithConfigure  = (*resourceConnection)(nil)
	_ resource.ResourceWithModifyPlan = (*resourceConnection)(nil)
)

type resourceConnection struct {
	pConfigData *pconfig.ProviderData
	client      *fabcore.ConnectionsClient
}

func NewResourceConnection() resource.Resource {
	return &resourceConnection{}
}

func (r *resourceConnection) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + ItemTFName
}

func (r *resourceConnection) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manage a Fabric " + ItemName + ".\n\n" +
			"See [" + ItemName + "](" + ItemDocsURL + ") for more information.\n\n" +
			ItemDocsSPNSupport,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The object ID of the connection.",
				Computed:            true,
				CustomType:          customtypes.UUIDType{},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "The display name of the connection.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(123),
				},
			},
			"connectivity_type": schema.StringAttribute{
				MarkdownDescription: "The connectivity type of the connection. Accepted values: " + utils.ConvertStringSlicesToString(fabcore.PossibleConnectivityTypeValues(), true, true),
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf(utils.ConvertEnumsToStringSlices(fabcore.PossibleConnectivityTypeValues(), false)...),
				},
			},
			"privacy_level": schema.StringAttribute{
				MarkdownDescription: "The privacy level of the connection. Accepted values: " + utils.ConvertStringSlicesToString(fabcore.PossiblePrivacyLevelValues(), true, true),
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf(utils.ConvertEnumsToStringSlices(fabcore.PossiblePrivacyLevelValues(), false)...),
				},
			},
			"gateway_id": schema.StringAttribute{
				MarkdownDescription: "The gateway object ID of the connection.",
				Optional:            true,
				CustomType:          customtypes.UUIDType{},
				Validators: []validator.String{
					superstringvalidator.RequireIfAttributeIsOneOf(path.MatchRoot("connectivity_type"),
						[]attr.Value{
							types.StringValue(string(fabcore.ConnectivityTypeVirtualNetworkGateway)),
							types.StringValue(string(fabcore.ConnectivityTypeOnPremisesGateway)),
							types.StringValue(string(fabcore.ConnectivityTypeOnPremisesGatewayPersonal)),
						}),
					superstringvalidator.NullIfAttributeIsOneOf(path.MatchRoot("connectivity_type"),
						[]attr.Value{
							types.StringValue(string(fabcore.ConnectivityTypeAutomatic)),
							types.StringValue(string(fabcore.ConnectivityTypeNone)),
							types.StringValue(string(fabcore.ConnectivityTypePersonalCloud)),
							types.StringValue(string(fabcore.ConnectivityTypeShareableCloud)),
						}),
				},
			},
			"connection_details": schema.SingleNestedAttribute{
				MarkdownDescription: "The connection details of the connection.",
				Required:            true,
				CustomType:          supertypes.NewSingleNestedObjectTypeOf[rsConnectionDetailsModel](ctx),
				Attributes: map[string]schema.Attribute{
					"path": schema.StringAttribute{
						MarkdownDescription: "The path of the connection.",
						Computed:            true,
					},
					"type": schema.StringAttribute{
						MarkdownDescription: "The type of the connection. Accepted values: " + utils.ConvertStringSlicesToString(possibleConnectionTypeValues(), true, true),
						Required:            true,
						Validators: []validator.String{
							stringvalidator.OneOf(possibleConnectionTypeValues()...),
						},
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					"creation_method": schema.StringAttribute{
						MarkdownDescription: "The creation method used to create the connection.",
						Required:            true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					"parameters": schema.MapAttribute{
						MarkdownDescription: "A map of key/value pairs of connection parameters.",
						Optional:            true,
						CustomType:          supertypes.NewMapTypeOf[string](ctx),
					},
				},
			},
			"credential_details": schema.SingleNestedAttribute{
				MarkdownDescription: "The credential details of the connection.",
				Required:            true,
				CustomType:          supertypes.NewSingleNestedObjectTypeOf[rsCredentialDetailsModel](ctx),
				Attributes: map[string]schema.Attribute{
					"connection_encryption": schema.StringAttribute{
						MarkdownDescription: "The connection encryption type of the connection. Accepted values: " + utils.ConvertStringSlicesToString(fabcore.PossibleConnectionEncryptionValues(), true, true),
						Required:            true,
						Validators: []validator.String{
							stringvalidator.OneOf(utils.ConvertEnumsToStringSlices(fabcore.PossibleConnectionEncryptionValues(), false)...),
						},
					},
					"single_sign_on_type": schema.StringAttribute{
						MarkdownDescription: "The single sign-on type of the connection. Accepted values: " + utils.ConvertStringSlicesToString(fabcore.PossibleSingleSignOnTypeValues(), true, true),
						Required:            true,
						Validators: []validator.String{
							stringvalidator.OneOf(utils.ConvertEnumsToStringSlices(fabcore.PossibleSingleSignOnTypeValues(), false)...),
						},
					},
					"skip_test_connection": schema.BoolAttribute{
						MarkdownDescription: "Whether the connection should skip the test connection during creation and update. `True` - Skip the test connection, `False` - Do not skip the test connection.",
						Required:            true,
					},
					"credential_type": schema.StringAttribute{
						MarkdownDescription: "The credential type of the connection. Possible values: " + utils.ConvertStringSlicesToString(fabcore.PossibleCredentialTypeValues(), true, true),
						Required:            true,
						Validators: []validator.String{
							stringvalidator.OneOf(utils.ConvertEnumsToStringSlices(fabcore.PossibleCredentialTypeValues(), false)...),
						},
					},
					"basic_credentials": schema.SingleNestedAttribute{
						MarkdownDescription: "The basic credentials.",
						Optional:            true,
						Sensitive:           true,
						CustomType:          supertypes.NewSingleNestedObjectTypeOf[credentialsBasicModel](ctx),
						Attributes: map[string]schema.Attribute{
							"username": schema.StringAttribute{
								MarkdownDescription: "The username.",
								Required:            true,
								Sensitive:           true,
							},
							"password": schema.StringAttribute{
								MarkdownDescription: "The password.",
								Required:            true,
								Sensitive:           true,
							},
						},
						Validators: []validator.Object{
							objectvalidator.ConflictsWith(
								path.MatchRelative().AtParent().AtName("key_credentials"),
								path.MatchRelative().AtParent().AtName("service_principal_credentials"),
								path.MatchRelative().AtParent().AtName("shared_access_signature_credentials"),
								path.MatchRelative().AtParent().AtName("windows_credentials"),
							),
						},
					},
					"key_credentials": schema.SingleNestedAttribute{
						MarkdownDescription: "The key credentials.",
						Optional:            true,
						Sensitive:           true,
						CustomType:          supertypes.NewSingleNestedObjectTypeOf[credentialsKeyModel](ctx),
						Attributes: map[string]schema.Attribute{
							"key": schema.StringAttribute{
								MarkdownDescription: "The key.",
								Required:            true,
								Sensitive:           true,
							},
						},
						Validators: []validator.Object{
							objectvalidator.ConflictsWith(
								path.MatchRelative().AtParent().AtName("basic_credentials"),
								path.MatchRelative().AtParent().AtName("service_principal_credentials"),
								path.MatchRelative().AtParent().AtName("shared_access_signature_credentials"),
								path.MatchRelative().AtParent().AtName("windows_credentials"),
							),
						},
					},
					"service_principal_credentials": schema.SingleNestedAttribute{
						MarkdownDescription: "The service principal credentials.",
						Optional:            true,
						Sensitive:           true,
						CustomType:          supertypes.NewSingleNestedObjectTypeOf[credentialsServicePrincipalModel](ctx),
						Attributes: map[string]schema.Attribute{
							"tenant_id": schema.StringAttribute{
								MarkdownDescription: "The tenant ID.",
								Required:            true,
								Sensitive:           true,
							},
							"client_id": schema.StringAttribute{
								MarkdownDescription: "The client ID.",
								Required:            true,
								Sensitive:           true,
							},
							"client_secret": schema.StringAttribute{
								MarkdownDescription: "The client secret.",
								Required:            true,
								Sensitive:           true,
							},
						},
						Validators: []validator.Object{
							objectvalidator.ConflictsWith(
								path.MatchRelative().AtParent().AtName("basic_credentials"),
								path.MatchRelative().AtParent().AtName("key_credentials"),
								path.MatchRelative().AtParent().AtName("shared_access_signature_credentials"),
								path.MatchRelative().AtParent().AtName("windows_credentials"),
							),
						},
					},
					"shared_access_signature_credentials": schema.SingleNestedAttribute{
						MarkdownDescription: "The shared access signature credentials.",
						Optional:            true,
						Sensitive:           true,
						CustomType:          supertypes.NewSingleNestedObjectTypeOf[credentialsSharedAccessSignatureModel](ctx),
						Attributes: map[string]schema.Attribute{
							"token": schema.StringAttribute{
								MarkdownDescription: "The token.",
								Required:            true,
								Sensitive:           true,
							},
						},
						Validators: []validator.Object{
							objectvalidator.ConflictsWith(
								path.MatchRelative().AtParent().AtName("basic_credentials"),
								path.MatchRelative().AtParent().AtName("key_credentials"),
								path.MatchRelative().AtParent().AtName("service_principal_credentials"),
								path.MatchRelative().AtParent().AtName("windows_credentials"),
							),
						},
					},
					"windows_credentials": schema.SingleNestedAttribute{
						MarkdownDescription: "The Windows credentials.",
						Optional:            true,
						Sensitive:           true,
						CustomType:          supertypes.NewSingleNestedObjectTypeOf[credentialsWindowsModel](ctx),
						Attributes: map[string]schema.Attribute{
							"username": schema.StringAttribute{
								MarkdownDescription: "The username.",
								Required:            true,
								Sensitive:           true,
							},
							"password": schema.StringAttribute{
								MarkdownDescription: "The password.",
								Required:            true,
								Sensitive:           true,
							},
						},
						Validators: []validator.Object{
							objectvalidator.ConflictsWith(
								path.MatchRelative().AtParent().AtName("basic_credentials"),
								path.MatchRelative().AtParent().AtName("key_credentials"),
								path.MatchRelative().AtParent().AtName("service_principal_credentials"),
								path.MatchRelative().AtParent().AtName("shared_access_signature_credentials"),
							),
						},
					},
				},
			},
			"timeouts": timeouts.AttributesAll(ctx),
		},
	}
}

func (r *resourceConnection) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	pConfigData, ok := req.ProviderData.(*pconfig.ProviderData)
	if !ok {
		resp.Diagnostics.AddError(
			common.ErrorResourceConfigType,
			fmt.Sprintf(common.ErrorFabricClientType, req.ProviderData),
		)

		return
	}

	r.pConfigData = pConfigData
	r.client = fabcore.NewClientFactoryWithClient(*pConfigData.FabricClient).NewConnectionsClient()
}

func (r *resourceConnection) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	var plan resourceConnectionModel

	if resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...); resp.Diagnostics.HasError() {
		return
	}

	connectionDetails, diags := plan.getConnectionDetails(ctx)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	var supportedConnectionType fabcore.ConnectionCreationMetadata

	if resp.Diagnostics.Append(r.getConnectionTypeMetadata(ctx, *connectionDetails, &supportedConnectionType)...); resp.Diagnostics.HasError() {
		return
	}

	if resp.Diagnostics.Append(r.validateCreationMethod(*connectionDetails, supportedConnectionType.CreationMethods)...); resp.Diagnostics.HasError() {
		return
	}

	credentialDetails, diags := plan.getCredentialDetails(ctx)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	if resp.Diagnostics.Append(r.validateConnectionEncryption(*credentialDetails, supportedConnectionType.SupportedConnectionEncryptionTypes)...); resp.Diagnostics.HasError() {
		return
	}

	if resp.Diagnostics.Append(r.validateCredentialType(*credentialDetails, supportedConnectionType.SupportedCredentialTypes)...); resp.Diagnostics.HasError() {
		return
	}

	if resp.Diagnostics.Append(r.validateSkipTestConnection(*credentialDetails, *supportedConnectionType.SupportsSkipTestConnection)...); resp.Diagnostics.HasError() {
		return
	}

	if resp.Diagnostics.Append(r.validateCreationMethodParameters(ctx, *connectionDetails, supportedConnectionType.CreationMethods)...); resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceConnection) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "CREATE", map[string]any{
		"action": "start",
	})
	tflog.Trace(ctx, "CREATE", map[string]any{
		"config": req.Config,
		"plan":   req.Plan,
	})

	var plan resourceConnectionModel

	if resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...); resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := plan.Timeouts.Create(ctx, r.pConfigData.Timeout)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	connectionDetails, diags := plan.getConnectionDetails(ctx)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	var supportedConnectionType fabcore.ConnectionCreationMetadata

	if resp.Diagnostics.Append(r.getConnectionTypeMetadata(ctx, *connectionDetails, &supportedConnectionType)...); resp.Diagnostics.HasError() {
		return
	}

	var reqCreate requestCreateConnection

	if resp.Diagnostics.Append(reqCreate.set(ctx, plan, supportedConnectionType)...); resp.Diagnostics.HasError() {
		return
	}

	respCreate, err := r.client.CreateConnection(ctx, reqCreate.CreateConnectionRequestClassification, nil)
	if resp.Diagnostics.Append(utils.GetDiagsFromError(ctx, err, utils.OperationCreate, nil)...); resp.Diagnostics.HasError() {
		return
	}

	plan.set(respCreate.Connection)

	if resp.Diagnostics.Append(plan.setConnectionDetails(ctx, respCreate.Connection)...); resp.Diagnostics.HasError() {
		return
	}

	if resp.Diagnostics.Append(plan.setCredentialDetails(ctx, respCreate.Connection)...); resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)

	tflog.Debug(ctx, "CREATE", map[string]any{
		"action": "end",
	})

	if resp.Diagnostics.HasError() {
		return
	}

	// r.client.CreateConnection(ctx, createConnectionRequest fabcore.CreateConnectionRequestClassification, options *fabcore.ConnectionsClientCreateConnectionOptions)
}

func (r *resourceConnection) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
}

func (r *resourceConnection) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

func (r *resourceConnection) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}

func (r *resourceConnection) getConnectionTypeMetadata(ctx context.Context, model rsConnectionDetailsModel, supportedConnectionType *fabcore.ConnectionCreationMetadata) diag.Diagnostics {
	respList, err := r.client.ListSupportedConnectionTypes(ctx, &fabcore.ConnectionsClientListSupportedConnectionTypesOptions{
		ShowAllCreationMethods: azto.Ptr(true),
	})

	if diags := utils.GetDiagsFromError(ctx, err, utils.OperationList, nil); diags.HasError() {
		return diags
	}

	var vNames []string

	for _, v := range respList {
		if *v.Type == model.Type.ValueString() {
			*supportedConnectionType = v

			return nil
		}

		vNames = append(vNames, string(*v.Type))
	}

	var diags diag.Diagnostics

	diags.AddAttributeError(
		path.Root("connection_details").AtName("type"),
		"Unsupported connection type",
		fmt.Sprintf("The connection type '%s' is not supported. Supported values: %s", model.Type.ValueString(), utils.ConvertStringSlicesToString(vNames, true, true)),
	)

	return diags
}

func (r *resourceConnection) validateCreationMethod(model rsConnectionDetailsModel, elements []fabcore.ConnectionCreationMethod) diag.Diagnostics {
	var vNames []string

	for _, v := range elements {
		if *v.Name == model.CreationMethod.ValueString() {
			return nil
		}

		vNames = append(vNames, *v.Name)
	}

	var diags diag.Diagnostics

	diags.AddAttributeError(
		path.Root("connection_details").AtName("creation_method"),
		"Unsupported creation method",
		fmt.Sprintf("The creation method '%s' is not supported. Supported values: %s", model.CreationMethod.ValueString(), utils.ConvertStringSlicesToString(vNames, true, true)),
	)

	return diags
}

func (r *resourceConnection) validateConnectionEncryption(model rsCredentialDetailsModel, elements []fabcore.ConnectionEncryption) diag.Diagnostics {
	for _, v := range elements {
		if v == fabcore.ConnectionEncryption(model.ConnectionEncryption.ValueString()) {
			return nil
		}
	}

	var diags diag.Diagnostics

	diags.AddAttributeError(
		path.Root("credential_details").AtName("connection_encryption"),
		"Unsupported connection encryption",
		fmt.Sprintf("The connection encryption '%s' is not supported. Supported values: %s", model.ConnectionEncryption.ValueString(), utils.ConvertStringSlicesToString(utils.ConvertEnumsToStringSlices(elements, true), true, false)),
	)

	return diags
}

func (r *resourceConnection) validateCredentialType(model rsCredentialDetailsModel, elements []fabcore.CredentialType) diag.Diagnostics {
	for _, v := range elements {
		if v == fabcore.CredentialType(model.CredentialType.ValueString()) {
			return nil
		}
	}

	var diags diag.Diagnostics

	diags.AddAttributeError(
		path.Root("credential_details").AtName("credential_type"),
		"Unsupported credential type",
		fmt.Sprintf("The credential type '%s' is not supported. Supported values: %s", model.CredentialType.ValueString(), utils.ConvertStringSlicesToString(utils.ConvertEnumsToStringSlices(elements, true), true, false)),
	)

	return diags
}

func (r *resourceConnection) validateSkipTestConnection(model rsCredentialDetailsModel, supportsSkipTestConnection bool) diag.Diagnostics {
	if model.SkipTestConnection.ValueBool() != supportsSkipTestConnection {

		var diags diag.Diagnostics

		diags.AddAttributeError(
			path.Root("credential_details").AtName("skip_test_connection"),
			"Unsupported skip test connection",
			"The skip test connection value is not supported.",
		)

		return diags
	}

	return nil
}

func (r *resourceConnection) validateCreationMethodParameters(ctx context.Context, model rsConnectionDetailsModel, elements []fabcore.ConnectionCreationMethod) diag.Diagnostics {
	var diags diag.Diagnostics
	var vParameters []fabcore.ConnectionCreationParameter

	for _, v := range elements {
		if *v.Name == model.CreationMethod.ValueString() {
			vParameters = v.Parameters

			break
		}
	}

	// in reality, this should never happen because the creation method is already validated
	if len(vParameters) == 0 {
		diags.AddAttributeError(
			path.Root("connection_details").AtName("creation_method"),
			"Unsupported creation method",
			"The creation method is not supported.",
		)

		return diags
	}

	connectionDetailsParams, diags := model.getParameters(ctx)
	if diags.HasError() {
		return diags
	}

	// check if all keys of connectionDetailsParams are in vParameters
	var vNames []string

	for k := range connectionDetailsParams {
		var found bool

		for _, v := range vParameters {
			if *v.Name == k {
				found = true

				break
			} else {
				found = false
			}

			vNames = append(vNames, *v.Name)
		}

		if !found {
			diags.AddAttributeError(
				path.Root("connection_details").AtName("parameters"),
				"Unsupported connection parameter",
				fmt.Sprintf("The connection parameter '%s' is not supported. Supported parameters: %s", k, utils.ConvertStringSlicesToString(vNames, true, true)),
			)
		}
	}

	if diags.HasError() {
		return diags
	}

	// check if all required keys of vParameters are in connectionDetailsParams
	for _, v := range vParameters {
		if *v.Required {
			if _, ok := connectionDetailsParams[*v.Name]; !ok {
				diags.AddAttributeError(
					path.Root("connection_details").AtName("parameters"),
					"Missing connection parameter",
					fmt.Sprintf("The required connection parameter '%s' is missing.", *v.Name),
				)
			}

			// if connectionDetailsParams[*v.Name] == "" {
			// 	diags.AddAttributeError(
			// 		path.Root("connection_details").AtName("parameters"),
			// 		"Missing required connection parameter",
			// 		fmt.Sprintf("The required connection parameter '%s' is missing.", *v.Name),
			// 	)

			// 	return diags
			// }
		}
	}

	if diags.HasError() {
		return diags
	}

	return nil
}
