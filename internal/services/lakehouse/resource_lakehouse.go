// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package lakehouse

import (
	"context"
	"net/http"
	"time"

	azto "github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoft/fabric-sdk-go/fabric"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	fablakehouse "github.com/microsoft/fabric-sdk-go/fabric/lakehouse"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
)

func NewResourceLakehouse(ctx context.Context) resource.Resource {
	creationPayloadSetter := func(_ context.Context, from lakehouseConfigurationModel) (*fablakehouse.CreationPayload, diag.Diagnostics) {
		if from.EnableSchemas.ValueBool() {
			cp := &fablakehouse.CreationPayload{
				EnableSchemas: from.EnableSchemas.ValueBoolPointer(),
			}

			return cp, nil
		}

		return nil, nil
	}

	propertiesSetter := func(ctx context.Context, from *fablakehouse.Properties, to *fabricitem.ResourceFabricItemConfigPropertiesModel[lakehousePropertiesModel, fablakehouse.Properties, lakehouseConfigurationModel, fablakehouse.CreationPayload]) diag.Diagnostics {
		properties := supertypes.NewSingleNestedObjectValueOfNull[lakehousePropertiesModel](ctx)

		if from != nil {
			propertiesModel := &lakehousePropertiesModel{}

			if diags := propertiesModel.set(ctx, from); diags.HasError() {
				return diags
			}

			if diags := properties.Set(ctx, propertiesModel); diags.HasError() {
				return diags
			}
		}

		to.Properties = properties

		return nil
	}

	itemGetter := func(ctx context.Context, fabricClient fabric.Client, model fabricitem.ResourceFabricItemConfigPropertiesModel[lakehousePropertiesModel, fablakehouse.Properties, lakehouseConfigurationModel, fablakehouse.CreationPayload], fabricItem *fabricitem.FabricItemProperties[fablakehouse.Properties]) error {
		client := fablakehouse.NewClientFactoryWithClient(fabricClient).NewItemsClient()

		for {
			respGet, err := client.GetLakehouse(ctx, model.WorkspaceID.ValueString(), model.ID.ValueString(), nil)
			if err != nil {
				return err
			}

			if respGet.Properties == nil || respGet.Properties.SQLEndpointProperties == nil {
				tflog.Info(ctx, "Lakehouse provisioning not done, waiting 30 seconds before retrying")
				time.Sleep(30 * time.Second) // lintignore:R018

				continue
			}

			switch *respGet.Properties.SQLEndpointProperties.ProvisioningStatus {
			case fablakehouse.SQLEndpointProvisioningStatusFailed:
				return &fabcore.ResponseError{
					ErrorCode:  (string)(fablakehouse.SQLEndpointProvisioningStatusFailed),
					StatusCode: http.StatusBadRequest,
					ErrorResponse: &fabcore.ErrorResponse{
						ErrorCode: azto.Ptr((string)(fablakehouse.SQLEndpointProvisioningStatusFailed)),
						Message:   azto.Ptr("Lakehouse SQL endpoint provisioning failed"),
					},
				}

			case fablakehouse.SQLEndpointProvisioningStatusSuccess:
				fabricItem.Set(respGet.Lakehouse)

				return nil
			default:
				tflog.Info(ctx, "Lakehouse provisioning in progress, waiting 30 seconds before retrying")
				time.Sleep(30 * time.Second) // lintignore:R018
			}
		}
	}

	config := fabricitem.ResourceFabricItemConfigProperties[lakehousePropertiesModel, fablakehouse.Properties, lakehouseConfigurationModel, fablakehouse.CreationPayload]{
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
		},
		IsConfigRequired:      false,
		ConfigAttributes:      getResourceLakehouseConfigurationAttributes(),
		CreationPayloadSetter: creationPayloadSetter,
		PropertiesAttributes:  getResourceLakehousePropertiesAttributes(ctx),
		PropertiesSetter:      propertiesSetter,
		ItemGetter:            itemGetter,
	}

	return fabricitem.NewResourceFabricItemConfigProperties(config)
}
