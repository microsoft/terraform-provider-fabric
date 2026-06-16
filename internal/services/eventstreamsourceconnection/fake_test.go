// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package eventstreamsourceconnection_test

import (
	"context"
	"net/http"

	azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	"github.com/microsoft/fabric-sdk-go/fabric/eventstream"
	fabfake "github.com/microsoft/fabric-sdk-go/fabric/fake"

	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
)

func fakeGetEventstreamSourceConnection(
	fakeWorkspaceID, fakeEventstreamID, fakeSourceID string,
	entity eventstream.TopologyClientGetEventstreamSourceConnectionResponse,
) func(
	ctx context.Context,
	workspaceID, eventstreamID, sourceID string,
	options *eventstream.TopologyClientGetEventstreamSourceConnectionOptions,
) (resp azfake.Responder[eventstream.TopologyClientGetEventstreamSourceConnectionResponse], errResp azfake.ErrorResponder) {
	return func(_ context.Context, workspaceID, eventstreamID, sourceID string, _ *eventstream.TopologyClientGetEventstreamSourceConnectionOptions) (resp azfake.Responder[eventstream.TopologyClientGetEventstreamSourceConnectionResponse], errResp azfake.ErrorResponder) {
		resp = azfake.Responder[eventstream.TopologyClientGetEventstreamSourceConnectionResponse]{}

		if sourceID != fakeSourceID || workspaceID != fakeWorkspaceID || eventstreamID != fakeEventstreamID {
			resp.SetResponse(http.StatusNotFound, eventstream.TopologyClientGetEventstreamSourceConnectionResponse{}, nil)
			errResp.SetError(fabfake.SetResponseError(http.StatusNotFound, fabcore.ErrItem.ItemNotFound.Error(), "Item not found"))
		} else {
			resp.SetResponse(http.StatusOK, entity, nil)
		}

		return resp, errResp
	}
}

func NewRandomEventstreamSourceConnection() eventstream.TopologyClientGetEventstreamSourceConnectionResponse {
	return eventstream.TopologyClientGetEventstreamSourceConnectionResponse{
		SourceConnectionResponse: eventstream.SourceConnectionResponse{
			EventHubName:            new(testhelp.RandomName()),
			FullyQualifiedNamespace: new(testhelp.RandomName()),
			AccessKeys: &eventstream.AccessKeys{
				PrimaryConnectionString:   new(testhelp.RandomName()),
				SecondaryConnectionString: new(testhelp.RandomName()),
				PrimaryKey:                new(testhelp.RandomName()),
				SecondaryKey:              new(testhelp.RandomName()),
			},
		},
	}
}
