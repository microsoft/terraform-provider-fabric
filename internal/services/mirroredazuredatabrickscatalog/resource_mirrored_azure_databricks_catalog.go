// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package mirroredazuredatabrickscatalog

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/microsoft/fabric-sdk-go/fabric"
	fabmirroredazuredatabrickscatalog "github.com/microsoft/fabric-sdk-go/fabric/mirroredazuredatabrickscatalog"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
)

func NewResourceMirroredAzureDatabricksCatalog(ctx context.Context) resource.Resource {
	creationPayloadSetter := func(_ context.Context, from mirroredAzureDatabricksCatalogConfigurationModel) (*fabmirroredazuredatabrickscatalog.CreationPayload, diag.Diagnostics) {
		creationPayload := fabmirroredazuredatabrickscatalog.CreationPayload{}
		// TBD
		return &creationPayload, nil
	}

	propertiesSetter := func(ctx context.Context, from *fabmirroredazuredatabrickscatalog.Properties, to *fabricitem.ResourceFabricItemConfigDefinitionPropertiesModel[mirroredAzureDatabricksCatalogPropertiesModel, fabmirroredazuredatabrickscatalog.Properties, mirroredAzureDatabricksCatalogConfigurationModel, fabmirroredazuredatabrickscatalog.CreationPayload]) diag.Diagnostics {
		properties := supertypes.NewSingleNestedObjectValueOfNull[mirroredAzureDatabricksCatalogPropertiesModel](ctx)

		if from != nil {
			propertiesModel := &mirroredAzureDatabricksCatalogPropertiesModel{}

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

	itemGetter := func(ctx context.Context, fabricClient fabric.Client, model fabricitem.ResourceFabricItemConfigDefinitionPropertiesModel[mirroredAzureDatabricksCatalogPropertiesModel, fabmirroredazuredatabrickscatalog.Properties, mirroredAzureDatabricksCatalogConfigurationModel, fabmirroredazuredatabrickscatalog.CreationPayload], fabricItem *fabricitem.FabricItemProperties[fabmirroredazuredatabrickscatalog.Properties]) error {
		client := fabmirroredazuredatabrickscatalog.NewClientFactoryWithClient(fabricClient).NewItemsClient()

		respGet, err := client.GetMirroredAzureDatabricksCatalog(ctx, model.WorkspaceID.ValueString(), model.ID.ValueString(), nil)
		if err != nil {
			return err
		}

		fabricItem.Set(respGet.MirroredAzureDatabricksCatalog)

		return nil
	}

	config := fabricitem.ResourceFabricItemConfigDefinitionProperties[mirroredAzureDatabricksCatalogPropertiesModel, fabmirroredazuredatabrickscatalog.Properties, mirroredAzureDatabricksCatalogConfigurationModel, fabmirroredazuredatabrickscatalog.CreationPayload]{
		ResourceFabricItemDefinition: fabricitem.ResourceFabricItemDefinition{
			TypeInfo:              ItemTypeInfo,
			FabricItemType:        FabricItemType,
			NameRenameAllowed:     true,
			DisplayNameMaxLength:  123,
			DescriptionMaxLength:  256,
			DefinitionPathDocsURL: ItemDefinitionPathDocsURL,
			DefinitionPathKeysValidator: []validator.Map{
				mapvalidator.SizeAtMost(len(itemDefinitionFormats)),
				mapvalidator.KeysAre(stringvalidator.OneOf(fabricitem.GetDefinitionFormatPaths(itemDefinitionFormats, fabricitem.DefinitionFormatDefault)...)),
			},
			DefinitionRequired: false,
			DefinitionEmpty:    ItemDefinitionEmpty,
			DefinitionFormats:  itemDefinitionFormats,
		},
		ConfigRequired:             false,
		ConfigOrDefinitionRequired: true,
		ConfigAttributes:           getResourceMirroredAzureDatabricksCatalogConfigurationAttributes(),
		CreationPayloadSetter:      creationPayloadSetter,
		PropertiesAttributes:       getResourceMirroredAzureDatabricksCatalogPropertiesAttributes(ctx),
		PropertiesSetter:           propertiesSetter,
		ItemGetter:                 itemGetter,
	}

	return fabricitem.NewResourceFabricItemConfigDefinitionProperties(config)
}
