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

var _ basetypes.MapTypable = (*MapTypeOf[attr.Value])(nil)

var MapOfStringType = MapTypeOf[basetypes.StringValue]{basetypes.MapType{ElemType: basetypes.StringType{}}} //nolint:gochecknoglobals

type MapTypeOf[T attr.Value] struct {
	basetypes.MapType
}

func NewMapTypeOf[T attr.Value](ctx context.Context) MapTypeOf[T] {
	return MapTypeOf[T]{basetypes.MapType{ElemType: newAttrTypeOf[T](ctx)}}
}

func (t MapTypeOf[T]) Equal(o attr.Type) bool {
	other, ok := o.(MapTypeOf[T])

	if !ok {
		return false
	}

	return t.MapType.Equal(other.MapType)
}

func (t MapTypeOf[T]) String() string {
	var zero T

	return fmt.Sprintf("MapTypeOf[%T]", zero)
}

func (t MapTypeOf[T]) ValueFromMap(ctx context.Context, in basetypes.MapValue) (basetypes.MapValuable, diag.Diagnostics) {
	var diags diag.Diagnostics

	if in.IsNull() {
		return NewMapValueOfNull[T](ctx), diags
	}

	if in.IsUnknown() {
		return NewMapValueOfUnknown[T](ctx), diags
	}

	mapValue, d := basetypes.NewMapValue(newAttrTypeOf[T](ctx), in.Elements())
	diags.Append(d...)

	if diags.HasError() {
		return NewMapValueOfUnknown[T](ctx), diags
	}

	return MapValueOf[T]{MapValue: mapValue}, diags
}

func (t MapTypeOf[T]) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	attrValue, err := t.MapType.ValueFromTerraform(ctx, in)
	if err != nil {
		return nil, err
	}

	mapValue, ok := attrValue.(basetypes.MapValue)
	if !ok {
		return nil, fmt.Errorf("unexpected value type of %T", attrValue)
	}

	mapValuable, diags := t.ValueFromMap(ctx, mapValue)
	if diags.HasError() {
		return nil, fmt.Errorf("unexpected error converting MapValue to MapValuable: %v", diags)
	}

	return mapValuable, nil
}

func (t MapTypeOf[T]) ValueType(_ context.Context) attr.Value {
	return MapValueOf[T]{}
}

func newAttrTypeOf[T attr.Value](ctx context.Context) attr.Type {
	var zero T

	return zero.Type(ctx)
}
