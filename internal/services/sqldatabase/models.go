// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package sqldatabase

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	fabsqldatabase "github.com/microsoft/fabric-sdk-go/fabric/sqldatabase"
)

type sqlDatabasePropertiesModel struct {
	ConnectionString types.String `tfsdk:"connection_string"`
	DatabaseName     types.String `tfsdk:"database_name"`
	ServerFqdn       types.String `tfsdk:"server_fqdn"`
}

func (to *sqlDatabasePropertiesModel) set(from fabsqldatabase.Properties) {
	to.ConnectionString = types.StringPointerValue(from.ConnectionString)
	to.DatabaseName = types.StringPointerValue(from.DatabaseName)
	to.ServerFqdn = types.StringPointerValue(from.ServerFqdn)
}
