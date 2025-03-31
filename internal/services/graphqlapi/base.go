// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package graphqlapi

import (
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"
)

const FabricItemType = fabcore.ItemTypeGraphQLAPI

var ItemTypeInfo = tftypeinfo.TFTypeInfo{ //nolint:gochecknoglobals
	Name:           "GraphQL API",
	Type:           "graphql_api",
	Names:          "GraphQL APIs",
	Types:          "graphql_apis",
	DocsURL:        "https://learn.microsoft.com/fabric/data-engineering/api-graphql-overview",
	IsPreview:      false,
	IsSPNSupported: true,
}
