// // Copyright (c) Microsoft Corporation
// // SPDX-License-Identifier: MPL-2.0

package eventstreamsourceconnection_test

import (
	"regexp"
	"testing"

	at "github.com/dcarbone/terraform-plugin-framework-utils/v3/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp/fakes"
)

var (
	testEphemeralItemFQN, testEphemeralItemHeader         = testhelp.TFEphemeral(common.ProviderTypeName, itemTypeInfo.Type, "test")
	testEphemeralItemEchoFQN, testEphemeralItemEchoConfig = testhelp.TFEphemeralEcho(testEphemeralItemFQN)
)

func TestUnit_EventstreamSourceConnectionEphemeralResource(t *testing.T) {
	fakeWorkspaceID := testhelp.RandomUUID()
	fakeEventstreamID := testhelp.RandomUUID()
	fakeSourceID := testhelp.RandomUUID()

	entity := NewRandomEventstreamSourceConnection()
	fakes.FakeServer.ServerFactory.Eventstream.TopologyServer.GetEventstreamSourceConnection = fakeGetEventstreamSourceConnection(entity)

	resource.Test(t, testhelp.NewTestUnitCase(t, &testEphemeralItemFQN, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		// error - no attributes
		{
			Config: at.CompileConfig(
				testEphemeralItemHeader,
				map[string]any{},
			),
			ExpectError: regexp.MustCompile(`The argument "workspace_id" is required, but no definition was found`),
		},
		// error - workspace_id - invalid UUID
		{
			Config: at.CompileConfig(
				testEphemeralItemHeader,
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
				testEphemeralItemHeader,
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
				testEphemeralItemHeader,
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
				testEphemeralItemHeader,
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
				testEphemeralItemHeader,
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
				testEphemeralItemHeader,
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
				testEphemeralItemHeader,
				map[string]any{
					"workspace_id":   fakeWorkspaceID,
					"eventstream_id": fakeEventstreamID,
				},
			),
			ExpectError: regexp.MustCompile(`The argument "source_id" is required, but no definition was found`),
		},
		// read
		{
			Config: at.JoinConfigs(
				at.CompileConfig(
					testEphemeralItemHeader,
					map[string]any{
						"workspace_id":   fakeWorkspaceID,
						"eventstream_id": fakeEventstreamID,
						"source_id":      fakeSourceID,
					}),
				testEphemeralItemEchoConfig,
			),
			ConfigStateChecks: []statecheck.StateCheck{
				statecheck.ExpectKnownValue(testEphemeralItemEchoFQN, tfjsonpath.New("data").AtMapKey("workspace_id"), knownvalue.StringExact(fakeWorkspaceID)),
				statecheck.ExpectKnownValue(testEphemeralItemEchoFQN, tfjsonpath.New("data").AtMapKey("eventstream_id"), knownvalue.StringExact(fakeEventstreamID)),
				statecheck.ExpectKnownValue(testEphemeralItemEchoFQN, tfjsonpath.New("data").AtMapKey("source_id"), knownvalue.StringExact(fakeSourceID)),
				statecheck.ExpectKnownValue(testEphemeralItemEchoFQN, tfjsonpath.New("data").AtMapKey("event_hub_name"), knownvalue.StringExact(*entity.SourceConnectionResponse.EventHubName)),
				statecheck.ExpectKnownValue(
					testEphemeralItemEchoFQN,
					tfjsonpath.New("data").AtMapKey("fully_qualified_namespace"),
					knownvalue.StringExact(*entity.SourceConnectionResponse.FullyQualifiedNamespace),
				),
				statecheck.ExpectKnownValue(
					testEphemeralItemEchoFQN,
					tfjsonpath.New("data").AtMapKey("access_keys").AtMapKey("primary_key"),
					knownvalue.StringExact(*entity.SourceConnectionResponse.AccessKeys.PrimaryKey),
				),
				statecheck.ExpectKnownValue(
					testEphemeralItemEchoFQN,
					tfjsonpath.New("data").AtMapKey("access_keys").AtMapKey("secondary_key"),
					knownvalue.StringExact(*entity.SourceConnectionResponse.AccessKeys.SecondaryKey),
				),
				statecheck.ExpectKnownValue(
					testEphemeralItemEchoFQN,
					tfjsonpath.New("data").AtMapKey("access_keys").AtMapKey("primary_connection_string"),
					knownvalue.StringExact(*entity.SourceConnectionResponse.AccessKeys.PrimaryConnectionString),
				),
				statecheck.ExpectKnownValue(
					testEphemeralItemEchoFQN,
					tfjsonpath.New("data").AtMapKey("access_keys").AtMapKey("secondary_connection_string"),
					knownvalue.StringExact(*entity.SourceConnectionResponse.AccessKeys.SecondaryConnectionString),
				),
			},
		},
	}))
}

func TestAcc_EventstreamSourceConnectionEphemeralResource(t *testing.T) {
	workspace := testhelp.WellKnown()["WorkspaceDS"].(map[string]any)
	workspaceID := workspace["id"].(string)

	evenstream := testhelp.WellKnown()["Eventstream"].(map[string]any)
	eventstreamID := evenstream["id"].(string)

	sourceConnection := evenstream["sourceConnection"].(map[string]any)
	sourceID := sourceConnection["sourceId"].(string)
	eventHubName := sourceConnection["eventHubName"].(string)
	fullyQualifiedNamespace := sourceConnection["fullyQualifiedNamespace"].(string)

	resource.ParallelTest(t, testhelp.NewTestAccCase(t, &testEphemeralItemFQN, nil, []resource.TestStep{
		// Test error - source not found
		{
			Config: at.CompileConfig(
				testEphemeralItemHeader,
				map[string]any{
					"source_id":      testhelp.RandomUUID(),
					"eventstream_id": eventstreamID,
					"workspace_id":   workspaceID,
				},
			),
			ExpectError: regexp.MustCompile(common.ErrorOpenHeader),
		},
		// Test success - valid configuration with echo validation
		{
			Config: at.JoinConfigs(
				at.CompileConfig(
					testEphemeralItemHeader,
					map[string]any{
						"source_id":      sourceID,
						"eventstream_id": eventstreamID,
						"workspace_id":   workspaceID,
					}),
				testEphemeralItemEchoConfig,
			),
			ConfigStateChecks: []statecheck.StateCheck{
				statecheck.ExpectKnownValue(testEphemeralItemEchoFQN, tfjsonpath.New("data").AtMapKey("source_id"), knownvalue.StringExact(sourceID)),
				statecheck.ExpectKnownValue(testEphemeralItemEchoFQN, tfjsonpath.New("data").AtMapKey("eventstream_id"), knownvalue.StringExact(eventstreamID)),
				statecheck.ExpectKnownValue(testEphemeralItemEchoFQN, tfjsonpath.New("data").AtMapKey("workspace_id"), knownvalue.StringExact(workspaceID)),
				statecheck.ExpectKnownValue(testEphemeralItemEchoFQN, tfjsonpath.New("data").AtMapKey("event_hub_name"), knownvalue.StringExact(eventHubName)),
				statecheck.ExpectKnownValue(testEphemeralItemEchoFQN, tfjsonpath.New("data").AtMapKey("fully_qualified_namespace"), knownvalue.StringExact(fullyQualifiedNamespace)),
				statecheck.ExpectKnownValue(testEphemeralItemEchoFQN, tfjsonpath.New("data").AtMapKey("access_keys").AtMapKey("primary_key"), knownvalue.NotNull()),
				statecheck.ExpectKnownValue(testEphemeralItemEchoFQN, tfjsonpath.New("data").AtMapKey("access_keys").AtMapKey("secondary_key"), knownvalue.NotNull()),
				statecheck.ExpectKnownValue(testEphemeralItemEchoFQN, tfjsonpath.New("data").AtMapKey("access_keys").AtMapKey("primary_connection_string"), knownvalue.NotNull()),
				statecheck.ExpectKnownValue(testEphemeralItemEchoFQN, tfjsonpath.New("data").AtMapKey("access_keys").AtMapKey("secondary_connection_string"), knownvalue.NotNull()),
			},
		},
	}))
}
