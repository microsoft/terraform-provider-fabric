// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package mirroreddatabase

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/microsoft/fabric-sdk-go/fabric"
	"github.com/microsoft/fabric-sdk-go/fabric/mirroreddatabase"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
)

func NewResourceMirroredDatabase() resource.Resource {
	propertiesSetter := func(ctx context.Context, from *mirroreddatabase.Properties, to *fabricitem.ResourceFabricItemDefinitionPropertiesModel[mirroredDatabasePropertiesModel, mirroreddatabase.Properties]) diag.Diagnostics {
		properties := supertypes.NewSingleNestedObjectValueOfNull[mirroredDatabasePropertiesModel](ctx)

		if from == nil {
			return nil
		}

		propertiesModel := &mirroredDatabasePropertiesModel{}
		if diags := propertiesModel.set(ctx, *from); diags.HasError() {
			return diags
		}

		diags := properties.Set(ctx, propertiesModel)
		if diags.HasError() {
			return diags
		}

		to.Properties = properties
		return nil
	}

	itemGetter := func(ctx context.Context, client fabric.Client, model fabricitem.ResourceFabricItemDefinitionPropertiesModel[mirroredDatabasePropertiesModel, mirroreddatabase.Properties], fabricItem *fabricitem.FabricItemProperties[mirroreddatabase.Properties]) error {
		mirroredDbClient := mirroreddatabase.NewClientFactoryWithClient(client).NewItemsClient()

		resp, err := mirroredDbClient.GetMirroredDatabase(ctx, model.WorkspaceID.ValueString(), model.ID.ValueString(), nil)
		if err != nil {
			return err
		}

		fabricItem.Set(resp.MirroredDatabase)
		return nil
	}

	config := fabricitem.ResourceFabricItemDefinitionProperties[mirroredDatabasePropertiesModel, mirroreddatabase.Properties]{
		ResourceFabricItemDefinition: fabricitem.ResourceFabricItemDefinition{
			Type:              ItemType,
			Name:              ItemName,
			NameRenameAllowed: true,
			TFName:            ItemTFName,
			MarkdownDescription: "Manages a Fabric " + ItemName + ".\n\n" +
				"Use this resource to create and manage a [" + ItemName + "](" + ItemDocsURL + ").\n\n" +
				ItemDocsSPNSupport,
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
		PropertiesAttributes: getResourceMirroredDatabasePropertiesAttributes(),
		PropertiesSetter:     propertiesSetter,
		ItemGetter:           itemGetter,
	}

	return fabricitem.NewResourceFabricItemDefinitionProperties(config)
}
