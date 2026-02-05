// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package mirroredazuredatabrickscatalog_test

import (
	"regexp"
	"testing"

	at "github.com/dcarbone/terraform-plugin-framework-utils/v3/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp/fakes"
)

var testDataSourceItemsFQN, testDataSourceItemsHeader = testhelp.TFDataSource(common.ProviderTypeName, itemTypeInfo.Types, "test")

func TestUnit_MirroredAzureDatabricksCatalogsDataSource(t *testing.T) {
	workspaceID := testhelp.RandomUUID()
	entity := fakes.NewRandomMirroredAzureDatabricksCatalogWithWorkspace(workspaceID)

	fakes.FakeServer.Upsert(fakes.NewRandomMirroredAzureDatabricksCatalogWithWorkspace(workspaceID))
	fakes.FakeServer.Upsert(entity)
	fakes.FakeServer.Upsert(fakes.NewRandomMirroredAzureDatabricksCatalogWithWorkspace(workspaceID))

	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, nil, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		// error - no attributes
		{
			Config: at.CompileConfig(
				testDataSourceItemsHeader,
				map[string]any{},
			),
			ExpectError: regexp.MustCompile(`The argument "workspace_id" is required, but no definition was found`),
		},
		// error - workspace_id - invalid UUID
		{
			Config: at.CompileConfig(
				testDataSourceItemsHeader,
				map[string]any{
					"workspace_id": "invalid uuid",
				},
			),
			ExpectError: regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
		// error - unexpected_attr
		{
			Config: at.CompileConfig(
				testDataSourceItemsHeader,
				map[string]any{
					"workspace_id":    workspaceID,
					"unexpected_attr": "test",
				},
			),
			ExpectError: regexp.MustCompile(`An argument named "unexpected_attr" is not expected here`),
		},
		// read
		{
			Config: at.CompileConfig(
				testDataSourceItemsHeader,
				map[string]any{
					"workspace_id": workspaceID,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrPtr(testDataSourceItemsFQN, "workspace_id", entity.WorkspaceID),
				resource.TestCheckResourceAttrSet(testDataSourceItemsFQN, "values.1.properties.auto_sync"),
				resource.TestCheckResourceAttrSet(testDataSourceItemsFQN, "values.1.properties.catalog_name"),
				resource.TestCheckResourceAttrSet(testDataSourceItemsFQN, "values.1.properties.databricks_workspace_connection_id"),
				resource.TestCheckResourceAttrSet(testDataSourceItemsFQN, "values.1.properties.mirror_status"),
				resource.TestCheckResourceAttrSet(testDataSourceItemsFQN, "values.1.properties.mirroring_mode"),
				resource.TestCheckResourceAttrSet(testDataSourceItemsFQN, "values.1.properties.onelake_tables_path"),
				resource.TestCheckResourceAttrSet(testDataSourceItemsFQN, "values.1.properties.storage_connection_id"),
				resource.TestCheckResourceAttrSet(testDataSourceItemsFQN, "values.1.properties.sync_details.status"),
				resource.TestCheckResourceAttrSet(testDataSourceItemsFQN, "values.1.properties.sync_details.last_sync_date_time"),
				resource.TestCheckResourceAttrSet(testDataSourceItemsFQN, "values.1.properties.sql_endpoint_properties.connection_string"),
			),
			ConfigStateChecks: []statecheck.StateCheck{
				statecheck.ExpectKnownValue(
					testDataSourceItemsFQN,
					tfjsonpath.New("values"),
					knownvalue.SetPartial([]knownvalue.Check{
						knownvalue.ObjectPartial(map[string]knownvalue.Check{
							"id":           knownvalue.StringExact(*entity.ID),
							"display_name": knownvalue.StringExact(*entity.DisplayName),
							"description":  knownvalue.StringExact(*entity.Description),
						}),
					}),
				),
			},
		},
	}))
}

func TestAcc_MirroredAzureDatabricksCatalogsDataSource(t *testing.T) {
	workspace := testhelp.WellKnown()["WorkspaceDS"].(map[string]any)
	workspaceID := workspace["id"].(string)

	resource.ParallelTest(t, testhelp.NewTestAccCase(t, nil, nil, []resource.TestStep{
		// read
		{
			Config: at.CompileConfig(
				testDataSourceItemsHeader,
				map[string]any{
					"workspace_id": workspaceID,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testDataSourceItemsFQN, "workspace_id", workspaceID),
				resource.TestCheckResourceAttrSet(testDataSourceItemsFQN, "values.0.id"),
				resource.TestCheckResourceAttrSet(testDataSourceItemsFQN, "values.0.display_name"),
				resource.TestCheckResourceAttrSet(testDataSourceItemsFQN, "values.0.properties.auto_sync"),
				resource.TestCheckResourceAttrSet(testDataSourceItemsFQN, "values.0.properties.databricks_workspace_connection_id"),
				resource.TestCheckResourceAttrSet(testDataSourceItemsFQN, "values.0.properties.mirror_status"),
				resource.TestCheckResourceAttrSet(testDataSourceItemsFQN, "values.0.properties.mirroring_mode"),
				resource.TestCheckResourceAttrSet(testDataSourceItemsFQN, "values.0.properties.onelake_tables_path"),
				resource.TestCheckResourceAttrSet(testDataSourceItemsFQN, "values.0.properties.sql_endpoint_properties.connection_string"),
			),
		},
	},
	))
}
