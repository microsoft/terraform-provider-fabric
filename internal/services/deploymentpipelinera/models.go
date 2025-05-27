// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package deploymentpipelinera

import (
	"context"

	timeoutsD "github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts" //revive:disable-line:import-alias-naming
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
	// ID                   customtypes.UUID                                     `tfsdk:"id"`
	DeploymentPipelineID customtypes.UUID                                     `tfsdk:"deployment_pipeline_id"`
	Role                 types.String                                         `tfsdk:"role"`
	Principal            supertypes.SingleNestedObjectValueOf[principalModel] `tfsdk:"principal"`
}

func (to *baseDeploymentPipelineRoleAssignmentModel) set(ctx context.Context, deploymentPipelineID string, from fabcore.DeploymentPipelineRoleAssignment) diag.Diagnostics {
	to.DeploymentPipelineID = customtypes.NewUUIDValue(deploymentPipelineID)
	to.Role = types.StringPointerValue((*string)(from.Role))

	principal := supertypes.NewSingleNestedObjectValueOfNull[principalModel](ctx)

	if from.Principal != nil {
		principalModel := &principalModel{}

		principalModel.set(*from.Principal)

		if diags := principal.Set(ctx, principalModel); diags.HasError() {
			return diags
		}
	}

	to.Principal = principal

	return nil
}

/*
DATA SOURCE
*/

type dataSourceDeploymentPipelineRoleAssignmentsModel struct {
	DeploymentPipelineID customtypes.UUID                                                             `tfsdk:"deployment_pipeline_id"`
	Values               supertypes.SetNestedObjectValueOf[baseDeploymentPipelineRoleAssignmentModel] `tfsdk:"values"`
	Timeouts             timeoutsD.Value                                                              `tfsdk:"timeouts"`
}

func (to *dataSourceDeploymentPipelineRoleAssignmentsModel) set(ctx context.Context, deploymentPipelineID string, from []fabcore.DeploymentPipelineRoleAssignment) diag.Diagnostics {
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
