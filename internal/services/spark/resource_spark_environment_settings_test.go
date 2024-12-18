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
	testResourceSparkEnvironmentSettingsFQN    = testhelp.ResourceFQN("fabric", sparkEnvironmentSettingsTFName, "test")
	testResourceSparkEnvironmentSettingsHeader = at.ResourceHeader(testhelp.TypeName("fabric", sparkEnvironmentSettingsTFName), "test")
)

func TestAcc_SparkEnvironmentSettingsResource_CRUD(t *testing.T) {
	capacity := testhelp.WellKnown()["Capacity"].(map[string]any)
	capacityID := capacity["id"].(string)

	workspaceResourceHCL, workspaceResourceFQN := testhelp.TestAccWorkspaceResource(t, capacityID)
	environmentResourceHCL, environmentResourceFQN := environmentResource(t, testhelp.RefByFQN(workspaceResourceFQN, "id"))

	resource.Test(t, testhelp.NewTestAccCase(t, &testResourceSparkEnvironmentSettingsFQN, nil, []resource.TestStep{
		// Create and Read
		{
			ResourceName: testResourceSparkEnvironmentSettingsFQN,
			Config: at.JoinConfigs(
				workspaceResourceHCL,
				environmentResourceHCL,
				at.CompileConfig(
					testResourceSparkEnvironmentSettingsHeader,
					map[string]any{
						"workspace_id":       testhelp.RefByFQN(workspaceResourceFQN, "id"),
						"environment_id":     testhelp.RefByFQN(environmentResourceFQN, "id"),
						"publication_status": "Published",
						"driver_cores":       4,
					},
				)),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceSparkEnvironmentSettingsFQN, "pool.name", "Starter Pool"),
				resource.TestCheckResourceAttr(testResourceSparkEnvironmentSettingsFQN, "driver_cores", "4"),
			),
		},
		// Update and Read
		{
			ResourceName: testResourceSparkEnvironmentSettingsFQN,
			Config: at.JoinConfigs(
				workspaceResourceHCL,
				environmentResourceHCL,
				at.CompileConfig(
					testResourceSparkEnvironmentSettingsHeader,
					map[string]any{
						"workspace_id":       testhelp.RefByFQN(workspaceResourceFQN, "id"),
						"environment_id":     testhelp.RefByFQN(environmentResourceFQN, "id"),
						"publication_status": "Published",
						"driver_cores":       8,
					},
				)),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceSparkEnvironmentSettingsFQN, "pool.name", "Starter Pool"),
				resource.TestCheckResourceAttr(testResourceSparkEnvironmentSettingsFQN, "driver_cores", "8"),
			),
		},
		// Update and Read (Spark properties)
		{
			ResourceName: testResourceSparkEnvironmentSettingsFQN,
			Config: at.JoinConfigs(
				workspaceResourceHCL,
				environmentResourceHCL,
				at.CompileConfig(
					testResourceSparkEnvironmentSettingsHeader,
					map[string]any{
						"workspace_id":       testhelp.RefByFQN(workspaceResourceFQN, "id"),
						"environment_id":     testhelp.RefByFQN(environmentResourceFQN, "id"),
						"publication_status": "Published",
						"driver_cores":       8,
						"spark_properties": map[string]any{
							`"spark.acls.enable"`: "true",
						},
					},
				)),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceSparkEnvironmentSettingsFQN, "pool.name", "Starter Pool"),
				resource.TestCheckResourceAttr(testResourceSparkEnvironmentSettingsFQN, "driver_cores", "8"),
				resource.TestCheckResourceAttrSet(testResourceSparkEnvironmentSettingsFQN, "spark_properties.spark.acls.enable"),
			),
		},
	},
	))
}
