// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package datapipeline

import (
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
)

const (
	ItemName                  = "Data Pipeline"
	ItemTFName                = "data_pipeline"
	ItemsName                 = "Data Pipelines"
	ItemsTFName               = "data_pipelines"
	ItemType                  = fabcore.ItemTypeDataPipeline
	ItemDocsSPNSupport        = common.DocsSPNSupported
	ItemDocsURL               = "https://learn.microsoft.com/rest/api/fabric/articles/item-management/definitions/datapipeline-definition"
	ItemDefinitionEmpty       = `{"properties":{"activities":[]}}`
	ItemDefinitionPathDocsURL = "https://learn.microsoft.com/fabric/data-factory/pipeline-rest-api"
	ItemPreview               = true
)

var itemDefinitionFormats = []fabricitem.DefinitionFormat{ //nolint:gochecknoglobals
	{
		Type:  fabricitem.DefinitionFormatDefault,
		API:   "",
		Paths: []string{"pipeline-content.json"},
	},
}
