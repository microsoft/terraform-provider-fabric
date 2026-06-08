// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package mirroredcatalog

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/microsoft/fabric-sdk-go/fabric"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	fabmirroredcatalog "github.com/microsoft/fabric-sdk-go/fabric/mirroredcatalog"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
)

func NewDataSourceMirroredCatalog(ctx context.Context) datasource.DataSource {
	propertiesSetter := func(ctx context.Context, from *fabmirroredcatalog.Properties, to *fabricitem.DataSourceFabricItemDefinitionPropertiesModel[mirroredCatalogPropertiesModel, fabmirroredcatalog.Properties]) diag.Diagnostics {
		properties := supertypes.NewSingleNestedObjectValueOfNull[mirroredCatalogPropertiesModel](ctx)

		if from != nil {
			propertiesModel := &mirroredCatalogPropertiesModel{}

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

	itemGetter := func(ctx context.Context, fabricClient fabric.Client, model fabricitem.DataSourceFabricItemDefinitionPropertiesModel[mirroredCatalogPropertiesModel, fabmirroredcatalog.Properties], fabricItem *fabricitem.FabricItemProperties[fabmirroredcatalog.Properties]) error {
		client := fabmirroredcatalog.NewClientFactoryWithClient(fabricClient).NewItemsClient()

		respGet, err := client.GetMirroredCatalog(ctx, model.WorkspaceID.ValueString(), model.ID.ValueString(), nil)
		if err != nil {
			return err
		}

		fabricItem.Set(respGet.MirroredCatalog)

		return nil
	}

	itemListGetter := func(ctx context.Context, fabricClient fabric.Client, model fabricitem.DataSourceFabricItemDefinitionPropertiesModel[mirroredCatalogPropertiesModel, fabmirroredcatalog.Properties], errNotFound fabcore.ResponseError, fabricItem *fabricitem.FabricItemProperties[fabmirroredcatalog.Properties]) error {
		client := fabmirroredcatalog.NewClientFactoryWithClient(fabricClient).NewItemsClient()

		pager := client.NewListMirroredCatalogsPager(model.WorkspaceID.ValueString(), nil)
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

	config := fabricitem.DataSourceFabricItemDefinitionProperties[mirroredCatalogPropertiesModel, fabmirroredcatalog.Properties]{
		DataSourceFabricItemDefinition: fabricitem.DataSourceFabricItemDefinition{
			TypeInfo:            ItemTypeInfo,
			FabricItemType:      FabricItemType,
			IsDisplayNameUnique: true,
			DefinitionFormats:   itemDefinitionFormats,
		},
		PropertiesAttributes: getDataSourceMirroredCatalogPropertiesAttributes(ctx),
		PropertiesSetter:     propertiesSetter,
		ItemGetter:           itemGetter,
		ItemListGetter:       itemListGetter,
	}

	return fabricitem.NewDataSourceFabricItemDefinitionProperties(config)
}
