// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package kqldashboard

import (
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"
)

const (
	FabricItemType            = fabcore.ItemTypeKQLDashboard
	ItemDefinitionEmpty       = `{}`
	ItemDefinitionPathDocsURL = "https://learn.microsoft.com/rest/api/fabric/articles/item-management/definitions/kql-dashboard-definition"
)

var ItemTypeInfo = tftypeinfo.TFTypeInfo{ //nolint:gochecknoglobals
	Name:           "KQL Dashboard",
	Type:           "kql_dashboard",
	Names:          "KQL Dashboards",
	Types:          "kql_dashboards",
	DocsURL:        "https://learn.microsoft.com/fabric/real-time-intelligence/dashboard-real-time-create",
	IsPreview:      false,
	IsSPNSupported: true,
}

var itemDefinitionFormats = []fabricitem.DefinitionFormat{ //nolint:gochecknoglobals
	{
		Type:  fabricitem.DefinitionFormatDefault,
		API:   "",
		Paths: []string{"RealTimeDashboard.json"},
	},
}
