// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package workspace_test

import (
	"regexp"
	"testing"

	at "github.com/dcarbone/terraform-plugin-framework-utils/v3/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp/fakes"
)

var testDataSourceItemsFQN, testDataSourceItemsHeader = testhelp.TFDataSource(common.ProviderTypeName, itemTypeInfo.Types, "test")

func TestUnit_WorkspacesDataSource(t *testing.T) {
	capacity := fakes.NewRandomCapacity()
	entity := fakes.NewRandomWorkspaceInfoWithType(fabcore.WorkspaceTypeWorkspace, capacity.ID)

	fakes.FakeServer.Upsert(capacity)
	fakes.FakeServer.Upsert(fakes.NewRandomWorkspaceInfoWithType(fabcore.WorkspaceTypeWorkspace, capacity.ID))
	fakes.FakeServer.Upsert(entity)
	fakes.FakeServer.Upsert(fakes.NewRandomWorkspaceInfoWithType(fabcore.WorkspaceTypePersonal, capacity.ID))
	fakes.FakeServer.Upsert(fakes.NewRandomWorkspaceInfoWithType(fabcore.WorkspaceTypeWorkspace, capacity.ID))
	fakes.FakeServer.Upsert(fakes.NewRandomWorkspaceInfoWithType(fabcore.WorkspaceTypeAdminWorkspace, capacity.ID))

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
