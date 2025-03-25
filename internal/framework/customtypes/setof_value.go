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

var _ basetypes.SetValuable = (*SetValueOf[attr.Value])(nil)

type SetValueOf[T attr.Value] struct {
	basetypes.SetValue
}

func (v SetValueOf[T]) Equal(o attr.Value) bool {
	other, ok := o.(SetValueOf[T])

	if !ok {
		return false
	}

	return v.SetValue.Equal(other.SetValue)
}

func (v SetValueOf[T]) Type(ctx context.Context) attr.Type {
	return NewSetTypeOf[T](ctx)
}

func NewSetValueOfNull[T attr.Value](ctx context.Context) SetValueOf[T] {
	return SetValueOf[T]{SetValue: basetypes.NewSetNull(newAttrTypeOf[T](ctx))}
}

func NewSetValueOfUnknown[T attr.Value](ctx context.Context) SetValueOf[T] {
	return SetValueOf[T]{SetValue: basetypes.NewSetUnknown(newAttrTypeOf[T](ctx))}
}

// NewSetValueOfSlice returns a new SetValueOf with the given slice value.
func NewSetValueOfSlice[T attr.Value](ctx context.Context, elements []T) (SetValueOf[T], diag.Diagnostics) {
	return newSetValueOf[T](ctx, elements)
}

// NewSetValueOfSlicePtr returns a new SetValueOf with the given slice value.
func NewSetValueOfSlicePtr[T attr.Value](ctx context.Context, elements []*T) (SetValueOf[T], diag.Diagnostics) {
	return newSetValueOf[T](ctx, elements)
}

func newSetValueOf[T attr.Value](ctx context.Context, elements any) (SetValueOf[T], diag.Diagnostics) {
	var diags diag.Diagnostics

	v, d := basetypes.NewSetValueFrom(ctx, newAttrTypeOf[T](ctx), elements)
	diags.Append(d...)

	if diags.HasError() {
		return NewSetValueOfUnknown[T](ctx), diags
	}

	return SetValueOf[T]{SetValue: v}, diags
}

func NewSetValueOfMust[T attr.Value](ctx context.Context, elements []T) SetValueOf[T] {
	return supertypes.MustDiag(newSetValueOf[T](ctx, elements))
}

func (v SetValueOf[T]) Get(ctx context.Context) (values []T, diags diag.Diagnostics) { //nolint:nonamedreturns
	values = make([]T, len(v.Elements()))

	diags.Append(v.ElementsAs(ctx, &values, false)...)

	return
}

func (v *SetValueOf[T]) Set(ctx context.Context, elements []T) diag.Diagnostics {
	var d diag.Diagnostics

	v.SetValue, d = types.SetValueFrom(ctx, v.ElementType(ctx), elements)

	return d
}

func (v SetValueOf[T]) IsKnown() bool {
	return !v.IsNull() && !v.IsUnknown()
}

func (v *SetValueOf[T]) SetNull(ctx context.Context) {
	v.SetValue = basetypes.NewSetNull(v.ElementType(ctx))
}

func (v *SetValueOf[T]) SetUnknown(ctx context.Context) {
	v.SetValue = basetypes.NewSetUnknown(v.ElementType(ctx))
}
