// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package connectionra_test

import (
	"context"
	"net/http"

	azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	azto "github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
)

func fakeGetConnectionRoleAssignment(
	entityOrProvider any,
) func(ctx context.Context, connectionID, connectionRoleAssignmentID string, options *fabcore.ConnectionsClientGetConnectionRoleAssignmentOptions) (resp azfake.Responder[fabcore.ConnectionsClientGetConnectionRoleAssignmentResponse], errResp azfake.ErrorResponder) {
	var entityProvider func() fabcore.ConnectionRoleAssignment

	switch v := entityOrProvider.(type) {
	case fabcore.ConnectionRoleAssignment:
		entityProvider = func() fabcore.ConnectionRoleAssignment {
			return v
		}
	case func() fabcore.ConnectionRoleAssignment:
		entityProvider = v
	default:
		panic("fakeGetConnectionRoleAssignment: entityOrProvider must be either fabcore.ConnectionRoleAssignment or func() fabcore.ConnectionRoleAssignment")
	}

	return func(_ context.Context, _, _ string, _ *fabcore.ConnectionsClientGetConnectionRoleAssignmentOptions) (resp azfake.Responder[fabcore.ConnectionsClientGetConnectionRoleAssignmentResponse], errResp azfake.ErrorResponder) {
		resp = azfake.Responder[fabcore.ConnectionsClientGetConnectionRoleAssignmentResponse]{}
		resp.SetResponse(http.StatusOK, fabcore.ConnectionsClientGetConnectionRoleAssignmentResponse{ConnectionRoleAssignment: entityProvider()}, nil)

		return resp, errResp
	}
}

func fakeListConnectionRoleAssignments(
	exampleResp fabcore.ConnectionRoleAssignments,
) func(connectionID string, options *fabcore.ConnectionsClientListConnectionRoleAssignmentsOptions) (resp azfake.PagerResponder[fabcore.ConnectionsClientListConnectionRoleAssignmentsResponse]) {
	return func(_ string, _ *fabcore.ConnectionsClientListConnectionRoleAssignmentsOptions) (resp azfake.PagerResponder[fabcore.ConnectionsClientListConnectionRoleAssignmentsResponse]) {
		resp = azfake.PagerResponder[fabcore.ConnectionsClientListConnectionRoleAssignmentsResponse]{}
		resp.AddPage(http.StatusOK, fabcore.ConnectionsClientListConnectionRoleAssignmentsResponse{ConnectionRoleAssignments: exampleResp}, nil)

		return resp
	}
}

func fakeAddConnectionRoleAssignment(
	exampleResp fabcore.ConnectionRoleAssignment,
) func(ctx context.Context, connectionID string, addConnectionRoleAssignmentRequest fabcore.AddConnectionRoleAssignmentRequest, options *fabcore.ConnectionsClientAddConnectionRoleAssignmentOptions) (resp azfake.Responder[fabcore.ConnectionsClientAddConnectionRoleAssignmentResponse], errResp azfake.ErrorResponder) {
	return func(_ context.Context, _ string, _ fabcore.AddConnectionRoleAssignmentRequest, _ *fabcore.ConnectionsClientAddConnectionRoleAssignmentOptions) (resp azfake.Responder[fabcore.ConnectionsClientAddConnectionRoleAssignmentResponse], errResp azfake.ErrorResponder) {
		resp = azfake.Responder[fabcore.ConnectionsClientAddConnectionRoleAssignmentResponse]{}
		resp.SetResponse(http.StatusCreated, fabcore.ConnectionsClientAddConnectionRoleAssignmentResponse{ConnectionRoleAssignment: exampleResp}, nil)

		return resp, errResp
	}
}

func fakeUpdateConnectionRoleAssignment(
	handlerOrEntity any,
) func(ctx context.Context, connectionID, connectionRoleAssignmentID string, updateConnectionRoleAssignmentRequest fabcore.UpdateConnectionRoleAssignmentRequest, options *fabcore.ConnectionsClientUpdateConnectionRoleAssignmentOptions) (resp azfake.Responder[fabcore.ConnectionsClientUpdateConnectionRoleAssignmentResponse], errResp azfake.ErrorResponder) {
	var updateHandler func(req fabcore.UpdateConnectionRoleAssignmentRequest) fabcore.ConnectionRoleAssignment

	switch v := handlerOrEntity.(type) {
	case fabcore.ConnectionRoleAssignment:
		updateHandler = func(_ fabcore.UpdateConnectionRoleAssignmentRequest) fabcore.ConnectionRoleAssignment {
			return v
		}
	case func(req fabcore.UpdateConnectionRoleAssignmentRequest) fabcore.ConnectionRoleAssignment:
		updateHandler = v
	default:
		panic("fakeUpdateConnectionRoleAssignment: handlerOrEntity must be either fabcore.ConnectionRoleAssignment or func(req fabcore.UpdateConnectionRoleAssignmentRequest) fabcore.ConnectionRoleAssignment")
	}

	return func(_ context.Context, _, _ string, req fabcore.UpdateConnectionRoleAssignmentRequest, _ *fabcore.ConnectionsClientUpdateConnectionRoleAssignmentOptions) (resp azfake.Responder[fabcore.ConnectionsClientUpdateConnectionRoleAssignmentResponse], errResp azfake.ErrorResponder) {
		updatedEntity := updateHandler(req)
		resp = azfake.Responder[fabcore.ConnectionsClientUpdateConnectionRoleAssignmentResponse]{}
		resp.SetResponse(http.StatusOK, fabcore.ConnectionsClientUpdateConnectionRoleAssignmentResponse{ConnectionRoleAssignment: updatedEntity}, nil)

		return resp, errResp
	}
}

func fakeDeleteConnectionRoleAssignment() func(ctx context.Context, connectionID, connectionRoleAssignmentID string, options *fabcore.ConnectionsClientDeleteConnectionRoleAssignmentOptions) (resp azfake.Responder[fabcore.ConnectionsClientDeleteConnectionRoleAssignmentResponse], errResp azfake.ErrorResponder) {
	return func(_ context.Context, _, _ string, _ *fabcore.ConnectionsClientDeleteConnectionRoleAssignmentOptions) (resp azfake.Responder[fabcore.ConnectionsClientDeleteConnectionRoleAssignmentResponse], errResp azfake.ErrorResponder) {
		resp = azfake.Responder[fabcore.ConnectionsClientDeleteConnectionRoleAssignmentResponse]{}
		resp.SetResponse(http.StatusOK, fabcore.ConnectionsClientDeleteConnectionRoleAssignmentResponse{}, nil)

		return resp, errResp
	}
}

func NewRandomConnectionRoleAssignment() fabcore.ConnectionRoleAssignment {
	itemID := testhelp.RandomUUID()

	return fabcore.ConnectionRoleAssignment{
		ID: azto.Ptr(itemID),
		Principal: &fabcore.Principal{
			ID:   azto.Ptr(itemID),
			Type: azto.Ptr(testhelp.RandomElement(fabcore.PossiblePrincipalTypeValues())),
		},
		Role: azto.Ptr(testhelp.RandomElement(fabcore.PossibleConnectionRoleValues())),
	}
}

func fakeStatefulConnectionRoleAssignmentCRUD(
	initialEntity fabcore.ConnectionRoleAssignment,
	updatedEntity fabcore.ConnectionRoleAssignment,
) (
	getFn func(ctx context.Context, connectionID, connectionRoleAssignmentID string, options *fabcore.ConnectionsClientGetConnectionRoleAssignmentOptions) (resp azfake.Responder[fabcore.ConnectionsClientGetConnectionRoleAssignmentResponse], errResp azfake.ErrorResponder),
	updateFn func(ctx context.Context, connectionID, connectionRoleAssignmentID string, updateConnectionRoleAssignmentRequest fabcore.UpdateConnectionRoleAssignmentRequest, options *fabcore.ConnectionsClientUpdateConnectionRoleAssignmentOptions) (resp azfake.Responder[fabcore.ConnectionsClientUpdateConnectionRoleAssignmentResponse], errResp azfake.ErrorResponder),
) {
	currentEntity := initialEntity

	getFn = fakeGetConnectionRoleAssignment(func() fabcore.ConnectionRoleAssignment {
		return currentEntity
	})

	updateFn = fakeUpdateConnectionRoleAssignment(func(_ fabcore.UpdateConnectionRoleAssignmentRequest) fabcore.ConnectionRoleAssignment {
		currentEntity = updatedEntity

		return updatedEntity
	})

	return getFn, updateFn
}

func NewRandomConnectionRoleAssignments() fabcore.ConnectionRoleAssignments {
	principal0ID := testhelp.RandomUUID()
	principal1ID := testhelp.RandomUUID()
	principal2ID := testhelp.RandomUUID()

	return fabcore.ConnectionRoleAssignments{
		Value: []fabcore.ConnectionRoleAssignment{
			{
				ID:   azto.Ptr(principal0ID),
				Role: azto.Ptr(testhelp.RandomElement(fabcore.PossibleConnectionRoleValues())),
				Principal: &fabcore.Principal{
					ID:   azto.Ptr(principal0ID),
					Type: azto.Ptr(testhelp.RandomElement(fabcore.PossiblePrincipalTypeValues())),
				},
			},
			{
				ID:   azto.Ptr(principal1ID),
				Role: azto.Ptr(testhelp.RandomElement(fabcore.PossibleConnectionRoleValues())),
				Principal: &fabcore.Principal{
					ID:   azto.Ptr(principal1ID),
					Type: azto.Ptr(testhelp.RandomElement(fabcore.PossiblePrincipalTypeValues())),
				},
			},
			{
				ID:   azto.Ptr(principal2ID),
				Role: azto.Ptr(testhelp.RandomElement(fabcore.PossibleConnectionRoleValues())),
				Principal: &fabcore.Principal{
					ID:   azto.Ptr(principal2ID),
					Type: azto.Ptr(testhelp.RandomElement(fabcore.PossiblePrincipalTypeValues())),
				},
			},
		},
	}
}
