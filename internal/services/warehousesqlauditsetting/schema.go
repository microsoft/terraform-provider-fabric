// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package warehousesqlauditsetting

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema" //revive:disable-line:import-alias-naming
	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema"   //revive:disable-line:import-alias-naming
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
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
			"item_id": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The item ID.",
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
			"state": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "Audit settings state type. Possible values: " + utils.ConvertStringSlicesToString(fabwarehouse.PossibleAuditSettingsStateValues(), true, true) + ".",
					Validators: []validator.String{
						stringvalidator.OneOf(utils.ConvertEnumsToStringSlices(fabwarehouse.PossibleAuditSettingsStateValues(), true)...),
					},
				},
				Resource: &schemaR.StringAttribute{
					Optional: true,
					Computed: true,
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"retention_days": superschema.Int32Attribute{
				Common: &schemaR.Int32Attribute{
					MarkdownDescription: "Retention days. `0` indicates indefinite retention period.",
				},
				Resource: &schemaR.Int32Attribute{
					Optional: true,
					Computed: true,
					Validators: []validator.Int32{
						int32validator.AtLeast(0),
					},
				},
				DataSource: &schemaD.Int32Attribute{
					Computed: true,
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
					ElementType: types.StringType,
				},
				Resource: &schemaR.SetAttribute{
					Required: true,
				},
				DataSource: &schemaD.SetAttribute{
					Computed: true,
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
