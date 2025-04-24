// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package onelakeshortcut_test

import (
	"testing"

	at "github.com/dcarbone/terraform-plugin-framework-utils/v3/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp/fakes"
)

var testDataSourceItemFQN, testDataSourceItemHeader = testhelp.TFDataSource(common.ProviderTypeName, itemTypeInfo.Type, "test")

func TestUnit_OneLakeShortcutDataSource(t *testing.T) {
	workspaceID := testhelp.RandomUUID()
	itemId := testhelp.RandomUUID()
	entity := fakes.NewRandomOnelakeShortcut()

	fakes.FakeServer.Upsert(fakes.NewRandomOnelakeShortcut())
	fakes.FakeServer.Upsert(entity)
	fakes.FakeServer.Upsert(fakes.NewRandomOnelakeShortcut())

	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, nil, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		// read by name
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"item_id":      itemId,
					"workspace_id": workspaceID,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrPtr(testDataSourceItemFQN, "display_name", entity.Name),
			),
		},
	}))
}

func TestAcc_WorkspaceDataSource(t *testing.T) {
	workspace := testhelp.WellKnown()["WorkspaceDS"].(map[string]any)
	workspaceID := workspace["id"].(string)
	itemID := testhelp.WellKnown()["Lakehouse"].(map[string]any)["id"].(string)
	// onelake := testhelp.WellKnown()["OneLakeShortcut"].(map[string]any)

	resource.ParallelTest(t, testhelp.NewTestAccCase(t, nil, nil, []resource.TestStep{
		// read by id
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"workspace_id": workspaceID,
					"item_id":      itemID,
					"name":         "publicholidays_1",
					"path":         "/Tables",
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "item_id", itemID),
			),
		},
		// // read by id - not found
		// {
		// 	Config: at.CompileConfig(
		// 		testDataSourceItemHeader,
		// 		map[string]any{
		// 			"id": testhelp.RandomUUID(),
		// 		},
		// 	),
		// 	ExpectError: regexp.MustCompile(common.ErrorReadHeader),
		// },
		// // read by name
		// {
		// 	Config: at.CompileConfig(
		// 		testDataSourceItemHeader,
		// 		map[string]any{
		// 			"display_name": entityDisplayName,
		// 		},
		// 	),
		// 	Check: resource.ComposeAggregateTestCheckFunc(
		// 		resource.TestCheckResourceAttr(testDataSourceItemFQN, "id", entityID),
		// 		resource.TestCheckResourceAttr(testDataSourceItemFQN, "display_name", entityDisplayName),
		// 		resource.TestCheckResourceAttr(testDataSourceItemFQN, "description", entityDescription),
		// 	),
		// },
		// // read by name - not found
		// {
		// 	Config: at.CompileConfig(
		// 		testDataSourceItemHeader,
		// 		map[string]any{
		// 			"display_name": testhelp.RandomName(),
		// 		},
		// 	),
		// 	ExpectError: regexp.MustCompile(common.ErrorReadHeader),
		// },
	}))
}
