// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package onelakeshortcut_test

import (
	"errors"
	"fmt"
	"regexp"
	"testing"

	at "github.com/dcarbone/terraform-plugin-framework-utils/v3/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp/fakes"
)

var testResourceItemFQN, testResourceItemHeader = testhelp.TFResource(common.ProviderTypeName, itemTypeInfo.Type, "test")

func TestUnit_OneLakeShortcutResource_Attributes(t *testing.T) {
	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		// error - no attributes
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{},
			),
			ExpectError: regexp.MustCompile(`Missing required argument`),
		},
		// error - unexpected attribute
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id":    testhelp.RandomUUID(),
					"item_id":         testhelp.RandomUUID(),
					"unexpected_attr": "test",
				},
			),
			ExpectError: regexp.MustCompile(`An argument named "unexpected_attr" is not expected here`),
		},
		// error - no required attribute item_id
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": testhelp.RandomUUID(),
				},
			),
			ExpectError: regexp.MustCompile(`The argument "item_id" is required, but no definition was found.`),
		},
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": "invalid uuid",
					"item_id":      testhelp.RandomUUID(),
					"name":         testhelp.RandomName(),
					"path":         testhelp.RandomName(),
				},
			),
			ExpectError: regexp.MustCompile(`The argument "target" is required, but no definition was found.`),
		},
		// error - no required attribute target
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": "invalid uuid",
					"item_id":      testhelp.RandomUUID(),
					"name":         testhelp.RandomName(),
					"path":         testhelp.RandomName(),
					"target": map[string]any{
						"onelake": map[string]any{
							"workspace_id": testhelp.RandomUUID(),
							"item_id":      testhelp.RandomUUID(),
							"path":         testhelp.RandomName(),
						},
					},
				},
			),
			ExpectError: regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
		// error - unexpected data source attribute for target
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": testhelp.RandomUUID(),
					"item_id":      testhelp.RandomUUID(),
					"name":         testhelp.RandomName(),
					"path":         testhelp.RandomName(),
					"target": map[string]any{
						"unexpected_attr": map[string]any{
							"workspace_id": testhelp.RandomUUID(),
							"item_id":      testhelp.RandomUUID(),
							"path":         testhelp.RandomName(),
						},
					},
				},
			),
			ExpectError: regexp.MustCompile(`Exactly one target type must be specified`),
		},
		// error - multiple target data sources
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": testhelp.RandomUUID(),
					"item_id":      testhelp.RandomUUID(),
					"name":         testhelp.RandomName(),
					"path":         testhelp.RandomName(),
					"target": map[string]any{
						"onelake": map[string]any{
							"workspace_id": testhelp.RandomUUID(),
							"item_id":      testhelp.RandomUUID(),
							"path":         testhelp.RandomName(),
						},
						"dataverse": map[string]any{
							"connection_id":      testhelp.RandomUUID(),
							"table_name":         testhelp.RandomName(),
							"deltalake_folder":   testhelp.RandomName(),
							"path":               testhelp.RandomName(),
							"environment_domain": testhelp.RandomName(),
						},
					},
				},
			),
			ExpectError: regexp.MustCompile(`Exactly one target type must be specified`),
		},
		// error - no target data source
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": testhelp.RandomUUID(),
					"item_id":      testhelp.RandomUUID(),
					"name":         testhelp.RandomName(),
					"path":         testhelp.RandomName(),
					"target":       map[string]any{},
				},
			),
			ExpectError: regexp.MustCompile(`Exactly one target type must be specified`),
		},
	}))
}

func TestUnit_OneLakeResource_ImportState(t *testing.T) {
	workspaceID := testhelp.RandomUUID()
	itemID := testhelp.RandomUUID()
	entity := NewRandomOnelakeShortcut()

	fakeTestUpsert(workspaceID, itemID, entity)

	fakes.FakeServer.ServerFactory.Core.OneLakeShortcutsServer.GetShortcut = fakeGetOneLakeShortcutFunc()
	fakes.FakeServer.ServerFactory.Core.OneLakeShortcutsServer.DeleteShortcut = fakeDeleteOneLakeShortcutFunc()

	testCase := at.CompileConfig(
		testResourceItemHeader,
		map[string]any{
			"workspace_id": workspaceID,
			"item_id":      itemID,
			"path":         *entity.Path,
			"name":         *entity.Name,
			"target": map[string]any{
				"onelake": map[string]any{
					"workspace_id": *entity.Target.OneLake.WorkspaceID,
					"item_id":      *entity.Target.OneLake.ItemID,
					"path":         *entity.Target.OneLake.Path,
				},
			},
		},
	)

	resource.Test(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		{
			ResourceName:  testResourceItemFQN,
			Config:        testCase,
			ImportStateId: "not-valid",
			ImportState:   true,
			ExpectError:   regexp.MustCompile(fmt.Sprintf(common.ErrorImportIdentifierDetails, "WorkspaceID/ItemID/Path/Name")),
		},
		{
			ResourceName:  testResourceItemFQN,
			Config:        testCase,
			ImportStateId: "test/id/test/test",
			ImportState:   true,
			ExpectError:   regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
		{
			ResourceName:  testResourceItemFQN,
			Config:        testCase,
			ImportStateId: fmt.Sprintf("%s/%s/%s/%s", "test", itemID, *entity.Path, *entity.Name),
			ImportState:   true,
			ExpectError:   regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
		// Import state testing
		{
			ResourceName:       testResourceItemFQN,
			Config:             testCase,
			ImportStateId:      fmt.Sprintf("%s/%s/%s/%s", workspaceID, itemID, *entity.Path, *entity.Name),
			ImportState:        true,
			ImportStatePersist: true,
			ImportStateCheck: func(is []*terraform.InstanceState) error {
				if len(is) != 1 {
					return errors.New("expected one instance state")
				}
				expectedID := fmt.Sprintf("%s%s%s%s", workspaceID, itemID, *entity.Path, *entity.Name)

				if is[0].ID != expectedID {
					return fmt.Errorf("%s: unexpected ID — got %q, want %q", testResourceItemFQN, is[0].ID, expectedID)
				}

				return nil
			},
		},
	}))
}

func TestUnit_OneLakeShortcuResource_CRUD(t *testing.T) {
	workspaceID := testhelp.RandomUUID()
	itemID := testhelp.RandomUUID()
	entityExist := NewRandomOnelakeShortcut()
	entityBefore := NewRandomOnelakeShortcut()
	entityAfter := NewRandomOnelakeShortcut()

	fakeTestUpsert(workspaceID, itemID, entityExist)

	fakes.FakeServer.ServerFactory.Core.OneLakeShortcutsServer.GetShortcut = fakeGetOneLakeShortcutFunc()
	fakes.FakeServer.ServerFactory.Core.OneLakeShortcutsServer.CreateShortcut = fakeCreateOneLakeShortcutFunc()
	fakes.FakeServer.ServerFactory.Core.OneLakeShortcutsServer.DeleteShortcut = fakeDeleteOneLakeShortcutFunc()

	resource.Test(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		// error - create - existing entity
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": workspaceID,
					"item_id":      itemID,
					"name":         *entityExist.Name,
					"path":         *entityExist.Path,
					"target": map[string]any{
						"onelake": map[string]any{
							"workspace_id": *entityExist.Target.OneLake.WorkspaceID,
							"item_id":      *entityExist.Target.OneLake.ItemID,
							"path":         *entityExist.Target.OneLake.Path,
						},
					},
				},
			),
			ExpectError: regexp.MustCompile(common.ErrorCreateHeader),
		},
		// Create and Read
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": workspaceID,
					"item_id":      itemID,
					"name":         *entityBefore.Name,
					"path":         *entityBefore.Path,
					"target": map[string]any{
						"onelake": map[string]any{
							"workspace_id": *entityBefore.Target.OneLake.WorkspaceID,
							"item_id":      *entityBefore.Target.OneLake.ItemID,
							"path":         *entityBefore.Target.OneLake.Path,
						},
					},
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "name", entityBefore.Name),
				resource.TestCheckResourceAttr(testResourceItemFQN, "path", *entityBefore.Path),
			),
		},
		// Update and Read
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": workspaceID,
					"item_id":      itemID,
					"name":         *entityBefore.Name,
					"path":         *entityBefore.Path,
					"target": map[string]any{
						"onelake": map[string]any{
							"workspace_id": *entityBefore.Target.OneLake.WorkspaceID,
							"item_id":      *entityAfter.Target.OneLake.ItemID,
							"path":         *entityBefore.Target.OneLake.Path,
						},
					},
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "name", entityBefore.Name),
				resource.TestCheckResourceAttr(testResourceItemFQN, "path", *entityBefore.Path),
			),
		},
		// Delete testing automatically occurs in TestCase
	}))
}

func TestAcc_OneLakeShortcutResource_CRUD(t *testing.T) {
	entityCreateDisplayName := testhelp.RandomName()
	entityTargetPath := "Files/images"
	entityUpdatedTargetPath := "Files/sample_dataset"
	workspaceID := testhelp.WellKnown()["WorkspaceDS"].(map[string]any)["id"].(string)
	lakehouseID := testhelp.WellKnown()["Lakehouse"].(map[string]any)["id"].(string)
	resource.Test(t, testhelp.NewTestAccCase(t, &testResourceItemFQN, nil, []resource.TestStep{
		// Create and Read
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"item_id":      lakehouseID,
					"workspace_id": workspaceID,
					"name":         entityCreateDisplayName,
					"path":         "/Files",
					"target": map[string]any{
						"onelake": map[string]any{
							"workspace_id": workspaceID,
							"item_id":      lakehouseID,
							"path":         entityTargetPath,
						},
					},
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "name", entityCreateDisplayName),
				resource.TestCheckResourceAttr(testResourceItemFQN, "target.onelake.item_id", lakehouseID),
				resource.TestCheckResourceAttr(testResourceItemFQN, "target.onelake.workspace_id", workspaceID),
				resource.TestCheckResourceAttr(testResourceItemFQN, "target.onelake.path", entityTargetPath),
				resource.TestCheckResourceAttr(testResourceItemFQN, "target.type", "OneLake"),
			),
		},
		// Update and Read
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"item_id":      lakehouseID,
					"workspace_id": workspaceID,
					"name":         entityCreateDisplayName,
					"path":         "/Files",
					"target": map[string]any{
						"onelake": map[string]any{
							"workspace_id": workspaceID,
							"item_id":      lakehouseID,
							"path":         entityUpdatedTargetPath,
						},
					},
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "name", entityCreateDisplayName),
				resource.TestCheckResourceAttr(testResourceItemFQN, "target.onelake.item_id", lakehouseID),
				resource.TestCheckResourceAttr(testResourceItemFQN, "target.onelake.workspace_id", workspaceID),
				resource.TestCheckResourceAttr(testResourceItemFQN, "target.onelake.path", entityUpdatedTargetPath),
				resource.TestCheckResourceAttr(testResourceItemFQN, "target.type", "OneLake"),
			),
		},
	},
	))
}
