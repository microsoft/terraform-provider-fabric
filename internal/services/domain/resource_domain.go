// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package domain

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	fabadmin "github.com/microsoft/fabric-sdk-go/fabric/admin"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
	pconfig "github.com/microsoft/terraform-provider-fabric/internal/provider/config"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.ResourceWithConfigure   = (*resourceDomain)(nil)
	_ resource.ResourceWithImportState = (*resourceDomain)(nil)
)

type resourceDomain struct {
	pConfigData         *pconfig.ProviderData
	client              *fabadmin.DomainsClient
	Name                string
	TFName              string
	MarkdownDescription string
	IsPreview           bool
}

func NewResourceDomain() resource.Resource {
	markdownDescription := "Manage a Fabric " + ItemName + ".\n\n" +
		"Use this resource to manage [" + ItemName + "](" + ItemDocsURL + ").\n\n" +
		ItemDocsSPNSupport

	return &resourceDomain{
		Name:                ItemName,
		TFName:              ItemTFName,
		MarkdownDescription: fabricitem.GetResourcePreviewNote(markdownDescription, ItemPreview),
		IsPreview:           ItemPreview,
	}
}

func (r *resourceDomain) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + r.TFName
}

func (r *resourceDomain) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: r.MarkdownDescription,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The " + r.Name + " ID.",
				Computed:            true,
				CustomType:          customtypes.UUIDType{},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "The " + r.Name + " display name.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(40),
				},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The " + r.Name + " description.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
				Validators: []validator.String{
					stringvalidator.LengthAtMost(256),
				},
			},
			"parent_domain_id": schema.StringAttribute{
				MarkdownDescription: "The " + r.Name + " parent ID.",
				Optional:            true,
				CustomType:          customtypes.UUIDType{},
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRoot("contributors_scope")),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"contributors_scope": schema.StringAttribute{
				MarkdownDescription: "The " + r.Name + " contributors scope. Possible values: " + utils.ConvertStringSlicesToString(fabadmin.PossibleContributorsScopeTypeValues(), true, true) + ".\n\n" +
					"-> Contributors scope can only be set at the root domain level.",
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.OneOf(utils.ConvertEnumsToStringSlices(fabadmin.PossibleContributorsScopeTypeValues(), false)...),
				},
			},
			"timeouts": timeouts.AttributesAll(ctx),
		},
	}
}

func (r *resourceDomain) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	if resp.Diagnostics.Append(fabricitem.IsPreviewMode(r.Name, r.IsPreview, r.pConfigData.Preview)...); resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceDomain) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "CREATE", map[string]any{
		"action": "start",
	})

	var plan, state resourceDomainModel

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

	var reqCreate requestCreateDomain

	reqCreate.set(plan)

	respCreate, err := r.client.CreateDomain(ctx, reqCreate.CreateDomainRequest, nil)
	if resp.Diagnostics.Append(utils.GetDiagsFromError(ctx, err, utils.OperationCreate, nil)...); resp.Diagnostics.HasError() {
		return
	}

	state.set(respCreate.Domain)

	if !(plan.ContributorsScope.IsNull() || plan.ContributorsScope.IsUnknown()) && plan.ContributorsScope.ValueString() != string(fabadmin.ContributorsScopeTypeAllTenant) && plan.ParentDomainID.IsNull() {
		tflog.Trace(ctx, "setting "+ItemName+" contributors scope")

		var reqUpdate requestUpdateDomain

		reqUpdate.set(plan)

		respUpdate, err := r.client.UpdateDomain(ctx, state.ID.ValueString(), reqUpdate.UpdateDomainRequest, nil)
		if resp.Diagnostics.Append(utils.GetDiagsFromError(ctx, err, utils.OperationCreate, nil)...); resp.Diagnostics.HasError() {
			return
		}

		state.set(respUpdate.Domain)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)

	tflog.Debug(ctx, "CREATE", map[string]any{
		"action": "end",
	})

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceDomain) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "READ", map[string]any{
		"action": "start",
	})
	tflog.Trace(ctx, "READ", map[string]any{
		"state": req.State,
	})

	var state resourceDomainModel

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

func (r *resourceDomain) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "UPDATE", map[string]any{
		"action": "start",
	})

	var plan resourceDomainModel

	if resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...); resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := plan.Timeouts.Update(ctx, r.pConfigData.Timeout)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var reqUpdate requestUpdateDomain

	reqUpdate.set(plan)

	respUpdate, err := r.client.UpdateDomain(ctx, plan.ID.ValueString(), reqUpdate.UpdateDomainRequest, nil)
	if resp.Diagnostics.Append(utils.GetDiagsFromError(ctx, err, utils.OperationUpdate, nil)...); resp.Diagnostics.HasError() {
		return
	}

	plan.set(respUpdate.Domain)

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)

	tflog.Debug(ctx, "UPDATE", map[string]any{
		"action": "end",
	})

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceDomain) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "DELETE", map[string]any{
		"action": "start",
	})

	var state resourceDomainModel

	if resp.Diagnostics.Append(req.State.Get(ctx, &state)...); resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := state.Timeouts.Delete(ctx, r.pConfigData.Timeout)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	_, err := r.client.DeleteDomain(ctx, state.ID.ValueString(), nil)
	if resp.Diagnostics.Append(utils.GetDiagsFromError(ctx, err, utils.OperationDelete, nil)...); resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "DELETE", map[string]any{
		"action": "end",
	})
}

func (r *resourceDomain) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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

func (r *resourceDomain) get(ctx context.Context, model *resourceDomainModel) error {
	tflog.Trace(ctx, "getting "+ItemName)

	respGet, err := r.client.GetDomain(ctx, model.ID.ValueString(), nil)
	if err != nil {
		return err
	}

	model.set(respGet.Domain)

	return nil
}
