// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package tags

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator" //revive:disable-line:import-alias-naming
	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema"    //revive:disable-line:import-alias-naming
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
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The " + ItemTypeInfo.Name + " ID.",
					CustomType:          customtypes.UUIDType{},
				},
				Resource: &schemaR.StringAttribute{
					Computed: true,
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
				Resource: &schemaR.StringAttribute{
					Computed: true,
					Optional: true,
					Validators: []validator.String{
						stringvalidator.LengthAtMost(40),
					},
				},
			},
			"tags": superschema.SuperSetNestedAttributeOf[baseTagModel]{
				Common: &schemaR.SetNestedAttribute{
					MarkdownDescription: "List of tags associated with the resource.",
				},
				Resource: &schemaR.SetNestedAttribute{
					Optional: true,
				},
				Attributes: map[string]superschema.Attribute{
					"id": superschema.StringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "The " + ItemTypeInfo.Name + " ID.",
						},
						Resource: &schemaR.StringAttribute{
							Computed: true,
							Optional: true,
						},
					},
					"display_name": superschema.StringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "The " + ItemTypeInfo.Name + " display name.",
						},
						Resource: &schemaR.StringAttribute{
							Optional: true,
							Computed: true,
						},
					},
					"scope": superschema.SuperSingleNestedAttributeOf[scopeModel]{
						Common: &schemaR.SingleNestedAttribute{
							MarkdownDescription: "Represents a tag scope.",
						},
						Resource: &schemaR.SingleNestedAttribute{
							Optional: true,
							Computed: true,
						},
						Attributes: map[string]superschema.Attribute{
							"type": superschema.StringAttribute{
								Common: &schemaR.StringAttribute{
									MarkdownDescription: "Scope Type.",
									Validators: []validator.String{
										stringvalidator.OneOf(utils.ConvertEnumsToStringSlices(fabadmin.PossibleTagScopeTypeValues(), true)...),
									},
								},
								Resource: &schemaR.StringAttribute{
									Optional: true,
									Computed: true,
								},
							},
						},
					},
				},
			},
			"scope": superschema.SuperSingleNestedAttributeOf[scopeModel]{
				Common: &schemaR.SingleNestedAttribute{
					MarkdownDescription: "Represents a tag scope.",
				},
				Resource: &schemaR.SingleNestedAttribute{
					Optional: true,
					Computed: true,
				},
				Attributes: map[string]superschema.Attribute{
					"type": superschema.StringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "Scope Type.",
							Validators: []validator.String{
								stringvalidator.OneOf(utils.ConvertEnumsToStringSlices(fabadmin.PossibleTagScopeTypeValues(), true)...),
							},
						},
						Resource: &schemaR.StringAttribute{
							Optional: true,
							Computed: true,
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
