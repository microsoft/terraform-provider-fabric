// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package operationsagent

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/microsoft/fabric-sdk-go/fabric"
	faboperationsagent "github.com/microsoft/fabric-sdk-go/fabric/operationsagent"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
)

func NewResourceOperationsAgent() resource.Resource {
	propertiesSetter := func(ctx context.Context, from *faboperationsagent.Properties, to *fabricitem.ResourceFabricItemDefinitionPropertiesModel[operationsAgentPropertiesModel, faboperationsagent.Properties]) diag.Diagnostics {
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

	itemGetter := func(ctx context.Context, fabricClient fabric.Client, model fabricitem.ResourceFabricItemDefinitionPropertiesModel[operationsAgentPropertiesModel, faboperationsagent.Properties], fabricItem *fabricitem.FabricItemProperties[faboperationsagent.Properties]) error {
		client := faboperationsagent.NewClientFactoryWithClient(fabricClient).NewItemsClient()

		respGet, err := client.GetOperationsAgent(ctx, model.WorkspaceID.ValueString(), model.ID.ValueString(), nil)
		if err != nil {
			return err
		}

		fabricItem.Set(respGet.OperationsAgent)

		return nil
	}

	config := fabricitem.ResourceFabricItemDefinitionProperties[operationsAgentPropertiesModel, faboperationsagent.Properties]{
		ResourceFabricItemDefinition: fabricitem.ResourceFabricItemDefinition{
			TypeInfo:              ItemTypeInfo,
			FabricItemType:        FabricItemType,
			NameRenameAllowed:     true,
			DisplayNameMaxLength:  123,
			DescriptionMaxLength:  256,
			DefinitionPathDocsURL: ItemDefinitionPathDocsURL,
			DefinitionPathKeysValidator: []validator.Map{
				mapvalidator.SizeAtMost(len(itemDefinitionFormats)),
				mapvalidator.KeysAre(fabricitem.DefinitionPathKeysValidator(itemDefinitionFormats)...),
			},
			DefinitionRequired: false,
			DefinitionEmpty:    ItemDefinitionEmpty,
			DefinitionFormats:  itemDefinitionFormats,
		},
		PropertiesAttributes: getResourceOperationsAgentPropertiesAttributes(),
		PropertiesSetter:     propertiesSetter,
		ItemGetter:           itemGetter,
	}

	return fabricitem.NewResourceFabricItemDefinitionProperties(config)
}
