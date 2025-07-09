// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package deploymentpipelinera

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

type baseDeploymentPipelineRoleAssignmentModel struct {
	ID                   customtypes.UUID                                     `tfsdk:"id"`
	DeploymentPipelineID customtypes.UUID                                     `tfsdk:"deployment_pipeline_id"`
	Role                 types.String                                         `tfsdk:"role"`
	Principal            supertypes.SingleNestedObjectValueOf[principalModel] `tfsdk:"principal"`
}

func (to *baseDeploymentPipelineRoleAssignmentModel) set(ctx context.Context, deploymentPipelineID string, from fabcore.DeploymentPipelineRoleAssignment) diag.Diagnostics {
	to.ID = customtypes.NewUUIDPointerValue(from.ID)
	to.DeploymentPipelineID = customtypes.NewUUIDValue(deploymentPipelineID)
	to.Role = types.StringPointerValue((*string)(from.Role))

	if from.Principal != nil {
		principalModel := &principalModel{}

		principalModel.set(*from.Principal)

		if diags := to.Principal.Set(ctx, principalModel); diags.HasError() {
			return diags
		}
	}

	return nil
}

/*
DATA-SOURCE (list)
*/

type dataSourceDeploymentPipelineRoleAssignmentsModel struct {
	DeploymentPipelineID customtypes.UUID                                                             `tfsdk:"deployment_pipeline_id"`
	Values               supertypes.SetNestedObjectValueOf[baseDeploymentPipelineRoleAssignmentModel] `tfsdk:"values"`
	Timeouts             timeoutsD.Value                                                              `tfsdk:"timeouts"`
}

func (to *dataSourceDeploymentPipelineRoleAssignmentsModel) setValues(ctx context.Context, deploymentPipelineID string, from []fabcore.DeploymentPipelineRoleAssignment) diag.Diagnostics {
	slice := make([]*baseDeploymentPipelineRoleAssignmentModel, 0, len(from))

	for _, entity := range from {
		var entityModel baseDeploymentPipelineRoleAssignmentModel

		if diags := entityModel.set(ctx, deploymentPipelineID, entity); diags.HasError() {
			return diags
		}

		slice = append(slice, &entityModel)
	}

	return to.Values.Set(ctx, slice)
}

/*
RESOURCE
*/

type resourceDeploymentPipelineRoleAssignmentModel struct {
	baseDeploymentPipelineRoleAssignmentModel

	Timeouts timeoutsR.Value `tfsdk:"timeouts"`
}

type requestCreateDeploymentPipelineRoleAssignment struct {
	fabcore.AddDeploymentPipelineRoleAssignmentRequest
}

func (to *requestCreateDeploymentPipelineRoleAssignment) set(ctx context.Context, from resourceDeploymentPipelineRoleAssignmentModel) diag.Diagnostics {
	principal, diags := from.Principal.Get(ctx)
	if diags.HasError() {
		return diags
	}

	to.Principal = &fabcore.Principal{
		ID:   principal.ID.ValueStringPointer(),
		Type: (*fabcore.PrincipalType)(principal.Type.ValueStringPointer()),
	}
	to.Role = (*fabcore.DeploymentPipelineRole)(from.Role.ValueStringPointer())

	return nil
}

/*
HELPER MODELS
*/

type principalModel struct {
	ID   customtypes.UUID `tfsdk:"id"`
	Type types.String     `tfsdk:"type"`
}

func (to *principalModel) set(from fabcore.Principal) {
	to.ID = customtypes.NewUUIDPointerValue(from.ID)
	to.Type = types.StringPointerValue((*string)(from.Type))
}
