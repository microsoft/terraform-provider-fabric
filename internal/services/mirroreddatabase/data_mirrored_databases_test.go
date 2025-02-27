// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package mirroreddatabase_test

import (
	"regexp"
	"testing"

	at "github.com/dcarbone/terraform-plugin-framework-utils/v3/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp/fakes"
)

var (
	testDataSourceMirroredDatabasesFQN    = testhelp.DataSourceFQN("fabric", "mirrored_databases", "test")
	testDataSourceMirroredDatabasesHeader = at.DataSourceHeader(testhelp.TypeName("fabric", "mirrored_databases"), "test")
)

func TestUnit_MirroredDatabasesDataSource(t *testing.T) {
	workspaceID := testhelp.RandomUUID()
	entity := fakes.NewRandomMirroredDatabaseWithWorkspace(workspaceID)

	fakes.FakeServer.Upsert(fakes.NewRandomMirroredDatabaseWithWorkspace(workspaceID))
	fakes.FakeServer.Upsert(entity)
	fakes.FakeServer.Upsert(fakes.NewRandomMirroredDatabaseWithWorkspace(workspaceID))

	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, nil, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		// error - no attributes
		{
			Config: at.CompileConfig(
				testDataSourceMirroredDatabasesHeader,
				map[string]any{},
			),
			ExpectError: regexp.MustCompile(`The argument "workspace_id" is required, but no definition was found`),
		},
		// error - workspace_id - invalid UUID
		{
			Config: at.CompileConfig(
				testDataSourceMirroredDatabasesHeader,
				map[string]any{
					"workspace_id": "invalid uuid",
				},
			),
			ExpectError: regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
		// error - unexpected_attr
		{
			Config: at.CompileConfig(
				testDataSourceMirroredDatabasesHeader,
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
				testDataSourceMirroredDatabasesHeader,
				map[string]any{
					"workspace_id": workspaceID,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrPtr(testDataSourceMirroredDatabasesFQN, "workspace_id", entity.WorkspaceID),
				resource.TestCheckResourceAttrPtr(testDataSourceMirroredDatabasesFQN, "values.1.id", entity.ID),
				resource.TestCheckResourceAttrSet(testDataSourceMirroredDatabasesFQN, "values.1.properties.default_schema"),
				resource.TestCheckResourceAttrSet(testDataSourceMirroredDatabasesFQN, "values.1.properties.onelake_tables_path"),
				resource.TestCheckResourceAttrSet(testDataSourceMirroredDatabasesFQN, "values.1.properties.sql_endpoint_properties.connection_string"),
			),
		},
	}))
}

func TestAcc_MirroredDatabasesDataSource(t *testing.T) {
	workspace := testhelp.WellKnown()["WorkspaceDS"].(map[string]any)
	workspaceID := workspace["id"].(string)

	resource.ParallelTest(t, testhelp.NewTestAccCase(t, nil, nil, []resource.TestStep{
		// read
		{
			Config: at.CompileConfig(
				testDataSourceMirroredDatabasesHeader,
				map[string]any{
					"workspace_id": workspaceID,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testDataSourceMirroredDatabasesFQN, "workspace_id", workspaceID),
				resource.TestCheckResourceAttrSet(testDataSourceMirroredDatabasesFQN, "values.0.id"),
				resource.TestCheckResourceAttrSet(testDataSourceMirroredDatabasesFQN, "values.0.properties.default_schema"),
				resource.TestCheckResourceAttrSet(testDataSourceMirroredDatabasesFQN, "values.0.properties.onelake_tables_path"),
			),
		},
	}))
}
