// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package workspaceocr

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema" //revive:disable-line:import-alias-naming
	"github.com/hashicorp/terraform-plugin-framework/path"
	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema" //revive:disable-line:import-alias-naming
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
			"default_action": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "Defines the default behavior for all cloud connection types that are not explicitly listed in the rules array. If set to \"Allow\", all unspecified connection types are permitted by default. If set to \"Deny\", all unspecified connection types are blocked by default unless explicitly allowed. This setting acts as a global fallback policy and is critical for enforcing a secure default posture in environments where only known and trusted connections should be permitted.",
					Validators: []validator.String{
						stringvalidator.OneOf(utils.ConvertEnumsToStringSlices(fabcore.PossibleConnectionAccessActionTypeValues(), true)...),
					},
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"rules": rulesAttribute(),
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

func rulesAttribute() superschema.SuperListNestedAttributeOf[rulesModel] {
	return superschema.SuperListNestedAttributeOf[rulesModel]{
		Common: &schemaR.ListNestedAttribute{
			MarkdownDescription: "A list of rules that define outbound access behavior for specific cloud connection types. Each rule may include endpoint-based or workspace-based restrictions depending on supported connection types.",
			Computed:            true,
		},
		Resource: &schemaR.ListNestedAttribute{
			Optional: true,
			Default: listdefault.StaticValue(types.ListValueMust(
				types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"connection_type": customtypes.CaseInsensitiveStringType{},
						"default_action":  types.StringType,
						"allowed_endpoints": types.ListType{
							ElemType: types.ObjectType{
								AttrTypes: map[string]attr.Type{
									"hostname_pattern": types.StringType,
								},
							},
						},
						"allowed_workspaces": types.ListType{
							ElemType: types.ObjectType{
								AttrTypes: map[string]attr.Type{
									"workspace_id": customtypes.UUIDType{},
								},
							},
						},
					},
				},
				[]attr.Value{},
			)),
		},
		Attributes: superschema.Attributes{
			"connection_type": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "Specifies the cloud connection type to which the rule applies. The behavior and applicability of other rule properties (such as allowedEndpoints or allowedWorkspaces) may vary depending on the capabilities of connection type.",
					CustomType:          customtypes.CaseInsensitiveStringType{},
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"default_action": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "Defines the default outbound access behavior for the connectionType. This field determines whether connections of this type are permitted or blocked by default, unless further refined by allowedEndpoints or allowedWorkspaces. If set to \"Allow\": All connections of this type are permitted unless explicitly denied by a more specific rule. This field provides fine-grained control over each connection type and complements the global fallback behavior defined by defaultAction.",
					Validators: []validator.String{
						stringvalidator.OneOf(utils.ConvertEnumsToStringSlices(fabcore.PossibleConnectionAccessActionTypeValues(), true)...),
					},
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"allowed_endpoints":  allowedEndpointsAttribute(),
			"allowed_workspaces": allowedWorkspacesAttribute(),
		},
	}
}

func allowedEndpointsAttribute() superschema.SuperListNestedAttributeOf[endpointModel] {
	return superschema.SuperListNestedAttributeOf[endpointModel]{
		Common: &schemaR.ListNestedAttribute{
			MarkdownDescription: "Defines a list of explicitly permitted external endpoints for the connectionType. Each entry in the array represents a hostname pattern that is allowed for outbound communication from the workspace. This field is applicable only to connection types that support endpoint-based filtering (e.g., SQL, MySQL, Web, etc.). If defaultAction is set to \"Deny\" for the connection type, only the endpoints listed here will be allowed; all others will be blocked.",
			Computed:            true,
		},
		Resource: &schemaR.ListNestedAttribute{
			Optional: true,
			Validators: []validator.List{
				listvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("allowed_workspaces")),
			},
			Default: listdefault.StaticValue(types.ListValueMust(
				types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"hostname_pattern": types.StringType,
					},
				},
				[]attr.Value{},
			)),
		},
		Attributes: superschema.Attributes{
			"hostname_pattern": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "A wildcard-supported pattern that defines the allowed external endpoint. Examples include *.microsoft.com, api.contoso.com, or data.partner.org.",
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
		},
	}
}

func allowedWorkspacesAttribute() superschema.SuperListNestedAttributeOf[workspaceModel] {
	return superschema.SuperListNestedAttributeOf[workspaceModel]{
		Common: &schemaR.ListNestedAttribute{
			MarkdownDescription: "Specifies a list of workspace IDs that are explicitly permitted for outbound communication for the given fabric connectionType. This field is applicable only to fabric connection types that support workspace-based filtering, limited to Lakehouse, Warehouse, FabricSql, and PowerPlatformDataflows. When defaultAction is set to \"Deny\" for a connection type, only the workspaces listed in allowedWorkspaces will be allowed for outbound access; all others will be blocked.",
			Computed:            true,
		},
		Resource: &schemaR.ListNestedAttribute{
			Optional: true,
			Validators: []validator.List{
				listvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("allowed_endpoints")),
			},
			Default: listdefault.StaticValue(types.ListValueMust(
				types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"workspace_id": customtypes.UUIDType{},
					},
				},
				[]attr.Value{},
			)),
		},
		Attributes: superschema.Attributes{
			"workspace_id": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The unique identifier (GUID) of the target workspace that is allowed to be connected from current workspace.",
					CustomType:          customtypes.UUIDType{},
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
		},
	}
}
