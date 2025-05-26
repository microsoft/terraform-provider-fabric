// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package customtypes

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var (
	_ basetypes.StringValuable                   = (*PathStringValue)(nil)
	_ basetypes.StringValuableWithSemanticEquals = (*PathStringValue)(nil)
)

type PathString = PathStringValue

type PathStringValue struct {
	basetypes.StringValue
}

func (PathStringValue) Type(context.Context) attr.Type {
	return PathStringType{}
}

func (v PathStringValue) Equal(o attr.Value) bool {
	other, ok := o.(PathString)
	if !ok {
		return false
	}

	return v.StringValue.Equal(other.StringValue)
}

func (v PathStringValue) StringSemanticEquals(ctx context.Context, newValuable basetypes.StringValuable) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	newValue, ok := newValuable.(PathStringValue)
	if !ok {
		diags.AddError(
			"Semantic Equality Check Error",
			"An unexpected value type was received while performing semantic equality checks. "+
				"Please report this to the provider developers.\n\n"+
				"Expected Value Type: "+fmt.Sprintf("%T", v)+"\n"+
				"Got Value Type: "+fmt.Sprintf("%T", newValuable),
		)

		return false, diags
	}

	oldV, d := v.ToStringValue(ctx)
	diags.Append(d...)

	if diags.HasError() {
		return false, diags
	}

	newV, d := newValue.ToStringValue(ctx)
	diags.Append(d...)

	if diags.HasError() {
		return false, diags
	}

	return strings.TrimPrefix(oldV.ValueString(), "/") == strings.TrimPrefix(newV.ValueString(), "/"), diags
}
