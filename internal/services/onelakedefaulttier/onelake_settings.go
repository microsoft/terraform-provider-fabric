// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package onelakedefaulttier

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/microsoft/fabric-sdk-go/fabric"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
)

type oneLakeSettingsResponse struct {
	Lifecycle *oneLakeSettingsLifecycle `json:"lifecycle"`
}

type oneLakeSettingsLifecycle struct {
	DefaultTier *fabcore.OneLakeAccessTier `json:"defaultTier"`
}

func getOneLakeSettings(ctx context.Context, client *fabric.Client, workspaceID string) (*oneLakeSettingsResponse, error) {
	const operationName = "core.OneLakeSettingsClient.GetSettings"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.Internal.Tracer(), nil)

	var err error
	defer func() {
		endSpan(err)
	}()

	// The service returns accessTimeTracking as a boolean, while SDK v0.23.0 models it as an enum.
	urlPath := "/v1/workspaces/" + url.PathEscape(workspaceID) + "/onelake/settings"

	request, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.Endpoint, urlPath))
	if err != nil {
		return nil, err
	}

	request.Raw().Header.Set("Accept", "application/json")

	response, err := client.Internal.Pipeline().Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if !runtime.HasStatusCode(response, http.StatusOK) {
		err = fabcore.NewResponseError(response)

		return nil, err
	}

	var settings oneLakeSettingsResponse

	err = json.NewDecoder(response.Body).Decode(&settings)
	if err != nil {
		return nil, fmt.Errorf("unmarshal OneLake settings response: %w", err)
	}

	return &settings, nil
}
