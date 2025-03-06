// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package workspace

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
)

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
