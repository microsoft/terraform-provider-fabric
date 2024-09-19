// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package semanticmodel

import (
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
)

const (
	ItemName                  = "Semantic Model"
	ItemTFName                = "semantic_model"
	ItemsName                 = "Semantic Models"
	ItemsTFName               = "semantic_models"
	ItemType                  = fabcore.ItemTypeSemanticModel
	ItemDocsSPNSupport        = common.DocsSPNSupported
	ItemDocsURL               = "https://learn.microsoft.com/power-bi/developer/projects/projects-dataset"
	ItemFormatTypeDefault     = "TMSL"
	ItemDefinitionPathDocsURL = "https://learn.microsoft.com/rest/api/fabric/articles/item-management/definitions/semantic-model-definition"
)

var (
	ItemFormatTypes         = []string{"TMSL"}                                                //nolint:gochecknoglobals
	ItemDefinitionPathsTMSL = []string{"model.bim", "definition.pbism", "diagramLayout.json"} //nolint:gochecknoglobals
)
