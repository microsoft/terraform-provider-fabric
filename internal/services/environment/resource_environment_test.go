// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package environment_test

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

func TestUnit_EnvironmentResource_Attributes(t *testing.T) {
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
	}))
}

func TestUnit_EnvironmentResource_ImportState(t *testing.T) {
	workspaceID := testhelp.RandomUUID()
	entity := fakes.NewRandomEnvironmentWithWorkspace(workspaceID)

	fakes.FakeServer.Upsert(fakes.NewRandomEnvironmentWithWorkspace(workspaceID))
	fakes.FakeServer.Upsert(entity)
	fakes.FakeServer.Upsert(fakes.NewRandomEnvironmentWithWorkspace(workspaceID))

	testCase := at.CompileConfig(
		testResourceItemHeader,
		map[string]any{
			"workspace_id": *entity.WorkspaceID,
			"display_name": *entity.DisplayName,
		},
	)

	resource.Test(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		{
			ResourceName:  testResourceItemFQN,
			Config:        testCase,
			ImportStateId: "not-valid",
			ImportState:   true,
			ExpectError:   regexp.MustCompile(fmt.Sprintf(common.ErrorImportIdentifierDetails, "WorkspaceID/EnvironmentID")),
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

func TestUnit_EnvironmentResource_CRUD(t *testing.T) {
	workspaceID := testhelp.RandomUUID()
	entityExist := fakes.NewRandomEnvironmentWithWorkspace(workspaceID)
	entityBefore := fakes.NewRandomEnvironmentWithWorkspace(workspaceID)
	entityAfter := fakes.NewRandomEnvironmentWithWorkspace(workspaceID)

	fakes.FakeServer.Upsert(fakes.NewRandomEnvironmentWithWorkspace(workspaceID))
	fakes.FakeServer.Upsert(entityExist)
	fakes.FakeServer.Upsert(entityAfter)
	fakes.FakeServer.Upsert(fakes.NewRandomEnvironmentWithWorkspace(workspaceID))

	resource.Test(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		// error - create - existing entity
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": *entityExist.WorkspaceID,
					"display_name": *entityExist.DisplayName,
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
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "display_name", entityBefore.DisplayName),
				resource.TestCheckResourceAttr(testResourceItemFQN, "description", ""),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.publish_details.state"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.publish_details.component_publish_info.spark_libraries.state"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.publish_details.component_publish_info.spark_settings.state"),
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
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "display_name", entityAfter.DisplayName),
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "description", entityAfter.Description),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.publish_details.state"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.publish_details.component_publish_info.spark_libraries.state"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.publish_details.component_publish_info.spark_settings.state"),
			),
		},
		// Delete testing automatically occurs in TestCase
	}))
}

func TestAcc_EnvironmentResource_CRUD(t *testing.T) {
	workspace := testhelp.WellKnown()["WorkspaceRS"].(map[string]any)
	workspaceID := workspace["id"].(string)

	entityCreateDisplayName := testhelp.RandomName()
	entityUpdateDisplayName := testhelp.RandomName()
	entityUpdateDescription := testhelp.RandomName()

	resource.Test(t, testhelp.NewTestAccCase(t, &testResourceItemFQN, nil, []resource.TestStep{
		// Create and Read
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": workspaceID,
					"display_name": entityCreateDisplayName,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "display_name", entityCreateDisplayName),
				resource.TestCheckResourceAttr(testResourceItemFQN, "description", ""),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.publish_details.state"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.publish_details.component_publish_info.spark_libraries.state"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.publish_details.component_publish_info.spark_settings.state"),
			),
		},
		// Update and Read
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": workspaceID,
					"display_name": entityUpdateDisplayName,
					"description":  entityUpdateDescription,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "display_name", entityUpdateDisplayName),
				resource.TestCheckResourceAttr(testResourceItemFQN, "description", entityUpdateDescription),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.publish_details.state"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.publish_details.component_publish_info.spark_libraries.state"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.publish_details.component_publish_info.spark_settings.state"),
			),
		},
	},
	))
}
