// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package itemjob_test

import (
	"regexp"
	"testing"

	at "github.com/dcarbone/terraform-plugin-framework-utils/v3/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp/fakes"
)

var testDataSourceItemFQN, testDataSourceItemHeader = testhelp.TFDataSource(common.ProviderTypeName, itemTypeInfo.Type, "test")

func TestUnit_ItemJobDataSource_Attributes(t *testing.T) {
	resource.Test(t, testhelp.NewTestUnitCase(t, nil, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		// error - no attributes
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{},
			),
			ExpectError: regexp.MustCompile(`Missing required argument`),
		},
		// error - unexpected attribute
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"workspace_id":    testhelp.RandomUUID(),
					"item_id":         testhelp.RandomUUID(),
					"id":              testhelp.RandomUUID(),
					"unexpected_attr": "test",
				},
			),
			ExpectError: regexp.MustCompile(`An argument named "unexpected_attr" is not expected here`),
		},
	}))
}

func TestUnit_ItemJobDataSource_Read(t *testing.T) {
	workspaceID := testhelp.RandomUUID()
	itemID := testhelp.RandomUUID()
	jobType := "Execute"
	entity := NewRandomItemJobInstance(workspaceID, itemID, jobType)

	fakes.FakeServer.ServerFactory.Core.JobSchedulerServer.GetItemJobInstance = fakeGetItemJobInstanceFunc()

	resource.Test(t, testhelp.NewTestUnitCase(t, nil, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		// read
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"workspace_id": workspaceID,
					"item_id":      itemID,
					"id":           *entity.ID,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "workspace_id", workspaceID),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "item_id", itemID),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "id", *entity.ID),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "job_type", jobType),
				resource.TestCheckResourceAttrSet(testDataSourceItemFQN, "status"),
			),
		},
	}))
}

func TestUnit_ItemJobDataSource_FailureReason(t *testing.T) {
	workspaceID := testhelp.RandomUUID()
	itemID := testhelp.RandomUUID()
	jobType := "Execute"
	entity := NewRandomFailedItemJobInstance(workspaceID, itemID, jobType)

	fakes.FakeServer.ServerFactory.Core.JobSchedulerServer.GetItemJobInstance = fakeGetItemJobInstanceFunc()

	resource.Test(t, testhelp.NewTestUnitCase(t, nil, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"workspace_id": workspaceID,
					"item_id":      itemID,
					"id":           *entity.ID,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "status", string(fabcore.ItemJobStatusFailed)),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "failure_reason.error_code", *entity.FailureReason.ErrorCode),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "failure_reason.message", *entity.FailureReason.Message),
			),
		},
	}))
}

func TestUnit_ItemJobDataSource_NotFound(t *testing.T) {
	workspaceID := testhelp.RandomUUID()
	itemID := testhelp.RandomUUID()

	fakes.FakeServer.ServerFactory.Core.JobSchedulerServer.GetItemJobInstance = fakeGetItemJobInstanceFunc()

	resource.Test(t, testhelp.NewTestUnitCase(t, nil, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"workspace_id": workspaceID,
					"item_id":      itemID,
					"id":           testhelp.RandomUUID(),
				},
			),
			ExpectError: regexp.MustCompile(common.ErrorReadHeader),
		},
	}))
}
