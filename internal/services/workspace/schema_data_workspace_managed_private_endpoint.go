// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package workspace

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
)

func getDataSourceWorkspaceManagedPrivateEndpointAttributes(ctx context.Context, name string) map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			MarkdownDescription: fmt.Sprintf("The %s ID.", name),
			Required:            true,
			CustomType:          customtypes.UUIDType{},
		},
		"name": schema.StringAttribute{
			MarkdownDescription: fmt.Sprintf("The %s name.", name),
			Computed:            true,
		},
		"provisioning_state": schema.StringAttribute{
			MarkdownDescription: "Provisioning state of endpoint. Possible values: " + utils.ConvertStringSlicesToString(fabcore.PossiblePrivateEndpointProvisioningStateValues(), true, true),
			Computed:            true,
		},
		"target_private_link_resource_id": schema.StringAttribute{
			MarkdownDescription: "Resource Id of data source for which private endpoint is created.",
			Computed:            true,
		},
		"target_subresource_type": schema.StringAttribute{
			MarkdownDescription: "Sub-resource pointing to [Private-link resource](https://learn.microsoft.com/azure/private-link/private-endpoint-overview#private-link-resource).",
			Computed:            true,
		},
		"connection_state": schema.SingleNestedAttribute{
			MarkdownDescription: "Endpoint connection state of provisioned endpoints.",
			Computed:            true,
			CustomType:          supertypes.NewSingleNestedObjectTypeOf[connectionStateModel](ctx),
			Attributes: map[string]schema.Attribute{
				"actions_required": schema.StringAttribute{
					MarkdownDescription: "Actions required to establish connection.",
					Computed:            true,
				},
				"status": schema.StringAttribute{
					MarkdownDescription: "Connection status. Possible values: " + utils.ConvertStringSlicesToString(fabcore.PossibleConnectionStatusValues(), true, true),
					Computed:            true,
				},
				"description": schema.StringAttribute{
					MarkdownDescription: "Description message provided on approving or rejecting the end point.",
					Computed:            true,
				},
			},
		},
	}
}
