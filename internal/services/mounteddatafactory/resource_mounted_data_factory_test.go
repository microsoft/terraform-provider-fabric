// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package mounteddatafactory_test

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

var (
	testResourceItemFQN    = testhelp.ResourceFQN("fabric", itemTFName, "test")
	testResourceItemHeader = at.ResourceHeader(testhelp.TypeName("fabric", itemTFName), "test")
)

var testHelperLocals = at.CompileLocalsConfig(map[string]any{
	"path": testhelp.GetFixturesDirPath("mounted_data_factory"),
})

var (
	virtualNetworkAzureResource01 = testhelp.WellKnown()["VirtualNetwork01"].(map[string]any)
	vNET01SubscriptionID          = virtualNetworkAzureResource01["subscriptionId"].(string)
	vNET01ResourceGroupName       = virtualNetworkAzureResource01["resourceGroupName"].(string)

	factoryName = testhelp.WellKnown()["MountedDataFactory"].(map[string]any)["displayName"].(string)

	testHelperDefinition = map[string]any{
		`"mountedDataFactory-content.json"`: map[string]any{
			"source": "${local.path}/mountedDataFactory-content.json",
			"tokens": map[string]any{
				"subscriptionId":    vNET01SubscriptionID,
				"resourceGroupName": vNET01ResourceGroupName,
				"factoryName":       factoryName,
			},
		},
	}
)

func TestUnit_MountedDataFactoryResource_Attributes(t *testing.T) {
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
		// error - workspace_id - invalid UUID
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				testHelperLocals,
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"workspace_id": "invalid uuid",
						"display_name": "test",
						"format":       "Default",
						"definition":   testHelperDefinition,
					},
				),
			),
			ExpectError: regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
		// error - unexpected attribute
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				testHelperLocals,
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"workspace_id":    "00000000-0000-0000-0000-000000000000",
						"unexpected_attr": "test",
						"format":          "Default",
						"definition":      testHelperDefinition,
					},
				),
			),
			ExpectError: regexp.MustCompile(`An argument named "unexpected_attr" is not expected here`),
		},
		// error - no required attributes
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				testHelperLocals,
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"display_name": "test",
						"format":       "Default",
						"definition":   testHelperDefinition,
					},
				),
			),
			ExpectError: regexp.MustCompile(`The argument "workspace_id" is required, but no definition was found.`),
		},
		// error - no required attributes
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				testHelperLocals,
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"workspace_id": "00000000-0000-0000-0000-000000000000",
						"format":       "Default",
						"definition":   testHelperDefinition,
					},
				),
			),
			ExpectError: regexp.MustCompile(`The argument "display_name" is required, but no definition was found.`),
		},
		// error - no required attributes
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				testHelperLocals,
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"workspace_id": "00000000-0000-0000-0000-000000000000",
						"display_name": "test",
						"definition":   testHelperDefinition,
					},
				),
			),
			ExpectError: regexp.MustCompile(`The argument "format" is required, but no definition was found.`),
		},
		// error - no required attributes
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": "00000000-0000-0000-0000-000000000000",
					"display_name": "test",
					"format":       "Default",
				},
			),
			ExpectError: regexp.MustCompile(`The argument "definition" is required, but no definition was found.`),
		},
	}))
}

func TestUnit_MountedDataFactoryResource_ImportState(t *testing.T) {
	workspaceID := testhelp.RandomUUID()
	entity := fakes.NewRandomItemWithWorkspace(itemType, workspaceID)

	fakes.FakeServer.Upsert(fakes.NewRandomItemWithWorkspace(itemType, workspaceID))
	fakes.FakeServer.Upsert(entity)
	fakes.FakeServer.Upsert(fakes.NewRandomItemWithWorkspace(itemType, workspaceID))

	testCase := at.JoinConfigs(
		testHelperLocals,
		at.CompileConfig(
			testResourceItemHeader,
			map[string]any{
				"workspace_id": *entity.WorkspaceID,
				"display_name": *entity.DisplayName,
				"format":       "Default",
				"definition":   testHelperDefinition,
			},
		))

	resource.Test(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		{
			ResourceName:  testResourceItemFQN,
			Config:        testCase,
			ImportStateId: "not-valid",
			ImportState:   true,
			ExpectError:   regexp.MustCompile(fmt.Sprintf(common.ErrorImportIdentifierDetails, fmt.Sprintf("WorkspaceID/%sID", string(itemType)))),
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

func TestUnit_MountedDataFactoryResource_CRUD(t *testing.T) {
	workspaceID := testhelp.RandomUUID()
	entityExist := fakes.NewRandomItemWithWorkspace(itemType, workspaceID)
	entityBefore := fakes.NewRandomItemWithWorkspace(itemType, workspaceID)
	entityAfter := fakes.NewRandomItemWithWorkspace(itemType, workspaceID)

	fakes.FakeServer.Upsert(fakes.NewRandomItemWithWorkspace(itemType, workspaceID))
	fakes.FakeServer.Upsert(entityExist)
	fakes.FakeServer.Upsert(entityAfter)
	fakes.FakeServer.Upsert(fakes.NewRandomItemWithWorkspace(itemType, workspaceID))

	resource.Test(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		// error - create - existing entity
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				testHelperLocals,
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"workspace_id": *entityExist.WorkspaceID,
						"display_name": *entityExist.DisplayName,
						"format":       "Default",
						"definition":   testHelperDefinition,
					},
				),
			),
			ExpectError: regexp.MustCompile(common.ErrorCreateHeader),
		},
		// Create and Read - no definition
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				testHelperLocals,
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"workspace_id": *entityExist.WorkspaceID,
						"display_name": *entityExist.DisplayName,
						"format":       "Default",
					},
				),
			),
			ExpectError: regexp.MustCompile(`The argument "definition" is required, but no definition was found.`),
		},
		// Create and Read with definition
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				testHelperLocals,
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"workspace_id": *entityBefore.WorkspaceID,
						"display_name": *entityBefore.DisplayName,
						"description":  *entityBefore.Description,
						"format":       "Default",
						"definition":   testHelperDefinition,
					},
				)),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "display_name", entityBefore.DisplayName),
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "description", entityBefore.Description),
				resource.TestCheckResourceAttr(testResourceItemFQN, "definition_update_enabled", "true"),
			),
		},
		// Update and Read
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				testHelperLocals,
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"workspace_id": *entityBefore.WorkspaceID,
						"display_name": *entityAfter.DisplayName,
						"description":  *entityAfter.Description,
						"format":       "Default",
						"definition":   testHelperDefinition,
					},
				),
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "display_name", entityAfter.DisplayName),
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "description", entityAfter.Description),
				resource.TestCheckResourceAttr(testResourceItemFQN, "definition_update_enabled", "true"),
			),
		},
		// Delete testing automatically occurs in TestCase
	}))
}

func TestAcc_MountedDataFactoryResource_CRUD(t *testing.T) {
	workspace := testhelp.WellKnown()["WorkspaceRS"].(map[string]any)
	workspaceID := workspace["id"].(string)

	entityCreateDisplayName := testhelp.RandomName()
	entityUpdateDisplayName := testhelp.RandomName()
	entityUpdateDescription := testhelp.RandomName()

	resource.Test(t, testhelp.NewTestAccCase(t, &testResourceItemFQN, nil, []resource.TestStep{
		// Create and Read
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				testHelperLocals,
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"workspace_id": workspaceID,
						"display_name": entityCreateDisplayName,
						"format":       "Default",
						"definition":   testHelperDefinition,
					},
				),
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "display_name", entityCreateDisplayName),
				resource.TestCheckResourceAttr(testResourceItemFQN, "description", ""),
			),
		},
		// Update and Read
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				testHelperLocals,
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"workspace_id": workspaceID,
						"display_name": entityUpdateDisplayName,
						"description":  entityUpdateDescription,
						"format":       "Default",
						"definition":   testHelperDefinition,
					},
				),
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "display_name", entityUpdateDisplayName),
				resource.TestCheckResourceAttr(testResourceItemFQN, "description", entityUpdateDescription),
			),
		},
	},
	))
}
