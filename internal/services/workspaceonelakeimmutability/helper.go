// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package workspaceonelakeimmutability

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
)

func getImmutabilityPolicy(ctx context.Context, client *fabcore.OneLakeSettingsClient, workspaceID, scope string) (fabcore.ImmutabilityPolicy, bool, diag.Diagnostics) {
	response, err := client.GetSettings(ctx, workspaceID, nil)
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
