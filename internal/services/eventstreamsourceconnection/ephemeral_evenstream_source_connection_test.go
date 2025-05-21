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
	"github.com/microsoft/terraform-provider-fabric/internal/services/eventstreamsourceconnection"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
)

var (
	testEphemeralItemFQN, testEphemeralItemHeader         = testhelp.TFEphemeral(common.ProviderTypeName, eventstreamsourceconnection.ItemTypeInfo.Type, "test")
	testEphemeralItemEchoFQN, testEphemeralItemEchoConfig = testhelp.TFEphemeralEcho(testEphemeralItemFQN)
)

func TestAcc_EventstreamSourceConnectionEphemeralResource(t *testing.T) {
	workspace := testhelp.WellKnown()["WorkspaceDS"].(map[string]any)
	workspaceID := workspace["id"].(string)

	evenstream := testhelp.WellKnown()["Eventstream"].(map[string]any)
	eventstreamID := evenstream["id"].(string)
	sourceID := evenstream["sourceId"].(string)

	resource.ParallelTest(t, testhelp.NewTestAccCase(t, nil, nil, []resource.TestStep{
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
				statecheck.ExpectKnownValue(testEphemeralItemEchoFQN, tfjsonpath.New("data").AtMapKey("event_hub_name"), knownvalue.NotNull()),
				statecheck.ExpectKnownValue(testEphemeralItemEchoFQN, tfjsonpath.New("data").AtMapKey("access_keys"), knownvalue.NotNull()),
				statecheck.ExpectKnownValue(testEphemeralItemEchoFQN, tfjsonpath.New("data").AtMapKey("access_keys").AtMapKey("primary_key"), knownvalue.NotNull()),
				statecheck.ExpectKnownValue(testEphemeralItemEchoFQN, tfjsonpath.New("data").AtMapKey("access_keys").AtMapKey("secondary_key"), knownvalue.NotNull()),
				statecheck.ExpectKnownValue(testEphemeralItemEchoFQN, tfjsonpath.New("data").AtMapKey("access_keys").AtMapKey("primary_connection_string"), knownvalue.NotNull()),
				statecheck.ExpectKnownValue(testEphemeralItemEchoFQN, tfjsonpath.New("data").AtMapKey("access_keys").AtMapKey("secondary_connection_string"), knownvalue.NotNull()),
			},
		},
	}))
}
