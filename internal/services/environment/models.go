// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package environment

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	fabenvironment "github.com/microsoft/fabric-sdk-go/fabric/environment"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
)

type environmentPropertiesModel struct {
	PublishDetails supertypes.SingleNestedObjectValueOf[environmentPublishDetailsModel] `tfsdk:"publish_details"`
}

func (to *environmentPropertiesModel) set(ctx context.Context, from fabenvironment.Properties) diag.Diagnostics {
	publishDetails := supertypes.NewSingleNestedObjectValueOfNull[environmentPublishDetailsModel](ctx)

	if from.PublishDetails != nil {
		publishDetailsModel := &environmentPublishDetailsModel{}

		if diags := publishDetailsModel.set(ctx, *from.PublishDetails); diags.HasError() {
			return diags
		}

		if diags := publishDetails.Set(ctx, publishDetailsModel); diags.HasError() {
			return diags
		}
	}

	to.PublishDetails = publishDetails

	return nil
}

type environmentPublishDetailsModel struct {
	State                types.String                                                               `tfsdk:"state"`
	TargetVersion        customtypes.UUID                                                           `tfsdk:"target_version"`
	StartTime            timetypes.RFC3339                                                          `tfsdk:"start_time"`
	EndTime              timetypes.RFC3339                                                          `tfsdk:"end_time"`
	ComponentPublishInfo supertypes.SingleNestedObjectValueOf[environmentComponentPublishInfoModel] `tfsdk:"component_publish_info"`
}

func (to *environmentPublishDetailsModel) set(ctx context.Context, from fabenvironment.PublishDetails) diag.Diagnostics {
	to.State = types.StringPointerValue((*string)(from.State))
	to.TargetVersion = customtypes.NewUUIDPointerValue(from.TargetVersion)
	to.StartTime = timetypes.NewRFC3339TimePointerValue(from.StartTime)
	to.EndTime = timetypes.NewRFC3339TimePointerValue(from.EndTime)

	componentPublishInfo := supertypes.NewSingleNestedObjectValueOfNull[environmentComponentPublishInfoModel](ctx)

	if from.ComponentPublishInfo != nil {
		publishDetailsModel := &environmentComponentPublishInfoModel{}

		if diags := publishDetailsModel.set(ctx, *from.ComponentPublishInfo); diags.HasError() {
			return diags
		}

		if diags := componentPublishInfo.Set(ctx, publishDetailsModel); diags.HasError() {
			return diags
		}
	}

	to.ComponentPublishInfo = componentPublishInfo

	return nil
}

type environmentComponentPublishInfoModel struct {
	SparkLibraries supertypes.SingleNestedObjectValueOf[environmentSparkLibrariesModel] `tfsdk:"spark_libraries"`
	SparkSettings  supertypes.SingleNestedObjectValueOf[environmentSparkSettingsModel]  `tfsdk:"spark_settings"`
}

func (to *environmentComponentPublishInfoModel) set(ctx context.Context, from fabenvironment.ComponentPublishInfo) diag.Diagnostics {
	sparkLibraries := supertypes.NewSingleNestedObjectValueOfNull[environmentSparkLibrariesModel](ctx)

	if from.SparkLibraries != nil {
		sparkLibrariesModel := &environmentSparkLibrariesModel{}
		sparkLibrariesModel.set(*from.SparkLibraries)

		if diags := sparkLibraries.Set(ctx, sparkLibrariesModel); diags.HasError() {
			return diags
		}
	}

	to.SparkLibraries = sparkLibraries

	sparkSettings := supertypes.NewSingleNestedObjectValueOfNull[environmentSparkSettingsModel](ctx)

	if from.SparkSettings != nil {
		sparkSettingsModel := &environmentSparkSettingsModel{}
		sparkSettingsModel.set(*from.SparkSettings)

		if diags := sparkSettings.Set(ctx, sparkSettingsModel); diags.HasError() {
			return diags
		}
	}

	to.SparkSettings = sparkSettings

	return nil
}

type environmentSparkLibrariesModel struct {
	State types.String `tfsdk:"state"`
}

func (to *environmentSparkLibrariesModel) set(from fabenvironment.SparkLibraries) {
	to.State = types.StringPointerValue((*string)(from.State))
}

type environmentSparkSettingsModel struct {
	State types.String `tfsdk:"state"`
}

func (to *environmentSparkSettingsModel) set(from fabenvironment.SparkSettings) {
	to.State = types.StringPointerValue((*string)(from.State))
}
