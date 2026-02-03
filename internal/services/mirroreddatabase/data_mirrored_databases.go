// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package mirroreddatabase

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/microsoft/fabric-sdk-go/fabric"
	fabmirroreddatabase "github.com/microsoft/fabric-sdk-go/fabric/mirroreddatabase"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
)

func NewDataSourceMirroredDatabases(ctx context.Context) datasource.DataSource {
	propertiesSetter := func(ctx context.Context, from *fabmirroreddatabase.Properties, to *fabricitem.FabricItemPropertiesModel[mirroredDatabasePropertiesModel, fabmirroreddatabase.Properties]) diag.Diagnostics {
		properties := supertypes.NewSingleNestedObjectValueOfNull[mirroredDatabasePropertiesModel](ctx)

		if from != nil {
			propertiesModel := &mirroredDatabasePropertiesModel{}
			propertiesModel.set(ctx, *from)

			if diags := properties.Set(ctx, propertiesModel); diags.HasError() {
				return diags
			}
		}

		to.Properties = properties

		return nil
	}

	itemListGetter := func(ctx context.Context, fabricClient fabric.Client, model fabricitem.DataSourceFabricItemsPropertiesModel[mirroredDatabasePropertiesModel, fabmirroreddatabase.Properties], fabricItems *[]fabricitem.FabricItemProperties[fabmirroreddatabase.Properties]) error {
		client := fabmirroreddatabase.NewClientFactoryWithClient(fabricClient).NewItemsClient()

		fabItems := make([]fabricitem.FabricItemProperties[fabmirroreddatabase.Properties], 0)

		respList, err := client.ListMirroredDatabases(ctx, model.WorkspaceID.ValueString(), nil)
		if err != nil {
			return err
		}

		for _, entity := range respList {
			var fabricItem fabricitem.FabricItemProperties[fabmirroreddatabase.Properties]
			fabricItem.Set(entity)
			fabItems = append(fabItems, fabricItem)
		}

		*fabricItems = fabItems

		return nil
	}

	config := fabricitem.DataSourceFabricItemsProperties[mirroredDatabasePropertiesModel, fabmirroreddatabase.Properties]{
		DataSourceFabricItems: fabricitem.DataSourceFabricItems{
			TypeInfo:       ItemTypeInfo,
			FabricItemType: FabricItemType,
		},
		PropertiesAttributes: getDataSourceMirroredDatabasePropertiesAttributes(ctx),
		PropertiesSetter:     propertiesSetter,
		ItemListGetter:       itemListGetter,
	}

	return fabricitem.NewDataSourceFabricItemsProperties(config)
}
