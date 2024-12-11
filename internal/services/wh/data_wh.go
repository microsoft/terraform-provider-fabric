// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package wh

import (
	"context"

	supertypes "github.com/FrangipaneTeam/terraform-plugin-framework-supertypes"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/microsoft/fabric-sdk-go/fabric"
	fabwarehouse "github.com/microsoft/fabric-sdk-go/fabric/warehouse"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
)

func NewDataSourceWH(ctx context.Context) datasource.DataSource {
	properties := schema.SingleNestedAttribute{
		MarkdownDescription: "The " + ItemName + " properties.",
		Computed:            true,
		CustomType:          supertypes.NewSingleNestedObjectTypeOf[warehousePropertiesModel](ctx),
		Attributes: map[string]schema.Attribute{
			"connection_string": schema.StringAttribute{
				MarkdownDescription: "Connection String",
				Computed:            true,
			},
			"created_date": schema.StringAttribute{
				MarkdownDescription: "Created Date",
				Computed:            true,
				CustomType:          timetypes.RFC3339Type{},
			},
			"last_updated_time": schema.StringAttribute{
				MarkdownDescription: "Last Updated Time",
				Computed:            true,
				CustomType:          timetypes.RFC3339Type{},
			},
		},
	}

	propertiesSetter := func(ctx context.Context, from *fabwarehouse.Properties, to *fabricitem.DataSourceFabricItemPropertiesModel[warehousePropertiesModel, fabwarehouse.Properties]) diag.Diagnostics {
		properties := supertypes.NewSingleNestedObjectValueOfNull[warehousePropertiesModel](ctx)

		if from != nil {
			propertiesModel := &warehousePropertiesModel{}
			propertiesModel.set(from)
			diags := properties.Set(ctx, propertiesModel)
			if diags.HasError() {
				return diags
			}
		}

		to.Properties = properties

		return nil
	}

	itemGetter := func(ctx context.Context, fabClient fabric.Client, model fabricitem.DataSourceFabricItemPropertiesModel[warehousePropertiesModel, fabwarehouse.Properties], fabItem *fabricitem.FabricItem[fabwarehouse.Properties]) diag.Diagnostics {
		client := fabwarehouse.NewClientFactoryWithClient(fabClient).NewItemsClient()

		respGet, err := client.GetWarehouse(ctx, model.WorkspaceID.ValueString(), model.ID.ValueString(), nil)
		if diags := utils.GetDiagsFromError(ctx, err, utils.OperationRead, nil); diags.HasError() {
			return diags
		}

		fabItem.Set(respGet.Warehouse)

		return nil
	}

	config := fabricitem.DataSourceFabricItemProperties[warehousePropertiesModel, fabwarehouse.Properties]{
		Type:   ItemType,
		Name:   ItemName,
		TFName: ItemTFName,
		MarkdownDescription: "Get a Fabric " + ItemName + ".\n\n" +
			"Use this data source to fetch a [" + ItemName + "](" + ItemDocsURL + ").\n\n" +
			ItemDocsSPNSupport,
		IsDisplayNameUnique: true,
		PropertiesSchema:    &properties,
		PropertiesSetter:    &propertiesSetter,
		ItemGetter:          &itemGetter,
	}

	return fabricitem.NewDataSourceFabricItemProperties(config)
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
