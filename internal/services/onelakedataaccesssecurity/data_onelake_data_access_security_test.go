// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package onelakedataaccesssecurity_test

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

func TestUnit_OneLakeDataAccessSecurityDataSource(t *testing.T) {
	workspaceID := testhelp.RandomUUID()
	itemID := testhelp.RandomUUID()
	entity := fakes.NewRandomOneLakeDataAccessesSecurityClient(itemID)

	fakes.FakeServer.Upsert(fakes.NewRandomOneLakeDataAccessesSecurityClient(testhelp.RandomUUID()))
	fakes.FakeServer.Upsert(entity)
	fakes.FakeServer.Upsert(fakes.NewRandomOneLakeDataAccessesSecurityClient(testhelp.RandomUUID()))

	fakes.FakeServer.ServerFactory.Core.OneLakeDataAccessSecurityServer.ListDataAccessRoles = fakeListOneLakeDataAccessSecurity(entity)

	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, nil, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		// error - no required attributes - item_id
		{
			ResourceName: testDataSourceItemsFQN,
			Config: at.CompileConfig(
				testDataSourceItemsHeader,
				map[string]any{
					"workspace_id": workspaceID,
				},
			),
			ExpectError: regexp.MustCompile(`The argument "item_id" is required, but no definition was found.`),
		},
		// error - no required attributes - workspace_id
		{
			ResourceName: testDataSourceItemsFQN,
			Config: at.CompileConfig(
				testDataSourceItemsHeader,
				map[string]any{
					"item_id": itemID,
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
					"item_id":      itemID,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(),
		},
	}))
}

func TestAcc_OneLakeDataAccessSecurityDataSource(t *testing.T) {
	randomWorkspaceID := testhelp.RandomUUID()
	itemID := testhelp.RandomUUID()

	resource.ParallelTest(t, testhelp.NewTestAccCase(t, nil, nil, []resource.TestStep{
		// Read - not found
		{
			Config: at.CompileConfig(
				testDataSourceItemsHeader,
				map[string]any{
					"workspace_id": randomWorkspaceID,
					"item_id":      itemID,
				},
			),
			ExpectError: regexp.MustCompile(common.ErrorListHeader),
		},
	}))
}
