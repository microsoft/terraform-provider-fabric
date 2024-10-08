// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package environment

import (
	"context"

	supertypes "github.com/FrangipaneTeam/terraform-plugin-framework-supertypes"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	fabenvironment "github.com/microsoft/fabric-sdk-go/fabric/environment"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
)

type baseEnvironmentModel struct {
	WorkspaceID customtypes.UUID `tfsdk:"workspace_id"`
	ID          customtypes.UUID `tfsdk:"id"`
	DisplayName types.String     `tfsdk:"display_name"`
	Description types.String     `tfsdk:"description"`
}

func (to *baseEnvironmentModel) set(from fabenvironment.Environment) {
	to.WorkspaceID = customtypes.NewUUIDPointerValue(from.WorkspaceID)
	to.ID = customtypes.NewUUIDPointerValue(from.ID)
	to.DisplayName = types.StringPointerValue(from.DisplayName)
	to.Description = types.StringPointerValue(from.Description)
}

type baseEnvironmentPropertiesModel struct {
	baseEnvironmentModel
	Properties supertypes.SingleNestedObjectValueOf[environmentPropertiesModel] `tfsdk:"properties"`
}

func (to *baseEnvironmentPropertiesModel) setProperties(ctx context.Context, from fabenvironment.Environment) diag.Diagnostics {
	properties := supertypes.NewSingleNestedObjectValueOfNull[environmentPropertiesModel](ctx)

	if from.Properties != nil {
		propertiesModel := &environmentPropertiesModel{}
		propertiesModel.set(ctx, from.Properties)

		if diags := properties.Set(ctx, propertiesModel); diags.HasError() {
			return diags
		}
	}

	to.Properties = properties

	return nil
}

type environmentPropertiesModel struct {
	PublishDetails supertypes.SingleNestedObjectValueOf[environmentPublishDetailsModel] `tfsdk:"publish_details"`
}

func (to *environmentPropertiesModel) set(ctx context.Context, from *fabenvironment.PublishInfo) {
	publishDetails := supertypes.NewSingleNestedObjectValueOfNull[environmentPublishDetailsModel](ctx)

	if from.PublishDetails != nil {
		publishDetailsModel := &environmentPublishDetailsModel{}
		publishDetailsModel.set(ctx, from.PublishDetails)
		publishDetails.Set(ctx, publishDetailsModel)
	}

	to.PublishDetails = publishDetails
}

type environmentPublishDetailsModel struct {
	State                types.String                                                               `tfsdk:"state"`
	TargetVersion        customtypes.UUID                                                           `tfsdk:"target_version"`
	StartTime            timetypes.RFC3339                                                          `tfsdk:"start_time"`
	EndTime              timetypes.RFC3339                                                          `tfsdk:"end_time"`
	ComponentPublishInfo supertypes.SingleNestedObjectValueOf[environmentComponentPublishInfoModel] `tfsdk:"component_publish_info"`
}

func (to *environmentPublishDetailsModel) set(ctx context.Context, from *fabenvironment.PublishDetails) {
	to.State = types.StringPointerValue((*string)(from.State))
	to.TargetVersion = customtypes.NewUUIDPointerValue(from.TargetVersion)
	to.StartTime = timetypes.NewRFC3339TimePointerValue(from.StartTime)
	to.EndTime = timetypes.NewRFC3339TimePointerValue(from.EndTime)

	componentPublishInfo := supertypes.NewSingleNestedObjectValueOfNull[environmentComponentPublishInfoModel](ctx)

	if from.ComponentPublishInfo != nil {
		publishDetailsModel := &environmentComponentPublishInfoModel{}
		publishDetailsModel.set(ctx, from.ComponentPublishInfo)
		componentPublishInfo.Set(ctx, publishDetailsModel)
	}

	to.ComponentPublishInfo = componentPublishInfo
}

type environmentComponentPublishInfoModel struct {
	SparkLibraries supertypes.SingleNestedObjectValueOf[environmentSparkLibrariesModel] `tfsdk:"spark_libraries"`
	SparkSettings  supertypes.SingleNestedObjectValueOf[environmentSparkSettingsModel]  `tfsdk:"spark_settings"`
}

func (to *environmentComponentPublishInfoModel) set(ctx context.Context, from *fabenvironment.ComponentPublishInfo) {
	sparkLibraries := supertypes.NewSingleNestedObjectValueOfNull[environmentSparkLibrariesModel](ctx)

	if from.SparkLibraries != nil {
		sparkLibrariesModel := &environmentSparkLibrariesModel{}
		sparkLibrariesModel.set(from.SparkLibraries)
		sparkLibraries.Set(ctx, sparkLibrariesModel)
	}

	to.SparkLibraries = sparkLibraries

	sparkSettings := supertypes.NewSingleNestedObjectValueOfNull[environmentSparkSettingsModel](ctx)

	if from.SparkSettings != nil {
		sparkSettingsModel := &environmentSparkSettingsModel{}
		sparkSettingsModel.set(from.SparkSettings)
		sparkSettings.Set(ctx, sparkSettingsModel)
	}

	to.SparkSettings = sparkSettings
}

type environmentSparkLibrariesModel struct {
	State types.String `tfsdk:"state"`
}

func (to *environmentSparkLibrariesModel) set(from *fabenvironment.SparkLibraries) {
	to.State = types.StringPointerValue((*string)(from.State))
}

type environmentSparkSettingsModel struct {
	State types.String `tfsdk:"state"`
}

func (to *environmentSparkSettingsModel) set(from *fabenvironment.SparkSettings) {
	to.State = types.StringPointerValue((*string)(from.State))
}
