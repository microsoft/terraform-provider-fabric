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
		// Update and Read (Spark properties)
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
						"spark_properties": map[string]any{
							`"spark.acls.enable"`: "true",
						},
					},
				)),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "pool.name", "Starter Pool"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "driver_cores", "8"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "spark_properties.spark.acls.enable"),
			),
		},
	},
	))
}
