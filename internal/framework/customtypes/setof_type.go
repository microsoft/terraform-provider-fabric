// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package customtypes

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var _ basetypes.SetTypable = (*SetTypeOf[attr.Value])(nil)

type SetTypeOf[T attr.Value] struct {
	basetypes.SetType
}

func NewSetTypeOf[T attr.Value](ctx context.Context) SetTypeOf[T] {
	return SetTypeOf[T]{basetypes.SetType{ElemType: newAttrTypeOf[T](ctx)}}
}

func (t SetTypeOf[T]) Equal(o attr.Type) bool {
	other, ok := o.(SetTypeOf[T])

	if !ok {
		return false
	}

	return t.SetType.Equal(other.SetType)
}

func (t SetTypeOf[T]) String() string {
	var zero T

	return fmt.Sprintf("SetTypeOf[%T]", zero)
}

func (t SetTypeOf[T]) ValueFromSet(ctx context.Context, in basetypes.SetValue) (basetypes.SetValuable, diag.Diagnostics) {
	var diags diag.Diagnostics

	if in.IsNull() {
		return NewSetValueOfNull[T](ctx), diags
	}

	if in.IsUnknown() {
		return NewSetValueOfUnknown[T](ctx), diags
	}

	mapValue, d := basetypes.NewSetValue(newAttrTypeOf[T](ctx), in.Elements())
	diags.Append(d...)

	if diags.HasError() {
		return NewSetValueOfUnknown[T](ctx), diags
	}

	return SetValueOf[T]{SetValue: mapValue}, diags
}

func (t SetTypeOf[T]) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	attrValue, err := t.SetType.ValueFromTerraform(ctx, in)
	if err != nil {
		return nil, err
	}

	mapValue, ok := attrValue.(basetypes.SetValue)
	if !ok {
		return nil, fmt.Errorf("unexpected value type of %T", attrValue)
	}

	mapValuable, diags := t.ValueFromSet(ctx, mapValue)
	if diags.HasError() {
		return nil, fmt.Errorf("unexpected error converting SetValue to MapValuable: %v", diags)
	}

	return mapValuable, nil
}

func (t SetTypeOf[T]) ValueType(_ context.Context) attr.Value {
	return SetValueOf[T]{}
}

// func newAttrTypeOf[T attr.Value](ctx context.Context) attr.Type {
// 	var zero T

// 	return zero.Type(ctx)
// }
