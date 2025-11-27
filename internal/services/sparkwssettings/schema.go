// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package sparkwssettings

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
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
	fabspark "github.com/microsoft/fabric-sdk-go/fabric/spark"
	superschema "github.com/orange-cloudavenue/terraform-plugin-framework-superschema"

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
			"automatic_log": superschema.SuperSingleNestedAttributeOf[automaticLogPropertiesModel]{
				Common: &schemaR.SingleNestedAttribute{
					MarkdownDescription: "Automatic Log properties.",
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
				Attributes: map[string]superschema.Attribute{
					"enabled": superschema.SuperBoolAttribute{
						Common: &schemaR.BoolAttribute{
							MarkdownDescription: "The status of the automatic log: `false` - Disabled, `true` - Enabled.",
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
				},
			},
			"environment": superschema.SuperSingleNestedAttributeOf[environmentPropertiesModel]{
				Common: &schemaR.SingleNestedAttribute{
					MarkdownDescription: "Environment properties.",
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
				Attributes: map[string]superschema.Attribute{
					"name": superschema.SuperStringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "The name of the environment.",
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
					"runtime_version": superschema.SuperStringAttribute{
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
				},
			},
			"high_concurrency": superschema.SuperSingleNestedAttributeOf[highConcurrencyPropertiesModel]{
				Common: &schemaR.SingleNestedAttribute{
					MarkdownDescription: "High Concurrency properties.",
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
				Attributes: map[string]superschema.Attribute{
					"notebook_interactive_run_enabled": superschema.SuperBoolAttribute{
						Common: &schemaR.BoolAttribute{
							MarkdownDescription: "The status of the high concurrency for notebook interactive run: `false` - Disabled, `true` - Enabled.",
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
					"notebook_pipeline_run_enabled": superschema.SuperBoolAttribute{
						Common: &schemaR.BoolAttribute{
							MarkdownDescription: "The status of the high concurrency for notebook pipeline run: `false` - Disabled, `true` - Enabled.",
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
				},
			},
			"job": superschema.SuperSingleNestedAttributeOf[jobPropertiesModel]{
				Common: &schemaR.SingleNestedAttribute{
					MarkdownDescription: "Job properties.",
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
				Attributes: map[string]superschema.Attribute{
					"conservative_job_admission_enabled": superschema.SuperBoolAttribute{
						Common: &schemaR.BoolAttribute{
							MarkdownDescription: "Reserve maximum cores for active Spark jobs. When this setting is enabled, your Fabric capacity reserves the maximum number of cores needed for active Spark jobs, ensuring job reliability by making sure that cores are available if a job scales up. When this setting is disabled, jobs are started based on the minimum number of cores needed, letting more jobs run at the same time. `false` - Disabled, `true` - Enabled.",
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
					"session_timeout_in_minutes": superschema.SuperInt32Attribute{
						Common: &schemaR.Int32Attribute{
							MarkdownDescription: "Time to terminate inactive Spark sessions. The maximum is 14 days (20160 minutes).",
						},
						Resource: &schemaR.Int32Attribute{
							Optional: true,
							Computed: true,
							Validators: []validator.Int32{
								int32validator.AtMost(20160),
							},
							PlanModifiers: []planmodifier.Int32{
								int32planmodifier.UseStateForUnknown(),
							},
						},
						DataSource: &schemaD.Int32Attribute{
							Computed: true,
						},
					},
				},
			},
			"pool": superschema.SuperSingleNestedAttributeOf[poolPropertiesModel]{
				Common: &schemaR.SingleNestedAttribute{
					MarkdownDescription: "Pool properties.",
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
				Attributes: map[string]superschema.Attribute{
					"customize_compute_enabled": superschema.SuperBoolAttribute{
						Common: &schemaR.BoolAttribute{
							MarkdownDescription: "Customize compute configurations for items. `false` - Disabled, `true` - Enabled.",
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
					"default_pool": superschema.SuperSingleNestedAttributeOf[defaultPoolPropertiesModel]{
						Common: &schemaR.SingleNestedAttribute{
							MarkdownDescription: "Default pool for workspace.",
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
						Attributes: map[string]superschema.Attribute{
							"id": superschema.SuperStringAttribute{
								Common: &schemaR.StringAttribute{
									MarkdownDescription: "The Pool ID. `00000000-0000-0000-0000-000000000000` means use the starter pool.",
									CustomType:          customtypes.UUIDType{},
								},
								Resource: &schemaR.StringAttribute{
									Optional: true,
									Computed: true,
									Validators: []validator.String{
										stringvalidator.ConflictsWith(
											path.MatchRelative().AtParent().AtName("name"),
											path.MatchRelative().AtParent().AtName("type"),
										),
									},
								},
								DataSource: &schemaD.StringAttribute{
									Computed: true,
								},
							},
							"name": superschema.SuperStringAttribute{
								Common: &schemaR.StringAttribute{
									MarkdownDescription: "The Pool name. `Starter Pool` means use the starting pool.",
								},
								Resource: &schemaR.StringAttribute{
									Optional: true,
									Computed: true,
									Validators: []validator.String{
										stringvalidator.AlsoRequires(path.MatchRelative().AtParent().AtName("type")),
										stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("id")),
									},
								},
								DataSource: &schemaD.StringAttribute{
									Computed: true,
								},
							},
							"type": superschema.SuperStringAttribute{
								Common: &schemaR.StringAttribute{
									MarkdownDescription: "The Pool type.",
									Validators: []validator.String{
										stringvalidator.OneOf(utils.ConvertEnumsToStringSlices(fabspark.PossibleCustomPoolTypeValues(), false)...),
									},
								},
								Resource: &schemaR.StringAttribute{
									Optional: true,
									Computed: true,
									Validators: []validator.String{
										stringvalidator.AlsoRequires(path.MatchRelative().AtParent().AtName("name")),
										stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("id")),
									},
								},
								DataSource: &schemaD.StringAttribute{
									Computed: true,
								},
							},
						},
					},
					"starter_pool": superschema.SuperSingleNestedAttributeOf[starterPoolPropertiesModel]{
						Common: &schemaR.SingleNestedAttribute{
							MarkdownDescription: "Starter pool for workspace. For more information about configuring starter pool, see [configuring starter pool](https://review.learn.microsoft.com/fabric/data-engineering/configure-starter-pools).",
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
						Attributes: map[string]superschema.Attribute{
							"max_node_count": superschema.SuperInt32Attribute{
								Common: &schemaR.Int32Attribute{
									MarkdownDescription: "The maximum node count.",
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
							"max_executors": superschema.SuperInt32Attribute{
								Common: &schemaR.Int32Attribute{
									MarkdownDescription: "The maximum executors count.",
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
