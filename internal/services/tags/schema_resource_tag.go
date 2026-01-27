// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package tags

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"   //revive:disable-line:import-alias-naming
	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator" //revive:disable-line:import-alias-naming
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator" //revive:disable-line:import-alias-naming
	"github.com/hashicorp/terraform-plugin-framework/path"
	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema" //revive:disable-line:import-alias-naming
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	fabadmin "github.com/microsoft/fabric-sdk-go/fabric/admin"
	superschema "github.com/orange-cloudavenue/terraform-plugin-framework-superschema"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
)

func resourceItemSchema() superschema.Schema { //revive:disable-line:flag-parameter
	return superschema.Schema{
		Resource: superschema.SchemaDetails{
			MarkdownDescription: fabricitem.NewResourceMarkdownDescription(ItemTypeInfo, false),
		},
		Attributes: map[string]superschema.Attribute{
			"id": superschema.SuperStringAttribute{
				Resource: &schemaR.StringAttribute{
					MarkdownDescription: "The " + ItemTypeInfo.Name + " ID.",
					CustomType:          customtypes.UUIDType{},
					Computed:            true,
					Optional:            true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
					Validators: []validator.String{
						stringvalidator.ConflictsWith(path.MatchRoot("tags")),
						stringvalidator.AlsoRequires(
							path.MatchRoot("display_name"),
							path.MatchRoot("scope"),
						),
					},
				},
			},
			"display_name": superschema.StringAttribute{
				Resource: &schemaR.StringAttribute{
					Computed:            true,
					Optional:            true,
					MarkdownDescription: "The " + ItemTypeInfo.Name + " display name.",
					Validators: []validator.String{
						stringvalidator.LengthAtMost(40),
						stringvalidator.ConflictsWith(path.MatchRoot("tags")),
						stringvalidator.AlsoRequires(
							path.MatchRoot("id"),
							path.MatchRoot("scope"),
						),
					},
				},
			},
			"tags": superschema.SuperListNestedAttributeOf[baseTagModel]{
				Resource: &schemaR.ListNestedAttribute{
					MarkdownDescription: "List of tags associated with the resource.",
					Optional:            true,
					Validators: []validator.List{
						listvalidator.ConflictsWith(
							path.MatchRoot("id"),
							path.MatchRoot("display_name"),
							path.MatchRoot("scope"),
						),
					},
				},
				Attributes: map[string]superschema.Attribute{
					"id": superschema.StringAttribute{
						Resource: &schemaR.StringAttribute{
							MarkdownDescription: "The " + ItemTypeInfo.Name + " ID.",
							Computed:            true,
						},
					},
					"display_name": superschema.StringAttribute{
						Resource: &schemaR.StringAttribute{
							MarkdownDescription: "The " + ItemTypeInfo.Name + " display name.",
							Required:            true,
						},
					},
					"scope": superschema.SuperSingleNestedAttributeOf[scopeModel]{
						Resource: &schemaR.SingleNestedAttribute{
							Computed:            true,
							MarkdownDescription: "Represents a tag scope.",
						},
						Attributes: map[string]superschema.Attribute{
							"type": superschema.StringAttribute{
								Resource: &schemaR.StringAttribute{
									Required:            true,
									MarkdownDescription: "Scope Type.",
									Validators: []validator.String{
										stringvalidator.OneOf(utils.ConvertEnumsToStringSlices(fabadmin.PossibleTagScopeTypeValues(), true)...),
									},
								},
							},
						},
					},
				},
			},
			"scope": superschema.SuperSingleNestedAttributeOf[scopeModel]{
				Resource: &schemaR.SingleNestedAttribute{
					MarkdownDescription: "Represents a tag scope.",
					Optional:            true,
					Computed:            true,
					Validators: []validator.Object{
						objectvalidator.ConflictsWith(path.MatchRoot("tags")),
						objectvalidator.AlsoRequires(
							path.MatchRoot("id"),
							path.MatchRoot("display_name"),
						),
					},
				},
				Attributes: map[string]superschema.Attribute{
					"type": superschema.StringAttribute{
						Resource: &schemaR.StringAttribute{
							Optional:            true,
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
				Resource: &superschema.ResourceTimeoutAttribute{
					Create: true,
					Read:   true,
					Update: true,
					Delete: true,
				},
			},
		},
	}
}
