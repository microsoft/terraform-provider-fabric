// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package gateway

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
	superint32validator "github.com/orange-cloudavenue/terraform-plugin-framework-validators/int32validator"
	supersetvalidator "github.com/orange-cloudavenue/terraform-plugin-framework-validators/setvalidator"
	superstringvalidator "github.com/orange-cloudavenue/terraform-plugin-framework-validators/stringvalidator"

	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
	pconfig "github.com/microsoft/terraform-provider-fabric/internal/provider/config"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.ResourceWithConfigure = (*resourceGateway)(nil)
	// _ resource.ResourceWithImportState = (*resourceGateway)(nil)
)

type resourceGateway struct {
	pConfigData *pconfig.ProviderData
	client      *fabcore.GatewaysClient
}

func NewResourceGateway() resource.Resource {
	return &resourceGateway{}
}

func (r *resourceGateway) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + ItemTFName
}

func (r *resourceGateway) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "This resource manages a Fabric " + ItemName + ".\n\n" +
			"See [" + ItemName + "s](" + ItemDocsURL + ") for more information.\n\n" +
			ItemDocsSPNSupport,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The " + ItemName + " ID.",
				Computed:            true,
				CustomType:          customtypes.UUIDType{},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "The " + ItemName + " display name.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(200),
					superstringvalidator.RequireIfAttributeIsOneOf(path.MatchRoot("type"),
						[]attr.Value{
							types.StringValue(string(fabcore.GatewayTypeVirtualNetwork)),
						}),
				},
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "The " + ItemName + " type. Accepted values: " + utils.ConvertStringSlicesToString(fabcore.PossibleGatewayTypeValues(), true, true),
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf(utils.ConvertEnumsToStringSlices(fabcore.PossibleGatewayTypeValues(), false)...),
				},
			},
			"capacity_id": schema.StringAttribute{
				MarkdownDescription: "The " + ItemName + " capacity ID.",
				Optional:            true,
				CustomType:          customtypes.UUIDType{},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					superstringvalidator.RequireIfAttributeIsOneOf(path.MatchRoot("type"),
						[]attr.Value{
							types.StringValue(string(fabcore.GatewayTypeVirtualNetwork)),
						}),
				},
			},
			"inactivity_minutes_before_sleep": schema.Int32Attribute{
				MarkdownDescription: "The " + ItemName + " inactivity minutes before sleep.",
				Required:            true,
				Validators: []validator.Int32{
					int32validator.OneOf(PossibleInactivityMinutesBeforeSleepValues...),
					superint32validator.RequireIfAttributeIsOneOf(path.MatchRoot("type"),
						[]attr.Value{
							types.StringValue(string(fabcore.GatewayTypeVirtualNetwork)),
						}),
				},
			},
			"number_of_member_gateways": schema.Int32Attribute{
				MarkdownDescription: "The " + ItemName + " number of member gateways.",
				Required:            true,
				Validators: []validator.Int32{
					int32validator.Between(MinNumberOfMemberGatewaysValues, MaxNumberOfMemberGatewaysValues),
					superint32validator.RequireIfAttributeIsOneOf(path.MatchRoot("type"),
						[]attr.Value{
							types.StringValue(string(fabcore.GatewayTypeVirtualNetwork)),
						}),
				},
			},
			"virtual_network_azure_resource": schema.SetNestedAttribute{
				MarkdownDescription: "The " + ItemName + " virtual network Azure resource.",
				Optional:            true,
				CustomType:          supertypes.NewSetNestedObjectTypeOf[virtualNetworkAzureResourceModel](ctx),
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.RequiresReplace(),
				},
				Validators: []validator.Set{
					supersetvalidator.RequireIfAttributeIsOneOf(path.MatchRoot("type"),
						[]attr.Value{
							types.StringValue(string(fabcore.GatewayTypeVirtualNetwork)),
						}),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"virtual_network_name": schema.StringAttribute{
							MarkdownDescription: "The virtual network name.",
							Required:            true,
						},
						"subnet_name": schema.StringAttribute{
							MarkdownDescription: "The subnet name.",
							Required:            true,
						},
						"resource_group_name": schema.StringAttribute{
							MarkdownDescription: "The resource group name.",
							Required:            true,
						},
						"subscription_id": schema.StringAttribute{
							MarkdownDescription: "The subscription ID.",
							Required:            true,
							CustomType:          customtypes.UUIDType{},
						},
					},
				},
			},
			"timeouts": timeouts.AttributesAll(ctx),
		},
	}
}

func (r *resourceGateway) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
	r.client = fabcore.NewClientFactoryWithClient(*pConfigData.FabricClient).NewGatewaysClient()
}

func (r *resourceGateway) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "CREATE", map[string]any{
		"action": "start",
	})
	tflog.Trace(ctx, "CREATE", map[string]any{
		"config": req.Config,
		"plan":   req.Plan,
	})

	var plan, state resourceGatewayModel

	if resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...); resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := plan.Timeouts.Create(ctx, r.pConfigData.Timeout)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	state.Timeouts = plan.Timeouts

	var reqCreate requestCreateGateway

	reqCreate.set(ctx, plan)

	respCreate, err := r.client.CreateGateway(ctx, reqCreate.CreateGatewayRequestClassification, nil)
	if resp.Diagnostics.Append(utils.GetDiagsFromError(ctx, err, utils.OperationCreate, nil)...); resp.Diagnostics.HasError() {
		return
	}

	if resp.Diagnostics.Append(state.set(ctx, respCreate.GatewayClassification)...); resp.Diagnostics.HasError() {
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

func (r *resourceGateway) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "READ", map[string]any{
		"action": "start",
	})
	tflog.Trace(ctx, "READ", map[string]any{
		"state": req.State,
	})

	var state resourceGatewayModel

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

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)

	tflog.Debug(ctx, "READ", map[string]any{
		"action": "end",
	})

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceGateway) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "UPDATE", map[string]any{
		"action": "start",
	})
	tflog.Trace(ctx, "UPDATE", map[string]any{
		"config": req.Config,
		"plan":   req.Plan,
		"state":  req.State,
	})

	var plan resourceGatewayModel

	if resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...); resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := plan.Timeouts.Update(ctx, r.pConfigData.Timeout)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var reqUpdate requestUpdateGateway

	reqUpdate.set(plan)

	respUpdate, err := r.client.UpdateGateway(ctx, plan.ID.ValueString(), reqUpdate.UpdateGatewayRequestClassification, nil)
	if resp.Diagnostics.Append(utils.GetDiagsFromError(ctx, err, utils.OperationUpdate, nil)...); resp.Diagnostics.HasError() {
		return
	}

	if resp.Diagnostics.Append(plan.set(ctx, respUpdate.GatewayClassification)...); resp.Diagnostics.HasError() {
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

func (r *resourceGateway) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "DELETE", map[string]any{
		"action": "start",
	})
	tflog.Trace(ctx, "DELETE", map[string]any{
		"state": req.State,
	})

	var state resourceGatewayModel

	if resp.Diagnostics.Append(req.State.Get(ctx, &state)...); resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := state.Timeouts.Delete(ctx, r.pConfigData.Timeout)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	_, err := r.client.DeleteGateway(ctx, state.ID.ValueString(), nil)
	if resp.Diagnostics.Append(utils.GetDiagsFromError(ctx, err, utils.OperationDelete, nil)...); resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "DELETE", map[string]any{
		"action": "end",
	})
}

// func (r *resourceGateway) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
// 	tflog.Debug(ctx, "IMPORT", map[string]any{
// 		"action": "start",
// 	})
// 	tflog.Trace(ctx, "IMPORT", map[string]any{
// 		"id": req.ID,
// 	})

// 	_, diags := customtypes.NewUUIDValueMust(req.ID)
// 	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
// 		return
// 	}

// 	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)

// 	tflog.Debug(ctx, "IMPORT", map[string]any{
// 		"action": "end",
// 	})

// 	if resp.Diagnostics.HasError() {
// 		return
// 	}
// }

func (r *resourceGateway) get(ctx context.Context, model *resourceGatewayModel) diag.Diagnostics {
	tflog.Trace(ctx, "getting "+ItemName)

	respGet, err := r.client.GetGateway(ctx, model.ID.ValueString(), nil)
	if diags := utils.GetDiagsFromError(ctx, err, utils.OperationRead, nil); diags.HasError() {
		return diags
	}

	return model.set(ctx, respGet.GatewayClassification)
}
