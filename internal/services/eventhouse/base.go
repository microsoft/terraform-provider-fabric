// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package eventhouse

import (
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"
)

const (
	ItemDefinitionEmpty       = `{}`
	ItemDefinitionPathDocsURL = "https://learn.microsoft.com/rest/api/fabric/articles/item-management/definitions/eventhouse-definition"
	FabricItemType            = fabcore.ItemTypeEventhouse
)

var ItemTypeInfo = tftypeinfo.TFTypeInfo{ //nolint:gochecknoglobals
	Name:           "Eventhouse",
	Type:           "eventhouse",
	Names:          "Eventhouses",
	Types:          "eventhouses",
	DocsURL:        "https://learn.microsoft.com/fabric/real-time-intelligence/eventhouse",
	IsPreview:      false,
	IsSPNSupported: true,
}

var itemDefinitionFormats = []fabricitem.DefinitionFormat{ //nolint:gochecknoglobals
	{
		Type:  fabricitem.DefinitionFormatDefault,
		API:   "",
		Paths: []string{"EventhouseProperties.json"},
	},
}

var (
	possibleMinimumConsumptionUnitsValues = []float64{0, 2.25, 4.25, 8.5, 13, 18, 26, 34, 50} //nolint:gochecknoglobals
	minimumConsumptionUnitsMin            = float64(51)                                       //nolint:gochecknoglobals
	minimumConsumptionUnitsMax            = float64(322)                                      //nolint:gochecknoglobals
)
