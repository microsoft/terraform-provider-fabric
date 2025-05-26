// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package customtypes

import (
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func NewPathStringNull() PathStringValue {
	return PathStringValue{
		StringValue: basetypes.NewStringNull(),
	}
}

func NewPathStringUnknown() PathStringValue {
	return PathStringValue{
		StringValue: basetypes.NewStringUnknown(),
	}
}

func NewPathStringValue(value string) PathStringValue {
	return PathStringValue{
		StringValue: basetypes.NewStringValue(value),
	}
}

func NewPathStringPointerValue(value *string) PathStringValue {
	return PathStringValue{
		StringValue: basetypes.NewStringPointerValue(value),
	}
}
