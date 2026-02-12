// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package digitaltwinbuilderflow

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/microsoft/fabric-sdk-go/fabric"
	fabdigitaltwinbuilderflow "github.com/microsoft/fabric-sdk-go/fabric/digitaltwinbuilderflow"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
)

func NewDataSourceDigitalTwinBuilderFlows(ctx context.Context) datasource.DataSource {
	propertiesSetter := func(ctx context.Context, from *fabdigitaltwinbuilderflow.Properties, to *fabricitem.FabricItemPropertiesModel[digitalTwinBuilderFlowConfigPropertiesModel, fabdigitaltwinbuilderflow.Properties]) diag.Diagnostics {
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

	itemListGetter := func(ctx context.Context, fabricClient fabric.Client, model fabricitem.DataSourceFabricItemsPropertiesModel[digitalTwinBuilderFlowConfigPropertiesModel, fabdigitaltwinbuilderflow.Properties], fabricItems *[]fabricitem.FabricItemProperties[fabdigitaltwinbuilderflow.Properties]) error {
		client := fabdigitaltwinbuilderflow.NewClientFactoryWithClient(fabricClient).NewItemsClient()

		fabItems := make([]fabricitem.FabricItemProperties[fabdigitaltwinbuilderflow.Properties], 0)

		respList, err := client.ListDigitalTwinBuilderFlows(ctx, model.WorkspaceID.ValueString(), nil)
		if err != nil {
			return err
		}

		for _, entity := range respList {
			var fabricItem fabricitem.FabricItemProperties[fabdigitaltwinbuilderflow.Properties]

			fabricItem.Set(entity)

			fabItems = append(fabItems, fabricItem)
		}

		*fabricItems = fabItems

		return nil
	}

	config := fabricitem.DataSourceFabricItemsProperties[digitalTwinBuilderFlowConfigPropertiesModel, fabdigitaltwinbuilderflow.Properties]{
		DataSourceFabricItems: fabricitem.DataSourceFabricItems{
			TypeInfo:       ItemTypeInfo,
			FabricItemType: FabricItemType,
		},
		PropertiesAttributes: getDataSourceDigitalTwinBuilderFlowProperties(ctx),
		PropertiesSetter:     propertiesSetter,
		ItemListGetter:       itemListGetter,
	}

	return fabricitem.NewDataSourceFabricItemsProperties(config)
}
