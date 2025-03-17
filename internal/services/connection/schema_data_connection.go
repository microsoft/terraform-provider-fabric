// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package connection

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
)

func getDataSourceConnectionAttributes(ctx context.Context) map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			MarkdownDescription: "The object ID of the connection.",
			Optional:            true,
			Computed:            true,
			CustomType:          customtypes.UUIDType{},
		},
		"display_name": schema.StringAttribute{
			MarkdownDescription: "The display name of the connection.",
			Optional:            true,
			Computed:            true,
		},
		"connectivity_type": schema.StringAttribute{
			MarkdownDescription: "The connectivity type of the connection. Possible values: " + utils.ConvertStringSlicesToString(fabcore.PossibleConnectivityTypeValues(), true, true),
			Computed:            true,
		},
		"privacy_level": schema.StringAttribute{
			MarkdownDescription: "The privacy level of the connection. Possible values: " + utils.ConvertStringSlicesToString(fabcore.PossiblePrivacyLevelValues(), true, true),
			Computed:            true,
		},
		"gateway_id": schema.StringAttribute{
			MarkdownDescription: "The gateway object ID of the connection.",
			Computed:            true,
			CustomType:          customtypes.UUIDType{},
		},
		"connection_details": schema.SingleNestedAttribute{
			MarkdownDescription: "The connection details of the connection.",
			Computed:            true,
			CustomType:          supertypes.NewSingleNestedObjectTypeOf[dsConnectionDetailsModel](ctx),
			Attributes: map[string]schema.Attribute{
				"path": schema.StringAttribute{
					MarkdownDescription: "The path of the connection.",
					Computed:            true,
				},
				"type": schema.StringAttribute{
					MarkdownDescription: "The type of the connection.",
					Computed:            true,
				},
			},
		},
		"credential_details": schema.SingleNestedAttribute{
			MarkdownDescription: "The credential details of the connection.",
			Computed:            true,
			CustomType:          supertypes.NewSingleNestedObjectTypeOf[dsCredentialDetailsModel](ctx),
			Attributes: map[string]schema.Attribute{
				"connection_encryption": schema.StringAttribute{
					MarkdownDescription: "The connection encryption type of the connection. Possible values: " + utils.ConvertStringSlicesToString(fabcore.PossibleConnectionEncryptionValues(), true, true),
					Computed:            true,
				},
				"single_sign_on_type": schema.StringAttribute{
					MarkdownDescription: "The single sign-on type of the connection. Possible values: " + utils.ConvertStringSlicesToString(fabcore.PossibleSingleSignOnTypeValues(), true, true),
					Computed:            true,
				},
				"skip_test_connection": schema.BoolAttribute{
					MarkdownDescription: "Whether the connection should skip the test connection during creation and update. `True` - Skip the test connection, `False` - Do not skip the test connection.",
					Computed:            true,
				},
				"credential_type": schema.StringAttribute{
					MarkdownDescription: "The credential type of the connection. Possible values: " + utils.ConvertStringSlicesToString(fabcore.PossibleCredentialTypeValues(), true, true),
					Computed:            true,
				},
			},
		},
	}
}
