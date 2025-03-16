// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package domain

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema" //revive:disable-line:import-alias-naming
	"github.com/hashicorp/terraform-plugin-framework/path"
	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema" //revive:disable-line:import-alias-naming
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	fabadmin "github.com/microsoft/fabric-sdk-go/fabric/admin"
	superschema "github.com/orange-cloudavenue/terraform-plugin-framework-superschema"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
)

func domainSchema(dsList bool) superschema.Schema { //revive:disable-line:flag-parameter
	markdownDescriptionR := "The " + ItemName + " resource allows you to manage [" + ItemName + "](" + ItemDocsURL + ")."
	markdownDescriptionR = fabricitem.GetResourceSPNSupportNote(markdownDescriptionR, ItemSPNSupport)
	markdownDescriptionR = fabricitem.GetResourcePreviewNote(markdownDescriptionR, ItemPreview)

	var dsTimeout *superschema.DatasourceTimeoutAttribute
	var markdownDescriptionD string

	if !dsList {
		dsTimeout = &superschema.DatasourceTimeoutAttribute{
			Read: true,
		}

		markdownDescriptionD = "The " + ItemName + " data-source allows you to read [" + ItemName + "](" + ItemDocsURL + ") details."
	} else {
		markdownDescriptionD = "The " + ItemName + " data-source allows you to read a list of [" + ItemName + "](" + ItemDocsURL + ") details."
	}

	markdownDescriptionD = fabricitem.GetDataSourceSPNSupportNote(markdownDescriptionD, ItemSPNSupport)
	markdownDescriptionD = fabricitem.GetDataSourcePreviewNote(markdownDescriptionD, ItemPreview)

	return superschema.Schema{
		Resource: superschema.SchemaDetails{
			MarkdownDescription: markdownDescriptionR,
		},
		DataSource: superschema.SchemaDetails{
			MarkdownDescription: markdownDescriptionD,
		},
		Attributes: map[string]superschema.Attribute{
			"id": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The " + ItemName + " ID.",
					CustomType:          customtypes.UUIDType{},
				},
				Resource: &schemaR.StringAttribute{
					Computed: true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
				},
				DataSource: &schemaD.StringAttribute{
					Required: !dsList,
					Computed: dsList,
				},
			},
			"display_name": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The " + ItemName + " display name.",
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
					Validators: []validator.String{
						stringvalidator.LengthAtMost(40),
					},
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"description": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The " + ItemName + " description.",
				},
				Resource: &schemaR.StringAttribute{
					Optional: true,
					Computed: true,
					Default:  stringdefault.StaticString(""),
					Validators: []validator.String{
						stringvalidator.LengthAtMost(256),
					},
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"parent_domain_id": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The " + ItemName + " parent ID.",
					CustomType:          customtypes.UUIDType{},
				},
				Resource: &schemaR.StringAttribute{
					Optional: true,
					Validators: []validator.String{
						stringvalidator.ConflictsWith(path.MatchRoot("contributors_scope")),
					},
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"contributors_scope": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The " + ItemName + " contributors scope.",
					Validators: []validator.String{
						stringvalidator.OneOf(utils.ConvertEnumsToStringSlices(fabadmin.PossibleContributorsScopeTypeValues(), true)...),
					},
				},
				Resource: &schemaR.StringAttribute{
					Optional: true,
					Computed: true,
					Validators: []validator.String{
						stringvalidator.ConflictsWith(path.MatchRoot("parent_domain_id")),
					},
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"timeouts": superschema.TimeoutAttribute{
				Resource: &superschema.ResourceTimeoutAttribute{
					Create: true,
					Read:   true,
					Update: true,
					Delete: true,
				},
				DataSource: dsTimeout,
			},
		},
	}
}
