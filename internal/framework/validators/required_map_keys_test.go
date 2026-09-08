// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package validators_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	fwvalidators "github.com/microsoft/terraform-provider-fabric/internal/framework/validators"
)

func TestUnit_RequiredMapKeysValidator(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		value    types.Map
		required []string
		hasError bool
	}{
		"null": {
			value:    types.MapNull(types.StringType),
			required: []string{"definition.json"},
		},
		"unknown": {
			value:    types.MapUnknown(types.StringType),
			required: []string{"definition.json"},
		},
		"required key present": {
			value: types.MapValueMust(types.StringType, map[string]attr.Value{
				"definition.json": types.StringValue("content"),
			}),
			required: []string{"definition.json"},
		},
		"required key missing": {
			value: types.MapValueMust(types.StringType, map[string]attr.Value{
				".platform": types.StringValue("content"),
			}),
			required: []string{"definition.json"},
			hasError: true,
		},
		"one of multiple required keys missing": {
			value: types.MapValueMust(types.StringType, map[string]attr.Value{
				"definition.json": types.StringValue("content"),
			}),
			required: []string{"definition.json", ".platform"},
			hasError: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			resp := &validator.MapResponse{}
			fwvalidators.RequiredMapKeys(test.required...).ValidateMap(t.Context(), validator.MapRequest{
				ConfigValue: test.value,
				Path:        path.Root("definition"),
			}, resp)

			if test.hasError != resp.Diagnostics.HasError() {
				t.Fatalf("expected HasError() to be %t, got diagnostics: %s", test.hasError, resp.Diagnostics)
			}
		})
	}
}
