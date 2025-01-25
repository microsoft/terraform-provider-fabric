// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package kqldatabase

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/microsoft/fabric-sdk-go/fabric"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	fabkqldatabase "github.com/microsoft/fabric-sdk-go/fabric/kqldatabase"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
)

func NewDataSourceKQLDatabase() datasource.DataSource {
	propertiesSetter := func(ctx context.Context, from *fabkqldatabase.Properties, to *fabricitem.DataSourceFabricItemDefinitionPropertiesModel[kqlDatabasePropertiesModel, fabkqldatabase.Properties]) diag.Diagnostics {
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

	itemGetter := func(ctx context.Context, fabricClient fabric.Client, model fabricitem.DataSourceFabricItemDefinitionPropertiesModel[kqlDatabasePropertiesModel, fabkqldatabase.Properties], fabricItem *fabricitem.FabricItemProperties[fabkqldatabase.Properties]) error {
		client := fabkqldatabase.NewClientFactoryWithClient(fabricClient).NewItemsClient()

		respGet, err := client.GetKQLDatabase(ctx, model.WorkspaceID.ValueString(), model.ID.ValueString(), nil)
		if err != nil {
			return err
		}

		fabricItem.Set(respGet.KQLDatabase)

		return nil
	}

	itemListGetter := func(ctx context.Context, fabricClient fabric.Client, model fabricitem.DataSourceFabricItemDefinitionPropertiesModel[kqlDatabasePropertiesModel, fabkqldatabase.Properties], errNotFound fabcore.ResponseError, fabricItem *fabricitem.FabricItemProperties[fabkqldatabase.Properties]) error {
		client := fabkqldatabase.NewClientFactoryWithClient(fabricClient).NewItemsClient()

		pager := client.NewListKQLDatabasesPager(model.WorkspaceID.ValueString(), nil)
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

	config := fabricitem.DataSourceFabricItemDefinitionProperties[kqlDatabasePropertiesModel, fabkqldatabase.Properties]{
		DataSourceFabricItemDefinition: fabricitem.DataSourceFabricItemDefinition{
			Type:   ItemType,
			Name:   ItemName,
			TFName: ItemTFName,
			MarkdownDescription: "Get a Fabric " + ItemName + ".\n\n" +
				"Use this data source to fetch a [" + ItemName + "](" + ItemDocsURL + ").\n\n" +
				ItemDocsSPNSupport,
			IsDisplayNameUnique: true,
			DefinitionFormats:   itemDefinitionFormats,
		},
		PropertiesAttributes: getDataSourceKQLDatabasePropertiesAttributes(),
		PropertiesSetter:     propertiesSetter,
		ItemGetter:           itemGetter,
		ItemListGetter:       itemListGetter,
	}

	return fabricitem.NewDataSourceFabricItemDefinitionProperties(config)
}
