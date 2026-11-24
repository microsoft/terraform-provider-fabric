// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package workspaceonelakeimmutability

import (
	"context"
	"net/http"
	"net/url"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/microsoft/fabric-sdk-go/fabric"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
)

func getImmutabilityPolicy(ctx context.Context, client *fabric.Client, workspaceID, scope string) (fabcore.ImmutabilityPolicy, bool, diag.Diagnostics) {
	const operationName = "core.OneLakeSettingsClient.GetSettings"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)

	urlPath := "/v1/workspaces/" + url.PathEscape(workspaceID) + "/onelake/settings"

	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.Endpoint, urlPath))
	if err != nil {
		return fabcore.ImmutabilityPolicy{}, false, utils.GetDiagsFromError(ctx, err, utils.OperationRead, nil)
	}

	req.Raw().Header["Accept"] = []string{"application/json"}

	httpResp, err := client.Internal.Pipeline().Do(req)
	if httpResp != nil && httpResp.Body != nil {
		defer httpResp.Body.Close()
	}

	if err != nil {
		return fabcore.ImmutabilityPolicy{}, false, utils.GetDiagsFromError(ctx, err, utils.OperationRead, fabcore.ErrCommon.EntityNotFound)
	}

	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = fabcore.NewResponseError(httpResp)

		return fabcore.ImmutabilityPolicy{}, false, utils.GetDiagsFromError(ctx, err, utils.OperationRead, fabcore.ErrCommon.EntityNotFound)
	}

	var response struct {
		ImmutabilityPolicies []fabcore.ImmutabilityPolicy `json:"immutabilityPolicies"`
	}

	err = runtime.UnmarshalAsJSON(httpResp, &response)
	if diags := utils.GetDiagsFromError(ctx, err, utils.OperationRead, fabcore.ErrCommon.EntityNotFound); diags.HasError() {
		return fabcore.ImmutabilityPolicy{}, false, diags
	}

	policy, found := findImmutabilityPolicy(response.ImmutabilityPolicies, scope)
	if !found {
		return fabcore.ImmutabilityPolicy{}, false, nil
	}

	if policy.RetentionDays == nil {
		var diags diag.Diagnostics
		diags.AddError(
			"Invalid OneLake immutability policy response",
			"The API response did not include the policy retention period.",
		)

		return fabcore.ImmutabilityPolicy{}, false, diags
	}

	if *policy.RetentionDays < 1 {
		var diags diag.Diagnostics
		diags.AddError(
			"Invalid OneLake immutability policy response",
			"The API response included a retention period less than one day.",
		)

		return fabcore.ImmutabilityPolicy{}, false, diags
	}

	return policy, true, nil
}
