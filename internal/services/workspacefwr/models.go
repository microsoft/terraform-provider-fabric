// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package workspacefwr

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

type baseWorkspaceFirewallRulesModel struct {
	WorkspaceID customtypes.UUID                                     `tfsdk:"workspace_id"`
	Rules       supertypes.SetNestedObjectValueOf[firewallRuleModel] `tfsdk:"rules"`
}

type firewallRuleModel struct {
	DisplayName types.String `tfsdk:"display_name"`
	Value       types.String `tfsdk:"value"`
}

func (to *baseWorkspaceFirewallRulesModel) set(ctx context.Context, workspaceID string, from fabcore.InboundFirewallConfiguration) diag.Diagnostics {
	to.WorkspaceID = customtypes.NewUUIDValue(workspaceID)

	slice := make([]*firewallRuleModel, 0, len(from.Rules))

	for _, prop := range from.Rules {
		firewallRule := &firewallRuleModel{}
		firewallRule.set(prop)
		slice = append(slice, firewallRule)
	}

	if diags := to.Rules.Set(ctx, slice); diags.HasError() {
		return diags
	}

	return nil
}

func (to *firewallRuleModel) set(from fabcore.FirewallRule) {
	to.DisplayName = types.StringPointerValue(from.DisplayName)
	to.Value = types.StringPointerValue(from.Value)
}

/*
DATA-SOURCE
*/

type dataSourceWorkspaceFirewallRulesModel struct {
	baseWorkspaceFirewallRulesModel

	Timeouts timeoutsD.Value `tfsdk:"timeouts"`
}

/*
RESOURCE
*/

type resourceWorkspaceFirewallRulesModel struct {
	baseWorkspaceFirewallRulesModel

	Timeouts timeoutsR.Value `tfsdk:"timeouts"`
}

type requestSetWorkspaceFirewallRules struct {
	fabcore.InboundFirewallConfiguration
}

func (to *requestSetWorkspaceFirewallRules) set(ctx context.Context, from resourceWorkspaceFirewallRulesModel) diag.Diagnostics {
	rules, diags := from.Rules.Get(ctx)
	if diags.HasError() {
		return diags
	}

	rulesSlice := make([]fabcore.FirewallRule, 0, len(rules))

	for _, prop := range rules {
		rulesSlice = append(rulesSlice, fabcore.FirewallRule{
			DisplayName: prop.DisplayName.ValueStringPointer(),
			Value:       prop.Value.ValueStringPointer(),
		})
	}

	to.Rules = rulesSlice

	return nil
}
