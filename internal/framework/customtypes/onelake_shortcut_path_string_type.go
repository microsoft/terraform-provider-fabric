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

var _ basetypes.StringTypable = (*PathStringType)(nil)

type PathStringType struct {
	basetypes.StringType
}

func (PathStringType) String() string {
	return "PathType"
}

func (t PathStringType) Equal(o attr.Type) bool {
	other, ok := o.(PathStringType)
	if !ok {
		return false
	}

	return t.StringType.Equal(other.StringType)
}

func (PathStringType) ValueType(context.Context) attr.Value {
	return PathString{}
}

func (t PathStringType) ValueFromString(_ context.Context, in basetypes.StringValue) (basetypes.StringValuable, diag.Diagnostics) {
	var diags diag.Diagnostics

	if in.IsNull() {
		return NewPathStringNull(), diags
	}

	if in.IsUnknown() {
		return NewPathStringUnknown(), diags
	}

	return PathStringValue{
		StringValue: in,
	}, diags
}

func (t PathStringType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	attrValue, err := t.StringType.ValueFromTerraform(ctx, in)
	if err != nil {
		return nil, err
	}

	stringValue, ok := attrValue.(basetypes.StringValue)
	if !ok {
		return nil, fmt.Errorf("unexpected value type of %T", attrValue)
	}

	stringValuable, diags := t.ValueFromString(ctx, stringValue)
	if diags.HasError() {
		return nil, fmt.Errorf("unexpected error converting StringValue to StringValuable: %v", diags)
	}

	return stringValuable, nil
}
