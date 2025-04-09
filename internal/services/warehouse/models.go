// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package warehouse

import (
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
	fabwarehouse "github.com/microsoft/fabric-sdk-go/fabric/warehouse"
)

type warehousePropertiesModel struct {
	CollationType    types.String      `tfsdk:"collation_type"`
	ConnectionString types.String      `tfsdk:"connection_string"`
	CreatedDate      timetypes.RFC3339 `tfsdk:"created_date"`
	LastUpdatedTime  timetypes.RFC3339 `tfsdk:"last_updated_time"`
}

type warehouseConfigurationModel struct {
	CollationType types.String `tfsdk:"collation_type"`
}

func (to *warehousePropertiesModel) set(from fabwarehouse.Properties) {
	to.CollationType = types.StringPointerValue((*string)(from.CollationType))
	to.ConnectionString = types.StringPointerValue(from.ConnectionString)
	to.CreatedDate = timetypes.NewRFC3339TimePointerValue(from.CreatedDate)
	to.LastUpdatedTime = timetypes.NewRFC3339TimePointerValue(from.LastUpdatedTime)
}
