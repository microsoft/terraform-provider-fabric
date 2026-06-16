// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package dataagent

import (
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"
)

const (
	FabricItemType            = fabcore.ItemTypeDataAgent
	ItemDefinitionEmpty       = `{}`
	ItemDefinitionPathDocsURL = "https://learn.microsoft.com/rest/api/fabric/articles/item-management/definitions/data-agent-definition"
)

var ItemTypeInfo = tftypeinfo.TFTypeInfo{ //nolint:gochecknoglobals
	Name:           "Data Agent",
	Type:           "data_agent",
	Names:          "Data Agents",
	Types:          "data_agents",
	DocsURL:        "https://learn.microsoft.com/fabric/data-science/concept-data-agent",
	IsPreview:      false,
	IsSPNSupported: true,
}

var itemDefinitionFormats = []fabricitem.DefinitionFormat{ //nolint:gochecknoglobals
	{
		Type: fabricitem.DefinitionFormatDefault,
		API:  "",
		Paths: []string{
			"Files/Config/data_agent.json",
			"Files/Config/draft/stage_config.json",
			"Files/Config/draft/*/datasource.json",
			"Files/Config/draft/*/fewshots.json",
			"Files/Config/publish_info.json",
			"Files/Config/published/stage_config.json",
			"Files/Config/published/*/datasource.json",
			"Files/Config/published/*/fewshots.json",
		},
	},
}
