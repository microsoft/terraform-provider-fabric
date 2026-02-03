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

var testDataSourceItemFQN, testDataSourceItemHeader = testhelp.TFDataSource(common.ProviderTypeName, itemTypeInfo.Type, "test")

func TestAcc_SparkEnvironmentSettingsDataSource(t *testing.T) {
	workspace := testhelp.WellKnown()["WorkspaceDS"].(map[string]any)
	workspaceID := workspace["id"].(string)

	environment := testhelp.WellKnown()["Environment"].(map[string]any)
	environmentID := environment["id"].(string)

	resource.ParallelTest(t, testhelp.NewTestAccCase(t, &testDataSourceItemFQN, nil, []resource.TestStep{
		// read - Published
		{
			ResourceName: testDataSourceItemFQN,
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"workspace_id":       workspaceID,
					"environment_id":     environmentID,
					"publication_status": "Published",
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrSet(testDataSourceItemFQN, "workspace_id"),
				resource.TestCheckResourceAttrSet(testDataSourceItemFQN, "id"),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "pool.name", "Starter Pool"),
			),
		},
		// read - Staging
		{
			ResourceName: testDataSourceItemFQN,
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"workspace_id":       workspaceID,
					"environment_id":     environmentID,
					"publication_status": "Staging",
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrSet(testDataSourceItemFQN, "workspace_id"),
				resource.TestCheckResourceAttrSet(testDataSourceItemFQN, "id"),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "pool.name", "Starter Pool"),
			),
		},
	},
	))
}
