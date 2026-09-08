// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package orgapp

import "github.com/hashicorp/terraform-plugin-framework/resource/schema"

func getResourceOrgAppPropertiesAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"onelake_root_path": schema.StringAttribute{
			MarkdownDescription: "OneLake path to the Org App root directory.",
			Computed:            true,
		},
	}
}
