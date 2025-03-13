// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package customtypes

import (
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func NewCaseInsensitiveStringNull() CaseInsensitiveStringValue {
	return CaseInsensitiveString{
		StringValue: basetypes.NewStringNull(),
	}
}

func NewCaseInsensitiveStringUnknown() CaseInsensitiveString {
	return CaseInsensitiveString{
		StringValue: basetypes.NewStringUnknown(),
	}
}

func NewCaseInsensitiveStringValue(value string) CaseInsensitiveString {
	return CaseInsensitiveString{
		StringValue: basetypes.NewStringValue(value),
	}
}

func NewCaseInsensitiveStringPointerValue(value *string) CaseInsensitiveString {
	return CaseInsensitiveString{
		StringValue: basetypes.NewStringPointerValue(value),
	}
}
