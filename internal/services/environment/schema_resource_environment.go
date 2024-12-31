// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package environment

import (
	"context"

	supertypes "github.com/FrangipaneTeam/terraform-plugin-framework-supertypes"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	fabenvironment "github.com/microsoft/fabric-sdk-go/fabric/environment"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
)

func getResourceEnvironmentPropertiesAttributes(ctx context.Context) map[string]schema.Attribute {
	publishStatePossibleValuesMarkdown := "Publish state. Possible values: " + utils.ConvertStringSlicesToString(fabenvironment.PossibleEnvironmentPublishStateValues(), true, true) + "."

	result := map[string]schema.Attribute{
		"publish_details": schema.SingleNestedAttribute{
			MarkdownDescription: "Environment publish operation details.",
			Computed:            true,
			CustomType:          supertypes.NewSingleNestedObjectTypeOf[environmentPublishDetailsModel](ctx),
			Attributes: map[string]schema.Attribute{
				"state": schema.StringAttribute{
					MarkdownDescription: publishStatePossibleValuesMarkdown,
					Computed:            true,
				},
				"target_version": schema.StringAttribute{
					MarkdownDescription: "Target version to be published.",
					Computed:            true,
					CustomType:          customtypes.UUIDType{},
				},
				"start_time": schema.StringAttribute{
					MarkdownDescription: "Start time of publish operation.",
					Computed:            true,
					CustomType:          timetypes.RFC3339Type{},
				},
				"end_time": schema.StringAttribute{
					MarkdownDescription: "End time of publish operation.",
					Computed:            true,
					CustomType:          timetypes.RFC3339Type{},
				},
				"component_publish_info": schema.SingleNestedAttribute{
					MarkdownDescription: "Environment component publish information.",
					Computed:            true,
					CustomType:          supertypes.NewSingleNestedObjectTypeOf[environmentComponentPublishInfoModel](ctx),
					Attributes: map[string]schema.Attribute{
						"spark_libraries": schema.SingleNestedAttribute{
							MarkdownDescription: "Spark libraries publish information.",
							Computed:            true,
							CustomType:          supertypes.NewSingleNestedObjectTypeOf[environmentSparkLibrariesModel](ctx),
							Attributes: map[string]schema.Attribute{
								"state": schema.StringAttribute{
									MarkdownDescription: publishStatePossibleValuesMarkdown,
									Computed:            true,
								},
							},
						},
						"spark_settings": schema.SingleNestedAttribute{
							MarkdownDescription: "Spark settings publish information.",
							Computed:            true,
							CustomType:          supertypes.NewSingleNestedObjectTypeOf[environmentSparkSettingsModel](ctx),
							Attributes: map[string]schema.Attribute{
								"state": schema.StringAttribute{
									MarkdownDescription: publishStatePossibleValuesMarkdown,
									Computed:            true,
								},
							},
						},
					},
				},
			},
		},
	}

	return result
}
