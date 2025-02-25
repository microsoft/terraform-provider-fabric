// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package connection

import (
	"context"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"

	azto "github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
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
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
	superobjectvalidator "github.com/orange-cloudavenue/terraform-plugin-framework-validators/objectvalidator"
	superstringvalidator "github.com/orange-cloudavenue/terraform-plugin-framework-validators/stringvalidator"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
	pconfig "github.com/microsoft/terraform-provider-fabric/internal/provider/config"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.ResourceWithConfigure  = (*resourceConnection)(nil)
	_ resource.ResourceWithModifyPlan = (*resourceConnection)(nil)
)

type resourceConnection struct {
	Name        string
	TFName      string
	IsPreview   bool
	pConfigData *pconfig.ProviderData
	client      *fabcore.ConnectionsClient
}

func NewResourceConnection() resource.Resource {
	return &resourceConnection{
		Name:      ItemName,
		TFName:    ItemTFName,
		IsPreview: ItemPreview,
	}
}

func (r *resourceConnection) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + r.TFName
}

func (r *resourceConnection) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	markdownDescription := "Manage a Fabric " + r.Name + ".\n\n" +
		"See [" + r.Name + "](" + ItemDocsURL + ") for more information.\n\n" +
		ItemDocsSPNSupport

	possibleConnectivityTypeValues := utils.RemoveSlicesByValues(fabcore.PossibleConnectivityTypeValues(), []fabcore.ConnectivityType{fabcore.ConnectivityTypeOnPremisesGateway, fabcore.ConnectivityTypeOnPremisesGatewayPersonal, fabcore.ConnectivityTypePersonalCloud})

	resp.Schema = schema.Schema{
		MarkdownDescription: fabricitem.GetResourcePreviewNote(markdownDescription, r.IsPreview),
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
				MarkdownDescription: "The connectivity type of the connection. Accepted values: " + utils.ConvertStringSlicesToString(possibleConnectivityTypeValues, true, true),
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf(utils.ConvertEnumsToStringSlices(possibleConnectivityTypeValues, false)...),
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
							// types.StringValue(string(fabcore.ConnectivityTypeOnPremisesGatewayPersonal)),
						}),
					superstringvalidator.NullIfAttributeIsOneOf(path.MatchRoot("connectivity_type"),
						[]attr.Value{
							types.StringValue(string(fabcore.ConnectivityTypeAutomatic)),
							types.StringValue(string(fabcore.ConnectivityTypeNone)),
							// types.StringValue(string(fabcore.ConnectivityTypePersonalCloud)),
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
						MarkdownDescription: "The type of the connection.",
						Required:            true,
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
						Optional:            true,
						Validators: []validator.String{
							stringvalidator.OneOf(utils.ConvertEnumsToStringSlices(fabcore.PossibleCredentialTypeValues(), false)...),
							superstringvalidator.NullIfAttributeIsOneOf(
								path.MatchRoot("connectivity_type"),
								[]attr.Value{types.StringValue(string(fabcore.ConnectivityTypeOnPremisesGateway))},
							),
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
								// path.MatchRelative().AtParent().AtName("encrypted_credentials"),
							),
							superobjectvalidator.RequireIfAttributeIsOneOf(
								path.MatchRelative().AtParent().AtName("credential_type"),
								[]attr.Value{
									types.StringValue(string(fabcore.CredentialTypeBasic)),
								},
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
								// path.MatchRelative().AtParent().AtName("encrypted_credentials"),
							),
							superobjectvalidator.RequireIfAttributeIsOneOf(
								path.MatchRelative().AtParent().AtName("credential_type"),
								[]attr.Value{
									types.StringValue(string(fabcore.CredentialTypeKey)),
								},
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
								// path.MatchRelative().AtParent().AtName("encrypted_credentials"),
							),
							superobjectvalidator.RequireIfAttributeIsOneOf(
								path.MatchRelative().AtParent().AtName("credential_type"),
								[]attr.Value{
									types.StringValue(string(fabcore.CredentialTypeServicePrincipal)),
								},
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
								// path.MatchRelative().AtParent().AtName("encrypted_credentials"),
							),
							superobjectvalidator.RequireIfAttributeIsOneOf(
								path.MatchRelative().AtParent().AtName("credential_type"),
								[]attr.Value{
									types.StringValue(string(fabcore.CredentialTypeSharedAccessSignature)),
								},
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
								// path.MatchRelative().AtParent().AtName("encrypted_credentials"),
							),
							superobjectvalidator.RequireIfAttributeIsOneOf(
								path.MatchRelative().AtParent().AtName("credential_type"),
								[]attr.Value{
									types.StringValue(string(fabcore.CredentialTypeWindows)),
								},
							),
						},
					},
					// "encrypted_credentials": schema.SingleNestedAttribute{
					// 	MarkdownDescription: "The encrypted serialized .json of the list of name value pairs. Name is a credential name and value is a credential value. Encryption is performed using the Rivest-Shamir-Adleman (RSA) encryption algorithm with the on-premises gateway member's public key.",
					// 	Optional:            true,
					// 	Sensitive:           true,
					// 	CustomType:          supertypes.NewSingleNestedObjectTypeOf[credentialsEncryptedModel](ctx),
					// 	Attributes: map[string]schema.Attribute{
					// 		"value": schema.StringAttribute{
					// 			MarkdownDescription: "The value.",
					// 			Required:            true,
					// 			Sensitive:           true,
					// 		},
					// 	},
					// 	Validators: []validator.Object{
					// 		objectvalidator.ConflictsWith(
					// 			path.MatchRelative().AtParent().AtName("basic_credentials"),
					// 			path.MatchRelative().AtParent().AtName("key_credentials"),
					// 			path.MatchRelative().AtParent().AtName("service_principal_credentials"),
					// 			path.MatchRelative().AtParent().AtName("windows_credentials"),
					// 			path.MatchRelative().AtParent().AtName("shared_access_signature_credentials"),
					// 		),
					// 		superobjectvalidator.RequireIfAttributeIsOneOf(
					// 			path.MatchRoot("connectivity_type"),
					// 			[]attr.Value{
					// 				types.StringValue(string(fabcore.ConnectivityTypeOnPremisesGateway)),
					// 				// types.StringValue(string(fabcore.ConnectivityTypeOnPremisesGatewayPersonal)),
					// 			},
					// 		),
					// 	},
					// },
				},
			},
			"timeouts": timeouts.AttributesAll(ctx),
		},
	}
}

func (r *resourceConnection) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	if resp.Diagnostics.Append(fabricitem.IsPreviewMode(r.Name, r.IsPreview, r.pConfigData.Preview)...); resp.Diagnostics.HasError() {
		return
	}
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

	// connectionDetails, diags := plan.getConnectionDetails(ctx)
	// if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
	// 	return
	// }

	// var supportedConnectionType fabcore.ConnectionCreationMetadata

	// if resp.Diagnostics.Append(r.getConnectionTypeMetadata(ctx, *connectionDetails, &supportedConnectionType)...); resp.Diagnostics.HasError() {
	// 	return
	// }

	var reqCreate requestCreateConnection

	if resp.Diagnostics.Append(reqCreate.set(ctx, plan)...); resp.Diagnostics.HasError() {
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
}

func (r *resourceConnection) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "READ", map[string]any{
		"action": "start",
	})
	tflog.Trace(ctx, "READ", map[string]any{
		"state": req.State,
	})

	var state resourceConnectionModel

	if resp.Diagnostics.Append(req.State.Get(ctx, &state)...); resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := state.Timeouts.Read(ctx, r.pConfigData.Timeout)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	diags = r.get(ctx, &state)
	if utils.IsErrNotFound(state.ID.ValueString(), &diags, fabcore.ErrWorkspace.WorkspaceNotFound) {
		resp.State.RemoveResource(ctx)

		resp.Diagnostics.Append(diags...)

		return
	}

	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)

	tflog.Debug(ctx, "READ", map[string]any{
		"action": "end",
	})

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceConnection) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "UPDATE", map[string]any{
		"action": "start",
	})
	tflog.Trace(ctx, "UPDATE", map[string]any{
		"config": req.Config,
		"plan":   req.Plan,
		"state":  req.State,
	})

	var plan resourceConnectionModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := plan.Timeouts.Update(ctx, r.pConfigData.Timeout)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var reqUpdate requestUpdateConnection

	if resp.Diagnostics.Append(reqUpdate.set(ctx, plan)...); resp.Diagnostics.HasError() {
		return
	}

	respUpdate, err := r.client.UpdateConnection(ctx, plan.ID.ValueString(), reqUpdate, nil)
	if resp.Diagnostics.Append(utils.GetDiagsFromError(ctx, err, utils.OperationUpdate, nil)...); resp.Diagnostics.HasError() {
		return
	}

	plan.set(respUpdate.Connection)

	if resp.Diagnostics.Append(plan.setConnectionDetails(ctx, respUpdate.Connection)...); resp.Diagnostics.HasError() {
		return
	}

	if resp.Diagnostics.Append(plan.setCredentialDetails(ctx, respUpdate.Connection)...); resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)

	tflog.Debug(ctx, "UPDATE", map[string]any{
		"action": "end",
	})

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceConnection) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "DELETE", map[string]any{
		"action": "start",
	})
	tflog.Trace(ctx, "DELETE", map[string]any{
		"state": req.State,
	})

	var state resourceConnectionModel

	if resp.Diagnostics.Append(req.State.Get(ctx, &state)...); resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := state.Timeouts.Delete(ctx, r.pConfigData.Timeout)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	_, err := r.client.DeleteConnection(ctx, state.ID.ValueString(), nil)
	if resp.Diagnostics.Append(utils.GetDiagsFromError(ctx, err, utils.OperationDelete, nil)...); resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "DELETE", map[string]any{
		"action": "end",
	})
}

func (r *resourceConnection) get(ctx context.Context, model *resourceConnectionModel) diag.Diagnostics {
	tflog.Trace(ctx, "GET", map[string]any{
		"id": model.ID.ValueString(),
	})

	respGet, err := r.client.GetConnection(ctx, model.ID.ValueString(), nil)
	if diags := utils.GetDiagsFromError(ctx, err, utils.OperationRead, fabcore.ErrWorkspace.WorkspaceNotFound); diags.HasError() {
		return diags
	}

	model.set(respGet.Connection)

	if diags := model.setConnectionDetails(ctx, respGet.Connection); diags.HasError() {
		return diags
	}

	if diags := model.setCredentialDetails(ctx, respGet.Connection); diags.HasError() {
		return diags
	}

	return nil
}

func (r *resourceConnection) getConnectionTypeMetadata(ctx context.Context, model rsConnectionDetailsModel, supportedConnectionType *fabcore.ConnectionCreationMetadata) diag.Diagnostics {
	respList, err := r.client.ListSupportedConnectionTypes(ctx, &fabcore.ConnectionsClientListSupportedConnectionTypesOptions{
		ShowAllCreationMethods: azto.Ptr(true),
	})

	if diags := utils.GetDiagsFromError(ctx, err, utils.OperationList, nil); diags.HasError() {
		return diags
	}

	vNames := make([]string, 0, len(respList))

	for _, v := range respList {
		if *v.Type == model.Type.ValueString() {
			*supportedConnectionType = v

			return nil
		}

		vNames = append(vNames, *v.Type)
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
	vNames := make([]string, 0, len(elements))

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

func (r *resourceConnection) validateSkipTestConnection(model rsCredentialDetailsModel, supportsSkipTestConnection bool) diag.Diagnostics { //revive:disable-line:flag-parameter
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
			}

			vNames = append(vNames, *v.Name)
		}

		if !found {
			diags.AddAttributeError(
				path.Root("connection_details").AtName("parameters"),
				"Unsupported connection parameter key",
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
					"Missing connection parameter key",
					fmt.Sprintf("The required connection parameter '%s' is missing.", *v.Name),
				)
			}

			if connectionDetailsParams[*v.Name] == "" {
				diags.AddAttributeError(
					path.Root("connection_details").AtName("parameters"),
					"Missing connection parameter value",
					fmt.Sprintf("The required connection parameter '%s' value is missing.", *v.Name),
				)
			}
		}
	}

	if diags.HasError() {
		return diags
	}

	for k, v := range connectionDetailsParams {
		var dataType fabcore.DataType
		var allowedValues []string

		for _, v := range vParameters {
			if *v.Name == k {
				dataType = *v.DataType
				allowedValues = v.AllowedValues

				break
			}
		}

		switch dataType {
		case fabcore.DataTypeBoolean:
			// Use boolean as the parameter input value. False - the value is false, True - the value is true.
			if !strings.EqualFold(v, "true") && !strings.EqualFold(v, "false") {
				diags.AddAttributeError(
					path.Root("connection_details").AtName("parameters"),
					"Invalid connection parameter value",
					fmt.Sprintf("The connection parameter '%s' value is invalid. Supported values: `True`, `False`", k),
				)
			}
		case fabcore.DataTypeDate:
			// Use date as the parameter input value, using YYYY-MM-DD format.
			if _, err := time.Parse(time.DateOnly, v); err != nil {
				diags.AddAttributeError(
					path.Root("connection_details").AtName("parameters"),
					"Invalid connection parameter value",
					fmt.Sprintf("The connection parameter '%s' value is invalid. Supported format: `YYYY-MM-DD`", k),
				)
			}
		case fabcore.DataTypeDateTime:
			// Use date time as the parameter input value, using YYYY-MM-DDTHH:mm:ss.FFFZ (ISO 8601) format.
			if _, err := time.Parse("2006-01-02T15:04:05.000Z07:00", v); err != nil {
				diags.AddAttributeError(
					path.Root("connection_details").AtName("parameters"),
					"Invalid connection parameter value",
					fmt.Sprintf("The connection parameter '%s' value is invalid. Supported format: `YYYY-MM-DDTHH:mm:ss.FFFZ`", k),
				)
			}
		case fabcore.DataTypeDateTimeZone:
			// Use date time zone as the parameter input value, using YYYY-MM-DDTHH:mm:ss.FFF±hh:mm format.
			if _, err := time.Parse("2006-01-02T03:04:05.000-07:00", v); err != nil {
				diags.AddAttributeError(
					path.Root("connection_details").AtName("parameters"),
					"Invalid connection parameter value",
					fmt.Sprintf("The connection parameter '%s' value is invalid. Supported format: `YYYY-MM-DDTHH:mm:ss.FFF±hh:mm`", k),
				)
			}
		case fabcore.DataTypeDuration:
			// Use duration as the parameter input value, using [-]P(n)DT(n)H(n)M(n)S format. For example: P3DT4H30M10S (for 3 days, 4 hours, 30 minutes, and 10 seconds).
			if _, err := time.ParseDuration(v); err != nil {
				diags.AddAttributeError(
					path.Root("connection_details").AtName("parameters"),
					"Invalid connection parameter value",
					fmt.Sprintf("The connection parameter '%s' value is invalid. Supported format: `[-]P(n)DT(n)H(n)M(n)S`", k),
				)
			}
		case fabcore.DataTypeNumber:
			// Use number as the parameter input value (integer or floating point).
			if _, err := strconv.ParseFloat(v, 64); err != nil {
				diags.AddAttributeError(
					path.Root("connection_details").AtName("parameters"),
					"Invalid connection parameter value",
					fmt.Sprintf("The connection parameter '%s' value is invalid. It must be integer or floating point.", k),
				)
			}
		case fabcore.DataTypeText:
			// Use text as the parameter input value.
			if v == "" {
				diags.AddAttributeError(
					path.Root("connection_details").AtName("parameters"),
					"Invalid connection parameter value",
					fmt.Sprintf("The connection parameter '%s' value is invalid. It must not be empty.", k),
				)
			}
		case fabcore.DataTypeTime:
			// Use time as the parameter input value, using HH:mm:ss.FFFZ format.
			if _, err := time.Parse("15:04:05.000Z07:00", v); err != nil {
				diags.AddAttributeError(
					path.Root("connection_details").AtName("parameters"),
					"Invalid connection parameter value",
					fmt.Sprintf("The connection parameter '%s' value is invalid. Supported format: `HH:mm:ss.FFFZ`", k),
				)
			}
		}

		if len(allowedValues) > 0 {
			if !slices.Contains(allowedValues, v) {
				diags.AddAttributeError(
					path.Root("connection_details").AtName("parameters"),
					"Invalid connection parameter value",
					fmt.Sprintf("The connection parameter '%s' value is invalid. Supported values: %s", k, utils.ConvertStringSlicesToString(allowedValues, true, true)),
				)
			}
		}
	}

	return diags
}
