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

type Operation string

const (
	OperationCreate    Operation = "create"
	OperationRead      Operation = "read"
	OperationUpdate    Operation = "update"
	OperationDelete    Operation = "delete"
	OperationList      Operation = "list"
	OperationImport    Operation = "import"
	OperationUndefined Operation = "undefined"
)

func IsErr(diags diag.Diagnostics, err error) bool {
	if !diags.HasError() {
		return false
	}

	if diags.Errors().Contains(diag.NewErrorDiagnostic(err.Error(), err.Error())) {
		return true
	}

	return false
}

func IsErrNotFound(resourceID string, diags *diag.Diagnostics, err error) bool {
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

func GetDiagsFromError(ctx context.Context, err error, operation Operation, errIs error) diag.Diagnostics { //nolint:gocognit, gocognit
	if err == nil {
		return nil
	}

	var diagErrSummary, diagErrDetail string

	switch operation {
	case OperationCreate:
		diagErrSummary = common.ErrorCreateHeader
		diagErrDetail = common.ErrorCreateDetails
	case OperationRead:
		diagErrSummary = common.ErrorReadHeader
		diagErrDetail = common.ErrorReadDetails
	case OperationUpdate:
		diagErrSummary = common.ErrorUpdateHeader
		diagErrDetail = common.ErrorUpdateDetails
	case OperationDelete:
		diagErrSummary = common.ErrorDeleteHeader
		diagErrDetail = common.ErrorDeleteDetails
	case OperationList:
		diagErrSummary = common.ErrorListHeader
		diagErrDetail = common.ErrorListDetails
	case OperationImport:
		diagErrSummary = common.ErrorImportHeader
		diagErrDetail = common.ErrorImportDetails
	default:
		diagErrSummary = ""
		diagErrDetail = ""
	}

	var errRespAzCore *azcore.ResponseError
	if errors.As(err, &errRespAzCore) {
		err = fabcore.NewResponseError(errRespAzCore.RawResponse)
	}

	var errRespFabric *fabcore.ResponseError
	var errAuthFailed *azidentity.AuthenticationFailedError
	var errAuthRequired *azidentity.AuthenticationRequiredError

	switch {
	case errors.As(err, &errRespFabric):
		tflog.Debug(ctx, "FABRIC ERROR", map[string]any{
			"StatusCode": errRespFabric.StatusCode,
			"ErrorCode":  errRespFabric.ErrorResponse.ErrorCode,
			"Message":    errRespFabric.ErrorResponse.Message,
			"RequestID":  errRespFabric.ErrorResponse.RequestID,
		})

		if (errIs != nil && errors.Is(err, errIs)) || (errIs != nil && errRespFabric.RawResponse.StatusCode == http.StatusNotFound) {
			diagErrSummary = errIs.Error()
			diagErrDetail = errIs.Error()

			break
		}

		var errCodes []string
		var errMessages []string

		if errRespFabric.ErrorResponse.ErrorCode != nil {
			errCodes = append(errCodes, *errRespFabric.ErrorResponse.ErrorCode)
		}

		if errRespFabric.ErrorResponse.Message != nil {
			errMessages = append(errMessages, *errRespFabric.ErrorResponse.Message)
		}

		errRequestID := ""
		if errRespFabric.ErrorResponse.RequestID != nil {
			errRequestID = *errRespFabric.ErrorResponse.RequestID
		}

		if len(errRespFabric.ErrorResponse.MoreDetails) > 0 {
			for _, errMoreDetail := range errRespFabric.ErrorResponse.MoreDetails {
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

		if diagErrSummary == "" {
			diagErrSummary = errCode
		}

		if diagErrDetail == "" {
			diagErrDetail = fmt.Sprintf("%s\n\nError Code: %s\nRequest ID: %s", errMessage, errCode, errRequestID)
		} else {
			diagErrDetail = fmt.Sprintf("%s: %s\n\nError Code: %s\nRequest ID: %s", diagErrDetail, errMessage, errCode, errRequestID)
		}
	case errors.As(err, &errAuthFailed):
		var errAuthResp authErrorResponse

		err := errAuthResp.getErrFromResp(errAuthFailed.RawResponse)
		if err != nil {
			diagErrSummary = "Failed to parse authentication error response"
			diagErrDetail = err.Error()
		} else {
			tflog.Debug(ctx, "AUTH FAILED ERROR", map[string]any{
				"CorrelationID":    errAuthResp.CorrelationID,
				"Error":            errAuthResp.Error,
				"ErrorDescription": errAuthResp.ErrorDescription,
				"ErrorURI":         errAuthResp.ErrorURI,
				"ErrorCodes":       errAuthResp.ErrorCodes,
				"Timestamp":        errAuthResp.Timestamp,
				"TraceID":          errAuthResp.TraceID,
			})

			diagErrSummary = errAuthResp.Error

			errCodes := make([]string, len(errAuthResp.ErrorCodes))
			for i, code := range errAuthResp.ErrorCodes {
				errCodes[i] = strconv.Itoa(code)
			}

			diagErrDetail = fmt.Sprintf("%s\n\nErrorCode: %s\nErrorURI: %s", errAuthResp.ErrorDescription, strings.Join(errCodes, " / "), errAuthResp.ErrorURI)
		}
	case errors.As(err, &errAuthRequired):
		tflog.Debug(ctx, "AUTH REQUIRED ERROR", map[string]any{
			"Error": err.Error(),
		})

		diagErrSummary = "authentication required"
		diagErrDetail = err.Error()
	default:
		tflog.Debug(ctx, "UNKNOWN ERROR", map[string]any{
			"Error": err.Error(),
		})

		diagErrSummary = "unknown error"
		diagErrDetail = err.Error()
	}

	var diags diag.Diagnostics

	diags.AddError(
		diagErrSummary,
		diagErrDetail,
	)

	tflog.Debug(ctx, err.Error())

	return diags
}

type authErrorResponse struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"` //nolint:tagliatelle
	ErrorCodes       []int  `json:"error_codes"`       //nolint:tagliatelle
	Timestamp        string `json:"timestamp"`
	TraceID          string `json:"trace_id"`       //nolint:tagliatelle
	CorrelationID    string `json:"correlation_id"` //nolint:tagliatelle
	ErrorURI         string `json:"error_uri"`      //nolint:tagliatelle
}

func (e *authErrorResponse) getErrFromResp(resp *http.Response) error {
	if resp.Body == nil {
		// this shouldn't happen in real-world scenarios as a
		// response with no body should set it to http.NoBody
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
