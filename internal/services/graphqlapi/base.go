// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package graphqlapi

import (
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
)

const (
	ItemName                  = "GraphQL API"
	ItemTFName                = "graphql_api"
	ItemsName                 = "GraphQL APIs"
	ItemsTFName               = "graphql_apis"
	ItemType                  = fabcore.ItemTypeGraphQLAPI
	ItemDocsSPNSupport        = common.DocsSPNSupported
	ItemDocsURL               = "https://learn.microsoft.com/fabric/data-engineering/api-graphql-overview"
	ItemDefinitionEmpty       = `TODO`
	ItemDefinitionPathDocsURL = "https://learn.microsoft.com/rest/api/fabric/articles/item-management/definitions/graphql-api-definition"
	ItemPreview               = true
)

var itemDefinitionFormats = []fabricitem.DefinitionFormat{ //nolint:gochecknoglobals
	{
		Type:  fabricitem.DefinitionFormatDefault,
		API:   "",
		Paths: []string{"graphql-definition.json"},
	},
}
