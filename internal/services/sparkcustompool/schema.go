// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package sparkcustompool

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema" //revive:disable-line:import-alias-naming
	"github.com/hashicorp/terraform-plugin-framework/path"
	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema" //revive:disable-line:import-alias-naming
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	fabspark "github.com/microsoft/fabric-sdk-go/fabric/spark"
	superschema "github.com/orange-cloudavenue/terraform-plugin-framework-superschema"
	superint32validator "github.com/orange-cloudavenue/terraform-plugin-framework-validators/int32validator"

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
					Optional: true,
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
			"name": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The " + ItemTypeInfo.Name + " ID.",
					Validators: []validator.String{
						stringvalidator.LengthAtMost(64),
						stringvalidator.RegexMatches(regexp.MustCompile(`^[a-zA-Z0-9-_ ]+$`), "The name must contain only letters, numbers, dashes, underscores and spaces."),
					},
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
					Validators: []validator.String{
						stringvalidator.NoneOfCaseInsensitive("Starter Pool"),
					},
				},
				DataSource: &schemaD.StringAttribute{
					Optional: true,
					Computed: true,
				},
			},
			"type": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The " + ItemTypeInfo.Name + " type.",
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
					Validators: []validator.String{
						stringvalidator.OneOf(utils.ConvertEnumsToStringSlices(utils.RemoveSliceByValue(fabspark.PossibleCustomPoolTypeValues(), fabspark.CustomPoolTypeCapacity), true)...),
					},
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
					Validators: []validator.String{
						stringvalidator.OneOf(utils.ConvertEnumsToStringSlices(fabspark.PossibleCustomPoolTypeValues(), true)...),
					},
				},
			},
			"node_family": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The Node family.",
					Validators: []validator.String{
						stringvalidator.OneOf(utils.ConvertEnumsToStringSlices(fabspark.PossibleNodeFamilyValues(), true)...),
					},
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"node_size": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The Node size.",
					Validators: []validator.String{
						stringvalidator.OneOf(utils.ConvertEnumsToStringSlices(fabspark.PossibleNodeSizeValues(), true)...),
					},
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"auto_scale": superschema.SuperSingleNestedAttributeOf[sparkCustomPoolAutoScaleModel]{
				Common: &schemaR.SingleNestedAttribute{
					MarkdownDescription: "Auto-scale properties.",
				},
				Resource: &schemaR.SingleNestedAttribute{
					Required: true,
				},
				DataSource: &schemaD.SingleNestedAttribute{
					Computed: true,
				},
				Attributes: superschema.Attributes{
					"enabled": superschema.BoolAttribute{
						Common: &schemaR.BoolAttribute{
							MarkdownDescription: "The status of the auto scale: `false` - Disabled, `true` - Enabled.",
						},
						Resource: &schemaR.BoolAttribute{
							Required: true,
						},
						DataSource: &schemaD.BoolAttribute{
							Computed: true,
						},
					},
					"min_node_count": superschema.Int32Attribute{
						Common: &schemaR.Int32Attribute{
							MarkdownDescription: "The minimum node count.",
						},
						Resource: &schemaR.Int32Attribute{
							Required: true,
						},
						DataSource: &schemaD.Int32Attribute{
							Computed: true,
						},
					},
					"max_node_count": superschema.Int32Attribute{
						Common: &schemaR.Int32Attribute{
							MarkdownDescription: "The maximum node count.",
						},
						Resource: &schemaR.Int32Attribute{
							Required: true,
						},
						DataSource: &schemaD.Int32Attribute{
							Computed: true,
						},
					},
				},
			},
			"dynamic_executor_allocation": superschema.SuperSingleNestedAttributeOf[sparkCustomPoolDynamicExecutorAllocationModel]{
				Common: &schemaR.SingleNestedAttribute{
					MarkdownDescription: "Dynamic Executor Allocation properties.",
				},
				Resource: &schemaR.SingleNestedAttribute{
					Required: true,
				},
				DataSource: &schemaD.SingleNestedAttribute{
					Computed: true,
				},
				Attributes: superschema.Attributes{
					"enabled": superschema.BoolAttribute{
						Common: &schemaR.BoolAttribute{
							MarkdownDescription: "The status of the dynamic executor allocation: `false` - Disabled, `true` - Enabled.",
						},
						Resource: &schemaR.BoolAttribute{
							Required: true,
						},
						DataSource: &schemaD.BoolAttribute{
							Computed: true,
						},
					},
					"min_executors": superschema.Int32Attribute{
						Common: &schemaR.Int32Attribute{
							MarkdownDescription: "The minimum executors.",
						},
						Resource: &schemaR.Int32Attribute{
							Computed: true,
							Optional: true,
							Validators: []validator.Int32{
								superint32validator.NullIfAttributeIsOneOf(
									path.MatchRoot("dynamic_executor_allocation").AtName("enabled"),
									[]attr.Value{types.BoolValue(false)},
								),
								superint32validator.RequireIfAttributeIsOneOf(
									path.MatchRoot("dynamic_executor_allocation").AtName("enabled"),
									[]attr.Value{types.BoolValue(true)},
								),
							},
						},
						DataSource: &schemaD.Int32Attribute{
							Computed: true,
						},
					},
					"max_executors": superschema.Int32Attribute{
						Common: &schemaR.Int32Attribute{
							MarkdownDescription: "The maximum executors.",
						},
						Resource: &schemaR.Int32Attribute{
							Computed: true,
							Optional: true,
							Validators: []validator.Int32{
								superint32validator.NullIfAttributeIsOneOf(
									path.MatchRoot("dynamic_executor_allocation").AtName("enabled"),
									[]attr.Value{types.BoolValue(false)},
								),
								superint32validator.RequireIfAttributeIsOneOf(
									path.MatchRoot("dynamic_executor_allocation").AtName("enabled"),
									[]attr.Value{types.BoolValue(true)},
								),
							},
						},
						DataSource: &schemaD.Int32Attribute{
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
