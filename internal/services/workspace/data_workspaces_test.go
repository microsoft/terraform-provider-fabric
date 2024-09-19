// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package workspace_test

import (
	"regexp"
	"testing"

	at "github.com/dcarbone/terraform-plugin-framework-utils/v3/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp/fakes"
)

var (
	testDataSourceItemsFQN    = testhelp.DataSourceFQN("fabric", itemsTFName, "test")
	testDataSourceItemsHeader = at.DataSourceHeader(testhelp.TypeName("fabric", itemsTFName), "test")
)

func TestUnit_WorkspacesDataSource(t *testing.T) {
	entity := fakes.NewRandomWorkspaceInfoWithType(fabcore.WorkspaceTypeWorkspace)

	fakes.FakeServer.Upsert(fakes.NewRandomWorkspaceInfoWithType(fabcore.WorkspaceTypeWorkspace))
	fakes.FakeServer.Upsert(entity)
	fakes.FakeServer.Upsert(fakes.NewRandomWorkspaceInfoWithType(fabcore.WorkspaceTypePersonal))
	fakes.FakeServer.Upsert(fakes.NewRandomWorkspaceInfoWithType(fabcore.WorkspaceTypeWorkspace))
	fakes.FakeServer.Upsert(fakes.NewRandomWorkspaceInfoWithType(fabcore.WorkspaceTypeAdminWorkspace))

	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, nil, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		// error - unexpected_attr
		{
			Config: at.CompileConfig(
				testDataSourceItemsHeader,
				map[string]any{
					"unexpected_attr": "test",
				},
			),
			ExpectError: regexp.MustCompile(`An argument named "unexpected_attr" is not expected here`),
		},
		// read
		{
			Config: at.CompileConfig(
				testDataSourceItemsHeader,
				map[string]any{},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrSet(testDataSourceItemsFQN, "values.1.id"),
			),
		},
	}))
}

func TestAcc_WorkspacesDataSource(t *testing.T) {
	resource.ParallelTest(t, testhelp.NewTestAccCase(t, nil, nil, []resource.TestStep{
		// read
		{
			Config: at.CompileConfig(
				testDataSourceItemsHeader,
				map[string]any{},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrSet(testDataSourceItemsFQN, "values.0.id"),
			),
		},
	},
	))
}
