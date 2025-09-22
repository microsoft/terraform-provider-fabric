package itemjobscheduler_test

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"testing"
	"time"

	at "github.com/dcarbone/terraform-plugin-framework-utils/v3/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/services/itemjobscheduler"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp/fakes"
)

var testResourceItemFQN, testResourceItemHeader = testhelp.TFResource(common.ProviderTypeName, itemTypeInfo.Type, "test")

func TestUnit_ItemJobSchedulerResource_Attributes(t *testing.T) {
	fakes.FakeServer.ServerFactory.Core.ItemsServer.GetItem = fakeGetFabricItem("test")

	resource.Test(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		// error - no attributes
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{},
			),
			ExpectError: regexp.MustCompile(`Missing required argument`),
		},
		// error - no required attributes - item id
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
		// error - no required attributes - job type
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
		// error - no required attributes - enabled
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
		// error - no required attributes - configuration
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
		// missing required attribute for cron type - interval
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
						"start_date_time": time.Now().UTC().Format(time.RFC3339),
						"end_date_time":   time.Now().Add(24 * time.Hour).UTC().Format(time.RFC3339),
						"type":            string(fabcore.ScheduleTypeCron),
					},
				},
			),
			ExpectError: regexp.MustCompile("Invalid configuration for attribute configuration.interval"),
		},
		// missing required attribute for daily type - times
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
						"start_date_time": time.Now().UTC().Format(time.RFC3339),
						"end_date_time":   time.Now().Add(24 * time.Hour).UTC().Format(time.RFC3339),
						"type":            string(fabcore.ScheduleTypeDaily),
					},
				},
			),
			ExpectError: regexp.MustCompile("Invalid configuration for attribute configuration.times"),
		},
		// missing required attribute for weekly type - times
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
						"start_date_time": time.Now().UTC().Format(time.RFC3339),
						"end_date_time":   time.Now().Add(24 * time.Hour).UTC().Format(time.RFC3339),
						"type":            string(fabcore.ScheduleTypeWeekly),
					},
				},
			),
			ExpectError: regexp.MustCompile("Invalid configuration for attribute configuration.times"),
		},
		// missing required attribute for weekly type - weekdays
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
						"start_date_time": time.Now().UTC().Format(time.RFC3339),
						"end_date_time":   time.Now().Add(24 * time.Hour).UTC().Format(time.RFC3339),
						"type":            string(fabcore.ScheduleTypeWeekly),
						"times":           []string{"09:00"},
					},
				},
			),
			ExpectError: regexp.MustCompile("Invalid configuration for attribute configuration.weekdays"),
		},
		// error - times - invalid string format
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
						"start_date_time": time.Now().UTC().Format(time.RFC3339),
						"end_date_time":   time.Now().Add(24 * time.Hour).UTC().Format(time.RFC3339),
						"type":            string(fabcore.ScheduleTypeDaily),
						"times":           []string{"9:0:0"},
					},
				},
			),
			ExpectError: regexp.MustCompile(common.ErrorAttValueMatch),
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
						"start_date_time": time.Now().Format(time.RFC3339),
						"end_date_time":   time.Now().Add(24 * time.Hour).Format(time.RFC3339),
						"type":            string(fabcore.ScheduleTypeCron),
						"interval":        testhelp.RandomIntRange(1, 60),
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
						"start_date_time": time.Now().Format(time.RFC3339),
						"end_date_time":   time.Now().Add(24 * time.Hour).Format(time.RFC3339),
						"type":            string(fabcore.ScheduleTypeCron),
						"interval":        testhelp.RandomIntRange(1, 60),
					},
					"unexpected_attr": "test",
				},
			),
			ExpectError: regexp.MustCompile(`An argument named "unexpected_attr" is not expected here`),
		},
		// error  - not a valid date time - start date time with UTC+3 offset
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
						"start_date_time": time.Now().In(time.FixedZone("UTC+3", 3*60*60)).Format(time.RFC3339),
						"end_date_time":   time.Now().Add(24 * time.Hour).UTC().Format(time.RFC3339),
						"type":            string(fabcore.ScheduleTypeCron),
						"interval":        testhelp.RandomIntRange(1, 60),
					},
				},
			),
			ExpectError: regexp.MustCompile(common.ErrorAttValueMatch),
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
						"start_date_time": time.Now().UTC().Format(time.RFC3339),
						"end_date_time":   time.Now().Add(24 * time.Hour).UTC().Format(time.RFC3339),
						"type":            string(fabcore.ScheduleTypeCron),
						"interval":        testhelp.RandomIntRange(1, 60),
					},
				},
			),
			ExpectError: regexp.MustCompile(`Error: Invalid Item Type`),
		},
	}))
}

func TestUnit_ItemJobSchedulerResource_ImportState(t *testing.T) {
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
				"start_date_time": configuration.StartDateTime.Format(time.RFC3339),
				"end_date_time":   configuration.EndDateTime.Format(time.RFC3339),
				"type":            string(*configuration.Type),
				"interval":        int(*entity.Configuration.(*fabcore.CronScheduleConfig).Interval),
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

func TestUnit_ItemJobSchedulerResource_CRUD(t *testing.T) {
	itemType := "dataflow"
	workspaceID := testhelp.RandomUUID()
	itemID := testhelp.RandomUUID()
	jobType := itemjobscheduler.JobTypeActions[itemType][0]
	entity := NewRandomItemSchedule(fabcore.ScheduleTypeCron)
	entityUpdate := NewRandomItemSchedule(fabcore.ScheduleTypeWeekly)

	weekdays := entityUpdate.Configuration.(*fabcore.WeeklyScheduleConfig).Weekdays

	weekdaysStr := make([]string, len(weekdays))
	for i, d := range weekdays {
		weekdaysStr[i] = string(d)
	}

	fakeTestUpsert(workspaceID, entity)

	fakes.FakeServer.ServerFactory.Core.ItemsServer.GetItem = fakeGetFabricItem(itemType)
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
					"job_type":     jobType,
					"enabled":      *entity.Enabled,
					"configuration": map[string]any{
						"start_date_time": entity.Configuration.GetScheduleConfig().StartDateTime.Format(time.RFC3339),
						"end_date_time":   entity.Configuration.GetScheduleConfig().EndDateTime.Format(time.RFC3339),
						"type":            string(*entity.Configuration.GetScheduleConfig().Type),
						"interval":        int(*entity.Configuration.(*fabcore.CronScheduleConfig).Interval),
					},
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "enabled", strconv.FormatBool(*entity.Enabled)),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "created_date_time"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "owner.id"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "owner.type"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "configuration.start_date_time", entity.Configuration.GetScheduleConfig().StartDateTime.Format(time.RFC3339)),
				resource.TestCheckResourceAttr(testResourceItemFQN, "configuration.end_date_time", entity.Configuration.GetScheduleConfig().EndDateTime.Format(time.RFC3339)),
				resource.TestCheckResourceAttr(testResourceItemFQN, "configuration.type", string(fabcore.ScheduleTypeCron)),
				resource.TestCheckResourceAttr(testResourceItemFQN, "configuration.interval", strconv.Itoa(int(*entity.Configuration.(*fabcore.CronScheduleConfig).Interval))),
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
					"job_type":     jobType,
					"enabled":      *entityUpdate.Enabled,
					"configuration": map[string]any{
						"start_date_time": entityUpdate.Configuration.GetScheduleConfig().StartDateTime.Format(time.RFC3339),
						"end_date_time":   entityUpdate.Configuration.GetScheduleConfig().EndDateTime.Format(time.RFC3339),
						"type":            string(*entityUpdate.Configuration.GetScheduleConfig().Type),
						"times":           entityUpdate.Configuration.(*fabcore.WeeklyScheduleConfig).Times,
						"weekdays":        weekdaysStr,
					},
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "enabled", strconv.FormatBool(*entityUpdate.Enabled)),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "created_date_time"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "owner.id"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "owner.type"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "configuration.start_date_time", entityUpdate.Configuration.GetScheduleConfig().StartDateTime.Format(time.RFC3339)),
				resource.TestCheckResourceAttr(testResourceItemFQN, "configuration.end_date_time", entityUpdate.Configuration.GetScheduleConfig().EndDateTime.Format(time.RFC3339)),
				resource.TestCheckResourceAttr(testResourceItemFQN, "configuration.type", string(fabcore.ScheduleTypeWeekly)),
				resource.TestCheckResourceAttr(testResourceItemFQN, "configuration.times.0", entityUpdate.Configuration.(*fabcore.WeeklyScheduleConfig).Times[0]),
				resource.TestCheckResourceAttr(testResourceItemFQN, "configuration.weekdays.0", weekdaysStr[0]),
			),
		},
		// Delete testing automatically occurs in TestCase
	}))
}

func TestAcc_ItemJobSchedulerResource_CRUD(t *testing.T) {
	workspace := testhelp.WellKnown()["WorkspaceRS"].(map[string]any)
	workspaceID := workspace["id"].(string)
	jobType := itemjobscheduler.JobTypeActions["dataflow"][0]
	dataflowResourceHCL, dataflowResourceFQN := dataflowResource(t, workspaceID)
	entity := NewRandomItemSchedule(fabcore.ScheduleTypeCron)
	entityUpdate := NewRandomItemSchedule(fabcore.ScheduleTypeMonthly)

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
							"start_date_time": entity.Configuration.GetScheduleConfig().StartDateTime.Format(time.RFC3339),
							"end_date_time":   entity.Configuration.GetScheduleConfig().EndDateTime.Format(time.RFC3339),
							"type":            string(*entity.Configuration.GetScheduleConfig().Type),
							"interval":        int(*entity.Configuration.(*fabcore.CronScheduleConfig).Interval),
						},
					},
				)),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "enabled", strconv.FormatBool(true)),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "created_date_time"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "owner.id"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "owner.type"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "configuration.start_date_time", entity.Configuration.GetScheduleConfig().StartDateTime.Format(time.RFC3339)),
				resource.TestCheckResourceAttr(testResourceItemFQN, "configuration.end_date_time", entity.Configuration.GetScheduleConfig().EndDateTime.Format(time.RFC3339)),
				resource.TestCheckResourceAttr(testResourceItemFQN, "configuration.type", string(fabcore.ScheduleTypeCron)),
				resource.TestCheckResourceAttr(testResourceItemFQN, "configuration.interval", strconv.Itoa(int(*entity.Configuration.(*fabcore.CronScheduleConfig).Interval))),
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
							"start_date_time": entityUpdate.Configuration.GetScheduleConfig().StartDateTime.Format(time.RFC3339),
							"end_date_time":   entityUpdate.Configuration.GetScheduleConfig().EndDateTime.Format(time.RFC3339),
							"type":            string(*entityUpdate.Configuration.GetScheduleConfig().Type),
							"times":           []string{"10:00"},
							"recurrence":      1,
							"occurrence": map[string]any{
								"occurrence_type": string(fabcore.OccurrenceTypeDayOfMonth),
								"day_of_month":    15,
							},
						},
					},
				)),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "enabled", strconv.FormatBool(false)),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "created_date_time"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "owner.id"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "owner.type"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "configuration.start_date_time", entityUpdate.Configuration.GetScheduleConfig().StartDateTime.Format(time.RFC3339)),
				resource.TestCheckResourceAttr(testResourceItemFQN, "configuration.end_date_time", entityUpdate.Configuration.GetScheduleConfig().EndDateTime.Format(time.RFC3339)),
				resource.TestCheckResourceAttr(testResourceItemFQN, "configuration.type", string(fabcore.ScheduleTypeMonthly)),
				resource.TestCheckResourceAttr(testResourceItemFQN, "configuration.times.0", "10:00"),
			),
		},
		// Delete testing automatically occurs in TestCase
	}))
}
