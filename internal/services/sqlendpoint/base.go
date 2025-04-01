// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package sqlendpoint

import (
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"
)

const FabricItemType = fabcore.ItemTypeSQLEndpoint

var ItemTypeInfo = tftypeinfo.TFTypeInfo{ //nolint:gochecknoglobals
	Name:           "SQL Endpoint",
	Type:           "sql_endpoint",
	Names:          "SQL Endpoints",
	Types:          "sql_endpoints",
	DocsURL:        "https://learn.microsoft.com/fabric/data-warehouse/data-warehousing#sql-analytics-endpoint-of-the-lakehouse",
	IsPreview:      true,
	IsSPNSupported: false,
}
