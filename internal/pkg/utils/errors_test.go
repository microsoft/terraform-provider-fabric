// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package utils_test

import (
	"errors"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/stretchr/testify/assert"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
)

func TestUnit_IsErr(t *testing.T) {
	t.Parallel()

	var diags diag.Diagnostics

	diags.AddError(
		"Error",
		"Error",
	)

	testCases := map[string]struct {
		diagnostics diag.Diagnostics
		err         string
		assertion   assert.BoolAssertionFunc
	}{
		"no error": {
			diagnostics: diags,
			err:         "This is an error",
			assertion:   assert.False,
		},
		"has error": {
			diagnostics: diags,
			err:         "Error",
			assertion:   assert.True,
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			result := utils.IsErr(testCase.diagnostics, errors.New(testCase.err))
			testCase.assertion(t, result, "they should be equal")
		})
	}
}

func TestUnit_IsErrNotFound(t *testing.T) {
	resourceID := "resource-id"
	var diags diag.Diagnostics

	diags.AddError(
		"Error1",
		"Error1",
	)

	diags.AddError(
		"Error2",
		"Error2",
	)

	testCases := map[string]struct {
		diagnostics diag.Diagnostics
		err         string
		expected    bool
	}{
		"no error": {
			diagnostics: diags,
			err:         "",
			expected:    false,
		},
		"error not found": {
			diagnostics: diags,
			err:         "not found",
			expected:    false,
		},
		"error found": {
			diagnostics: diags,
			err:         "Error1",
			expected:    true,
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			diags := &testCase.diagnostics
			result := utils.IsErrNotFound(resourceID, diags, errors.New(testCase.err))
			assert.Equal(t, testCase.expected, result, "they should be equal")

			if testCase.expected {
				assert.Len(t, *diags, 1, "diagnostics should contain one warning")
			} else {
				assert.Len(t, *diags, 2, "diagnostics should contain two errors")
			}
		})
	}
}
