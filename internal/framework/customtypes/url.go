// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package customtypes

import (
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func NewURLNull() URL {
	return URL{
		StringValue: basetypes.NewStringNull(),
	}
}

func NewURLUnknown() URL {
	return URL{
		StringValue: basetypes.NewStringUnknown(),
	}
}

func NewURLValue(value string) URL {
	return URL{
		StringValue: basetypes.NewStringValue(value),
	}
}

func NewURLPointerValue(value *string) URL {
	return URL{
		StringValue: basetypes.NewStringPointerValue(value),
	}
}

func NewURLValueMust(value string) (URL, diag.Diagnostics) {
	v := NewURLValue(value)

	_, diags := v.ValueURL()
	if diags.HasError() {
		return URL{}, diags
	}

	return v, nil
}

func NewURLPointerValueMust(value *string) (URL, diag.Diagnostics) {
	v := NewURLPointerValue(value)

	_, diags := v.ValueURL()
	if diags.HasError() {
		return URL{}, diags
	}

	return v, nil
}
