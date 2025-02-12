// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package semanticmodel

import (
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
)

const (
	ItemName                  = "Semantic Model"
	ItemTFName                = "semantic_model"
	ItemsName                 = "Semantic Models"
	ItemsTFName               = "semantic_models"
	ItemType                  = fabcore.ItemTypeSemanticModel
	ItemDocsSPNSupport        = common.DocsSPNSupported
	ItemDocsURL               = "https://learn.microsoft.com/power-bi/developer/projects/projects-dataset"
	ItemDefinitionPathDocsURL = "https://learn.microsoft.com/rest/api/fabric/articles/item-management/definitions/semantic-model-definition"
)

var itemDefinitionFormats = []fabricitem.DefinitionFormat{ //nolint:gochecknoglobals
	{
		Type:  "TMSL",
		API:   "TMSL",
		Paths: []string{"model.bim", "definition.pbism", "diagramLayp.json"},
	},
	{
		Type:  "TMDL",
		API:   "TMDL",
		Paths: []string{"definition/database.tmdl", "definition/model.tmdl", "definition/expressions.tmdl", "definition/relationships.tmdl", "definition.pbism", "diagramLayp.json", "definition/tables/*.tmdl"},
	},
}
