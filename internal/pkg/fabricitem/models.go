// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package fabricitem

import (
	"reflect"

	"github.com/hashicorp/terraform-plugin-framework/types"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
)

type fabricItemModel struct {
	WorkspaceID customtypes.UUID `tfsdk:"workspace_id"`
	ID          customtypes.UUID `tfsdk:"id"`
	DisplayName types.String     `tfsdk:"display_name"`
	Description types.String     `tfsdk:"description"`
}

func (to *fabricItemModel) set(from fabcore.Item) {
	to.WorkspaceID = customtypes.NewUUIDPointerValue(from.WorkspaceID)
	to.ID = customtypes.NewUUIDPointerValue(from.ID)
	to.DisplayName = types.StringPointerValue(from.DisplayName)
	to.Description = types.StringPointerValue(from.Description)
}

type FabricItemPropertiesModel[Ttfprop, Titemprop any] struct { //revive:disable-line:exported
	WorkspaceID customtypes.UUID                              `tfsdk:"workspace_id"`
	ID          customtypes.UUID                              `tfsdk:"id"`
	DisplayName types.String                                  `tfsdk:"display_name"`
	Description types.String                                  `tfsdk:"description"`
	Properties  supertypes.SingleNestedObjectValueOf[Ttfprop] `tfsdk:"properties"`
}

func (to *FabricItemPropertiesModel[Ttfprop, Titemprop]) set(from FabricItemProperties[Titemprop]) {
	to.WorkspaceID = customtypes.NewUUIDPointerValue(from.WorkspaceID)
	to.ID = customtypes.NewUUIDPointerValue(from.ID)
	to.DisplayName = types.StringPointerValue(from.DisplayName)
	to.Description = types.StringPointerValue(from.Description)
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
