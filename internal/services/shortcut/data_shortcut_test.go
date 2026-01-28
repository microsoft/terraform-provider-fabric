// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package shortcut_test

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

func TestUnit_ShortcutDataSource(t *testing.T) {
	workspaceID := testhelp.RandomUUID()
	itemID := testhelp.RandomUUID()
	entity := NewRandomShortcut()

	fakeTestUpsert(workspaceID, itemID, entity)

	fakes.FakeServer.ServerFactory.Core.OneLakeShortcutsServer.GetShortcut = fakeGetShortcutFunc()

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
					"item_id":      itemID,
				},
			),
			ExpectError: regexp.MustCompile(`The argument "path" is required, but no definition was found.`),
		},
		// missing name attribute
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"workspace_id": workspaceID,
					"item_id":      itemID,
					"path":         *entity.Path,
				},
			),
			ExpectError: regexp.MustCompile(`The argument "name" is required, but no definition was found.`),
		},
		// read
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"workspace_id": workspaceID,
					"item_id":      itemID,
					"name":         *entity.Name,
					"path":         *entity.Path,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "name", *entity.Name),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "actual_name", *entity.Name),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "path", *entity.Path),
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
					"path":         *entity.Path,
				},
			),
			ExpectError: regexp.MustCompile(common.ErrorReadHeader),
		},
	}))
}

func TestAcc_ShortcutDataSource(t *testing.T) {
	onelakeShortcut := testhelp.WellKnown()["Shortcut"]
	workspaceID := onelakeShortcut.(map[string]any)["workspaceId"].(string)
	itemID := onelakeShortcut.(map[string]any)["lakehouseId"].(string)
	shortcutName := onelakeShortcut.(map[string]any)["shortcutName"].(string)
	shortcutPath := onelakeShortcut.(map[string]any)["shortcutPath"].(string)

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
					"name":         shortcutName,
					"path":         shortcutPath,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "name", shortcutName),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "actual_name", shortcutName),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "path", shortcutPath),
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
					"path":         shortcutPath,
				},
			),
			ExpectError: regexp.MustCompile(common.ErrorReadHeader),
		},
	}))
}
