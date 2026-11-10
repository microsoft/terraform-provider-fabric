// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package monitoringdatacollection

import "github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"

var ItemTypeInfo = tftypeinfo.TFTypeInfo{ //nolint:gochecknoglobals
	Name:           "Monitoring Data Collection",
	Type:           "monitoring_data_collection",
	DocsURL:        "https://learn.microsoft.com/rest/api/fabric/monitoring/data-collection",
	IsPreview:      true,
	IsSPNSupported: true,
}
