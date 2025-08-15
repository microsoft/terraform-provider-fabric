// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package warehousesnapshot

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/microsoft/fabric-sdk-go/fabric"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	fabwarehousesnapshot "github.com/microsoft/fabric-sdk-go/fabric/warehousesnapshot"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
)

func NewDataSourceWarehouseSnapshot() datasource.DataSource {
	propertiesSetter := func(ctx context.Context, from *fabwarehousesnapshot.Properties, to *fabricitem.DataSourceFabricItemPropertiesModel[warehouseSnapshotPropertiesModel, fabwarehousesnapshot.Properties]) diag.Diagnostics {
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

	itemGetter := func(ctx context.Context, fabricClient fabric.Client, model fabricitem.DataSourceFabricItemPropertiesModel[warehouseSnapshotPropertiesModel, fabwarehousesnapshot.Properties], fabricItem *fabricitem.FabricItemProperties[fabwarehousesnapshot.Properties]) error {
		client := fabwarehousesnapshot.NewClientFactoryWithClient(fabricClient).NewItemsClient()

		respGet, err := client.GetWarehouseSnapshot(ctx, model.WorkspaceID.ValueString(), model.ID.ValueString(), nil)
		if err != nil {
			return err
		}

		fabricItem.Set(respGet.WarehouseSnapshot)

		return nil
	}

	itemListGetter := func(ctx context.Context, fabricClient fabric.Client, model fabricitem.DataSourceFabricItemPropertiesModel[warehouseSnapshotPropertiesModel, fabwarehousesnapshot.Properties], errNotFound fabcore.ResponseError, fabricItem *fabricitem.FabricItemProperties[fabwarehousesnapshot.Properties]) error {
		client := fabwarehousesnapshot.NewClientFactoryWithClient(fabricClient).NewItemsClient()

		pager := client.NewListWarehouseSnapshotsPager(model.WorkspaceID.ValueString(), nil)
		for pager.More() {
			page, err := pager.NextPage(ctx)
			if err != nil {
				return err
			}

			for _, entity := range page.Value {
				if *entity.DisplayName == model.DisplayName.ValueString() {
					fabricItem.Set(entity)

					return nil
				}
			}
		}

		return &errNotFound
	}

	config := fabricitem.DataSourceFabricItemProperties[warehouseSnapshotPropertiesModel, fabwarehousesnapshot.Properties]{
		DataSourceFabricItem: fabricitem.DataSourceFabricItem{
			TypeInfo:            ItemTypeInfo,
			FabricItemType:      FabricItemType,
			IsDisplayNameUnique: true,
		},
		PropertiesAttributes: getDataSourceWarehouseSnapshotPropertiesAttributes(),
		PropertiesSetter:     propertiesSetter,
		ItemGetter:           itemGetter,
		ItemListGetter:       itemListGetter,
	}

	return fabricitem.NewDataSourceFabricItemProperties(config)
}
