// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package lakehouse

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/microsoft/fabric-sdk-go/fabric"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	fablakehouse "github.com/microsoft/fabric-sdk-go/fabric/lakehouse"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
)

func NewDataSourceLakehouse(ctx context.Context) datasource.DataSource {
	propertiesSetter := func(ctx context.Context, from *fablakehouse.Properties, to *fabricitem.DataSourceFabricItemPropertiesModel[lakehousePropertiesModel, fablakehouse.Properties]) diag.Diagnostics {
		properties := supertypes.NewSingleNestedObjectValueOfNull[lakehousePropertiesModel](ctx)

		if from != nil {
			propertiesModel := &lakehousePropertiesModel{}

			if diags := propertiesModel.set(ctx, from); diags.HasError() {
				return diags
			}

			if diags := properties.Set(ctx, propertiesModel); diags.HasError() {
				return diags
			}
		}

		to.Properties = properties

		return nil
	}

	itemGetter := func(ctx context.Context, fabricClient fabric.Client, model fabricitem.DataSourceFabricItemPropertiesModel[lakehousePropertiesModel, fablakehouse.Properties], fabricItem *fabricitem.FabricItemProperties[fablakehouse.Properties]) error {
		client := fablakehouse.NewClientFactoryWithClient(fabricClient).NewItemsClient()

		respGet, err := client.GetLakehouse(ctx, model.WorkspaceID.ValueString(), model.ID.ValueString(), nil)
		if err != nil {
			return err
		}

		fabricItem.Set(respGet.Lakehouse)

		return nil
	}

	itemListGetter := func(ctx context.Context, fabricClient fabric.Client, model fabricitem.DataSourceFabricItemPropertiesModel[lakehousePropertiesModel, fablakehouse.Properties], errNotFound fabcore.ResponseError, fabricItem *fabricitem.FabricItemProperties[fablakehouse.Properties]) error {
		client := fablakehouse.NewClientFactoryWithClient(fabricClient).NewItemsClient()

		pager := client.NewListLakehousesPager(model.WorkspaceID.ValueString(), nil)
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

	config := fabricitem.DataSourceFabricItemProperties[lakehousePropertiesModel, fablakehouse.Properties]{
		DataSourceFabricItem: fabricitem.DataSourceFabricItem{
			Type:   ItemType,
			Name:   ItemName,
			TFName: ItemTFName,
			MarkdownDescription: "Get a Fabric " + ItemName + ".\n\n" +
				"Use this data source to fetch a [" + ItemName + "](" + ItemDocsURL + ").\n\n" +
				ItemDocsSPNSupport,
			IsDisplayNameUnique: true,
		},
		PropertiesAttributes: getDataSourceLakehousePropertiesAttributes(ctx),
		PropertiesSetter:     propertiesSetter,
		ItemGetter:           itemGetter,
		ItemListGetter:       itemListGetter,
	}

	return fabricitem.NewDataSourceFabricItemProperties(config)
}
