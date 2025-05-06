// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package onelakeshortcut_test

import (
	//"errors"
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
					"unexpected_attr": "test",
				},
			),
			ExpectError: regexp.MustCompile(`An argument named "unexpected_attr" is not expected here`),
		},
		// // error - no required attributes
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
		// // error - invalid uuid - capacity_id
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": "invalid uuid",
				},
			),
			ExpectError: regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
	}))
}

func TestUnit_LakehouseResource_ImportState(t *testing.T) {
	workspaceID := testhelp.RandomUUID()
	itemId := testhelp.RandomUUID()
	entity := fakes.NewRandomOnelakeShortcut()

	fakes.FakeServer.Upsert(fakes.NewRandomOnelakeShortcut())
	fakes.FakeServer.Upsert(entity)
	fakes.FakeServer.Upsert(fakes.NewRandomOnelakeShortcut())

	testCase := at.CompileConfig(
		testResourceItemHeader,
		map[string]any{
			"workspace_id": workspaceID,
			"item_id":      itemId,
			"path":         *entity.Path,
			"name":         *entity.Name,
		},
	)

	resource.Test(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		// {
		// 	ResourceName:  testResourceItemFQN,
		// 	Config:        testCase,
		// 	ImportStateId: "not-valid",
		// 	ImportState:   true,
		// 	ExpectError:   regexp.MustCompile(fmt.Sprintf(common.ErrorImportIdentifierDetails, "<WorkspaceID>/<ItemID>/<Path>/<Name>")),
		// },
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
			ImportStateId: fmt.Sprintf("%s/%s/%s/%s", "test", itemId, *entity.Path, *entity.Name),
			ImportState:   true,
			ExpectError:   regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
		// Import state testing
		{
			ResourceName:       testResourceItemFQN,
			Config:             testCase,
			ImportStateId:      fmt.Sprintf("%s/%s/%s/%s", workspaceID, itemId, *entity.Path, *entity.Name),
			ImportState:        true,
			ImportStatePersist: true,
			ImportStateCheck: func(is []*terraform.InstanceState) error {
				if len(is) != 1 {
					return errors.New("expected one instance state")
				}

				if is[0].ID != workspaceID {
					return errors.New(testResourceItemFQN + ": unexpected ID")
				}

				return nil
			},
		},
	}))
}

func TestAcc_OneLakeShortcutResource_CRUD(t *testing.T) {
	entityCreateDisplayName := testhelp.RandomName()

	entityTargetPath := "Files/images"
	entityUpdatedTargetPath := "Files/sample_dataset"
	workspaceId := testhelp.WellKnown()["WorkspaceDS"].(map[string]any)["id"].(string)
	lakehouseId := testhelp.WellKnown()["Lakehouse"].(map[string]any)["id"].(string)
	resource.Test(t, testhelp.NewTestAccCase(t, &testResourceItemFQN, nil, []resource.TestStep{
		// Create and Read
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"item_id":      lakehouseId,
					"workspace_id": workspaceId,
					"name":         entityCreateDisplayName,
					"path":         "Files",
					"target": map[string]any{
						"onelake": map[string]any{
							"workspace_id": workspaceId,
							"item_id":      lakehouseId,
							"path":         entityTargetPath,
						},
					},
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "name", entityCreateDisplayName),
				resource.TestCheckResourceAttr(testResourceItemFQN, "target.onelake.item_id", lakehouseId),
				resource.TestCheckResourceAttr(testResourceItemFQN, "target.onelake.workspace_id", workspaceId),
				resource.TestCheckResourceAttr(testResourceItemFQN, "target.onelake.path", entityTargetPath),
			),
		},
		// Update - Create with OverwriteOnly and Read
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"item_id":      lakehouseId,
					"workspace_id": workspaceId,
					"name":         entityCreateDisplayName,
					"path":         "Files",
					"target": map[string]any{
						"onelake": map[string]any{
							"workspace_id": workspaceId,
							"item_id":      lakehouseId,
							"path":         entityUpdatedTargetPath,
						},
					},
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "name", entityCreateDisplayName),
				resource.TestCheckResourceAttr(testResourceItemFQN, "target.onelake.item_id", lakehouseId),
				resource.TestCheckResourceAttr(testResourceItemFQN, "target.onelake.workspace_id", workspaceId),
				resource.TestCheckResourceAttr(testResourceItemFQN, "target.onelake.path", entityUpdatedTargetPath),
			),
		},
	},
	))
}
