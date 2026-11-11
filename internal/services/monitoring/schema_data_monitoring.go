// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package monitoring

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func getDataSourceMonitoringPropertiesAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"data_collection_state": schema.StringAttribute{
			MarkdownDescription: "The observed data collection state of the Monitoring item.",
			Computed:            true,
		},
	}
}
