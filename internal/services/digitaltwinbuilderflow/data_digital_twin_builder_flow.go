// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package digitaltwinbuilderflow

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/microsoft/fabric-sdk-go/fabric"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	fabdigitaltwinbuilderflow "github.com/microsoft/fabric-sdk-go/fabric/digitaltwinbuilderflow"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
)

func NewDataSourceDigitalTwinBuilderFlow(ctx context.Context) datasource.DataSource {
	propertiesSetter := func(ctx context.Context, from *fabdigitaltwinbuilderflow.Properties, to *fabricitem.DataSourceFabricItemDefinitionPropertiesModel[digitalTwinBuilderFlowConfigPropertiesModel, fabdigitaltwinbuilderflow.Properties]) diag.Diagnostics {
		properties := supertypes.NewSingleNestedObjectValueOfNull[digitalTwinBuilderFlowConfigPropertiesModel](ctx)

		if from != nil {
			propertiesModel := &digitalTwinBuilderFlowConfigPropertiesModel{}
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

	itemGetter := func(ctx context.Context, fabricClient fabric.Client, model fabricitem.DataSourceFabricItemDefinitionPropertiesModel[digitalTwinBuilderFlowConfigPropertiesModel, fabdigitaltwinbuilderflow.Properties], fabricItem *fabricitem.FabricItemProperties[fabdigitaltwinbuilderflow.Properties]) error {
		client := fabdigitaltwinbuilderflow.NewClientFactoryWithClient(fabricClient).NewItemsClient()

		respGet, err := client.GetDigitalTwinBuilderFlow(ctx, model.WorkspaceID.ValueString(), model.ID.ValueString(), nil)
		if err != nil {
			return err
		}

		fabricItem.Set(respGet.DigitalTwinBuilderFlow)

		return nil
	}

	itemListGetter := func(ctx context.Context, fabricClient fabric.Client, model fabricitem.DataSourceFabricItemDefinitionPropertiesModel[digitalTwinBuilderFlowConfigPropertiesModel, fabdigitaltwinbuilderflow.Properties], errNotFound fabcore.ResponseError, fabricItem *fabricitem.FabricItemProperties[fabdigitaltwinbuilderflow.Properties]) error {
		client := fabdigitaltwinbuilderflow.NewClientFactoryWithClient(fabricClient).NewItemsClient()

		pager := client.NewListDigitalTwinBuilderFlowsPager(model.WorkspaceID.ValueString(), nil)
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

	config := fabricitem.DataSourceFabricItemDefinitionProperties[digitalTwinBuilderFlowConfigPropertiesModel, fabdigitaltwinbuilderflow.Properties]{
		DataSourceFabricItemDefinition: fabricitem.DataSourceFabricItemDefinition{
			TypeInfo:            ItemTypeInfo,
			FabricItemType:      FabricItemType,
			IsDisplayNameUnique: true,
			DefinitionFormats:   itemDefinitionFormats,
		},
		PropertiesAttributes: getDataSourceDigitalTwinBuilderFlowProperties(ctx),
		PropertiesSetter:     propertiesSetter,
		ItemGetter:           itemGetter,
		ItemListGetter:       itemListGetter,
	}

	return fabricitem.NewDataSourceFabricItemDefinitionProperties(config)
}
