// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package workspacempe

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

type baseWorkspaceManagedPrivateEndpointModel struct {
	WorkspaceID                 customtypes.UUID                                           `tfsdk:"workspace_id"`
	ID                          customtypes.UUID                                           `tfsdk:"id"`
	Name                        types.String                                               `tfsdk:"name"`
	ProvisioningState           types.String                                               `tfsdk:"provisioning_state"`
	TargetPrivateLinkResourceID customtypes.CaseInsensitiveString                          `tfsdk:"target_private_link_resource_id"`
	TargetSubresourceType       types.String                                               `tfsdk:"target_subresource_type"`
	ConnectionState             supertypes.SingleNestedObjectValueOf[connectionStateModel] `tfsdk:"connection_state"`
}

func (to *baseWorkspaceManagedPrivateEndpointModel) set(ctx context.Context, workspaceID string, from fabcore.ManagedPrivateEndpoint) diag.Diagnostics {
	to.WorkspaceID = customtypes.NewUUIDValue(workspaceID)

	to.ID = customtypes.NewUUIDPointerValue(from.ID)
	to.Name = types.StringPointerValue(from.Name)
	to.ProvisioningState = types.StringPointerValue((*string)(from.ProvisioningState))
	to.TargetPrivateLinkResourceID = customtypes.NewCaseInsensitiveStringPointerValue(from.TargetPrivateLinkResourceID)
	to.TargetSubresourceType = types.StringPointerValue((from.TargetSubresourceType))

	connectionState := supertypes.NewSingleNestedObjectValueOfNull[connectionStateModel](ctx)

	if from.ConnectionState != nil {
		connectionStateModel := &connectionStateModel{}
		connectionStateModel.set(*from.ConnectionState)

		if diags := connectionState.Set(ctx, connectionStateModel); diags.HasError() {
			return diags
		}
	}

	to.ConnectionState = connectionState

	return nil
}

type connectionStateModel struct {
	ActionsRequired types.String `tfsdk:"actions_required"`
	Status          types.String `tfsdk:"status"`
	Description     types.String `tfsdk:"description"`
}

func (to *connectionStateModel) set(from fabcore.PrivateEndpointConnectionState) {
	to.ActionsRequired = types.StringPointerValue(from.ActionsRequired)
	to.Status = types.StringPointerValue((*string)(from.Status))
	to.Description = types.StringPointerValue(from.Description)
}

/*
DATA-SOURCE
*/

type dataSourceWorkspaceManagedPrivateEndpointModel struct {
	baseWorkspaceManagedPrivateEndpointModel

	Timeouts timeoutsD.Value `tfsdk:"timeouts"`
}

/*
DATA-SOURCE (list)
*/

type dataSourceWorkspaceManagedPrivateEndpointsModel struct {
	WorkspaceID customtypes.UUID                                                            `tfsdk:"workspace_id"`
	Values      supertypes.SetNestedObjectValueOf[baseWorkspaceManagedPrivateEndpointModel] `tfsdk:"values"`
	Timeouts    timeoutsD.Value                                                             `tfsdk:"timeouts"`
}

func (to *dataSourceWorkspaceManagedPrivateEndpointsModel) setValues(ctx context.Context, workspaceID string, from []fabcore.ManagedPrivateEndpoint) diag.Diagnostics {
	slice := make([]*baseWorkspaceManagedPrivateEndpointModel, 0, len(from))

	for _, entity := range from {
		var entityModel baseWorkspaceManagedPrivateEndpointModel
		if diags := entityModel.set(ctx, workspaceID, entity); diags.HasError() {
			return diags
		}

		slice = append(slice, &entityModel)
	}

	return to.Values.Set(ctx, slice)
}

/*
RESOURCE
*/

type resourceWorkspaceManagedPrivateEndpointModel struct {
	baseWorkspaceManagedPrivateEndpointModel

	RequestMessage types.String    `tfsdk:"request_message"`
	Timeouts       timeoutsR.Value `tfsdk:"timeouts"`
}

type requestCreateWorkspaceManagedPrivateEndpoint struct {
	fabcore.CreateManagedPrivateEndpointRequest
}

func (to *requestCreateWorkspaceManagedPrivateEndpoint) set(from resourceWorkspaceManagedPrivateEndpointModel) {
	to.Name = from.Name.ValueStringPointer()
	to.TargetPrivateLinkResourceID = from.TargetPrivateLinkResourceID.ValueStringPointer()
	to.TargetSubresourceType = from.TargetSubresourceType.ValueStringPointer()
	to.RequestMessage = from.RequestMessage.ValueStringPointer()
}
