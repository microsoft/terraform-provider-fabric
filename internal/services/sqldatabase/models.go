// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package sqldatabase

import (
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
	fabsqldatabase "github.com/microsoft/fabric-sdk-go/fabric/sqldatabase"
)

type sqlDatabasePropertiesModel struct {
	ConnectionString     types.String      `tfsdk:"connection_string"`
	DatabaseName         types.String      `tfsdk:"database_name"`
	ServerFqdn           types.String      `tfsdk:"server_fqdn"`
	BackupRetentionDays  types.Int32       `tfsdk:"backup_retention_days"`
	Collation            types.String      `tfsdk:"collation"`
	EarliestRestorePoint timetypes.RFC3339 `tfsdk:"earliest_restore_point"`
	LatestRestorePoint   timetypes.RFC3339 `tfsdk:"latest_restore_point"`
}

func (to *sqlDatabasePropertiesModel) set(from fabsqldatabase.Properties) {
	to.ConnectionString = types.StringPointerValue(from.ConnectionString)
	to.DatabaseName = types.StringPointerValue(from.DatabaseName)
	to.ServerFqdn = types.StringPointerValue(from.ServerFqdn)
	to.BackupRetentionDays = types.Int32PointerValue(from.BackupRetentionDays)
	to.Collation = types.StringPointerValue(from.Collation)
	to.EarliestRestorePoint = timetypes.NewRFC3339TimePointerValue(from.EarliestRestorePoint)
	to.LatestRestorePoint = timetypes.NewRFC3339TimePointerValue(from.LatestRestorePoint)
}
