// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package workspacegit

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
	pconfig "github.com/microsoft/terraform-provider-fabric/internal/provider/config"
)

var _ resource.ResourceWithConfigure = (*resourceWorkspaceGit)(nil)

type resourceWorkspaceGit struct {
	pConfigData *pconfig.ProviderData
	client      *fabcore.GitClient
	TypeInfo    tftypeinfo.TFTypeInfo
}

func NewResourceWorkspaceGit() resource.Resource {
	return &resourceWorkspaceGit{
		TypeInfo: ItemTypeInfo,
	}
}

func (r *resourceWorkspaceGit) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = r.TypeInfo.FullTypeName(false)
}

func (r *resourceWorkspaceGit) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = itemSchema().GetResource(ctx)
	// possibleInitializationStrategyValues := utils.RemoveSliceByValue(fabcore.PossibleInitializationStrategyValues(), fabcore.InitializationStrategyNone)
	// gitProviderTypeAttPath := path.MatchRoot("git_provider_details").AtName("git_provider_type")
	// gitProviderTypeAzureDevOps := types.StringValue(string(fabcore.GitProviderTypeAzureDevOps))
	// gitProviderTypeGitHub := types.StringValue(string(fabcore.GitProviderTypeGitHub))
	//
	//	resp.Schema = schema.Schema{
	//		MarkdownDescription: fabricitem.NewResourceMarkdownDescription(r.TypeInfo, false),
	//		Attributes: map[string]schema.Attribute{
	//			"id": schema.StringAttribute{
	//				Computed:   true,
	//				CustomType: customtypes.UUIDType{},
	//				PlanModifiers: []planmodifier.String{
	//					stringplanmodifier.UseStateForUnknown(),
	//				},
	//			},
	//			"workspace_id": schema.StringAttribute{
	//				MarkdownDescription: "The Workspace ID.",
	//				Required:            true,
	//				CustomType:          customtypes.UUIDType{},
	//				PlanModifiers: []planmodifier.String{
	//					stringplanmodifier.RequiresReplace(),
	//				},
	//			},
	//			"initialization_strategy": schema.StringAttribute{
	//				MarkdownDescription: "The initialization strategy. Accepted values: " + utils.ConvertStringSlicesToString(possibleInitializationStrategyValues, true, true),
	//				Required:            true,
	//				PlanModifiers: []planmodifier.String{
	//					stringplanmodifier.RequiresReplace(),
	//				},
	//				Validators: []validator.String{
	//					stringvalidator.OneOf(utils.ConvertEnumsToStringSlices(possibleInitializationStrategyValues, true)...),
	//				},
	//			},
	//			"git_sync_details": schema.SingleNestedAttribute{
	//				MarkdownDescription: "The git sync details.",
	//				Computed:            true,
	//				CustomType:          supertypes.NewSingleNestedObjectTypeOf[gitSyncDetailsModel](ctx),
	//				Attributes: map[string]schema.Attribute{
	//					"head": schema.StringAttribute{
	//						MarkdownDescription: "The git head.",
	//						Computed:            true,
	//					},
	//					"last_sync_time": schema.StringAttribute{
	//						MarkdownDescription: "The last sync time.",
	//						Computed:            true,
	//						CustomType:          timetypes.RFC3339Type{},
	//					},
	//				},
	//			},
	//			"git_provider_details": schema.SingleNestedAttribute{
	//				MarkdownDescription: "The Git provider details.",
	//				Required:            true,
	//				CustomType:          supertypes.NewSingleNestedObjectTypeOf[gitProviderDetailsModel](ctx),
	//				PlanModifiers: []planmodifier.Object{
	//					objectplanmodifier.RequiresReplace(),
	//				},
	//				Attributes: map[string]schema.Attribute{
	//					"git_provider_type": schema.StringAttribute{
	//						MarkdownDescription: "The Git provider type. Accepted values: " + utils.ConvertStringSlicesToString(fabcore.PossibleGitProviderTypeValues(), true, true),
	//						Required:            true,
	//						Validators: []validator.String{
	//							stringvalidator.OneOf(utils.ConvertEnumsToStringSlices(fabcore.PossibleGitProviderTypeValues(), true)...),
	//						},
	//					},
	//					"organization_name": schema.StringAttribute{
	//						MarkdownDescription: "The Azure DevOps organization name.",
	//						Computed:            true,
	//						Optional:            true,
	//						Validators: []validator.String{
	//							stringvalidator.LengthAtMost(100),
	//							superstringvalidator.NullIfAttributeIsOneOf(
	//								gitProviderTypeAttPath,
	//								[]attr.Value{gitProviderTypeGitHub},
	//							),
	//							superstringvalidator.RequireIfAttributeIsOneOf(
	//								gitProviderTypeAttPath,
	//								[]attr.Value{gitProviderTypeAzureDevOps},
	//							),
	//						},
	//					},
	//					"project_name": schema.StringAttribute{
	//						MarkdownDescription: "The Azure DevOps project name.",
	//						Computed:            true,
	//						Optional:            true,
	//						Validators: []validator.String{
	//							stringvalidator.LengthAtMost(100),
	//							superstringvalidator.NullIfAttributeIsOneOf(
	//								gitProviderTypeAttPath,
	//								[]attr.Value{gitProviderTypeGitHub},
	//							),
	//							superstringvalidator.RequireIfAttributeIsOneOf(
	//								gitProviderTypeAttPath,
	//								[]attr.Value{gitProviderTypeAzureDevOps},
	//							),
	//						},
	//					},
	//					"owner_name": schema.StringAttribute{
	//						MarkdownDescription: "The GitHub owner name.",
	//						Computed:            true,
	//						Optional:            true,
	//						Validators: []validator.String{
	//							stringvalidator.LengthAtMost(100),
	//							superstringvalidator.NullIfAttributeIsOneOf(
	//								gitProviderTypeAttPath,
	//								[]attr.Value{gitProviderTypeAzureDevOps},
	//							),
	//							superstringvalidator.RequireIfAttributeIsOneOf(
	//								gitProviderTypeAttPath,
	//								[]attr.Value{gitProviderTypeGitHub},
	//							),
	//						},
	//					},
	//					"repository_name": schema.StringAttribute{
	//						MarkdownDescription: "The repository name.",
	//						Required:            true,
	//						Validators: []validator.String{
	//							stringvalidator.LengthAtMost(128),
	//						},
	//					},
	//					"branch_name": schema.StringAttribute{
	//						MarkdownDescription: "The branch name.",
	//						Required:            true,
	//						Validators: []validator.String{
	//							stringvalidator.LengthAtMost(250),
	//						},
	//					},
	//					"directory_name": schema.StringAttribute{
	//						MarkdownDescription: "The directory name.",
	//						Required:            true,
	//						Validators: []validator.String{
	//							stringvalidator.LengthAtMost(256),
	//							stringvalidator.RegexMatches(
	//								regexp.MustCompile(`^/.*`),
	//								"Directory name path must starts with forward slash '/'.",
	//							),
	//						},
	//					},
	//				},
	//			},
	//			"git_credentials": schema.SingleNestedAttribute{
	//				MarkdownDescription: "The Git credentials details.",
	//				Computed:            true,
	//				Optional:            true,
	//				Validators: []validator.Object{
	//					superobjectvalidator.NullIfAttributeIsOneOf(
	//						gitProviderTypeAttPath,
	//						[]attr.Value{gitProviderTypeAzureDevOps},
	//					),
	//					superobjectvalidator.RequireIfAttributeIsOneOf(
	//						gitProviderTypeAttPath,
	//						[]attr.Value{gitProviderTypeGitHub},
	//					),
	//				},
	//				CustomType: supertypes.NewSingleNestedObjectTypeOf[gitCredentialsModel](ctx),
	//				Attributes: map[string]schema.Attribute{
	//					"source": schema.StringAttribute{
	//						MarkdownDescription: "The Git credentials source. Possible values: " + utils.ConvertStringSlicesToString(fabcore.PossibleGitCredentialsSourceValues(), true, true),
	//						Computed:            true,
	//					},
	//					"connection_id": schema.StringAttribute{
	//						MarkdownDescription: "The object ID of the connection.",
	//						Computed:            true,
	//						Optional:            true,
	//						CustomType:          customtypes.UUIDType{},
	//						Validators: []validator.String{
	//							superstringvalidator.NullIfAttributeIsOneOf(
	//								gitProviderTypeAttPath,
	//								[]attr.Value{gitProviderTypeAzureDevOps},
	//							),
	//							superstringvalidator.RequireIfAttributeIsOneOf(
	//								gitProviderTypeAttPath,
	//								[]attr.Value{gitProviderTypeGitHub},
	//							),
	//						},
	//					},
	//				},
	//			},
	//			"git_connection_state": schema.StringAttribute{
	//				MarkdownDescription: "The git connection state. Possible values: " + utils.ConvertStringSlicesToString(fabcore.PossibleGitConnectionStateValues(), true, true),
	//				Computed:            true,
	//			},
	//			"timeouts": timeouts.AttributesAll(ctx),
	//		},
	//	}
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

	if resp.Diagnostics.Append(fabricitem.IsPreviewMode(r.TypeInfo.Name, r.TypeInfo.IsPreview, r.pConfigData.Preview)...); resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceWorkspaceGit) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "CREATE", map[string]any{
		"action": "start",
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

	plan.ID = plan.WorkspaceID

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

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)

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

	var plan resourceWorkspaceGitModel

	if resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...); resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := plan.Timeouts.Update(ctx, r.pConfigData.Timeout)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var reqUpdate requestUpdateGitCredentials

	if resp.Diagnostics.Append(reqUpdate.set(ctx, plan)...); resp.Diagnostics.HasError() {
		return
	}

	respUpdate, err := r.client.UpdateMyGitCredentials(ctx, plan.WorkspaceID.ValueString(), reqUpdate, nil)
	if resp.Diagnostics.Append(utils.GetDiagsFromError(ctx, err, utils.OperationUpdate, nil)...); resp.Diagnostics.HasError() {
		return
	}

	if resp.Diagnostics.Append(plan.setCredentials(ctx, respUpdate.GitCredentialsConfigurationResponseClassification)...); resp.Diagnostics.HasError() {
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

func (r *resourceWorkspaceGit) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "DELETE", map[string]any{
		"action": "start",
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
	respGet, err := r.client.GetConnection(ctx, model.WorkspaceID.ValueString(), nil)
	if diags := utils.GetDiagsFromError(ctx, err, utils.OperationRead, fabcore.ErrGit.GitProviderResourceNotFound); diags.HasError() {
		return diags
	}

	if diags := model.set(ctx, respGet.GitConnection); diags.HasError() {
		return diags
	}

	respGetCredentials, err := r.client.GetMyGitCredentials(ctx, model.WorkspaceID.ValueString(), nil)
	if diags := utils.GetDiagsFromError(ctx, err, utils.OperationRead, nil); diags.HasError() {
		return diags
	}

	if diags := model.setCredentials(ctx, respGetCredentials.GitCredentialsConfigurationResponseClassification); diags.HasError() {
		return diags
	}

	return nil
}
