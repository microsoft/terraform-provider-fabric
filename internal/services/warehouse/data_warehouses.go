// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package warehouse

import (
	"context"

	supertypes "github.com/FrangipaneTeam/terraform-plugin-framework-supertypes"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/microsoft/fabric-sdk-go/fabric"
	fabwarehouse "github.com/microsoft/fabric-sdk-go/fabric/warehouse"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
)

func NewDataSourceWarehouses(ctx context.Context) datasource.DataSource {
	propertiesSchema := schema.SingleNestedAttribute{
		MarkdownDescription: "The " + ItemName + " properties.",
		Computed:            true,
		CustomType:          supertypes.NewSingleNestedObjectTypeOf[warehousePropertiesModel](ctx),
		Attributes: map[string]schema.Attribute{
			"connection_string": schema.StringAttribute{
				MarkdownDescription: "The SQL connection string connected to the workspace containing this warehouse.",
				Computed:            true,
			},
			"created_date": schema.StringAttribute{
				MarkdownDescription: "The date and time the warehouse was created.",
				Computed:            true,
				CustomType:          timetypes.RFC3339Type{},
			},
			"last_updated_time": schema.StringAttribute{
				MarkdownDescription: "The date and time the warehouse was last updated.",
				Computed:            true,
				CustomType:          timetypes.RFC3339Type{},
			},
		},
	}

	propertiesSetter := func(ctx context.Context, from *fabwarehouse.Properties, to *fabricitem.FabricItemPropertiesModel[warehousePropertiesModel, fabwarehouse.Properties]) diag.Diagnostics {
		properties := supertypes.NewSingleNestedObjectValueOfNull[warehousePropertiesModel](ctx)

		if from != nil {
			propertiesModel := &warehousePropertiesModel{}
			propertiesModel.set(from)

			if diags := properties.Set(ctx, propertiesModel); diags.HasError() {
				return diags
			}
		}

		to.Properties = properties

		return nil
	}

	itemListGetter := func(ctx context.Context, fabricClient fabric.Client, model fabricitem.DataSourceFabricItemsPropertiesModel[warehousePropertiesModel, fabwarehouse.Properties], fabricItems *[]fabricitem.FabricItemProperties[fabwarehouse.Properties]) error {
		client := fabwarehouse.NewClientFactoryWithClient(fabricClient).NewItemsClient()

		fabItems := make([]fabricitem.FabricItemProperties[fabwarehouse.Properties], 0)

		respList, err := client.ListWarehouses(ctx, model.WorkspaceID.ValueString(), nil)
		if err != nil {
			return err
		}

		for _, entity := range respList {
			var fabricItem fabricitem.FabricItemProperties[fabwarehouse.Properties]

			fabricItem.Set(entity)

			fabItems = append(fabItems, fabricItem)
		}

		*fabricItems = fabItems

		return nil
	}

	config := fabricitem.DataSourceFabricItemsProperties[warehousePropertiesModel, fabwarehouse.Properties]{
		DataSourceFabricItems: fabricitem.DataSourceFabricItems{
			Type:   ItemType,
			Name:   ItemName,
			Names:  ItemsName,
			TFName: ItemsTFName,
			MarkdownDescription: "List a Fabric " + ItemsName + ".\n\n" +
				"Use this data source to list [" + ItemsName + "](" + ItemDocsURL + ").\n\n" +
				ItemDocsSPNSupport,
		},
		PropertiesSchema: propertiesSchema,
		PropertiesSetter: propertiesSetter,
		ItemListGetter:   itemListGetter,
	}

	return fabricitem.NewDataSourceFabricItemsProperties(config)
}
