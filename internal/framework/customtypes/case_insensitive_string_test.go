// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package customtypes_test

import (
	"testing"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
)

func TestUnit_CaseInsensitiveStringSemanticEquals(t *testing.T) {
	t.Parallel()

	type testCase struct {
		val1, val2 customtypes.CaseInsensitiveString
		equals     bool
	}
	tests := map[string]testCase{
		"both lowercase, equal": {
			val1:   customtypes.NewCaseInsensitiveStringValue("test"),
			val2:   customtypes.NewCaseInsensitiveStringValue("test"),
			equals: true,
		},
		"both uppercase, equal": {
			val1:   customtypes.NewCaseInsensitiveStringValue("TEST"),
			val2:   customtypes.NewCaseInsensitiveStringValue("TEST"),
			equals: true,
		},
		"first uppercase, second lowercase, equal": {
			val1:   customtypes.NewCaseInsensitiveStringValue("TEST"),
			val2:   customtypes.NewCaseInsensitiveStringValue("test"),
			equals: true,
		},
		"first lowercase, second uppercase, equal": {
			val1:   customtypes.NewCaseInsensitiveStringValue("test"),
			val2:   customtypes.NewCaseInsensitiveStringValue("TEST"),
			equals: true,
		},
		"not equal": {
			val1:   customtypes.NewCaseInsensitiveStringValue("Test1"),
			val2:   customtypes.NewCaseInsensitiveStringValue("Test2"),
			equals: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()

			equals, _ := test.val1.StringSemanticEquals(ctx, test.val2)

			if got, want := equals, test.equals; got != want {
				t.Errorf("StringSemanticEquals(%q, %q) = %v, want %v", test.val1, test.val2, got, want)
			}
		})
	}
}
