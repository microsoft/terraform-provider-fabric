// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package onelakeshortcut_test

import (
	"regexp"
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
		// error - no attributes
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{},
			),
			ExpectError: regexp.MustCompile(`The argument "workspace_id" is required, but no definition was found`),
		},
		// error - unexpected attribute
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"workspace_id":    workspaceID,
					"unexpected_attr": "test",
				},
			),
			ExpectError: regexp.MustCompile(`An argument named "unexpected_attr" is not expected here`),
		},
		// missing item_id attribute
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"workspace_id": workspaceID,
				},
			),
			ExpectError: regexp.MustCompile(`The argument "item_id" is required, but no definition was found`),
		},
		// missing path attribute
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"workspace_id": workspaceID,
					"item_id":      itemId,
				},
			),
			ExpectError: regexp.MustCompile(`parameter shortcutPath cannot be empty`),
		},
		// missing name attribute
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"workspace_id": workspaceID,
					"item_id":      itemId,
					"path":         "Files",
				},
			),
			ExpectError: regexp.MustCompile(`These attributes must be configured together: \[path,name\]`),
		},
		// read
		// {
		// 	Config: at.CompileConfig(
		// 		testDataSourceItemHeader,
		// 		map[string]any{
		// 			"workspace_id": workspaceID,
		// 			"item_id":      itemId,
		// 			"name":         *entity.Name,
		// 			"path":         *entity.Path,
		// 		},
		// 	),
		// Check: resource.ComposeAggregateTestCheckFunc(
		// 	resource.TestCheckResourceAttr(testDataSourceItemFQN, "name", entity.Name),
		// 	resource.TestCheckResourceAttr(testDataSourceItemFQN, "path", entity.Path),
		// 	resource.TestCheckResourceAttrSet(testDataSourceItemFQN, "target.onelake.path"),
		// 	resource.TestCheckResourceAttrSet(testDataSourceItemFQN, "target.onelake.workspace_id"),
		// 	resource.TestCheckResourceAttrSet(testDataSourceItemFQN, "target.onelake.item_id"),
		// ),
		// },
		// read - not found
		// {
		// 	Config: at.CompileConfig(
		// 		testDataSourceItemHeader,
		// 		map[string]any{
		// 			"workspace_id": workspaceID,
		// 			"item_id":      testhelp.RandomUUID(),
		// 			"name":         *entity.Name,
		// 			"path":         *entity.Path,
		// 		},
		// 	),
		// 	ExpectError: regexp.MustCompile(common.ErrorReadHeader),
		// },
	}))
}

func TestAcc_WorkspaceDataSource(t *testing.T) {
	workspaceID := testhelp.WellKnown()["WorkspaceDS"].(map[string]any)["id"].(string)
	itemID := testhelp.WellKnown()["Lakehouse"].(map[string]any)["id"].(string)
	tableName := testhelp.WellKnown()["Lakehouse"].(map[string]any)["tableName"].(string)

	resource.ParallelTest(t, testhelp.NewTestAccCase(t, nil, nil, []resource.TestStep{
		// error - no attributes
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{},
			),
			ExpectError: regexp.MustCompile(`The argument "workspace_id" is required, but no definition was found`),
		},
		// read
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"workspace_id": workspaceID,
					"item_id":      itemID,
					"name":         tableName,
					"path":         "Tables",
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "name", tableName),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "path", "Tables"),
				resource.TestCheckResourceAttrSet(testDataSourceItemFQN, "target.onelake.path"),
				resource.TestCheckResourceAttrSet(testDataSourceItemFQN, "target.onelake.workspace_id"),
				resource.TestCheckResourceAttrSet(testDataSourceItemFQN, "target.onelake.item_id"),
			),
		},
		// read - not found
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"workspace_id": workspaceID,
					"item_id":      itemID,
					"name":         testhelp.RandomName(),
					"path":         "Tables",
				},
			),
			ExpectError: regexp.MustCompile(common.ErrorReadHeader),
		},
	}))
}
