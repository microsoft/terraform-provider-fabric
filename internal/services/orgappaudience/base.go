// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package orgappaudience

import (
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"
)

const (
	FabricItemType = fabcore.ItemTypeOrgAppAudience

	ItemDefinitionPathDocsURL = "https://learn.microsoft.com/rest/api/fabric/articles/item-management/definitions/orgappaudience-definition"
)

var ItemTypeInfo = tftypeinfo.TFTypeInfo{ //nolint:gochecknoglobals
	Name:           "Org App Audience",
	Type:           "org_app_audience",
	Names:          "Org App Audiences",
	Types:          "org_app_audiences",
	DocsURL:        "https://learn.microsoft.com/power-bi/explore-reports/org-app-items#use-audiences-in-org-apps",
	IsPreview:      false,
	IsSPNSupported: true,
}

var itemDefinitionFormats = []fabricitem.DefinitionFormat{ //nolint:gochecknoglobals
	{
		Type:  fabricitem.DefinitionFormatDefault,
		API:   "",
		Paths: []string{"definition.json"},
	},
}
