// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package workspacencp

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

type baseWorkspaceNetworkCommunicationPolicyModel struct {
	ID          customtypes.UUID                                 `tfsdk:"id"`
	WorkspaceID customtypes.UUID                                 `tfsdk:"workspace_id"`
	Inbound     supertypes.SingleNestedObjectValueOf[rulesModel] `tfsdk:"inbound"`
	Outbound    supertypes.SingleNestedObjectValueOf[rulesModel] `tfsdk:"outbound"`
}

func (to *baseWorkspaceNetworkCommunicationPolicyModel) set(ctx context.Context, workspaceID string, from fabcore.WorkspaceNetworkingCommunicationPolicy) diag.Diagnostics {
	to.ID = customtypes.NewUUIDValue(workspaceID)
	to.WorkspaceID = customtypes.NewUUIDValue(workspaceID)

	inbound := supertypes.NewSingleNestedObjectValueOfNull[rulesModel](ctx)

	if from.Inbound != nil && from.Inbound.PublicAccessRules != nil {
		inboundModel := &rulesModel{}
		if diags := inboundModel.set(ctx, *from.Inbound.PublicAccessRules); diags.HasError() {
			return diags
		}

		if diags := inbound.Set(ctx, inboundModel); diags.HasError() {
			return diags
		}
	}

	to.Inbound = inbound

	outbound := supertypes.NewSingleNestedObjectValueOfNull[rulesModel](ctx)

	if from.Outbound != nil && from.Outbound.PublicAccessRules != nil {
		outboundModel := &rulesModel{}
		if diags := outboundModel.set(ctx, *from.Outbound.PublicAccessRules); diags.HasError() {
			return diags
		}

		if diags := outbound.Set(ctx, outboundModel); diags.HasError() {
			return diags
		}
	}

	to.Outbound = outbound

	return nil
}

type rulesModel struct {
	PublicAccessRules supertypes.SingleNestedObjectValueOf[networkRulesModel] `tfsdk:"public_access_rules"`
}

type networkRulesModel struct {
	DefaultAction types.String `tfsdk:"default_action"`
}

func (to *rulesModel) set(ctx context.Context, from fabcore.NetworkRules) diag.Diagnostics {
	publicAccessRules := supertypes.NewSingleNestedObjectValueOfNull[networkRulesModel](ctx)

	publicAccessRulesModel := &networkRulesModel{}

	publicAccessRulesModel.DefaultAction = types.StringPointerValue((*string)(from.DefaultAction))
	if diags := publicAccessRules.Set(ctx, publicAccessRulesModel); diags.HasError() {
		return diags
	}

	to.PublicAccessRules = publicAccessRules

	return nil
}

/*
DATA-SOURCE
*/

type dataSourceWorkspaceNetworkCommunicationPolicyModel struct {
	baseWorkspaceNetworkCommunicationPolicyModel

	Timeouts timeoutsD.Value `tfsdk:"timeouts"`
}

/*
RESOURCE
*/

type resourceWorkspaceNetworkCommunicationPolicyModel struct {
	baseWorkspaceNetworkCommunicationPolicyModel

	Timeouts timeoutsR.Value `tfsdk:"timeouts"`
}

type requestSetWorkspaceNetworkCommunicationPolicy struct {
	fabcore.WorkspaceNetworkingCommunicationPolicy
}

func setNetworkRules(ctx context.Context, rules supertypes.SingleNestedObjectValueOf[rulesModel]) (*fabcore.NetworkRules, diag.Diagnostics) {
	if rules.IsNull() || rules.IsUnknown() {
		return nil, nil
	}

	rulesVal, diags := rules.Get(ctx)
	if diags.HasError() {
		return nil, diags
	}

	if rulesVal.PublicAccessRules.IsNull() || rulesVal.PublicAccessRules.IsUnknown() {
		return nil, nil
	}

	publicAccessRulesModel, diags := rulesVal.PublicAccessRules.Get(ctx)
	if diags.HasError() {
		return nil, diags
	}

	networkRules := &fabcore.NetworkRules{}

	if !publicAccessRulesModel.DefaultAction.IsNull() && !publicAccessRulesModel.DefaultAction.IsUnknown() {
		networkRules.DefaultAction = (*fabcore.NetworkAccessRule)(publicAccessRulesModel.DefaultAction.ValueStringPointer())
	}

	return networkRules, nil
}

func (to *requestSetWorkspaceNetworkCommunicationPolicy) set(ctx context.Context, from resourceWorkspaceNetworkCommunicationPolicyModel) diag.Diagnostics {
	inboundRules, diags := setNetworkRules(ctx, from.Inbound)
	if diags.HasError() {
		return diags
	}

	if inboundRules != nil {
		to.Inbound = &fabcore.InboundRules{PublicAccessRules: inboundRules}
	}

	outboundRules, diags := setNetworkRules(ctx, from.Outbound)
	if diags.HasError() {
		return diags
	}

	if outboundRules != nil {
		to.Outbound = &fabcore.OutboundRules{PublicAccessRules: outboundRules}
	}

	return nil
}
