// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package monitoring

import (
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"
)

const (
	FabricItemType            = fabcore.ItemTypeMonitoring
	ItemDefinitionEmpty       = `{"$schema":"https://developer.microsoft.com/json-schemas/fabric/item/monitoring/definition/monitoringProperties/1.0.0/schema.json"}`
	ItemDefinitionPathDocsURL = "https://learn.microsoft.com/rest/api/fabric/articles/item-management/definitions/monitoring-definition"
)

var ItemTypeInfo = tftypeinfo.TFTypeInfo{ //nolint:gochecknoglobals
	Name:           "Monitoring",
	Type:           "monitoring",
	Names:          "Monitoring Items",
	Types:          "monitoring_items",
	DocsURL:        ItemDefinitionPathDocsURL,
	IsPreview:      true,
	IsSPNSupported: true,
}

var itemDefinitionFormats = []fabricitem.DefinitionFormat{ //nolint:gochecknoglobals
	{
		Type:  fabricitem.DefinitionFormatDefault,
		API:   "MonitoringV1",
		Paths: []string{"MonitoringProperties.json"},
	},
}
