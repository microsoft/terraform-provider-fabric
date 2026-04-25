// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package eventstreamdestinationconnection_test

import (
	"context"
	"net/http"

	azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	"github.com/microsoft/fabric-sdk-go/fabric/eventstream"
	fabfake "github.com/microsoft/fabric-sdk-go/fabric/fake"

	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
)

func fakeGetEventstreamDestinationConnection(
	fakeWorkspaceID, fakeEventstreamID, fakeDestinationID string,
	entity eventstream.TopologyClientGetEventstreamDestinationConnectionResponse,
) func(
	ctx context.Context,
	workspaceID, eventstreamID, destinationID string,
	options *eventstream.TopologyClientGetEventstreamDestinationConnectionOptions,
) (resp azfake.Responder[eventstream.TopologyClientGetEventstreamDestinationConnectionResponse], errResp azfake.ErrorResponder) {
	return func(_ context.Context, workspaceID, eventstreamID, destinationID string, _ *eventstream.TopologyClientGetEventstreamDestinationConnectionOptions) (resp azfake.Responder[eventstream.TopologyClientGetEventstreamDestinationConnectionResponse], errResp azfake.ErrorResponder) {
		resp = azfake.Responder[eventstream.TopologyClientGetEventstreamDestinationConnectionResponse]{}

		if destinationID != fakeDestinationID || workspaceID != fakeWorkspaceID || eventstreamID != fakeEventstreamID {
			resp.SetResponse(http.StatusNotFound, eventstream.TopologyClientGetEventstreamDestinationConnectionResponse{}, nil)
			errResp.SetError(fabfake.SetResponseError(http.StatusNotFound, fabcore.ErrItem.ItemNotFound.Error(), "Item not found"))
		} else {
			resp.SetResponse(http.StatusOK, entity, nil)
		}

		return resp, errResp
	}
}

func NewRandomEventstreamDestinationConnection() eventstream.TopologyClientGetEventstreamDestinationConnectionResponse {
	return eventstream.TopologyClientGetEventstreamDestinationConnectionResponse{
		DestinationConnectionResponse: eventstream.DestinationConnectionResponse{
			EventHubName:            new(testhelp.RandomName()),
			FullyQualifiedNamespace: new(testhelp.RandomName()),
			ConsumerGroupName:       new(testhelp.RandomName()),
			AccessKeys: &eventstream.AccessKeys{
				PrimaryConnectionString:   new(testhelp.RandomName()),
				SecondaryConnectionString: new(testhelp.RandomName()),
				PrimaryKey:                new(testhelp.RandomName()),
				SecondaryKey:              new(testhelp.RandomName()),
			},
		},
	}
}
