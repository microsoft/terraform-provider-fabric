// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package workspacegit

import (
	"context"
	"fmt"

	azto "github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	timeoutsD "github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts" //revive:disable-line:import-alias-naming
	timeoutsR "github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"   //revive:disable-line:import-alias-naming
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
)

/*
BASE MODEL
*/

type baseWorkspaceGitModel struct {
	ID                 customtypes.UUID                                              `tfsdk:"id"`
	WorkspaceID        customtypes.UUID                                              `tfsdk:"workspace_id"`
	GitConnectionState types.String                                                  `tfsdk:"git_connection_state"`
	GitSyncDetails     supertypes.SingleNestedObjectValueOf[gitSyncDetailsModel]     `tfsdk:"git_sync_details"`
	GitProviderDetails supertypes.SingleNestedObjectValueOf[gitProviderDetailsModel] `tfsdk:"git_provider_details"`
	GitCredentials     supertypes.SingleNestedObjectValueOf[gitCredentialsModel]     `tfsdk:"git_credentials"`
}

func (to *baseWorkspaceGitModel) set(ctx context.Context, from fabcore.GitConnection) diag.Diagnostics {
	to.GitConnectionState = types.StringPointerValue((*string)(from.GitConnectionState))

	syncDetails := supertypes.NewSingleNestedObjectValueOfNull[gitSyncDetailsModel](ctx)

	if from.GitSyncDetails != nil {
		syncDetailsModel := &gitSyncDetailsModel{}
		syncDetailsModel.set(*from.GitSyncDetails)

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

func (to *baseWorkspaceGitModel) setCredentials(ctx context.Context, from fabcore.GitCredentialsConfigurationResponseClassification) diag.Diagnostics {
	gitCredentials := supertypes.NewSingleNestedObjectValueOfNull[gitCredentialsModel](ctx)

	if from != nil {
		credentialsModel := &gitCredentialsModel{}
		if diags := credentialsModel.set(from); diags.HasError() {
			return diags
		}

		if diags := gitCredentials.Set(ctx, credentialsModel); diags.HasError() {
			return diags
		}
	}

	to.GitCredentials = gitCredentials

	return nil
}

/*
DATA-SOURCE
*/

type dataSourceWorkspaceGitModel struct {
	baseWorkspaceGitModel
	Timeouts timeoutsD.Value `tfsdk:"timeouts"`
}

/*
RESOURCE
*/

type resourceWorkspaceGitModel struct {
	InitializationStrategy types.String `tfsdk:"initialization_strategy"`
	baseWorkspaceGitModel
	Timeouts timeoutsR.Value `tfsdk:"timeouts"`
}

type requestGitConnect struct {
	fabcore.GitConnectRequest
}

func (to *requestGitConnect) set(ctx context.Context, from resourceWorkspaceGitModel) diag.Diagnostics {
	gitProviderDetails, diags := from.GitProviderDetails.Get(ctx)
	if diags.HasError() {
		return diags
	}

	var reqGitProviderDetails fabcore.GitProviderDetailsClassification
	var reqGitCredentials fabcore.GitCredentialsClassification

	gitProviderType := (fabcore.GitProviderType)(gitProviderDetails.GitProviderType.ValueString())

	switch gitProviderType {
	case fabcore.GitProviderTypeAzureDevOps:
		reqGitProviderDetails = &fabcore.AzureDevOpsDetails{
			GitProviderType:  &gitProviderType,
			OrganizationName: gitProviderDetails.OrganizationName.ValueStringPointer(),
			ProjectName:      gitProviderDetails.ProjectName.ValueStringPointer(),
			RepositoryName:   gitProviderDetails.RepositoryName.ValueStringPointer(),
			BranchName:       gitProviderDetails.BranchName.ValueStringPointer(),
			DirectoryName:    gitProviderDetails.DirectoryName.ValueStringPointer(),
		}

		reqGitCredentials = &fabcore.AutomaticGitCredentials{
			Source: azto.Ptr(fabcore.GitCredentialsSourceAutomatic),
		}

	case fabcore.GitProviderTypeGitHub:
		reqGitProviderDetails = &fabcore.GitHubDetails{
			GitProviderType: &gitProviderType,
			OwnerName:       gitProviderDetails.OwnerName.ValueStringPointer(),
			RepositoryName:  gitProviderDetails.RepositoryName.ValueStringPointer(),
			BranchName:      gitProviderDetails.BranchName.ValueStringPointer(),
			DirectoryName:   gitProviderDetails.DirectoryName.ValueStringPointer(),
		}

		gitGitCredentials, diags := from.GitCredentials.Get(ctx)
		if diags.HasError() {
			return diags
		}

		reqGitCredentials = &fabcore.ConfiguredConnectionGitCredentials{
			Source:       azto.Ptr(fabcore.GitCredentialsSourceConfiguredConnection),
			ConnectionID: gitGitCredentials.ConnectionID.ValueStringPointer(),
		}

	default:
		diags.AddError(
			"Unsupported Git provider type",
			fmt.Sprintf("The Git provider type '%T' is not supported.", gitProviderType),
		)

		return diags
	}

	to.GitProviderDetails = reqGitProviderDetails
	to.MyGitCredentials = reqGitCredentials

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

type requestUpdateGitCredentials struct {
	fabcore.UpdateGitCredentialsRequestClassification
}

func (to *requestUpdateGitCredentials) set(ctx context.Context, from resourceWorkspaceGitModel) diag.Diagnostics {
	gitCredentials, diags := from.GitCredentials.Get(ctx)
	if diags.HasError() {
		return diags
	}

	var reqUpdate fabcore.UpdateGitCredentialsRequestClassification

	gitCredentialsSource := fabcore.GitCredentialsSource(gitCredentials.Source.ValueString())

	switch gitCredentialsSource {
	case fabcore.GitCredentialsSourceAutomatic:
		reqUpdate = &fabcore.UpdateGitCredentialsToAutomaticRequest{
			Source: azto.Ptr(fabcore.GitCredentialsSourceAutomatic),
		}
	case fabcore.GitCredentialsSourceConfiguredConnection:
		reqUpdate = &fabcore.UpdateGitCredentialsToConfiguredConnectionRequest{
			Source:       azto.Ptr(fabcore.GitCredentialsSourceConfiguredConnection),
			ConnectionID: gitCredentials.ConnectionID.ValueStringPointer(),
		}
	case fabcore.GitCredentialsSourceNone:
		reqUpdate = &fabcore.UpdateGitCredentialsToNoneRequest{
			Source: azto.Ptr(fabcore.GitCredentialsSourceNone),
		}
	}

	to.UpdateGitCredentialsRequestClassification = reqUpdate

	return nil
}

/*
HELPER MODELS
*/

type gitSyncDetailsModel struct {
	Head         types.String      `tfsdk:"head"`
	LastSyncTime timetypes.RFC3339 `tfsdk:"last_sync_time"`
}

func (to *gitSyncDetailsModel) set(from fabcore.GitSyncDetails) {
	to.Head = types.StringPointerValue(from.Head)
	to.LastSyncTime = timetypes.NewRFC3339TimePointerValue(from.LastSyncTime)
}

type gitCredentialsModel struct {
	Source       types.String     `tfsdk:"source"`
	ConnectionID customtypes.UUID `tfsdk:"connection_id"`
}

func (to *gitCredentialsModel) set(from fabcore.GitCredentialsConfigurationResponseClassification) diag.Diagnostics {
	switch gitCredentials := from.(type) {
	case *fabcore.AutomaticGitCredentialsResponse:
		to.Source = types.StringPointerValue((*string)(gitCredentials.Source))
		to.ConnectionID = customtypes.NewUUIDNull()

		return nil

	case *fabcore.ConfiguredConnectionGitCredentialsResponse:
		to.Source = types.StringPointerValue((*string)(gitCredentials.Source))
		to.ConnectionID = customtypes.NewUUIDPointerValue(gitCredentials.ConnectionID)

		return nil

	case *fabcore.NoneGitCredentialsResponse:
		to.Source = types.StringPointerValue((*string)(gitCredentials.Source))
		to.ConnectionID = customtypes.NewUUIDNull()

		return nil

	default:
		var diags diag.Diagnostics

		diags.AddError(
			"Unsupported Git credentials type",
			fmt.Sprintf("The Git credentials type '%T' is not supported.", gitCredentials),
		)

		return diags
	}
}

type gitProviderDetailsModel struct {
	GitProviderType  types.String `tfsdk:"git_provider_type"`
	OrganizationName types.String `tfsdk:"organization_name"`
	ProjectName      types.String `tfsdk:"project_name"`
	OwnerName        types.String `tfsdk:"owner_name"`
	RepositoryName   types.String `tfsdk:"repository_name"`
	BranchName       types.String `tfsdk:"branch_name"`
	DirectoryName    types.String `tfsdk:"directory_name"`
}

func (to *gitProviderDetailsModel) set(from fabcore.GitProviderDetailsClassification) diag.Diagnostics {
	switch gitProviderDetails := from.(type) {
	case *fabcore.AzureDevOpsDetails:
		// Azure DevOps does not have an owner name
		to.OwnerName = types.StringNull()

		to.GitProviderType = types.StringPointerValue((*string)(gitProviderDetails.GitProviderType))
		to.OrganizationName = types.StringPointerValue(gitProviderDetails.OrganizationName)
		to.ProjectName = types.StringPointerValue(gitProviderDetails.ProjectName)
		to.RepositoryName = types.StringPointerValue(gitProviderDetails.RepositoryName)
		to.BranchName = types.StringPointerValue(gitProviderDetails.BranchName)
		to.DirectoryName = types.StringPointerValue(gitProviderDetails.DirectoryName)

		return nil

	case *fabcore.GitHubDetails:
		// GitHub does not have an organization name or project name
		to.OrganizationName = types.StringNull()
		to.ProjectName = types.StringNull()

		to.GitProviderType = types.StringPointerValue((*string)(gitProviderDetails.GitProviderType))
		to.OwnerName = types.StringPointerValue(gitProviderDetails.OwnerName)
		to.RepositoryName = types.StringPointerValue(gitProviderDetails.RepositoryName)
		to.BranchName = types.StringPointerValue(gitProviderDetails.BranchName)
		to.DirectoryName = types.StringPointerValue(gitProviderDetails.DirectoryName)

		return nil

	default:
		var diags diag.Diagnostics

		diags.AddError(
			"Unsupported Git provider type",
			fmt.Sprintf("The Git provider type '%T' is not supported.", gitProviderDetails),
		)

		return diags
	}
}
