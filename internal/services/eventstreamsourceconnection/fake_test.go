// // Copyright (c) Microsoft Corporation
// // SPDX-License-Identifier: MPL-2.0

package eventstreamsourceconnection_test

import (
	"context"
	"net/http"

	azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/microsoft/fabric-sdk-go/fabric/eventstream"

	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
)

func fakeGetEventstreamSourceConnection(
	entity eventstream.TopologyClientGetEventstreamSourceConnectionResponse,
) func(
	ctx context.Context,
	workspaceID, eventstreamID, sourceID string,
	options *eventstream.TopologyClientGetEventstreamSourceConnectionOptions,
) (resp azfake.Responder[eventstream.TopologyClientGetEventstreamSourceConnectionResponse], errResp azfake.ErrorResponder) {
	return func(_ context.Context, _, _, _ string, _ *eventstream.TopologyClientGetEventstreamSourceConnectionOptions) (resp azfake.Responder[eventstream.TopologyClientGetEventstreamSourceConnectionResponse], errResp azfake.ErrorResponder) {
		resp = azfake.Responder[eventstream.TopologyClientGetEventstreamSourceConnectionResponse]{}
		resp.SetResponse(http.StatusOK, entity, nil)

		return
	}
}

func NewRandomEventstreamSourceConnection() eventstream.TopologyClientGetEventstreamSourceConnectionResponse {
	return eventstream.TopologyClientGetEventstreamSourceConnectionResponse{
		SourceConnectionResponse: eventstream.SourceConnectionResponse{
			EventHubName:            to.Ptr(testhelp.RandomName()),
			FullyQualifiedNamespace: to.Ptr(testhelp.RandomName()),
			AccessKeys: &eventstream.AccessKeys{
				PrimaryConnectionString:   to.Ptr(testhelp.RandomName()),
				SecondaryConnectionString: to.Ptr(testhelp.RandomName()),
				PrimaryKey:                to.Ptr(testhelp.RandomName()),
				SecondaryKey:              to.Ptr(testhelp.RandomName()),
			},
		},
	}
}
