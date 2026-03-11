// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package itemjobscheduler_test

import (
	"regexp"
	"testing"

	at "github.com/dcarbone/terraform-plugin-framework-utils/v3/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp/fakes"
)

var testDataSourceItemsFQN, testDataSourceItemsHeader = testhelp.TFDataSource(common.ProviderTypeName, itemTypeInfo.Types, "test")

func TestUnit_ItemJobSchedulersDataSource(t *testing.T) {
	workspaceID := testhelp.RandomUUID()
	itemID := testhelp.RandomUUID()
	jobType := testhelp.RandomName()

	fakeTestUpsert(workspaceID, NewRandomItemSchedule(fabcore.ScheduleTypeCron))
	fakeTestUpsert(workspaceID, NewRandomItemSchedule(fabcore.ScheduleTypeCron))

	fakes.FakeServer.ServerFactory.Core.JobSchedulerServer.NewListItemSchedulesPager = fakeItemSchedulesFunc()

	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, nil, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		// error - no attributes
		{
			Config: at.CompileConfig(
				testDataSourceItemsHeader,
				map[string]any{},
			),
			ExpectError: regexp.MustCompile(`The argument "workspace_id" is required, but no definition was found`),
		},
		{
			Config: at.CompileConfig(
				testDataSourceItemsHeader,
				map[string]any{
					"workspace_id": testhelp.RandomUUID(),
				},
			),
			ExpectError: regexp.MustCompile(`The argument "item_id" is required, but no definition was found`),
		},
		{
			Config: at.CompileConfig(
				testDataSourceItemsHeader,
				map[string]any{
					"workspace_id": testhelp.RandomUUID(),
					"item_id":      itemID,
				},
			),
			ExpectError: regexp.MustCompile(`The argument "job_type" is required, but no definition was found`),
		},
		// error - workspace_id - invalid UUID
		{
			Config: at.CompileConfig(
				testDataSourceItemsHeader,
				map[string]any{
					"workspace_id": "invalid uuid",
					"item_id":      itemID,
					"job_type":     jobType,
				},
			),
			ExpectError: regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
		// error - unexpected attribute
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
					"item_id":      itemID,
					"job_type":     jobType,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrSet(testDataSourceItemsFQN, "values.0.enabled"),
				resource.TestCheckResourceAttrSet(testDataSourceItemsFQN, "values.0.configuration.type"),
				resource.TestCheckResourceAttrSet(testDataSourceItemsFQN, "values.0.id"),
				resource.TestCheckResourceAttrSet(testDataSourceItemsFQN, "values.0.configuration.type"),
				resource.TestCheckResourceAttrSet(testDataSourceItemsFQN, "values.0.created_date_time"),
				resource.TestCheckResourceAttrSet(testDataSourceItemsFQN, "values.0.owner.type"),
				resource.TestCheckResourceAttrSet(testDataSourceItemsFQN, "values.0.owner.id"),
			),
		},
	}))
}

func TestAcc_ItemJobSchedulersDataSource(t *testing.T) {
	workspace := testhelp.WellKnown()["WorkspaceDS"].(map[string]any)
	workspaceID := workspace["id"].(string)

	entity := testhelp.WellKnown()["ItemJobScheduler"].(map[string]any)
	entityJobType := entity["jobType"].(string)
	entityItemID := entity["itemId"].(string)

	resource.ParallelTest(t, testhelp.NewTestAccCase(t, nil, nil, []resource.TestStep{
		// read
		{
			Config: at.CompileConfig(
				testDataSourceItemsHeader,
				map[string]any{
					"workspace_id": workspaceID,
					"item_id":      entityItemID,
					"job_type":     entityJobType,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testDataSourceItemsFQN, "workspace_id", workspaceID),
				resource.TestCheckResourceAttrSet(testDataSourceItemsFQN, "values.0.id"),
				resource.TestCheckResourceAttrSet(testDataSourceItemsFQN, "values.0.configuration.type"),
				resource.TestCheckResourceAttrSet(testDataSourceItemsFQN, "values.0.configuration.type"),
				resource.TestCheckResourceAttrSet(testDataSourceItemsFQN, "values.0.created_date_time"),
				resource.TestCheckResourceAttrSet(testDataSourceItemsFQN, "values.0.owner.type"),
				resource.TestCheckResourceAttrSet(testDataSourceItemsFQN, "values.0.owner.id"),
			),
		},
	},
	))
}
