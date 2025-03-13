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

var _ basetypes.StringTypable = (*URLType)(nil)

type URLType struct {
	basetypes.StringType
}

func (t URLType) String() string {
	return "URLType"
}

func (t URLType) Equal(o attr.Type) bool {
	other, ok := o.(URLType)
	if !ok {
		return false
	}

	return t.StringType.Equal(other.StringType)
}

func (t URLType) ValueType(_ context.Context) attr.Value {
	return URLValue{}
}

func (t URLType) ValueFromString(_ context.Context, in basetypes.StringValue) (basetypes.StringValuable, diag.Diagnostics) {
	var diags diag.Diagnostics

	if in.IsNull() {
		return NewURLNull(), diags
	}

	if in.IsUnknown() {
		return NewURLUnknown(), diags
	}

	return URLValue{
		StringValue: in,
	}, nil
}

func (t URLType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
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
