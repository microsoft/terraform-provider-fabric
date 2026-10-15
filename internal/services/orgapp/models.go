// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package orgapp

import (
	faborgapp "github.com/microsoft/fabric-sdk-go/fabric/orgapp"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
)

type orgAppPropertiesModel struct {
	OneLakeRootPath customtypes.URL `tfsdk:"onelake_root_path"`
}

func (to *orgAppPropertiesModel) set(from faborgapp.Properties) {
	to.OneLakeRootPath = customtypes.NewURLPointerValue(from.OneLakeRootPath)
}
