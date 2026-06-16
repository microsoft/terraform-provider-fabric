// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package workspacegop

import (
	timeoutsD "github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts" //revive:disable-line:import-alias-naming
	timeoutsR "github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"   //revive:disable-line:import-alias-naming
	"github.com/hashicorp/terraform-plugin-framework/types"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
)

type baseWorkspaceGitOutboundPolicyModel struct {
	WorkspaceID   customtypes.UUID `tfsdk:"workspace_id"`
	DefaultAction types.String     `tfsdk:"default_action"`
}

func (to *baseWorkspaceGitOutboundPolicyModel) set(workspaceID string, from fabcore.NetworkRules) {
	to.WorkspaceID = customtypes.NewUUIDValue(workspaceID)
	to.DefaultAction = types.StringPointerValue((*string)(from.DefaultAction))
}

/*
DATA-SOURCE
*/

type dataSourceWorkspaceGitOutboundPolicyModel struct {
	baseWorkspaceGitOutboundPolicyModel

	Timeouts timeoutsD.Value `tfsdk:"timeouts"`
}

/*
RESOURCE
*/

type resourceWorkspaceGitOutboundPolicyModel struct {
	baseWorkspaceGitOutboundPolicyModel

	Timeouts timeoutsR.Value `tfsdk:"timeouts"`
}

type requestSetWorkspaceGitOutboundPolicy struct {
	fabcore.NetworkRules
}

func (to *requestSetWorkspaceGitOutboundPolicy) set(from resourceWorkspaceGitOutboundPolicyModel) {
	to.DefaultAction = (*fabcore.NetworkAccessRule)(from.DefaultAction.ValueStringPointer())
}
