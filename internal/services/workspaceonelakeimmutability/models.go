// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package workspaceonelakeimmutability

import (
	"context"

	timeoutsD "github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts" //revive:disable-line:import-alias-naming
	timeoutsR "github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"   //revive:disable-line:import-alias-naming
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
)

type baseWorkspaceOneLakeImmutabilityPolicyModel struct {
	WorkspaceID   customtypes.UUID `tfsdk:"workspace_id"`
	Scope         types.String     `tfsdk:"scope"`
	RetentionDays types.Int64      `tfsdk:"retention_days"`
}

func (to *baseWorkspaceOneLakeImmutabilityPolicyModel) set(_ context.Context, workspaceID string, from fabcore.ImmutabilityPolicy) diag.Diagnostics {
	var diags diag.Diagnostics

	if from.Scope == nil {
		diags.AddError(
			"Invalid OneLake immutability policy response",
			"The API response did not include the policy scope.",
		)

		return diags
	}

	if from.RetentionDays == nil {
		diags.AddError(
			"Invalid OneLake immutability policy response",
			"The API response did not include the policy retention period.",
		)

		return diags
	}

	to.WorkspaceID = customtypes.NewUUIDValue(workspaceID)
	to.Scope = types.StringValue(string(*from.Scope))
	to.RetentionDays = types.Int64Value(int64(*from.RetentionDays))

	return nil
}

/*
DATA-SOURCE
*/

type dataSourceWorkspaceOneLakeImmutabilityPolicyModel struct {
	baseWorkspaceOneLakeImmutabilityPolicyModel

	Timeouts timeoutsD.Value `tfsdk:"timeouts"`
}

/*
RESOURCE
*/

type resourceWorkspaceOneLakeImmutabilityPolicyModel struct {
	baseWorkspaceOneLakeImmutabilityPolicyModel

	Timeouts timeoutsR.Value `tfsdk:"timeouts"`
}

type requestModifyWorkspaceOneLakeImmutabilityPolicy struct {
	fabcore.ImmutabilityPolicyRequest
}

func (to *requestModifyWorkspaceOneLakeImmutabilityPolicy) set(from resourceWorkspaceOneLakeImmutabilityPolicyModel) diag.Diagnostics {
	retentionDaysValue := from.RetentionDays.ValueInt64()
	if retentionDaysValue < 1 || retentionDaysValue > maxRetentionDays {
		var diags diag.Diagnostics
		diags.AddAttributeError(
			path.Root("retention_days"),
			"Invalid retention period",
			"The retention period must be between one day and the maximum value supported by the Fabric API.",
		)

		return diags
	}

	scope := fabcore.ImmutabilityScope(from.Scope.ValueString())
	retentionDays := int32(retentionDaysValue)

	to.ImmutabilityPolicyRequest = fabcore.ImmutabilityPolicyRequest{
		Scope:         &scope,
		RetentionDays: &retentionDays,
	}

	return nil
}

func findImmutabilityPolicy(policies []fabcore.ImmutabilityPolicy, scope string) (fabcore.ImmutabilityPolicy, bool) {
	for _, policy := range policies {
		if policy.Scope != nil && string(*policy.Scope) == scope {
			return policy, true
		}
	}

	return fabcore.ImmutabilityPolicy{}, false
}
