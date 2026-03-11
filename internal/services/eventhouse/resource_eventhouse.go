// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package eventhouse

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/microsoft/fabric-sdk-go/fabric"
	fabeventhouse "github.com/microsoft/fabric-sdk-go/fabric/eventhouse"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
)

func NewResourceEventhouse(ctx context.Context) resource.Resource {
	creationPayloadSetter := func(_ context.Context, from eventhouseConfigurationModel) (*fabeventhouse.CreationPayload, diag.Diagnostics) {
		creationPayload := fabeventhouse.CreationPayload{}

		if !from.MinimumConsumptionUnits.IsNull() && !from.MinimumConsumptionUnits.IsUnknown() {
			creationPayload.MinimumConsumptionUnits = from.MinimumConsumptionUnits.ValueFloat64Pointer()
		}

		return &creationPayload, nil
	}

	propertiesSetter := func(ctx context.Context, from *fabeventhouse.Properties, to *fabricitem.ResourceFabricItemConfigDefinitionPropertiesModel[eventhousePropertiesModel, fabeventhouse.Properties, eventhouseConfigurationModel, fabeventhouse.CreationPayload]) diag.Diagnostics {
		properties := supertypes.NewSingleNestedObjectValueOfNull[eventhousePropertiesModel](ctx)

		if from != nil {
			propertiesModel := &eventhousePropertiesModel{}
			propertiesModel.set(ctx, *from)

			if diags := properties.Set(ctx, propertiesModel); diags.HasError() {
				return diags
			}
		}

		to.Properties = properties

		return nil
	}

	itemGetter := func(ctx context.Context, fabricClient fabric.Client, model fabricitem.ResourceFabricItemConfigDefinitionPropertiesModel[eventhousePropertiesModel, fabeventhouse.Properties, eventhouseConfigurationModel, fabeventhouse.CreationPayload], fabricItem *fabricitem.FabricItemProperties[fabeventhouse.Properties]) error {
		client := fabeventhouse.NewClientFactoryWithClient(fabricClient).NewItemsClient()

		respGet, err := client.GetEventhouse(ctx, model.WorkspaceID.ValueString(), model.ID.ValueString(), nil)
		if err != nil {
			return err
		}

		fabricItem.Set(respGet.Eventhouse)

		return nil
	}

	config := fabricitem.ResourceFabricItemConfigDefinitionProperties[eventhousePropertiesModel, fabeventhouse.Properties, eventhouseConfigurationModel, fabeventhouse.CreationPayload]{
		ResourceFabricItemDefinition: fabricitem.ResourceFabricItemDefinition{
			TypeInfo:       ItemTypeInfo,
			FabricItemType: FabricItemType, NameRenameAllowed: true,
			DisplayNameMaxLength:  123,
			DescriptionMaxLength:  256,
			DefinitionPathDocsURL: ItemDefinitionPathDocsURL,
			DefinitionPathKeysValidator: []validator.Map{
				mapvalidator.SizeAtMost(1),
				mapvalidator.KeysAre(fabricitem.DefinitionPathKeysValidator(itemDefinitionFormats)...),
			},
			DefinitionRequired: false,
			DefinitionEmpty:    ItemDefinitionEmpty,
			DefinitionFormats:  itemDefinitionFormats,
		},
		ConfigRequired:             false,
		ConfigOrDefinitionRequired: false,
		ConfigAttributes:           getResourceEventhouseConfigurationAttributes(),
		CreationPayloadSetter:      creationPayloadSetter,
		PropertiesAttributes:       getResourceEventhousePropertiesAttributes(ctx),
		PropertiesSetter:           propertiesSetter,
		ItemGetter:                 itemGetter,
	}

	return fabricitem.NewResourceFabricItemConfigDefinitionProperties(config)
}
