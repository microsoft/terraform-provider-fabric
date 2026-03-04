// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package tags

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema" //revive:disable-line:import-alias-naming
	"github.com/hashicorp/terraform-plugin-framework/path"
	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema" //revive:disable-line:import-alias-naming
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	fabadmin "github.com/microsoft/fabric-sdk-go/fabric/admin"
	superschema "github.com/orange-cloudavenue/terraform-plugin-framework-superschema"
	superstringvalidator "github.com/orange-cloudavenue/terraform-plugin-framework-validators/stringvalidator"

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
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The " + ItemTypeInfo.Name + " ID.",
					CustomType:          customtypes.UUIDType{},
					Computed:            true,
				},
				DataSource: &schemaD.StringAttribute{
					Optional: !isList,
				},
				Resource: &schemaR.StringAttribute{
					Optional: true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
				},
			},
			"display_name": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The " + ItemTypeInfo.Name + " display name.",
				},
				DataSource: &schemaD.StringAttribute{
					Optional: !isList,
					Computed: true,
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
					Validators: []validator.String{
						stringvalidator.LengthAtMost(40),
					},
				},
			},
			"scope": superschema.SuperSingleNestedAttributeOf[scopeModel]{
				Common: &schemaR.SingleNestedAttribute{
					MarkdownDescription: "Represents a tag scope.",
					Computed:            true,
				},
				Resource: &schemaR.SingleNestedAttribute{
					Optional: true,
					PlanModifiers: []planmodifier.Object{
						objectplanmodifier.RequiresReplace(),
					},
				},
				Attributes: map[string]superschema.Attribute{
					"domain_id": superschema.StringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "Domain object ID.",
							Computed:            true,
						},
						Resource: &schemaR.StringAttribute{
							Optional: true,
							Validators: []validator.String{
								superstringvalidator.RequireIfAttributeIsOneOf(
									path.MatchRelative().AtParent().AtName("type"),
									[]attr.Value{types.StringValue(string(fabadmin.TagScopeTypeDomain))},
								),
							},
						},
					},
					"type": superschema.StringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "Scope Type.",
							Computed:            true,
						},
						Resource: &schemaR.StringAttribute{
							Optional: true,
							Validators: []validator.String{
								stringvalidator.OneOf(utils.ConvertEnumsToStringSlices(fabadmin.PossibleTagScopeTypeValues(), true)...),
							},
						},
					},
				},
			},
			"timeouts": superschema.TimeoutAttribute{
				DataSource: dsTimeout,
				Resource: &superschema.ResourceTimeoutAttribute{
					Create: true,
					Read:   true,
					Delete: true,
					Update: true,
				},
			},
		},
	}
}
