// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package sparkjobdefinition

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func getDataSourceSparkJobDefinitionPropertiesAttributes() map[string]schema.Attribute {
	result := map[string]schema.Attribute{
		"onelake_root_path": schema.StringAttribute{
			MarkdownDescription: "OneLake path to the Spark Job Definition root directory.",
			Computed:            true,
		},
	}

	return result
}
