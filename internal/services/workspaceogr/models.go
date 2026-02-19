// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package workspaceogr

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

type baseWorkspaceOutboundGatewayRulesModel struct {
	ID              customtypes.UUID                                                   `tfsdk:"id"`
	WorkspaceID     customtypes.UUID                                                   `tfsdk:"workspace_id"`
	AllowedGateways supertypes.ListNestedObjectValueOf[gatewayAccessRuleMetadataModel] `tfsdk:"allowed_gateways"`
	DefaultAction   types.String                                                       `tfsdk:"default_action"`
}

type gatewayAccessRuleMetadataModel struct {
	ID customtypes.UUID `tfsdk:"id"`
}

func (to *baseWorkspaceOutboundGatewayRulesModel) set(ctx context.Context, workspaceID string, from fabcore.WorkspaceOutboundGateways) diag.Diagnostics {
	to.WorkspaceID = customtypes.NewUUIDValue(workspaceID)
	to.DefaultAction = types.StringPointerValue((*string)(from.DefaultAction))

	slice := make([]*gatewayAccessRuleMetadataModel, 0, len(from.AllowedGateways))

	for _, prop := range from.AllowedGateways {
		gatewayRuleMetadata := &gatewayAccessRuleMetadataModel{}
		gatewayRuleMetadata.set(prop)
		slice = append(slice, gatewayRuleMetadata)
	}

	if diags := to.AllowedGateways.Set(ctx, slice); diags.HasError() {
		return diags
	}

	return nil
}

func (to *gatewayAccessRuleMetadataModel) set(from fabcore.GatewayAccessRuleMetadata) {
	to.ID = customtypes.NewUUIDPointerValue(from.ID)
}

/*
DATA-SOURCE
*/

type dataSourceWorkspaceOutboundGatewayRulesModel struct {
	baseWorkspaceOutboundGatewayRulesModel

	Timeouts timeoutsD.Value `tfsdk:"timeouts"`
}

/*
RESOURCE
*/

type resourceWorkspaceOutboundGatewayRulesModel struct {
	baseWorkspaceOutboundGatewayRulesModel

	Timeouts timeoutsR.Value `tfsdk:"timeouts"`
}

type requestSetWorkspaceOutboundGatewayRules struct {
	fabcore.WorkspaceOutboundGateways
}

func (to *requestSetWorkspaceOutboundGatewayRules) set(ctx context.Context, from resourceWorkspaceOutboundGatewayRulesModel) diag.Diagnostics {
	to.DefaultAction = (*fabcore.GatewayAccessActionType)(from.DefaultAction.ValueStringPointer())

	allowedGateways, diags := from.AllowedGateways.Get(ctx)
	if diags.HasError() {
		return diags
	}

	allowedGatewaysSlice := make([]fabcore.GatewayAccessRuleMetadata, 0, len(allowedGateways))

	for _, prop := range allowedGateways {
		allowedGatewaysSlice = append(allowedGatewaysSlice, fabcore.GatewayAccessRuleMetadata{
			ID: prop.ID.ValueStringPointer(),
		})
	}

	to.AllowedGateways = allowedGatewaysSlice

	return nil
}
