// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package paginatedreport

import (
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
)

const (
	ItemName           = "Paginated Report"
	ItemTFName         = "paginated_report"
	ItemsName          = "Paginated Reports"
	ItemsTFName        = "paginated_reports"
	ItemType           = fabcore.ItemTypePaginatedReport
	ItemDocsSPNSupport = common.DocsSPNNotSupported
	ItemDocsURL        = "https://learn.microsoft.com/power-bi/paginated-reports/web-authoring/get-started-paginated-formatted-table"
	ItemPreview        = true
)
