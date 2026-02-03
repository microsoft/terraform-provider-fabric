// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package paginatedreport

import (
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"
)

const FabricItemType = fabcore.ItemTypePaginatedReport

var ItemTypeInfo = tftypeinfo.TFTypeInfo{ //nolint:gochecknoglobals
	Name:           "Paginated Report",
	Type:           "paginated_report",
	Names:          "Paginated Reports",
	Types:          "paginated_reports",
	DocsURL:        "https://learn.microsoft.com/power-bi/paginated-reports/web-authoring/get-started-paginated-formatted-table",
	IsPreview:      true,
	IsSPNSupported: false,
}
