// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package externaldatashare_test

import (
	"regexp"
	"testing"

	at "github.com/dcarbone/terraform-plugin-framework-utils/v3/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp/fakes"
)

var testDataSourceItemsFQN, testDataSourceItemsHeader = testhelp.TFDataSource(common.ProviderTypeName, itemTypeInfo.Types, "test")

func TestUnit_ExternalDataSharesDataSource(t *testing.T) {
	workspaceID := testhelp.RandomUUID()
	entity := NewRandomExternalDataShare(workspaceID)

	fakeTestUpsert(NewRandomExternalDataShare(workspaceID))
	fakeTestUpsert(entity)
	fakeTestUpsert(NewRandomExternalDataShare(workspaceID))

	fakes.FakeServer.ServerFactory.Core.ExternalDataSharesProviderServer.NewListExternalDataSharesInItemPager = fakeListExternalDataSharesProvider()

	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, nil, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		// error - no required attributes
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"workspace_id": workspaceID,
				},
			),
			ExpectError: regexp.MustCompile(`The argument "item_id" is required, but no definition was found.`),
		},
		// error - no required attributes
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"item_id": *entity.ItemID,
				},
			),
			ExpectError: regexp.MustCompile(`The argument "workspace_id" is required, but no definition was found.`),
		},
		// Read
		{
			Config: at.CompileConfig(
				testDataSourceItemsHeader,
				map[string]any{
					"workspace_id": workspaceID,
					"item_id":      *entity.ItemID,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrSet(testDataSourceItemsFQN, "values.0.id"),
				resource.TestCheckResourceAttrSet(testDataSourceItemsFQN, "values.0.status"),
				resource.TestCheckResourceAttrSet(testDataSourceItemsFQN, "values.0.invitation_url"),
				resource.TestCheckResourceAttrSet(testDataSourceItemsFQN, "values.0.recipient.user_principal_name"),
				resource.TestCheckResourceAttrSet(testDataSourceItemsFQN, "values.0.paths.0"),
			),
		},
	}))
}

func TestAcc_ExternalDataSharesDataSource(t *testing.T) {
	workspace := testhelp.WellKnown()["WorkspaceDS"].(map[string]any)
	workspaceID := workspace["id"].(string)

	lakehouse := testhelp.WellKnown()["Lakehouse"].(map[string]any)
	lakehouseID := lakehouse["id"].(string)

	resource.ParallelTest(t, testhelp.NewTestAccCase(t, nil, nil, []resource.TestStep{
		// read
		{
			Config: at.CompileConfig(
				testDataSourceItemsHeader,
				map[string]any{
					"workspace_id": workspaceID,
					"item_id":      lakehouseID,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrSet(testDataSourceItemsFQN, "values.0.id"),
				resource.TestCheckResourceAttrSet(testDataSourceItemsFQN, "values.0.status"),
				resource.TestCheckResourceAttrSet(testDataSourceItemsFQN, "values.0.invitation_url"),
				resource.TestCheckResourceAttrSet(testDataSourceItemsFQN, "values.0.recipient.user_principal_name"),
				resource.TestCheckResourceAttrSet(testDataSourceItemsFQN, "values.0.paths.0"),
			),
		},
	},
	))
}
