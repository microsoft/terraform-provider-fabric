// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package eventstreamsourceconnection_test

import (
	"regexp"
	"testing"

	at "github.com/dcarbone/terraform-plugin-framework-utils/v3/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp/fakes"
)

var testDataSourceItemFQN, testDataSourceItemHeader = testhelp.TFDataSource(common.ProviderTypeName, itemTypeInfo.Type, "test")

func TestUnit_EventstreamDataSource(t *testing.T) {
	fakeWorkspaceID := testhelp.RandomUUID()
	fakeEventstreamID := testhelp.RandomUUID()
	fakeSourceID := testhelp.RandomUUID()

	entity := NewRandomEventstreamSourceConnection()
	fakes.FakeServer.ServerFactory.Eventstream.TopologyServer.GetEventstreamSourceConnection = fakeGetEventstreamSourceConnection(
		fakeWorkspaceID,
		fakeEventstreamID,
		fakeSourceID,
		entity)

	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, nil, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
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
				resource.TestCheckResourceAttrPtr(testDataSourceItemFQN, "event_hub_name", entity.EventHubName),
				resource.TestCheckResourceAttrPtr(testDataSourceItemFQN, "fully_qualified_namespace", entity.FullyQualifiedNamespace),
				resource.TestCheckResourceAttrPtr(testDataSourceItemFQN, "access_keys.primary_connection_string", entity.AccessKeys.PrimaryConnectionString),
				resource.TestCheckResourceAttrPtr(testDataSourceItemFQN, "access_keys.secondary_connection_string", entity.AccessKeys.SecondaryConnectionString),
				resource.TestCheckResourceAttrPtr(testDataSourceItemFQN, "access_keys.primary_key", entity.AccessKeys.PrimaryKey),
				resource.TestCheckResourceAttrPtr(testDataSourceItemFQN, "access_keys.secondary_key", entity.AccessKeys.SecondaryKey),
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
	fullyQualifiedNamespace := sourceConnection["fullyQualifiedNamespace"].(string)

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
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "fully_qualified_namespace", fullyQualifiedNamespace),
				resource.TestCheckResourceAttrSet(testDataSourceItemFQN, "access_keys.primary_connection_string"),
				resource.TestCheckResourceAttrSet(testDataSourceItemFQN, "access_keys.secondary_connection_string"),
				resource.TestCheckResourceAttrSet(testDataSourceItemFQN, "access_keys.primary_key"),
				resource.TestCheckResourceAttrSet(testDataSourceItemFQN, "access_keys.secondary_key"),
			),
		},
	}))
}
