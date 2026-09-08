// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package orgapp

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/microsoft/fabric-sdk-go/fabric"
	faborgapp "github.com/microsoft/fabric-sdk-go/fabric/orgapp"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	fwvalidators "github.com/microsoft/terraform-provider-fabric/internal/framework/validators"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
)

func NewResourceOrgApp() resource.Resource {
	propertiesSetter := func(ctx context.Context, from *faborgapp.Properties, to *fabricitem.ResourceFabricItemDefinitionPropertiesModel[orgAppPropertiesModel, faborgapp.Properties]) diag.Diagnostics {
		properties := supertypes.NewSingleNestedObjectValueOfNull[orgAppPropertiesModel](ctx)

		if from != nil {
			propertiesModel := &orgAppPropertiesModel{}
			propertiesModel.set(*from)

			if diags := properties.Set(ctx, propertiesModel); diags.HasError() {
				return diags
			}
		}

		to.Properties = properties

		return nil
	}

	itemGetter := func(ctx context.Context, fabricClient fabric.Client, model fabricitem.ResourceFabricItemDefinitionPropertiesModel[orgAppPropertiesModel, faborgapp.Properties], fabricItem *fabricitem.FabricItemProperties[faborgapp.Properties]) error {
		client := faborgapp.NewClientFactoryWithClient(fabricClient).NewItemsClient()

		respGet, err := client.GetOrgApp(ctx, model.WorkspaceID.ValueString(), model.ID.ValueString(), nil)
		if err != nil {
			return err
		}

		fabricItem.Set(respGet.OrgApp)

		return nil
	}

	config := fabricitem.ResourceFabricItemDefinitionProperties[orgAppPropertiesModel, faborgapp.Properties]{
		TypeInfo:              ItemTypeInfo,
		FabricItemType:        FabricItemType,
		NameRenameAllowed:     true,
		DisplayNameMaxLength:  123,
		DescriptionMaxLength:  256,
		DefinitionPathDocsURL: ItemDefinitionPathDocsURL,
		DefinitionPathKeysValidator: []validator.Map{
			mapvalidator.SizeBetween(1, 2),
			mapvalidator.KeysAre(fabricitem.DefinitionPathKeysValidator(itemDefinitionFormats)...),
			fwvalidators.RequiredMapKeys("definition.json"),
		},
		DefinitionRequired:   false,
		DefinitionEmpty:      ItemDefinitionEmpty,
		DefinitionFormats:    itemDefinitionFormats,
		PropertiesAttributes: getResourceOrgAppPropertiesAttributes(),
		PropertiesSetter:     propertiesSetter,
		ItemGetter:           itemGetter,
	}

	return fabricitem.NewResourceFabricItemDefinitionProperties(config)
}
