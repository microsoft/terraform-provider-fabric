// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package gateway

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
	pconfig "github.com/microsoft/terraform-provider-fabric/internal/provider/config"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.ResourceWithConfigure   = (*resourceGatewayRoleAssignment)(nil)
	_ resource.ResourceWithImportState = (*resourceGatewayRoleAssignment)(nil)
)

type resourceGatewayRoleAssignment struct {
	pConfigData *pconfig.ProviderData
	client      *fabcore.GatewaysClient
	Name        string
	IsPreview   bool
}

func NewResourceGatewayRoleAssignment() resource.Resource {
	return &resourceGatewayRoleAssignment{
		Name:      GatewayRoleAssignmentName,
		IsPreview: ItemPreview,
	}
}

func (r *resourceGatewayRoleAssignment) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + GatewayRoleAssignmentTFName
}

func (r *resourceGatewayRoleAssignment) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: fabricitem.GetResourcePreviewNote("Manage a "+GatewayRoleAssignmentName+".\n\n"+
			ItemDocsSPNSupport, r.IsPreview),
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The " + GatewayRoleAssignmentName + " ID.",
				Computed:            true,
				CustomType:          customtypes.UUIDType{},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"principal_id": schema.StringAttribute{
				MarkdownDescription: "The Principal ID.",
				Required:            true,
				CustomType:          customtypes.UUIDType{},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"principal_type": schema.StringAttribute{
				MarkdownDescription: "The type of the principal. Accepted values: " + utils.ConvertStringSlicesToString(fabcore.PossiblePrincipalTypeValues(), true, true) + ".",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf(utils.ConvertEnumsToStringSlices(fabcore.PossiblePrincipalTypeValues(), false)...),
				},
			},
			"role": schema.StringAttribute{
				MarkdownDescription: "The Gateway Role of the principal. Accepted values: " + utils.ConvertStringSlicesToString(fabcore.PossibleGatewayRoleValues(), true, true) + ".",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf(utils.ConvertEnumsToStringSlices(fabcore.PossibleGatewayRoleValues(), false)...),
				},
			},
			"gateway_id": schema.StringAttribute{
				MarkdownDescription: "The " + ItemName + " ID.",
				Required:            true,
				CustomType:          customtypes.UUIDType{},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"timeouts": timeouts.AttributesAll(ctx),
		},
	}
}

func (r *resourceGatewayRoleAssignment) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	if resp.Diagnostics.Append(fabricitem.IsPreviewMode(r.Name, r.IsPreview, r.pConfigData.Preview)...); resp.Diagnostics.HasError() {
		return
	}

	r.client = fabcore.NewClientFactoryWithClient(*pConfigData.FabricClient).NewGatewaysClient()
}

func (r *resourceGatewayRoleAssignment) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "CREATE", map[string]any{
		"action": "start",
	})
	tflog.Trace(ctx, "CREATE", map[string]any{
		"config": req.Config,
		"plan":   req.Plan,
	})

	var plan resourceGatewayRoleAssignmentModel

	if resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...); resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := plan.Timeouts.Create(ctx, r.pConfigData.Timeout)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var reqCreate requestCreateGatewayRoleAssignment

	reqCreate.set(plan)

	respCreate, err := r.client.AddGatewayRoleAssignment(ctx, plan.GatewayID.ValueString(), reqCreate.AddGatewayRoleAssignmentRequest, nil)
	if resp.Diagnostics.Append(utils.GetDiagsFromError(ctx, err, utils.OperationCreate, nil)...); resp.Diagnostics.HasError() {
		return
	}

	plan.set(respCreate.GatewayRoleAssignment)

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)

	tflog.Debug(ctx, "CREATE", map[string]any{
		"action": "end",
	})

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceGatewayRoleAssignment) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "READ", map[string]any{
		"action": "start",
	})
	tflog.Trace(ctx, "READ", map[string]any{
		"state": req.State,
	})

	var state resourceGatewayRoleAssignmentModel

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
	if diags := utils.GetDiagsFromError(ctx, err, utils.OperationRead, fabcore.ErrCommon.EntityNotFound); resp.Diagnostics.HasError() {
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

func (r *resourceGatewayRoleAssignment) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "UPDATE", map[string]any{
		"action": "start",
	})
	tflog.Trace(ctx, "UPDATE", map[string]any{
		"config": req.Config,
		"plan":   req.Plan,
		"state":  req.State,
	})

	var plan resourceGatewayRoleAssignmentModel

	if resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...); resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := plan.Timeouts.Update(ctx, r.pConfigData.Timeout)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var reqUpdate requestUpdateGatewayRoleAssignment

	reqUpdate.set(plan)

	respUpdate, err := r.client.UpdateGatewayRoleAssignment(ctx, plan.GatewayID.ValueString(), plan.ID.ValueString(), reqUpdate.UpdateGatewayRoleAssignmentRequest, nil)
	if resp.Diagnostics.Append(utils.GetDiagsFromError(ctx, err, utils.OperationUpdate, nil)...); resp.Diagnostics.HasError() {
		return
	}

	plan.set(respUpdate.GatewayRoleAssignment)

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)

	tflog.Debug(ctx, "UPDATE", map[string]any{
		"action": "end",
	})

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceGatewayRoleAssignment) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "DELETE", map[string]any{
		"action": "start",
	})
	tflog.Trace(ctx, "DELETE", map[string]any{
		"state": req.State,
	})

	var state resourceGatewayRoleAssignmentModel

	if resp.Diagnostics.Append(req.State.Get(ctx, &state)...); resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := state.Timeouts.Delete(ctx, r.pConfigData.Timeout)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	_, err := r.client.DeleteGatewayRoleAssignment(ctx, state.GatewayID.ValueString(), state.ID.ValueString(), nil)
	if resp.Diagnostics.Append(utils.GetDiagsFromError(ctx, err, utils.OperationDelete, nil)...); resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "DELETE", map[string]any{
		"action": "end",
	})
}

func (r *resourceGatewayRoleAssignment) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Debug(ctx, "IMPORT", map[string]any{
		"action": "start",
	})
	tflog.Trace(ctx, "IMPORT", map[string]any{
		"id": req.ID,
	})

	gatewayID, gatewayRoleAssignmentID, found := strings.Cut(req.ID, "/")
	if !found {
		resp.Diagnostics.AddError(
			common.ErrorImportIdentifierHeader,
			fmt.Sprintf(common.ErrorImportIdentifierDetails, "GatewayID/GatewayRoleAssignmentID"),
		)

		return
	}

	uuidGatewayID, diags := customtypes.NewUUIDValueMust(gatewayID)
	resp.Diagnostics.Append(diags...)

	uuidGatewayRoleAssignmentID, diags := customtypes.NewUUIDValueMust(gatewayRoleAssignmentID)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	var timeout timeouts.Value
	if resp.Diagnostics.Append(resp.State.GetAttribute(ctx, path.Root("timeouts"), &timeout)...); resp.Diagnostics.HasError() {
		return
	}

	state := resourceGatewayRoleAssignmentModel{
		ID:        uuidGatewayRoleAssignmentID,
		GatewayID: uuidGatewayID,
		Timeouts:  timeout,
	}

	err := r.get(ctx, &state)
	if resp.Diagnostics.Append(utils.GetDiagsFromError(ctx, err, utils.OperationImport, nil)...); resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)

	tflog.Debug(ctx, "IMPORT", map[string]any{
		"action": "end",
	})

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceGatewayRoleAssignment) get(ctx context.Context, model *resourceGatewayRoleAssignmentModel) error {
	tflog.Trace(ctx, "getting Gateway Role Assignment")

	respGetInfo, err := r.client.GetGatewayRoleAssignment(ctx, model.GatewayID.ValueString(), model.ID.ValueString(), nil)
	if err != nil {
		return err
	}

	model.set(respGetInfo.GatewayRoleAssignment)

	return nil
}
