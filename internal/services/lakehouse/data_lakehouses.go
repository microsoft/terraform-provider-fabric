// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package lakehouse

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/microsoft/fabric-sdk-go/fabric"
	fablakehouse "github.com/microsoft/fabric-sdk-go/fabric/lakehouse"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
)

func NewDataSourceLakehouses(ctx context.Context) datasource.DataSource {
	propertiesSetter := func(ctx context.Context, from *fablakehouse.Properties, to *fabricitem.FabricItemPropertiesModel[lakehousePropertiesModel, fablakehouse.Properties]) diag.Diagnostics {
		properties := supertypes.NewSingleNestedObjectValueOfNull[lakehousePropertiesModel](ctx)

		if from != nil {
			propertiesModel := &lakehousePropertiesModel{}

			if diags := propertiesModel.set(ctx, from); diags.HasError() {
				return diags
			}

			if diags := properties.Set(ctx, propertiesModel); diags.HasError() {
				return diags
			}
		}

		to.Properties = properties

		return nil
	}

	itemListGetter := func(ctx context.Context, fabricClient fabric.Client, model fabricitem.DataSourceFabricItemsPropertiesModel[lakehousePropertiesModel, fablakehouse.Properties], fabricItems *[]fabricitem.FabricItemProperties[fablakehouse.Properties]) error {
		client := fablakehouse.NewClientFactoryWithClient(fabricClient).NewItemsClient()

		fabItems := make([]fabricitem.FabricItemProperties[fablakehouse.Properties], 0)

		respList, err := client.ListLakehouses(ctx, model.WorkspaceID.ValueString(), nil)
		if err != nil {
			return err
		}

		for _, entity := range respList {
			var fabricItem fabricitem.FabricItemProperties[fablakehouse.Properties]

			fabricItem.Set(entity)

			fabItems = append(fabItems, fabricItem)
		}

		*fabricItems = fabItems

		return nil
	}

	config := fabricitem.DataSourceFabricItemsProperties[lakehousePropertiesModel, fablakehouse.Properties]{
		DataSourceFabricItems: fabricitem.DataSourceFabricItems{
			TypeInfo:       ItemTypeInfo,
			FabricItemType: FabricItemType,
		},
		PropertiesAttributes: getDataSourceLakehousePropertiesAttributes(ctx),
		PropertiesSetter:     propertiesSetter,
		ItemListGetter:       itemListGetter,
	}

	return fabricitem.NewDataSourceFabricItemsProperties(config)
}
