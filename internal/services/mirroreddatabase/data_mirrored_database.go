// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package mirroreddatabase

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/microsoft/fabric-sdk-go/fabric"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	"github.com/microsoft/fabric-sdk-go/fabric/mirroreddatabase"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
)

func NewDataSourceMirroredDatabase() datasource.DataSource {
	propertiesSetter := func(ctx context.Context, from *mirroreddatabase.Properties, to *fabricitem.DataSourceFabricItemDefinitionPropertiesModel[mirroredDatabasePropertiesModel, mirroreddatabase.Properties]) diag.Diagnostics {
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

	// itemGetter retrieves a single mirrored database item.
	itemGetter := func(ctx context.Context, fabricClient fabric.Client, model fabricitem.DataSourceFabricItemDefinitionPropertiesModel[mirroredDatabasePropertiesModel, mirroreddatabase.Properties], fabricItem *fabricitem.FabricItemProperties[mirroreddatabase.Properties]) error {
		client := mirroreddatabase.NewClientFactoryWithClient(fabricClient).NewItemsClient()
		respGet, err := client.GetMirroredDatabase(ctx, model.WorkspaceID.ValueString(), model.ID.ValueString(), nil)
		if err != nil {
			return err
		}
		fabricItem.Set(respGet.MirroredDatabase)
		return nil
	}

	// itemListGetter searches for a mirrored database by its display name.
	itemListGetter := func(ctx context.Context, fabricClient fabric.Client, model fabricitem.DataSourceFabricItemDefinitionPropertiesModel[mirroredDatabasePropertiesModel, mirroreddatabase.Properties], errNotFound fabcore.ResponseError, fabricItem *fabricitem.FabricItemProperties[mirroreddatabase.Properties]) error {
		client := mirroreddatabase.NewClientFactoryWithClient(fabricClient).NewItemsClient()
		pager := client.NewListMirroredDatabasesPager(model.WorkspaceID.ValueString(), nil)
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

	config := fabricitem.DataSourceFabricItemDefinitionProperties[mirroredDatabasePropertiesModel, mirroreddatabase.Properties]{
		DataSourceFabricItemDefinition: fabricitem.DataSourceFabricItemDefinition{
			Type:   "MirroredDatabase",
			Name:   "Mirrored Database",
			TFName: "mirrored_database",
			MarkdownDescription: "Get a Fabric Mirrored Database.\n\n" +
				"Use this data source to fetch a Mirrored Database.\n\n",
			IsDisplayNameUnique: true,
			DefinitionFormats:   itemDefinitionFormats,
		},
		PropertiesAttributes: getDataSourceMirroredDatabasePropertiesAttributes(), // define this function to return schema attributes
		PropertiesSetter:     propertiesSetter,
		ItemGetter:           itemGetter,
		ItemListGetter:       itemListGetter,
	}

	return fabricitem.NewDataSourceFabricItemDefinitionProperties(config)
}
