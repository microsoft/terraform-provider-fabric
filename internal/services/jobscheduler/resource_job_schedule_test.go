package jobscheduler_test

import (
	"errors"
	"fmt"
	"regexp"
	"testing"
	"time"

	at "github.com/dcarbone/terraform-plugin-framework-utils/v3/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp/fakes"
)

var testResourceItemFQN, testResourceItemHeader = testhelp.TFResource(common.ProviderTypeName, itemTypeInfo.Type, "test")

func TestUnit_JobSchedulerResource_Attributes(t *testing.T) {
	fakes.FakeServer.ServerFactory.Core.ItemsServer.GetItem = fakeGetFabricItem("test")

	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		// error - no attributes
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{},
			),
			ExpectError: regexp.MustCompile(`Missing required argument`),
		},
		// error - no required attributes
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": testhelp.RandomUUID(),
				},
			),
			ExpectError: regexp.MustCompile(`The argument "item_id" is required, but no definition was found.`),
		},
		// error - no required attributes
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": testhelp.RandomUUID(),
				},
			),
			ExpectError: regexp.MustCompile(`The argument "item_id" is required, but no definition was found.`),
		},
		// error - no required attributes
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": testhelp.RandomUUID(),
					"item_id":      testhelp.RandomUUID(),
				},
			),
			ExpectError: regexp.MustCompile(`The argument "job_type" is required, but no definition was found.`),
		},
		// error - no required attributes
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": testhelp.RandomUUID(),
					"item_id":      testhelp.RandomUUID(),
					"job_type":     testhelp.RandomName(),
				},
			),
			ExpectError: regexp.MustCompile(`The argument "enabled" is required, but no definition was found.`),
		},
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": testhelp.RandomUUID(),
					"item_id":      testhelp.RandomUUID(),
					"job_type":     testhelp.RandomName(),
					"enabled":      true,
				},
			),
			ExpectError: regexp.MustCompile(`The argument "configuration" is required, but no definition was found.`),
		},
		// error - workspace_id - invalid UUID
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": "invalid uuid",
					"item_id":      testhelp.RandomUUID(),
					"job_type":     "test",
					"enabled":      true,

					"configuration": map[string]any{
						"local_time_zone": testhelp.RandomName(),
						"start_date_time": time.Now().Format(time.RFC3339),
						"end_date_time":   time.Now().Add(24 * time.Hour).Format(time.RFC3339),
						"type":            string(fabcore.ScheduleTypeCron),
					},
				},
			),
			ExpectError: regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
		// error - unexpected attribute
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": "00000000-0000-0000-0000-000000000000",
					"item_id":      testhelp.RandomUUID(),
					"job_type":     "test",
					"enabled":      true,

					"configuration": map[string]any{
						"local_time_zone": testhelp.RandomName(),
						"start_date_time": time.Now().Format(time.RFC3339),
						"end_date_time":   time.Now().Add(24 * time.Hour).Format(time.RFC3339),
						"type":            string(fabcore.ScheduleTypeCron),
					},
					"unexpected_attr": "test",
				},
			),
			ExpectError: regexp.MustCompile(`An argument named "unexpected_attr" is not expected here`),
		},

		// error  - not a valid item type
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": testhelp.RandomUUID(),
					"item_id":      testhelp.RandomUUID(),
					"job_type":     "test",
					"enabled":      true,

					"configuration": map[string]any{
						"local_time_zone": testhelp.RandomName(),
						"start_date_time": time.Now().Format(time.RFC3339),
						"end_date_time":   time.Now().Add(24 * time.Hour).Format(time.RFC3339),
						"type":            string(fabcore.ScheduleTypeCron),
					},
				},
			),
			ExpectError: regexp.MustCompile(`Error: Invalid Job Type`),
		},
	}))
}

func TestUnit_JobSchedulerResource_ImportState(t *testing.T) {
	workspaceID := testhelp.RandomUUID()
	itemID := testhelp.RandomUUID()
	jobType := testhelp.RandomName()
	entity := NewRandomItemSchedule(fabcore.ScheduleTypeCron)
	configuration := entity.Configuration.GetScheduleConfig()

	fakeTestUpsert(workspaceID, entity)

	fakes.FakeServer.ServerFactory.Core.JobSchedulerServer.GetItemSchedule = fakeGetItemScheduleFunc()
	fakes.FakeServer.ServerFactory.Core.JobSchedulerServer.DeleteItemSchedule = fakeDeleteItemScheduleFunc()
	fakes.FakeServer.ServerFactory.Core.JobSchedulerServer.UpdateItemSchedule = fakeUpdateItemScheduleFunc()
	testCase := at.CompileConfig(
		testResourceItemHeader,
		map[string]any{
			"workspace_id": workspaceID,
			"item_id":      itemID,
			"job_type":     jobType,
			"enabled":      *entity.Enabled,
			"configuration": map[string]any{
				"local_time_zone": *configuration.LocalTimeZoneID,
				"start_date_time": (*configuration.StartDateTime).Format(time.RFC3339),
				"end_date_time":   (*configuration.EndDateTime).Format(time.RFC3339),
				"type":            string(*configuration.Type),
			},
		},
	)

	resource.Test(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		{
			ResourceName:  testResourceItemFQN,
			Config:        testCase,
			ImportStateId: "not-valid",
			ImportState:   true,
			ExpectError:   regexp.MustCompile("WorkspaceID/ItemID/JobType/ScheduleID"),
		},
		{
			ResourceName:  testResourceItemFQN,
			Config:        testCase,
			ImportStateId: "test/id/test/id",
			ImportState:   true,
			ExpectError:   regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
		{
			ResourceName:  testResourceItemFQN,
			Config:        testCase,
			ImportStateId: fmt.Sprintf("%s/%s/%s/%s", "test", itemID, "test", *entity.ID),
			ImportState:   true,
			ExpectError:   regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
		{
			ResourceName:  testResourceItemFQN,
			Config:        testCase,
			ImportStateId: fmt.Sprintf("%s/%s/%s/%s", workspaceID, "test", jobType, *entity.ID),
			ImportState:   true,
			ExpectError:   regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
		{
			ResourceName:  testResourceItemFQN,
			Config:        testCase,
			ImportStateId: fmt.Sprintf("%s/%s/%s/%s", workspaceID, itemID, jobType, "test"),
			ImportState:   true,
			ExpectError:   regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
		// Import state testing
		{
			ResourceName:       testResourceItemFQN,
			Config:             testCase,
			ImportStateId:      fmt.Sprintf("%s/%s/%s/%s", workspaceID, itemID, jobType, *entity.ID),
			ImportState:        true,
			ImportStatePersist: true,
			ImportStateCheck: func(is []*terraform.InstanceState) error {
				if len(is) != 1 {
					return errors.New("expected one instance state")
				}

				if is[0].ID != *entity.ID {
					return errors.New(testResourceItemFQN + ": unexpected ID")
				}

				return nil
			},
		},
	}))
}

func TestUnit_JobSchedulerResource_CRUD(t *testing.T) {
	workspaceID := testhelp.RandomUUID()
	itemID := testhelp.RandomUUID()
	// jobType := testhelp.RandomName()
	entity := NewRandomItemSchedule(fabcore.ScheduleTypeCron)

	fakeTestUpsert(workspaceID, entity)

	fakes.FakeServer.ServerFactory.Core.ItemsServer.GetItem = fakeGetFabricItem("Dataflow")
	fakes.FakeServer.ServerFactory.Core.JobSchedulerServer.GetItemSchedule = fakeGetItemScheduleFunc()
	fakes.FakeServer.ServerFactory.Core.JobSchedulerServer.DeleteItemSchedule = fakeDeleteItemScheduleFunc()
	fakes.FakeServer.ServerFactory.Core.JobSchedulerServer.UpdateItemSchedule = fakeUpdateItemScheduleFunc()
	fakes.FakeServer.ServerFactory.Core.JobSchedulerServer.CreateItemSchedule = fakeCreateItemScheduleFunc()

	resource.Test(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		// Create and Read
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": workspaceID,
					"item_id":      itemID,
					"job_type":     "Execute",
					"enabled":      *entity.Enabled,
					"configuration": map[string]any{
						"local_time_zone": *entity.Configuration.GetScheduleConfig().LocalTimeZoneID,
						"start_date_time": (*entity.Configuration.GetScheduleConfig().StartDateTime).Format(time.RFC3339),
						"end_date_time":   (*entity.Configuration.GetScheduleConfig().EndDateTime).Format(time.RFC3339),
						"type":            string(*entity.Configuration.GetScheduleConfig().Type),
						"interval":        12,
					},
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "configuration.local_time_zone", entity.Configuration.GetScheduleConfig().LocalTimeZoneID),
			),
		},
		// Update and Read
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": workspaceID,
					"item_id":      itemID,
					"job_type":     "Execute",
					"enabled":      *entity.Enabled,
					"configuration": map[string]any{
						"local_time_zone": *entity.Configuration.GetScheduleConfig().LocalTimeZoneID,
						"start_date_time": (*entity.Configuration.GetScheduleConfig().StartDateTime).Format(time.RFC3339),
						"end_date_time":   (*entity.Configuration.GetScheduleConfig().EndDateTime).Format(time.RFC3339),
						"type":            string(*entity.Configuration.GetScheduleConfig().Type),
					},
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "configuration.local_time_zone", entity.Configuration.GetScheduleConfig().LocalTimeZoneID),
			),
		},
		// Delete testing automatically occurs in TestCase
	}))
}

func TestAcc_JobScheduleResource_CRUD(t *testing.T) {
	workspace := testhelp.WellKnown()["WorkspaceRS"].(map[string]any)
	workspaceID := workspace["id"].(string)
	jobType := "Execute"
	dataflowResourceHCL, dataflowResourceFQN := dataflowResource(t, workspaceID)
	entity := NewRandomItemSchedule(fabcore.ScheduleTypeCron)

	resource.Test(t, testhelp.NewTestAccCase(t, &testResourceItemFQN, nil, []resource.TestStep{ // Create and Read
		// Create and Read
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				dataflowResourceHCL,
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"workspace_id": workspaceID,
						"item_id":      testhelp.RefByFQN(dataflowResourceFQN, "id"),
						"job_type":     jobType,
						"enabled":      true,
						"configuration": map[string]any{
							"local_time_zone": "Central Standard Time",
							"start_date_time": "2024-04-28T00:00:00Z",
							"end_date_time":   "2024-04-30T23:59:00Z",
							"type":            string(*entity.Configuration.GetScheduleConfig().Type),
							"interval":        12,
						},
					},
				)),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "configuration.local_time_zone", "Central Standard Time"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "configuration.start_date_time", "2024-04-28T00:00:00Z"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "configuration.end_date_time", "2024-04-30T23:59:00Z"),
			),
		},
		// Update and Read
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				dataflowResourceHCL,
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"workspace_id": workspaceID,
						"item_id":      testhelp.RefByFQN(dataflowResourceFQN, "id"),
						"job_type":     jobType,
						"enabled":      false,
						"configuration": map[string]any{
							"local_time_zone": "Central Standard Time",
							"start_date_time": "2024-04-28T00:00:00Z",
							"end_date_time":   "2024-04-30T23:59:00Z",
							"type":            string(*entity.Configuration.GetScheduleConfig().Type),
							"interval":        12,
						},
					},
				)),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "configuration.local_time_zone", "Central Standard Time"),
			),
		},
		// Delete testing automatically occurs in TestCase
	}))
}
