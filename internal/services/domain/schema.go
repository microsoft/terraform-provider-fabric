// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package domain

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema" //revive:disable-line:import-alias-naming
	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema"   //revive:disable-line:import-alias-naming
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	superschema "github.com/orange-cloudavenue/terraform-plugin-framework-superschema"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
)

func itemSchema(isList bool) superschema.Schema { //revive:disable-line:flag-parameter
	markdownDescriptionR := fabricitem.NewResourceMarkdownDescription(ItemTypeInfo, false)
	markdownDescriptionD := fabricitem.NewDataSourceMarkdownDescription(ItemTypeInfo, isList)

	var dsTimeout *superschema.DatasourceTimeoutAttribute

	if !isList {
		dsTimeout = &superschema.DatasourceTimeoutAttribute{
			Read: true,
		}
	}

	return superschema.Schema{
		Resource: superschema.SchemaDetails{
			MarkdownDescription: markdownDescriptionR,
		},
		DataSource: superschema.SchemaDetails{
			MarkdownDescription: markdownDescriptionD,
		},
		Attributes: map[string]superschema.Attribute{
			"id": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The " + ItemTypeInfo.Name + " ID.",
					CustomType:          customtypes.UUIDType{},
				},
				Resource: &schemaR.StringAttribute{
					Computed: true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
				},
				DataSource: &schemaD.StringAttribute{
					Required: !isList,
					Computed: isList,
				},
			},
			"display_name": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The " + ItemTypeInfo.Name + " display name.",
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
					MarkdownDescription: "The " + ItemTypeInfo.Name + " description.",
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
			"parent_domain_id": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The " + ItemTypeInfo.Name + " parent ID.",
					CustomType:          customtypes.UUIDType{},
				},
				Resource: &schemaR.StringAttribute{
					Optional: true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"default_label_id": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The domain default sensitivity label. To remove the defaultLabelId from a domain, set its value to an empty UUID in your request: '00000000-0000-0000-0000-000000000000'.",
					CustomType:          customtypes.UUIDType{},
				},
				Resource: &schemaR.StringAttribute{
					Optional: true,
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
