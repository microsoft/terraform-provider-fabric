// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package connection

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
)

func ResourceConnectionSchema(ctx context.Context) schema.Schema {
	tflog.Info(ctx, "Building schema for connection resource")

	// Define attribute types for connection_details and credential_details
	connectionDetailsAttributeTypes := map[string]attr.Type{
		"type":            types.StringType,
		"creation_method": types.StringType,
		"parameters": types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{
			"name":      types.StringType,
			"data_type": types.StringType,
			"value":     types.StringType,
		}}},
	}

	connectivityTypeValues := []string{"ShareableCloud", "OnPremisesGateway", "VirtualNetworkGateway"}
	privacyLevelValues := []string{"Organizational", "Private", "Public", "None"}

	return schema.Schema{
		MarkdownDescription: "Manages a Microsoft Fabric Connection",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The Connection ID.",
				CustomType:          customtypes.UUIDType{},
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "The Connection display name.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(200),
				},
			},
			"connectivity_type": schema.StringAttribute{
				MarkdownDescription: "Connectivity type. Possible values: " + utils.ConvertStringSlicesToString(connectivityTypeValues, true, true) + ".",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf(connectivityTypeValues...),
				},
			},
			"gateway_id": schema.StringAttribute{
				MarkdownDescription: "Gateway ID. Required for OnPremisesGateway and VirtualNetworkGateway connectivity types.",
				CustomType:          customtypes.UUIDType{},
				Optional:            true,
			},
			"privacy_level": schema.StringAttribute{
				MarkdownDescription: "Privacy level. Possible values: " + utils.ConvertStringSlicesToString(privacyLevelValues, true, true) + ".",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("Organizational"),
				Validators: []validator.String{
					stringvalidator.OneOf(privacyLevelValues...),
				},
			},
			"connection_details": schema.ObjectAttribute{
				MarkdownDescription: "Connection details. Can be specified as a nested block or as an object.",
				Required:            true,
				AttributeTypes:      connectionDetailsAttributeTypes,
				// Connection details themselves are not sensitive, individual sensitive values should be marked separately
				Sensitive: false,
			},
			"credential_details": schema.SingleNestedAttribute{
				MarkdownDescription: "Credential details. Can be specified as a nested block or as an object.",
				Required:            true,
				Sensitive:           false,
				Attributes: map[string]schema.Attribute{
					"single_sign_on_type": schema.StringAttribute{
						MarkdownDescription: "Single sign-on type.",
						Required:            true,
					},
					"connection_encryption": schema.StringAttribute{
						MarkdownDescription: "Connection encryption type.",
						Required:            true,
					},
					"skip_test_connection": schema.BoolAttribute{
						MarkdownDescription: "Whether to skip test connection.",
						Required:            true,
					},
					"credentials": schema.SingleNestedAttribute{
						MarkdownDescription: "Credentials configuration.",
						Required:            true,
						Sensitive:           false,
						Attributes: map[string]schema.Attribute{
							"credential_type": schema.StringAttribute{
								MarkdownDescription: "Credential type.",
								Required:            true,
							},
							// Make all other credential fields optional
							"username": schema.StringAttribute{
								MarkdownDescription: "Username for Basic or Windows authentication.",
								Optional:            true,
								Sensitive:           true,
							},
							"password": schema.StringAttribute{
								MarkdownDescription: "Password for Basic or Windows authentication.",
								Optional:            true,
								Sensitive:           true,
							},
							"key": schema.StringAttribute{
								MarkdownDescription: "Key for Key authentication.",
								Optional:            true,
								Sensitive:           true,
							},
							"application_id": schema.StringAttribute{
								MarkdownDescription: "Application ID for Service Principal authentication.",
								Optional:            true,
								Sensitive:           true,
							},
							"application_secret": schema.StringAttribute{
								MarkdownDescription: "Application Secret for Service Principal authentication.",
								Optional:            true,
								Sensitive:           true,
							},
							"tenant_id": schema.StringAttribute{
								MarkdownDescription: "Tenant ID for Service Principal authentication.",
								Optional:            true,
								Sensitive:           true,
							},
							"sas_token": schema.StringAttribute{
								MarkdownDescription: "SAS Token for Shared Access Signature authentication.",
								Optional:            true,
								Sensitive:           true,
							},
							"domain": schema.StringAttribute{
								MarkdownDescription: "Domain for Windows authentication.",
								Optional:            true,
								Sensitive:           true,
							},
						},
					},
				},
			},
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{}),
		},
	}
}
