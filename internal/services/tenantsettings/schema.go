// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package tenantsettings

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema" //revive:disable-line:import-alias-naming
	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema"   //revive:disable-line:import-alias-naming
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
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
		Resource: superschema.SchemaDetails{
			MarkdownDescription: fabricitem.NewResourceMarkdownDescription(ItemTypeInfo, false),
		},
		DataSource: superschema.SchemaDetails{
			MarkdownDescription: fabricitem.NewDataSourceMarkdownDescription(ItemTypeInfo, false),
		},
		Attributes: map[string]superschema.Attribute{
			"setting_name": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The name of the tenant setting.",
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
				},
				DataSource: &schemaD.StringAttribute{
					Optional: true,
					Computed: true,
				},
			},
			"tenant_setting_group": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "Tenant setting group name.",
				},
				Resource: &schemaR.StringAttribute{
					Computed: true,
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"title": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The title of the tenant setting.",
				},
				Resource: &schemaR.StringAttribute{
					Computed: true,
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"delete_behaviour": superschema.StringAttribute{
				Resource: &schemaR.StringAttribute{
					MarkdownDescription: "Indicates whether the tenant setting is disabled when deleted. NoChange - The tenant setting is not disabled when deleted. Disable - The tenant setting is disabled when deleted.",
					Validators: []validator.String{
						stringvalidator.OneOf(utils.ConvertEnumsToStringSlices(PossibleDeleteBehaviourValues(), true)...),
					},
					Computed: true,
					Optional: true,
					Default:  stringdefault.StaticString(string(NoChange)),
				},
			},
			"can_specify_security_groups": superschema.BoolAttribute{
				Common: &schemaR.BoolAttribute{
					MarkdownDescription: "Indicates if the tenant setting is enabled for a security group. False - The tenant setting is enabled for the entire organization. True - The tenant setting is enabled for security groups.",
				},
				Resource: &schemaR.BoolAttribute{
					Computed: true,
				},
				DataSource: &schemaD.BoolAttribute{
					Computed: true,
				},
			},
			"enabled": superschema.BoolAttribute{
				Common: &schemaR.BoolAttribute{
					MarkdownDescription: "The status of the tenant setting. False - Disabled, True - Enabled.",
				},
				Resource: &schemaR.BoolAttribute{
					Required: true,
				},
				DataSource: &schemaD.BoolAttribute{
					Computed: true,
				},
			},
			"delegate_to_capacity": superschema.BoolAttribute{
				Common: &schemaR.BoolAttribute{
					MarkdownDescription: "Indicates whether the tenant setting can be delegated to a capacity admin. False - Capacity admin cannot override the tenant setting. True - Capacity admin can override the tenant setting.",
				},
				Resource: &schemaR.BoolAttribute{
					Optional: true,
					Computed: true,
				},
				DataSource: &schemaD.BoolAttribute{
					Computed: true,
				},
			},
			"delegate_to_domain": superschema.BoolAttribute{
				Common: &schemaR.BoolAttribute{
					MarkdownDescription: "Indicates whether the tenant setting can be delegated to a domain admin. False - Domain admin cannot override the tenant setting. True - Domain admin can override the tenant setting.",
				},
				Resource: &schemaR.BoolAttribute{
					Optional: true,
					Computed: true,
				},
				DataSource: &schemaD.BoolAttribute{
					Computed: true,
				},
			},
			"delegate_to_workspace": superschema.BoolAttribute{
				Common: &schemaR.BoolAttribute{
					MarkdownDescription: "Indicates whether the tenant setting can be delegated to a workspace admin. False - Workspace admin cannot override the tenant setting. True - Workspace admin can override the tenant setting.",
				},
				Resource: &schemaR.BoolAttribute{
					Optional: true,
					Computed: true,
				},
				DataSource: &schemaD.BoolAttribute{
					Computed: true,
				},
			},
			"enabled_security_groups": superschema.SuperSetNestedAttributeOf[tenantSettingsSecurityGroup]{
				Common: &schemaR.SetNestedAttribute{
					MarkdownDescription: "A list of enabled security groups.",
				},
				Resource: &schemaR.SetNestedAttribute{
					Optional: true,
					Computed: true,
				},
				DataSource: &schemaD.SetNestedAttribute{
					Computed: true,
				},
				Attributes: superschema.Attributes{
					"graph_id": superschema.StringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "The graph ID of the security group.",
							CustomType:          customtypes.UUIDType{},
						},
						Resource: &schemaR.StringAttribute{
							Optional: true,
							Computed: true,
						},
						DataSource: &schemaD.StringAttribute{
							Computed: true,
						},
					},
					"name": superschema.StringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "The name of the security group.",
						},
						Resource: &schemaR.StringAttribute{
							Optional: true,
							Computed: true,
						},
						DataSource: &schemaD.StringAttribute{
							Computed: true,
						},
					},
				},
			},
			"excluded_security_groups": superschema.SuperSetNestedAttributeOf[tenantSettingsSecurityGroup]{
				Common: &schemaR.SetNestedAttribute{
					MarkdownDescription: "A list of excluded security groups.",
				},
				Resource: &schemaR.SetNestedAttribute{
					Optional: true,
					Computed: true,
				},
				DataSource: &schemaD.SetNestedAttribute{
					Computed: true,
				},
				Attributes: superschema.Attributes{
					"graph_id": superschema.StringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "The graph ID of the security group.",
							CustomType:          customtypes.UUIDType{},
						},
						Resource: &schemaR.StringAttribute{
							Optional: true,
							Computed: true,
						},
						DataSource: &schemaD.StringAttribute{
							Computed: true,
						},
					},
					"name": superschema.StringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "The name of the security group.",
						},
						Resource: &schemaR.StringAttribute{
							Optional: true,
							Computed: true,
						},
						DataSource: &schemaD.StringAttribute{
							Computed: true,
						},
					},
				},
			},
			"properties": superschema.SuperSetNestedAttributeOf[tenantSettingsProperty]{
				Common: &schemaR.SetNestedAttribute{
					MarkdownDescription: "Tenant setting properties.",
				},
				Resource: &schemaR.SetNestedAttribute{
					Optional: true,
					Computed: true,
				},
				DataSource: &schemaD.SetNestedAttribute{
					Computed: true,
				},
				Attributes: superschema.Attributes{
					"name": superschema.StringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "The name of the property.",
						},
						Resource: &schemaR.StringAttribute{
							Optional: true,
							Computed: true,
						},
						DataSource: &schemaD.StringAttribute{
							Computed: true,
						},
					},
					"type": superschema.StringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "The type of the property.",
						},
						Resource: &schemaR.StringAttribute{
							Optional: true,
							Computed: true,
							Validators: []validator.String{
								stringvalidator.OneOf(utils.ConvertEnumsToStringSlices(fabadmin.PossibleTenantSettingPropertyTypeValues(), true)...),
							},
						},
						DataSource: &schemaD.StringAttribute{
							Computed: true,
						},
					},
					"value": superschema.StringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "The value of the property.",
						},
						Resource: &schemaR.StringAttribute{
							Optional: true,
							Computed: true,
						},
						DataSource: &schemaD.StringAttribute{
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
				DataSource: dsTimeout,
			},
		},
	}
}
