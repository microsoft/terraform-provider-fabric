// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package externaldatasharesprovider_test

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

var testResourceItemFQN, testResourceItemHeader = testhelp.TFResource(common.ProviderTypeName, itemTypeInfo.Type, "test")

func TestUnit_ExternalDataShareResource_Attributes(t *testing.T) {
	workspaceID := testhelp.RandomUUID()
	itemID := testhelp.RandomUUID()

	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		// error - no attributes
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{},
				),
			),
			ExpectError: regexp.MustCompile(`Missing required argument`),
		},
		// error - workspace_id - invalid UUID
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"workspace_id": "invalid uuid",
						"item_id":      itemID,
						"paths":        []string{"Files/MyFolder1"},
						"recipient": map[string]any{
							"user_principal_name": "test@example.com",
						},
					},
				)),
			ExpectError: regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
		// error - item_id - invalid UUID
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"workspace_id": workspaceID,
						"item_id":      "invalid uuid",
						"paths":        []string{"Files/MyFolder1"},
						"recipient": map[string]any{
							"user_principal_name": "test@example.com",
						},
					},
				)),
			ExpectError: regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
		// error - unexpected attribute
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"workspace_id":    workspaceID,
						"item_id":         itemID,
						"unexpected_attr": "test",
						"paths":           []string{"Files/MyFolder1"},
						"recipient": map[string]any{
							"user_principal_name": "test@example.com",
						},
					},
				)),
			ExpectError: regexp.MustCompile(`An argument named "unexpected_attr" is not expected here`),
		},
		// error - no required attributes
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"workspace_id": workspaceID,
						"item_id":      itemID,
						"paths":        []string{"Files/MyFolder1"},
					},
				)),
			ExpectError: regexp.MustCompile(`The argument "recipient" is required, but no definition was found.`),
		},
		// error - no required attributes
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"workspace_id": workspaceID,
						"item_id":      itemID,
						"recipient": map[string]any{
							"user_principal_name": "test@example.com",
						},
					},
				)),
			ExpectError: regexp.MustCompile(`The argument "paths" is required, but no definition was found.`),
		},
	}))
}

func TestUnit_ExternalDataShareResource_CRUD(t *testing.T) {
	workspaceID := testhelp.RandomUUID()
	entity := NewRandomExternalDataShare(workspaceID)

	fakes.FakeServer.ServerFactory.Core.ExternalDataSharesProviderServer.CreateExternalDataShare = fakeCreateExternalDataShareProvider(entity)
	fakes.FakeServer.ServerFactory.Core.ExternalDataSharesProviderServer.GetExternalDataShare = fakeGetExternalDataShareProvider(entity)
	fakes.FakeServer.ServerFactory.Core.ExternalDataSharesProviderServer.DeleteExternalDataShare = fakeDeleteExternalDataShareProvider()

	resource.Test(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		// Create and Read
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": workspaceID,
					"item_id":      *entity.ItemID,
					"paths":        entity.Paths,
					"recipient": map[string]any{
						"user_principal_name": *entity.Recipient.UserPrincipalName,
					},
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "id", *entity.ID),
				resource.TestCheckResourceAttr(testResourceItemFQN, "workspace_id", *entity.WorkspaceID),
				resource.TestCheckResourceAttr(testResourceItemFQN, "item_id", *entity.ItemID),
				resource.TestCheckResourceAttr(testResourceItemFQN, "status", string(*entity.Status)),
				resource.TestCheckResourceAttr(testResourceItemFQN, "invitation_url", *entity.InvitationURL),
				resource.TestCheckResourceAttr(testResourceItemFQN, "recipient.user_principal_name", *entity.Recipient.UserPrincipalName),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "creator_principal.id"),
			),
		},
	}))
}

func TestAcc_ExternalDataShareResource_CRUD(t *testing.T) {
	workspace := testhelp.WellKnown()["WorkspaceDS"].(map[string]any)
	workspaceID := workspace["id"].(string)

	lakehouse := testhelp.WellKnown()["Lakehouse"].(map[string]any)
	lakehouseID := lakehouse["id"].(string)

	userPrincipalName := "test@example.com"

	resource.ParallelTest(t, testhelp.NewTestAccCase(t, nil, nil, []resource.TestStep{
		// read
		{
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": workspaceID,
					"item_id":      lakehouseID,
					"paths":        []string{"Tables/publicholidays"},
					"recipient": map[string]any{
						"user_principal_name": userPrincipalName,
					},
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "workspace_id", workspaceID),
				resource.TestCheckResourceAttr(testResourceItemFQN, "item_id", lakehouseID),
				resource.TestCheckResourceAttr(testResourceItemFQN, "recipient.user_principal_name", userPrincipalName),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "id"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "status"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "invitation_url"),
			),
		},
	},
	))
}
