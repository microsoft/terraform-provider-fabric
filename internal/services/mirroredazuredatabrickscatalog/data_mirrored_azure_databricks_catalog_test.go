// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package mirroredazuredatabrickscatalog_test

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

var testDataSourceItemFQN, testDataSourceItemHeader = testhelp.TFDataSource(common.ProviderTypeName, itemTypeInfo.Type, "test")

func TestUnit_MirroredAzureDatabricksCatalogDataSource(t *testing.T) {
	workspaceID := testhelp.RandomUUID()
	entity := fakes.NewRandomMirroredAzureDatabricksCatalogWithWorkspace(workspaceID)

	fakes.FakeServer.Upsert(fakes.NewRandomMirroredAzureDatabricksCatalogWithWorkspace(workspaceID))
	fakes.FakeServer.Upsert(entity)
	fakes.FakeServer.Upsert(fakes.NewRandomMirroredAzureDatabricksCatalogWithWorkspace(workspaceID))

	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, nil, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		// error - no attributes
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{},
			),
			ExpectError: regexp.MustCompile(`The argument "workspace_id" is required, but no definition was found`),
		},
		// error - workspace_id - invalid UUID
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"workspace_id": "invalid uuid",
				},
			),
			ExpectError: regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
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
		// error - conflicting attributes
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"workspace_id": workspaceID,
					"id":           *entity.ID,
					"display_name": *entity.DisplayName,
				},
			),
			ExpectError: regexp.MustCompile(`These attributes cannot be configured together: \[id,display_name\]`),
		},
		// error - no required attributes
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"workspace_id": workspaceID,
				},
			),
			ExpectError: regexp.MustCompile(`Exactly one of these attributes must be configured: \[id,display_name\]`),
		},
		// error - no required attributes
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"id": *entity.ID,
				},
			),
			ExpectError: regexp.MustCompile(`The argument "workspace_id" is required, but no definition was found.`),
		},
		// read by id
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"workspace_id": workspaceID,
					"id":           *entity.ID,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrPtr(testDataSourceItemFQN, "workspace_id", entity.WorkspaceID),
				resource.TestCheckResourceAttrPtr(testDataSourceItemFQN, "id", entity.ID),
				resource.TestCheckResourceAttrPtr(testDataSourceItemFQN, "display_name", entity.DisplayName),
				resource.TestCheckResourceAttrPtr(testDataSourceItemFQN, "description", entity.Description),
			),
		},
		// read by id - not found
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"workspace_id": workspaceID,
					"id":           testhelp.RandomUUID(),
				},
			),
			ExpectError: regexp.MustCompile(common.ErrorReadHeader),
		},

		// read by name
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"workspace_id": workspaceID,
					"display_name": *entity.DisplayName,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrPtr(testDataSourceItemFQN, "workspace_id", entity.WorkspaceID),
				resource.TestCheckResourceAttrPtr(testDataSourceItemFQN, "id", entity.ID),
				resource.TestCheckResourceAttrPtr(testDataSourceItemFQN, "display_name", entity.DisplayName),
				resource.TestCheckResourceAttrPtr(testDataSourceItemFQN, "description", entity.Description),
				resource.TestCheckResourceAttrSet(testDataSourceItemFQN, "properties.auto_sync"),
				resource.TestCheckResourceAttrSet(testDataSourceItemFQN, "properties.catalog_name"),
				resource.TestCheckResourceAttrSet(testDataSourceItemFQN, "properties.databricks_workspace_connection_id"),
				resource.TestCheckResourceAttrSet(testDataSourceItemFQN, "properties.mirror_status"),
				resource.TestCheckResourceAttrSet(testDataSourceItemFQN, "properties.mirroring_mode"),
				resource.TestCheckResourceAttrSet(testDataSourceItemFQN, "properties.onelake_tables_path"),
				resource.TestCheckResourceAttrSet(testDataSourceItemFQN, "properties.storage_connection_id"),
				resource.TestCheckResourceAttrSet(testDataSourceItemFQN, "properties.sync_details.status"),
				resource.TestCheckResourceAttrSet(testDataSourceItemFQN, "properties.sync_details.last_sync_date_time"),
				resource.TestCheckResourceAttrSet(testDataSourceItemFQN, "properties.sql_endpoint_properties.connection_string"),
			),
		},
		// read by name - not found
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"workspace_id": workspaceID,
					"display_name": testhelp.RandomName(),
				},
			),
			ExpectError: regexp.MustCompile(common.ErrorReadHeader),
		},
	}))
}

func TestAcc_MirroredAzureDatabricksCatalogDataSource(t *testing.T) {
	workspace := testhelp.WellKnown()["WorkspaceDS"].(map[string]any)
	workspaceID := workspace["id"].(string)

	entity := testhelp.WellKnown()["MirroredAzureDatabricksCatalog"].(map[string]any)
	entityID := entity["id"].(string)
	entityDisplayName := entity["displayName"].(string)
	entityDescription := entity["description"].(string)

	resource.ParallelTest(t, testhelp.NewTestAccCase(t, &testDataSourceItemFQN, nil, []resource.TestStep{
		// read by id
		{
			ResourceName: testDataSourceItemFQN,
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"workspace_id": workspaceID,
					"id":           entityID,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "workspace_id", workspaceID),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "id", entityID),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "display_name", entityDisplayName),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "description", entityDescription),
				resource.TestCheckResourceAttrSet(testDataSourceItemFQN, "properties.auto_sync"),
				resource.TestCheckNoResourceAttr(testDataSourceItemFQN, "properties.catalog_name"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "properties.databricks_workspace_connection_id", "00000000-0000-0000-0000-000000000000"),
				resource.TestCheckResourceAttrSet(testDataSourceItemFQN, "properties.mirror_status"),
				resource.TestCheckResourceAttrSet(testDataSourceItemFQN, "properties.mirroring_mode"),
				resource.TestCheckResourceAttrSet(testDataSourceItemFQN, "properties.onelake_tables_path"),
				resource.TestCheckNoResourceAttr(testDataSourceItemFQN, "properties.sync_details"),
				resource.TestCheckNoResourceAttr(testResourceItemFQN, "properties.storage_connection_id"),
				resource.TestCheckResourceAttrSet(testDataSourceItemFQN, "properties.sql_endpoint_properties.connection_string"),
			),
		},
		// read by id - not found
		{
			ResourceName: testDataSourceItemFQN,
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"workspace_id": workspaceID,
					"id":           testhelp.RandomUUID(),
				},
			),
			ExpectError: regexp.MustCompile(common.ErrorReadHeader),
		},
		// read by name
		{
			ResourceName: testDataSourceItemFQN,
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"workspace_id": workspaceID,
					"display_name": entityDisplayName,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "workspace_id", workspaceID),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "id", entityID),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "display_name", entityDisplayName),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "description", entityDescription),
				resource.TestCheckResourceAttrSet(testDataSourceItemFQN, "properties.auto_sync"),
				resource.TestCheckNoResourceAttr(testDataSourceItemFQN, "properties.catalog_name"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "properties.databricks_workspace_connection_id", "00000000-0000-0000-0000-000000000000"),
				resource.TestCheckResourceAttrSet(testDataSourceItemFQN, "properties.mirror_status"),
				resource.TestCheckResourceAttrSet(testDataSourceItemFQN, "properties.mirroring_mode"),
				resource.TestCheckResourceAttrSet(testDataSourceItemFQN, "properties.onelake_tables_path"),
				resource.TestCheckNoResourceAttr(testDataSourceItemFQN, "properties.sync_details"),
				resource.TestCheckNoResourceAttr(testResourceItemFQN, "properties.storage_connection_id"),
				resource.TestCheckResourceAttrSet(testDataSourceItemFQN, "properties.sql_endpoint_properties.connection_string"),
			),
		},
		// read by id with definition - missing format error
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"workspace_id":      workspaceID,
					"id":                entityID,
					"output_definition": true,
				},
			),
			ExpectError: regexp.MustCompile("Invalid configuration for attribute format"),
		},
		// read by name - not found
		{
			ResourceName: testDataSourceItemFQN,
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"workspace_id": workspaceID,
					"display_name": testhelp.RandomName(),
				},
			),
			ExpectError: regexp.MustCompile(common.ErrorReadHeader),
		},
		//nolint:godox,godoclint
		// TODO: API issue - getDefinition returns 404 for item created without definition
		// // read by id with definition
		// {
		// 	ResourceName: testDataSourceItemFQN,
		// 	Config: at.CompileConfig(
		// 		testDataSourceItemHeader,
		// 		map[string]any{
		// 			"workspace_id":      workspaceID,
		// 			"id":                entityID,
		// 			"format":            "Default",
		// 			"output_definition": true,
		// 		},
		// 	),
		// 	Check: resource.ComposeAggregateTestCheckFunc(
		// 		resource.TestCheckResourceAttr(testDataSourceItemFQN, "workspace_id", workspaceID),
		// 		resource.TestCheckResourceAttr(testDataSourceItemFQN, "id", entityID),
		// 		resource.TestCheckResourceAttr(testDataSourceItemFQN, "display_name", entityDisplayName),
		// 		resource.TestCheckResourceAttr(testDataSourceItemFQN, "description", entityDescription),
		// 		resource.TestCheckResourceAttrSet(testDataSourceItemFQN, "definition.mirroringAzureDatabricksCatalog.json.content"),
		// 	),
		// },
	}))
}
