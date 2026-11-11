// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package plan

import (
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"
)

const (
	FabricItemType = fabcore.ItemTypePlan

	ItemDefinitionEmpty       = `{}`
	ItemDefinitionPathDocsURL = "https://learn.microsoft.com/rest/api/fabric/articles/item-management/definitions/plan-definition"
)

var ItemTypeInfo = tftypeinfo.TFTypeInfo{ //nolint:gochecknoglobals
	Name:           "Plan",
	Type:           "plan",
	Names:          "Plans",
	Types:          "plans",
	DocsURL:        "https://learn.microsoft.com/fabric/iq/plan/overview",
	IsPreview:      false,
	IsSPNSupported: true,
}

var itemDefinitionFormats = []fabricitem.DefinitionFormat{ //nolint:gochecknoglobals
	{
		Type: fabricitem.DefinitionFormatDefault,
		API:  "",
		Paths: []string{
			"definition.json",
			"planProperties.json",
			"connectedPlanning/infobridge.json",
			"cube/cube.json",
			"sheets/*/sheet.json",
			"sheets/*/commentSettings.json",
			"sheets/*/visuals/*/dataInput.json",
			"sheets/*/visuals/*/properties.json",
			"sheets/*/visuals/*/writeback.json",
			"sheets/*/visuals/*/insertRows.json",
			"sheets/*/visuals/*/scenarios.json",
			"sheets/*/visuals/*/modelTemplate.json",
			"sheets/*/visuals/*/columnConfigs.json",
			"sheets/*/visuals/*/source.json",
			"sheets/*/visuals/*/sourceSettings.json",
			"sheets/*/visuals/*/approvals.json",
			"sheets/*/visuals/*/automations.json",
			"sheets/*/visuals/*/forms.json",
		},
	},
}
