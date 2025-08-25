// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package itemjobscheduler_test

import (
	"regexp"
	"strconv"
	"testing"
	"time"

	at "github.com/dcarbone/terraform-plugin-framework-utils/v3/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp/fakes"
)

var testDataSourceItemFQN, testDataSourceItemHeader = testhelp.TFDataSource(common.ProviderTypeName, itemTypeInfo.Type, "test")

func TestUnit_ItemJobSchedulerDataSource(t *testing.T) {
	workspaceID := testhelp.RandomUUID()
	itemID := testhelp.RandomUUID()
	jobType := testhelp.RandomName()
	entity := NewRandomItemSchedule(fabcore.ScheduleTypeCron)

	fakeTestUpsert(workspaceID, entity)

	fakes.FakeServer.ServerFactory.Core.JobSchedulerServer.GetItemSchedule = fakeGetItemScheduleFunc()

	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, nil, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		// error - no attributes
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{},
			),
			ExpectError: regexp.MustCompile(`The argument "workspace_id" is required, but no definition was found`),
		},
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"workspace_id": testhelp.RandomUUID(),
				},
			),
			ExpectError: regexp.MustCompile(`The argument "item_id" is required, but no definition was found`),
		},
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"workspace_id": testhelp.RandomUUID(),
					"item_id":      itemID,
				},
			),
			ExpectError: regexp.MustCompile(`The argument "job_type" is required, but no definition was found`),
		},
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"workspace_id": testhelp.RandomUUID(),
					"item_id":      itemID,
					"job_type":     jobType,
				},
			),
			ExpectError: regexp.MustCompile(`The argument "id" is required, but no definition was found`),
		},
		// error - workspace_id - invalid UUID
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"workspace_id": "invalid uuid",
					"item_id":      itemID,
					"job_type":     jobType,
					"id":           *entity.ID,
				},
			),
			ExpectError: regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
		// error - unexpected attribute
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"workspace_id":    workspaceID,
					"unexpected_attr": "test",
				},
			),
			ExpectError: regexp.MustCompile(`An argument named "unexpected_attr" is not expected here`),
		},
		// read by id
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"workspace_id": workspaceID,
					"item_id":      itemID,
					"job_type":     jobType,
					"id":           *entity.ID,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "enabled", strconv.FormatBool(*entity.Enabled)),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "created_date_time", entity.CreatedDateTime.Format(time.RFC3339)),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "owner.id", *entity.Owner.ID),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "owner.display_name", *entity.Owner.DisplayName),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "owner.type", string(*entity.Owner.Type)),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "configuration.start_date_time", entity.Configuration.GetScheduleConfig().StartDateTime.Format(time.RFC3339)),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "configuration.start_date_time", entity.Configuration.GetScheduleConfig().StartDateTime.Format(time.RFC3339)),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "configuration.type", string(*entity.Configuration.GetScheduleConfig().Type)),
			),
		},
		// read by id - not found
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"workspace_id": workspaceID,
					"item_id":      itemID,
					"job_type":     jobType,
					"id":           testhelp.RandomUUID(),
				},
			),
			ExpectError: regexp.MustCompile(common.ErrorReadHeader),
		},
	}))
}

func TestAcc_ItemJobSchedulerDataSource(t *testing.T) {
	workspace := testhelp.WellKnown()["WorkspaceDS"].(map[string]any)
	workspaceID := workspace["id"].(string)

	entity := testhelp.WellKnown()["ItemJobScheduler"].(map[string]any)
	entityID := entity["id"].(string)
	entityJobType := entity["jobType"].(string)
	entityItemID := entity["itemId"].(string)
	createdDateTime, _ := time.Parse("2006-01-02T15:04:05.99", entity["createdDateTime"].(string))

	resource.ParallelTest(t, testhelp.NewTestAccCase(t, &testDataSourceItemFQN, nil, []resource.TestStep{
		// read by id
		{
			ResourceName: testDataSourceItemFQN,
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"workspace_id": workspaceID,
					"item_id":      entityItemID,
					"job_type":     entityJobType,
					"id":           entityID,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "workspace_id", workspaceID),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "id", entityID),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "configuration.type", entity["configurationType"].(string)),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "created_date_time", createdDateTime.Format(time.RFC3339)),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "owner.type", entity["ownerType"].(string)),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "owner.id", entity["ownerId"].(string)),
			),
		},
		// read by id - not found
		{
			ResourceName: testDataSourceItemFQN,
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"workspace_id": workspaceID,
					"item_id":      entityItemID,
					"job_type":     entityJobType,
					"id":           testhelp.RandomUUID(),
				},
			),
			ExpectError: regexp.MustCompile(common.ErrorReadHeader),
		},
	}))
}
