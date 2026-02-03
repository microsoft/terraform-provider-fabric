// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package sparkjobdefinition

import (
	fabsparkjobdefinition "github.com/microsoft/fabric-sdk-go/fabric/sparkjobdefinition"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
)

type sparkJobDefinitionPropertiesModel struct {
	OneLakeRootPath customtypes.URL `tfsdk:"onelake_root_path"`
}

func (to *sparkJobDefinitionPropertiesModel) set(from fabsparkjobdefinition.Properties) {
	to.OneLakeRootPath = customtypes.NewURLPointerValue(from.OneLakeRootPath)
}
