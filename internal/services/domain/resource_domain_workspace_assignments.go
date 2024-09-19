// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package domain

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	fabadmin "github.com/microsoft/fabric-sdk-go/fabric/admin"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
	pconfig "github.com/microsoft/terraform-provider-fabric/internal/provider/config"
)

// Ensure the implementation satisfies the expected interfaces.
var _ resource.ResourceWithConfigure = (*resourceDomainWorkspaceAssignments)(nil)

type resourceDomainWorkspaceAssignments struct {
	pConfigData *pconfig.ProviderData
	client      *fabadmin.DomainsClient
}

func NewResourceDomainWorkspaceAssignments() resource.Resource {
	return &resourceDomainWorkspaceAssignments{}
}

func (r *resourceDomainWorkspaceAssignments) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + DomainWorkspaceAssignmentsTFName
}

func (r *resourceDomainWorkspaceAssignments) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manage a Fabric " + DomainWorkspaceAssignmentsName + ".\n\n" +
			"See [" + ItemName + "](" + ItemDocsURL + ") for more information.\n\n" +
			ItemDocsSPNSupport,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The " + DomainWorkspaceAssignmentsName + " ID.",
				Computed:            true,
				CustomType:          customtypes.UUIDType{},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"domain_id": schema.StringAttribute{
				MarkdownDescription: "The Domain ID. " + common.DocsRequiresReplace,
				Required:            true,
				CustomType:          customtypes.UUIDType{},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"workspace_ids": schema.SetAttribute{
				MarkdownDescription: "The list of Workspaces.",
				Required:            true,
				ElementType:         customtypes.UUIDType{},
			},
			"timeouts": timeouts.AttributesAll(ctx),
		},
	}
}

func (r *resourceDomainWorkspaceAssignments) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
	r.client = fabadmin.NewClientFactoryWithClient(*pConfigData.FabricClient).NewDomainsClient()
}

func (r *resourceDomainWorkspaceAssignments) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "CREATE", map[string]any{
		"action": "start",
	})
	tflog.Trace(ctx, "CREATE", map[string]any{
		"config": req.Config,
		"plan":   req.Plan,
	})

	var plan resourceDomainWorkspaceAssignmentsModel

	if resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...); resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := plan.Timeouts.Create(ctx, r.pConfigData.Timeout)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var reqCreate requestCreateDomainWorkspaceAssignments

	if resp.Diagnostics.Append(reqCreate.set(ctx, plan)...); resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.AssignDomainWorkspacesByIDs(ctx, plan.DomainID.ValueString(), reqCreate.AssignDomainWorkspacesByIDsRequest, nil)
	if resp.Diagnostics.Append(utils.GetDiagsFromError(ctx, err, utils.OperationCreate, nil)...); resp.Diagnostics.HasError() {
		return
	}

	plan.ID = plan.DomainID

	if resp.Diagnostics.Append(r.list(ctx, &plan)...); resp.Diagnostics.HasError() {
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

func (r *resourceDomainWorkspaceAssignments) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "READ", map[string]any{
		"action": "start",
	})
	tflog.Trace(ctx, "READ", map[string]any{
		"state": req.State,
	})

	var state resourceDomainWorkspaceAssignmentsModel

	if resp.Diagnostics.Append(req.State.Get(ctx, &state)...); resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := state.Timeouts.Read(ctx, r.pConfigData.Timeout)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if resp.Diagnostics.Append(r.list(ctx, &state)...); resp.Diagnostics.HasError() {
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

func (r *resourceDomainWorkspaceAssignments) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "UPDATE", map[string]any{
		"action": "start",
	})
	tflog.Trace(ctx, "UPDATE", map[string]any{
		"config": req.Config,
		"plan":   req.Plan,
		"state":  req.State,
	})

	var plan, state resourceDomainWorkspaceAssignmentsModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := state.Timeouts.Update(ctx, r.pConfigData.Timeout)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	added, diags := r.diffWorkspaces(ctx, plan.WorkspaceIDs, state.WorkspaceIDs)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	removed, diags := r.diffWorkspaces(ctx, state.WorkspaceIDs, plan.WorkspaceIDs)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	if len(added.Elements()) > 0 {
		state.WorkspaceIDs = added

		var reqCreate requestCreateDomainWorkspaceAssignments

		if resp.Diagnostics.Append(reqCreate.set(ctx, state)...); resp.Diagnostics.HasError() {
			return
		}

		_, err := r.client.AssignDomainWorkspacesByIDs(ctx, plan.DomainID.ValueString(), reqCreate.AssignDomainWorkspacesByIDsRequest, nil)
		if resp.Diagnostics.Append(utils.GetDiagsFromError(ctx, err, utils.OperationUpdate, nil)...); resp.Diagnostics.HasError() {
			return
		}
	}

	if len(removed.Elements()) > 0 {
		state.WorkspaceIDs = removed

		var reqDelete requestDeleteDomainWorkspaceAssignments

		if resp.Diagnostics.Append(reqDelete.set(ctx, state)...); resp.Diagnostics.HasError() {
			return
		}

		_, err := r.client.UnassignDomainWorkspacesByIDs(ctx, plan.DomainID.ValueString(), &fabadmin.DomainsClientUnassignDomainWorkspacesByIDsOptions{
			UnassignDomainWorkspacesByIDsRequest: &reqDelete.UnassignDomainWorkspacesByIDsRequest,
		})
		if resp.Diagnostics.Append(utils.GetDiagsFromError(ctx, err, utils.OperationUpdate, nil)...); resp.Diagnostics.HasError() {
			return
		}
	}

	if resp.Diagnostics.Append(r.list(ctx, &plan)...); resp.Diagnostics.HasError() {
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

func (r *resourceDomainWorkspaceAssignments) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "DELETE", map[string]any{
		"action": "start",
	})
	tflog.Trace(ctx, "DELETE", map[string]any{
		"state": req.State,
	})

	var state resourceDomainWorkspaceAssignmentsModel

	if resp.Diagnostics.Append(req.State.Get(ctx, &state)...); resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := state.Timeouts.Delete(ctx, r.pConfigData.Timeout)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if len(state.WorkspaceIDs.Elements()) > 0 {
		var reqDelete requestDeleteDomainWorkspaceAssignments

		reqDelete.set(ctx, state)

		_, err := r.client.UnassignDomainWorkspacesByIDs(ctx, state.DomainID.ValueString(), &fabadmin.DomainsClientUnassignDomainWorkspacesByIDsOptions{
			UnassignDomainWorkspacesByIDsRequest: &reqDelete.UnassignDomainWorkspacesByIDsRequest,
		})
		if resp.Diagnostics.Append(utils.GetDiagsFromError(ctx, err, utils.OperationDelete, nil)...); resp.Diagnostics.HasError() {
			return
		}
	}

	resp.State.RemoveResource(ctx)

	tflog.Debug(ctx, "DELETE", map[string]any{
		"action": "end",
	})
}

func (r *resourceDomainWorkspaceAssignments) diffWorkspaces(ctx context.Context, slice1, slice2 types.Set) (types.Set, diag.Diagnostics) {
	s1 := make([]customtypes.UUID, 0, len(slice1.Elements()))

	diags := slice1.ElementsAs(ctx, &s1, false)
	if diags.HasError() {
		return types.SetNull(customtypes.UUIDType{}), diags
	}

	s2 := make([]customtypes.UUID, 0, len(slice1.Elements()))

	diags = slice2.ElementsAs(ctx, &s2, false)
	if diags.HasError() {
		return types.SetNull(customtypes.UUIDType{}), diags
	}

	m := make(map[string]bool)
	for _, item := range s2 {
		m[item.ValueString()] = true
	}

	elements := []customtypes.UUID{}

	for _, item := range s1 {
		if !m[item.ValueString()] {
			elements = append(elements, item)
		}
	}

	diff, diags := types.SetValueFrom(ctx, customtypes.UUIDType{}, elements)
	if diags.HasError() {
		return types.SetNull(customtypes.UUIDType{}), diags
	}

	return diff, nil
}

func (r *resourceDomainWorkspaceAssignments) list(ctx context.Context, model *resourceDomainWorkspaceAssignmentsModel) diag.Diagnostics {
	tflog.Trace(ctx, "getting Domain Workspace Assignments")

	respList, err := r.client.ListDomainWorkspaces(ctx, model.DomainID.ValueString(), nil)
	if diags := utils.GetDiagsFromError(ctx, err, utils.OperationList, nil); diags.HasError() {
		return diags
	}

	return model.setWorkspaces(ctx, respList)
}
