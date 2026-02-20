// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package sparkenvsettings_test

import (
	"testing"

	at "github.com/dcarbone/terraform-plugin-framework-utils/v3/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
)

var testResourceItemFQN, testResourceItemHeader = testhelp.TFResource(common.ProviderTypeName, itemTypeInfo.Type, "test")

func TestAcc_SparkEnvironmentSettingsResource_CRUD(t *testing.T) {
	capacity := testhelp.WellKnown()["Capacity"].(map[string]any)
	capacityID := capacity["id"].(string)

	workspaceResourceHCL, workspaceResourceFQN := testhelp.TestAccWorkspaceResource(t, capacityID)
	environmentResourceHCL, environmentResourceFQN := environmentResource(t, testhelp.RefByFQN(workspaceResourceFQN, "id"))

	resource.Test(t, testhelp.NewTestAccCase(t, &testResourceItemFQN, nil, []resource.TestStep{
		// Create and Read
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				workspaceResourceHCL,
				environmentResourceHCL,
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"workspace_id":       testhelp.RefByFQN(workspaceResourceFQN, "id"),
						"environment_id":     testhelp.RefByFQN(environmentResourceFQN, "id"),
						"publication_status": "Published",
						"driver_cores":       4,
					},
				)),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "pool.name", "Starter Pool"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "driver_cores", "4"),
			),
		},
		// Update and Read
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				workspaceResourceHCL,
				environmentResourceHCL,
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"workspace_id":       testhelp.RefByFQN(workspaceResourceFQN, "id"),
						"environment_id":     testhelp.RefByFQN(environmentResourceFQN, "id"),
						"publication_status": "Published",
						"driver_cores":       8,
					},
				)),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "pool.name", "Starter Pool"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "driver_cores", "8"),
			),
		},
	},
	))
}

func TestAcc_SparkEnvironmentSettingsSparkPropertiesResource_CRUD(t *testing.T) {
	capacity := testhelp.WellKnown()["Capacity"].(map[string]any)
	capacityID := capacity["id"].(string)

	workspaceResourceHCL, workspaceResourceFQN := testhelp.TestAccWorkspaceResource(t, capacityID)
	environmentResourceHCL, environmentResourceFQN := environmentResource(t, testhelp.RefByFQN(workspaceResourceFQN, "id"))

	resource.Test(t, testhelp.NewTestAccCase(t, &testResourceItemFQN, nil, []resource.TestStep{
		// Create and Read (Spark properties)
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				workspaceResourceHCL,
				environmentResourceHCL,
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"workspace_id":       testhelp.RefByFQN(workspaceResourceFQN, "id"),
						"environment_id":     testhelp.RefByFQN(environmentResourceFQN, "id"),
						"publication_status": "Published",
						"driver_cores":       8,
						"spark_properties": []map[string]any{
							{
								"key":   "spark.acls.enable",
								"value": "true",
							},
							{
								"key":   "spark.admin.acls.groups",
								"value": "test",
							},
						},
					},
				)),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "pool.name", "Starter Pool"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "driver_cores", "8"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "spark_properties.0.value"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "spark_properties.1.value"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "spark_properties.#", "2"),
			),
		},
		// Update and Read (test Spark properties sync)
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				workspaceResourceHCL,
				environmentResourceHCL,
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"workspace_id":       testhelp.RefByFQN(workspaceResourceFQN, "id"),
						"environment_id":     testhelp.RefByFQN(environmentResourceFQN, "id"),
						"publication_status": "Published",
						"driver_cores":       8,
						"spark_properties": []map[string]any{
							{
								"key":   "spark.cores.max",
								"value": "12",
							},
						},
					},
				)),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "pool.name", "Starter Pool"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "driver_cores", "8"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "spark_properties.0.value"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "spark_properties.0.key", "spark.cores.max"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "spark_properties.#", "1"),
			),
		},
		// Update and Read (remove Spark properties)
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				workspaceResourceHCL,
				environmentResourceHCL,
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"workspace_id":       testhelp.RefByFQN(workspaceResourceFQN, "id"),
						"environment_id":     testhelp.RefByFQN(environmentResourceFQN, "id"),
						"publication_status": "Published",
						"driver_cores":       8,
					},
				)),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "pool.name", "Starter Pool"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "driver_cores", "8"),
				resource.TestCheckNoResourceAttr(testResourceItemFQN, "spark_properties"),
			),
		},
	},
	))
}
