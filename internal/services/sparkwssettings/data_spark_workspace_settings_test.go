// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package sparkwssettings_test

import (
	"testing"

	at "github.com/dcarbone/terraform-plugin-framework-utils/v3/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
)

var testDataSourceItemFQN, testDataSourceItemHeader = testhelp.TFDataSource(common.ProviderTypeName, itemTypeInfo.Type, "test")

func TestAcc_SparkWorkspaceSettingsDataSource(t *testing.T) {
	workspace := testhelp.WellKnown()["WorkspaceDS"].(map[string]any)
	workspaceID := workspace["id"].(string)

	resource.ParallelTest(t, testhelp.NewTestAccCase(t, &testDataSourceItemFQN, nil, []resource.TestStep{
		// read
		{
			ResourceName: testDataSourceItemFQN,
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"workspace_id": workspaceID,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "workspace_id", workspaceID),
				resource.TestCheckResourceAttrSet(testDataSourceItemFQN, "id"),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "automatic_log.enabled", "true"),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "high_concurrency.notebook_interactive_run_enabled", "true"),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "high_concurrency.notebook_pipeline_run_enabled", "false"),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "pool.customize_compute_enabled", "true"),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "pool.default_pool.name", "Starter Pool"),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "pool.default_pool.type", "Workspace"),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "environment.runtime_version", "1.3"),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "job.conservative_job_admission_enabled", "false"),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "job.session_timeout_in_minutes", "20"),
			),
		},
	},
	))
}
