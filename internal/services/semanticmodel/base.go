// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package semanticmodel

import (
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"
)

const (
	FabricItemType            = fabcore.ItemTypeSemanticModel
	ItemDefinitionPathDocsURL = "https://learn.microsoft.com/rest/api/fabric/articles/item-management/definitions/semantic-model-definition"
)

var ItemTypeInfo = tftypeinfo.TFTypeInfo{ //nolint:gochecknoglobals
	Name:           "Semantic Model",
	Type:           "semantic_model",
	Names:          "Semantic Models",
	Types:          "semantic_models",
	DocsURL:        "https://learn.microsoft.com/power-bi/developer/projects/projects-dataset",
	IsPreview:      false,
	IsSPNSupported: true,
}

var itemDefinitionFormats = []fabricitem.DefinitionFormat{ //nolint:gochecknoglobals
	{
		Type:  "TMSL",
		API:   "TMSL",
		Paths: []string{"model.bim", "definition.pbism", "diagramLayout.json"},
	},
	{
		Type: "TMDL",
		API:  "TMDL",
		Paths: []string{
			"definition/database.tmdl",
			"definition/model.tmdl",
			"definition/expressions.tmdl",
			"definition/relationships.tmdl",
			"definition/dataSources.tmdl",
			"definition.pbism",
			"diagramLayout.json",
			"definition/tables/*.tmdl",
			"definition/roles/*.tmdl",
			"definition/perspectives/*.tmdl",
			"definition/cultures/*.tmdl",
		},
	},
}
