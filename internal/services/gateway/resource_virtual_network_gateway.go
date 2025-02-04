// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package gateway

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
	pconfig "github.com/microsoft/terraform-provider-fabric/internal/provider/config"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = (*resourceVirtualNetworkGateway)(nil)
	_ resource.ResourceWithConfigure   = (*resourceVirtualNetworkGateway)(nil)
	_ resource.ResourceWithImportState = (*resourceVirtualNetworkGateway)(nil)
)

type resourceVirtualNetworkGateway struct {
	pConfigData *pconfig.ProviderData
	client      *fabcore.GatewaysClient
}

func NewResourceVirtualNetworkGateway() resource.Resource {
	return &resourceVirtualNetworkGateway{}
}

func (r *resourceVirtualNetworkGateway) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + VirtualNetworkItemTFType
}

func (r *resourceVirtualNetworkGateway) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(200),
				},
			},
			"capacity_id": schema.StringAttribute{
				MarkdownDescription: "The " + ItemName + " capacity ID.",
				Required:            true,
				CustomType:          customtypes.UUIDType{},
			},
			"inactivity_minutes_before_sleep": schema.Int32Attribute{
				MarkdownDescription: "The " + ItemName + " inactivity minutes before sleep.",
				Required:            true,
				Validators: []validator.Int32{
					int32validator.OneOf(PossibleInactivityMinutesBeforeSleepValues...),
				},
			},
			"number_of_member_gateways": schema.Int32Attribute{
				MarkdownDescription: "The " + ItemName + " number of member gateways.",
				Required:            true,
				Validators: []validator.Int32{
					int32validator.Between(MinNumberOfMemberGatewaysValues, MaxNumberOfMemberGatewaysValues),
				},
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.RequiresReplace(),
				},
			},
			"virtual_network_azure_resource": schema.SingleNestedAttribute{
				MarkdownDescription: "The Azure resource of the virtual network.",
				Required:            true,
				CustomType:          supertypes.NewSingleNestedObjectTypeOf[virtualNetworkAzureResourceModel](ctx),
				Attributes: map[string]schema.Attribute{
					"subscription_id": schema.StringAttribute{
						MarkdownDescription: "The subscription ID.",
						Required:            true,
						CustomType:          customtypes.UUIDType{},
					},
					"resource_group_name": schema.StringAttribute{
						MarkdownDescription: "The name of the resource group.",
						Required:            true,
					},
					"virtual_network_name": schema.StringAttribute{
						MarkdownDescription: "The name of the virtual network.",
						Required:            true,
					},
					"subnet_name": schema.StringAttribute{
						MarkdownDescription: "The name of the subnet.",
						Required:            true,
					},
				},
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.RequiresReplace(),
				},
			},
			"timeouts": timeouts.AttributesAll(ctx),
		},
	}
}

func (r *resourceVirtualNetworkGateway) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *resourceVirtualNetworkGateway) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "CREATE", map[string]any{
		"action": "start",
	})
	tflog.Trace(ctx, "CREATE", map[string]any{
		"config": req.Config,
		"plan":   req.Plan,
	})

	var plan, state ResourceVirtualNetworkGatewayModel

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

	// get the virtual gateway from the classifiaction
	vng := respCreate.GatewayClassification.(*fabcore.VirtualNetworkGateway)
	state.set(ctx, *vng)

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)

	tflog.Debug(ctx, "CREATE", map[string]any{
		"action": "end",
	})

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceVirtualNetworkGateway) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "READ", map[string]any{
		"action": "start",
	})
	tflog.Trace(ctx, "READ", map[string]any{
		"state": req.State,
	})

	var state ResourceVirtualNetworkGatewayModel

	if resp.Diagnostics.Append(req.State.Get(ctx, &state)...); resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := state.Timeouts.Read(ctx, r.pConfigData.Timeout)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	err := r.get(ctx, &state)
	if diags := utils.GetDiagsFromError(ctx, err, utils.OperationRead, fabcore.ErrCommon.EntityNotFound); diags.HasError() {
		if utils.IsErrNotFound(state.ID.ValueString(), &diags, fabcore.ErrCommon.EntityNotFound) {
			resp.State.RemoveResource(ctx)
		}

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

func (r *resourceVirtualNetworkGateway) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "UPDATE", map[string]any{
		"action": "start",
	})
	tflog.Trace(ctx, "UPDATE", map[string]any{
		"config": req.Config,
		"plan":   req.Plan,
		"state":  req.State,
	})

	var plan ResourceVirtualNetworkGatewayModel

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

	// get the virtual gateway from the classifiaction
	vng := respUpdate.GatewayClassification.(*fabcore.VirtualNetworkGateway)
	plan.set(ctx, *vng)

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)

	tflog.Debug(ctx, "UPDATE", map[string]any{
		"action": "end",
	})

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceVirtualNetworkGateway) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "DELETE", map[string]any{
		"action": "start",
	})
	tflog.Trace(ctx, "DELETE", map[string]any{
		"state": req.State,
	})

	var state ResourceVirtualNetworkGatewayModel

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

func (r *resourceVirtualNetworkGateway) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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

func (r *resourceVirtualNetworkGateway) get(ctx context.Context, model *ResourceVirtualNetworkGatewayModel) error {
	tflog.Trace(ctx, "getting "+ItemName)

	respGet, err := r.client.GetGateway(ctx, model.ID.ValueString(), nil)
	if err != nil {
		return err
	}

	vng := respGet.GatewayClassification.(*fabcore.VirtualNetworkGateway)
	model.set(ctx, *vng)

	return nil
}
