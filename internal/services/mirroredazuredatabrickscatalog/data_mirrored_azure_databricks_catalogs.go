// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package mirroredazuredatabrickscatalog

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/microsoft/fabric-sdk-go/fabric"
	fabmirroredazuredatabrickscatalog "github.com/microsoft/fabric-sdk-go/fabric/mirroredazuredatabrickscatalog"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
)

func NewDataSourceMirroredAzureDatabricksCatalogs(ctx context.Context) datasource.DataSource {
	propertiesSetter := func(ctx context.Context, from *fabmirroredazuredatabrickscatalog.Properties, to *fabricitem.FabricItemPropertiesModel[mirroredAzureDatabricksCatalogPropertiesModel, fabmirroredazuredatabrickscatalog.Properties]) diag.Diagnostics {
		properties := supertypes.NewSingleNestedObjectValueOfNull[mirroredAzureDatabricksCatalogPropertiesModel](ctx)

		if from != nil {
			propertiesModel := &mirroredAzureDatabricksCatalogPropertiesModel{}

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

	itemListGetter := func(ctx context.Context, fabricClient fabric.Client, model fabricitem.DataSourceFabricItemsPropertiesModel[mirroredAzureDatabricksCatalogPropertiesModel, fabmirroredazuredatabrickscatalog.Properties], fabricItems *[]fabricitem.FabricItemProperties[fabmirroredazuredatabrickscatalog.Properties]) error {
		client := fabmirroredazuredatabrickscatalog.NewClientFactoryWithClient(fabricClient).NewItemsClient()

		fabItems := make([]fabricitem.FabricItemProperties[fabmirroredazuredatabrickscatalog.Properties], 0)

		respList, err := client.ListMirroredAzureDatabricksCatalogs(ctx, model.WorkspaceID.ValueString(), nil)
		if err != nil {
			return err
		}

		for _, entity := range respList {
			var fabricItem fabricitem.FabricItemProperties[fabmirroredazuredatabrickscatalog.Properties]

			fabricItem.Set(entity)

			fabItems = append(fabItems, fabricItem)
		}

		*fabricItems = fabItems

		return nil
	}

	config := fabricitem.DataSourceFabricItemsProperties[mirroredAzureDatabricksCatalogPropertiesModel, fabmirroredazuredatabrickscatalog.Properties]{
		DataSourceFabricItems: fabricitem.DataSourceFabricItems{
			TypeInfo:       ItemTypeInfo,
			FabricItemType: FabricItemType,
		},
		PropertiesAttributes: getDataSourceMirroredAzureDatabricksCatalogPropertiesAttributes(ctx),
		PropertiesSetter:     propertiesSetter,
		ItemListGetter:       itemListGetter,
	}

	return fabricitem.NewDataSourceFabricItemsProperties(config)
}
