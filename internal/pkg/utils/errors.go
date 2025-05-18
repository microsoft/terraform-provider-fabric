// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package utils

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
)

// Operation represents the type of operation being performed by the provider.
type Operation string

// Supported operation types.
const (
	OperationCreate    Operation = "create"
	OperationRead      Operation = "read"
	OperationUpdate    Operation = "update"
	OperationDelete    Operation = "delete"
	OperationList      Operation = "list"
	OperationImport    Operation = "import"
	OperationUndefined Operation = "undefined"
)

// DefaultErrorHandler implements the ErrorHandler interface with standard error handling logic.
type ErrorHandler struct{}

// NewErrorHandler creates a new default error handler.
func NewErrorHandler() *ErrorHandler {
	return &ErrorHandler{}
}

// HasError checks if diagnostics contains a specific error.
func (h *ErrorHandler) HasError(diags diag.Diagnostics, err error) bool {
	if !diags.HasError() {
		return false
	}

	return diags.Errors().Contains(diag.NewErrorDiagnostic(err.Error(), err.Error()))
}

// IsNotFoundError checks if an error represents a resource not found condition and updates diagnostics accordingly.
func (h *ErrorHandler) IsNotFoundError(resourceID string, diags *diag.Diagnostics, err error) bool {
	if !diags.HasError() || err == nil {
		return false
	}

	if diags.Errors().Contains(diag.NewErrorDiagnostic(err.Error(), err.Error())) {
		var d diag.Diagnostics

		d.AddWarning(
			"Resource not found",
			fmt.Sprintf("Resource with ID %s not found. It may have been deleted outside of Terraform. Removing object from state.",
				resourceID,
			),
		)

		*diags = d

		return true
	}

	return false
}

// GetDiagsFromError converts an error to Terraform diagnostic messages.
func (h *ErrorHandler) GetDiagsFromError(ctx context.Context, err error, operation Operation, errIs error) diag.Diagnostics {
	if err == nil {
		return nil
	}

	diagErrSummary, diagErrDetail := h.getOperationErrorMessages(operation)

	// Check for known error types and extract appropriate details
	diagErrSummary, diagErrDetail = h.processError(ctx, err, errIs, diagErrSummary, diagErrDetail)

	var diags diag.Diagnostics
	diags.AddError(diagErrSummary, diagErrDetail)

	tflog.Debug(ctx, err.Error())

	return diags
}

// getOperationErrorMessages returns appropriate error messages based on operation type.
func (h *ErrorHandler) getOperationErrorMessages(operation Operation) (summary, detail string) { //nolint:nonamedreturns
	switch operation {
	case OperationCreate:
		return common.ErrorCreateHeader, common.ErrorCreateDetails
	case OperationRead:
		return common.ErrorReadHeader, common.ErrorReadDetails
	case OperationUpdate:
		return common.ErrorUpdateHeader, common.ErrorUpdateDetails
	case OperationDelete:
		return common.ErrorDeleteHeader, common.ErrorDeleteDetails
	case OperationList:
		return common.ErrorListHeader, common.ErrorListDetails
	case OperationImport:
		return common.ErrorImportHeader, common.ErrorImportDetails
	default:
		return "", ""
	}
}

// processError examines an error and returns appropriate diagnostic messages.
func (h *ErrorHandler) processError(ctx context.Context, err, errIs error, defaultSummary, defaultDetail string) (summary, detail string) { //nolint:nonamedreturns
	// Convert Azure Core error to Fabric error if needed
	var errRespAzCore *azcore.ResponseError
	if errors.As(err, &errRespAzCore) {
		err = fabcore.NewResponseError(errRespAzCore.RawResponse)
	}

	// Handle different error types
	var errRespFabric *fabcore.ResponseError
	var errAuthFailed *azidentity.AuthenticationFailedError
	var errAuthRequired *azidentity.AuthenticationRequiredError

	switch {
	case errors.As(err, &errRespFabric):
		return h.processFabricError(ctx, errRespFabric, errIs, defaultSummary, defaultDetail)
	case errors.As(err, &errAuthFailed):
		return h.processAuthFailedError(ctx, errAuthFailed)
	case errors.As(err, &errAuthRequired):
		return h.processAuthRequiredError(ctx, errAuthRequired)
	default:
		return h.processGenericError(ctx, err)
	}
}

// processFabricError handles Fabric-specific response errors.
func (h *ErrorHandler) processFabricError(
	ctx context.Context,
	errResp *fabcore.ResponseError,
	errIs error,
	defaultSummary, defaultDetail string,
) (string, string) { //revive:disable-line:confusing-results
	tflog.Debug(ctx, "FABRIC ERROR", map[string]any{
		"StatusCode": errResp.StatusCode,
		"ErrorCode":  errResp.ErrorResponse.ErrorCode,
		"Message":    errResp.ErrorResponse.Message,
		"RequestID":  errResp.ErrorResponse.RequestID,
	})

	// Handle special case for error identity check or not found errors
	if errIs != nil && (errors.Is(errResp, errIs) || (errResp.RawResponse != nil && errResp.RawResponse.StatusCode == http.StatusNotFound)) {
		return errIs.Error(), errIs.Error()
	}

	var errCodes []string
	var errMessages []string

	if errResp.ErrorResponse.ErrorCode != nil {
		errCodes = append(errCodes, *errResp.ErrorResponse.ErrorCode)
	}

	if errResp.ErrorResponse.Message != nil {
		errMessages = append(errMessages, *errResp.ErrorResponse.Message)
	}

	errRequestID := ""
	if errResp.ErrorResponse.RequestID != nil {
		errRequestID = *errResp.ErrorResponse.RequestID
	}

	// Collect additional error details
	if len(errResp.ErrorResponse.MoreDetails) > 0 {
		for _, errMoreDetail := range errResp.ErrorResponse.MoreDetails {
			if errMoreDetail.ErrorCode != nil {
				errCodes = append(errCodes, *errMoreDetail.ErrorCode)
			}

			if errMoreDetail.Message != nil {
				errMessages = append(errMessages, *errMoreDetail.Message)
			}
		}
	}

	errCode := strings.Join(errCodes, " / ")
	errMessage := strings.Join(errMessages, " / ")

	summary := defaultSummary
	if summary == "" {
		summary = errCode
	}

	detail := defaultDetail
	if detail == "" {
		detail = fmt.Sprintf("%s\n\nError Code: %s\nRequest ID: %s", errMessage, errCode, errRequestID)
	} else {
		detail = fmt.Sprintf("%s: %s\n\nError Code: %s\nRequest ID: %s", detail, errMessage, errCode, errRequestID)
	}

	return summary, detail
}

// processAuthFailedError handles Azure authentication failure errors.
func (h *ErrorHandler) processAuthFailedError(ctx context.Context, errAuthFailed *azidentity.AuthenticationFailedError) (string, string) { //revive:disable-line:confusing-results
	var errAuthResp authErrorResponse

	// Check if RawResponse is nil to avoid panic
	if errAuthFailed.RawResponse == nil {
		return "AuthenticationFailedError", errAuthFailed.Error()
	}

	err := errAuthResp.getErrFromResp(errAuthFailed.RawResponse)
	if err != nil {
		return "Failed to parse authentication error response", err.Error()
	}

	tflog.Debug(ctx, "AUTH FAILED ERROR", map[string]any{
		"CorrelationID":    errAuthResp.CorrelationID,
		"Error":            errAuthResp.Error,
		"ErrorDescription": errAuthResp.ErrorDescription,
		"ErrorURI":         errAuthResp.ErrorURI,
		"ErrorCodes":       errAuthResp.ErrorCodes,
		"Timestamp":        errAuthResp.Timestamp,
		"TraceID":          errAuthResp.TraceID,
	})

	summary := errAuthResp.Error

	errCodes := make([]string, len(errAuthResp.ErrorCodes))
	for i, code := range errAuthResp.ErrorCodes {
		errCodes[i] = strconv.Itoa(code)
	}

	detail := fmt.Sprintf("%s\n\nErrorCode: %s\nErrorURI: %s",
		errAuthResp.ErrorDescription,
		strings.Join(errCodes, " / "),
		errAuthResp.ErrorURI)

	return summary, detail
}

// processAuthRequiredError handles Azure authentication required errors.
func (h *ErrorHandler) processAuthRequiredError(ctx context.Context, errAuthRequired *azidentity.AuthenticationRequiredError) (string, string) { //revive:disable-line:confusing-results
	tflog.Debug(ctx, "AUTH REQUIRED ERROR", map[string]any{
		"Error": errAuthRequired.Error(),
	})

	return "authentication required", errAuthRequired.Error()
}

// processGenericError handles unknown error types.
func (h *ErrorHandler) processGenericError(ctx context.Context, err error) (string, string) { //revive:disable-line:confusing-results
	tflog.Debug(ctx, "UNKNOWN ERROR", map[string]any{
		"Error": err.Error(),
	})

	return "unknown error", err.Error()
}

// authErrorResponse represents the structure of an Azure authentication error response.
type authErrorResponse struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"` //nolint:tagliatelle
	ErrorCodes       []int  `json:"error_codes"`       //nolint:tagliatelle
	Timestamp        string `json:"timestamp"`
	TraceID          string `json:"trace_id"`       //nolint:tagliatelle
	CorrelationID    string `json:"correlation_id"` //nolint:tagliatelle
	ErrorURI         string `json:"error_uri"`      //nolint:tagliatelle
}

// getErrFromResp parses an HTTP response body into an authErrorResponse.
func (e *authErrorResponse) getErrFromResp(resp *http.Response) error {
	if resp == nil || resp.Body == nil {
		return nil
	}

	respBody, err := io.ReadAll(resp.Body)

	errClose := resp.Body.Close()
	if errClose != nil {
		return errClose
	}

	if err != nil {
		return err
	}

	if err := json.Unmarshal(respBody, &e); err != nil {
		return err
	}

	return nil
}

// For backward compatibility.
func IsErr(diags diag.Diagnostics, err error) bool {
	return NewErrorHandler().HasError(diags, err)
}

func IsErrNotFound(resourceID string, diags *diag.Diagnostics, err error) bool {
	return NewErrorHandler().IsNotFoundError(resourceID, diags, err)
}

func GetDiagsFromError(ctx context.Context, err error, operation Operation, errIs error) diag.Diagnostics {
	return NewErrorHandler().GetDiagsFromError(ctx, err, operation, errIs)
}
