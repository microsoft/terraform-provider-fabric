// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package fabricitem

import (
	"reflect"

	"github.com/hashicorp/terraform-plugin-framework/types"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
)

type baseFabricItemModel struct {
	WorkspaceID customtypes.UUID `tfsdk:"workspace_id"`
	ID          customtypes.UUID `tfsdk:"id"`
	DisplayName types.String     `tfsdk:"display_name"`
	Description types.String     `tfsdk:"description"`
}

func (to *baseFabricItemModel) set(from fabcore.Item) {
	to.WorkspaceID = customtypes.NewUUIDPointerValue(from.WorkspaceID)
	to.ID = customtypes.NewUUIDPointerValue(from.ID)
	to.DisplayName = types.StringPointerValue(from.DisplayName)
	to.Description = types.StringPointerValue(from.Description)
}

type baseFabricItemModel1[T any, Tm any] struct {
	WorkspaceID customtypes.UUID `tfsdk:"workspace_id"`
	ID          customtypes.UUID `tfsdk:"id"`
	DisplayName types.String     `tfsdk:"display_name"`
	Description types.String     `tfsdk:"description"`
}

func (to *baseFabricItemModel1[T, Tm]) set1(from FabricItem[Tm]) {
	to.WorkspaceID = customtypes.NewUUIDPointerValue(from.WorkspaceID)
	to.ID = customtypes.NewUUIDPointerValue(from.ID)
	to.DisplayName = types.StringPointerValue(from.DisplayName)
	to.Description = types.StringPointerValue(from.Description)
}

type FabricItem[Tm any] struct {
	fabcore.Item
	Properties *Tm
}

func (to *FabricItem[Tm]) Set(from any) {
	fromValue := reflect.ValueOf(from)
	if fromValue.Kind() == reflect.Pointer {
		fromValue = fromValue.Elem()
	}

	to.WorkspaceID = getFieldStringValue(fromValue, "WorkspaceID")
	to.ID = getFieldStringValue(fromValue, "ID")
	to.DisplayName = getFieldStringValue(fromValue, "DisplayName")
	to.Description = getFieldStringValue(fromValue, "Description")
	to.Properties = getFieldStructValue[Tm](fromValue, "Properties")
}

func getFieldStringValue(v reflect.Value, fieldName string) *string {
	field := v.FieldByName(fieldName)
	if field.Kind() == reflect.Pointer {
		field = field.Elem()
	}

	if field.IsValid() && field.Kind() == reflect.String {
		str := field.Interface().(string)

		return &str
	}

	return nil
}

func getFieldStructValue[T any](v reflect.Value, fieldName string) *T {
	field := v.FieldByName(fieldName)

	if field.IsValid() && field.CanInterface() {
		if value, ok := field.Interface().(T); ok {
			return &value
		}
	}
	return nil

	// if field.Kind() == reflect.Pointer {
	// 	field = field.Elem()
	// }

	// if field.IsValid() && field.Kind() == reflect.Struct {
	// 	return field.Interface()
	// }

	// return nil
}
