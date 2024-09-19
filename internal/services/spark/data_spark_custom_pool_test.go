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
	testDataSourceSparkCustomPoolFQN    = testhelp.DataSourceFQN("fabric", sparkCustomPoolTFName, "test")
	testDataSourceSparkCustomPoolHeader = at.DataSourceHeader(testhelp.TypeName("fabric", sparkCustomPoolTFName), "test")
)

func TestAcc_SparkCustomPoolDataSource(t *testing.T) {
	workspaceResourceHCL, workspaceResourceFQN := testhelp.TestAccWorkspaceResource(t, *testhelp.WellKnown().Capacity.ID)
	entityName := testhelp.RandomName()

	resource.ParallelTest(t, testhelp.NewTestAccCase(t, &testDataSourceSparkCustomPoolFQN, nil, []resource.TestStep{
		// read
		{
			ResourceName: testDataSourceSparkCustomPoolFQN,
			Config: at.JoinConfigs(
				workspaceResourceHCL,
				at.CompileConfig(
					testResourceSparkCustomPoolHeader,
					getSparkCustomPoolResourceAttr(t, testhelp.RefByFQN(workspaceResourceFQN, "id"), entityName),
				),
				at.CompileConfig(
					testDataSourceSparkCustomPoolHeader,
					map[string]any{
						"workspace_id": testhelp.RefByFQN(workspaceResourceFQN, "id"),
						"id":           testhelp.RefByFQN(testResourceSparkCustomPoolFQN, "id"),
					},
				)),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrSet(testDataSourceSparkCustomPoolFQN, "workspace_id"),
				resource.TestCheckResourceAttrSet(testDataSourceSparkCustomPoolFQN, "id"),
				resource.TestCheckResourceAttr(testDataSourceSparkCustomPoolFQN, "name", entityName),
			),
		},
	},
	))
}
