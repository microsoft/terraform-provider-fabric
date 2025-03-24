// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package mounteddatafactory

import (
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
)

const (
	ItemName                  = "Mounted Data Factory"
	ItemTFName                = "mounteddatafactory"
	ItemsName                 = "Mounted Data Factories"
	ItemsTFName               = "mounteddatafactories"
	ItemType                  = fabcore.ItemTypeMountedDataFactory
	ItemDocsSPNSupport        = common.DocsSPNSupported
	ItemDocsURL               = "https://learn.microsoft.com/en-us/fabric/data-factory/data-factory-overview"
	ItemDefinitionEmpty       = `{}`
	ItemDefinitionPathDocsURL = "https://learn.microsoft.com/en-us/rest/api/fabric/articles/item-management/definitions/mounted-data-factory-definition"
	ItemPreview               = false
)

var itemDefinitionFormats = []fabricitem.DefinitionFormat{ //nolint:gochecknoglobals
	{
		Type:  fabricitem.DefinitionFormatDefault,
		API:   "",
		Paths: []string{"MountedDataFactoryProperties.json"},
	},
}
