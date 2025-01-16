// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package sparkjobdefinition

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/microsoft/fabric-sdk-go/fabric"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	fabsparkjobdefinition "github.com/microsoft/fabric-sdk-go/fabric/sparkjobdefinition"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
)

func NewDataSourceSparkJobDefinition() datasource.DataSource {
	propertiesSetter := func(ctx context.Context, from *fabsparkjobdefinition.Properties, to *fabricitem.DataSourceFabricItemDefinitionPropertiesModel[sparkJobDefinitionPropertiesModel, fabsparkjobdefinition.Properties]) diag.Diagnostics {
		properties := supertypes.NewSingleNestedObjectValueOfNull[sparkJobDefinitionPropertiesModel](ctx)

		if from != nil {
			propertiesModel := &sparkJobDefinitionPropertiesModel{}
			propertiesModel.set(from)

			if diags := properties.Set(ctx, propertiesModel); diags.HasError() {
				return diags
			}
		}

		to.Properties = properties

		return nil
	}

	itemGetter := func(ctx context.Context, fabricClient fabric.Client, model fabricitem.DataSourceFabricItemDefinitionPropertiesModel[sparkJobDefinitionPropertiesModel, fabsparkjobdefinition.Properties], fabricItem *fabricitem.FabricItemProperties[fabsparkjobdefinition.Properties]) error {
		client := fabsparkjobdefinition.NewClientFactoryWithClient(fabricClient).NewItemsClient()

		respGet, err := client.GetSparkJobDefinition(ctx, model.WorkspaceID.ValueString(), model.ID.ValueString(), nil)
		if err != nil {
			return err
		}

		fabricItem.Set(respGet.SparkJobDefinition)

		return nil
	}

	itemListGetter := func(ctx context.Context, fabricClient fabric.Client, model fabricitem.DataSourceFabricItemDefinitionPropertiesModel[sparkJobDefinitionPropertiesModel, fabsparkjobdefinition.Properties], errNotFound fabcore.ResponseError, fabricItem *fabricitem.FabricItemProperties[fabsparkjobdefinition.Properties]) error {
		client := fabsparkjobdefinition.NewClientFactoryWithClient(fabricClient).NewItemsClient()

		pager := client.NewListSparkJobDefinitionsPager(model.WorkspaceID.ValueString(), nil)
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

	config := fabricitem.DataSourceFabricItemDefinitionProperties[sparkJobDefinitionPropertiesModel, fabsparkjobdefinition.Properties]{
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
		PropertiesAttributes: getDataSourceSparkJobDefinitionPropertiesAttributes(),
		PropertiesSetter:     propertiesSetter,
		ItemGetter:           itemGetter,
		ItemListGetter:       itemListGetter,
	}

	return fabricitem.NewDataSourceFabricItemDefinitionProperties(config)
}
