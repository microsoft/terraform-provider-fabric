// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package sqldatabase

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/microsoft/fabric-sdk-go/fabric"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	fabsqldatabase "github.com/microsoft/fabric-sdk-go/fabric/sqldatabase"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
)

func NewDataSourceSQLDatabase() datasource.DataSource {
	propertiesSetter := func(ctx context.Context, from *fabsqldatabase.Properties, to *fabricitem.DataSourceFabricItemPropertiesModel[sqlDatabasePropertiesModel, fabsqldatabase.Properties]) diag.Diagnostics {
		properties := supertypes.NewSingleNestedObjectValueOfNull[sqlDatabasePropertiesModel](ctx)

		if from != nil {
			propertiesModel := &sqlDatabasePropertiesModel{}
			propertiesModel.set(*from)

			if diags := properties.Set(ctx, propertiesModel); diags.HasError() {
				return diags
			}
		}

		to.Properties = properties

		return nil
	}

	itemGetter := func(ctx context.Context, fabricClient fabric.Client, model fabricitem.DataSourceFabricItemPropertiesModel[sqlDatabasePropertiesModel, fabsqldatabase.Properties], fabricItem *fabricitem.FabricItemProperties[fabsqldatabase.Properties]) error {
		client := fabsqldatabase.NewClientFactoryWithClient(fabricClient).NewItemsClient()

		respGet, err := client.GetSQLDatabase(ctx, model.WorkspaceID.ValueString(), model.ID.ValueString(), nil)
		if err != nil {
			return err
		}

		fabricItem.Set(respGet.SQLDatabase)

		return nil
	}

	itemListGetter := func(ctx context.Context, fabricClient fabric.Client, model fabricitem.DataSourceFabricItemPropertiesModel[sqlDatabasePropertiesModel, fabsqldatabase.Properties], errNotFound fabcore.ResponseError, fabricItem *fabricitem.FabricItemProperties[fabsqldatabase.Properties]) error {
		client := fabsqldatabase.NewClientFactoryWithClient(fabricClient).NewItemsClient()

		pager := client.NewListSQLDatabasesPager(model.WorkspaceID.ValueString(), nil)
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

	config := fabricitem.DataSourceFabricItemProperties[sqlDatabasePropertiesModel, fabsqldatabase.Properties]{
		DataSourceFabricItem: fabricitem.DataSourceFabricItem{
			TypeInfo:            ItemTypeInfo,
			FabricItemType:      FabricItemType,
			IsDisplayNameUnique: true,
		},
		PropertiesAttributes: getDataSourceSQLDatabasePropertiesAttributes(),
		PropertiesSetter:     propertiesSetter,
		ItemGetter:           itemGetter,
		ItemListGetter:       itemListGetter,
	}

	return fabricitem.NewDataSourceFabricItemProperties(config)
}
