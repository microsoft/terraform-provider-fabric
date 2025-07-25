// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package report

import (
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"
)

const (
	FabricItemType            = fabcore.ItemTypeReport
	ItemDefinitionPathDocsURL = "https://learn.microsoft.com/rest/api/fabric/articles/item-management/definitions/report-definition"
)

var ItemTypeInfo = tftypeinfo.TFTypeInfo{ //nolint:gochecknoglobals
	Name:           "Report",
	Type:           "report",
	Names:          "Reports",
	Types:          "reports",
	DocsURL:        "https://learn.microsoft.com/power-bi/developer/projects/projects-report",
	IsPreview:      false,
	IsSPNSupported: true,
}

var itemDefinitionFormats = []fabricitem.DefinitionFormat{ //nolint:gochecknoglobals
	{
		Type:  "PBIR-Legacy",
		API:   "PBIR-Legacy",
		Paths: []string{"report.json", "definition.pbir", "StaticResources/RegisteredResources/*", "StaticResources/SharedResources/*"},
	},
	{
		Type: "PBIR",
		API:  "PBIR",
		Paths: []string{
			"definition/report.json",
			"definition/version.json",
			"definition.pbir",
			"definition/bookmarks/*.json",
			"definition/pages/*.json",
			"StaticResources/RegisteredResources/*",
			"StaticResources/SharedResources/*",
		},
	},
}
