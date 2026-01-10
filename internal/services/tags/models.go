// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package tags

import (
	"context"

	timeoutsD "github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts" //revive:disable-line:import-alias-naming
	timeoutsR "github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"   //revive:disable-line:import-alias-naming
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	fabadmin "github.com/microsoft/fabric-sdk-go/fabric/admin"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
)

/*
BASE MODEL
*/

type baseTagModel struct {
	ID          customtypes.UUID                                 `tfsdk:"id"`
	DisplayName types.String                                     `tfsdk:"display_name"`
	Scope       supertypes.SingleNestedObjectValueOf[scopeModel] `tfsdk:"scope"`
}

type scopeModel struct {
	Type types.String `tfsdk:"type"`
}

type dataSourceTagsModel struct {
	Values   supertypes.SetNestedObjectValueOf[baseTagModel] `tfsdk:"values"`
	Timeouts timeoutsD.Value                                 `tfsdk:"timeouts"`
}

type dataSourceTagModel struct {
	baseTagModel

	Timeouts timeoutsD.Value `tfsdk:"timeouts"`
}

func (to *dataSourceTagsModel) setValues(ctx context.Context, from []fabadmin.TagInfo) diag.Diagnostics {
	slice := make([]*baseTagModel, 0, len(from))

	for _, entity := range from {
		var entityModel baseTagModel

		if diags := entityModel.set(ctx, entity); diags.HasError() {
			return diags
		}

		slice = append(slice, &entityModel)
	}

	return to.Values.Set(ctx, slice)
}

func (to *baseTagModel) set(ctx context.Context, from fabadmin.TagInfo) diag.Diagnostics {
	to.ID = customtypes.NewUUIDPointerValue(from.ID)
	to.DisplayName = types.StringPointerValue(from.DisplayName)

	to.Scope = supertypes.NewSingleNestedObjectValueOfNull[scopeModel](ctx)
	scope := supertypes.NewSingleNestedObjectValueOfNull[scopeModel](ctx)

	if from.Scope != nil {
		scopeModel := &scopeModel{}
		scopeModel.set(from.Scope)

		if diags := scope.Set(ctx, scopeModel); diags.HasError() {
			return diags
		}
	}

	to.Scope = scope

	return nil
}

func (to *scopeModel) set(from fabadmin.TagScopeClassification) {
	scope := from.GetTagScope()
	to.Type = types.StringPointerValue((*string)(scope.Type))
}

/*
RESOURCE
*/

type resourceTagsModel struct {
	baseTagModel

	Tags supertypes.ListNestedObjectValueOf[baseTagModel] `tfsdk:"tags"`

	Timeouts timeoutsR.Value `tfsdk:"timeouts"`
}

type requestCreateTags struct {
	fabadmin.CreateTagsRequest
}

func (to *resourceTagsModel) set(ctx context.Context, from []fabadmin.Tag) diag.Diagnostics {
	to.Tags = supertypes.NewListNestedObjectValueOfNull[baseTagModel](ctx)
	to.Scope = supertypes.NewSingleNestedObjectValueOfNull[scopeModel](ctx)

	slice := make([]*baseTagModel, 0, len(from))

	for _, entity := range from {
		item := &baseTagModel{}
		if diags := item.setValue(ctx, entity); diags.HasError() {
			return diags
		}

		slice = append(slice, item)
	}

	return to.Tags.Set(ctx, slice)
}

func (to *baseTagModel) setValue(ctx context.Context, from fabadmin.Tag) diag.Diagnostics {
	to.ID = customtypes.NewUUIDPointerValue(from.ID)
	to.DisplayName = types.StringPointerValue(from.DisplayName)

	to.Scope = supertypes.NewSingleNestedObjectValueOfNull[scopeModel](ctx)

	scope := &scopeModel{
		Type: types.StringPointerValue((*string)(from.Scope.GetTagScope().Type)),
	}
	if diags := to.Scope.Set(ctx, scope); diags.HasError() {
		return diags
	}

	return nil
}

func (to *requestCreateTags) set(ctx context.Context, from resourceTagsModel) diag.Diagnostics {
	tags, diags := from.Tags.Get(ctx)
	if diags.HasError() {
		return diags
	}

	for _, tag := range tags {
		to.CreateTagsRequest.CreateTagsRequest = append(to.CreateTagsRequest.CreateTagsRequest, fabadmin.CreateTagRequest{
			DisplayName: tag.DisplayName.ValueStringPointer(),
		})
	}

	return nil
}

type requestUpdateTags struct {
	fabadmin.UpdateTagRequest
}

func (to *requestUpdateTags) set(from resourceTagsModel) {
	to.DisplayName = from.DisplayName.ValueStringPointer()
}
