// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package workspace

import (
	"context"
	"fmt"

	azto "github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	timeoutsd "github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"
	timeoutsr "github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
)

type dataSourceWorkspaceGitModel struct {
	baseWorkspaceGitModel
	Timeouts timeoutsd.Value `tfsdk:"timeouts"`
}

type resourceWorkspaceGitModel struct {
	ID                     customtypes.UUID `tfsdk:"id"`
	InitializationStrategy types.String     `tfsdk:"initialization_strategy"`
	baseWorkspaceGitModel
	Timeouts timeoutsr.Value `tfsdk:"timeouts"`
}

type baseWorkspaceGitModel struct {
	WorkspaceID        customtypes.UUID                                              `tfsdk:"workspace_id"`
	GitConnectionState types.String                                                  `tfsdk:"git_connection_state"`
	GitSyncDetails     supertypes.SingleNestedObjectValueOf[gitSyncDetailsModel]     `tfsdk:"git_sync_details"`
	GitProviderDetails supertypes.SingleNestedObjectValueOf[gitProviderDetailsModel] `tfsdk:"git_provider_details"`
}

func (to *baseWorkspaceGitModel) set(ctx context.Context, from fabcore.GitConnection) diag.Diagnostics {
	to.GitConnectionState = types.StringPointerValue((*string)(from.GitConnectionState))

	syncDetails := supertypes.NewSingleNestedObjectValueOfNull[gitSyncDetailsModel](ctx)

	if from.GitSyncDetails != nil {
		syncDetailsModel := &gitSyncDetailsModel{}
		syncDetailsModel.set(from.GitSyncDetails)

		if diags := syncDetails.Set(ctx, syncDetailsModel); diags.HasError() {
			return diags
		}
	}

	to.GitSyncDetails = syncDetails

	providerDetails := supertypes.NewSingleNestedObjectValueOfNull[gitProviderDetailsModel](ctx)

	if from.GitProviderDetails != nil {
		providerDetailsModel := &gitProviderDetailsModel{}
		if diags := providerDetailsModel.set(from.GitProviderDetails); diags.HasError() {
			return diags
		}

		if diags := providerDetails.Set(ctx, providerDetailsModel); diags.HasError() {
			return diags
		}
	}

	to.GitProviderDetails = providerDetails

	return nil
}

type gitSyncDetailsModel struct {
	Head         types.String      `tfsdk:"head"`
	LastSyncTime timetypes.RFC3339 `tfsdk:"last_sync_time"`
}

func (to *gitSyncDetailsModel) set(from *fabcore.GitSyncDetails) {
	to.Head = types.StringPointerValue(from.Head)
	to.LastSyncTime = timetypes.NewRFC3339TimePointerValue(from.LastSyncTime)
}

type gitProviderDetailsModel struct {
	GitProviderType  types.String `tfsdk:"git_provider_type"`
	OrganizationName types.String `tfsdk:"organization_name"`
	ProjectName      types.String `tfsdk:"project_name"`
	RepositoryName   types.String `tfsdk:"repository_name"`
	BranchName       types.String `tfsdk:"branch_name"`
	DirectoryName    types.String `tfsdk:"directory_name"`
}

func (to *gitProviderDetailsModel) set(from fabcore.GitProviderDetailsClassification) diag.Diagnostics {
	var diags diag.Diagnostics

	switch gitProviderDetails := from.(type) {
	case *fabcore.AzureDevOpsDetails:
		to.GitProviderType = types.StringPointerValue((*string)(gitProviderDetails.GitProviderType))
		to.OrganizationName = types.StringPointerValue(gitProviderDetails.OrganizationName)
		to.ProjectName = types.StringPointerValue(gitProviderDetails.ProjectName)
		to.RepositoryName = types.StringPointerValue(gitProviderDetails.RepositoryName)
		to.BranchName = types.StringPointerValue(gitProviderDetails.BranchName)
		to.DirectoryName = types.StringPointerValue(gitProviderDetails.DirectoryName)

		return nil

	default:
		diags.AddError("Unsupported Git provider type", fmt.Sprintf("The Git provider type '%T' is not supported.", gitProviderDetails))

		return diags
	}
}

type requestGitConnect struct {
	fabcore.GitConnectRequest
}

func (to *requestGitConnect) set(ctx context.Context, from resourceWorkspaceGitModel) diag.Diagnostics {
	gitProviderDetails, diags := from.GitProviderDetails.Get(ctx)
	if diags.HasError() {
		return diags
	}

	gitProviderType := (fabcore.GitProviderType)(gitProviderDetails.GitProviderType.ValueString())

	switch gitProviderType {
	case fabcore.GitProviderTypeAzureDevOps:
		to.GitProviderDetails = &fabcore.AzureDevOpsDetails{
			GitProviderType:  &gitProviderType,
			OrganizationName: gitProviderDetails.OrganizationName.ValueStringPointer(),
			ProjectName:      gitProviderDetails.ProjectName.ValueStringPointer(),
			RepositoryName:   gitProviderDetails.RepositoryName.ValueStringPointer(),
			BranchName:       gitProviderDetails.BranchName.ValueStringPointer(),
			DirectoryName:    gitProviderDetails.DirectoryName.ValueStringPointer(),
		}
	default:
		diags.AddError(
			"Unsupported Git provider type",
			fmt.Sprintf("The Git provider type '%s' is not supported.", string(gitProviderType)),
		)

		return diags
	}

	return nil
}

type requestGitCommitTo struct {
	fabcore.CommitToGitRequest
}

func (to *requestGitCommitTo) set(workspaceHead *string) {
	to.Mode = azto.Ptr(fabcore.CommitModeAll)
	to.WorkspaceHead = workspaceHead
}

type requestGitUpdateFrom struct {
	fabcore.UpdateFromGitRequest
}

func (to *requestGitUpdateFrom) set(remoteCommitHash, conflictResolutionPolicy *string) {
	policy := fabcore.ConflictResolutionPolicyPreferWorkspace
	if *conflictResolutionPolicy != "None" {
		policy = fabcore.ConflictResolutionPolicy(*conflictResolutionPolicy)
	}

	to.RemoteCommitHash = remoteCommitHash
	to.Options = &fabcore.UpdateOptions{}
	to.ConflictResolution = &fabcore.WorkspaceConflictResolution{
		ConflictResolutionPolicy: azto.Ptr(policy),
		ConflictResolutionType:   azto.Ptr(fabcore.ConflictResolutionTypeWorkspace),
	}
}

type requestGitInitialize struct {
	fabcore.InitializeGitConnectionRequest
}

func (to *requestGitInitialize) set(from resourceWorkspaceGitModel) {
	to.InitializationStrategy = (*fabcore.InitializationStrategy)(from.InitializationStrategy.ValueStringPointer())
}
