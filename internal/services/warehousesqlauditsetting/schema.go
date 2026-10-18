// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package warehousesqlauditsetting

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"                    //revive:disable-line:import-alias-naming
	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema" //revive:disable-line:import-alias-naming
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	fabwarehouse "github.com/microsoft/fabric-sdk-go/fabric/warehouse"
	superschema "github.com/orange-cloudavenue/terraform-plugin-framework-superschema"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

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
			"warehouse_id": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The Warehouse ID.",
					CustomType:          customtypes.UUIDType{},
					Required:            true,
				},
				Resource: &schemaR.StringAttribute{
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
				},
			},
			"state": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "Audit settings state type. Possible values: " + utils.ConvertStringSlicesToString(fabwarehouse.PossibleAuditSettingsStateValues(), true, true) + ".",
					Validators: []validator.String{
						stringvalidator.OneOf(utils.ConvertEnumsToStringSlices(fabwarehouse.PossibleAuditSettingsStateValues(), true)...),
					},
					Computed: true,
				},
				Resource: &schemaR.StringAttribute{
					Optional: true,
					Default:  stringdefault.StaticString(string(fabwarehouse.AuditSettingsStateDisabled)),
				},
			},
			"retention_days": superschema.Int32Attribute{
				Common: &schemaR.Int32Attribute{
					MarkdownDescription: "Retention days. `0` indicates indefinite retention period.",
					Computed:            true,
				},
				Resource: &schemaR.Int32Attribute{
					Optional: true,
					Default:  int32default.StaticInt32(0),
					Validators: []validator.Int32{
						int32validator.AtLeast(0),
					},
				},
			},
			"audit_actions_and_groups": superschema.SuperSetAttribute{
				Common: &schemaR.SetAttribute{
					MarkdownDescription: "Audit actions and groups.",
					CustomType: supertypes.SetTypeOf[types.String]{
						SetType: basetypes.SetType{
							ElemType: types.StringType,
						},
					},
					Computed: true,
				},
				Resource: &schemaR.SetAttribute{
					Optional: true,
					Default: setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{
						types.StringValue("SUCCESSFUL_DATABASE_AUTHENTICATION_GROUP"),
						types.StringValue("FAILED_DATABASE_AUTHENTICATION_GROUP"),
						types.StringValue("BATCH_COMPLETED_GROUP"),
					})),
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
