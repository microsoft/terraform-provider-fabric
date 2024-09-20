// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package domain

import (
	"context"
	"fmt"

	supertypes "github.com/FrangipaneTeam/terraform-plugin-framework-supertypes"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	fabadmin "github.com/microsoft/fabric-sdk-go/fabric/admin"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
	pconfig "github.com/microsoft/terraform-provider-fabric/internal/provider/config"
)

const (
	DomainRoleAssignmentsName   = "Domain Role Assignments"
	DomainRoleAssignmentsTFName = "domain_role_assignments"
)

// Ensure the implementation satisfies the expected interfaces.
var _ resource.ResourceWithConfigure = (*resourceDomainRoleAssignments)(nil)

type resourceDomainRoleAssignments struct {
	pConfigData *pconfig.ProviderData
	client      *fabadmin.DomainsClient
}

func NewResourceDomainRoleAssignments() resource.Resource {
	return &resourceDomainRoleAssignments{}
}

func (r *resourceDomainRoleAssignments) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + DomainRoleAssignmentsTFName
}

func (r *resourceDomainRoleAssignments) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	possiblePrincipalTypeValues := utils.RemoveSlicesByValues(fabadmin.PossiblePrincipalTypeValues(), []fabadmin.PrincipalType{fabadmin.PrincipalTypeServicePrincipal, fabadmin.PrincipalTypeServicePrincipalProfile})

	resp.Schema = schema.Schema{
		MarkdownDescription: "Manage a Fabric " + DomainRoleAssignmentsName + ".\n\n" +
			"See [" + ItemName + "](" + ItemDocsURL + ") for more information.\n\n" +
			ItemDocsSPNSupport,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The " + DomainRoleAssignmentsName + " ID.",
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
			"role": schema.StringAttribute{
				MarkdownDescription: "The Role of the principals. Accepted values: " + utils.ConvertStringSlicesToString(fabadmin.PossibleDomainRoleValues(), true, true) + ". " + common.DocsRequiresReplace,
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf(utils.ConvertEnumsToStringSlices(fabadmin.PossibleDomainRoleValues(), false)...),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"principals": schema.SetNestedAttribute{
				MarkdownDescription: "The list of Principals.",
				Required:            true,
				CustomType:          supertypes.NewSetNestedObjectTypeOf[principalModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "The Principal ID. " + common.DocsRequiresReplace,
							Required:            true,
							CustomType:          customtypes.UUIDType{},
						},
						"type": schema.StringAttribute{
							MarkdownDescription: "ThePrincipal type. Accepted values: " + utils.ConvertStringSlicesToString(possiblePrincipalTypeValues, true, true) + ". " + common.DocsRequiresReplace,
							Required:            true,
							Validators: []validator.String{
								stringvalidator.OneOf(utils.ConvertEnumsToStringSlices(possiblePrincipalTypeValues, false)...),
							},
						},
					},
				},
			},
			"timeouts": timeouts.AttributesAll(ctx),
		},
	}
}

func (r *resourceDomainRoleAssignments) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *resourceDomainRoleAssignments) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "CREATE", map[string]any{
		"action": "start",
	})
	tflog.Trace(ctx, "CREATE", map[string]any{
		"config": req.Config,
		"plan":   req.Plan,
	})

	var plan resourceDomainRoleAssignmentsModel

	if resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...); resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := plan.Timeouts.Create(ctx, r.pConfigData.Timeout)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if resp.Diagnostics.Append(r.checkDomainSupport(ctx, plan)...); resp.Diagnostics.HasError() {
		return
	}

	var reqCreate requestCreateDomainRoleAssignments

	if resp.Diagnostics.Append(reqCreate.set(ctx, plan)...); resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.RoleAssignmentsBulkAssign(ctx, plan.DomainID.ValueString(), reqCreate.DomainRoleAssignmentRequest, nil)
	if resp.Diagnostics.Append(utils.GetDiagsFromError(ctx, err, utils.OperationCreate, nil)...); resp.Diagnostics.HasError() {
		return
	}

	plan.ID = plan.DomainID

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)

	tflog.Debug(ctx, "CREATE", map[string]any{
		"action": "end",
	})

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceDomainRoleAssignments) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "READ", map[string]any{
		"action": "start",
	})
	tflog.Trace(ctx, "READ", map[string]any{
		"state": req.State,
	})

	var state resourceDomainRoleAssignmentsModel

	if resp.Diagnostics.Append(req.State.Get(ctx, &state)...); resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := state.Timeouts.Read(ctx, r.pConfigData.Timeout)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// DO NOTHING
	// This resource does not have get/list API

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)

	tflog.Debug(ctx, "READ", map[string]any{
		"action": "end",
	})

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceDomainRoleAssignments) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "UPDATE", map[string]any{
		"action": "start",
	})
	tflog.Trace(ctx, "UPDATE", map[string]any{
		"config": req.Config,
		"plan":   req.Plan,
		"state":  req.State,
	})

	var plan, state resourceDomainRoleAssignmentsModel

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

	added, diags := r.diffPrincipals(ctx, plan.Principals, state.Principals)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	removed, diags := r.diffPrincipals(ctx, state.Principals, plan.Principals)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	if len(added) > 0 {
		state.Principals = supertypes.NewSetNestedObjectValueOfValueSlice(ctx, added)

		var reqCreate requestCreateDomainRoleAssignments

		if resp.Diagnostics.Append(reqCreate.set(ctx, state)...); resp.Diagnostics.HasError() {
			return
		}

		_, err := r.client.RoleAssignmentsBulkAssign(ctx, plan.DomainID.ValueString(), reqCreate.DomainRoleAssignmentRequest, nil)
		if resp.Diagnostics.Append(utils.GetDiagsFromError(ctx, err, utils.OperationUpdate, nil)...); resp.Diagnostics.HasError() {
			return
		}
	}

	if len(removed) > 0 {
		state.Principals = supertypes.NewSetNestedObjectValueOfValueSlice(ctx, removed)

		var reqDelete requestDeleteDomainRoleAssignments

		if resp.Diagnostics.Append(reqDelete.set(ctx, state)...); resp.Diagnostics.HasError() {
			return
		}

		_, err := r.client.RoleAssignmentsBulkUnassign(ctx, plan.DomainID.ValueString(), reqDelete.DomainRoleUnassignmentRequest, nil)
		if resp.Diagnostics.Append(utils.GetDiagsFromError(ctx, err, utils.OperationDelete, nil)...); resp.Diagnostics.HasError() {
			return
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)

	tflog.Debug(ctx, "UPDATE", map[string]any{
		"action": "end",
	})

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceDomainRoleAssignments) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "DELETE", map[string]any{
		"action": "start",
	})
	tflog.Trace(ctx, "DELETE", map[string]any{
		"state": req.State,
	})

	var state resourceDomainRoleAssignmentsModel

	if resp.Diagnostics.Append(req.State.Get(ctx, &state)...); resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := state.Timeouts.Delete(ctx, r.pConfigData.Timeout)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var reqDelete requestDeleteDomainRoleAssignments

	reqDelete.set(ctx, state)

	_, err := r.client.RoleAssignmentsBulkUnassign(ctx, state.DomainID.ValueString(), reqDelete.DomainRoleUnassignmentRequest, nil)
	if resp.Diagnostics.Append(utils.GetDiagsFromError(ctx, err, utils.OperationDelete, nil)...); resp.Diagnostics.HasError() {
		return
	}

	resp.State.RemoveResource(ctx)

	tflog.Debug(ctx, "DELETE", map[string]any{
		"action": "end",
	})
}

func (r *resourceDomainRoleAssignments) checkDomainSupport(ctx context.Context, model resourceDomainRoleAssignmentsModel) diag.Diagnostics {
	var diags diag.Diagnostics

	respGet, err := r.client.GetDomain(ctx, model.DomainID.ValueString(), nil)
	if diags := utils.GetDiagsFromError(ctx, err, utils.OperationDelete, nil); diags.HasError() {
		return diags
	}

	if respGet.ParentDomainID != nil {
		diags.AddError(
			"Subdomains are not supported",
			"Role Assignment is not supported for subdomains. Please use the root-level domain.",
		)

		return diags
	}

	if *respGet.ContributorsScope != fabadmin.ContributorsScopeTypeSpecificUsersAndGroups {
		diags.AddError(
			"Unsupported Domain Contributors Scope",
			"Only '"+string(fabadmin.ContributorsScopeTypeSpecificUsersAndGroups)+"' Domain Contributors Scope is supported for Role Assignment.",
		)

		return diags
	}

	return nil
}

func (r *resourceDomainRoleAssignments) diffPrincipals(ctx context.Context, slice1, slice2 supertypes.SetNestedObjectValueOf[principalModel]) ([]principalModel, diag.Diagnostics) {
	s1, diags := slice1.Get(ctx)
	if diags.HasError() {
		return nil, diags
	}

	s2, diags := slice2.Get(ctx)
	if diags.HasError() {
		return nil, diags
	}

	m := make(map[string]bool)
	for _, item := range s2 {
		m[item.ID.ValueString()] = true
	}

	var diff []principalModel

	for _, item := range s1 {
		if !m[item.ID.ValueString()] {
			slice := principalModel{
				ID:   item.ID,
				Type: item.Type,
			}

			diff = append(diff, slice)
		}
	}

	return diff, nil
}