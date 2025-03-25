// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package workspace

import (
	"context"
	"fmt"

	timeoutsD "github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts" //revive:disable-line:import-alias-naming
	timeoutsR "github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"   //revive:disable-line:import-alias-naming
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
)

/*
BASE MODEL
*/

type baseWorkspaceModel struct {
	ID          customtypes.UUID `tfsdk:"id"`
	DisplayName types.String     `tfsdk:"display_name"`
	Description types.String     `tfsdk:"description"`
	Type        types.String     `tfsdk:"type"`
	CapacityID  customtypes.UUID `tfsdk:"capacity_id"`
}

func (to *baseWorkspaceModel) set(from fabcore.Workspace) {
	to.ID = customtypes.NewUUIDPointerValue(from.ID)
	to.DisplayName = types.StringPointerValue(from.DisplayName)
	to.Description = types.StringPointerValue(from.Description)
	to.Type = types.StringPointerValue((*string)(from.Type))
	to.CapacityID = customtypes.NewUUIDPointerValue(from.CapacityID)
}

type baseWorkspaceInfoModel struct {
	baseWorkspaceModel
	CapacityAssignmentProgress types.String                                                 `tfsdk:"capacity_assignment_progress"`
	CapacityRegion             types.String                                                 `tfsdk:"capacity_region"`
	OneLakeEndpoints           supertypes.SingleNestedObjectValueOf[oneLakeEndpointsModel]  `tfsdk:"onelake_endpoints"`
	Identity                   supertypes.SingleNestedObjectValueOf[workspaceIdentityModel] `tfsdk:"identity"`
}

func (to *baseWorkspaceInfoModel) set(ctx context.Context, from fabcore.WorkspaceInfo) diag.Diagnostics {
	to.ID = customtypes.NewUUIDPointerValue(from.ID)
	to.DisplayName = types.StringPointerValue(from.DisplayName)
	to.Description = types.StringPointerValue(from.Description)
	to.Type = types.StringPointerValue((*string)(from.Type))
	to.CapacityID = customtypes.NewUUIDPointerValue(from.CapacityID)
	to.CapacityAssignmentProgress = types.StringPointerValue((*string)(from.CapacityAssignmentProgress))
	to.CapacityRegion = types.StringPointerValue((*string)(from.CapacityRegion))

	oneLakeEndpoints := supertypes.NewSingleNestedObjectValueOfNull[oneLakeEndpointsModel](ctx)

	if from.OneLakeEndpoints != nil {
		oneLakeEndpointsModel := &oneLakeEndpointsModel{}
		oneLakeEndpointsModel.set(*from.OneLakeEndpoints)

		if diags := oneLakeEndpoints.Set(ctx, oneLakeEndpointsModel); diags.HasError() {
			return diags
		}
	}

	to.OneLakeEndpoints = oneLakeEndpoints

	workspaceIdentity := supertypes.NewSingleNestedObjectValueOfNull[workspaceIdentityModel](ctx)

	if from.WorkspaceIdentity != nil {
		workspaceIdentityModel := &workspaceIdentityModel{}
		workspaceIdentityModel.set(*from.WorkspaceIdentity)

		if diags := workspaceIdentity.Set(ctx, workspaceIdentityModel); diags.HasError() {
			return diags
		}
	}

	to.Identity = workspaceIdentity

	return nil
}

/*
DATA-SOURCE
*/

type dataSourceWorkspaceModel struct {
	baseWorkspaceInfoModel
	Timeouts timeoutsD.Value `tfsdk:"timeouts"`
}

/*
DATA-SOURCE (list)
*/

type dataSourceWorkspacesModel struct {
	Values   supertypes.ListNestedObjectValueOf[baseWorkspaceModel] `tfsdk:"values"`
	Timeouts timeoutsD.Value                                        `tfsdk:"timeouts"`
}

func (to *dataSourceWorkspacesModel) setValues(ctx context.Context, from []fabcore.Workspace) diag.Diagnostics {
	slice := make([]*baseWorkspaceModel, 0, len(from))

	for _, entity := range from {
		var entityModel baseWorkspaceModel
		entityModel.set(entity)
		slice = append(slice, &entityModel)
	}

	return to.Values.Set(ctx, slice)
}

/*
RESOURCE
*/

type resourceWorkspaceModel struct {
	baseWorkspaceInfoModel
	Timeouts timeoutsR.Value `tfsdk:"timeouts"`
}

type requestCreateWorkspace struct {
	fabcore.CreateWorkspaceRequest
}

func (to *requestCreateWorkspace) set(from resourceWorkspaceModel) {
	to.DisplayName = from.DisplayName.ValueStringPointer()
	to.Description = from.Description.ValueStringPointer()
	to.CapacityID = from.CapacityID.ValueStringPointer()
}

type requestUpdateWorkspace struct {
	fabcore.UpdateWorkspaceRequest
}

func (to *requestUpdateWorkspace) set(from resourceWorkspaceModel) {
	to.DisplayName = from.DisplayName.ValueStringPointer()
	to.Description = from.Description.ValueStringPointer()
}

type assignWorkspaceToCapacityRequest struct {
	fabcore.AssignWorkspaceToCapacityRequest
}

func (to *assignWorkspaceToCapacityRequest) set(from resourceWorkspaceModel) {
	to.CapacityID = from.CapacityID.ValueStringPointer()
}

/*
HELPER MODELS
*/

type workspaceIdentityModel struct {
	Type               types.String     `tfsdk:"type"`
	ApplicationID      customtypes.UUID `tfsdk:"application_id"`
	ServicePrincipalID customtypes.UUID `tfsdk:"service_principal_id"`
}

func (to *workspaceIdentityModel) set(from fabcore.WorkspaceIdentity) {
	to.Type = types.StringValue(workspaceIdentityTypes[0])
	to.ApplicationID = customtypes.NewUUIDPointerValue(from.ApplicationID)
	to.ServicePrincipalID = customtypes.NewUUIDPointerValue(from.ServicePrincipalID)
}

type oneLakeEndpointsModel struct {
	BlobEndpoint customtypes.URL `tfsdk:"blob_endpoint"`
	DfsEndpoint  customtypes.URL `tfsdk:"dfs_endpoint"`
}

func (to *oneLakeEndpointsModel) set(from fabcore.OneLakeEndpoints) {
	to.BlobEndpoint = customtypes.NewURLPointerValue(from.BlobEndpoint)
	to.DfsEndpoint = customtypes.NewURLPointerValue(from.DfsEndpoint)
}

func checkWorkspaceType(entity fabcore.WorkspaceInfo) diag.Diagnostics {
	var diags diag.Diagnostics

	switch *entity.Type {
	case fabcore.WorkspaceTypePersonal:
		diags.AddError(
			common.ErrorWorkspaceNotSupportedHeader,
			fmt.Sprintf(common.ErrorWorkspaceNotSupportedDetails, string(fabcore.WorkspaceTypePersonal)),
		)

		return diags
	case fabcore.WorkspaceTypeAdminWorkspace:
		diags.AddError(
			common.ErrorWorkspaceNotSupportedHeader,
			fmt.Sprintf(common.ErrorWorkspaceNotSupportedDetails, string(fabcore.WorkspaceTypeAdminWorkspace)),
		)

		return diags
	default:
		return nil
	}
}
