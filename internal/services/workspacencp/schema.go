// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package workspacencp

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema" //revive:disable-line:import-alias-naming
	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema"   //revive:disable-line:import-alias-naming
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	superschema "github.com/orange-cloudavenue/terraform-plugin-framework-superschema"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
)

func itemSchema() superschema.Schema {
	return superschema.Schema{
		Resource: superschema.SchemaDetails{
			MarkdownDescription: fabricitem.NewResourceMarkdownDescription(ItemTypeInfo, false),
		},
		DataSource: superschema.SchemaDetails{
			MarkdownDescription: fabricitem.NewDataSourceMarkdownDescription(ItemTypeInfo, false),
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
					Computed: true,
				},
			},
			"workspace_id": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The Workspace ID.",
					CustomType:          customtypes.UUIDType{},
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
				},
				DataSource: &schemaD.StringAttribute{
					Required: true,
				},
			},
			"inbound": superschema.SuperSingleNestedAttributeOf[rulesModel]{
				Common: &schemaR.SingleNestedAttribute{
					MarkdownDescription: "The policy for all inbound communications to a workspace.",
					Computed:            true,
				},
				Resource: &schemaR.SingleNestedAttribute{
					Optional: true,
				},
				Attributes: networkRulesAttributes("The policy for inbound communications to a workspace from public networks."),
			},
			"outbound": superschema.SuperSingleNestedAttributeOf[rulesModel]{
				Common: &schemaR.SingleNestedAttribute{
					MarkdownDescription: "The policy for all outbound communications from a workspace.",
					Computed:            true,
				},
				Resource: &schemaR.SingleNestedAttribute{
					Optional: true,
				},
				Attributes: networkRulesAttributes("The policy for outbound communications to public networks from a workspace."),
			},
			"timeouts": superschema.TimeoutAttribute{
				Resource: &superschema.ResourceTimeoutAttribute{
					Create: true,
					Read:   true,
					Update: true,
					Delete: true,
				},
				DataSource: &superschema.DatasourceTimeoutAttribute{
					Read: true,
				},
			},
		},
	}
}

func networkRulesAttributes(publicAccessRulesMarkdown string) superschema.Attributes {
	return superschema.Attributes{
		"public_access_rules": superschema.SuperSingleNestedAttributeOf[networkRulesModel]{
			Common: &schemaR.SingleNestedAttribute{
				MarkdownDescription: publicAccessRulesMarkdown,
				Computed:            true,
			},
			Resource: &schemaR.SingleNestedAttribute{
				Optional: true,
			},
			Attributes: superschema.Attributes{
				"default_action": superschema.SuperStringAttribute{
					Common: &schemaR.StringAttribute{
						MarkdownDescription: "Default policy for workspace access from public networks.",
						Computed:            true,
						Validators: []validator.String{
							stringvalidator.OneOf(utils.ConvertEnumsToStringSlices(fabcore.PossibleNetworkAccessRuleValues(), true)...),
						},
					},
					Resource: &schemaR.StringAttribute{
						Optional: true,
						Default:  stringdefault.StaticString(string(fabcore.NetworkAccessRuleAllow)),
					},
				},
			},
		},
	}
}
