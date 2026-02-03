// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package sparkenvsettings

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema" //revive:disable-line:import-alias-naming
	"github.com/hashicorp/terraform-plugin-framework/path"
	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema" //revive:disable-line:import-alias-naming
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	fabenvironment "github.com/microsoft/fabric-sdk-go/fabric/environment"
	superschema "github.com/orange-cloudavenue/terraform-plugin-framework-superschema"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
)

func itemSchema() superschema.Schema { //nolint:maintidx
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
			"environment_id": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The Environment ID.",
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
			"publication_status": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "Publication status.",
					Validators: []validator.String{
						stringvalidator.OneOf(SparkEnvironmentPublicationStatusValues...),
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
			"driver_cores": superschema.Int32Attribute{
				Common: &schemaR.Int32Attribute{
					MarkdownDescription: "Publication status.",
					Validators: []validator.Int32{
						int32validator.OneOf(SparkEnvironmentDriverCoresValues...),
					},
				},
				Resource: &schemaR.Int32Attribute{
					Optional: true,
					Computed: true,
					PlanModifiers: []planmodifier.Int32{
						int32planmodifier.UseStateForUnknown(),
					},
				},
				DataSource: &schemaD.Int32Attribute{
					Computed: true,
				},
			},
			"driver_memory": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "Spark driver memory.",
					Validators: []validator.String{
						stringvalidator.OneOf(SparkEnvironmentDriverMemoryValues...),
					},
				},
				Resource: &schemaR.StringAttribute{
					Optional: true,
					Computed: true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"executor_cores": superschema.Int32Attribute{
				Common: &schemaR.Int32Attribute{
					MarkdownDescription: "Spark executor core.",
					Validators: []validator.Int32{
						int32validator.OneOf(SparkEnvironmentExecutorCoresValues...),
					},
				},
				Resource: &schemaR.Int32Attribute{
					Optional: true,
					Computed: true,
					PlanModifiers: []planmodifier.Int32{
						int32planmodifier.UseStateForUnknown(),
					},
				},
				DataSource: &schemaD.Int32Attribute{
					Computed: true,
				},
			},
			"executor_memory": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "Spark executor memory.",
					Validators: []validator.String{
						stringvalidator.OneOf(SparkEnvironmentExecutorMemoryValues...),
					},
				},
				Resource: &schemaR.StringAttribute{
					Optional: true,
					Computed: true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"runtime_version": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "[Runtime](https://review.learn.microsoft.com/fabric/data-engineering/runtime) version.",
					Validators: []validator.String{
						stringvalidator.OneOf(SparkRuntimeVersionValues...),
					},
				},
				Resource: &schemaR.StringAttribute{
					Optional: true,
					Computed: true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"spark_properties": superschema.SuperMapAttribute{
				Common: &schemaR.MapAttribute{
					MarkdownDescription: "A map of key/value pairs of Spark properties.",
					CustomType:          supertypes.MapTypeOf[types.String]{MapType: types.MapType{ElemType: types.StringType}},
					ElementType:         types.StringType,
				},
				Resource: &schemaR.MapAttribute{
					Optional: true,
					Computed: true,
					Validators: []validator.Map{
						mapvalidator.KeysAre(stringvalidator.RegexMatches(
							regexp.MustCompile(`^spark\.[a-zA-Z0-9]+([\.]?[a-zA-Z0-9]+)*$`),
							"Spark properties:\n"+
								"- must starts with 'spark.'\n"+
								"- cannot contains any white spaces\n"+
								"- dot '.' is allowed but not at the start or end of the property key",
						)),
					},
				},
				DataSource: &schemaD.MapAttribute{
					Computed: true,
				},
			},
			"dynamic_executor_allocation": superschema.SuperSingleNestedAttributeOf[dynamicExecutorAllocationPropertiesModel]{
				Common: &schemaR.SingleNestedAttribute{
					MarkdownDescription: "Dynamic Executor Allocation properties.",
				},
				Resource: &schemaR.SingleNestedAttribute{
					Optional: true,
					Computed: true,
					PlanModifiers: []planmodifier.Object{
						objectplanmodifier.UseStateForUnknown(),
					},
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
							Optional: true,
							Computed: true,
							PlanModifiers: []planmodifier.Bool{
								boolplanmodifier.UseStateForUnknown(),
							},
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
							Optional: true,
							Computed: true,
							PlanModifiers: []planmodifier.Int32{
								int32planmodifier.UseStateForUnknown(),
							},
							Validators: []validator.Int32{
								int32validator.AtLeast(1),
								int32validator.AlsoRequires(path.MatchRelative().AtParent().AtName("max_executors")),
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
							Optional: true,
							Computed: true,
							PlanModifiers: []planmodifier.Int32{
								int32planmodifier.UseStateForUnknown(),
							},
							Validators: []validator.Int32{
								int32validator.AtLeast(1),
								int32validator.AlsoRequires(path.MatchRelative().AtParent().AtName("min_executors")),
							},
						},
						DataSource: &schemaD.Int32Attribute{
							Computed: true,
						},
					},
				},
			},
			"pool": superschema.SuperSingleNestedAttributeOf[instancePoolPropertiesModel]{
				Common: &schemaR.SingleNestedAttribute{
					MarkdownDescription: "Environment pool.",
				},
				Resource: &schemaR.SingleNestedAttribute{
					Optional: true,
					Computed: true,
					PlanModifiers: []planmodifier.Object{
						objectplanmodifier.UseStateForUnknown(),
					},
				},
				DataSource: &schemaD.SingleNestedAttribute{
					Computed: true,
				},
				Attributes: superschema.Attributes{
					"id": superschema.SuperStringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "The Pool ID.",
							CustomType:          customtypes.UUIDType{},
						},
						Resource: &schemaR.StringAttribute{
							Computed: true,
						},
						DataSource: &schemaD.StringAttribute{
							Computed: true,
						},
					},
					"name": superschema.StringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "The Pool name. `Starter Pool` means using the starting pool.",
						},
						Resource: &schemaR.StringAttribute{
							Optional: true,
							Computed: true,
							Validators: []validator.String{
								stringvalidator.AlsoRequires(path.MatchRelative().AtParent().AtName("type")),
							},
						},
						DataSource: &schemaD.StringAttribute{
							Computed: true,
						},
					},
					"type": superschema.StringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "The Pool type.",
							Validators: []validator.String{
								stringvalidator.OneOf(utils.ConvertEnumsToStringSlices(fabenvironment.PossibleCustomPoolTypeValues(), true)...),
							},
						},
						Resource: &schemaR.StringAttribute{
							Optional: true,
							Computed: true,
							Validators: []validator.String{
								stringvalidator.AlsoRequires(path.MatchRelative().AtParent().AtName("name")),
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
