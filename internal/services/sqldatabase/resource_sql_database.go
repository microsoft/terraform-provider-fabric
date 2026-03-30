// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package sqldatabase

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/microsoft/fabric-sdk-go/fabric"
	fabsqldatabase "github.com/microsoft/fabric-sdk-go/fabric/sqldatabase"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	fwvalidators "github.com/microsoft/terraform-provider-fabric/internal/framework/validators"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
)

func NewResourceSQLDatabase(ctx context.Context) resource.Resource {
	creationPayloadSetter := func(ctx context.Context, from sqlDatabaseConfigurationModel) (*fabsqldatabase.CreationPayloadClassification, diag.Diagnostics) {
		var reqCreate requestCreateSQLDatabasePayload
		if diags := reqCreate.set(ctx, from); diags.HasError() {
			return nil, diags
		}

		return &reqCreate.CreationPayloadClassification, nil
	}

	propertiesSetter := func(ctx context.Context, from *fabsqldatabase.Properties, to *fabricitem.ResourceFabricItemConfigDefinitionPropertiesModel[sqlDatabasePropertiesModel, fabsqldatabase.Properties, sqlDatabaseConfigurationModel, fabsqldatabase.CreationPayloadClassification]) diag.Diagnostics {
		properties := supertypes.NewSingleNestedObjectValueOfNull[sqlDatabasePropertiesModel](ctx)

		if from != nil {
			propertiesModel := &sqlDatabasePropertiesModel{}

			propertiesModel.set(*from)

			if diags := properties.Set(ctx, propertiesModel); diags.HasError() {
				return diags
			}
		}

		to.Properties = properties

		return nil
	}

	itemGetter := func(ctx context.Context, fabricClient fabric.Client, model fabricitem.ResourceFabricItemConfigDefinitionPropertiesModel[sqlDatabasePropertiesModel, fabsqldatabase.Properties, sqlDatabaseConfigurationModel, fabsqldatabase.CreationPayloadClassification], fabricItem *fabricitem.FabricItemProperties[fabsqldatabase.Properties]) error {
		client := fabsqldatabase.NewClientFactoryWithClient(fabricClient).NewItemsClient()

		respGet, err := client.GetSQLDatabase(ctx, model.WorkspaceID.ValueString(), model.ID.ValueString(), nil)
		if err != nil {
			return err
		}

		fabricItem.Set(respGet.SQLDatabase)

		return nil
	}

	config := fabricitem.ResourceFabricItemConfigDefinitionProperties[sqlDatabasePropertiesModel, fabsqldatabase.Properties, sqlDatabaseConfigurationModel, fabsqldatabase.CreationPayloadClassification]{
		ResourceFabricItemDefinition: fabricitem.ResourceFabricItemDefinition{
			TypeInfo:              ItemTypeInfo,
			FabricItemType:        FabricItemType,
			NameRenameAllowed:     true,
			DisplayNameMaxLength:  123,
			DescriptionMaxLength:  256,
			DefinitionPathDocsURL: ItemDefinitionPathDocsURL,
			DefinitionFormats:     itemDefinitionFormats,
			DefinitionPathKeysValidator: []validator.Map{
				mapvalidator.SizeAtLeast(1),
				mapvalidator.KeysAre(
					fwvalidators.PatternsIfAttributeIsOneOf(
						path.MatchRoot("format"),
						[]attr.Value{types.StringValue("dacpac")},
						fabricitem.GetDefinitionFormatPaths(itemDefinitionFormats, "dacpac"),
						"Definition path must match one of the following: "+utils.ConvertStringSlicesToString(fabricitem.GetDefinitionFormatPaths(itemDefinitionFormats, "dacpac"), true, false),
					),
					fwvalidators.PatternsIfAttributeIsOneOf(
						path.MatchRoot("format"),
						[]attr.Value{types.StringValue("sqlproj")},
						fabricitem.GetDefinitionFormatPaths(itemDefinitionFormats, "sqlproj"),
						"Definition path must match one of the following: "+utils.ConvertStringSlicesToString(fabricitem.GetDefinitionFormatPaths(itemDefinitionFormats, "sqlproj"), true, false),
					),
				),
			},
			DefinitionRequired: false,
			DefinitionEmpty:    "",
		},
		ConfigRequired:             false,
		ConfigOrDefinitionRequired: false,
		ConfigAttributes:           getResourceSQLDatabaseConfigurationAttributes(ctx),
		CreationPayloadSetter:      creationPayloadSetter,
		PropertiesAttributes:       getResourceSQLDatabasePropertiesAttributes(),
		PropertiesSetter:           propertiesSetter,
		ItemGetter:                 itemGetter,
	}

	return fabricitem.NewResourceFabricItemConfigDefinitionProperties(config)
}
