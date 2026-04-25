// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package snowflakedatabase

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/microsoft/fabric-sdk-go/fabric"
	fabsnowflakedatabase "github.com/microsoft/fabric-sdk-go/fabric/snowflakedatabase"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
)

func NewDataSourceSnowflakeDatabases(ctx context.Context) datasource.DataSource {
	propertiesSetter := func(ctx context.Context, from *fabsnowflakedatabase.Properties, to *fabricitem.FabricItemPropertiesModel[snowflakeDatabasePropertiesModel, fabsnowflakedatabase.Properties]) diag.Diagnostics {
		properties := supertypes.NewSingleNestedObjectValueOfNull[snowflakeDatabasePropertiesModel](ctx)

		if from != nil {
			propertiesModel := &snowflakeDatabasePropertiesModel{}

			if diags := propertiesModel.set(ctx, *from); diags.HasError() {
				return diags
			}

			if diags := properties.Set(ctx, propertiesModel); diags.HasError() {
				return diags
			}
		}

		to.Properties = properties

		return nil
	}

	itemListGetter := func(ctx context.Context, fabricClient fabric.Client, model fabricitem.DataSourceFabricItemsPropertiesModel[snowflakeDatabasePropertiesModel, fabsnowflakedatabase.Properties], fabricItems *[]fabricitem.FabricItemProperties[fabsnowflakedatabase.Properties]) error {
		client := fabsnowflakedatabase.NewClientFactoryWithClient(fabricClient).NewItemsClient()

		fabItems := make([]fabricitem.FabricItemProperties[fabsnowflakedatabase.Properties], 0)

		respList, err := client.ListSnowflakeDatabases(ctx, model.WorkspaceID.ValueString(), nil)
		if err != nil {
			return err
		}

		for _, entity := range respList {
			var fabricItem fabricitem.FabricItemProperties[fabsnowflakedatabase.Properties]

			fabricItem.Set(entity)

			fabItems = append(fabItems, fabricItem)
		}

		*fabricItems = fabItems

		return nil
	}

	config := fabricitem.DataSourceFabricItemsProperties[snowflakeDatabasePropertiesModel, fabsnowflakedatabase.Properties]{
		DataSourceFabricItems: fabricitem.DataSourceFabricItems{
			TypeInfo:       ItemTypeInfo,
			FabricItemType: FabricItemType,
		},
		PropertiesAttributes: getDataSourceSnowflakeDatabasePropertiesAttributes(ctx),
		PropertiesSetter:     propertiesSetter,
		ItemListGetter:       itemListGetter,
	}

	return fabricitem.NewDataSourceFabricItemsProperties(config)
}
