// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package spark_test

import (
	"testing"

	at "github.com/dcarbone/terraform-plugin-framework-utils/v3/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
)

var (
	testResourceSparkWorkspaceSettingsFQN    = testhelp.ResourceFQN("fabric", sparkWorkspaceSettingsTFName, "test")
	testResourceSparkWorkspaceSettingsHeader = at.ResourceHeader(testhelp.TypeName("fabric", sparkWorkspaceSettingsTFName), "test")
)

func TestAcc_SparkWorkspaceSettingsResource_CRUD(t *testing.T) {
	workspaceResourceHCL, workspaceResourceFQN := testhelp.TestAccWorkspaceResource(t, *testhelp.WellKnown().Capacity.ID)

	resource.Test(t, testhelp.NewTestAccCase(t, &testResourceSparkWorkspaceSettingsFQN, nil, []resource.TestStep{
		// Create and Read
		{
			ResourceName: testResourceSparkWorkspaceSettingsFQN,
			Config: at.JoinConfigs(
				workspaceResourceHCL,
				at.CompileConfig(
					testResourceSparkWorkspaceSettingsHeader,
					map[string]any{
						"workspace_id": testhelp.RefByFQN(workspaceResourceFQN, "id"),
						"automatic_log": map[string]any{
							"enabled": false,
						},
					},
				)),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceSparkWorkspaceSettingsFQN, "pool.default_pool.name", "Starter Pool"),
				resource.TestCheckResourceAttr(testResourceSparkWorkspaceSettingsFQN, "automatic_log.enabled", "false"),
			),
		},
		// Update and Read
		{
			ResourceName: testResourceSparkWorkspaceSettingsFQN,
			Config: at.JoinConfigs(
				workspaceResourceHCL,
				at.CompileConfig(
					testResourceSparkWorkspaceSettingsHeader,
					map[string]any{
						"workspace_id": testhelp.RefByFQN(workspaceResourceFQN, "id"),
						"automatic_log": map[string]any{
							"enabled": true,
						},
					},
				)),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceSparkWorkspaceSettingsFQN, "pool.default_pool.name", "Starter Pool"),
				resource.TestCheckResourceAttr(testResourceSparkWorkspaceSettingsFQN, "automatic_log.enabled", "true"),
			),
		},
	},
	))
}
