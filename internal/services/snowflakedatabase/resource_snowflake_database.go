// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package snowflakedatabase

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/microsoft/fabric-sdk-go/fabric"
	fabsnowflakedatabase "github.com/microsoft/fabric-sdk-go/fabric/snowflakedatabase"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
)

func NewResourceSnowflakeDatabase(ctx context.Context) resource.Resource {
	creationPayloadSetter := func(_ context.Context, from snowflakeDatabaseConfigurationModel) (*fabsnowflakedatabase.CreationPayload, diag.Diagnostics) {
		creationPayload := &fabsnowflakedatabase.CreationPayload{
			ConnectionID:          from.ConnectionID.ValueStringPointer(),
			SnowflakeDatabaseName: from.SnowflakeDatabaseName.ValueStringPointer(),
		}

		return creationPayload, nil
	}

	propertiesSetter := func(ctx context.Context, from *fabsnowflakedatabase.Properties, to *fabricitem.ResourceFabricItemConfigDefinitionPropertiesModel[snowflakeDatabasePropertiesModel, fabsnowflakedatabase.Properties, snowflakeDatabaseConfigurationModel, fabsnowflakedatabase.CreationPayload]) diag.Diagnostics {
		properties := supertypes.NewSingleNestedObjectValueOfNull[snowflakeDatabasePropertiesModel](ctx)

		if from != nil {
			propertiesModel := &snowflakeDatabasePropertiesModel{}

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

	itemGetter := func(ctx context.Context, fabricClient fabric.Client, model fabricitem.ResourceFabricItemConfigDefinitionPropertiesModel[snowflakeDatabasePropertiesModel, fabsnowflakedatabase.Properties, snowflakeDatabaseConfigurationModel, fabsnowflakedatabase.CreationPayload], fabricItem *fabricitem.FabricItemProperties[fabsnowflakedatabase.Properties]) error {
		client := fabsnowflakedatabase.NewClientFactoryWithClient(fabricClient).NewItemsClient()

		respGet, err := client.GetSnowflakeDatabase(ctx, model.WorkspaceID.ValueString(), model.ID.ValueString(), nil)
		if err != nil {
			return err
		}

		fabricItem.Set(respGet.SnowflakeDatabase)

		return nil
	}

	config := fabricitem.ResourceFabricItemConfigDefinitionProperties[snowflakeDatabasePropertiesModel, fabsnowflakedatabase.Properties, snowflakeDatabaseConfigurationModel, fabsnowflakedatabase.CreationPayload]{
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
		ConfigOrDefinitionRequired: false,
		ConfigAttributes:           getResourceSnowflakeDatabaseConfigurationAttributes(),
		CreationPayloadSetter:      creationPayloadSetter,
		PropertiesAttributes:       getResourceSnowflakeDatabasePropertiesAttributes(ctx),
		PropertiesSetter:           propertiesSetter,
		ItemGetter:                 itemGetter,
	}

	return fabricitem.NewResourceFabricItemConfigDefinitionProperties(config)
}
