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

var _ basetypes.StringTypable = (*CaseInsensitiveStringType)(nil)

type CaseInsensitiveStringType struct {
	basetypes.StringType
}

func (CaseInsensitiveStringType) String() string {
	return "CaseInsensitiveStringType"
}

func (t CaseInsensitiveStringType) Equal(o attr.Type) bool {
	other, ok := o.(CaseInsensitiveStringType)
	if !ok {
		return false
	}

	return t.StringType.Equal(other.StringType)
}

func (CaseInsensitiveStringType) ValueType(context.Context) attr.Value {
	return CaseInsensitiveString{}
}

func (t CaseInsensitiveStringType) ValueFromString(_ context.Context, in basetypes.StringValue) (basetypes.StringValuable, diag.Diagnostics) {
	var diags diag.Diagnostics

	if in.IsNull() {
		return NewCaseInsensitiveStringNull(), diags
	}

	if in.IsUnknown() {
		return NewCaseInsensitiveStringUnknown(), diags
	}

	return CaseInsensitiveStringValue{
		StringValue: in,
	}, diags
}

func (t CaseInsensitiveStringType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
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
