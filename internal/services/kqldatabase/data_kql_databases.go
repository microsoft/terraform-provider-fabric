// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package kqldatabase

import (
	"context"

	supertypes "github.com/FrangipaneTeam/terraform-plugin-framework-supertypes"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/microsoft/fabric-sdk-go/fabric"
	fabkqldatabase "github.com/microsoft/fabric-sdk-go/fabric/kqldatabase"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
)

func NewDataSourceKQLDatabases() datasource.DataSource {
	propertiesSetter := func(ctx context.Context, from *fabkqldatabase.Properties, to *fabricitem.FabricItemPropertiesModel[kqlDatabasePropertiesModel, fabkqldatabase.Properties]) diag.Diagnostics {
		properties := supertypes.NewSingleNestedObjectValueOfNull[kqlDatabasePropertiesModel](ctx)

		if from != nil {
			propertiesModel := &kqlDatabasePropertiesModel{}
			propertiesModel.set(from)

			if diags := properties.Set(ctx, propertiesModel); diags.HasError() {
				return diags
			}
		}

		to.Properties = properties

		return nil
	}

	itemListGetter := func(ctx context.Context, fabricClient fabric.Client, model fabricitem.DataSourceFabricItemsPropertiesModel[kqlDatabasePropertiesModel, fabkqldatabase.Properties], fabricItems *[]fabricitem.FabricItemProperties[fabkqldatabase.Properties]) error {
		client := fabkqldatabase.NewClientFactoryWithClient(fabricClient).NewItemsClient()

		fabItems := make([]fabricitem.FabricItemProperties[fabkqldatabase.Properties], 0)

		respList, err := client.ListKQLDatabases(ctx, model.WorkspaceID.ValueString(), nil)
		if err != nil {
			return err
		}

		for _, entity := range respList {
			var fabricItem fabricitem.FabricItemProperties[fabkqldatabase.Properties]

			fabricItem.Set(entity)

			fabItems = append(fabItems, fabricItem)
		}

		*fabricItems = fabItems

		return nil
	}

	config := fabricitem.DataSourceFabricItemsProperties[kqlDatabasePropertiesModel, fabkqldatabase.Properties]{
		DataSourceFabricItems: fabricitem.DataSourceFabricItems{
			Type:   ItemType,
			Name:   ItemName,
			Names:  ItemsName,
			TFName: ItemsTFName,
			MarkdownDescription: "List a Fabric " + ItemsName + ".\n\n" +
				"Use this data source to list [" + ItemsName + "](" + ItemDocsURL + ").\n\n" +
				ItemDocsSPNSupport,
		},
		PropertiesAttributes: getDataSourceKQLDatabasePropertiesAttributes(),
		PropertiesSetter:     propertiesSetter,
		ItemListGetter:       itemListGetter,
	}

	return fabricitem.NewDataSourceFabricItemsProperties(config)
}
