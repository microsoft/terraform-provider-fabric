// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package eventstreamsourceconnection_test

import (
	"context"
	"net/http"
	"regexp"
	"testing"

	azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	at "github.com/dcarbone/terraform-plugin-framework-utils/v3/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	"github.com/microsoft/fabric-sdk-go/fabric/eventstream"
	fabfake "github.com/microsoft/fabric-sdk-go/fabric/fake"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
)

var testDataSourceItemFQN, testDataSourceItemHeader = testhelp.TFDataSource(common.ProviderTypeName, itemTypeInfo.Type, "test")

func TestUnit_EventstreamDataSource(t *testing.T) {
	fakeWorkspaceID := testhelp.RandomUUID()
	fakeEventstreamID := testhelp.RandomUUID()
	fakeSourceID := testhelp.RandomUUID()

	eventstreamSourceConnection := eventstream.TopologyClientGetEventstreamSourceConnectionResponse{
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

	serverFactory := &fabfake.ServerFactory{}
	serverFactory.Eventstream.TopologyServer.GetEventstreamSourceConnection = func(ctx context.Context, workspaceID, eventstreamID, sourceID string, options *eventstream.TopologyClientGetEventstreamSourceConnectionOptions) (resp azfake.Responder[eventstream.TopologyClientGetEventstreamSourceConnectionResponse], errResp azfake.ErrorResponder) {
		if sourceID != fakeSourceID || workspaceID != fakeWorkspaceID || eventstreamID != fakeEventstreamID {
			resp.SetResponse(http.StatusNotFound, eventstream.TopologyClientGetEventstreamSourceConnectionResponse{}, nil)
			errResp.SetError(fabfake.SetResponseError(http.StatusNotFound, fabcore.ErrItem.ItemNotFound.Error(), "Item not found"))
		} else {
			resp.SetResponse(200, eventstreamSourceConnection, nil)
		}

		return
	}

	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, nil, serverFactory, nil, []resource.TestStep{
		// error - no attributes
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{},
			),
			ExpectError: regexp.MustCompile(`The argument "workspace_id" is required, but no definition was found`),
		},
		// error - workspace_id - invalid UUID
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"workspace_id":   "invalid uuid",
					"eventstream_id": fakeEventstreamID,
					"source_id":      fakeSourceID,
				},
			),
			ExpectError: regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
		// error - eventstream_id - invalid UUID
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"workspace_id":   fakeWorkspaceID,
					"eventstream_id": "invalid uuid",
					"source_id":      fakeSourceID,
				},
			),
			ExpectError: regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
		// error - source_id - invalid UUID
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"workspace_id":   fakeWorkspaceID,
					"eventstream_id": fakeEventstreamID,
					"source_id":      "invalid uuid",
				},
			),
			ExpectError: regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
		// error - unexpected attribute
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"workspace_id":    fakeWorkspaceID,
					"eventstream_id":  fakeEventstreamID,
					"source_id":       fakeSourceID,
					"unexpected_attr": "test",
				},
			),
			ExpectError: regexp.MustCompile(`An argument named "unexpected_attr" is not expected here`),
		},
		// error - no required attributes workspace_id
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"eventstream_id": fakeEventstreamID,
					"source_id":      fakeSourceID,
				},
			),
			ExpectError: regexp.MustCompile(`The argument "workspace_id" is required, but no definition was found`),
		},
		// error - no required attributes eventstream_id
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"workspace_id": fakeWorkspaceID,
					"source_id":    fakeSourceID,
				},
			),
			ExpectError: regexp.MustCompile(`The argument "eventstream_id" is required, but no definition was found`),
		},
		// error - no required attributes source_id
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"workspace_id":   fakeWorkspaceID,
					"eventstream_id": fakeEventstreamID,
				},
			),
			ExpectError: regexp.MustCompile(`The argument "source_id" is required, but no definition was found`),
		},
		// invalid workspace_id
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"workspace_id":   testhelp.RandomUUID(),
					"eventstream_id": fakeEventstreamID,
					"source_id":      fakeSourceID,
				},
			),
			ExpectError: regexp.MustCompile(common.ErrorReadHeader),
		},
		// invalid eventstream_id
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"workspace_id":   fakeWorkspaceID,
					"eventstream_id": testhelp.RandomUUID(),
					"source_id":      fakeSourceID,
				},
			),
			ExpectError: regexp.MustCompile(common.ErrorReadHeader),
		},
		// invalid source_id
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"workspace_id":   fakeWorkspaceID,
					"eventstream_id": fakeEventstreamID,
					"source_id":      testhelp.RandomUUID(),
				},
			),
			ExpectError: regexp.MustCompile(common.ErrorReadHeader),
		},
		// read
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"workspace_id":   fakeWorkspaceID,
					"eventstream_id": fakeEventstreamID,
					"source_id":      fakeSourceID,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrPtr(testDataSourceItemFQN, "workspace_id", &fakeWorkspaceID),
				resource.TestCheckResourceAttrPtr(testDataSourceItemFQN, "source_id", &fakeSourceID),
				resource.TestCheckResourceAttrPtr(testDataSourceItemFQN, "eventstream_id", &fakeEventstreamID),
				resource.TestCheckResourceAttrPtr(testDataSourceItemFQN, "event_hub_name", eventstreamSourceConnection.EventHubName),
				resource.TestCheckResourceAttrPtr(testDataSourceItemFQN, "fully_qualified_namespace", eventstreamSourceConnection.FullyQualifiedNamespace),
				resource.TestCheckResourceAttrPtr(testDataSourceItemFQN, "access_keys.primary_connection_string", eventstreamSourceConnection.AccessKeys.PrimaryConnectionString),
				resource.TestCheckResourceAttrPtr(testDataSourceItemFQN, "access_keys.secondary_connection_string", eventstreamSourceConnection.AccessKeys.SecondaryConnectionString),
				resource.TestCheckResourceAttrPtr(testDataSourceItemFQN, "access_keys.primary_key", eventstreamSourceConnection.AccessKeys.PrimaryKey),
				resource.TestCheckResourceAttrPtr(testDataSourceItemFQN, "access_keys.secondary_key", eventstreamSourceConnection.AccessKeys.SecondaryKey),
			),
		},
	}))
}

func TestAcc_EventstreamSourceConnectionDataSource(t *testing.T) {
	workspace := testhelp.WellKnown()["WorkspaceDS"].(map[string]any)
	workspaceID := workspace["id"].(string)

	evenstream := testhelp.WellKnown()["Eventstream"].(map[string]any)
	eventstreamID := evenstream["id"].(string)

	sourceConnection := evenstream["sourceConnection"].(map[string]any)
	sourceID := sourceConnection["sourceId"].(string)
	eventHubName := sourceConnection["eventHubName"].(string)
	FullyQualifiedNamespace := sourceConnection["fullyQualifiedNamespace"].(string)

	resource.ParallelTest(t, testhelp.NewTestAccCase(t, nil, nil, []resource.TestStep{
		// read by source id - not found
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"workspace_id":   workspaceID,
					"eventstream_id": eventstreamID,
					"source_id":      testhelp.RandomUUID(),
				},
			),
			ExpectError: regexp.MustCompile(common.ErrorReadHeader),
		},
		// read by id
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"workspace_id":   workspaceID,
					"eventstream_id": eventstreamID,
					"source_id":      sourceID,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "source_id", sourceID),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "eventstream_id", eventstreamID),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "workspace_id", workspaceID),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "event_hub_name", eventHubName),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "fully_qualified_namespace", FullyQualifiedNamespace),
				resource.TestCheckResourceAttrSet(testDataSourceItemFQN, "access_keys.primary_connection_string"),
				resource.TestCheckResourceAttrSet(testDataSourceItemFQN, "access_keys.secondary_connection_string"),
				resource.TestCheckResourceAttrSet(testDataSourceItemFQN, "access_keys.primary_key"),
				resource.TestCheckResourceAttrSet(testDataSourceItemFQN, "access_keys.secondary_key"),
			),
		},
	}))
}
