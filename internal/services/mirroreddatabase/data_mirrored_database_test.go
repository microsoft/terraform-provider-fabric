// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package mirroreddatabase_test

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

var (
	testDataSourceMirroredDatabaseFQN    = testhelp.DataSourceFQN("fabric", itemTFName, "test")
	testDataSourceMirroredDatabaseHeader = at.DataSourceHeader(testhelp.TypeName("fabric", itemTFName), "test")
)

func TestUnit_MirroredDatabaseDataSource(t *testing.T) {
	workspaceID := testhelp.RandomUUID()
	entity := fakes.NewRandomItemWithWorkspace(itemTFName, workspaceID)

	fakes.FakeServer.Upsert(fakes.NewRandomItemWithWorkspace(itemTFName, workspaceID))
	fakes.FakeServer.Upsert(entity)
	fakes.FakeServer.Upsert(fakes.NewRandomItemWithWorkspace(itemTFName, workspaceID))

	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, nil, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		// error - no attributes provided
		{
			Config: at.CompileConfig(
				testDataSourceMirroredDatabaseHeader,
				map[string]any{},
			),
			ExpectError: regexp.MustCompile(`The argument "workspace_id" is required, but no definition was found`),
		},
		// error - invalid workspace_id UUID
		{
			Config: at.CompileConfig(
				testDataSourceMirroredDatabaseHeader,
				map[string]any{
					"workspace_id": "invalid uuid",
				},
			),
			ExpectError: regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
		// error - unexpected attribute
		{
			Config: at.CompileConfig(
				testDataSourceMirroredDatabaseHeader,
				map[string]any{
					"workspace_id":    workspaceID,
					"unexpected_attr": "test",
				},
			),
			ExpectError: regexp.MustCompile(`An argument named "unexpected_attr" is not expected here`),
		},
		// error - conflicting attributes provided together
		{
			Config: at.CompileConfig(
				testDataSourceMirroredDatabaseHeader,
				map[string]any{
					"workspace_id": workspaceID,
					"id":           *entity.ID,
					"display_name": *entity.DisplayName,
				},
			),
			ExpectError: regexp.MustCompile(`These attributes cannot be configured together: \[id,display_name\]`),
		},
		// error - missing one of required attributes (neither id nor display_name provided)
		{
			Config: at.CompileConfig(
				testDataSourceMirroredDatabaseHeader,
				map[string]any{
					"workspace_id": workspaceID,
				},
			),
			ExpectError: regexp.MustCompile(`Exactly one of these attributes must be configured: \[id,display_name\]`),
		},
		// error - id provided without workspace_id
		{
			Config: at.CompileConfig(
				testDataSourceMirroredDatabaseHeader,
				map[string]any{
					"id": *entity.ID,
				},
			),
			ExpectError: regexp.MustCompile(`The argument "workspace_id" is required, but no definition was found`),
		},
		// read by id - success
		{
			Config: at.CompileConfig(
				testDataSourceMirroredDatabaseHeader,
				map[string]any{
					"workspace_id": workspaceID,
					"id":           *entity.ID,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrPtr(testDataSourceMirroredDatabaseFQN, "workspace_id", entity.WorkspaceID),
				resource.TestCheckResourceAttrPtr(testDataSourceMirroredDatabaseFQN, "id", entity.ID),
				resource.TestCheckResourceAttrPtr(testDataSourceMirroredDatabaseFQN, "display_name", entity.DisplayName),
				resource.TestCheckResourceAttrPtr(testDataSourceMirroredDatabaseFQN, "description", entity.Description),
			),
		},
		// read by id - not found
		{
			Config: at.CompileConfig(
				testDataSourceMirroredDatabaseHeader,
				map[string]any{
					"workspace_id": workspaceID,
					"id":           testhelp.RandomUUID(),
				},
			),
			ExpectError: regexp.MustCompile(common.ErrorReadHeader),
		},
		// read by display name - success
		{
			Config: at.CompileConfig(
				testDataSourceMirroredDatabaseHeader,
				map[string]any{
					"workspace_id": workspaceID,
					"display_name": *entity.DisplayName,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrPtr(testDataSourceMirroredDatabaseFQN, "workspace_id", entity.WorkspaceID),
				resource.TestCheckResourceAttrPtr(testDataSourceMirroredDatabaseFQN, "id", entity.ID),
				resource.TestCheckResourceAttrPtr(testDataSourceMirroredDatabaseFQN, "display_name", entity.DisplayName),
				resource.TestCheckResourceAttrPtr(testDataSourceMirroredDatabaseFQN, "description", entity.Description),
			),
		},
		// read by display name - not found
		{
			Config: at.CompileConfig(
				testDataSourceMirroredDatabaseHeader,
				map[string]any{
					"workspace_id": workspaceID,
					"display_name": testhelp.RandomName(),
				},
			),
			ExpectError: regexp.MustCompile(common.ErrorReadHeader),
		},
		// read by id with definition - missing required format configuration
		{
			Config: at.CompileConfig(
				testDataSourceMirroredDatabaseHeader,
				map[string]any{
					"workspace_id":      workspaceID,
					"id":                *entity.ID,
					"output_definition": true,
				},
			),
			ExpectError: regexp.MustCompile("Invalid configuration for attribute format"),
		},
		// read by id with definition - success
		{
			Config: at.CompileConfig(
				testDataSourceMirroredDatabaseHeader,
				map[string]any{
					"workspace_id":      workspaceID,
					"id":                *entity.ID,
					"output_definition": true,
					"format":            "Default",
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrPtr(testDataSourceMirroredDatabaseFQN, "workspace_id", entity.WorkspaceID),
				resource.TestCheckResourceAttrPtr(testDataSourceMirroredDatabaseFQN, "id", entity.ID),
				resource.TestCheckResourceAttrPtr(testDataSourceMirroredDatabaseFQN, "display_name", entity.DisplayName),
				resource.TestCheckResourceAttrPtr(testDataSourceMirroredDatabaseFQN, "description", entity.Description),
				resource.TestCheckResourceAttrSet(testDataSourceMirroredDatabaseFQN, "definition.RealTimeDefinition.json.content"),
			),
		},
	}))
}

func TestAcc_MirroredDatabaseDataSource(t *testing.T) {
	workspace := testhelp.WellKnown()["WorkspaceDS"].(map[string]any)
	workspaceID := workspace["id"].(string)

	entity := testhelp.WellKnown()["MirroredDatabase"].(map[string]any)
	entityID := entity["id"].(string)
	entityDisplayName := entity["displayName"].(string)
	entityDescription := entity["description"].(string)

	resource.ParallelTest(t, testhelp.NewTestAccCase(t, nil, nil, []resource.TestStep{
		// read by id - success
		{
			Config: at.CompileConfig(
				testDataSourceMirroredDatabaseHeader,
				map[string]any{
					"workspace_id": workspaceID,
					"id":           entityID,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testDataSourceMirroredDatabaseFQN, "workspace_id", workspaceID),
				resource.TestCheckResourceAttr(testDataSourceMirroredDatabaseFQN, "id", entityID),
				resource.TestCheckResourceAttr(testDataSourceMirroredDatabaseFQN, "display_name", entityDisplayName),
				resource.TestCheckResourceAttr(testDataSourceMirroredDatabaseFQN, "description", entityDescription),
			),
		},
		// read by id - not found
		{
			Config: at.CompileConfig(
				testDataSourceMirroredDatabaseHeader,
				map[string]any{
					"workspace_id": workspaceID,
					"id":           testhelp.RandomUUID(),
				},
			),
			ExpectError: regexp.MustCompile(common.ErrorReadHeader),
		},
		// read by display name - success
		{
			Config: at.CompileConfig(
				testDataSourceMirroredDatabaseHeader,
				map[string]any{
					"workspace_id": workspaceID,
					"display_name": entityDisplayName,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testDataSourceMirroredDatabaseFQN, "workspace_id", workspaceID),
				resource.TestCheckResourceAttr(testDataSourceMirroredDatabaseFQN, "id", entityID),
				resource.TestCheckResourceAttr(testDataSourceMirroredDatabaseFQN, "display_name", entityDisplayName),
				resource.TestCheckResourceAttr(testDataSourceMirroredDatabaseFQN, "description", entityDescription),
				resource.TestCheckNoResourceAttr(testDataSourceMirroredDatabaseFQN, "definition"),
			),
		},
		// read by display name - not found
		{
			Config: at.CompileConfig(
				testDataSourceMirroredDatabaseHeader,
				map[string]any{
					"workspace_id": workspaceID,
					"display_name": testhelp.RandomName(),
				},
			),
			ExpectError: regexp.MustCompile(common.ErrorReadHeader),
		},
		// read by id with definition - missing format error
		{
			Config: at.CompileConfig(
				testDataSourceMirroredDatabaseHeader,
				map[string]any{
					"workspace_id":      workspaceID,
					"id":                entityID,
					"output_definition": true,
				},
			),
			ExpectError: regexp.MustCompile("Invalid configuration for attribute format"),
		},
		// read by id with definition - success
		{
			Config: at.CompileConfig(
				testDataSourceMirroredDatabaseHeader,
				map[string]any{
					"workspace_id":      workspaceID,
					"id":                entityID,
					"output_definition": true,
					"format":            "Default",
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testDataSourceMirroredDatabaseFQN, "workspace_id", workspaceID),
				resource.TestCheckResourceAttr(testDataSourceMirroredDatabaseFQN, "id", entityID),
				resource.TestCheckResourceAttr(testDataSourceMirroredDatabaseFQN, "display_name", entityDisplayName),
				resource.TestCheckResourceAttr(testDataSourceMirroredDatabaseFQN, "description", entityDescription),
				resource.TestCheckResourceAttrSet(testDataSourceMirroredDatabaseFQN, "definition.RealTimeDefinition.json.content"),
			),
		},
	}))
}
