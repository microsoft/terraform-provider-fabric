// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package workspacefwr

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema" //revive:disable-line:import-alias-naming
	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema"   //revive:disable-line:import-alias-naming
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	superschema "github.com/orange-cloudavenue/terraform-plugin-framework-superschema"
	superstringvalidator "github.com/orange-cloudavenue/terraform-plugin-framework-validators/stringvalidator"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
)

// Service-side limits, see https://learn.microsoft.com/fabric/security/security-workspace-level-firewall-overview
const (
	maxRules           = 256
	maxRuleDisplayName = 128
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
			"workspace_id": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The Workspace ID.",
					CustomType:          customtypes.UUIDType{},
					Required:            true,
				},
				Resource: &schemaR.StringAttribute{
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
				},
			},
			"rules": superschema.SuperSetNestedAttributeOf[firewallRuleModel]{
				Common: &schemaR.SetNestedAttribute{
					MarkdownDescription: "A set of rules that define the IP addresses permitted for inbound access to the Workspace. " +
						"Rules are only enforced when `inbound.public_access_rules.default_action` of the `fabric_workspace_network_communication_policy` is set to `Deny`. " +
						"Only public IP addresses are supported.",
					Computed: true,
				},
				Resource: &schemaR.SetNestedAttribute{
					Optional: true,
					Validators: []validator.Set{
						setvalidator.SizeAtMost(maxRules),
					},
					Default: setdefault.StaticValue(types.SetValueMust(
						types.ObjectType{
							AttrTypes: map[string]attr.Type{
								"display_name": types.StringType,
								"value":        types.StringType,
							},
						},
						[]attr.Value{},
					)),
				},
				Attributes: superschema.Attributes{
					"display_name": superschema.SuperStringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "The name of the rule. Must be unique within the Workspace.",
						},
						Resource: &schemaR.StringAttribute{
							Required: true,
							Validators: []validator.String{
								stringvalidator.LengthBetween(1, maxRuleDisplayName),
							},
						},
						DataSource: &schemaD.StringAttribute{
							Computed: true,
						},
					},
					"value": superschema.SuperStringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "The IP range in start-end format (for example, `192.0.2.1-192.0.2.10`), CIDR notation (for example, `203.0.113.0/24`), or a single IP address.",
						},
						Resource: &schemaR.StringAttribute{
							Required: true,
							Validators: []validator.String{
								superstringvalidator.IsNetwork([]superstringvalidator.NetworkValidatorType{
									superstringvalidator.IPV4,
									superstringvalidator.IPV4WithCIDR,
									superstringvalidator.IPV4Range,
								}, true),
							},
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
				DataSource: &superschema.DatasourceTimeoutAttribute{
					Read: true,
				},
			},
		},
	}
}
