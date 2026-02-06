// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package customtypes_test

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-framework/attr/xattr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/path"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
)

func TestUnit_UUIDValidateAttribute(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		uuidValue     customtypes.UUID
		expectedDiags diag.Diagnostics
	}{
		"empty-struct": {
			uuidValue: customtypes.UUID{},
		},
		"null": {
			uuidValue: customtypes.NewUUIDNull(),
		},
		"unknown": {
			uuidValue: customtypes.NewUUIDUnknown(),
		},
		"valid UUID v1": {
			uuidValue: customtypes.NewUUIDValue("f40ac97c-6641-11ef-bc09-eae371e4bb1e"),
		},
		"valid UUID v4": {
			uuidValue: customtypes.NewUUIDValue("e429c573-efd5-403d-888d-b92b1b9efaf5"),
		},
		"valid UUID v7": {
			uuidValue: customtypes.NewUUIDValue("01919fbe-6c52-78d3-bc27-acb170c701da"),
		},
		"invalid UUID - no hyphen": {
			uuidValue: customtypes.NewUUIDValue("04b2abf753d946cca681f1c64802abf0"),
			expectedDiags: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test"),
					customtypes.UUIDTypeErrorInvalidStringHeader,
					fmt.Sprintf(customtypes.UUIDTypeErrorInvalidStringDetails, "04b2abf753d946cca681f1c64802abf0", "uuid string is wrong length"),
				),
			},
		},
		"invalid UUID - underscore": {
			uuidValue: customtypes.NewUUIDValue("4ea3c6b1_d119_4020_b54e_7a8706222613"),
			expectedDiags: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test"),
					customtypes.UUIDTypeErrorInvalidStringHeader,
					fmt.Sprintf(customtypes.UUIDTypeErrorInvalidStringDetails, "4ea3c6b1_d119_4020_b54e_7a8706222613", "uuid is improperly formatted"),
				),
			},
		},
		"invalid UUID - invalid characters": {
			uuidValue: customtypes.NewUUIDValue("dd52ef4c.b623-4ef8-8c3b-626d87b69d28"),
			expectedDiags: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test"),
					customtypes.UUIDTypeErrorInvalidStringHeader,
					fmt.Sprintf(customtypes.UUIDTypeErrorInvalidStringDetails, "dd52ef4c.b623-4ef8-8c3b-626d87b69d28", "uuid is improperly formatted"),
				),
			},
		},
		"invalid UUID - invalid length": {
			uuidValue: customtypes.NewUUIDValue("a8f514f9-8fa0-4b93-886d-3425919fc67"),
			expectedDiags: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test"),
					customtypes.UUIDTypeErrorInvalidStringHeader,
					fmt.Sprintf(customtypes.UUIDTypeErrorInvalidStringDetails, "a8f514f9-8fa0-4b93-886d-3425919fc67", "uuid string is wrong length"),
				),
			},
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			resp := xattr.ValidateAttributeResponse{}

			testCase.uuidValue.ValidateAttribute(
				t.Context(),
				xattr.ValidateAttributeRequest{Path: path.Root("test")},
				&resp,
			)

			if diff := cmp.Diff(resp.Diagnostics, testCase.expectedDiags); diff != "" {
				t.Errorf("Unexpected diagnostics (-got, +expected): %s", diff)
			}
		})
	}
}

func TestUnit_UUIDValidateParameter(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		uuidValue       customtypes.UUID
		expectedFuncErr *function.FuncError
	}{
		"empty-struct": {
			uuidValue: customtypes.UUID{},
		},
		"null": {
			uuidValue: customtypes.NewUUIDNull(),
		},
		"unknown": {
			uuidValue: customtypes.NewUUIDUnknown(),
		},
		"valid UUID v1": {
			uuidValue: customtypes.NewUUIDValue("f40ac97c-6641-11ef-bc09-eae371e4bb1e"),
		},
		"valid UUID v4": {
			uuidValue: customtypes.NewUUIDValue("e429c573-efd5-403d-888d-b92b1b9efaf5"),
		},
		"valid UUID v7": {
			uuidValue: customtypes.NewUUIDValue("01919fbe-6c52-78d3-bc27-acb170c701da"),
		},
		"invalid UUID - no hyphen": {
			uuidValue: customtypes.NewUUIDValue("04b2abf753d946cca681f1c64802abf0"),
			expectedFuncErr: function.NewArgumentFuncError(
				0,
				customtypes.UUIDTypeErrorInvalidStringHeader+": "+fmt.Sprintf(customtypes.UUIDTypeErrorInvalidStringDetails, "04b2abf753d946cca681f1c64802abf0", "uuid string is wrong length"),
			),
		},
		"invalid UUID - underscore": {
			uuidValue: customtypes.NewUUIDValue("4ea3c6b1_d119_4020_b54e_7a8706222613"),
			expectedFuncErr: function.NewArgumentFuncError(
				0,
				customtypes.UUIDTypeErrorInvalidStringHeader+": "+fmt.Sprintf(customtypes.UUIDTypeErrorInvalidStringDetails, "4ea3c6b1_d119_4020_b54e_7a8706222613", "uuid is improperly formatted"),
			),
		},
		"invalid UUID - invalid characters": {
			uuidValue: customtypes.NewUUIDValue("dd52ef4c.b623-4ef8-8c3b-626d87b69d28"),
			expectedFuncErr: function.NewArgumentFuncError(
				0,
				customtypes.UUIDTypeErrorInvalidStringHeader+": "+fmt.Sprintf(customtypes.UUIDTypeErrorInvalidStringDetails, "dd52ef4c.b623-4ef8-8c3b-626d87b69d28", "uuid is improperly formatted"),
			),
		},
		"invalid UUID - invalid length": {
			uuidValue: customtypes.NewUUIDValue("a8f514f9-8fa0-4b93-886d-3425919fc67"),
			expectedFuncErr: function.NewArgumentFuncError(
				0,
				customtypes.UUIDTypeErrorInvalidStringHeader+": "+fmt.Sprintf(customtypes.UUIDTypeErrorInvalidStringDetails, "a8f514f9-8fa0-4b93-886d-3425919fc67", "uuid string is wrong length"),
			),
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			resp := function.ValidateParameterResponse{}

			testCase.uuidValue.ValidateParameter(
				t.Context(),
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

func TestUnit_UUIDValueUUID(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		uuidValue     customtypes.UUID
		expectedUUID  string
		expectedDiags diag.Diagnostics
	}{
		"UUID value is null ": {
			uuidValue: customtypes.NewUUIDNull(),
			expectedDiags: diag.Diagnostics{
				diag.NewErrorDiagnostic(
					customtypes.UUIDTypeErrorInvalidStringHeader,
					"UUID string value is null",
				),
			},
		},
		"UUID value is unknown ": {
			uuidValue: customtypes.NewUUIDUnknown(),
			expectedDiags: diag.Diagnostics{
				diag.NewErrorDiagnostic(
					customtypes.UUIDTypeErrorInvalidStringHeader,
					"UUID string value is unknown",
				),
			},
		},
		"valid UUID ": {
			uuidValue:    customtypes.NewUUIDValue("6977f7bd-266f-4c24-a51c-0c353185267a"),
			expectedUUID: "6977f7bd-266f-4c24-a51c-0c353185267a",
		},
	}
	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			uuidValue, diags := testCase.uuidValue.ValueUUID()
			expectedUUID, _ := uuid.ParseUUID(testCase.expectedUUID)

			if uuidValue != (string)(expectedUUID) {
				t.Errorf("Unexpected difference in UUID, got: %s, expected: %s", uuidValue, testCase.expectedUUID)
			}

			if diff := cmp.Diff(diags, testCase.expectedDiags); diff != "" {
				t.Errorf("Unexpected diagnostics (-got, +expected): %s", diff)
			}
		})
	}
}
