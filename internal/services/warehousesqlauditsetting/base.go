// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package warehousesqlauditsetting

import (
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"
)

var ItemTypeInfo = tftypeinfo.TFTypeInfo{ //nolint:gochecknoglobals
	Name:           "Warehouse SQL Audit Settings",
	Type:           "warehouse_sql_audit_settings",
	DocsURL:        "https://learn.microsoft.com/fabric/data-warehouse/sql-audit-logs",
	IsPreview:      false,
	IsSPNSupported: true,
}
