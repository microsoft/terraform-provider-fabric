// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package mirroredcatalog

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/microsoft/fabric-sdk-go/fabric"
	fabmirroredcatalog "github.com/microsoft/fabric-sdk-go/fabric/mirroredcatalog"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
)

func NewResourceMirroredCatalog(ctx context.Context) resource.Resource {
	propertiesSetter := func(ctx context.Context, from *fabmirroredcatalog.Properties, to *fabricitem.ResourceFabricItemDefinitionPropertiesModel[mirroredCatalogPropertiesModel, fabmirroredcatalog.Properties]) diag.Diagnostics {
		properties := supertypes.NewSingleNestedObjectValueOfNull[mirroredCatalogPropertiesModel](ctx)

		if from != nil {
			propertiesModel := &mirroredCatalogPropertiesModel{}

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

	itemGetter := func(ctx context.Context, fabricClient fabric.Client, model fabricitem.ResourceFabricItemDefinitionPropertiesModel[mirroredCatalogPropertiesModel, fabmirroredcatalog.Properties], fabricItem *fabricitem.FabricItemProperties[fabmirroredcatalog.Properties]) error {
		client := fabmirroredcatalog.NewClientFactoryWithClient(fabricClient).NewItemsClient()

		respGet, err := client.GetMirroredCatalog(ctx, model.WorkspaceID.ValueString(), model.ID.ValueString(), nil)
		if err != nil {
			return err
		}

		fabricItem.Set(respGet.MirroredCatalog)

		return nil
	}

	config := fabricitem.ResourceFabricItemDefinitionProperties[mirroredCatalogPropertiesModel, fabmirroredcatalog.Properties]{
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
			DefinitionRequired: true,
			DefinitionEmpty:    ItemDefinitionEmpty,
			DefinitionFormats:  itemDefinitionFormats,
		},
		PropertiesAttributes: getResourceMirroredCatalogPropertiesAttributes(ctx),
		PropertiesSetter:     propertiesSetter,
		ItemGetter:           itemGetter,
	}

	return fabricitem.NewResourceFabricItemDefinitionProperties(config)
}
