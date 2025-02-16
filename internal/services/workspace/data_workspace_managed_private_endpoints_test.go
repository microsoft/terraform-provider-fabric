// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package workspace_test

import (
	"testing"

	at "github.com/dcarbone/terraform-plugin-framework-utils/v3/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
)

var (
	testWorkspaceManagedPrivateEndpointsFQN    = testhelp.DataSourceFQN("fabric", workspaceManagedPrivateEndpointsTFName, "test")
	testWorkspaceManagedPrivateEndpointsHeader = at.DataSourceHeader(testhelp.TypeName("fabric", workspaceManagedPrivateEndpointsTFName), "test")
)

func TestAcc_WorkspaceManagedPrivateEndpointsDataSource(t *testing.T) {
	workspace := testhelp.WellKnown()["WorkspaceDS"].(map[string]any)
	workspaceID := workspace["id"].(string)

	resource.ParallelTest(t, testhelp.NewTestAccCase(t, nil, nil, []resource.TestStep{
		// read
		{
			Config: at.CompileConfig(
				testWorkspaceManagedPrivateEndpointsHeader,
				map[string]any{
					"workspace_id": workspaceID,
				},
			),
			ConfigStateChecks: []statecheck.StateCheck{
				statecheck.ExpectKnownValue(
					testWorkspaceManagedPrivateEndpointsFQN,
					tfjsonpath.New("values").AtSliceIndex(0).AtMapKey("id"),
					knownvalue.NotNull(),
				),
			},
		},
	},
	))
}
