// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package folder_test

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

func TestUnit_FolderResource_Attributes(t *testing.T) {
	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		// error - no attributes
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{},
				),
			),
			ExpectError: regexp.MustCompile(`Missing required argument`),
		},
		// error - workspace_id - invalid UUID
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"workspace_id": "invalid uuid",
						"display_name": "test",
					},
				)),
			ExpectError: regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
		// error - unexpected attribute
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"workspace_id":    "00000000-0000-0000-0000-000000000000",
						"display_name":    "test",
						"unexpected_attr": "test",
					},
				)),
			ExpectError: regexp.MustCompile(`An argument named "unexpected_attr" is not expected here`),
		},
		// error - no required attributes
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"display_name": "test",
					},
				)),
			ExpectError: regexp.MustCompile(`The argument "workspace_id" is required, but no definition was found.`),
		},
		// error - no required attributes
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"workspace_id": "00000000-0000-0000-0000-000000000000",
					},
				)),
			ExpectError: regexp.MustCompile(`The argument "display_name" is required, but no definition was found.`),
		},
	}))
}

func TestUnit_FolderResource_ImportState(t *testing.T) {
	workspaceID := testhelp.RandomUUID()
	entity := fakes.NewRandomFolderWithWorkspace(workspaceID)

	fakes.FakeServer.Upsert(fakes.NewRandomFolderWithWorkspace(workspaceID))
	fakes.FakeServer.Upsert(entity)
	fakes.FakeServer.Upsert(fakes.NewRandomFolderWithWorkspace(workspaceID))

	testCase := at.JoinConfigs(
		at.CompileConfig(
			testResourceItemHeader,
			map[string]any{
				"workspace_id": *entity.WorkspaceID,
				"display_name": *entity.DisplayName,
			},
		))

	resource.Test(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		{
			ResourceName:  testResourceItemFQN,
			Config:        testCase,
			ImportStateId: "not-valid",
			ImportState:   true,
			ExpectError:   regexp.MustCompile(fmt.Sprintf(common.ErrorImportIdentifierDetails, "WorkspaceID/FolderID")),
		},
		{
			ResourceName:  testResourceItemFQN,
			Config:        testCase,
			ImportStateId: "test/id",
			ImportState:   true,
			ExpectError:   regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
		{
			ResourceName:  testResourceItemFQN,
			Config:        testCase,
			ImportStateId: fmt.Sprintf("%s/%s", "test", *entity.ID),
			ImportState:   true,
			ExpectError:   regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
		{
			ResourceName:  testResourceItemFQN,
			Config:        testCase,
			ImportStateId: fmt.Sprintf("%s/%s", *entity.WorkspaceID, "test"),
			ImportState:   true,
			ExpectError:   regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
		// Import state testing
		{
			ResourceName:       testResourceItemFQN,
			Config:             testCase,
			ImportStateId:      fmt.Sprintf("%s/%s", *entity.WorkspaceID, *entity.ID),
			ImportState:        true,
			ImportStatePersist: true,
			ImportStateCheck: func(is []*terraform.InstanceState) error {
				if len(is) != 1 {
					return errors.New("expected one instance state")
				}

				if is[0].ID != *entity.ID {
					return errors.New(testResourceItemFQN + ": unexpected ID")
				}

				return nil
			},
		},
	}))
}

func TestUnit_FolderResource_CRUD(t *testing.T) {
	workspaceID := testhelp.RandomUUID()
	entityExist := fakes.NewRandomFolderWithWorkspace(workspaceID)
	entityBefore := fakes.NewRandomFolderWithWorkspace(workspaceID)
	entityAfter := fakes.NewRandomFolderWithWorkspace(workspaceID)

	fakes.FakeServer.Upsert(fakes.NewRandomFolderWithWorkspace(workspaceID))
	fakes.FakeServer.Upsert(entityExist)
	fakes.FakeServer.Upsert(fakes.NewRandomFolderWithWorkspace(workspaceID))

	resource.Test(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		// error - create - existing entity
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"workspace_id":     *entityExist.WorkspaceID,
						"display_name":     *entityExist.DisplayName,
						"parent_folder_id": *entityExist.ParentFolderID,
					},
				)),
			ExpectError: regexp.MustCompile(common.ErrorCreateHeader),
		},
		// Create and Read
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"workspace_id":     *entityBefore.WorkspaceID,
						"display_name":     *entityBefore.DisplayName,
						"parent_folder_id": *entityBefore.ParentFolderID,
					},
				)),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "display_name", entityBefore.DisplayName),
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "parent_folder_id", entityBefore.ParentFolderID),
			),
		},
		// Move and Read
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"workspace_id":     *entityBefore.WorkspaceID,
						"display_name":     *entityBefore.DisplayName,
						"parent_folder_id": *entityAfter.ParentFolderID,
					},
				)),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "display_name", entityBefore.DisplayName),
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "parent_folder_id", entityAfter.ParentFolderID),
			),
		},

		// Update and Read
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"workspace_id":     *entityBefore.WorkspaceID,
						"display_name":     *entityAfter.DisplayName,
						"parent_folder_id": *entityAfter.ParentFolderID,
					},
				)),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "display_name", entityAfter.DisplayName),
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "parent_folder_id", entityAfter.ParentFolderID),
			),
		},
		// Update, Move and Read
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"workspace_id":     *entityBefore.WorkspaceID,
						"display_name":     *entityBefore.DisplayName,
						"parent_folder_id": *entityBefore.ParentFolderID,
					},
				)),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "display_name", entityBefore.DisplayName),
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "parent_folder_id", entityBefore.ParentFolderID),
			),
		},
		// Delete testing automatically occurs in TestCase
	}))
}

func TestAcc_FolderResource_CRUD(t *testing.T) {
	workspace := testhelp.WellKnown()["WorkspaceDS"].(map[string]any)
	workspaceID := workspace["id"].(string)

	entityCreateDisplayName := testhelp.RandomName()
	entityUpdateDisplayName := testhelp.RandomName()
	folder := testhelp.WellKnown()["Folder"].(map[string]any)
	folderID := folder["id"].(string)

	resource.Test(t, testhelp.NewTestAccCase(t, &testResourceItemFQN, nil, []resource.TestStep{
		// Create and Read
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"workspace_id": workspaceID,
						"display_name": entityCreateDisplayName,
					},
				)),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "display_name", entityCreateDisplayName),
			),
		},
		// Update, Move and Read
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"workspace_id":     workspaceID,
						"display_name":     entityUpdateDisplayName,
						"parent_folder_id": folderID,
					},
				)),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "display_name", entityUpdateDisplayName),
				resource.TestCheckResourceAttr(testResourceItemFQN, "parent_folder_id", folderID),
			),
		},
		// Update and Read
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"workspace_id":     workspaceID,
						"display_name":     entityCreateDisplayName,
						"parent_folder_id": folderID,
					},
				)),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "display_name", entityCreateDisplayName),
				resource.TestCheckResourceAttr(testResourceItemFQN, "parent_folder_id", folderID),
			),
		},
		// Move and Read
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"workspace_id": workspaceID,
						"display_name": entityCreateDisplayName,
					},
				)),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "display_name", entityCreateDisplayName),
				resource.TestCheckNoResourceAttr(testResourceItemFQN, "parent_folder_id"),
			),
		},
	},
	))
}
