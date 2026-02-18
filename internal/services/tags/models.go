// Copyright Microsoft Corporation 2026
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
	DomainID customtypes.UUID `tfsdk:"domain_id"`
	Type     types.String     `tfsdk:"type"`
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
	switch s := from.(type) {
	case *fabadmin.DomainTagScope:
		to.Type = types.StringPointerValue((*string)(s.Type))
		to.DomainID = customtypes.NewUUIDPointerValue(s.DomainID)
	case *fabadmin.TenantTagScope:
		to.Type = types.StringPointerValue((*string)(s.Type))
	default:
		scope := from.GetTagScope()
		to.Type = types.StringPointerValue((*string)(scope.Type))
	}
}

/*
RESOURCE
*/

type resourceTagsModel struct {
	baseTagModel

	Timeouts timeoutsR.Value `tfsdk:"timeouts"`
}

type requestCreateTags struct {
	fabadmin.CreateTagsRequest
}

func (to *resourceTagsModel) setValue(ctx context.Context, from []fabadmin.Tag) diag.Diagnostics {
	to.Scope = supertypes.NewSingleNestedObjectValueOfNull[scopeModel](ctx)
	entity := from[0]

	to.ID = customtypes.NewUUIDPointerValue(entity.ID)
	to.DisplayName = types.StringPointerValue(entity.DisplayName)

	scope := &scopeModel{}
	scope.set(entity.Scope)

	if diags := to.Scope.Set(ctx, scope); diags.HasError() {
		return diags
	}

	return nil
}

func (to *baseTagModel) setTag(ctx context.Context, from fabadmin.Tag) diag.Diagnostics {
	to.ID = customtypes.NewUUIDPointerValue(from.ID)
	to.DisplayName = types.StringPointerValue(from.DisplayName)

	to.Scope = supertypes.NewSingleNestedObjectValueOfNull[scopeModel](ctx)

	var scope *scopeModel

	switch s := from.Scope.(type) {
	case *fabadmin.DomainTagScope:
		scope = &scopeModel{
			Type:     types.StringPointerValue((*string)(s.Type)),
			DomainID: customtypes.NewUUIDPointerValue(s.DomainID),
		}
	case *fabadmin.TenantTagScope:
		scope = &scopeModel{
			Type: types.StringPointerValue((*string)(s.Type)),
		}
	default:
		scope = &scopeModel{
			Type: types.StringPointerValue((*string)(from.Scope.GetTagScope().Type)),
		}
	}

	if diags := to.Scope.Set(ctx, scope); diags.HasError() {
		return diags
	}

	return nil
}

func (to *baseTagModel) setTagInfo(ctx context.Context, from fabadmin.TagInfo) diag.Diagnostics {
	to.ID = customtypes.NewUUIDPointerValue(from.ID)
	to.DisplayName = types.StringPointerValue(from.DisplayName)

	to.Scope = supertypes.NewSingleNestedObjectValueOfNull[scopeModel](ctx)

	var scope *scopeModel

	switch s := from.Scope.(type) {
	case *fabadmin.DomainTagScope:
		scope = &scopeModel{
			Type:     types.StringPointerValue((*string)(s.Type)),
			DomainID: customtypes.NewUUIDPointerValue(s.DomainID),
		}
	case *fabadmin.TenantTagScope:
		scope = &scopeModel{
			Type: types.StringPointerValue((*string)(s.Type)),
		}
	default:
		scope = &scopeModel{
			Type: types.StringPointerValue((*string)(from.Scope.GetTagScope().Type)),
		}
	}

	if diags := to.Scope.Set(ctx, scope); diags.HasError() {
		return diags
	}

	return nil
}

func (to *requestCreateTags) set(ctx context.Context, from resourceTagsModel) diag.Diagnostics {
	to.CreateTagsRequest.CreateTagsRequest = append(to.CreateTagsRequest.CreateTagsRequest, fabadmin.CreateTagRequest{
		DisplayName: from.DisplayName.ValueStringPointer(),
	})

	if from.Scope.IsKnown() {
		scope, diags := from.Scope.Get(ctx)
		if diags.HasError() {
			return diags
		}

		if scope.Type.ValueString() == string(fabadmin.TagScopeTypeTenant) {
			to.Scope = &fabadmin.TenantTagScope{
				Type: (*fabadmin.TagScopeType)(scope.Type.ValueStringPointer()),
			}
		} else if scope.Type.ValueString() == string(fabadmin.TagScopeTypeDomain) {
			to.Scope = &fabadmin.DomainTagScope{
				Type:     (*fabadmin.TagScopeType)(scope.Type.ValueStringPointer()),
				DomainID: scope.DomainID.ValueStringPointer(),
			}
		}
	}

	return nil
}

type requestUpdateTags struct {
	fabadmin.UpdateTagRequest
}

func (to *requestUpdateTags) set(from resourceTagsModel) {
	to.DisplayName = from.DisplayName.ValueStringPointer()
}
