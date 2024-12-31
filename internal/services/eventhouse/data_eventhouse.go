// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package eventhouse

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/microsoft/fabric-sdk-go/fabric"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	fabeventhouse "github.com/microsoft/fabric-sdk-go/fabric/eventhouse"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
)

func NewDataSourceEventhouse(ctx context.Context) datasource.DataSource {
	propertiesSetter := func(ctx context.Context, from *fabeventhouse.Properties, to *fabricitem.DataSourceFabricItemDefinitionPropertiesModel[eventhousePropertiesModel, fabeventhouse.Properties]) diag.Diagnostics {
		properties := supertypes.NewSingleNestedObjectValueOfNull[eventhousePropertiesModel](ctx)

		if from != nil {
			propertiesModel := &eventhousePropertiesModel{}
			propertiesModel.set(ctx, from)

			if diags := properties.Set(ctx, propertiesModel); diags.HasError() {
				return diags
			}
		}

		to.Properties = properties

		return nil
	}

	itemGetter := func(ctx context.Context, fabricClient fabric.Client, model fabricitem.DataSourceFabricItemDefinitionPropertiesModel[eventhousePropertiesModel, fabeventhouse.Properties], fabricItem *fabricitem.FabricItemProperties[fabeventhouse.Properties]) error {
		client := fabeventhouse.NewClientFactoryWithClient(fabricClient).NewItemsClient()

		respGet, err := client.GetEventhouse(ctx, model.WorkspaceID.ValueString(), model.ID.ValueString(), nil)
		if err != nil {
			return err
		}

		fabricItem.Set(respGet.Eventhouse)

		return nil
	}

	itemListGetter := func(ctx context.Context, fabricClient fabric.Client, model fabricitem.DataSourceFabricItemDefinitionPropertiesModel[eventhousePropertiesModel, fabeventhouse.Properties], errNotFound fabcore.ResponseError, fabricItem *fabricitem.FabricItemProperties[fabeventhouse.Properties]) error {
		client := fabeventhouse.NewClientFactoryWithClient(fabricClient).NewItemsClient()

		pager := client.NewListEventhousesPager(model.WorkspaceID.ValueString(), nil)
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

	config := fabricitem.DataSourceFabricItemDefinitionProperties[eventhousePropertiesModel, fabeventhouse.Properties]{
		DataSourceFabricItemDefinition: fabricitem.DataSourceFabricItemDefinition{
			Type:   ItemType,
			Name:   ItemName,
			TFName: ItemTFName,
			MarkdownDescription: "Get a Fabric " + ItemName + ".\n\n" +
				"Use this data source to fetch an [" + ItemName + "](" + ItemDocsURL + ").\n\n" +
				ItemDocsSPNSupport,
			IsDisplayNameUnique: true,
			FormatTypeDefault:   ItemFormatTypeDefault,
			FormatTypes:         ItemFormatTypes,
			DefinitionPathKeys:  ItemDefinitionPaths,
		},
		PropertiesAttributes: getDataSourceEventhousePropertiesAttributes(ctx),
		PropertiesSetter:     propertiesSetter,
		ItemGetter:           itemGetter,
		ItemListGetter:       itemListGetter,
	}

	return fabricitem.NewDataSourceFabricItemDefinitionProperties(config)
}
