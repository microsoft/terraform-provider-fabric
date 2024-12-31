// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package workspace

import (
	"context"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
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

var _ resource.ResourceWithConfigure = (*resourceWorkspaceGit)(nil)

type resourceWorkspaceGit struct {
	pConfigData *pconfig.ProviderData
	client      *fabcore.GitClient
}

func NewResourceWorkspaceGit() resource.Resource {
	return &resourceWorkspaceGit{}
}

func (r *resourceWorkspaceGit) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + WorkspaceGitTFName
}

func (r *resourceWorkspaceGit) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	possibleInitializationStrategyValues := utils.RemoveSliceByValue(fabcore.PossibleInitializationStrategyValues(), fabcore.InitializationStrategyNone)

	resp.Schema = schema.Schema{
		MarkdownDescription: "Manage a Fabric " + WorkspaceGitName + ".\n\n" +
			"See [" + WorkspaceGitName + "](" + WorkspaceGitDocsURL + ") for more information.\n\n" +
			common.DocsSPNNotSupported,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:   true,
				CustomType: customtypes.UUIDType{},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"workspace_id": schema.StringAttribute{
				MarkdownDescription: "The Workspace ID.",
				Required:            true,
				CustomType:          customtypes.UUIDType{},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"initialization_strategy": schema.StringAttribute{
				MarkdownDescription: "The initialization strategy. Accepted values: " + utils.ConvertStringSlicesToString(possibleInitializationStrategyValues, true, true),
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf(utils.ConvertEnumsToStringSlices(possibleInitializationStrategyValues, true)...),
				},
			},
			"git_sync_details": schema.SingleNestedAttribute{
				MarkdownDescription: "The git sync details.",
				Computed:            true,
				CustomType:          supertypes.NewSingleNestedObjectTypeOf[gitSyncDetailsModel](ctx),
				Attributes: map[string]schema.Attribute{
					"head": schema.StringAttribute{
						MarkdownDescription: "The git head.",
						Computed:            true,
					},
					"last_sync_time": schema.StringAttribute{
						MarkdownDescription: "The last sync time.",
						Computed:            true,
						CustomType:          timetypes.RFC3339Type{},
					},
				},
			},
			"git_provider_details": schema.SingleNestedAttribute{
				MarkdownDescription: "The Git provider details.",
				Required:            true,
				CustomType:          supertypes.NewSingleNestedObjectTypeOf[gitProviderDetailsModel](ctx),
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.RequiresReplace(),
				},
				Attributes: map[string]schema.Attribute{
					"git_provider_type": schema.StringAttribute{
						MarkdownDescription: "The Git provider type. Accepted values: " + utils.ConvertStringSlicesToString(utils.RemoveSliceByValue(fabcore.PossibleGitProviderTypeValues(), fabcore.GitProviderTypeGitHub), true, true),
						Required:            true,
						Validators: []validator.String{
							stringvalidator.OneOf(utils.ConvertEnumsToStringSlices(utils.RemoveSliceByValue(fabcore.PossibleGitProviderTypeValues(), fabcore.GitProviderTypeGitHub), true)...),
						},
					},
					"organization_name": schema.StringAttribute{
						MarkdownDescription: "The organization name.",
						Required:            true,
						Validators: []validator.String{
							stringvalidator.LengthAtMost(100),
						},
					},
					"project_name": schema.StringAttribute{
						MarkdownDescription: "The project name.",
						Required:            true,
						Validators: []validator.String{
							stringvalidator.LengthAtMost(100),
						},
					},
					"repository_name": schema.StringAttribute{
						MarkdownDescription: "The repository name.",
						Required:            true,
						Validators: []validator.String{
							stringvalidator.LengthAtMost(128),
						},
					},
					"branch_name": schema.StringAttribute{
						MarkdownDescription: "The branch name.",
						Required:            true,
						Validators: []validator.String{
							stringvalidator.LengthAtMost(250),
						},
					},
					"directory_name": schema.StringAttribute{
						MarkdownDescription: "The directory name.",
						Required:            true,
						Validators: []validator.String{
							stringvalidator.LengthAtMost(256),
							stringvalidator.RegexMatches(
								regexp.MustCompile(`^/.*`),
								"Directory name path must starts with forward slash '/'.",
							),
						},
					},
				},
			},
			"git_connection_state": schema.StringAttribute{
				MarkdownDescription: "The git connection state. Possible values: " + utils.ConvertStringSlicesToString(fabcore.PossibleGitConnectionStateValues(), true, true),
				Computed:            true,
			},
			"timeouts": timeouts.AttributesAll(ctx),
		},
	}
}

func (r *resourceWorkspaceGit) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
	r.client = fabcore.NewClientFactoryWithClient(*pConfigData.FabricClient).NewGitClient()
}

func (r *resourceWorkspaceGit) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "CREATE", map[string]any{
		"action": "start",
	})
	tflog.Trace(ctx, "CREATE", map[string]any{
		"config": req.Config,
		"plan":   req.Plan,
	})

	var plan resourceWorkspaceGitModel

	if resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...); resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := plan.Timeouts.Create(ctx, r.pConfigData.Timeout)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Connect.
	var reqGitConnect requestGitConnect

	if resp.Diagnostics.Append(reqGitConnect.set(ctx, plan)...); resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.Connect(ctx, plan.WorkspaceID.ValueString(), reqGitConnect.GitConnectRequest, nil)
	if resp.Diagnostics.Append(utils.GetDiagsFromError(ctx, err, utils.OperationCreate, nil)...); resp.Diagnostics.HasError() {
		return
	}

	// Initialize.
	var reqGitInitialize requestGitInitialize

	reqGitInitialize.set(plan)

	gitInitResp, err := r.client.InitializeConnection(ctx, plan.WorkspaceID.ValueString(), &fabcore.GitClientBeginInitializeConnectionOptions{
		GitInitializeConnectionRequest: &reqGitInitialize.InitializeGitConnectionRequest,
	})
	if resp.Diagnostics.Append(utils.GetDiagsFromError(ctx, err, utils.OperationCreate, nil)...); resp.Diagnostics.HasError() {
		return
	}

	// Git commit.
	switch *gitInitResp.RequiredAction {
	case fabcore.RequiredActionCommitToGit: // Commit to Git.
		var reqGitCommitTo requestGitCommitTo

		reqGitCommitTo.set(gitInitResp.WorkspaceHead)

		_, err = r.client.CommitToGit(ctx, plan.WorkspaceID.ValueString(), reqGitCommitTo.CommitToGitRequest, nil)

	case fabcore.RequiredActionUpdateFromGit: // Update from Git.
		var reqGitUpdateFrom requestGitUpdateFrom

		reqGitUpdateFrom.set(gitInitResp.RemoteCommitHash, plan.InitializationStrategy.ValueStringPointer())

		_, err = r.client.UpdateFromGit(ctx, plan.WorkspaceID.ValueString(), reqGitUpdateFrom.UpdateFromGitRequest, nil)
	case fabcore.RequiredActionNone:
		// Do nothing.
	default:
		resp.Diagnostics.AddError(
			common.ErrorCreateHeader,
			fmt.Sprintf("Unsupported required git action '%s'.", *gitInitResp.RequiredAction),
		)
	}

	if resp.Diagnostics.Append(utils.GetDiagsFromError(ctx, err, utils.OperationCreate, nil)...); resp.Diagnostics.HasError() {
		return
	}

	if resp.Diagnostics.Append(r.get(ctx, &plan)...); resp.Diagnostics.HasError() {
		return
	}

	plan.ID = customtypes.NewUUIDValue(plan.WorkspaceID.ValueString())

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)

	tflog.Debug(ctx, "CREATE", map[string]any{
		"action": "end",
	})

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceWorkspaceGit) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "READ", map[string]any{
		"action": "start",
	})
	tflog.Trace(ctx, "READ", map[string]any{
		"state": req.State,
	})

	var state resourceWorkspaceGitModel

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
	if utils.IsErrNotFound(state.ID.ValueString(), &diags, fabcore.ErrGit.GitProviderResourceNotFound) {
		resp.State.RemoveResource(ctx)

		resp.Diagnostics.Append(diags...)

		return
	}

	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	if !state.GitConnectionState.IsNull() && !state.GitConnectionState.IsUnknown() && state.GitConnectionState.ValueString() != (string)(fabcore.GitConnectionStateConnectedAndInitialized) {
		resp.Diagnostics.AddWarning(
			"Unexpected Git connection state",
			fmt.Sprintf("Git connection state is '%s'.\nIt may have been deleted outside of Terraform. Removing object from state.", state.GitConnectionState.ValueString()),
		)

		resp.State.RemoveResource(ctx)

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)

	tflog.Debug(ctx, "READ", map[string]any{
		"action": "end",
	})

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceWorkspaceGit) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "UPDATE", map[string]any{
		"action": "start",
	})
	tflog.Trace(ctx, "UPDATE", map[string]any{
		"config": req.Config,
		"plan":   req.Plan,
		"state":  req.State,
	})

	// in real world, this should not reach here
	resp.Diagnostics.AddError(
		common.ErrorUpdateHeader,
		"Update is not supported. Requires delete and recreate.",
	)

	tflog.Debug(ctx, "UPDATE", map[string]any{
		"action": "end",
	})

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceWorkspaceGit) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "DELETE", map[string]any{
		"action": "start",
	})
	tflog.Trace(ctx, "DELETE", map[string]any{
		"state": req.State,
	})

	var state resourceWorkspaceGitModel

	if resp.Diagnostics.Append(req.State.Get(ctx, &state)...); resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := state.Timeouts.Delete(ctx, r.pConfigData.Timeout)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	_, err := r.client.Disconnect(ctx, state.WorkspaceID.ValueString(), nil)

	diags = utils.GetDiagsFromError(ctx, err, utils.OperationDelete, fabcore.ErrGit.WorkspaceNotConnectedToGit)

	if diags.HasError() && !utils.IsErr(diags, fabcore.ErrGit.WorkspaceNotConnectedToGit) {
		resp.Diagnostics.Append(diags...)

		return
	}

	resp.State.RemoveResource(ctx)

	tflog.Debug(ctx, "DELETE", map[string]any{
		"action": "end",
	})
}

func (r *resourceWorkspaceGit) get(ctx context.Context, model *resourceWorkspaceGitModel) diag.Diagnostics {
	tflog.Trace(ctx, fmt.Sprintf("getting %s for Workspace ID: %s", WorkspaceGitName, model.WorkspaceID.ValueString()))

	respGet, err := r.client.GetConnection(ctx, model.WorkspaceID.ValueString(), nil)
	if diags := utils.GetDiagsFromError(ctx, err, utils.OperationRead, fabcore.ErrGit.GitProviderResourceNotFound); diags.HasError() {
		return diags
	}

	return model.set(ctx, respGet.GitConnection)
}
