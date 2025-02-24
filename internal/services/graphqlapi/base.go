// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package graphqlapi

import (
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
)

const (
	ItemName           = "GraphQL API"
	ItemTFName         = "graphql_api"
	ItemsName          = "GraphQL APIs"
	ItemsTFName        = "graphql_apis"
	ItemType           = fabcore.ItemTypeGraphQLAPI
	ItemDocsSPNSupport = common.DocsSPNSupported
	ItemDocsURL        = "https://learn.microsoft.com/fabric/data-engineering/api-graphql-overview"
	ItemPreview        = true
)
