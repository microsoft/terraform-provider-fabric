// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package kqldatabase

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/microsoft/fabric-sdk-go/fabric"
	fabkqldatabase "github.com/microsoft/fabric-sdk-go/fabric/kqldatabase"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
)

func NewResourceKQLDatabase() resource.Resource {
	creationPayloadSetter := func(_ context.Context, from kqlDatabaseConfigurationModel) (*fabkqldatabase.CreationPayloadClassification, diag.Diagnostics) {
		kqlDatabaseType := (fabkqldatabase.Type)(from.DatabaseType.ValueString())

		var cp fabkqldatabase.CreationPayloadClassification

		switch kqlDatabaseType {
		case fabkqldatabase.TypeReadWrite:
			creationPayload := fabkqldatabase.ReadWriteDatabaseCreationPayload{
				DatabaseType:           &kqlDatabaseType,
				ParentEventhouseItemID: from.EventhouseID.ValueStringPointer(),
			}

			cp = &creationPayload
		case fabkqldatabase.TypeShortcut:
			creationPayload := fabkqldatabase.ShortcutDatabaseCreationPayload{}
			creationPayload.DatabaseType = &kqlDatabaseType
			creationPayload.ParentEventhouseItemID = from.EventhouseID.ValueStringPointer()

			if !from.InvitationToken.IsNull() && !from.InvitationToken.IsUnknown() {
				creationPayload.InvitationToken = from.InvitationToken.ValueStringPointer()
			}

			if !from.SourceClusterURI.IsNull() && !from.SourceClusterURI.IsUnknown() {
				creationPayload.SourceClusterURI = from.SourceClusterURI.ValueStringPointer()
			}

			creationPayload.SourceDatabaseName = from.SourceDatabaseName.ValueStringPointer()

			cp = &creationPayload
		default:
			var diags diag.Diagnostics

			diags.AddError(
				"Unsupported KQL database type",
				fmt.Sprintf("The KQL database type '%s' is not supported.", string(kqlDatabaseType)),
			)

			return nil, diags
		}

		return &cp, nil
	}

	propertiesSetter := func(ctx context.Context, from *fabkqldatabase.Properties, to *fabricitem.ResourceFabricItemConfigDefinitionPropertiesModel[kqlDatabasePropertiesModel, fabkqldatabase.Properties, kqlDatabaseConfigurationModel, fabkqldatabase.CreationPayloadClassification]) diag.Diagnostics {
		properties := supertypes.NewSingleNestedObjectValueOfNull[kqlDatabasePropertiesModel](ctx)

		if from != nil {
			propertiesModel := &kqlDatabasePropertiesModel{}
			propertiesModel.set(from)

			if diags := properties.Set(ctx, propertiesModel); diags.HasError() {
				return diags
			}
		}

		to.Properties = properties

		return nil
	}

	itemGetter := func(ctx context.Context, fabricClient fabric.Client, model fabricitem.ResourceFabricItemConfigDefinitionPropertiesModel[kqlDatabasePropertiesModel, fabkqldatabase.Properties, kqlDatabaseConfigurationModel, fabkqldatabase.CreationPayloadClassification], fabricItem *fabricitem.FabricItemProperties[fabkqldatabase.Properties]) error {
		client := fabkqldatabase.NewClientFactoryWithClient(fabricClient).NewItemsClient()

		respGet, err := client.GetKQLDatabase(ctx, model.WorkspaceID.ValueString(), model.ID.ValueString(), nil)
		if err != nil {
			return err
		}

		fabricItem.Set(respGet.KQLDatabase)

		return nil
	}

	config := fabricitem.ResourceFabricItemConfigDefinitionProperties[kqlDatabasePropertiesModel, fabkqldatabase.Properties, kqlDatabaseConfigurationModel, fabkqldatabase.CreationPayloadClassification]{
		ResourceFabricItemDefinition: fabricitem.ResourceFabricItemDefinition{
			Type:              ItemType,
			Name:              ItemName,
			NameRenameAllowed: true,
			TFName:            ItemTFName,
			MarkdownDescription: "Manage a Fabric " + ItemName + ".\n\n" +
				"Use this resource to manage a [" + ItemName + "](" + ItemDocsURL + ").\n\n" +
				ItemDocsSPNSupport,
			DisplayNameMaxLength:  123,
			DescriptionMaxLength:  256,
			DefinitionPathDocsURL: ItemDefinitionPathDocsURL,
			DefinitionFormats:     itemDefinitionFormats,
			DefinitionPathKeysValidator: []validator.Map{
				mapvalidator.SizeAtLeast(2),
				mapvalidator.SizeAtMost(2),
				mapvalidator.KeysAre(fabricitem.DefinitionPathKeysValidator(itemDefinitionFormats)...),
			},
			DefinitionRequired: false,
			DefinitionEmpty:    "",
		},
		ConfigRequired:             false,
		ConfigOrDefinitionRequired: true,
		ConfigAttributes:           getResourceKQLDatabaseConfigurationAttributes(),
		CreationPayloadSetter:      creationPayloadSetter,
		PropertiesAttributes:       getResourceKQLDatabasePropertiesAttributes(),
		PropertiesSetter:           propertiesSetter,
		ItemGetter:                 itemGetter,
	}

	return fabricitem.NewResourceFabricItemConfigDefinitionProperties(config)
}
