// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package onelakediagnosticsetting

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/microsoft/fabric-sdk-go/fabric"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
)

func getOneLakeDiagnosticSettings(
	ctx context.Context,
	client *fabric.Client,
	workspaceID string,
) (*fabcore.OneLakeDiagnosticSettings, error) {
	if client == nil || client.Internal == nil {
		return nil, errors.New("fabric client is not configured")
	}

	const operationName = "core.OneLakeSettingsClient.GetSettings"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.Internal.Tracer(), nil)

	var spanError error
	defer func() { endSpan(spanError) }()

	urlPath := "/v1/workspaces/{workspaceId}/onelake/settings"
	urlPath = strings.ReplaceAll(urlPath, "{workspaceId}", url.PathEscape(workspaceID))

	request, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.Endpoint, urlPath))
	if err != nil {
		spanError = err

		return nil, fmt.Errorf("creating OneLake settings request: %w", err)
	}

	request.Raw().Header.Set("Accept", "application/json")

	response, err := client.Internal.Pipeline().Do(request)
	if err != nil {
		spanError = err

		return nil, fmt.Errorf("sending OneLake settings request: %w", err)
	}

	if response == nil {
		spanError = errors.New("OneLake settings request returned a nil HTTP response")

		return nil, spanError
	}

	if response.StatusCode != http.StatusOK {
		spanError = fabcore.NewResponseError(response)
		if spanError != nil {
			return nil, spanError
		}

		spanError = fmt.Errorf("OneLake settings request returned HTTP status %d without an error response body", response.StatusCode)

		return nil, spanError
	}

	if response.Body == nil {
		spanError = errors.New("OneLake settings response has no body")

		return nil, spanError
	}
	defer response.Body.Close()

	// Decode only the setting this resource owns; unrelated OneLake fields can
	// have wire types that the pinned SDK does not accept.
	var result struct {
		Diagnostics *fabcore.OneLakeDiagnosticSettings `json:"diagnostics"`
	}

	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		spanError = err

		return nil, fmt.Errorf("unmarshalling OneLake diagnostics response: %w", err)
	}

	return result.Diagnostics, nil
}
