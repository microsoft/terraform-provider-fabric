// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package operationsagent

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	faboperationsagent "github.com/microsoft/fabric-sdk-go/fabric/operationsagent"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
)

func getResourceOperationsAgentPropertiesAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"state": schema.StringAttribute{
			MarkdownDescription: "The current state of the OperationsAgent. Possible values:" + utils.ConvertStringSlicesToString(faboperationsagent.PossibleAgentStateValues(), true, true) + ".",
			Computed:            true,
		},
	}
}
