// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package datapipeline

import (
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
)

const (
	ItemName                  = "Data Pipeline"
	ItemTFName                = "data_pipeline"
	ItemsName                 = "Data Pipelines"
	ItemsTFName               = "data_pipelines"
	ItemType                  = fabcore.ItemTypeDataPipeline
	ItemDocsSPNSupport        = common.DocsSPNNotSupported
	ItemDocsURL               = "https://learn.microsoft.com/fabric/data-factory/data-factory-overview#data-pipelines"
	ItemDefinitionEmpty       = `{"properties":{"activities":[]}}`
	ItemDefinitionPathDocsURL = "https://learn.microsoft.com/fabric/data-factory/pipeline-rest-api"
)

var ItemDefinitionPaths = []string{"pipeline-content.json"} //nolint:gochecknoglobals
