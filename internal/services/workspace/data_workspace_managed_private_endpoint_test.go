// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package workspace_test

import (
	"regexp"
	"testing"

	at "github.com/dcarbone/terraform-plugin-framework-utils/v3/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
)

var (
	testWorkspaceManagedPrivateEndpointFQN    = testhelp.DataSourceFQN("fabric", workspaceManagedPrivateEndpointsTFName, "test")
	testWorkspaceManagedPrivateEndpointHeader = at.DataSourceHeader(testhelp.TypeName("fabric", workspaceManagedPrivateEndpointsTFName), "test")
)

func TestAcc_WorkspaceManagedPrivateEndpointDataSource(t *testing.T) {
	workspace := testhelp.WellKnown()["WorkspaceMPE"].(map[string]any)
	workspaceID := workspace["id"].(string)

	entity := testhelp.WellKnown()["WorkspaceManagedPrivateEndpoint"].(map[string]any)
	entityID := entity["id"].(string)
	entityName := entity["name"].(string)

	resource.ParallelTest(t, testhelp.NewTestAccCase(t, nil, nil, []resource.TestStep{
		// read by id
		{
			Config: at.CompileConfig(
				testWorkspaceManagedPrivateEndpointHeader,
				map[string]any{
					"workspace_id": workspaceID,
					"id":           entityID,
				},
			),
			ConfigStateChecks: []statecheck.StateCheck{
				statecheck.ExpectKnownValue(
					testWorkspaceManagedPrivateEndpointFQN,
					tfjsonpath.New("id"),
					knownvalue.StringExact(entityID),
				),
				statecheck.ExpectKnownValue(
					testWorkspaceManagedPrivateEndpointFQN,
					tfjsonpath.New("name"),
					knownvalue.StringExact(entityName),
				),
			},
		},
		// read by id - not found
		{
			Config: at.CompileConfig(
				testWorkspaceManagedPrivateEndpointHeader,
				map[string]any{
					"workspace_id": workspaceID,
					"id":           testhelp.RandomUUID(),
				},
			),
			ExpectError: regexp.MustCompile(common.ErrorReadHeader),
		},
		// read by name
		{
			Config: at.CompileConfig(
				testWorkspaceManagedPrivateEndpointHeader,
				map[string]any{
					"workspace_id": workspaceID,
					"name":         entityName,
				},
			),
			ConfigStateChecks: []statecheck.StateCheck{
				statecheck.ExpectKnownValue(
					testWorkspaceManagedPrivateEndpointFQN,
					tfjsonpath.New("id"),
					knownvalue.StringExact(entityID),
				),
				statecheck.ExpectKnownValue(
					testWorkspaceManagedPrivateEndpointFQN,
					tfjsonpath.New("name"),
					knownvalue.StringExact(entityName),
				),
			},
		},
		// read by name - not found
		{
			Config: at.CompileConfig(
				testWorkspaceManagedPrivateEndpointHeader,
				map[string]any{
					"workspace_id": workspaceID,
					"name":         testhelp.RandomName(),
				},
			),
			ExpectError: regexp.MustCompile(common.ErrorReadHeader),
		},
	}))
}
