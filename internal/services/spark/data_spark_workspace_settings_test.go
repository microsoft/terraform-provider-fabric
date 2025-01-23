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
	testDataSourceSparkWorkspaceSettingsFQN    = testhelp.DataSourceFQN("fabric", sparkWorkspaceSettingsTFName, "test")
	testDataSourceSparkWorkspaceSettingsHeader = at.DataSourceHeader(testhelp.TypeName("fabric", sparkWorkspaceSettingsTFName), "test")
)

func TestAcc_SparkWorkspaceSettingsDataSource(t *testing.T) {
	workspace := testhelp.WellKnown()["WorkspaceDS"].(map[string]any)
	workspaceID := workspace["id"].(string)

	resource.ParallelTest(t, testhelp.NewTestAccCase(t, &testDataSourceSparkWorkspaceSettingsFQN, nil, []resource.TestStep{
		// read
		{
			ResourceName: testDataSourceSparkWorkspaceSettingsFQN,
			Config: at.CompileConfig(
				testDataSourceSparkWorkspaceSettingsHeader,
				map[string]any{
					"workspace_id": workspaceID,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testDataSourceSparkWorkspaceSettingsFQN, "workspace_id", workspaceID),
				resource.TestCheckResourceAttrSet(testDataSourceSparkWorkspaceSettingsFQN, "id"),
				resource.TestCheckResourceAttr(testDataSourceSparkWorkspaceSettingsFQN, "automatic_log.enabled", "true"),
				resource.TestCheckResourceAttr(testDataSourceSparkWorkspaceSettingsFQN, "high_concurrency.notebook_interactive_run_enabled", "true"),
				resource.TestCheckResourceAttr(testDataSourceSparkWorkspaceSettingsFQN, "high_concurrency.notebook_pipeline_run_enabled", "false"),
				resource.TestCheckResourceAttr(testDataSourceSparkWorkspaceSettingsFQN, "pool.customize_compute_enabled", "true"),
				resource.TestCheckResourceAttr(testDataSourceSparkWorkspaceSettingsFQN, "pool.default_pool.name", "Starter Pool"),
				resource.TestCheckResourceAttr(testDataSourceSparkWorkspaceSettingsFQN, "pool.default_pool.type", "Workspace"),
				resource.TestCheckResourceAttr(testDataSourceSparkWorkspaceSettingsFQN, "environment.runtime_version", "1.3"),
				resource.TestCheckResourceAttr(testDataSourceSparkWorkspaceSettingsFQN, "job.conservative_job_admission_enabled", "false"),
				resource.TestCheckResourceAttr(testDataSourceSparkWorkspaceSettingsFQN, "job.session_timeout_in_minutes", "20"),
			),
		},
	},
	))
}
