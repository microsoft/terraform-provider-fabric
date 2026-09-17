// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package semanticmodelcb_test

import (
	"context"
	"net/http"

	azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	fabsemanticmodel "github.com/microsoft/fabric-sdk-go/fabric/semanticmodel"
)

func fakeBindSemanticModelConnection() func(ctx context.Context, workspaceID, semanticModelID string, bindSemanticModelConnectionRequest fabsemanticmodel.BindSemanticModelConnectionRequest, options *fabsemanticmodel.ItemsClientBindSemanticModelConnectionOptions) (resp azfake.Responder[fabsemanticmodel.ItemsClientBindSemanticModelConnectionResponse], errResp azfake.ErrorResponder) {
	return func(_ context.Context, _, _ string, _ fabsemanticmodel.BindSemanticModelConnectionRequest, _ *fabsemanticmodel.ItemsClientBindSemanticModelConnectionOptions) (resp azfake.Responder[fabsemanticmodel.ItemsClientBindSemanticModelConnectionResponse], errResp azfake.ErrorResponder) {
		resp = azfake.Responder[fabsemanticmodel.ItemsClientBindSemanticModelConnectionResponse]{}
		resp.SetResponse(http.StatusOK, fabsemanticmodel.ItemsClientBindSemanticModelConnectionResponse{}, nil)

		return resp, errResp
	}
}
