// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package connection

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
)

func DataSourceConnectionSchema(ctx context.Context) schema.Schema {
	tflog.Info(ctx, "Building schema for connection data source")
	
	return schema.Schema{
		MarkdownDescription: "Use this data source to retrieve details of a Microsoft Fabric Connection by ID or display name.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The Connection ID.",
				CustomType:          customtypes.UUIDType{},
				Optional:            true,
				Computed:            true,
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "The Connection display name. Either id or display_name must be provided to retrieve a connection.",
				Optional:            true,
				Computed:            true,
			},
			"connectivity_type": schema.StringAttribute{
				MarkdownDescription: "Connectivity type.",
				Computed:            true,
			},
			"gateway_id": schema.StringAttribute{
				MarkdownDescription: "Gateway ID.",
				CustomType:          customtypes.UUIDType{},
				Computed:            true,
			},
			"privacy_level": schema.StringAttribute{
				MarkdownDescription: "Privacy level.",
				Computed:            true,
			},
			"connection_details": schema.SingleNestedAttribute{
				MarkdownDescription: "Connection details.",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						MarkdownDescription: "Connection type.",
						Computed:            true,
					},
					"creation_method": schema.StringAttribute{
						MarkdownDescription: "Creation method.",
						Computed:            true,
					},
					"parameters": schema.ListNestedAttribute{
						MarkdownDescription: "Connection parameters.",
						Computed:            true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									MarkdownDescription: "Parameter name.",
									Computed:            true,
								},
								"data_type": schema.StringAttribute{
									MarkdownDescription: "Data type.",
									Computed:            true,
								},
								"value": schema.StringAttribute{
									MarkdownDescription: "Parameter value.",
									Computed:            true,
								},
							},
						},
					},
				},
			},
			"credential_details": schema.SingleNestedAttribute{
				MarkdownDescription: "Credential details.",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"single_sign_on_type": schema.StringAttribute{
						MarkdownDescription: "Single sign-on type.",
						Computed:            true,
					},
					"connection_encryption": schema.StringAttribute{
						MarkdownDescription: "Connection encryption.",
						Computed:            true,
					},
					"skip_test_connection": schema.BoolAttribute{
						MarkdownDescription: "Skip test connection.",
						Computed:            true,
					},
					"credentials": schema.SingleNestedAttribute{
						MarkdownDescription: "Credentials.",
						Computed:            true,
						Attributes: map[string]schema.Attribute{
							"credential_type": schema.StringAttribute{
								MarkdownDescription: "Credential type.",
								Computed:            true,
							},
							"username": schema.StringAttribute{
								MarkdownDescription: "Username for Basic authentication.",
								Computed:            true,
								Sensitive:           true,
							},
							"password": schema.StringAttribute{
								MarkdownDescription: "Password for Basic authentication.",
								Computed:            true,
								Sensitive:           true,
							},
							"key": schema.StringAttribute{
								MarkdownDescription: "Key for Key authentication.",
								Computed:            true,
								Sensitive:           true,
							},
							"application_id": schema.StringAttribute{
								MarkdownDescription: "Application ID for Service Principal authentication.",
								Computed:            true,
								Sensitive:           true,
							},
							"application_secret": schema.StringAttribute{
								MarkdownDescription: "Application Secret for Service Principal authentication.",
								Computed:            true,
								Sensitive:           true,
							},
							"tenant_id": schema.StringAttribute{
								MarkdownDescription: "Tenant ID for Service Principal authentication.",
								Computed:            true,
								Sensitive:           true,
							},
							"sas_token": schema.StringAttribute{
								MarkdownDescription: "SAS Token for Shared Access Signature authentication.",
								Computed:            true,
								Sensitive:           true,
							},
							"domain": schema.StringAttribute{
								MarkdownDescription: "Domain for Windows authentication.",
								Computed:            true,
							},
						},
					},
				},
			},
			"timeouts": timeouts.Attributes(ctx),
		},
	}
}