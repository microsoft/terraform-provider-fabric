// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package mirroredcatalog

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/microsoft/fabric-sdk-go/fabric"
	fabmirroredcatalog "github.com/microsoft/fabric-sdk-go/fabric/mirroredcatalog"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
)

func NewDataSourceMirroredCatalogs(ctx context.Context) datasource.DataSource {
	propertiesSetter := func(ctx context.Context, from *fabmirroredcatalog.Properties, to *fabricitem.FabricItemPropertiesModel[mirroredCatalogPropertiesModel, fabmirroredcatalog.Properties]) diag.Diagnostics {
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

	itemListGetter := func(ctx context.Context, fabricClient fabric.Client, model fabricitem.DataSourceFabricItemsPropertiesModel[mirroredCatalogPropertiesModel, fabmirroredcatalog.Properties], fabricItems *[]fabricitem.FabricItemProperties[fabmirroredcatalog.Properties]) error {
		client := fabmirroredcatalog.NewClientFactoryWithClient(fabricClient).NewItemsClient()

		fabItems := make([]fabricitem.FabricItemProperties[fabmirroredcatalog.Properties], 0)

		respList, err := client.ListMirroredCatalogs(ctx, model.WorkspaceID.ValueString(), nil)
		if err != nil {
			return err
		}

		for _, entity := range respList {
			var fabricItem fabricitem.FabricItemProperties[fabmirroredcatalog.Properties]
			fabricItem.Set(entity)
			fabItems = append(fabItems, fabricItem)
		}

		*fabricItems = fabItems

		return nil
	}

	config := fabricitem.DataSourceFabricItemsProperties[mirroredCatalogPropertiesModel, fabmirroredcatalog.Properties]{
		DataSourceFabricItems: fabricitem.DataSourceFabricItems{
			TypeInfo:       ItemTypeInfo,
			FabricItemType: FabricItemType,
		},
		PropertiesAttributes: getDataSourceMirroredCatalogPropertiesAttributes(ctx),
		PropertiesSetter:     propertiesSetter,
		ItemListGetter:       itemListGetter,
	}

	return fabricitem.NewDataSourceFabricItemsProperties(config)
}
