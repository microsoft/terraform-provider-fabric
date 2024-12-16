// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package utils_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	azto "github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	"github.com/stretchr/testify/assert"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
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

func TestUnit_GetDiagsFromError(t *testing.T) {
	ctx := context.Background()

	t.Run("nil error", func(t *testing.T) {
		diags := utils.GetDiagsFromError(ctx, nil, utils.OperationRead, nil)

		assert.False(t, diags.HasError())
	})

	t.Run("fabcore.ResponseError", func(t *testing.T) {
		requestID := testhelp.RandomUUID()
		err := &fabcore.ResponseError{
			ErrorCode:  "ErrorCode",
			StatusCode: http.StatusNotFound,
			ErrorResponse: &fabcore.ErrorResponse{
				ErrorCode: azto.Ptr("ErrorCode"),
				Message:   azto.Ptr("Message"),
				RequestID: &requestID,
			},
		}
		diags := utils.GetDiagsFromError(ctx, err, utils.OperationRead, nil)

		assert.Len(t, diags, 1)
		assert.Equal(t, common.ErrorReadHeader, diags[0].Summary())
		assert.Equal(t, fmt.Sprintf("%s: %s\n\nErrorCode: %s\nRequestID: %s", common.ErrorReadDetails, "Message", "ErrorCode", requestID), diags[0].Detail())
	})

	t.Run("azcore.ResponseError", func(t *testing.T) {
		requestID := testhelp.RandomUUID()
		respBody := map[string]any{
			"errorCode": "ErrorCode",
			"message":   "Message",
			"moreDetails": []map[string]any{
				{
					"errorCode": "ErrorCodeMoreDetails",
					"message":   "MessageMoreDetails",
					"relatedResource": map[string]any{
						"resourceId":   testhelp.RandomUUID(),
						"resourceType": "ResourceType",
					},
				},
			},
			"relatedResource": map[string]any{
				"resourceId":   testhelp.RandomUUID(),
				"resourceType": "ResourceType",
			},
			"requestId": requestID,
		}
		respBodyJSON, _ := json.Marshal(respBody)

		err := &azcore.ResponseError{
			ErrorCode:  "ErrorCode",
			StatusCode: http.StatusNotFound,
			RawResponse: &http.Response{
				StatusCode: http.StatusNotFound,
				Status:     "404 Not Found",
				Body:       io.NopCloser(strings.NewReader(string(respBodyJSON))),
				Request: &http.Request{
					Method: http.MethodGet,
					URL: &url.URL{
						Scheme: "https",
						Host:   "example.com",
					},
				},
			},
		}

		diags := utils.GetDiagsFromError(ctx, err, utils.OperationRead, nil)

		assert.Len(t, diags, 1)
		assert.Equal(t, common.ErrorReadHeader, diags[0].Summary())
		assert.Equal(t, fmt.Sprintf("%s: %s\n\nErrorCode: %s\nRequestID: %s", common.ErrorReadDetails, "Message / MessageMoreDetails", "ErrorCode / ErrorCodeMoreDetails", requestID), diags[0].Detail())
	})

	t.Run("azidentity.AuthenticationFailedError", func(t *testing.T) {
		respBody := map[string]any{
			"error":             "invalid_client",
			"error_description": "AADSTS7000215: Invalid client secret is provided.",
			"error_codes":       []int{7000215},
			"timestamp":         "2024-12-12 18:00:00Z",
			"trace_id":          testhelp.RandomUUID(),
			"correlation_id":    testhelp.RandomUUID(),
			"error_uri":         "https://login.microsoftonline.com/error?code=7000215",
		}
		respBodyJSON, _ := json.Marshal(respBody)
		err := &azidentity.AuthenticationFailedError{
			RawResponse: &http.Response{
				StatusCode: http.StatusNotFound,
				Status:     "404 Not Found",
				Body:       io.NopCloser(strings.NewReader(string(respBodyJSON))),
				Request: &http.Request{
					Method: http.MethodGet,
					URL: &url.URL{
						Scheme: "https",
						Host:   "example.com",
					},
				},
			},
		}
		diags := utils.GetDiagsFromError(ctx, err, utils.OperationRead, nil)

		assert.Len(t, diags, 1)
		assert.Equal(t, "invalid_client", diags[0].Summary())
		assert.Equal(t, "AADSTS7000215: Invalid client secret is provided.\n\nErrorCode: 7000215\nErrorURI: https://login.microsoftonline.com/error?code=7000215", diags[0].Detail())
	})

	t.Run("azidentity.AuthenticationRequiredError", func(t *testing.T) {
		err := &azidentity.AuthenticationRequiredError{}
		diags := utils.GetDiagsFromError(ctx, err, utils.OperationRead, nil)

		assert.Len(t, diags, 1)
		assert.Equal(t, "authentication required", diags[0].Summary())
		assert.Equal(t, "", diags[0].Detail())
	})

	t.Run("unexpected error", func(t *testing.T) {
		err := errors.New("unexpected error")
		diags := utils.GetDiagsFromError(ctx, err, utils.OperationRead, nil)

		assert.Len(t, diags, 1)
		assert.Equal(t, "unknown error", diags[0].Summary())
		assert.Equal(t, "unexpected error", diags[0].Detail())
	})
}
