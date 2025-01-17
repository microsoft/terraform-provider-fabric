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
	capacity := testhelp.WellKnown()["Capacity"].(map[string]any)
	capacityID := capacity["id"].(string)

	workspaceResourceHCL, workspaceResourceFQN := testhelp.TestAccWorkspaceResource(t, capacityID)

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
						"high_concurrency": map[string]any{
							"notebook_interactive_run_enabled": false,
							"notebook_pipeline_run_enabled":    true,
						},
						"job": map[string]any{
							"conservative_job_admission_enabled": true,
							"session_timeout_in_minutes":         60,
						},
					},
				)),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrSet(testDataSourceSparkWorkspaceSettingsFQN, "workspace_id"),
				resource.TestCheckResourceAttrSet(testDataSourceSparkWorkspaceSettingsFQN, "id"),
				resource.TestCheckResourceAttr(testDataSourceSparkWorkspaceSettingsFQN, "automatic_log.enabled", "false"),
				resource.TestCheckResourceAttr(testDataSourceSparkWorkspaceSettingsFQN, "high_concurrency.notebook_interactive_run_enabled", "false"),
				resource.TestCheckResourceAttr(testDataSourceSparkWorkspaceSettingsFQN, "high_concurrency.notebook_pipeline_run_enabled", "true"),
				resource.TestCheckResourceAttr(testDataSourceSparkWorkspaceSettingsFQN, "pool.customize_compute_enabled", "true"),
				resource.TestCheckResourceAttr(testDataSourceSparkWorkspaceSettingsFQN, "pool.default_pool.name", "Starter Pool"),
				resource.TestCheckResourceAttr(testDataSourceSparkWorkspaceSettingsFQN, "pool.default_pool.type", "Workspace"),
				resource.TestCheckResourceAttr(testDataSourceSparkWorkspaceSettingsFQN, "environment.runtime_version", "1.3"),
				resource.TestCheckResourceAttr(testDataSourceSparkWorkspaceSettingsFQN, "job.conservative_job_admission_enabled", "true"),
				resource.TestCheckResourceAttr(testDataSourceSparkWorkspaceSettingsFQN, "job.session_timeout_in_minutes", "60"),
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
