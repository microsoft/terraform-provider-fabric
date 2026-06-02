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

type fabricItemModel struct {
	WorkspaceID customtypes.UUID `tfsdk:"workspace_id"`
	ID          customtypes.UUID `tfsdk:"id"`
	DisplayName types.String     `tfsdk:"display_name"`
	Description types.String     `tfsdk:"description"`
	FolderID    customtypes.UUID `tfsdk:"folder_id"`
	Tags        types.Set        `tfsdk:"tags"`
}

func (to *fabricItemModel) set(ctx context.Context, from fabcore.Item) diag.Diagnostics {
	to.WorkspaceID = customtypes.NewUUIDPointerValue(from.WorkspaceID)
	to.ID = customtypes.NewUUIDPointerValue(from.ID)
	to.DisplayName = types.StringPointerValue(from.DisplayName)
	to.Description = types.StringPointerValue(from.Description)
	to.FolderID = customtypes.NewUUIDPointerValue(from.FolderID)

	return SetResourceTagsFromItem(ctx, &to.Tags, from.Tags)
}

type FabricItemPropertiesModel[Ttfprop, Titemprop any] struct { //revive:disable-line:exported
	WorkspaceID customtypes.UUID                              `tfsdk:"workspace_id"`
	ID          customtypes.UUID                              `tfsdk:"id"`
	DisplayName types.String                                  `tfsdk:"display_name"`
	Description types.String                                  `tfsdk:"description"`
	FolderID    customtypes.UUID                              `tfsdk:"folder_id"`
	Properties  supertypes.SingleNestedObjectValueOf[Ttfprop] `tfsdk:"properties"`
	Tags        types.Set                                     `tfsdk:"tags"`
}

func (to *FabricItemPropertiesModel[Ttfprop, Titemprop]) set(ctx context.Context, from FabricItemProperties[Titemprop]) diag.Diagnostics {
	to.WorkspaceID = customtypes.NewUUIDPointerValue(from.WorkspaceID)
	to.ID = customtypes.NewUUIDPointerValue(from.ID)
	to.DisplayName = types.StringPointerValue(from.DisplayName)
	to.Description = types.StringPointerValue(from.Description)
	to.FolderID = customtypes.NewUUIDPointerValue(from.FolderID)

	return SetResourceTagsFromItem(ctx, &to.Tags, from.Tags)
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
	to.Properties = getFieldStructValue[Titemprop](fromValue, "Properties")
	to.Tags = getFieldSliceValue[fabcore.ItemTag](fromValue, "Tags")
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

func getFieldSliceValue[T any](v reflect.Value, fieldName string) []T {
	field := v.FieldByName(fieldName)
	if !field.IsValid() {
		return nil
	}

	if field.Kind() == reflect.Pointer {
		field = field.Elem()
	}

	if !field.IsValid() || field.Kind() != reflect.Slice || field.Len() == 0 {
		return nil
	}

	if value, ok := field.Interface().([]T); ok {
		return value
	}

	// Convert element-by-element for structurally identical but differently-named types
	// (e.g., fabwarehouse.ItemTag → fabcore.ItemTag)
	targetType := reflect.TypeFor[T]()
	result := make([]T, 0, field.Len())

	for i := range field.Len() {
		elem := field.Index(i)
		if elem.Kind() == reflect.Pointer {
			elem = elem.Elem()
		}

		if elem.Type().ConvertibleTo(targetType) {
			converted, ok := elem.Convert(targetType).Interface().(T)
			if ok {
				result = append(result, converted)
			}
		}
	}

	if len(result) == 0 {
		return nil
	}

	return result
}
