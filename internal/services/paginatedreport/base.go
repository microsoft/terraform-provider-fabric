// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package paginatedreport

import (
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"
)

const (
	FabricItemType            = fabcore.ItemTypePaginatedReport
	ItemDefinitionPathDocsURL = "https://learn.microsoft.com/rest/api/fabric/articles/item-management/definitions/paginatedreport-definition"
)

var ItemTypeInfo = tftypeinfo.TFTypeInfo{ //nolint:gochecknoglobals
	Name:           "Paginated Report",
	Type:           "paginated_report",
	Names:          "Paginated Reports",
	Types:          "paginated_reports",
	DocsURL:        "https://learn.microsoft.com/rest/api/fabric/paginatedreport/items",
	IsPreview:      false,
	IsSPNSupported: true,
}

var itemDefinitionFormats = []fabricitem.DefinitionFormat{ //nolint:gochecknoglobals
	{
		Type:  fabricitem.DefinitionFormatDefault,
		API:   "",
		Paths: []string{"*.rdl"},
	},
}
