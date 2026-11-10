// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package monitoring

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func getResourceMonitoringPropertiesAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"data_collection_state": schema.StringAttribute{
			MarkdownDescription: "The observed data collection state. Manage enablement separately with `fabric_monitoring_data_collection`; changing the definition does not enable collection.",
			Computed:            true,
		},
	}
}
