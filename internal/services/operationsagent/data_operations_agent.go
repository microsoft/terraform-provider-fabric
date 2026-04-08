// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package operationsagent

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/microsoft/fabric-sdk-go/fabric"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	faboperationsagent "github.com/microsoft/fabric-sdk-go/fabric/operationsagent"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
)

func NewDataSourceOperationsAgent() datasource.DataSource {
	propertiesSetter := func(ctx context.Context, from *faboperationsagent.Properties, to *fabricitem.DataSourceFabricItemDefinitionPropertiesModel[operationsAgentPropertiesModel, faboperationsagent.Properties]) diag.Diagnostics {
		properties := supertypes.NewSingleNestedObjectValueOfNull[operationsAgentPropertiesModel](ctx)

		if from != nil {
			propertiesModel := &operationsAgentPropertiesModel{}
			propertiesModel.set(*from)

			if diags := properties.Set(ctx, propertiesModel); diags.HasError() {
				return diags
			}
		}

		to.Properties = properties

		return nil
	}

	itemGetter := func(ctx context.Context, fabricClient fabric.Client, model fabricitem.DataSourceFabricItemDefinitionPropertiesModel[operationsAgentPropertiesModel, faboperationsagent.Properties], fabricItem *fabricitem.FabricItemProperties[faboperationsagent.Properties]) error {
		client := faboperationsagent.NewClientFactoryWithClient(fabricClient).NewItemsClient()

		respGet, err := client.GetOperationsAgent(ctx, model.WorkspaceID.ValueString(), model.ID.ValueString(), nil)
		if err != nil {
			return err
		}

		fabricItem.Set(respGet.OperationsAgent)

		return nil
	}

	itemListGetter := func(ctx context.Context, fabricClient fabric.Client, model fabricitem.DataSourceFabricItemDefinitionPropertiesModel[operationsAgentPropertiesModel, faboperationsagent.Properties], errNotFound fabcore.ResponseError, fabricItem *fabricitem.FabricItemProperties[faboperationsagent.Properties]) error {
		client := faboperationsagent.NewClientFactoryWithClient(fabricClient).NewItemsClient()

		pager := client.NewListOperationsAgentsPager(model.WorkspaceID.ValueString(), nil)
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

	config := fabricitem.DataSourceFabricItemDefinitionProperties[operationsAgentPropertiesModel, faboperationsagent.Properties]{
		DataSourceFabricItemDefinition: fabricitem.DataSourceFabricItemDefinition{
			TypeInfo:            ItemTypeInfo,
			FabricItemType:      FabricItemType,
			IsDisplayNameUnique: true,
			DefinitionFormats:   itemDefinitionFormats,
		},
		PropertiesAttributes: getDataSourceOperationsAgentPropertiesAttributes(),
		PropertiesSetter:     propertiesSetter,
		ItemGetter:           itemGetter,
		ItemListGetter:       itemListGetter,
	}

	return fabricitem.NewDataSourceFabricItemDefinitionProperties(config)
}
