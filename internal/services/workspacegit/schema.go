// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package workspacegit

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema" //revive:disable-line:import-alias-naming
	"github.com/hashicorp/terraform-plugin-framework/path"
	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema" //revive:disable-line:import-alias-naming
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	superschema "github.com/orange-cloudavenue/terraform-plugin-framework-superschema"
	superobjectvalidator "github.com/orange-cloudavenue/terraform-plugin-framework-validators/objectvalidator"
	superstringvalidator "github.com/orange-cloudavenue/terraform-plugin-framework-validators/stringvalidator"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
)

func itemSchema() superschema.Schema { //nolint:maintidx
	gitProviderTypeAttPath := path.MatchRoot("git_provider_details").AtName("git_provider_type")
	gitProviderTypeAzureDevOps := types.StringValue(string(fabcore.GitProviderTypeAzureDevOps))
	gitProviderTypeGitHub := types.StringValue(string(fabcore.GitProviderTypeGitHub))
	possibleInitializationStrategyValues := utils.RemoveSliceByValue(fabcore.PossibleInitializationStrategyValues(), fabcore.InitializationStrategyNone)

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
			"initialization_strategy": superschema.StringAttribute{
				Resource: &schemaR.StringAttribute{
					MarkdownDescription: "The initialization strategy.",
					Required:            true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
					Validators: []validator.String{
						stringvalidator.OneOf(utils.ConvertEnumsToStringSlices(possibleInitializationStrategyValues, true)...),
					},
				},
			},
			"git_connection_state": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The Git connection state",
					Validators: []validator.String{
						stringvalidator.OneOf(utils.ConvertEnumsToStringSlices(fabcore.PossibleGitConnectionStateValues(), true)...),
					},
				},
				Resource: &schemaR.StringAttribute{
					Computed: true,
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"git_sync_details": superschema.SuperSingleNestedAttributeOf[gitSyncDetailsModel]{
				Common: &schemaR.SingleNestedAttribute{
					MarkdownDescription: "The Git sync details.",
				},
				Resource: &schemaR.SingleNestedAttribute{
					Computed: true,
				},
				DataSource: &schemaD.SingleNestedAttribute{
					Computed: true,
				},
				Attributes: map[string]superschema.Attribute{
					"head": superschema.StringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "The git head.",
						},
						Resource: &schemaR.StringAttribute{
							Computed: true,
						},
						DataSource: &schemaD.StringAttribute{
							Computed: true,
						},
					},
					"last_sync_time": superschema.SuperStringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "The last sync time.",
							CustomType:          timetypes.RFC3339Type{},
						},
						Resource: &schemaR.StringAttribute{
							Computed: true,
						},
						DataSource: &schemaD.StringAttribute{
							Computed: true,
						},
					},
				},
			},
			"git_provider_details": superschema.SuperSingleNestedAttributeOf[gitProviderDetailsModel]{
				Common: &schemaR.SingleNestedAttribute{
					MarkdownDescription: "The Git provider details.",
				},
				Resource: &schemaR.SingleNestedAttribute{
					Required: true,
					PlanModifiers: []planmodifier.Object{
						objectplanmodifier.RequiresReplace(),
					},
				},
				DataSource: &schemaD.SingleNestedAttribute{
					Computed: true,
				},
				Attributes: map[string]superschema.Attribute{
					"git_provider_type": superschema.StringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "The git provider type.",
							Validators: []validator.String{
								stringvalidator.OneOf(utils.ConvertEnumsToStringSlices(fabcore.PossibleGitProviderTypeValues(), true)...),
							},
						},
						Resource: &schemaR.StringAttribute{
							Required: true,
						},
						DataSource: &schemaD.StringAttribute{
							Computed: true,
						},
					},
					"organization_name": superschema.StringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "The Azure DevOps organization name.",
						},
						Resource: &schemaR.StringAttribute{
							Computed: true,
							Optional: true,
							Validators: []validator.String{
								stringvalidator.LengthAtMost(100),
								superstringvalidator.NullIfAttributeIsOneOf(
									gitProviderTypeAttPath,
									[]attr.Value{gitProviderTypeGitHub},
								),
								superstringvalidator.RequireIfAttributeIsOneOf(
									gitProviderTypeAttPath,
									[]attr.Value{gitProviderTypeAzureDevOps},
								),
							},
						},
						DataSource: &schemaD.StringAttribute{
							Computed: true,
						},
					},
					"project_name": superschema.StringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "The Azure DevOps project name.",
						},
						Resource: &schemaR.StringAttribute{
							Computed: true,
							Optional: true,
							Validators: []validator.String{
								stringvalidator.LengthAtMost(100),
								superstringvalidator.NullIfAttributeIsOneOf(
									gitProviderTypeAttPath,
									[]attr.Value{gitProviderTypeGitHub},
								),
								superstringvalidator.RequireIfAttributeIsOneOf(
									gitProviderTypeAttPath,
									[]attr.Value{gitProviderTypeAzureDevOps},
								),
							},
						},
						DataSource: &schemaD.StringAttribute{
							Computed: true,
						},
					},
					"owner_name": superschema.StringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "The GitHub owner name.",
						},
						Resource: &schemaR.StringAttribute{
							Computed: true,
							Optional: true,
							Validators: []validator.String{
								stringvalidator.LengthAtMost(100),
								superstringvalidator.NullIfAttributeIsOneOf(
									gitProviderTypeAttPath,
									[]attr.Value{gitProviderTypeAzureDevOps},
								),
								superstringvalidator.RequireIfAttributeIsOneOf(
									gitProviderTypeAttPath,
									[]attr.Value{gitProviderTypeGitHub},
								),
							},
						},
						DataSource: &schemaD.StringAttribute{
							Computed: true,
						},
					},
					"repository_name": superschema.StringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "The repository name.",
						},
						Resource: &schemaR.StringAttribute{
							Required: true,
							Validators: []validator.String{
								stringvalidator.LengthAtMost(128),
							},
						},
						DataSource: &schemaD.StringAttribute{
							Computed: true,
						},
					},
					"branch_name": superschema.StringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "The branch name.",
						},
						Resource: &schemaR.StringAttribute{
							Required: true,
							Validators: []validator.String{
								stringvalidator.LengthAtMost(250),
							},
						},
						DataSource: &schemaD.StringAttribute{
							Computed: true,
						},
					},
					"directory_name": superschema.StringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "The directory name.",
						},
						Resource: &schemaR.StringAttribute{
							Required: true,
							Validators: []validator.String{
								stringvalidator.LengthAtMost(256),
								stringvalidator.RegexMatches(
									regexp.MustCompile(`^/.*`),
									"Directory name path must starts with forward slash '/'.",
								),
							},
						},
						DataSource: &schemaD.StringAttribute{
							Computed: true,
						},
					},
				},
			},
			"git_credentials": superschema.SuperSingleNestedAttributeOf[gitCredentialsModel]{
				Common: &schemaR.SingleNestedAttribute{
					MarkdownDescription: "The Git credentials details.",
				},
				Resource: &schemaR.SingleNestedAttribute{
					Computed: true,
					Optional: true,
					Validators: []validator.Object{
						superobjectvalidator.RequireIfAttributeIsOneOf(
							gitProviderTypeAttPath,
							[]attr.Value{gitProviderTypeGitHub},
						),
					},
				},
				DataSource: &schemaD.SingleNestedAttribute{
					Computed: true,
				},
				Attributes: map[string]superschema.Attribute{
					"source": superschema.StringAttribute{
						Common: &schemaR.StringAttribute{
							Validators: []validator.String{
								stringvalidator.OneOf(utils.ConvertEnumsToStringSlices(fabcore.PossibleGitCredentialsSourceValues(), true)...),
							},
						},
						Resource: &schemaR.StringAttribute{
							MarkdownDescription: "The Git credentials source. If the value of `git_provider_details.git_provider_type` attribute is `GitHub` this attribute MUST be `ConfiguredConnection`. If the value of `git_provider_details.git_provider_type` attribute is `AzureDevOps` this attribute MUST be one of `ConfiguredConnection`, `Automatic`, and it defaults to `Automatic`.",
							Computed:            true,
							Optional:            true,
						},
						DataSource: &schemaD.StringAttribute{
							MarkdownDescription: "The Git credentials source.",
							Computed:            true,
						},
					},
					"connection_id": superschema.SuperStringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "The connection ID.",
							CustomType:          customtypes.UUIDType{},
						},
						Resource: &schemaR.StringAttribute{
							Computed: true,
							Optional: true,
							Validators: []validator.String{
								superstringvalidator.RequireIfAttributeIsOneOf(
									gitProviderTypeAttPath,
									[]attr.Value{gitProviderTypeGitHub},
								),
								superstringvalidator.RequireIfAttributeIsOneOf(
									path.MatchRoot("git_credentials").AtName("source"),
									[]attr.Value{types.StringValue(string(fabcore.GitCredentialsSourceConfiguredConnection))},
								),
								superstringvalidator.NullIfAttributeIsOneOf(
									path.MatchRoot("git_credentials").AtName("source"),
									[]attr.Value{types.StringValue(string(fabcore.GitCredentialsSourceAutomatic))},
								),
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
