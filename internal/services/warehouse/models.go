// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package warehouse

import (
	"context"

	supertypes "github.com/FrangipaneTeam/terraform-plugin-framework-supertypes"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
	fabwarehouse "github.com/microsoft/fabric-sdk-go/fabric/warehouse"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
)

type baseWarehouseModel struct {
	WorkspaceID customtypes.UUID                                               `tfsdk:"workspace_id"`
	ID          customtypes.UUID                                               `tfsdk:"id"`
	DisplayName types.String                                                   `tfsdk:"display_name"`
	Description types.String                                                   `tfsdk:"description"`
	Properties  supertypes.SingleNestedObjectValueOf[warehousePropertiesModel] `tfsdk:"properties"`
}

func (to *baseWarehouseModel) set(ctx context.Context, from fabwarehouse.Warehouse) {
	to.WorkspaceID = customtypes.NewUUIDPointerValue(from.WorkspaceID)
	to.ID = customtypes.NewUUIDPointerValue(from.ID)
	to.DisplayName = types.StringPointerValue(from.DisplayName)
	to.Description = types.StringPointerValue(from.Description)

	properties := supertypes.NewSingleNestedObjectValueOfNull[warehousePropertiesModel](ctx)

	if from.Properties != nil {
		propertiesModel := &warehousePropertiesModel{}
		propertiesModel.set(from.Properties)
		properties.Set(ctx, propertiesModel)
	}

	to.Properties = properties
}

type warehousePropertiesModel struct {
	ConnectionString types.String      `tfsdk:"connection_string"`
	CreatedDate      timetypes.RFC3339 `tfsdk:"created_date"`
	LastUpdatedTime  timetypes.RFC3339 `tfsdk:"last_updated_time"`
}

func (to *warehousePropertiesModel) set(from *fabwarehouse.Properties) {
	to.ConnectionString = types.StringPointerValue(from.ConnectionString)
	to.CreatedDate = timetypes.NewRFC3339TimePointerValue(from.CreatedDate)
	to.LastUpdatedTime = timetypes.NewRFC3339TimePointerValue(from.LastUpdatedTime)
}
