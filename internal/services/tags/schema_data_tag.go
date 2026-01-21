// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package tags

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema" //revive:disable-line:import-alias-naming
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	fabadmin "github.com/microsoft/fabric-sdk-go/fabric/admin"
	superschema "github.com/orange-cloudavenue/terraform-plugin-framework-superschema"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
)

func itemSchema(isList bool) superschema.Schema { //revive:disable-line:flag-parameter
	var dsTimeout *superschema.DatasourceTimeoutAttribute

	if !isList {
		dsTimeout = &superschema.DatasourceTimeoutAttribute{
			Read: true,
		}
	}

	return superschema.Schema{
		DataSource: superschema.SchemaDetails{
			MarkdownDescription: fabricitem.NewDataSourceMarkdownDescription(ItemTypeInfo, isList),
		},
		Attributes: map[string]superschema.Attribute{
			"id": superschema.SuperStringAttribute{
				DataSource: &schemaD.StringAttribute{
					MarkdownDescription: "The " + ItemTypeInfo.Name + " ID.",
					CustomType:          customtypes.UUIDType{},
					Optional:            !isList,
					Computed:            true,
				},
			},
			"display_name": superschema.StringAttribute{
				DataSource: &schemaD.StringAttribute{
					MarkdownDescription: "The " + ItemTypeInfo.Name + " display name.",
					Optional:            !isList,
					Computed:            true,
				},
			},
			"scope": superschema.SuperSingleNestedAttributeOf[scopeModel]{
				DataSource: &schemaD.SingleNestedAttribute{
					MarkdownDescription: "Represents a tag scope.",
					Computed:            true,
				},
				Attributes: map[string]superschema.Attribute{
					"type": superschema.StringAttribute{
						DataSource: &schemaD.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Scope Type.",
							Validators: []validator.String{
								stringvalidator.OneOf(utils.ConvertEnumsToStringSlices(fabadmin.PossibleTagScopeTypeValues(), true)...),
							},
						},
					},
				},
			},
			"timeouts": superschema.TimeoutAttribute{
				DataSource: dsTimeout,
			},
		},
	}
}
