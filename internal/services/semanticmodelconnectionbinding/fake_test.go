// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package semanticmodelconnectionbinding_test

import (
	"context"
	"errors"
	"net/http"

	azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	fabfake "github.com/microsoft/fabric-sdk-go/fabric/fake"
	fabsemanticmodel "github.com/microsoft/fabric-sdk-go/fabric/semanticmodel"

	"github.com/microsoft/terraform-provider-fabric/internal/testhelp/fakes"
)

type semanticModelConnectionBindingKey struct {
	workspaceID     string
	semanticModelID string
	connectionPath  string
	connectionType  string
}

type semanticModelConnectionBindingFakeState struct {
	bindings map[semanticModelConnectionBindingKey]fabcore.ItemConnection
}

func newSemanticModelConnectionBindingFakeState() *semanticModelConnectionBindingFakeState {
	return &semanticModelConnectionBindingFakeState{
		bindings: make(map[semanticModelConnectionBindingKey]fabcore.ItemConnection),
	}
}

func (s *semanticModelConnectionBindingFakeState) configure() {
	fakes.FakeServer.ServerFactory.SemanticModel.ItemsServer.BindSemanticModelConnection = s.bind
	fakes.FakeServer.ServerFactory.Core.ItemsServer.NewListItemConnectionsPager = s.list
}

func (s *semanticModelConnectionBindingFakeState) add(workspaceID, semanticModelID string, binding fabcore.ItemConnection) {
	s.bindings[bindingKey(workspaceID, semanticModelID, *binding.ConnectionDetails.Path, *binding.ConnectionDetails.Type)] = binding
}

func (s *semanticModelConnectionBindingFakeState) unbind(workspaceID, semanticModelID, connectionPath, connectionType string) {
	key := bindingKey(workspaceID, semanticModelID, connectionPath, connectionType)
	binding := s.bindings[key]
	connectivityType := fabcore.ConnectivityTypeNone
	binding.ID = nil
	binding.ConnectivityType = &connectivityType
	s.bindings[key] = binding
}

func (s *semanticModelConnectionBindingFakeState) bind(
	_ context.Context,
	workspaceID string,
	semanticModelID string,
	request fabsemanticmodel.BindSemanticModelConnectionRequest,
	_ *fabsemanticmodel.ItemsClientBindSemanticModelConnectionOptions,
) (resp azfake.Responder[fabsemanticmodel.ItemsClientBindSemanticModelConnectionResponse], errResp azfake.ErrorResponder) {
	resp = azfake.Responder[fabsemanticmodel.ItemsClientBindSemanticModelConnectionResponse]{}

	if request.ConnectionBinding == nil || request.ConnectionBinding.ConnectionDetails == nil ||
		request.ConnectionBinding.ConnectionDetails.Path == nil || request.ConnectionBinding.ConnectionDetails.Type == nil ||
		request.ConnectionBinding.ConnectivityType == nil {
		errResp.SetError(fabfake.SetResponseError(http.StatusBadRequest, "InvalidRequest", "Connection details are required."))
		resp.SetResponse(http.StatusBadRequest, fabsemanticmodel.ItemsClientBindSemanticModelConnectionResponse{}, nil)

		return resp, errResp
	}

	connectionBinding := request.ConnectionBinding
	if connectionBinding.ID != nil &&
		(*connectionBinding.ConnectivityType == fabsemanticmodel.ConnectivityTypeAutomatic ||
			*connectionBinding.ConnectivityType == fabsemanticmodel.ConnectivityTypeNone) {
		errResp.SetError(fabfake.SetResponseError(http.StatusBadRequest, "InvalidRequest", "Connection ID must be omitted for Automatic and None connectivity."))
		resp.SetResponse(http.StatusBadRequest, fabsemanticmodel.ItemsClientBindSemanticModelConnectionResponse{}, nil)

		return resp, errResp
	}

	key := bindingKey(workspaceID, semanticModelID, *connectionBinding.ConnectionDetails.Path, *connectionBinding.ConnectionDetails.Type)
	if *connectionBinding.ConnectivityType == fabsemanticmodel.ConnectivityTypeNone {
		delete(s.bindings, key)
	} else {
		connectivityType := fabcore.ConnectivityType(*connectionBinding.ConnectivityType)
		s.bindings[key] = fabcore.ItemConnection{
			ConnectionDetails: &fabcore.ListConnectionDetails{
				Path: connectionBinding.ConnectionDetails.Path,
				Type: connectionBinding.ConnectionDetails.Type,
			},
			ConnectivityType: &connectivityType,
			DisplayName:      new("Test Connection"),
			ID:               connectionBinding.ID,
		}
	}

	resp.SetResponse(http.StatusOK, fabsemanticmodel.ItemsClientBindSemanticModelConnectionResponse{}, nil)

	return resp, errResp
}

func (s *semanticModelConnectionBindingFakeState) list(
	workspaceID string,
	semanticModelID string,
	_ *fabcore.ItemsClientListItemConnectionsOptions,
) (resp azfake.PagerResponder[fabcore.ItemsClientListItemConnectionsResponse]) {
	bindings := make([]fabcore.ItemConnection, 0, len(s.bindings))
	for key, binding := range s.bindings {
		if key.workspaceID == workspaceID && key.semanticModelID == semanticModelID {
			bindings = append(bindings, binding)
		}
	}

	resp = azfake.PagerResponder[fabcore.ItemsClientListItemConnectionsResponse]{}
	resp.AddPage(http.StatusOK, fabcore.ItemsClientListItemConnectionsResponse{
		Value: bindings,
	}, nil)

	return resp
}

func bindingKey(workspaceID, semanticModelID, connectionPath, connectionType string) semanticModelConnectionBindingKey {
	return semanticModelConnectionBindingKey{
		workspaceID:     workspaceID,
		semanticModelID: semanticModelID,
		connectionPath:  connectionPath,
		connectionType:  connectionType,
	}
}

func newRandomItemConnection(connectionID, connectionPath, connectivityType string) fabcore.ItemConnection {
	connectivity := fabcore.ConnectivityType(connectivityType)

	return fabcore.ItemConnection{
		ID: new(connectionID),
		ConnectionDetails: &fabcore.ListConnectionDetails{
			Path: new(connectionPath),
			Type: new("SQL"),
		},
		ConnectivityType: &connectivity,
		DisplayName:      new("Test Connection"),
	}
}

func newUnboundItemConnection(connectionPath string) fabcore.ItemConnection {
	connection := newRandomItemConnection("", connectionPath, string(fabsemanticmodel.ConnectivityTypeNone))
	connection.ID = nil

	return connection
}

func fakeListItemConnectionsError() func(
	workspaceID string,
	semanticModelID string,
	options *fabcore.ItemsClientListItemConnectionsOptions,
) (resp azfake.PagerResponder[fabcore.ItemsClientListItemConnectionsResponse]) {
	return func(
		_ string,
		_ string,
		_ *fabcore.ItemsClientListItemConnectionsOptions,
	) (resp azfake.PagerResponder[fabcore.ItemsClientListItemConnectionsResponse]) {
		resp.AddError(errors.New("list failed"))

		return resp
	}
}

func fakeBindSemanticModelConnectionError() func(
	ctx context.Context,
	workspaceID string,
	semanticModelID string,
	request fabsemanticmodel.BindSemanticModelConnectionRequest,
	options *fabsemanticmodel.ItemsClientBindSemanticModelConnectionOptions,
) (azfake.Responder[fabsemanticmodel.ItemsClientBindSemanticModelConnectionResponse], azfake.ErrorResponder) {
	return func(
		_ context.Context,
		_ string,
		_ string,
		_ fabsemanticmodel.BindSemanticModelConnectionRequest,
		_ *fabsemanticmodel.ItemsClientBindSemanticModelConnectionOptions,
	) (resp azfake.Responder[fabsemanticmodel.ItemsClientBindSemanticModelConnectionResponse], errResp azfake.ErrorResponder) {
		resp = azfake.Responder[fabsemanticmodel.ItemsClientBindSemanticModelConnectionResponse]{}

		errResp.SetError(fabfake.SetResponseError(http.StatusInternalServerError, "InternalError", "Bind failed."))
		resp.SetResponse(http.StatusInternalServerError, fabsemanticmodel.ItemsClientBindSemanticModelConnectionResponse{}, nil)

		return resp, errResp
	}
}
