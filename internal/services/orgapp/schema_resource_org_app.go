// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package orgapp

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

func getResourceOrgAppPropertiesAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"onelake_root_path": schema.StringAttribute{
			MarkdownDescription: "OneLake path to the Org App root directory.",
			Computed:            true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
	}
}
