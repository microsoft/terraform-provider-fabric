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

var testResourceItemFQN, testResourceItemHeader = testhelp.TFResource(common.ProviderTypeName, itemTypeInfo.Type, "test")

func TestAcc_SparkWorkspaceSettingsResource_CRUD(t *testing.T) {
	capacity := testhelp.WellKnown()["Capacity"].(map[string]any)
	capacityID := capacity["id"].(string)

	workspaceResourceHCL, workspaceResourceFQN := testhelp.TestAccWorkspaceResource(t, capacityID)

	resource.Test(t, testhelp.NewTestAccCase(t, &testResourceItemFQN, nil, []resource.TestStep{
		// Create and Read
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				workspaceResourceHCL,
				at.CompileConfig(
					testResourceItemHeader,
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
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "workspace_id"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "id"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "automatic_log.enabled", "false"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "high_concurrency.notebook_interactive_run_enabled", "false"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "high_concurrency.notebook_pipeline_run_enabled", "true"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "pool.customize_compute_enabled", "true"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "pool.default_pool.name", "Starter Pool"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "pool.default_pool.type", "Workspace"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "environment.runtime_version", "1.3"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "job.conservative_job_admission_enabled", "true"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "job.session_timeout_in_minutes", "60"),
			),
		},
		// Update and Read
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				workspaceResourceHCL,
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"workspace_id": testhelp.RefByFQN(workspaceResourceFQN, "id"),
						"automatic_log": map[string]any{
							"enabled": true,
						},
					},
				)),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "pool.default_pool.name", "Starter Pool"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "automatic_log.enabled", "true"),
			),
		},
		// Update and Read
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				workspaceResourceHCL,
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"workspace_id": testhelp.RefByFQN(workspaceResourceFQN, "id"),
						"automatic_log": map[string]any{
							"enabled": true,
						},
						"pool": map[string]any{
							"default_pool": map[string]any{
								"id": "00000000-0000-0000-0000-000000000000",
							},
						},
					},
				)),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "pool.default_pool.id", "00000000-0000-0000-0000-000000000000"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "pool.default_pool.name", "Starter Pool"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "automatic_log.enabled", "true"),
			),
		},
		// Update and Read
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				workspaceResourceHCL,
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"workspace_id": testhelp.RefByFQN(workspaceResourceFQN, "id"),
						"automatic_log": map[string]any{
							"enabled": true,
						},
						"pool": map[string]any{
							"default_pool": map[string]any{
								"name": "Starter Pool",
								"type": "Workspace",
							},
						},
					},
				)),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "pool.default_pool.id", "00000000-0000-0000-0000-000000000000"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "pool.default_pool.name", "Starter Pool"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "automatic_log.enabled", "true"),
			),
		},
	},
	))
}
