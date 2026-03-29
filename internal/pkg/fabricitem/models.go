// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package fabricitem

import (
	"context"
	"reflect"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
)

type sensitivityLabelSettingsModel struct {
	LabelID                       customtypes.UUID `tfsdk:"label_id"`
	SensitivityLabelApplyStrategy types.String     `tfsdk:"sensitivity_label_apply_strategy"`
}

type sensitivityLabelModel struct {
	LabelID customtypes.UUID `tfsdk:"label_id"`
}

func (to *sensitivityLabelModel) set(from fabcore.SensitivityLabel) {
	to.LabelID = customtypes.NewUUIDPointerValue(from.ID)
}

type fabricItemModel struct {
	WorkspaceID customtypes.UUID `tfsdk:"workspace_id"`
	ID          customtypes.UUID `tfsdk:"id"`
	DisplayName types.String     `tfsdk:"display_name"`
	Description types.String     `tfsdk:"description"`
	FolderID    customtypes.UUID `tfsdk:"folder_id"`
}

func (to *fabricItemModel) set(from fabcore.Item) {
	to.WorkspaceID = customtypes.NewUUIDPointerValue(from.WorkspaceID)
	to.ID = customtypes.NewUUIDPointerValue(from.ID)
	to.DisplayName = types.StringPointerValue(from.DisplayName)
	to.Description = types.StringPointerValue(from.Description)
	to.FolderID = customtypes.NewUUIDPointerValue(from.FolderID)
}

type dataSourceFabricItemBaseModel struct {
	fabricItemModel

	SensitivityLabel supertypes.SingleNestedObjectValueOf[sensitivityLabelModel] `tfsdk:"sensitivity_label"`
}

func (to *dataSourceFabricItemBaseModel) set(ctx context.Context, from fabcore.Item) diag.Diagnostics {
	to.fabricItemModel.set(from)

	to.SensitivityLabel = supertypes.NewSingleNestedObjectValueOfNull[sensitivityLabelModel](ctx)

	if from.SensitivityLabel != nil {
		sensitivityLabelModel := &sensitivityLabelModel{}
		sensitivityLabelModel.set(*from.SensitivityLabel)

		if diags := to.SensitivityLabel.Set(ctx, sensitivityLabelModel); diags.HasError() {
			return diags
		}
	}

	return nil
}

type resourceFabricItemBaseModel struct {
	fabricItemModel

	SensitivityLabelSettings supertypes.SingleNestedObjectValueOf[sensitivityLabelSettingsModel] `tfsdk:"sensitivity_label_settings"`
}

type FabricItemPropertiesModel[Ttfprop, Titemprop any] struct { //revive:disable-line:exported
	WorkspaceID customtypes.UUID                              `tfsdk:"workspace_id"`
	ID          customtypes.UUID                              `tfsdk:"id"`
	DisplayName types.String                                  `tfsdk:"display_name"`
	Description types.String                                  `tfsdk:"description"`
	FolderID    customtypes.UUID                              `tfsdk:"folder_id"`
	Properties  supertypes.SingleNestedObjectValueOf[Ttfprop] `tfsdk:"properties"`
}

func (to *FabricItemPropertiesModel[Ttfprop, Titemprop]) set(from FabricItemProperties[Titemprop]) {
	to.WorkspaceID = customtypes.NewUUIDPointerValue(from.WorkspaceID)
	to.ID = customtypes.NewUUIDPointerValue(from.ID)
	to.DisplayName = types.StringPointerValue(from.DisplayName)
	to.Description = types.StringPointerValue(from.Description)
	to.FolderID = customtypes.NewUUIDPointerValue(from.FolderID)
}

type DataSourceFabricItemPropertiesBaseModel[Ttfprop, Titemprop any] struct { //revive:disable-line:exported
	FabricItemPropertiesModel[Ttfprop, Titemprop]

	SensitivityLabel supertypes.SingleNestedObjectValueOf[sensitivityLabelModel] `tfsdk:"sensitivity_label"`
}

func (to *DataSourceFabricItemPropertiesBaseModel[Ttfprop, Titemprop]) set(ctx context.Context, from FabricItemProperties[Titemprop]) diag.Diagnostics {
	to.FabricItemPropertiesModel.set(from)

	to.SensitivityLabel = supertypes.NewSingleNestedObjectValueOfNull[sensitivityLabelModel](ctx)

	if from.SensitivityLabel != nil {
		sensitivityLabelModel := &sensitivityLabelModel{}
		sensitivityLabelModel.set(*from.SensitivityLabel)

		if diags := to.SensitivityLabel.Set(ctx, sensitivityLabelModel); diags.HasError() {
			return diags
		}
	}

	return nil
}

type ResourceFabricItemPropertiesBaseModel[Ttfprop, Titemprop any] struct { //revive:disable-line:exported
	FabricItemPropertiesModel[Ttfprop, Titemprop]

	SensitivityLabelSettings supertypes.SingleNestedObjectValueOf[sensitivityLabelSettingsModel] `tfsdk:"sensitivity_label_settings"`
}

type FabricItemProperties[Titemprop any] struct { //revive:disable-line:exported
	fabcore.Item

	Properties *Titemprop
}

func (to *FabricItemProperties[Titemprop]) Set(from any) {
	fromValue := reflect.ValueOf(from)
	if fromValue.Kind() == reflect.Pointer {
		fromValue = fromValue.Elem()
	}

	to.WorkspaceID = getFieldStringValue(fromValue, "WorkspaceID")
	to.ID = getFieldStringValue(fromValue, "ID")
	to.DisplayName = getFieldStringValue(fromValue, "DisplayName")
	to.Description = getFieldStringValue(fromValue, "Description")
	to.FolderID = getFieldStringValue(fromValue, "FolderID")
	to.SensitivityLabel = getFieldStructValue[fabcore.SensitivityLabel](fromValue, "SensitivityLabel")
	to.Properties = getFieldStructValue[Titemprop](fromValue, "Properties")
}

func getFieldStringValue(v reflect.Value, fieldName string) *string {
	field := v.FieldByName(fieldName)
	if field.Kind() == reflect.Pointer {
		field = field.Elem()
	}

	if field.IsValid() && field.Kind() == reflect.String {
		if str, ok := field.Interface().(string); ok {
			return &str
		}
	}

	return nil
}

func getFieldStructValue[Titemprop any](v reflect.Value, fieldName string) *Titemprop {
	field := v.FieldByName(fieldName)
	if field.Kind() == reflect.Pointer {
		field = field.Elem()
	}

	if field.IsValid() && field.CanInterface() {
		if value, ok := field.Interface().(Titemprop); ok {
			return &value
		}
	}

	return nil
}
