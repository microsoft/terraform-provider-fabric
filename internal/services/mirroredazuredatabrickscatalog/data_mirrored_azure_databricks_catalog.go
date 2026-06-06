// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package mirroredazuredatabrickscatalog

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/microsoft/fabric-sdk-go/fabric"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	fabmirroredazuredatabrickscatalog "github.com/microsoft/fabric-sdk-go/fabric/mirroredazuredatabrickscatalog"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
)

func NewDataSourceMirroredAzureDatabricksCatalog(ctx context.Context) datasource.DataSource {
	propertiesSetter := func(ctx context.Context, from *fabmirroredazuredatabrickscatalog.Properties, to *fabricitem.DataSourceFabricItemDefinitionPropertiesModel[mirroredAzureDatabricksCatalogPropertiesModel, fabmirroredazuredatabrickscatalog.Properties]) diag.Diagnostics {
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

	itemGetter := func(ctx context.Context, fabricClient fabric.Client, model fabricitem.DataSourceFabricItemDefinitionPropertiesModel[mirroredAzureDatabricksCatalogPropertiesModel, fabmirroredazuredatabrickscatalog.Properties], fabricItem *fabricitem.FabricItemProperties[fabmirroredazuredatabrickscatalog.Properties]) error {
		client := fabmirroredazuredatabrickscatalog.NewClientFactoryWithClient(fabricClient).NewItemsClient()

		respGet, err := client.GetMirroredAzureDatabricksCatalog(ctx, model.WorkspaceID.ValueString(), model.ID.ValueString(), nil)
		if err != nil {
			return err
		}

		fabricItem.Set(respGet.MirroredAzureDatabricksCatalog)

		return nil
	}

	itemListGetter := func(ctx context.Context, fabricClient fabric.Client, model fabricitem.DataSourceFabricItemDefinitionPropertiesModel[mirroredAzureDatabricksCatalogPropertiesModel, fabmirroredazuredatabrickscatalog.Properties], errNotFound fabcore.ResponseError, fabricItem *fabricitem.FabricItemProperties[fabmirroredazuredatabrickscatalog.Properties]) error {
		client := fabmirroredazuredatabrickscatalog.NewClientFactoryWithClient(fabricClient).NewItemsClient()

		pager := client.NewListMirroredAzureDatabricksCatalogsPager(model.WorkspaceID.ValueString(), nil)
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

	config := fabricitem.DataSourceFabricItemDefinitionProperties[mirroredAzureDatabricksCatalogPropertiesModel, fabmirroredazuredatabrickscatalog.Properties]{
		DataSourceFabricItemDefinition: fabricitem.DataSourceFabricItemDefinition{
			TypeInfo:            ItemTypeInfo,
			FabricItemType:      FabricItemType,
			IsDisplayNameUnique: true,
			DefinitionFormats:   itemDefinitionFormats,
		},
		PropertiesAttributes: getDataSourceMirroredAzureDatabricksCatalogPropertiesAttributes(ctx),
		PropertiesSetter:     propertiesSetter,
		ItemGetter:           itemGetter,
		ItemListGetter:       itemListGetter,
	}

	return fabricitem.NewDataSourceFabricItemDefinitionProperties(config)
}
