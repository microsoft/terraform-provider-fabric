// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package mirroreddatabase

import (
	"context"
	"net/http"
	"time"

	azto "github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoft/fabric-sdk-go/fabric"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	fabmirroreddatabase "github.com/microsoft/fabric-sdk-go/fabric/mirroreddatabase"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
)

func NewResourceMirroredDatabase(ctx context.Context) resource.Resource {
	propertiesSetter := func(ctx context.Context, from *fabmirroreddatabase.Properties, to *fabricitem.ResourceFabricItemDefinitionPropertiesModel[mirroredDatabasePropertiesModel, fabmirroreddatabase.Properties]) diag.Diagnostics {
		properties := supertypes.NewSingleNestedObjectValueOfNull[mirroredDatabasePropertiesModel](ctx)

		if from != nil {
			propertiesModel := &mirroredDatabasePropertiesModel{}

			if diags := propertiesModel.set(ctx, *from); diags.HasError() {
				return diags
			}

			diags := properties.Set(ctx, propertiesModel)
			if diags.HasError() {
				return diags
			}
		}

		to.Properties = properties

		return nil
	}

	itemGetter := func(ctx context.Context, fabricClient fabric.Client, model fabricitem.ResourceFabricItemDefinitionPropertiesModel[mirroredDatabasePropertiesModel, fabmirroreddatabase.Properties], fabricItem *fabricitem.FabricItemProperties[fabmirroreddatabase.Properties]) error {
		client := fabmirroreddatabase.NewClientFactoryWithClient(fabricClient).NewItemsClient()

		for {
			respGet, err := client.GetMirroredDatabase(ctx, model.WorkspaceID.ValueString(), model.ID.ValueString(), nil)
			if err != nil {
				return err
			}

			if respGet.Properties == nil || respGet.Properties.SQLEndpointProperties == nil {
				tflog.Info(ctx, "Mirrored Database provisioning not done, waiting 30 seconds before retrying")
				time.Sleep(30 * time.Second) // lintignore:R018

				continue
			}

			switch *respGet.Properties.SQLEndpointProperties.ProvisioningStatus {
			case fabmirroreddatabase.SQLEndpointProvisioningStatusFailed:
				return &fabcore.ResponseError{
					ErrorCode:  (string)(fabmirroreddatabase.SQLEndpointProvisioningStatusFailed),
					StatusCode: http.StatusBadRequest,
					ErrorResponse: &fabcore.ErrorResponse{
						ErrorCode: azto.Ptr((string)(fabmirroreddatabase.SQLEndpointProvisioningStatusFailed)),
						Message:   azto.Ptr("Mirrored Database SQL endpoint provisioning failed"),
					},
				}

			case fabmirroreddatabase.SQLEndpointProvisioningStatusSuccess:
				fabricItem.Set(respGet.MirroredDatabase)

				return nil
			default:
				tflog.Info(ctx, "Mirrored Database provisioning in progress, waiting 30 seconds before retrying")
				time.Sleep(30 * time.Second) // lintignore:R018
			}
		}
	}

	config := fabricitem.ResourceFabricItemDefinitionProperties[mirroredDatabasePropertiesModel, fabmirroreddatabase.Properties]{
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
		PropertiesAttributes: getResourceMirroredDatabasePropertiesAttributes(ctx),
		PropertiesSetter:     propertiesSetter,
		ItemGetter:           itemGetter,
	}

	return fabricitem.NewResourceFabricItemDefinitionProperties(config)
}
