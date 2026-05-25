// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package snowflakedatabase

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/microsoft/fabric-sdk-go/fabric"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	fabsnowflakedatabase "github.com/microsoft/fabric-sdk-go/fabric/snowflakedatabase"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
)

func NewDataSourceSnowflakeDatabase(ctx context.Context) datasource.DataSource {
	propertiesSetter := func(ctx context.Context, from *fabsnowflakedatabase.Properties, to *fabricitem.DataSourceFabricItemDefinitionPropertiesModel[snowflakeDatabasePropertiesModel, fabsnowflakedatabase.Properties]) diag.Diagnostics {
		properties := supertypes.NewSingleNestedObjectValueOfNull[snowflakeDatabasePropertiesModel](ctx)

		if from != nil {
			propertiesModel := &snowflakeDatabasePropertiesModel{}

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

	itemGetter := func(ctx context.Context, fabricClient fabric.Client, model fabricitem.DataSourceFabricItemDefinitionPropertiesModel[snowflakeDatabasePropertiesModel, fabsnowflakedatabase.Properties], fabricItem *fabricitem.FabricItemProperties[fabsnowflakedatabase.Properties]) error {
		client := fabsnowflakedatabase.NewClientFactoryWithClient(fabricClient).NewItemsClient()

		respGet, err := client.GetSnowflakeDatabase(ctx, model.WorkspaceID.ValueString(), model.ID.ValueString(), nil)
		if err != nil {
			return err
		}

		fabricItem.Set(respGet.SnowflakeDatabase)

		return nil
	}

	itemListGetter := func(ctx context.Context, fabricClient fabric.Client, model fabricitem.DataSourceFabricItemDefinitionPropertiesModel[snowflakeDatabasePropertiesModel, fabsnowflakedatabase.Properties], errNotFound fabcore.ResponseError, fabricItem *fabricitem.FabricItemProperties[fabsnowflakedatabase.Properties]) error {
		client := fabsnowflakedatabase.NewClientFactoryWithClient(fabricClient).NewItemsClient()

		pager := client.NewListSnowflakeDatabasesPager(model.WorkspaceID.ValueString(), nil)
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

	config := fabricitem.DataSourceFabricItemDefinitionProperties[snowflakeDatabasePropertiesModel, fabsnowflakedatabase.Properties]{
		DataSourceFabricItemDefinition: fabricitem.DataSourceFabricItemDefinition{
			TypeInfo:            ItemTypeInfo,
			FabricItemType:      FabricItemType,
			IsDisplayNameUnique: true,
			DefinitionFormats:   itemDefinitionFormats,
		},
		PropertiesAttributes: getDataSourceSnowflakeDatabasePropertiesAttributes(ctx),
		PropertiesSetter:     propertiesSetter,
		ItemGetter:           itemGetter,
		ItemListGetter:       itemListGetter,
	}

	return fabricitem.NewDataSourceFabricItemDefinitionProperties(config)
}
