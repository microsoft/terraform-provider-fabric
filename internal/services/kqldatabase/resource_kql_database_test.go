// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package kqldatabase_test

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

func TestUnit_KQLDatabaseResource_Attributes(t *testing.T) {
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
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": "invalid uuid",
					"display_name": "test",
					"configuration": map[string]any{
						"database_type": "ReadWrite",
						"eventhouse_id": "00000000-0000-0000-0000-000000000000",
					},
				},
			),
			ExpectError: regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
		// error - unexpected attribute
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id":    "00000000-0000-0000-0000-000000000000",
					"unexpected_attr": "test",
				},
			),
			ExpectError: regexp.MustCompile(`An argument named "unexpected_attr" is not expected here`),
		},
		// error - no required attributes
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"display_name": "test",
				},
			),
			ExpectError: regexp.MustCompile(`The argument "workspace_id" is required, but no definition was found.`),
		},
		// error - no required attributes
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": "00000000-0000-0000-0000-000000000000",
				},
			),
			ExpectError: regexp.MustCompile(`The argument "display_name" is required, but no definition was found.`),
		},
		// error - no required attributes (configuration)
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id":  "00000000-0000-0000-0000-000000000000",
					"display_name":  "test",
					"configuration": map[string]any{},
				},
			),
			ExpectError: regexp.MustCompile(`Inappropriate value for attribute "configuration".`),
		},
		// error - database_type - invalid
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": "00000000-0000-0000-0000-000000000000",
					"display_name": "test",
					"configuration": map[string]any{
						"database_type": "Test",
						"eventhouse_id": "00000000-0000-0000-0000-000000000000",
					},
				},
			),
			ExpectError: regexp.MustCompile(common.ErrorAttValueMatch),
		},
		// error - source_database_name - unexpected attribute
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": "00000000-0000-0000-0000-000000000000",
					"display_name": "test",
					"configuration": map[string]any{
						"database_type":        "ReadWrite",
						"eventhouse_id":        "00000000-0000-0000-0000-000000000000",
						"source_database_name": "test",
					},
				},
			),
			ExpectError: regexp.MustCompile("Invalid configuration for attribute"),
		},
		// error - source_cluster_uri - unexpected attribute
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": "00000000-0000-0000-0000-000000000000",
					"display_name": "test",
					"configuration": map[string]any{
						"database_type":      "ReadWrite",
						"eventhouse_id":      "00000000-0000-0000-0000-000000000000",
						"source_cluster_uri": "https://test.westus.kusto.windows.net",
					},
				},
			),
			ExpectError: regexp.MustCompile("Invalid configuration for attribute"),
		},
		// error - invitation_token - unexpected attribute
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": "00000000-0000-0000-0000-000000000000",
					"display_name": "test",
					"configuration": map[string]any{
						"database_type":    "ReadWrite",
						"eventhouse_id":    "00000000-0000-0000-0000-000000000000",
						"invitation_token": "test",
					},
				},
			),
			ExpectError: regexp.MustCompile("Invalid configuration for attribute"),
		},
		// error - invitation_token/source_database_name - invalid combination
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": "00000000-0000-0000-0000-000000000000",
					"display_name": "test",
					"configuration": map[string]any{
						"database_type":        "Shortcut",
						"eventhouse_id":        "00000000-0000-0000-0000-000000000000",
						"invitation_token":     "test",
						"source_database_name": "test",
					},
				},
			),
			ExpectError: regexp.MustCompile("Invalid Attribute Combination"),
		},
		// error - invitation_token/source_cluster_uri - invalid combination
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": "00000000-0000-0000-0000-000000000000",
					"display_name": "test",
					"configuration": map[string]any{
						"database_type":      "Shortcut",
						"eventhouse_id":      "00000000-0000-0000-0000-000000000000",
						"invitation_token":   "test",
						"source_cluster_uri": "https://test.westus.kusto.windows.net",
					},
				},
			),
			ExpectError: regexp.MustCompile("Invalid Attribute Combination"),
		},
		// error - source_cluster_uri - invalid URL
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": "00000000-0000-0000-0000-000000000000",
					"display_name": "test",
					"configuration": map[string]any{
						"database_type":        "Shortcut",
						"eventhouse_id":        "00000000-0000-0000-0000-000000000000",
						"source_database_name": "test",
						"source_cluster_uri":   "test.westus.kusto.windows.net",
					},
				},
			),
			ExpectError: regexp.MustCompile(customtypes.URLTypeErrorInvalidStringHeader),
		},
	}))
}

func TestUnit_KQLDatabaseResource_ImportState(t *testing.T) {
	workspaceID := testhelp.RandomUUID()
	entity := fakes.NewRandomKQLDatabaseWithWorkspace(workspaceID)

	fakes.FakeServer.Upsert(fakes.NewRandomKQLDatabaseWithWorkspace(workspaceID))
	fakes.FakeServer.Upsert(entity)
	fakes.FakeServer.Upsert(fakes.NewRandomKQLDatabaseWithWorkspace(workspaceID))

	testCase := at.CompileConfig(
		testResourceItemHeader,
		map[string]any{
			"workspace_id": *entity.WorkspaceID,
			"display_name": *entity.DisplayName,
			"configuration": map[string]any{
				"database_type": "ReadWrite",
				"eventhouse_id": "00000000-0000-0000-0000-000000000000",
			},
		},
	)

	resource.Test(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		{
			ResourceName:  testResourceItemFQN,
			Config:        testCase,
			ImportStateId: "not-valid",
			ImportState:   true,
			ExpectError:   regexp.MustCompile(fmt.Sprintf(common.ErrorImportIdentifierDetails, "WorkspaceID/KQLDatabaseID")),
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

func TestUnit_KQLDatabaseResource_CRUD(t *testing.T) {
	workspaceID := testhelp.RandomUUID()
	eventhouse := fakes.NewRandomEventhouseWithWorkspace(workspaceID)
	fakes.FakeServer.Upsert(eventhouse)

	entityExist := fakes.NewRandomKQLDatabaseWithWorkspace(workspaceID)
	entityBefore := fakes.NewRandomKQLDatabaseWithWorkspace(workspaceID)
	entityAfter := fakes.NewRandomKQLDatabaseWithWorkspace(workspaceID)

	fakes.FakeServer.Upsert(fakes.NewRandomKQLDatabaseWithWorkspace(workspaceID))
	fakes.FakeServer.Upsert(entityExist)
	fakes.FakeServer.Upsert(entityAfter)
	fakes.FakeServer.Upsert(fakes.NewRandomKQLDatabaseWithWorkspace(workspaceID))

	resource.Test(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		// error - create - existing entity
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": *entityExist.WorkspaceID,
					"display_name": *entityExist.DisplayName,
					"configuration": map[string]any{
						"database_type": "ReadWrite",
						"eventhouse_id": *entityExist.Properties.ParentEventhouseItemID,
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
					"workspace_id": *entityBefore.WorkspaceID,
					"display_name": *entityBefore.DisplayName,
					"configuration": map[string]any{
						"database_type": "ReadWrite",
						"eventhouse_id": *eventhouse.ID,
					},
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "display_name", entityBefore.DisplayName),
				resource.TestCheckResourceAttr(testResourceItemFQN, "description", ""),
			),
		},
		// Update and Read
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": *entityBefore.WorkspaceID,
					"display_name": *entityAfter.DisplayName,
					"description":  *entityAfter.Description,
					"configuration": map[string]any{
						"database_type": "ReadWrite",
						"eventhouse_id": *eventhouse.ID,
					},
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "display_name", entityAfter.DisplayName),
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "description", entityAfter.Description),
			),
		},
		// Delete testing automatically occurs in TestCase
	}))
}

func TestAcc_KQLDatabaseResource_CRUD(t *testing.T) {
	workspace := testhelp.WellKnown()["Workspace"].(map[string]any)
	workspaceID := workspace["id"].(string)

	eventhouseResourceHCL, eventhouseResourceFQN := eventhouseResource(t, workspaceID)

	entityCreateDisplayName := testhelp.RandomName()
	entityUpdateDisplayName := testhelp.RandomName()
	entityUpdateDescription := testhelp.RandomName()

	resource.Test(t, testhelp.NewTestAccCase(t, &testResourceItemFQN, nil, []resource.TestStep{
		// Create and Read (Configuration)
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				eventhouseResourceHCL,
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"workspace_id": workspaceID,
						"display_name": entityCreateDisplayName,
						"configuration": map[string]any{
							"database_type": "ReadWrite",
							"eventhouse_id": testhelp.RefByFQN(eventhouseResourceFQN, "id"),
						},
					},
				),
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "display_name", entityCreateDisplayName),
				resource.TestCheckResourceAttr(testResourceItemFQN, "description", ""),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.eventhouse_id"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.query_service_uri"),
			),
		},
		// Update and Read (Configuration)
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				eventhouseResourceHCL,
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"workspace_id": workspaceID,
						"display_name": entityUpdateDisplayName,
						"description":  entityUpdateDescription,
						"configuration": map[string]any{
							"database_type": "ReadWrite",
							"eventhouse_id": testhelp.RefByFQN(eventhouseResourceFQN, "id"),
						},
					},
				),
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "display_name", entityUpdateDisplayName),
				resource.TestCheckResourceAttr(testResourceItemFQN, "description", entityUpdateDescription),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.eventhouse_id"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.query_service_uri"),
			),
		},
	},
	))
}
