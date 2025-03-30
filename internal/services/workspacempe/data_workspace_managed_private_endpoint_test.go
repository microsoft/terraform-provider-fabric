// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package workspacempe_test

import (
	"regexp"
	"testing"

	at "github.com/dcarbone/terraform-plugin-framework-utils/v3/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp/fakes"
)

var testDataSourceItemFQN, testDataSourceItemHeader = testhelp.TFDataSource(common.ProviderTypeName, itemTypeInfo.Type, "test")

func TestUnit_WorkspaceManagedPrivateEndpointDataSource(t *testing.T) {
	workspaceID := testhelp.RandomUUID()
	entity := fakes.NewRandomWorkspaceManagedPrivateEndpoint()

	fakes.FakeServer.Upsert(entity)

	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, nil, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		// read
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"workspace_id": workspaceID,
					"id":           *entity.ID,
				},
			),
			// Check: resource.ComposeAggregateTestCheckFunc(
			// 	resource.TestCheckResourceAttr(testDataSourceItemFQN, "workspace_id", workspaceID),
			// 	resource.TestCheckResourceAttrPtr(testDataSourceItemFQN, "values.0.id", entity.ID),
			// 	resource.TestCheckResourceAttrPtr(testDataSourceItemFQN, "values.0.name", entity.Name),
			// ),
			ConfigStateChecks: []statecheck.StateCheck{
				statecheck.ExpectKnownValue(
					testDataSourceItemFQN,
					tfjsonpath.New("workspace_id"),
					knownvalue.StringExact(workspaceID),
				),
				statecheck.ExpectKnownValue(
					testDataSourceItemFQN,
					tfjsonpath.New("id"),
					knownvalue.StringExact(*entity.ID),
				),
			},
		},
	}))
}

func TestAcc_WorkspaceManagedPrivateEndpointDataSource(t *testing.T) {
	workspace := testhelp.WellKnown()["WorkspaceMPE"].(map[string]any)
	workspaceID := workspace["id"].(string)

	entity := testhelp.WellKnown()["ManagedPrivateEndpoint"].(map[string]any)
	entityID := entity["id"].(string)
	entityName := entity["name"].(string)

	resource.ParallelTest(t, testhelp.NewTestAccCase(t, nil, nil, []resource.TestStep{
		// read by id
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"workspace_id": workspaceID,
					"id":           entityID,
				},
			),
			ConfigStateChecks: []statecheck.StateCheck{
				statecheck.ExpectKnownValue(
					testDataSourceItemFQN,
					tfjsonpath.New("id"),
					knownvalue.StringExact(entityID),
				),
				statecheck.ExpectKnownValue(
					testDataSourceItemFQN,
					tfjsonpath.New("name"),
					knownvalue.StringExact(entityName),
				),
			},
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
					"name":         entityName,
				},
			),
			ConfigStateChecks: []statecheck.StateCheck{
				statecheck.ExpectKnownValue(
					testDataSourceItemFQN,
					tfjsonpath.New("id"),
					knownvalue.StringExact(entityID),
				),
				statecheck.ExpectKnownValue(
					testDataSourceItemFQN,
					tfjsonpath.New("name"),
					knownvalue.StringExact(entityName),
				),
			},
		},
		// read by name - not found
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"workspace_id": workspaceID,
					"name":         testhelp.RandomName(),
				},
			),
			ExpectError: regexp.MustCompile(common.ErrorReadHeader),
		},
	}))
}
