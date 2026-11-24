// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package workspaceonelakeimmutability

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator" //revive:disable-line:import-alias-naming
	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema"  //revive:disable-line:import-alias-naming
	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema"    //revive:disable-line:import-alias-naming
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	superschema "github.com/orange-cloudavenue/terraform-plugin-framework-superschema"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
)

func itemSchema() superschema.Schema {
	return superschema.Schema{
		Resource: superschema.SchemaDetails{
			MarkdownDescription: fabricitem.NewResourceMarkdownDescription(
				ItemTypeInfo,
				false,
			) + "\n\nThe Fabric API does not support deleting or resetting this policy. Destroying the Terraform resource only removes it from state; the remote policy remains configured.",
		},
		DataSource: superschema.SchemaDetails{
			MarkdownDescription: fabricitem.NewDataSourceMarkdownDescription(ItemTypeInfo, false),
		},
		Attributes: superschema.Attributes{
			"workspace_id": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The workspace ID.",
					CustomType:          customtypes.UUIDType{},
					Required:            true,
				},
				Resource: &schemaR.StringAttribute{
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
				},
			},
			"scope": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The scope of the immutability policy. Currently, only `DiagnosticLogs` is supported.",
					CustomType:          basetypes.StringType{},
					Validators: []validator.String{
						stringvalidator.OneOf(utils.ConvertEnumsToStringSlices(fabcore.PossibleImmutabilityScopeValues(), true)...),
					},
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
			"retention_days": superschema.SuperInt64Attribute{
				Common: &schemaR.Int64Attribute{
					MarkdownDescription: "The number of days that diagnostic logs remain immutable. The API does not allow reducing this value.",
					CustomType:          basetypes.Int64Type{},
					Validators: []validator.Int64{
						int64validator.Between(1, maxRetentionDays),
					},
				},
				Resource: &schemaR.Int64Attribute{
					Required: true,
				},
				DataSource: &schemaD.Int64Attribute{
					Computed: true,
				},
			},
			"timeouts": superschema.TimeoutAttribute{
				Resource: &superschema.ResourceTimeoutAttribute{
					Create: true,
					Read:   true,
					Update: true,
				},
				DataSource: &superschema.DatasourceTimeoutAttribute{
					Read: true,
				},
			},
		},
	}
}
