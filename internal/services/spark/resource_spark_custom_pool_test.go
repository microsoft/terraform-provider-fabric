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
	testResourceSparkCustomPoolFQN    = testhelp.ResourceFQN("fabric", sparkCustomPoolTFName, "test")
	testResourceSparkCustomPoolHeader = at.ResourceHeader(testhelp.TypeName("fabric", sparkCustomPoolTFName), "test")
)

func TestAcc_SparkCustomPoolResource_CRUD(t *testing.T) {
	capacity := testhelp.WellKnown()["Capacity"].(map[string]any)
	capacityID := capacity["id"].(string)

	workspaceResourceHCL, workspaceResourceFQN := testhelp.TestAccWorkspaceResource(t, capacityID)
	entityCreateName := testhelp.RandomName()
	entityUpdateName := testhelp.RandomName()

	resource.Test(t, testhelp.NewTestAccCase(t, &testResourceSparkCustomPoolFQN, nil, []resource.TestStep{
		// Create and Read
		{
			ResourceName: testResourceSparkCustomPoolFQN,
			Config: at.JoinConfigs(
				workspaceResourceHCL,
				at.CompileConfig(
					testResourceSparkCustomPoolHeader,
					map[string]any{
						"workspace_id": testhelp.RefByFQN(workspaceResourceFQN, "id"),
						"name":         entityCreateName,
						"type":         "Workspace",
						"node_family":  "MemoryOptimized",
						"node_size":    "Small",
						"auto_scale": map[string]any{
							"enabled":        true,
							"min_node_count": 1,
							"max_node_count": 3,
						},
						"dynamic_executor_allocation": map[string]any{
							"enabled":       true,
							"min_executors": 1,
							"max_executors": 2,
						},
					},
				)),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceSparkCustomPoolFQN, "name", entityCreateName),
				resource.TestCheckResourceAttr(testResourceSparkCustomPoolFQN, "auto_scale.enabled", "true"),
				resource.TestCheckResourceAttr(testResourceSparkCustomPoolFQN, "auto_scale.min_node_count", "1"),
				resource.TestCheckResourceAttr(testResourceSparkCustomPoolFQN, "auto_scale.max_node_count", "3"),
				resource.TestCheckResourceAttr(testResourceSparkCustomPoolFQN, "dynamic_executor_allocation.enabled", "true"),
				resource.TestCheckResourceAttr(testResourceSparkCustomPoolFQN, "dynamic_executor_allocation.min_executors", "1"),
				resource.TestCheckResourceAttr(testResourceSparkCustomPoolFQN, "dynamic_executor_allocation.max_executors", "2"),
			),
		},
		// Update and Read
		{
			ResourceName: testResourceSparkCustomPoolFQN,
			Config: at.JoinConfigs(
				workspaceResourceHCL,
				at.CompileConfig(
					testResourceSparkCustomPoolHeader,
					map[string]any{
						"workspace_id": testhelp.RefByFQN(workspaceResourceFQN, "id"),
						"name":         entityUpdateName,
						"type":         "Workspace",
						"node_family":  "MemoryOptimized",
						"node_size":    "Small",
						"auto_scale": map[string]any{
							"enabled":        false,
							"min_node_count": 1,
							"max_node_count": 3,
						},
						"dynamic_executor_allocation": map[string]any{
							"enabled": false,
						},
					},
				)),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceSparkCustomPoolFQN, "name", entityUpdateName),
				resource.TestCheckResourceAttr(testResourceSparkCustomPoolFQN, "auto_scale.enabled", "false"),
				resource.TestCheckResourceAttr(testResourceSparkCustomPoolFQN, "auto_scale.min_node_count", "1"),
				resource.TestCheckResourceAttr(testResourceSparkCustomPoolFQN, "auto_scale.max_node_count", "3"),
				resource.TestCheckResourceAttr(testResourceSparkCustomPoolFQN, "dynamic_executor_allocation.enabled", "false"),
				resource.TestCheckNoResourceAttr(testResourceSparkCustomPoolFQN, "dynamic_executor_allocation.min_executors"),
				resource.TestCheckNoResourceAttr(testResourceSparkCustomPoolFQN, "dynamic_executor_allocation.max_executors"),
			),
		},
	},
	))
}
