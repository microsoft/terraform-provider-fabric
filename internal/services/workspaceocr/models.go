// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package workspaceocr

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

type baseWorkspaceOutboundCloudConnectionRulesModel struct {
	WorkspaceID   customtypes.UUID                               `tfsdk:"workspace_id"`
	Rules         supertypes.ListNestedObjectValueOf[rulesModel] `tfsdk:"rules"`
	DefaultAction types.String                                   `tfsdk:"default_action"`
}

type rulesModel struct {
	ConnectionType    customtypes.CaseInsensitiveString                  `tfsdk:"connection_type"`
	DefaultAction     types.String                                       `tfsdk:"default_action"`
	AllowedEndpoints  supertypes.ListNestedObjectValueOf[endpointModel]  `tfsdk:"allowed_endpoints"`
	AllowedWorkspaces supertypes.ListNestedObjectValueOf[workspaceModel] `tfsdk:"allowed_workspaces"`
}

type workspaceModel struct {
	WorkspaceID customtypes.UUID `tfsdk:"workspace_id"`
}

type endpointModel struct {
	HostnamePattern types.String `tfsdk:"hostname_pattern"`
}

func (to *baseWorkspaceOutboundCloudConnectionRulesModel) set(ctx context.Context, workspaceID string, from fabcore.WorkspaceOutboundConnections) diag.Diagnostics {
	to.WorkspaceID = customtypes.NewUUIDValue(workspaceID)
	to.DefaultAction = types.StringPointerValue((*string)(from.DefaultAction))

	slice := make([]*rulesModel, 0, len(from.Rules))

	for _, prop := range from.Rules {
		rule := &rulesModel{}

		if diags := rule.set(ctx, prop); diags.HasError() {
			return diags
		}

		slice = append(slice, rule)
	}

	if diags := to.Rules.Set(ctx, slice); diags.HasError() {
		return diags
	}

	return nil
}

func (to *rulesModel) set(ctx context.Context, from fabcore.OutboundConnectionRule) diag.Diagnostics {
	to.ConnectionType = customtypes.NewCaseInsensitiveStringPointerValue(from.ConnectionType)
	to.DefaultAction = types.StringPointerValue((*string)(from.DefaultAction))

	allowedEndpoints := make([]*endpointModel, 0, len(from.AllowedEndpoints))
	for _, endpoint := range from.AllowedEndpoints {
		endpointM := &endpointModel{}
		endpointM.set(endpoint)
		allowedEndpoints = append(allowedEndpoints, endpointM)
	}

	if diags := to.AllowedEndpoints.Set(ctx, allowedEndpoints); diags.HasError() {
		return diags
	}

	allowedWorkspaces := make([]*workspaceModel, 0, len(from.AllowedWorkspaces))
	for _, workspace := range from.AllowedWorkspaces {
		workspaceM := &workspaceModel{}
		workspaceM.set(workspace)
		allowedWorkspaces = append(allowedWorkspaces, workspaceM)
	}

	if diags := to.AllowedWorkspaces.Set(ctx, allowedWorkspaces); diags.HasError() {
		return diags
	}

	return nil
}

func (to *endpointModel) set(from fabcore.ConnectionRuleEndpointMetadata) {
	to.HostnamePattern = types.StringPointerValue(from.HostnamePattern)
}

func (to *workspaceModel) set(from fabcore.ConnectionRuleWorkspaceMetadata) {
	to.WorkspaceID = customtypes.NewUUIDPointerValue(from.WorkspaceID)
}

/*
DATA-SOURCE
*/

type dataSourceWorkspaceOutboundCloudConnectionRulesModel struct {
	baseWorkspaceOutboundCloudConnectionRulesModel

	Timeouts timeoutsD.Value `tfsdk:"timeouts"`
}

/*
RESOURCE
*/

type resourceWorkspaceOutboundCloudConnectionRulesModel struct {
	baseWorkspaceOutboundCloudConnectionRulesModel

	Timeouts timeoutsR.Value `tfsdk:"timeouts"`
}

type requestSetWorkspaceOutboundCloudConnectionRules struct {
	fabcore.WorkspaceOutboundConnections
}

func (to *requestSetWorkspaceOutboundCloudConnectionRules) set(ctx context.Context, from resourceWorkspaceOutboundCloudConnectionRulesModel) diag.Diagnostics {
	to.DefaultAction = (*fabcore.ConnectionAccessActionType)(from.DefaultAction.ValueStringPointer())

	rules, diags := from.Rules.Get(ctx)
	if diags.HasError() {
		return diags
	}

	rulesSlice := make([]fabcore.OutboundConnectionRule, 0, len(rules))
	for _, rule := range rules {
		ruleM := fabcore.OutboundConnectionRule{
			ConnectionType: rule.ConnectionType.ValueStringPointer(),
			DefaultAction:  (*fabcore.ConnectionAccessActionType)(rule.DefaultAction.ValueStringPointer()),
		}

		endpoints, diags := rule.AllowedEndpoints.Get(ctx)
		if diags.HasError() {
			return diags
		}

		endpointsSlice := make([]fabcore.ConnectionRuleEndpointMetadata, 0, len(endpoints))
		for _, endpoint := range endpoints {
			endpointM := fabcore.ConnectionRuleEndpointMetadata{
				HostnamePattern: endpoint.HostnamePattern.ValueStringPointer(),
			}
			endpointsSlice = append(endpointsSlice, endpointM)
		}

		ruleM.AllowedEndpoints = endpointsSlice

		workspaces, diags := rule.AllowedWorkspaces.Get(ctx)
		if diags.HasError() {
			return diags
		}

		workspacesSlice := make([]fabcore.ConnectionRuleWorkspaceMetadata, 0, len(workspaces))
		for _, workspace := range workspaces {
			workspaceM := fabcore.ConnectionRuleWorkspaceMetadata{
				WorkspaceID: workspace.WorkspaceID.ValueStringPointer(),
			}
			workspacesSlice = append(workspacesSlice, workspaceM)
		}

		ruleM.AllowedWorkspaces = workspacesSlice
		rulesSlice = append(rulesSlice, ruleM)
	}

	to.Rules = rulesSlice

	return nil
}
