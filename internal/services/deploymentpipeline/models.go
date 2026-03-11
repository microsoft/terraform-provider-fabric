// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package deploymentpipeline

import (
	"context"

	timeoutsD "github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts" //revive:disable-line:import-alias-naming
	timeoutsR "github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"   //revive:disable-line:import-alias-naming
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
)

/*
BASE MODEL
*/

type baseDeploymentPipelineModel struct {
	ID          customtypes.UUID `tfsdk:"id"`
	DisplayName types.String     `tfsdk:"display_name"`
	Description types.String     `tfsdk:"description"`
}

type baseDeploymentPipelineStageModel struct {
	ID          customtypes.UUID `tfsdk:"id"`
	DisplayName types.String     `tfsdk:"display_name"`
	Description types.String     `tfsdk:"description"`
	WorkspaceID types.String     `tfsdk:"workspace_id"`
	IsPublic    types.Bool       `tfsdk:"is_public"`
}

type baseDeploymentPipelineExtendedInfoModel struct {
	baseDeploymentPipelineModel

	Stages supertypes.ListNestedObjectValueOf[baseDeploymentPipelineStageModel] `tfsdk:"stages"`
}

func (to *baseDeploymentPipelineModel) set(from fabcore.DeploymentPipeline) {
	to.ID = customtypes.NewUUIDPointerValue(from.ID)
	to.DisplayName = types.StringPointerValue(from.DisplayName)
	to.Description = types.StringPointerValue(from.Description)
}

func (to *baseDeploymentPipelineExtendedInfoModel) setStages(ctx context.Context, from []*baseDeploymentPipelineStageModel) diag.Diagnostics {
	slice := make([]*baseDeploymentPipelineStageModel, 0, len(from))

	for _, entity := range from {
		var entityModel baseDeploymentPipelineStageModel

		entityModel.ID = entity.ID
		entityModel.WorkspaceID = entity.WorkspaceID
		entityModel.DisplayName = entity.DisplayName
		entityModel.Description = entity.Description
		entityModel.IsPublic = entity.IsPublic

		slice = append(slice, &entityModel)
	}

	if diags := to.Stages.Set(ctx, slice); diags.HasError() {
		return diags
	}

	return nil
}

func (to *baseDeploymentPipelineExtendedInfoModel) set(ctx context.Context, from fabcore.DeploymentPipelineExtendedInfo) diag.Diagnostics {
	to.ID = customtypes.NewUUIDPointerValue(from.ID)
	to.DisplayName = types.StringPointerValue(from.DisplayName)
	to.Description = types.StringPointerValue(from.Description)
	slice := make([]*baseDeploymentPipelineStageModel, 0, len(from.Stages))

	for _, entity := range from.Stages {
		entityModel := &baseDeploymentPipelineStageModel{
			ID:          customtypes.NewUUIDPointerValue(entity.ID),
			WorkspaceID: types.StringPointerValue(entity.WorkspaceID),
			DisplayName: types.StringPointerValue(entity.DisplayName),
			Description: types.StringPointerValue(entity.Description),
			IsPublic:    types.BoolPointerValue(entity.IsPublic),
		}
		slice = append(slice, entityModel)
	}

	if diags := to.Stages.Set(ctx, slice); diags.HasError() {
		return diags
	}

	return nil
}

/*
DATA-SOURCE
*/

type dataSourceDeploymentPipelineModel struct {
	baseDeploymentPipelineExtendedInfoModel

	Timeouts timeoutsD.Value `tfsdk:"timeouts"`
}

/*
DATA-SOURCE (list)
*/

type dataSourceDeploymentPipelinesModel struct {
	Values   supertypes.SetNestedObjectValueOf[baseDeploymentPipelineModel] `tfsdk:"values"`
	Timeouts timeoutsD.Value                                                `tfsdk:"timeouts"`
}

func (to *dataSourceDeploymentPipelinesModel) setValues(ctx context.Context, from []fabcore.DeploymentPipeline) diag.Diagnostics {
	slice := make([]*baseDeploymentPipelineModel, 0, len(from))

	for _, entity := range from {
		var entityModel baseDeploymentPipelineModel

		entityModel.set(entity)
		slice = append(slice, &entityModel)
	}

	return to.Values.Set(ctx, slice)
}

/*
RESOURCE
*/

type resourceDeploymentPipelineModel struct {
	baseDeploymentPipelineExtendedInfoModel

	Timeouts timeoutsR.Value `tfsdk:"timeouts"`
}

type requestCreateDeploymentPipeline struct {
	fabcore.CreateDeploymentPipelineRequest
}

type requestAssignStageToWorkspace struct {
	fabcore.DeploymentPipelineAssignWorkspaceRequest
}

func (to *requestAssignStageToWorkspace) set(from baseDeploymentPipelineStageModel) {
	to.WorkspaceID = from.WorkspaceID.ValueStringPointer()
	to.WorkspaceID = from.WorkspaceID.ValueStringPointer()
}

func (to *requestCreateDeploymentPipeline) set(ctx context.Context, from resourceDeploymentPipelineModel) diag.Diagnostics {
	to.DisplayName = from.DisplayName.ValueStringPointer()
	to.Description = from.Description.ValueStringPointer()

	entities, diags := from.Stages.Get(ctx)
	if diags.HasError() {
		return diags
	}

	to.Stages = make([]fabcore.DeploymentPipelineStageRequest, len(entities))
	for i, entity := range entities {
		to.Stages[i].DisplayName = entity.DisplayName.ValueStringPointer()
		to.Stages[i].Description = entity.Description.ValueStringPointer()
		to.Stages[i].IsPublic = entity.IsPublic.ValueBoolPointer()
	}

	return nil
}

type requestUpdateDeploymentPipeline struct {
	fabcore.UpdateDeploymentPipelineRequest
}

func (to *requestUpdateDeploymentPipeline) set(from resourceDeploymentPipelineModel) {
	to.DisplayName = from.DisplayName.ValueStringPointer()
	to.Description = from.Description.ValueStringPointer()
}

type requestUpdateDeploymentPipelineStage struct {
	fabcore.DeploymentPipelineStageRequest
}

func (to *requestUpdateDeploymentPipelineStage) set(from baseDeploymentPipelineStageModel) {
	to.DisplayName = from.DisplayName.ValueStringPointer()
	to.Description = from.Description.ValueStringPointer()
	to.IsPublic = from.IsPublic.ValueBoolPointer()
}

func (to *baseDeploymentPipelineStageModel) set(from fabcore.DeploymentPipelineStage) {
	to.ID = customtypes.NewUUIDPointerValue(from.ID)
	to.DisplayName = types.StringPointerValue(from.DisplayName)
	to.Description = types.StringPointerValue(from.Description)
	to.IsPublic = types.BoolPointerValue(from.IsPublic)
}
