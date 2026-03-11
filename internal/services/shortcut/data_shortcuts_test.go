// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package shortcut_test

import (
	"regexp"
	"testing"

	at "github.com/dcarbone/terraform-plugin-framework-utils/v3/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp/fakes"
)

var testDataSourceItemsFQN, testDataSourceItemsHeader = testhelp.TFDataSource(common.ProviderTypeName, itemTypeInfo.Types, "test")

func TestUnit_ShortcutsDataSource(t *testing.T) {
	workspaceID := testhelp.RandomUUID()
	itemID := testhelp.RandomUUID()

	fakeTestUpsert(workspaceID, itemID, NewRandomShortcut())
	fakeTestUpsert(workspaceID, itemID, NewRandomShortcut())

	fakes.FakeServer.ServerFactory.Core.OneLakeShortcutsServer.NewListShortcutsPager = fakeShortcutsFunc()

	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, nil, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		// error - no attributes
		{
			Config: at.CompileConfig(
				testDataSourceItemsHeader,
				map[string]any{},
			),
			ExpectError: regexp.MustCompile(`The argument "workspace_id" is required, but no definition was found`),
		},
		// error - workspace_id - invalid UUID
		{
			Config: at.CompileConfig(
				testDataSourceItemsHeader,
				map[string]any{
					"workspace_id": "invalid uuid",
					"item_id":      "invalid uuid",
				},
			),
			ExpectError: regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
		// error - unexpected_attr
		{
			Config: at.CompileConfig(
				testDataSourceItemsHeader,
				map[string]any{
					"workspace_id":    workspaceID,
					"item_id":         itemID,
					"unexpected_attr": "test",
				},
			),
			ExpectError: regexp.MustCompile(`An argument named "unexpected_attr" is not expected here`),
		},
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
				resource.TestCheckResourceAttrSet(testDataSourceItemsFQN, "values.0.actual_name"),
				resource.TestCheckResourceAttrSet(testDataSourceItemsFQN, "values.0.path"),
				resource.TestCheckResourceAttrSet(testDataSourceItemsFQN, "values.0.target.onelake.item_id"),
				resource.TestCheckResourceAttrSet(testDataSourceItemsFQN, "values.0.target.onelake.path"),
				resource.TestCheckResourceAttrSet(testDataSourceItemsFQN, "values.0.target.onelake.workspace_id"),
			),
		},
	}))
}

func TestAcc_ShortcutsDataSource(t *testing.T) {
	onelakeShortcut := testhelp.WellKnown()["Shortcut"]
	workspaceID := onelakeShortcut.(map[string]any)["workspaceId"].(string)
	itemID := onelakeShortcut.(map[string]any)["lakehouseId"].(string)

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
				resource.TestCheckResourceAttrSet(testDataSourceItemsFQN, "values.0.actual_name"),
				resource.TestCheckResourceAttrSet(testDataSourceItemsFQN, "values.0.path"),
				resource.TestCheckResourceAttrSet(testDataSourceItemsFQN, "values.0.target.onelake.item_id"),
				resource.TestCheckResourceAttrSet(testDataSourceItemsFQN, "values.0.target.onelake.path"),
				resource.TestCheckResourceAttrSet(testDataSourceItemsFQN, "values.0.target.onelake.workspace_id"),
			),
		},
	},
	))
}
