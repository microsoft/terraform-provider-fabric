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
	workspaceResourceHCL, workspaceResourceFQN := testhelp.TestAccWorkspaceResource(t, *testhelp.WellKnown().Capacity.ID)

	resource.ParallelTest(t, testhelp.NewTestAccCase(t, &testDataSourceSparkWorkspaceSettingsFQN, nil, []resource.TestStep{
		// read
		{
			ResourceName: testDataSourceSparkWorkspaceSettingsFQN,
			Config: at.JoinConfigs(
				workspaceResourceHCL,
				at.CompileConfig(
					testDataSourceSparkWorkspaceSettingsHeader,
					map[string]any{
						"workspace_id": testhelp.RefByFQN(workspaceResourceFQN, "id"),
					},
				)),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrSet(testDataSourceSparkWorkspaceSettingsFQN, "workspace_id"),
				resource.TestCheckResourceAttrSet(testDataSourceSparkWorkspaceSettingsFQN, "id"),
				resource.TestCheckResourceAttr(testDataSourceSparkWorkspaceSettingsFQN, "pool.default_pool.name", "Starter Pool"),
			),
		},
	},
	))
}
