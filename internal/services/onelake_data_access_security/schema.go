// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package onelake_data_access_security

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema" //revive:disable-line:import-alias-naming
	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema"   //revive:disable-line:import-alias-naming
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	superschema "github.com/orange-cloudavenue/terraform-plugin-framework-superschema"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
)

func itemSchema() superschema.Schema {
	return superschema.Schema{
		Attributes: map[string]superschema.Attribute{
			"value": superschema.SuperMapNestedAttributeOf[dataAccessRole]{
				Attributes: dataAccessRoleAttributes(),
			},
		},
	}
}

func dataAccessRoleAttributes() superschema.Attributes {
	return superschema.Attributes{
		"id": superschema.SuperStringAttribute{
			Common: &schemaR.StringAttribute{
				MarkdownDescription: "The unique ID of the Data Access Role.",
				CustomType:          customtypes.UUIDType{},
			},
			Resource: &schemaR.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			DataSource: &schemaD.StringAttribute{
				Computed: true,
			},
		},
		"name": superschema.StringAttribute{
			Common: &schemaR.StringAttribute{
				MarkdownDescription: "The name of the Data Access Role.",
			},
			Resource: &schemaR.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(256),
				},
			},
			DataSource: &schemaD.StringAttribute{
				Computed: true,
			},
		},
		"decision_rules": superschema.SuperMapNestedAttributeOf[decisionRule]{
			Attributes: superschema.Attributes{
				"effect": superschema.StringAttribute{
					Common: &schemaR.StringAttribute{
						MarkdownDescription: "The effect of the decision rule.",
					},
					Resource: &schemaR.StringAttribute{
						Required: true,
					},
					DataSource: &schemaD.StringAttribute{
						Computed: true,
					},
				},
				"permission": superschema.SuperMapNestedAttributeOf[permissionScope]{
					Attributes: superschema.Attributes{
						"attribute_name": superschema.SuperSingleNestedAttribute{
							Attributes: map[string]superschema.Attribute{
								"action": superschema.StringAttribute{
									Resource:   &schemaR.StringAttribute{Required: true},
									DataSource: &schemaD.StringAttribute{Computed: true},
								},
								"path": superschema.StringAttribute{
									Resource:   &schemaR.StringAttribute{Required: true},
									DataSource: &schemaD.StringAttribute{Computed: true},
								},
							},
						},
						"attribute_value_included_in": superschema.SuperMapNestedAttributeOf[types.String]{
							Attributes: superschema.Attributes{},
						},
					},
				},
			},
		},
		"member": superschema.SuperSingleNestedAttribute{
			Attributes: map[string]superschema.Attribute{
				"fabric_item_members": superschema.SuperListNestedAttributeOf[FabricItemMember]{
					Attributes: superschema.Attributes{
						"item_access": superschema.SuperMapNestedAttributeOf[types.String]{
							Attributes: superschema.Attributes{},
						},
						"source_path": superschema.StringAttribute{
							Resource:   &schemaR.StringAttribute{Required: true},
							DataSource: &schemaD.StringAttribute{Computed: true},
						},
					},
				},
				"microsoft_entra_members": superschema.SuperListNestedAttributeOf[MicrosoftEntraMember]{
					Attributes: superschema.Attributes{
						"object_id": superschema.StringAttribute{
							Common:     &schemaR.StringAttribute{CustomType: customtypes.UUIDType{}},
							Resource:   &schemaR.StringAttribute{Required: true},
							DataSource: &schemaD.StringAttribute{Computed: true},
						},
						"object_type": superschema.StringAttribute{
							Resource:   &schemaR.StringAttribute{Required: true},
							DataSource: &schemaD.StringAttribute{Computed: true},
						},
						"tenant_id": superschema.StringAttribute{
							Resource:   &schemaR.StringAttribute{Required: true},
							DataSource: &schemaD.StringAttribute{Computed: true},
						},
					},
				},
			},
		},
	}
}
