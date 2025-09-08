// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package warehousesnapshot

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/microsoft/fabric-sdk-go/fabric"
	fabwarehousesnapshot "github.com/microsoft/fabric-sdk-go/fabric/warehousesnapshot"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
)

func NewDataSourceWarehouseSnapshots() datasource.DataSource {
	propertiesSetter := func(ctx context.Context, from *fabwarehousesnapshot.Properties, to *fabricitem.FabricItemPropertiesModel[warehouseSnapshotPropertiesModel, fabwarehousesnapshot.Properties]) diag.Diagnostics {
		properties := supertypes.NewSingleNestedObjectValueOfNull[warehouseSnapshotPropertiesModel](ctx)

		if from != nil {
			propertiesModel := &warehouseSnapshotPropertiesModel{}

			propertiesModel.set(*from)

			if diags := properties.Set(ctx, propertiesModel); diags.HasError() {
				return diags
			}
		}

		to.Properties = properties

		return nil
	}

	itemListGetter := func(ctx context.Context, fabricClient fabric.Client, model fabricitem.DataSourceFabricItemsPropertiesModel[warehouseSnapshotPropertiesModel, fabwarehousesnapshot.Properties], fabricItems *[]fabricitem.FabricItemProperties[fabwarehousesnapshot.Properties]) error {
		client := fabwarehousesnapshot.NewClientFactoryWithClient(fabricClient).NewItemsClient()

		fabItems := make([]fabricitem.FabricItemProperties[fabwarehousesnapshot.Properties], 0)

		respList, err := client.ListWarehouseSnapshots(ctx, model.WorkspaceID.ValueString(), nil)
		if err != nil {
			return err
		}

		for _, entity := range respList {
			var fabricItem fabricitem.FabricItemProperties[fabwarehousesnapshot.Properties]

			fabricItem.Set(entity)

			fabItems = append(fabItems, fabricItem)
		}

		*fabricItems = fabItems

		return nil
	}

	config := fabricitem.DataSourceFabricItemsProperties[warehouseSnapshotPropertiesModel, fabwarehousesnapshot.Properties]{
		DataSourceFabricItems: fabricitem.DataSourceFabricItems{
			TypeInfo:       ItemTypeInfo,
			FabricItemType: FabricItemType,
		},
		PropertiesAttributes: getDataSourceWarehouseSnapshotPropertiesAttributes(),
		PropertiesSetter:     propertiesSetter,
		ItemListGetter:       itemListGetter,
	}

	return fabricitem.NewDataSourceFabricItemsProperties(config)
}
