// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package workspace_test

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

var testDataSourceItemFQN, testDataSourceItemHeader = testhelp.TFDataSource(common.ProviderTypeName, itemTypeInfo.Type, "test")

func TestUnit_WorkspaceDataSource(t *testing.T) {
	capacity := fakes.NewRandomCapacity()
	entity := fakes.NewRandomWorkspaceInfo(capacity.ID)
	entityTypePersonal := fakes.NewRandomWorkspaceInfoWithType(fabcore.WorkspaceTypePersonal, nil)
	entityTypeAdmin := fakes.NewRandomWorkspaceInfoWithType(fabcore.WorkspaceTypeAdminWorkspace, nil)


	fakes.FakeServer.Upsert(capacity)
	fakes.FakeServer.Upsert(fakes.NewRandomWorkspaceInfo(capacity.ID))
	fakes.FakeServer.Upsert(entity)
	fakes.FakeServer.Upsert(entityTypePersonal)
	fakes.FakeServer.Upsert(entityTypeAdmin)
	fakes.FakeServer.Upsert(fakes.NewRandomWorkspaceInfo(capacity.ID))

	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, nil, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		// error - no attributes
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{},
			),
			ExpectError: regexp.MustCompile(`Exactly one of these attributes must be configured: \[id,display_name\]`),
		},
		// error - id - invalid UUID
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"id": "invalid uuid",
				},
			),
			ExpectError: regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
		// error - unexpected attribute
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"id":              *entity.ID,
					"unexpected_attr": "test",
				},
			),
			ExpectError: regexp.MustCompile(`An argument named "unexpected_attr" is not expected here`),
		},
		// error - conflicting attributes
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"id":           *entity.ID,
					"display_name": *entity.DisplayName,
				},
			),
			ExpectError: regexp.MustCompile(`These attributes cannot be configured together: \[id,display_name\]`),
		},
		// error - personal type
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"id": *entityTypePersonal.ID,
				},
			),
			ExpectError: regexp.MustCompile(common.ErrorWorkspaceNotSupportedHeader),
		},
		// error - admin type
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"id": *entityTypeAdmin.ID,
				},
			),
			ExpectError: regexp.MustCompile(common.ErrorWorkspaceNotSupportedHeader),
		},
		// read by id
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"id": *entity.ID,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrPtr(testDataSourceItemFQN, "id", entity.ID),
				resource.TestCheckResourceAttrPtr(testDataSourceItemFQN, "display_name", entity.DisplayName),
				resource.TestCheckResourceAttrPtr(testDataSourceItemFQN, "description", entity.Description),
			),
		},
		// read by id - not found
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"id": testhelp.RandomUUID(),
				},
			),
			ExpectError: regexp.MustCompile(common.ErrorReadHeader),
		},

		// read by name
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"display_name": *entity.DisplayName,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrPtr(testDataSourceItemFQN, "id", entity.ID),
				resource.TestCheckResourceAttrPtr(testDataSourceItemFQN, "display_name", entity.DisplayName),
				resource.TestCheckResourceAttrPtr(testDataSourceItemFQN, "description", entity.Description),
			),
		},
		// read by name - not found
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"display_name": testhelp.RandomName(),
				},
			),
			ExpectError: regexp.MustCompile(common.ErrorReadHeader),
		},
	}))
}

func TestAcc_WorkspaceDataSource(t *testing.T) {
	entity := testhelp.WellKnown()["WorkspaceDS"].(map[string]any)
	entityID := entity["id"].(string)
	entityDisplayName := entity["displayName"].(string)
	entityDescription := entity["description"].(string)

	resource.ParallelTest(t, testhelp.NewTestAccCase(t, nil, nil, []resource.TestStep{
		// read by id
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"id": entityID,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "id", entityID),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "display_name", entityDisplayName),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "description", entityDescription),
			),
		},
		// read by id - not found
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"id": testhelp.RandomUUID(),
				},
			),
			ExpectError: regexp.MustCompile(common.ErrorReadHeader),
		},
		// read by name
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"display_name": entityDisplayName,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "id", entityID),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "display_name", entityDisplayName),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "description", entityDescription),
			),
		},
		// read by name - not found
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"display_name": testhelp.RandomName(),
				},
			),
			ExpectError: regexp.MustCompile(common.ErrorReadHeader),
		},
	}))
}
