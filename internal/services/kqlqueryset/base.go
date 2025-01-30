// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package kqlqueryset

import (
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
)

const (
	ItemName                  = "KQL Queryset"
	ItemTFName                = "kql_queryset"
	ItemsName                 = "KQL Querysets"
	ItemsTFName               = "kql_querysets"
	ItemType                  = fabcore.ItemTypeKQLQueryset
	ItemDocsSPNSupport        = common.DocsSPNSupported
	ItemDocsURL               = "https://learn.microsoft.com/fabric/real-time-intelligence/kusto-query-set"
	ItemDefinitionEmpty       = `{}`
	ItemDefinitionPathDocsURL = "https://learn.microsoft.com/rest/api/fabric/articles/item-management/definitions/kql-queryset-definition"
	ItemPreview               = true
)

var itemDefinitionFormats = []fabricitem.DefinitionFormat{ //nolint:gochecknoglobals
	{
		Type:  fabricitem.DefinitionFormatDefault,
		API:   "",
		Paths: []string{"RealTimeQueryset.json"},
	},
}
