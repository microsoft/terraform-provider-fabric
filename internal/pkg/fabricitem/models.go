// Copyright Microsoft Corporation 2026
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
	FolderID    customtypes.UUID `tfsdk:"folder_id"`
}

func (to *fabricItemModel) set(from fabcore.Item) {
	to.WorkspaceID = customtypes.NewUUIDPointerValue(from.WorkspaceID)
	to.ID = customtypes.NewUUIDPointerValue(from.ID)
	to.DisplayName = types.StringPointerValue(from.DisplayName)
	to.Description = types.StringPointerValue(from.Description)
	to.FolderID = customtypes.NewUUIDPointerValue(from.FolderID)
}

type FabricItemPropertiesModel[Ttfprop, Titemprop any] struct { //revive:disable-line:exported
	WorkspaceID customtypes.UUID                              `tfsdk:"workspace_id"`
	ID          customtypes.UUID                              `tfsdk:"id"`
	DisplayName types.String                                  `tfsdk:"display_name"`
	Description types.String                                  `tfsdk:"description"`
	FolderID    customtypes.UUID                              `tfsdk:"folder_id"`
	Properties  supertypes.SingleNestedObjectValueOf[Ttfprop] `tfsdk:"properties"`
	Tags        supertypes.SetValueOf[customtypes.UUID]       `tfsdk:"tags"`
}

func (to *FabricItemPropertiesModel[Ttfprop, Titemprop]) set(from FabricItemProperties[Titemprop]) {
	to.WorkspaceID = customtypes.NewUUIDPointerValue(from.WorkspaceID)
	to.ID = customtypes.NewUUIDPointerValue(from.ID)
	to.DisplayName = types.StringPointerValue(from.DisplayName)
	to.Description = types.StringPointerValue(from.Description)
	to.FolderID = customtypes.NewUUIDPointerValue(from.FolderID)
}

// ItemTagInfo is a provider-local DTO for the read-only `Tags` field exposed by the
// Fabric items API. It is populated reflectively from the source item type's `Tags`
// slice (e.g. `fabcore.Item.Tags` or `fabenvironment.Environment.Tags`), avoiding
// cross-package type assertions on the SDK's per-package `ItemTag` types.
type ItemTagInfo struct {
	ID          *string
	DisplayName *string
}

type FabricItemProperties[Titemprop any] struct { //revive:disable-line:exported
	fabcore.Item

	Properties *Titemprop
	Tags       []ItemTagInfo
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
	to.Tags = getFieldTagsValue(fromValue, "Tags")
}

// getFieldTagsValue reflectively reads a `Tags` slice field from the source struct and
// converts each element into ItemTagInfo by reading its `ID` and `DisplayName` string
// fields. Silently returns nil when the field is missing or its shape doesn't match —
// keeping the helper safe across SDK packages that each define their own `ItemTag`
// struct (e.g. `fabcore.ItemTag`, `fabenvironment.ItemTag`) without a shared interface.
func getFieldTagsValue(v reflect.Value, fieldName string) []ItemTagInfo {
	field := v.FieldByName(fieldName)
	if !field.IsValid() || field.Kind() != reflect.Slice {
		return nil
	}

	if field.Len() == 0 {
		return nil
	}

	tags := make([]ItemTagInfo, 0, field.Len())

	for i := 0; i < field.Len(); i++ {
		elem := field.Index(i)
		if elem.Kind() == reflect.Pointer {
			elem = elem.Elem()
		}

		if !elem.IsValid() || elem.Kind() != reflect.Struct {
			continue
		}

		tags = append(tags, ItemTagInfo{
			ID:          getFieldStringValue(elem, "ID"),
			DisplayName: getFieldStringValue(elem, "DisplayName"),
		})
	}

	return tags
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
