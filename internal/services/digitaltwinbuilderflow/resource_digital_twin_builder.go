// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package digitaltwinbuilderflow

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/microsoft/fabric-sdk-go/fabric"
	fabdigitaltwinbuilderflow "github.com/microsoft/fabric-sdk-go/fabric/digitaltwinbuilderflow"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
)

func NewResourceDigitalTwinBuilderFlow(ctx context.Context) resource.Resource {
	refType := fabdigitaltwinbuilderflow.ItemReferenceTypeByID

	creationPayloadSetter := func(ctx context.Context, from digitalTwinBuilderFlowConfigPropertiesModel) (*fabdigitaltwinbuilderflow.CreationPayload, diag.Diagnostics) {
		itemRef, diags := from.DigitalTwinBuilderItemReference.Get(ctx)
		if diags.HasError() {
			return nil, diags
		}

		cp := &fabdigitaltwinbuilderflow.CreationPayload{
			DigitalTwinBuilderItemReference: &fabdigitaltwinbuilderflow.ItemReferenceByID{
				ItemID:        itemRef.ItemID.ValueStringPointer(),
				ReferenceType: &refType,
				WorkspaceID:   itemRef.WorkspaceID.ValueStringPointer(),
			},
		}

		return cp, nil
	}
	propertiesSetter := func(ctx context.Context, from *fabdigitaltwinbuilderflow.Properties, to *fabricitem.ResourceFabricItemConfigDefinitionPropertiesModel[digitalTwinBuilderFlowConfigPropertiesModel, fabdigitaltwinbuilderflow.Properties, digitalTwinBuilderFlowConfigPropertiesModel, fabdigitaltwinbuilderflow.CreationPayload]) diag.Diagnostics {
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
	itemGetter := func(ctx context.Context, fabricClient fabric.Client, model fabricitem.ResourceFabricItemConfigDefinitionPropertiesModel[digitalTwinBuilderFlowConfigPropertiesModel, fabdigitaltwinbuilderflow.Properties, digitalTwinBuilderFlowConfigPropertiesModel, fabdigitaltwinbuilderflow.CreationPayload], fabricItem *fabricitem.FabricItemProperties[fabdigitaltwinbuilderflow.Properties]) error {
		client := fabdigitaltwinbuilderflow.NewClientFactoryWithClient(fabricClient).NewItemsClient()

		respGet, err := client.GetDigitalTwinBuilderFlow(ctx, model.WorkspaceID.ValueString(), model.ID.ValueString(), nil)
		if err != nil {
			return err
		}

		fabricItem.Set(respGet.DigitalTwinBuilderFlow)

		return nil
	}

	config := fabricitem.ResourceFabricItemConfigDefinitionProperties[digitalTwinBuilderFlowConfigPropertiesModel, fabdigitaltwinbuilderflow.Properties, digitalTwinBuilderFlowConfigPropertiesModel, fabdigitaltwinbuilderflow.CreationPayload]{
		ResourceFabricItemDefinition: fabricitem.ResourceFabricItemDefinition{
			TypeInfo:              ItemTypeInfo,
			FabricItemType:        FabricItemType,
			NameRenameAllowed:     true,
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
		ConfigOrDefinitionRequired: true,
		ConfigAttributes:           getResourceDigitalTwinBuilderFlowConfigurationProperties(ctx),
		CreationPayloadSetter:      creationPayloadSetter,
		PropertiesAttributes:       getResourceDigitalTwinBuilderFlowProperties(ctx),
		PropertiesSetter:           propertiesSetter,
		ItemGetter:                 itemGetter,
	}

	return fabricitem.NewResourceFabricItemConfigDefinitionProperties(config)
}
