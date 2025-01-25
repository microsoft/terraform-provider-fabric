// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package kqldashboard

import (
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
)

const (
	ItemName                  = "KQL Dashboard"
	ItemTFName                = "kql_dashboard"
	ItemsName                 = "KQL Dashboards"
	ItemsTFName               = "kql_dashboards"
	ItemType                  = fabcore.ItemTypeKQLDashboard
	ItemDocsSPNSupport        = common.DocsSPNSupported
	ItemDocsURL               = "https://learn.microsoft.com/fabric/real-time-intelligence/dashboard-real-time-create"
	ItemDefinitionEmpty       = `{}`
	ItemDefinitionPathDocsURL = "https://learn.microsoft.com/rest/api/fabric/articles/item-management/definitions/kql-dashboard-definition"
	ItemPreview               = true
)

var itemDefinitionFormats = []fabricitem.DefinitionFormat{ //nolint:gochecknoglobals
	{
		Type:  fabricitem.DefinitionFormatDefault,
		API:   "",
		Paths: []string{"RealTimeDashboard.json"},
	},
}
