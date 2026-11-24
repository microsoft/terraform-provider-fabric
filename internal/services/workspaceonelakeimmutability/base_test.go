// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package workspaceonelakeimmutability_test

import (
	"math"
	"testing"

	"github.com/microsoft/terraform-provider-fabric/internal/services/workspaceonelakeimmutability"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
)

var itemTypeInfo = workspaceonelakeimmutability.ItemTypeInfo

func oneLakeImmutabilityPolicyFromWellKnown(t *testing.T) (workspaceID, scope string, retentionDays int64) {
	t.Helper()

	if !testhelp.IsWellKnownDataAvailable() {
		t.Skip("well-known data is unavailable")
	}

	policy, ok := testhelp.WellKnown()["OneLakeImmutabilityPolicy"].(map[string]any)
	if !ok {
		t.Skip("OneLake immutability policy is not ready; rerun Set-WellKnown after diagnostic events have propagated")
	}

	workspaceID, ok = policy["workspaceId"].(string)
	if !ok || workspaceID == "" {
		t.Fatal("well-known OneLake immutability policy has no valid workspaceId")
	}

	scope, ok = policy["scope"].(string)
	if !ok || scope == "" {
		t.Fatal("well-known OneLake immutability policy has no valid scope")
	}

	retentionDaysValue, ok := policy["retentionDays"].(float64)
	if !ok || retentionDaysValue < 1 || retentionDaysValue > math.MaxInt32 || math.Trunc(retentionDaysValue) != retentionDaysValue {
		t.Fatal("well-known OneLake immutability policy has no valid retentionDays")
	}

	return workspaceID, scope, int64(retentionDaysValue)
}
