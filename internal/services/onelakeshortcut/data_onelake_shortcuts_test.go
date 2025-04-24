// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package onelakeshortcut_test

import (
	"testing"

	at "github.com/dcarbone/terraform-plugin-framework-utils/v3/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
)

var testDataSourceItemsFQN, testDataSourceItemsHeader = testhelp.TFDataSource(common.ProviderTypeName, itemTypeInfo.Types, "test")

func TestAcc_WorkspaceManagedPrivateEndpointsDataSource(t *testing.T) {
	workspace := testhelp.WellKnown()["WorkspaceDS"].(map[string]any)
	workspaceID := workspace["id"].(string)
	itemID := testhelp.WellKnown()["Lakehouse"].(map[string]any)["id"].(string)

	resource.ParallelTest(t, testhelp.NewTestAccCase(t, nil, nil, []resource.TestStep{
		// read
		{
			Config: at.CompileConfig(
				testDataSourceItemsHeader,
				map[string]any{
					"workspace_id": workspaceID,
					"item_id":      itemID,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrSet(testDataSourceItemsFQN, "values.0.name"),
				resource.TestCheckResourceAttrSet(testDataSourceItemsFQN, "values.0.path"),
				resource.TestCheckResourceAttrSet(testDataSourceItemsFQN, "values.0.target.type"),
				resource.TestCheckResourceAttrSet(testDataSourceItemsFQN, "values.0.target.onelake.item_id"),
				resource.TestCheckResourceAttrSet(testDataSourceItemsFQN, "values.0.target.onelake.path"),
				resource.TestCheckResourceAttrSet(testDataSourceItemsFQN, "values.0.target.onelake.workspace_id"),
			),
		},
	},
	))
}
