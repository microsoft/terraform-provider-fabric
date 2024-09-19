// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package customtypes_test

import (
	"context"
	"fmt"
	"net/url"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/attr/xattr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/path"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
)

func TestUnit_URLValidateAttribute(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		urlValue      customtypes.URL
		expectedDiags diag.Diagnostics
	}{
		"empty-struct": {
			urlValue: customtypes.URL{},
		},
		"null": {
			urlValue: customtypes.NewURLNull(),
		},
		"unknown": {
			urlValue: customtypes.NewURLUnknown(),
		},
		"valid URL - localhost": {
			urlValue: customtypes.NewURLValue("https://localhost"),
		},
		"valid URL - port": {
			urlValue: customtypes.NewURLValue("https://localhost:8080"),
		},
		"valid URL - domain": {
			urlValue: customtypes.NewURLValue("https://example.com"),
		},
		"valid URL - subdomain": {
			urlValue: customtypes.NewURLValue("https://www.example.com"),
		},
		"invalid URL - localhost no scheme": {
			urlValue: customtypes.NewURLValue("localhost"),
			expectedDiags: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test"),
					customtypes.URLTypeErrorInvalidStringHeader,
					fmt.Sprintf(customtypes.URLTypeErrorInvalidStringDetails, "localhost", "invalid URI for request"),
				),
			},
		},
		"invalid URL - domain no scheme": {
			urlValue: customtypes.NewURLValue("example.com"),
			expectedDiags: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test"),
					customtypes.URLTypeErrorInvalidStringHeader,
					fmt.Sprintf(customtypes.URLTypeErrorInvalidStringDetails, "example.com", "invalid URI for request"),
				),
			},
		},
		"invalid URL - invalid characters": {
			urlValue: customtypes.NewURLValue("https:/example.com"),
			expectedDiags: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test"),
					customtypes.URLTypeErrorInvalidStringHeader,
					fmt.Sprintf(customtypes.URLTypeErrorInvalidStringDetails, "https:/example.com", "invalid URI for request"),
				),
			},
		},
		"invalid URL - invalid only scheme": {
			urlValue: customtypes.NewURLValue("https://"),
			expectedDiags: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test"),
					customtypes.URLTypeErrorInvalidStringHeader,
					fmt.Sprintf(customtypes.URLTypeErrorInvalidStringDetails, "https://", "invalid URI for request"),
				),
			},
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			resp := xattr.ValidateAttributeResponse{}

			testCase.urlValue.ValidateAttribute(
				context.Background(),
				xattr.ValidateAttributeRequest{Path: path.Root("test")},
				&resp,
			)

			if diff := cmp.Diff(resp.Diagnostics, testCase.expectedDiags); diff != "" {
				t.Errorf("Unexpected diagnostics (-got, +expected): %s", diff)
			}
		})
	}
}

func TestUnit_URLValidateParameter(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		urlValue        customtypes.URL
		expectedFuncErr *function.FuncError
	}{
		"empty-struct": {
			urlValue: customtypes.URL{},
		},
		"null": {
			urlValue: customtypes.NewURLNull(),
		},
		"unknown": {
			urlValue: customtypes.NewURLUnknown(),
		},
		"valid URL - localhost": {
			urlValue: customtypes.NewURLValue("https://localhost"),
		},
		"valid URL - port": {
			urlValue: customtypes.NewURLValue("https://localhost:8080"),
		},
		"valid URL - domain": {
			urlValue: customtypes.NewURLValue("https://example.com"),
		},
		"valid URL - subdomain": {
			urlValue: customtypes.NewURLValue("https://www.example.com"),
		},
		"invalid URL - localhost no scheme": {
			urlValue: customtypes.NewURLValue("localhost"),
			expectedFuncErr: function.NewArgumentFuncError(
				0,
				customtypes.URLTypeErrorInvalidStringHeader+": "+fmt.Sprintf(customtypes.URLTypeErrorInvalidStringDetails, "localhost", "invalid URI for request"),
			),
		},
		"invalid URL - domain no scheme": {
			urlValue: customtypes.NewURLValue("example.com"),
			expectedFuncErr: function.NewArgumentFuncError(
				0,
				customtypes.URLTypeErrorInvalidStringHeader+": "+fmt.Sprintf(customtypes.URLTypeErrorInvalidStringDetails, "example.com", "invalid URI for request"),
			),
		},
		"invalid URL - invalid characters": {
			urlValue: customtypes.NewURLValue("https:/example.com"),
			expectedFuncErr: function.NewArgumentFuncError(
				0,
				customtypes.URLTypeErrorInvalidStringHeader+": "+fmt.Sprintf(customtypes.URLTypeErrorInvalidStringDetails, "https:/example.com", "invalid URI for request"),
			),
		},
		"invalid URL - invalid only scheme": {
			urlValue: customtypes.NewURLValue("https://"),
			expectedFuncErr: function.NewArgumentFuncError(
				0,
				customtypes.URLTypeErrorInvalidStringHeader+": "+fmt.Sprintf(customtypes.URLTypeErrorInvalidStringDetails, "https://", "invalid URI for request"),
			),
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			resp := function.ValidateParameterResponse{}

			testCase.urlValue.ValidateParameter(
				context.Background(),
				function.ValidateParameterRequest{
					Position: 0,
				},
				&resp,
			)

			if diff := cmp.Diff(resp.Error, testCase.expectedFuncErr); diff != "" {
				t.Errorf("Unexpected diagnostics (-got, +expected): %s", diff)
			}
		})
	}
}

func TestUnit_URLValueURL(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		urlValue      customtypes.URL
		expectedURL   string
		expectedDiags diag.Diagnostics
	}{
		"URL value is null ": {
			urlValue: customtypes.NewURLNull(),
			expectedDiags: diag.Diagnostics{
				diag.NewErrorDiagnostic(
					customtypes.URLTypeErrorInvalidStringHeader,
					"URL string value is null",
				),
			},
		},
		"URL value is unknown ": {
			urlValue: customtypes.NewURLUnknown(),
			expectedDiags: diag.Diagnostics{
				diag.NewErrorDiagnostic(
					customtypes.URLTypeErrorInvalidStringHeader,
					"URL string value is unknown",
				),
			},
		},
		"valid URL": {
			urlValue:    customtypes.NewURLValue("https://example.com"),
			expectedURL: "https://example.com",
		},
	}
	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			urlValue, diags := testCase.urlValue.ValueURL()

			if testCase.expectedURL != "" {
				expectedURLObj, _ := url.ParseRequestURI(testCase.expectedURL)
				expectedURL := expectedURLObj.String()

				if urlValue != expectedURL {
					t.Errorf("Unexpected difference in URL, got: %s, expected: %s", urlValue, testCase.expectedURL)
				}
			}

			if diff := cmp.Diff(diags, testCase.expectedDiags); diff != "" {
				t.Errorf("Unexpected diagnostics (-got, +expected): %s", diff)
			}
		})
	}
}
