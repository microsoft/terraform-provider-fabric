// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package kqldatabase

import (
	"context"
	"fmt"

	supertypes "github.com/FrangipaneTeam/terraform-plugin-framework-supertypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/microsoft/fabric-sdk-go/fabric"
	fabkqldatabase "github.com/microsoft/fabric-sdk-go/fabric/kqldatabase"

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

	propertiesSetter := func(ctx context.Context, from *fabkqldatabase.Properties, to *fabricitem.ResourceFabricItemConfigPropertiesModel[kqlDatabasePropertiesModel, fabkqldatabase.Properties, kqlDatabaseConfigurationModel, fabkqldatabase.CreationPayloadClassification]) diag.Diagnostics {
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

	itemGetter := func(ctx context.Context, fabricClient fabric.Client, model fabricitem.ResourceFabricItemConfigPropertiesModel[kqlDatabasePropertiesModel, fabkqldatabase.Properties, kqlDatabaseConfigurationModel, fabkqldatabase.CreationPayloadClassification], fabricItem *fabricitem.FabricItemProperties[fabkqldatabase.Properties]) error {
		client := fabkqldatabase.NewClientFactoryWithClient(fabricClient).NewItemsClient()

		respGet, err := client.GetKQLDatabase(ctx, model.WorkspaceID.ValueString(), model.ID.ValueString(), nil)
		if err != nil {
			return err
		}

		fabricItem.Set(respGet.KQLDatabase)

		return nil
	}

	config := fabricitem.ResourceFabricItemConfigProperties[kqlDatabasePropertiesModel, fabkqldatabase.Properties, kqlDatabaseConfigurationModel, fabkqldatabase.CreationPayloadClassification]{
		ResourceFabricItem: fabricitem.ResourceFabricItem{
			Type:              ItemType,
			Name:              ItemName,
			NameRenameAllowed: true,
			TFName:            ItemTFName,
			MarkdownDescription: "Manage a Fabric " + ItemName + ".\n\n" +
				"Use this resource to manage a [" + ItemName + "](" + ItemDocsURL + ").\n\n" +
				ItemDocsSPNSupport,
			DisplayNameMaxLength: 123,
			DescriptionMaxLength: 256,
			// FormatTypeDefault:     ItemFormatTypeDefault,
			// FormatTypes:           ItemFormatTypes,
			// DefinitionPathDocsURL: ItemDefinitionPathDocsURL,
			// DefinitionPathKeys:    ItemDefinitionPaths,
			// DefinitionPathKeysValidator: []validator.Map{
			// 	mapvalidator.SizeAtLeast(2),
			// 	mapvalidator.SizeAtMost(2),
			// 	mapvalidator.KeysAre(stringvalidator.OneOf(ItemDefinitionPaths...)),
			// },
			// DefinitionRequired: false,
			// DefinitionEmpty:    "",
		},
		IsConfigRequired:      false,
		ConfigAttributes:      getResourceKQLDatabaseConfigurationAttributes(),
		CreationPayloadSetter: creationPayloadSetter,
		PropertiesAttributes:  getResourceKQLDatabasePropertiesAttributes(),
		PropertiesSetter:      propertiesSetter,
		ItemGetter:            itemGetter,
	}

	return fabricitem.NewResourceFabricItemConfigProperties(config)
}
