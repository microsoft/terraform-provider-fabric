// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package operationsagent

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	faboperationsagent "github.com/microsoft/fabric-sdk-go/fabric/operationsagent"
)

type operationsAgentPropertiesModel struct {
	State types.String `tfsdk:"state"`
}

func (to *operationsAgentPropertiesModel) set(from faboperationsagent.Properties) {
	to.State = types.StringPointerValue((*string)(from.State))
}
