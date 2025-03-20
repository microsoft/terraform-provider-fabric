// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package customtypes

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
)

var _ basetypes.MapValuable = (*MapValueOf[attr.Value])(nil)

type MapValueOf[T attr.Value] struct {
	basetypes.MapValue
}

type (
	MapOfString = MapValueOf[basetypes.StringValue]
)

func (v MapValueOf[T]) Equal(o attr.Value) bool {
	other, ok := o.(MapValueOf[T])

	if !ok {
		return false
	}

	return v.MapValue.Equal(other.MapValue)
}

func (v MapValueOf[T]) Type(ctx context.Context) attr.Type {
	return NewMapTypeOf[T](ctx)
}

func NewMapValueOfNull[T attr.Value](ctx context.Context) MapValueOf[T] {
	return MapValueOf[T]{MapValue: basetypes.NewMapNull(newAttrTypeOf[T](ctx))}
}

func NewMapValueOfUnknown[T attr.Value](ctx context.Context) MapValueOf[T] {
	return MapValueOf[T]{MapValue: basetypes.NewMapUnknown(newAttrTypeOf[T](ctx))}
}

func NewMapValueOf[T attr.Value](ctx context.Context, elements map[string]attr.Value) (MapValueOf[T], diag.Diagnostics) {
	var diags diag.Diagnostics

	v, d := basetypes.NewMapValue(newAttrTypeOf[T](ctx), elements)
	diags.Append(d...)

	if diags.HasError() {
		return NewMapValueOfUnknown[T](ctx), diags
	}

	return MapValueOf[T]{MapValue: v}, diags
}

func NewMapValueOfMust[T attr.Value](ctx context.Context, elements map[string]attr.Value) MapValueOf[T] {
	return supertypes.MustDiag(NewMapValueOf[T](ctx, elements))
}

func (v MapValueOf[T]) Get(ctx context.Context) (values map[string]T, diags diag.Diagnostics) { //nolint:nonamedreturns
	values = make(map[string]T, len(v.MapValue.Elements()))

	diags.Append(v.MapValue.ElementsAs(ctx, &values, false)...)

	return
}

func (v *MapValueOf[T]) Set(ctx context.Context, elements map[string]T) diag.Diagnostics {
	var d diag.Diagnostics

	v.MapValue, d = types.MapValueFrom(ctx, v.ElementType(ctx), elements)

	return d
}

func (v MapValueOf[T]) IsKnown() bool {
	return !v.MapValue.IsNull() && !v.MapValue.IsUnknown()
}

func (v *MapValueOf[T]) SetNull(ctx context.Context) {
	v.MapValue = basetypes.NewMapNull(v.ElementType(ctx))
}

func (v *MapValueOf[T]) SetUnknown(ctx context.Context) {
	v.MapValue = basetypes.NewMapUnknown(v.ElementType(ctx))
}
