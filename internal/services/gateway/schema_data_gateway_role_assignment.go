// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package gateway

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
)

func getDataSourceGatewayRoleAssignmentAttributes(ctx context.Context, isList bool) map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			MarkdownDescription: "The " + GatewayRoleAssignmentName + " ID.",
			Required:            !isList, // If it's a 'get' data source then the ID is required, otherwise it's computed for a 'list' data source.
			Computed:            isList,
			CustomType:          customtypes.UUIDType{},
		},
		"role": schema.StringAttribute{
			MarkdownDescription: "The gateway role of the principal. Possible values: " + utils.ConvertStringSlicesToString(fabcore.PossibleGatewayRoleValues(), true, true) + ".",
			Computed:            true,
		},
		"principal": schema.SingleNestedAttribute{
			MarkdownDescription: "The principal.",
			Computed:            true,
			CustomType:          supertypes.NewSingleNestedObjectTypeOf[gatewayRoleAssignmentPrincipalModel](ctx),
			Attributes: map[string]schema.Attribute{
				"id": schema.StringAttribute{
					MarkdownDescription: "The principal ID.",
					Computed:            true,
					CustomType:          customtypes.UUIDType{},
				},
				"type": schema.StringAttribute{
					MarkdownDescription: "The principal type. Possible values: " + utils.ConvertStringSlicesToString(fabcore.PossiblePrincipalTypeValues(), true, true) + ".",
					Computed:            true,
				},
			},
		},
	}
}
