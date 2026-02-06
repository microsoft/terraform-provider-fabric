// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package utils_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
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
	ctx := t.Context()

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
		assert.Equal(t, fmt.Sprintf("%s: %s\n\nError Code: %s\nRequest ID: %s", common.ErrorReadDetails, "Message", "ErrorCode", requestID), diags[0].Detail())
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
		assert.Equal(
			t,
			fmt.Sprintf("%s: %s\n\nError Code: %s\nRequest ID: %s", common.ErrorReadDetails, "Message / MessageMoreDetails", "ErrorCode / ErrorCodeMoreDetails", requestID),
			diags[0].Detail(),
		)
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
		assert.Empty(t, diags[0].Detail())
	})

	t.Run("unexpected error", func(t *testing.T) {
		err := errors.New("unexpected error")
		diags := utils.GetDiagsFromError(ctx, err, utils.OperationRead, nil)

		assert.Len(t, diags, 1)
		assert.Equal(t, "unknown error", diags[0].Summary())
		assert.Equal(t, "unexpected error", diags[0].Detail())
	})
}

func TestUnit_HasError(t *testing.T) {
	handler := utils.NewErrorHandler()
	testErr := errors.New("test error")

	tests := []struct {
		name     string
		diags    diag.Diagnostics
		err      error
		expected bool
	}{
		{
			name:     "no errors in diagnostics",
			diags:    diag.Diagnostics{},
			err:      testErr,
			expected: false,
		},
		{
			name: "error matches",
			diags: func() diag.Diagnostics {
				var d diag.Diagnostics
				d.AddError(testErr.Error(), testErr.Error())

				return d
			}(),
			err:      testErr,
			expected: true,
		},
		{
			name: "error doesn't match",
			diags: func() diag.Diagnostics {
				var d diag.Diagnostics
				d.AddError("different error", "different error")

				return d
			}(),
			err:      testErr,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := handler.HasError(tt.diags, tt.err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestUnit_IsNotFoundError(t *testing.T) {
	handler := utils.NewErrorHandler()
	testErr := errors.New("resource not found")

	tests := []struct {
		name       string
		resourceID string
		diags      *diag.Diagnostics
		err        error
		expected   bool
		checkDiags bool
	}{
		{
			name:       "no errors in diagnostics",
			resourceID: "test-resource",
			diags:      &diag.Diagnostics{},
			err:        testErr,
			expected:   false,
			checkDiags: false,
		},
		{
			name:       "nil error",
			resourceID: "test-resource",
			diags: func() *diag.Diagnostics {
				var d diag.Diagnostics
				d.AddError("some error", "some error")

				return &d
			}(),
			err:        nil,
			expected:   false,
			checkDiags: false,
		},
		{
			name:       "error matches",
			resourceID: "test-resource",
			diags: func() *diag.Diagnostics {
				var d diag.Diagnostics
				d.AddError(testErr.Error(), testErr.Error())

				return &d
			}(),
			err:        testErr,
			expected:   true,
			checkDiags: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := handler.IsNotFoundError(tt.resourceID, tt.diags, tt.err)
			assert.Equal(t, tt.expected, result)

			if tt.checkDiags {
				// Check if diagnostics was updated with warning
				warnings := tt.diags.Warnings()
				assert.Len(t, warnings, 1)
				assert.Contains(t, warnings[0].Summary(), "Resource not found")
				assert.Contains(t, warnings[0].Detail(), tt.resourceID)
			}
		})
	}
}

func TestUnit_GetDiagsFromError_OperationMessages(t *testing.T) {
	handler := utils.NewErrorHandler()
	ctx := t.Context()
	testErr := errors.New("test error")

	tests := []struct {
		name           string
		operation      utils.Operation
		expectSummary  string
		expectDetailIn string
	}{
		{
			name:           "create operation",
			operation:      utils.OperationCreate,
			expectSummary:  "unknown error",
			expectDetailIn: testErr.Error(),
		},
		{
			name:           "read operation",
			operation:      utils.OperationRead,
			expectSummary:  "unknown error",
			expectDetailIn: testErr.Error(),
		},
		{
			name:           "update operation",
			operation:      utils.OperationUpdate,
			expectSummary:  "unknown error",
			expectDetailIn: testErr.Error(),
		},
		{
			name:           "delete operation",
			operation:      utils.OperationDelete,
			expectSummary:  "unknown error",
			expectDetailIn: testErr.Error(),
		},
		{
			name:           "list operation",
			operation:      utils.OperationList,
			expectSummary:  "unknown error",
			expectDetailIn: testErr.Error(),
		},
		{
			name:           "import operation",
			operation:      utils.OperationImport,
			expectSummary:  "unknown error",
			expectDetailIn: testErr.Error(),
		},
		{
			name:           "open operation",
			operation:      utils.OperationOpen,
			expectSummary:  "unknown error",
			expectDetailIn: testErr.Error(),
		},
		{
			name:           "undefined operation",
			operation:      utils.OperationUndefined,
			expectSummary:  "unknown error",
			expectDetailIn: testErr.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diags := handler.GetDiagsFromError(ctx, testErr, tt.operation, nil)

			assert.True(t, diags.HasError())
			errors := diags.Errors()
			assert.Len(t, errors, 1)
			assert.Equal(t, tt.expectSummary, errors[0].Summary())
			assert.Contains(t, errors[0].Detail(), tt.expectDetailIn)
		})
	}
}

func TestUnit_GetDiagsFromError_FabricError(t *testing.T) {
	handler := utils.NewErrorHandler()
	ctx := t.Context()

	tests := []struct {
		name         string
		fabricErr    *fabcore.ResponseError
		wantContains []string
	}{
		{
			name:      "normal error response",
			fabricErr: createFabricError(t, 400, "InvalidParameter", "The parameter is invalid", "req-123"),
			wantContains: []string{
				"InvalidParameter",
				"The parameter is invalid",
				"req-123",
			},
		},
		{
			name: "nil RawResponse & nil ErrorResponse.RequestID",
			fabricErr: &fabcore.ResponseError{
				StatusCode: 400,
				ErrorResponse: &fabcore.ErrorResponse{
					ErrorCode: azto.Ptr("InvalidParameter"),
					Message:   azto.Ptr("The parameter is invalid"),
				},
			},
			wantContains: []string{
				"InvalidParameter",
				"The parameter is invalid",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diags := handler.GetDiagsFromError(ctx, tt.fabricErr, utils.OperationCreate, nil)

			assert.True(t, diags.HasError())
			diagErrs := diags.Errors()
			assert.Len(t, diagErrs, 1)

			for _, want := range tt.wantContains {
				assert.Contains(t, diagErrs[0].Detail(), want)
			}
		})
	}
}

func TestUnit_GetDiagsFromError_AuthFailedError(t *testing.T) {
	handler := utils.NewErrorHandler()
	ctx := t.Context()
	authErr := createAuthFailedError(t)

	diags := handler.GetDiagsFromError(ctx, authErr, utils.OperationCreate, nil)

	assert.True(t, diags.HasError())
	diagErrs := diags.Errors()
	assert.Len(t, diagErrs, 1)
	assert.Equal(t, "invalid_client", diagErrs[0].Summary())
	assert.Contains(t, diagErrs[0].Detail(), "Client authentication failed")
	assert.Contains(t, diagErrs[0].Detail(), "700016")
}

func TestUnit_GetDiagsFromError_AuthRequiredError(t *testing.T) {
	handler := utils.NewErrorHandler()
	ctx := t.Context()

	authErr := &azidentity.AuthenticationRequiredError{}

	diags := handler.GetDiagsFromError(ctx, authErr, utils.OperationCreate, nil)

	assert.True(t, diags.HasError())
	diagErrs := diags.Errors()
	assert.Len(t, diagErrs, 1)
	assert.Equal(t, "authentication required", diagErrs[0].Summary())
}

// Test backward compatibility functions.
func TestUnit_BackwardCompatibility(t *testing.T) {
	ctx := t.Context()
	testErr := errors.New("test error")

	// Test IsErr
	var diags diag.Diagnostics
	diags.AddError(testErr.Error(), testErr.Error())
	assert.True(t, utils.IsErr(diags, testErr))

	// Test IsErrNotFound
	diags = diag.Diagnostics{}
	diags.AddError(testErr.Error(), testErr.Error())
	assert.True(t, utils.IsErrNotFound("test-resource", &diags, testErr))

	// Test GetDiagsFromError
	result := utils.GetDiagsFromError(ctx, testErr, utils.OperationCreate, nil)
	assert.True(t, result.HasError())
}

func createFabricError(t *testing.T, statusCode int, errorCode, message, requestID string) *fabcore.ResponseError {
	t.Helper()

	resp := httptest.NewRecorder().Result()
	resp.StatusCode = statusCode

	errCode := errorCode
	errMsg := message
	errReqID := requestID

	return &fabcore.ResponseError{
		RawResponse: resp,
		StatusCode:  statusCode,
		ErrorResponse: &fabcore.ErrorResponse{
			ErrorCode: &errCode,
			Message:   &errMsg,
			RequestID: &errReqID,
		},
	}
}

func createAuthFailedError(t *testing.T) *azidentity.AuthenticationFailedError {
	t.Helper()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)

		_, err := w.Write([]byte(`{
			"error": "invalid_client",
			"error_description": "Client authentication failed",
			"error_codes": [700016],
			"timestamp": "2023-08-01T12:00:00Z",
			"trace_id": "trace-123",
			"correlation_id": "corr-456",
			"error_uri": "https://login.microsoftonline.com/error"
		}`))
		if err != nil {
			t.Fatalf("failed to write response: %v", err)
		}
	}))

	req, _ := http.NewRequestWithContext(t.Context(), http.MethodGet, server.URL, nil)
	resp, _ := http.DefaultClient.Do(req) //nolint:bodyclose
	err := &azidentity.AuthenticationFailedError{
		RawResponse: resp,
	}

	return err
}
