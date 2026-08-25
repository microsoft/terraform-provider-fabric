// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package paginatedreport_test

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

func TestUnit_PaginatedReportDataSource(t *testing.T) {
	workspaceID := testhelp.RandomUUID()
	entity := fakes.NewRandomItemWithWorkspace(fabricItemType, workspaceID)
	fakes.FakeServer.Upsert(entity)

	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, nil, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		{
			Config:      at.CompileConfig(testDataSourceItemHeader, map[string]any{}),
			ExpectError: regexp.MustCompile(`The argument "workspace_id" is required`),
		},
		{
			Config: at.CompileConfig(testDataSourceItemHeader, map[string]any{
				"workspace_id": "invalid uuid",
				"id":           *entity.ID,
			}),
			ExpectError: regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
		{
			Config: at.CompileConfig(testDataSourceItemHeader, map[string]any{
				"workspace_id": workspaceID,
				"id":           *entity.ID,
				"display_name": *entity.DisplayName,
			}),
			ExpectError: regexp.MustCompile(`These attributes cannot be configured together`),
		},
		{
			Config: at.CompileConfig(testDataSourceItemHeader, map[string]any{
				"workspace_id": workspaceID,
				"id":           *entity.ID,
			}),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrPtr(testDataSourceItemFQN, "id", entity.ID),
				resource.TestCheckResourceAttrPtr(testDataSourceItemFQN, "display_name", entity.DisplayName),
			),
		},
		{
			Config: at.CompileConfig(testDataSourceItemHeader, map[string]any{
				"workspace_id": workspaceID,
				"id":           testhelp.RandomUUID(),
			}),
			ExpectError: regexp.MustCompile(common.ErrorReadHeader),
		},
		{
			Config: at.CompileConfig(testDataSourceItemHeader, map[string]any{
				"workspace_id": workspaceID,
				"display_name": *entity.DisplayName,
			}),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrPtr(testDataSourceItemFQN, "id", entity.ID),
				resource.TestCheckResourceAttrPtr(testDataSourceItemFQN, "display_name", entity.DisplayName),
			),
		},
		{
			Config: at.CompileConfig(testDataSourceItemHeader, map[string]any{
				"workspace_id": workspaceID,
				"display_name": testhelp.RandomName(),
			}),
			ExpectError: regexp.MustCompile(common.ErrorReadHeader),
		},
	}))
}

func TestAcc_PaginatedReportDataSource(t *testing.T) {
	workspace := testhelp.WellKnown()["WorkspaceDS"].(map[string]any)
	workspaceID := workspace["id"].(string)
	entity := testhelp.WellKnown()["PaginatedReport"].(map[string]any)
	entityID := entity["id"].(string)
	entityDisplayName := entity["displayName"].(string)
	entityDescription := entity["description"].(string)

	resource.ParallelTest(t, testhelp.NewTestAccCase(t, nil, nil, []resource.TestStep{
		{
			Config: at.CompileConfig(testDataSourceItemHeader, map[string]any{
				"workspace_id": workspaceID,
				"id":           entityID,
			}),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "id", entityID),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "display_name", entityDisplayName),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "description", entityDescription),
			),
		},
		{
			Config: at.CompileConfig(testDataSourceItemHeader, map[string]any{
				"workspace_id": workspaceID,
				"display_name": entityDisplayName,
			}),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "id", entityID),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "display_name", entityDisplayName),
			),
		},
	}))
}
