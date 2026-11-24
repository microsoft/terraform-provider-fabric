// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package onelakedefaulttier

import (
	timeoutsD "github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts" //revive:disable-line:import-alias-naming
	timeoutsR "github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"   //revive:disable-line:import-alias-naming
	"github.com/hashicorp/terraform-plugin-framework/types"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
)

type baseOneLakeDefaultTierModel struct {
	ID          customtypes.UUID `tfsdk:"id"`
	WorkspaceID customtypes.UUID `tfsdk:"workspace_id"`
	DefaultTier types.String     `tfsdk:"default_tier"`
}

func (to *baseOneLakeDefaultTierModel) set(workspaceID string, tier *fabcore.OneLakeAccessTier) {
	to.ID = customtypes.NewUUIDValue(workspaceID)
	to.WorkspaceID = customtypes.NewUUIDValue(workspaceID)
	to.DefaultTier = types.StringPointerValue((*string)(tier))
}

type dataSourceOneLakeDefaultTierModel struct {
	baseOneLakeDefaultTierModel

	Timeouts timeoutsD.Value `tfsdk:"timeouts"`
}

type resourceOneLakeDefaultTierModel struct {
	baseOneLakeDefaultTierModel

	Timeouts timeoutsR.Value `tfsdk:"timeouts"`
}
