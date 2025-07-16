// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package folder_test

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

func TestUnit_FoldersDataSource(t *testing.T) {
	workspaceID := testhelp.RandomUUID()
	entity := fakes.NewRandomFolderWithWorkspace(workspaceID)
	childEntity := fakes.NewRandomSubfolder(workspaceID, *entity.ID)
	fakes.FakeServer.Upsert(fakes.NewRandomFolderWithWorkspace(workspaceID))
	fakes.FakeServer.Upsert(entity)
	fakes.FakeServer.Upsert(childEntity)
	fakes.FakeServer.Upsert(fakes.NewRandomSubfolder(workspaceID, *childEntity.ID))

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
			),
			ConfigStateChecks: []statecheck.StateCheck{
				statecheck.ExpectKnownValue(
					testDataSourceItemsFQN,
					tfjsonpath.New("values"),
					knownvalue.SetPartial([]knownvalue.Check{
						knownvalue.ObjectPartial(map[string]knownvalue.Check{
							"id":               knownvalue.StringExact(*entity.ID),
							"display_name":     knownvalue.StringExact(*entity.DisplayName),
							"parent_folder_id": knownvalue.StringExact(*entity.ParentFolderID),
						}),
					}),
				),
			},
		},
		// read with options - recursive defaults to true
		{
			Config: at.CompileConfig(
				testDataSourceItemsHeader,
				map[string]any{
					"workspace_id":   workspaceID,
					"root_folder_id": *entity.ID,
				},
			),
			ConfigStateChecks: []statecheck.StateCheck{
				statecheck.ExpectKnownValue(
					testDataSourceItemsFQN,
					tfjsonpath.New("values"),
					knownvalue.SetSizeExact(2),
				),
			},
		},
		// read with options - recursive false
		{
			Config: at.CompileConfig(
				testDataSourceItemsHeader,
				map[string]any{
					"workspace_id":   workspaceID,
					"root_folder_id": *entity.ID,
					"recursive":      false,
				},
			),
			ConfigStateChecks: []statecheck.StateCheck{
				statecheck.ExpectKnownValue(
					testDataSourceItemsFQN,
					tfjsonpath.New("values"),
					knownvalue.SetSizeExact(1),
				),
			},
		},
	}))
}

func TestAcc_FoldersDataSource(t *testing.T) {
	workspace := testhelp.WellKnown()["WorkspaceDS"].(map[string]any)
	workspaceID := workspace["id"].(string)

	folder := testhelp.WellKnown()["Folder"].(map[string]any)
	folderID := folder["id"].(string)

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
			),
		},
		// read with options - recursive defaults to true
		{
			Config: at.CompileConfig(
				testDataSourceItemsHeader,
				map[string]any{
					"workspace_id":   workspaceID,
					"root_folder_id": folderID,
				},
			),
			ConfigStateChecks: []statecheck.StateCheck{
				statecheck.ExpectKnownValue(
					testDataSourceItemsFQN,
					tfjsonpath.New("values"),
					knownvalue.SetSizeExact(2),
				),
			},
		},
		// read with options - recursive = false
		{
			Config: at.CompileConfig(
				testDataSourceItemsHeader,
				map[string]any{
					"workspace_id":   workspaceID,
					"root_folder_id": folderID,
					"recursive":      false,
				},
			),
			ConfigStateChecks: []statecheck.StateCheck{
				statecheck.ExpectKnownValue(
					testDataSourceItemsFQN,
					tfjsonpath.New("values"),
					knownvalue.SetSizeExact(1),
				),
			},
		},
	},
	))
}
