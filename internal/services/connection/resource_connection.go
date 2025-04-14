// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package connection

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
	pconfig "github.com/microsoft/terraform-provider-fabric/internal/provider/config"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.ResourceWithConfigure   = (*resourceConnection)(nil)
	_ resource.ResourceWithImportState = (*resourceConnection)(nil)
)

type resourceConnection struct {
	pConfigData *pconfig.ProviderData
	client      *fabcore.ConnectionsClient
	TypeInfo    tftypeinfo.TFTypeInfo
}

func NewResourceConnection() resource.Resource {
	return &resourceConnection{
		TypeInfo: ItemTypeInfo,
	}
}

func (r *resourceConnection) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "fabric_connection"
}

func (r *resourceConnection) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	tflog.Info(ctx, "Building schema for connection_alt resource")
	// Define attribute types for connection_details and credential_details
	connectionDetailsAttributeTypes := map[string]attr.Type{
		"type":            types.StringType,
		"creation_method": types.StringType,
		"parameters": types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{
			"name":      types.StringType,
			"data_type": types.StringType,
			"value":     types.StringType,
		}}},
	}

	connectivityTypeValues := []string{"ShareableCloud", "OnPremisesGateway", "VirtualNetworkGateway"}
	privacyLevelValues := []string{"Organizational", "Private", "Public", "None"}

	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a Microsoft Fabric Connection",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The Connection ID.",
				CustomType:          customtypes.UUIDType{},
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "The Connection display name.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(200),
				},
			},
			"connectivity_type": schema.StringAttribute{
				MarkdownDescription: "Connectivity type. Possible values: " + utils.ConvertStringSlicesToString(connectivityTypeValues, true, true) + ".",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf(connectivityTypeValues...),
				},
			},
			"gateway_id": schema.StringAttribute{
				MarkdownDescription: "Gateway ID. Required for OnPremisesGateway and VirtualNetworkGateway connectivity types.",
				CustomType:          customtypes.UUIDType{},
				Optional:            true,
			},
			"privacy_level": schema.StringAttribute{
				MarkdownDescription: "Privacy level. Possible values: " + utils.ConvertStringSlicesToString(privacyLevelValues, true, true) + ".",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("Organizational"),
				Validators: []validator.String{
					stringvalidator.OneOf(privacyLevelValues...),
				},
			},
			"connection_details": schema.ObjectAttribute{
				MarkdownDescription: "Connection details. Can be specified as a nested block or as an object.",
				Required:            true,
				AttributeTypes:      connectionDetailsAttributeTypes,
			},
			"credential_details": schema.SingleNestedAttribute{
				MarkdownDescription: "Credential details. Can be specified as a nested block or as an object.",
				Required:            true,
				Attributes: map[string]schema.Attribute{
					"single_sign_on_type": schema.StringAttribute{
						MarkdownDescription: "Single sign-on type.",
						Required:            true,
					},
					"connection_encryption": schema.StringAttribute{
						MarkdownDescription: "Connection encryption type.",
						Required:            true,
					},
					"skip_test_connection": schema.BoolAttribute{
						MarkdownDescription: "Whether to skip test connection.",
						Required:            true,
					},
					"credentials": schema.SingleNestedAttribute{
						MarkdownDescription: "Credentials configuration.",
						Required:            true,
						Attributes: map[string]schema.Attribute{
							"credential_type": schema.StringAttribute{
								MarkdownDescription: "Credential type.",
								Required:            true,
							},
							// Make all other credential fields optional
							"username": schema.StringAttribute{
								MarkdownDescription: "Username for Basic or Windows authentication.",
								Optional:            true,
							},
							"password": schema.StringAttribute{
								MarkdownDescription: "Password for Basic or Windows authentication.",
								Optional:            true,
							},
							"key": schema.StringAttribute{
								MarkdownDescription: "Key for Key authentication.",
								Optional:            true,
							},
							"application_id": schema.StringAttribute{
								MarkdownDescription: "Application ID for Service Principal authentication.",
								Optional:            true,
							},
							"application_secret": schema.StringAttribute{
								MarkdownDescription: "Application Secret for Service Principal authentication.",
								Optional:            true,
							},
							"tenant_id": schema.StringAttribute{
								MarkdownDescription: "Tenant ID for Service Principal authentication.",
								Optional:            true,
							},
							"sas_token": schema.StringAttribute{
								MarkdownDescription: "SAS Token for Shared Access Signature authentication.",
								Optional:            true,
							},
							"domain": schema.StringAttribute{
								MarkdownDescription: "Domain for Windows authentication.",
								Optional:            true,
							},
						},
					},
				},
			},
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{}),
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

	if resp.Diagnostics.Append(fabricitem.IsPreviewMode(r.TypeInfo.Name, r.TypeInfo.IsPreview, r.pConfigData.Preview)...); resp.Diagnostics.HasError() {
		return
	}

	r.client = fabcore.NewClientFactoryWithClient(*pConfigData.FabricClient).NewConnectionsClient()
}

func (r *resourceConnection) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "CREATE", map[string]any{
		"action": "start",
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

	// Use the helper function from models_alt.go to build the request
	reqCreate, diags := buildConnectionRequest(
		ctx,
		plan.ConnectionDetails,
		plan.CredentialDetails,
		plan.DisplayName.ValueString(),
		plan.ConnectivityType.ValueString(),
		plan.PrivacyLevel.ValueString(),
		plan.GatewayID.ValueString(),
	)

	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	logRequestBody(ctx, reqCreate)

	respCreate, err := r.client.CreateConnection(ctx, reqCreate.CreateConnectionRequestClassification, nil)
	if resp.Diagnostics.Append(utils.GetDiagsFromError(ctx, err, utils.OperationCreate, nil)...); resp.Diagnostics.HasError() {
		return
	}

	var state resourceConnectionModel
	state.Timeouts = plan.Timeouts

	if resp.Diagnostics.Append(state.set(ctx, respCreate.Connection)...); resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)

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
	if utils.IsErrNotFound(state.ID.ValueString(), &diags, fabcore.ErrCommon.EntityNotFound) {
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

	var plan resourceConnectionModel

	if resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...); resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := plan.Timeouts.Update(ctx, r.pConfigData.Timeout)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var reqUpdate requestUpdateConnection

	if resp.Diagnostics.Append(reqUpdate.set(plan)...); resp.Diagnostics.HasError() {
		return
	}

	respUpdate, err := r.client.UpdateConnection(ctx, plan.ID.ValueString(), reqUpdate.UpdateConnectionRequestClassification, nil)
	if resp.Diagnostics.Append(utils.GetDiagsFromError(ctx, err, utils.OperationUpdate, nil)...); resp.Diagnostics.HasError() {
		return
	}

	if resp.Diagnostics.Append(plan.set(ctx, respUpdate.Connection)...); resp.Diagnostics.HasError() {
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

func (r *resourceConnection) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Debug(ctx, "IMPORT", map[string]any{
		"action": "start",
	})
	tflog.Trace(ctx, "IMPORT", map[string]any{
		"id": req.ID,
	})

	_, diags := customtypes.NewUUIDValueMust(req.ID)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)

	tflog.Debug(ctx, "IMPORT", map[string]any{
		"action": "end",
	})

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceConnection) get(ctx context.Context, model *resourceConnectionModel) diag.Diagnostics {
	respGet, err := r.client.GetConnection(ctx, model.ID.ValueString(), nil)
	if diags := utils.GetDiagsFromError(ctx, err, utils.OperationRead, fabcore.ErrCommon.EntityNotFound); diags.HasError() {
		return diags
	}

	return model.set(ctx, respGet.Connection)
}
