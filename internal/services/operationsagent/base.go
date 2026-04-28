// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package operationsagent

import (
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"
)

const (
	FabricItemType            = fabcore.ItemTypeOperationsAgent
	ItemDefinitionEmpty       = `{}`
	ItemDefinitionPathDocsURL = "https://learn.microsoft.com/rest/api/fabric/articles/item-management/definitions/operations-agent-definition#operationsagentconfiguration-contents"
)

var ItemTypeInfo = tftypeinfo.TFTypeInfo{ //nolint:gochecknoglobals
	Name:           "Operations Agent",
	Type:           "operations_agent",
	Names:          "Operations Agents",
	Types:          "operations_agents",
	DocsURL:        "https://learn.microsoft.com/fabric/real-time-intelligence/operations-agent",
	IsPreview:      true,
	IsSPNSupported: false,
}

var itemDefinitionFormats = []fabricitem.DefinitionFormat{ //nolint:gochecknoglobals
	{
		Type:  fabricitem.DefinitionFormatDefault,
		API:   "",
		Paths: []string{"Configurations.json"},
	},
}
