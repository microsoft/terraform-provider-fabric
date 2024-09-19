// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package customtypes

import (
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func NewUUIDNull() UUID {
	return UUID{
		StringValue: basetypes.NewStringNull(),
	}
}

func NewUUIDUnknown() UUID {
	return UUID{
		StringValue: basetypes.NewStringUnknown(),
	}
}

func NewUUIDValue(value string) UUID {
	return UUID{
		StringValue: basetypes.NewStringValue(value),
	}
}

func NewUUIDPointerValue(value *string) UUID {
	return UUID{
		StringValue: basetypes.NewStringPointerValue(value),
	}
}

func NewUUIDValueMust(value string) (UUID, diag.Diagnostics) {
	v := NewUUIDValue(value)

	_, diags := v.ValueUUID()
	if diags.HasError() {
		return UUID{}, diags
	}

	return v, nil
}

func NewUUIDPointerValueMust(value *string) (UUID, diag.Diagnostics) {
	v := NewUUIDPointerValue(value)

	_, diags := v.ValueUUID()
	if diags.HasError() {
		return UUID{}, diags
	}

	return v, nil
}
