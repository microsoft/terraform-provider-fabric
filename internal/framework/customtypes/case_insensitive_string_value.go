// Copyright Microsoft Corporation 2026
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
	_ basetypes.StringValuable                   = (*CaseInsensitiveStringValue)(nil)
	_ basetypes.StringValuableWithSemanticEquals = (*CaseInsensitiveStringValue)(nil)
)

type CaseInsensitiveString = CaseInsensitiveStringValue

type CaseInsensitiveStringValue struct {
	basetypes.StringValue
}

func (CaseInsensitiveStringValue) Type(context.Context) attr.Type {
	return CaseInsensitiveStringType{}
}

func (v CaseInsensitiveStringValue) Equal(o attr.Value) bool {
	other, ok := o.(CaseInsensitiveString)
	if !ok {
		return false
	}

	return v.StringValue.Equal(other.StringValue)
}

func (v CaseInsensitiveStringValue) StringSemanticEquals(ctx context.Context, newValuable basetypes.StringValuable) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	newValue, ok := newValuable.(CaseInsensitiveStringValue)
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

	return strings.EqualFold(oldV.ValueString(), newV.ValueString()), diags
}
