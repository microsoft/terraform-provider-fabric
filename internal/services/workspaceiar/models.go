// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package workspaceiar

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

type baseWorkspaceInboundAzureResourceRulesModel struct {
	WorkspaceID customtypes.UUID                                                 `tfsdk:"workspace_id"`
	Rules       supertypes.SetNestedObjectValueOf[inboundAzureResourceRuleModel] `tfsdk:"rules"`
}

type inboundAzureResourceRuleModel struct {
	DisplayName types.String `tfsdk:"display_name"`
	ResourceID  types.String `tfsdk:"resource_id"`
}

func (to *baseWorkspaceInboundAzureResourceRulesModel) set(ctx context.Context, workspaceID string, from fabcore.WorkspaceInboundAzureResourceRules) diag.Diagnostics {
	to.WorkspaceID = customtypes.NewUUIDValue(workspaceID)

	slice := make([]*inboundAzureResourceRuleModel, 0, len(from.Rules))

	for _, prop := range from.Rules {
		inboundAzureResourceRule := &inboundAzureResourceRuleModel{}
		inboundAzureResourceRule.set(prop)
		slice = append(slice, inboundAzureResourceRule)
	}

	if diags := to.Rules.Set(ctx, slice); diags.HasError() {
		return diags
	}

	return nil
}

func (to *inboundAzureResourceRuleModel) set(from fabcore.WorkspaceInboundAzureResourceRule) {
	to.DisplayName = types.StringPointerValue(from.DisplayName)
	to.ResourceID = types.StringPointerValue(from.ResourceID)
}

/*
DATA-SOURCE
*/

type dataSourceWorkspaceInboundAzureResourceRulesModel struct {
	baseWorkspaceInboundAzureResourceRulesModel

	Timeouts timeoutsD.Value `tfsdk:"timeouts"`
}

/*
RESOURCE
*/

type resourceWorkspaceInboundAzureResourceRulesModel struct {
	baseWorkspaceInboundAzureResourceRulesModel

	Timeouts timeoutsR.Value `tfsdk:"timeouts"`
}

type requestSetWorkspaceInboundAzureResourceRules struct {
	fabcore.WorkspaceInboundAzureResourceRules
}

func (to *requestSetWorkspaceInboundAzureResourceRules) set(ctx context.Context, from resourceWorkspaceInboundAzureResourceRulesModel) diag.Diagnostics {
	rules, diags := from.Rules.Get(ctx)
	if diags.HasError() {
		return diags
	}

	rulesSlice := make([]fabcore.WorkspaceInboundAzureResourceRule, 0, len(rules))

	for _, prop := range rules {
		rulesSlice = append(rulesSlice, fabcore.WorkspaceInboundAzureResourceRule{
			DisplayName: prop.DisplayName.ValueStringPointer(),
			ResourceID:  prop.ResourceID.ValueStringPointer(),
		})
	}

	to.Rules = rulesSlice

	return nil
}
