// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package workspaceonelakeimmutability_test

import (
	"context"
	"net/http"

	azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	azto "github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	fabfake "github.com/microsoft/fabric-sdk-go/fabric/fake"
)

type oneLakeImmutabilityPolicyStore struct {
	policy            *fabcore.ImmutabilityPolicy
	modifyError       bool
	modifyErrorOnCall int
	modifyCalls       int
	getErrorStatus    int
	getErrorCode      string
}

func fakeGetOneLakeSettings(
	store *oneLakeImmutabilityPolicyStore,
) func(ctx context.Context, workspaceID string, options *fabcore.OneLakeSettingsClientGetSettingsOptions) (resp azfake.Responder[fabcore.OneLakeSettingsClientGetSettingsResponse], errResp azfake.ErrorResponder) {
	return func(_ context.Context, _ string, _ *fabcore.OneLakeSettingsClientGetSettingsOptions) (resp azfake.Responder[fabcore.OneLakeSettingsClientGetSettingsResponse], errResp azfake.ErrorResponder) {
		if store.getErrorStatus != 0 {
			errorCode := store.getErrorCode
			if errorCode == "" {
				errorCode = "TestFailure"
			}

			errResp.SetError(fabfake.SetResponseError(store.getErrorStatus, errorCode, "The fake settings request failed"))

			resp = azfake.Responder[fabcore.OneLakeSettingsClientGetSettingsResponse]{}
			resp.SetResponse(store.getErrorStatus, fabcore.OneLakeSettingsClientGetSettingsResponse{}, nil)

			return resp, errResp
		}

		policies := make([]fabcore.ImmutabilityPolicy, 0, 1)
		if store.policy != nil {
			policies = append(policies, *store.policy)
		}

		resp = azfake.Responder[fabcore.OneLakeSettingsClientGetSettingsResponse]{}
		resp.SetResponse(http.StatusOK, fabcore.OneLakeSettingsClientGetSettingsResponse{
			Diagnostics: &fabcore.OneLakeDiagnosticSettings{
				Status: new("Disabled"),
			},
			ImmutabilityPolicies: policies,
		}, nil)

		return resp, errResp
	}
}

func fakeModifyOneLakeImmutabilityPolicy(
	store *oneLakeImmutabilityPolicyStore,
) func(ctx context.Context, workspaceID string, request fabcore.ImmutabilityPolicyRequest, options *fabcore.OneLakeSettingsClientModifyImmutabilityPolicyOptions) (resp azfake.Responder[fabcore.OneLakeSettingsClientModifyImmutabilityPolicyResponse], errResp azfake.ErrorResponder) {
	return func(_ context.Context, _ string, request fabcore.ImmutabilityPolicyRequest, _ *fabcore.OneLakeSettingsClientModifyImmutabilityPolicyOptions) (resp azfake.Responder[fabcore.OneLakeSettingsClientModifyImmutabilityPolicyResponse], errResp azfake.ErrorResponder) {
		store.modifyCalls++
		resp = azfake.Responder[fabcore.OneLakeSettingsClientModifyImmutabilityPolicyResponse]{}

		if store.modifyError || (store.modifyErrorOnCall > 0 && store.modifyCalls >= store.modifyErrorOnCall) {
			errResp.SetError(fabfake.SetResponseError(http.StatusBadRequest, "TestFailure", "The fake modify failed"))
			resp.SetResponse(http.StatusBadRequest, fabcore.OneLakeSettingsClientModifyImmutabilityPolicyResponse{}, nil)

			return resp, errResp
		}

		if store.policy != nil && request.RetentionDays != nil && *request.RetentionDays < *store.policy.RetentionDays {
			errResp.SetError(fabfake.SetResponseError(http.StatusBadRequest, "RetentionDaysReductionNotAllowed", "Retention days cannot be reduced"))
			resp.SetResponse(http.StatusBadRequest, fabcore.OneLakeSettingsClientModifyImmutabilityPolicyResponse{}, nil)

			return resp, errResp
		}

		store.policy = &fabcore.ImmutabilityPolicy{
			Scope:         request.Scope,
			RetentionDays: request.RetentionDays,
		}

		resp.SetResponse(http.StatusOK, fabcore.OneLakeSettingsClientModifyImmutabilityPolicyResponse{}, nil)

		return resp, errResp
	}
}

type setPolicyBeforeApply struct {
	store         *oneLakeImmutabilityPolicyStore
	retentionDays int32
}

func (c setPolicyBeforeApply) CheckPlan(context.Context, plancheck.CheckPlanRequest, *plancheck.CheckPlanResponse) {
	policy := newOneLakeImmutabilityPolicy(c.retentionDays)
	c.store.policy = &policy
}

func newOneLakeImmutabilityPolicy(retentionDays int32) fabcore.ImmutabilityPolicy {
	return fabcore.ImmutabilityPolicy{
		Scope:         azto.Ptr(fabcore.ImmutabilityScopeDiagnosticLogs),
		RetentionDays: new(retentionDays),
	}
}
