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
	testDataSourceSparkEnvironmentSettingsFQN    = testhelp.DataSourceFQN("fabric", sparkEnvironmentSettingsTFName, "test")
	testDataSourceSparkEnvironmentSettingsHeader = at.DataSourceHeader(testhelp.TypeName("fabric", sparkEnvironmentSettingsTFName), "test")
)

func TestAcc_SparkEnvironmentSettingsDataSource(t *testing.T) {
	environment := testhelp.WellKnown().Environment
	workspace := testhelp.WellKnown().Workspace

	resource.ParallelTest(t, testhelp.NewTestAccCase(t, &testDataSourceSparkEnvironmentSettingsFQN, nil, []resource.TestStep{
		// read - Published
		{
			ResourceName: testDataSourceSparkEnvironmentSettingsFQN,
			Config: at.CompileConfig(
				testDataSourceSparkEnvironmentSettingsHeader,
				map[string]any{
					"workspace_id":       *workspace.ID,
					"environment_id":     *environment.ID,
					"publication_status": "Published",
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrSet(testDataSourceSparkEnvironmentSettingsFQN, "workspace_id"),
				resource.TestCheckResourceAttrSet(testDataSourceSparkEnvironmentSettingsFQN, "id"),
				resource.TestCheckResourceAttr(testDataSourceSparkEnvironmentSettingsFQN, "pool.name", "Starter Pool"),
			),
		},
		// read - Staging
		{
			ResourceName: testDataSourceSparkEnvironmentSettingsFQN,
			Config: at.CompileConfig(
				testDataSourceSparkEnvironmentSettingsHeader,
				map[string]any{
					"workspace_id":       *workspace.ID,
					"environment_id":     *environment.ID,
					"publication_status": "Staging",
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrSet(testDataSourceSparkEnvironmentSettingsFQN, "workspace_id"),
				resource.TestCheckResourceAttrSet(testDataSourceSparkEnvironmentSettingsFQN, "id"),
				resource.TestCheckResourceAttr(testDataSourceSparkEnvironmentSettingsFQN, "pool.name", "Starter Pool"),
			),
		},
	},
	))
}
