// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package externaldatashareprovider_test

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

func TestUnit_ExternalDataShareDataSource(t *testing.T) {
	workspaceID := testhelp.RandomUUID()
	entity := NewRandomExternalDataShare(workspaceID)

	fakeTestUpsert(NewRandomExternalDataShare(workspaceID))
	fakeTestUpsert(entity)
	fakeTestUpsert(NewRandomExternalDataShare(workspaceID))

	fakes.FakeServer.ServerFactory.Core.ExternalDataSharesProviderServer.GetExternalDataShare = fakeGetExternalDataShareProvider()

	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, nil, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		// error - workspace id - invalid UUID
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"workspace_id": "invalid uuid",
					"item_id":      *entity.ItemID,
					"id":           *entity.ID,
				},
			),
			ExpectError: regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
		// error - item id - invalid UUID
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"workspace_id": *entity.WorkspaceID,
					"item_id":      "invalid uuid",
					"id":           *entity.ID,
				},
			),
			ExpectError: regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
		// error - external data share id - invalid UUID
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"workspace_id": *entity.WorkspaceID,
					"item_id":      *entity.ItemID,
					"id":           "invalid uuid",
				},
			),
			ExpectError: regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
		// error - no required attributes
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"item_id": *entity.ItemID,
					"id":      *entity.ID,
				},
			),
			ExpectError: regexp.MustCompile(`The argument "workspace_id" is required, but no definition was found.`),
		},
		// error - no required attributes
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"item_id":      *entity.ItemID,
					"workspace_id": *entity.WorkspaceID,
				},
			),
			ExpectError: regexp.MustCompile(`Missing required argument`),
		},
		// error - no required attributes
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"workspace_id": *entity.WorkspaceID,
					"id":           *entity.ID,
				},
			),
			ExpectError: regexp.MustCompile(`The argument "item_id" is required, but no definition was found.`),
		},
		// error - unexpected attribute
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"workspace_id":    *entity.WorkspaceID,
					"item_id":         *entity.ItemID,
					"id":              *entity.ID,
					"unexpected_attr": "test",
				},
			),
			ExpectError: regexp.MustCompile(`An argument named "unexpected_attr" is not expected here`),
		},
		// Read
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"workspace_id": *entity.WorkspaceID,
					"item_id":      *entity.ItemID,
					"id":           *entity.ID,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrSet(testDataSourceItemFQN, "workspace_id"),
				resource.TestCheckResourceAttrSet(testDataSourceItemFQN, "item_id"),
				resource.TestCheckResourceAttrSet(testDataSourceItemFQN, "id"),
				resource.TestCheckResourceAttrSet(testDataSourceItemFQN, "status"),
				resource.TestCheckResourceAttrSet(testDataSourceItemFQN, "invitation_url"),
			),
		},
	}))
}

func TestAcc_ExternalDataShareDataSource(t *testing.T) {
	workspace := testhelp.WellKnown()["WorkspaceDS"].(map[string]any)
	workspaceID := workspace["id"].(string)

	lakehouse := testhelp.WellKnown()["Lakehouse"].(map[string]any)
	lakehouseID := lakehouse["id"].(string)

	externalDataShare := testhelp.WellKnown()["ExternalDataShare"].(map[string]any)
	externalDataShareID := externalDataShare["id"].(string)

	resource.ParallelTest(t, testhelp.NewTestAccCase(t, nil, nil, []resource.TestStep{
		// read
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"workspace_id": workspaceID,
					"item_id":      lakehouseID,
					"id":           externalDataShareID,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "workspace_id", workspaceID),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "item_id", lakehouseID),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "id", externalDataShareID),
				resource.TestCheckResourceAttrSet(testDataSourceItemFQN, "id"),
				resource.TestCheckResourceAttrSet(testDataSourceItemFQN, "status"),
				resource.TestCheckResourceAttrSet(testDataSourceItemFQN, "invitation_url"),
				resource.TestCheckResourceAttrSet(testDataSourceItemFQN, "recipient.user_principal_name"),
				resource.TestCheckResourceAttrSet(testDataSourceItemFQN, "paths.0"),
			),
		},
	},
	))
}
