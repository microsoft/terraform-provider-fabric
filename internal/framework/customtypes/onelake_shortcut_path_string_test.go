// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package customtypes_test

import (
	"testing"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
)

func TestUnit_PathStringSemanticEquals(t *testing.T) {
	t.Parallel()

	type testCase struct {
		val1, val2 customtypes.PathString
		equals     bool
	}
	tests := map[string]testCase{
		"both with `/`, equal": {
			val1:   customtypes.NewPathStringValue("/test"),
			val2:   customtypes.NewPathStringValue("/test"),
			equals: true,
		},
		"both without `/`, equal": {
			val1:   customtypes.NewPathStringValue("test"),
			val2:   customtypes.NewPathStringValue("test"),
			equals: true,
		},
		"first with `/`, second without, equal": {
			val1:   customtypes.NewPathStringValue("/test"),
			val2:   customtypes.NewPathStringValue("test"),
			equals: true,
		},
		"first without, second with `/`, equal": {
			val1:   customtypes.NewPathStringValue("test"),
			val2:   customtypes.NewPathStringValue("/test"),
			equals: true,
		},
		"not equal": {
			val1:   customtypes.NewPathStringValue("test"),
			val2:   customtypes.NewPathStringValue("Test"),
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
