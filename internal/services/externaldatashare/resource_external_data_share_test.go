// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package externaldatashare_test

import (
	"regexp"
	"strings"
	"testing"

	at "github.com/dcarbone/terraform-plugin-framework-utils/v3/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

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
							"type":                "User",
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
							"type":                "User",
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
							"type":                "User",
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
							"type":                "User",
							"user_principal_name": "test@example.com",
						},
					},
				)),
			ExpectError: regexp.MustCompile(`The argument "paths" is required, but no definition was found.`),
		},
		// error - invalid path
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"workspace_id": workspaceID,
						"item_id":      itemID,
						"paths":        []string{"InvalidPath"},
						"recipient": map[string]any{
							"type":                "User",
							"user_principal_name": "test@example.com",
						},
					},
				)),
			ExpectError: regexp.MustCompile(`A valid path must start with 'Files/'\s+or 'Tables/'`),
		},
		// error - recipient.type - invalid value
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"workspace_id": workspaceID,
						"item_id":      itemID,
						"paths":        []string{"Files/MyFolder1"},
						"recipient": map[string]any{
							"type":                "InvalidType",
							"user_principal_name": "test@example.com",
						},
					},
				)),
			ExpectError: regexp.MustCompile(`Attribute recipient.type value must be one of`),
		},
		// error - recipient.user_principal_name - too long
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"workspace_id": workspaceID,
						"item_id":      itemID,
						"paths":        []string{"Files/MyFolder1"},
						"recipient": map[string]any{
							"type":                "User",
							"user_principal_name": strings.Repeat("a", 257),
						},
					},
				)),
			ExpectError: regexp.MustCompile(`Attribute recipient.user_principal_name string length must be at most\s+256`),
		},
		// error - recipient.user_principal_name - required when type is User
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"workspace_id": workspaceID,
						"item_id":      itemID,
						"paths":        []string{"Files/MyFolder1"},
						"recipient": map[string]any{
							"type": "User",
						},
					},
				)),
			ExpectError: regexp.MustCompile(`(?i)user_principal_name`),
		},
		// error - recipient.user_principal_name - must be null when type is ServicePrincipal
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"workspace_id": workspaceID,
						"item_id":      itemID,
						"paths":        []string{"Files/MyFolder1"},
						"recipient": map[string]any{
							"type":                "ServicePrincipal",
							"user_principal_name": "test@example.com",
							"tenant_id":           "00000000-0000-0000-0000-000000000000",
						},
					},
				)),
			ExpectError: regexp.MustCompile(`(?i)user_principal_name`),
		},
		// error - recipient.tenant_id - required when type is ServicePrincipal
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"workspace_id": workspaceID,
						"item_id":      itemID,
						"paths":        []string{"Files/MyFolder1"},
						"recipient": map[string]any{
							"type":         "ServicePrincipal",
							"principal_id": "00000000-0000-0000-0000-000000000000",
						},
					},
				)),
			ExpectError: regexp.MustCompile(`(?i)tenant_id`),
		},
	}))
}

func TestUnit_ExternalDataShareResource_CRUD(t *testing.T) {
	workspaceID := testhelp.RandomUUID()
	entity := NewRandomExternalDataShare(workspaceID)
	entityRecipient := entity.Recipient.(*fabcore.ExternalDataShareUserRecipient)

	fakeTestUpsert(NewRandomExternalDataShare(workspaceID))

	fakes.FakeServer.ServerFactory.Core.ExternalDataSharesProviderServer.CreateExternalDataShare = fakeCreateExternalDataShareProvider()
	fakes.FakeServer.ServerFactory.Core.ExternalDataSharesProviderServer.GetExternalDataShare = fakeGetExternalDataShareProvider()
	fakes.FakeServer.ServerFactory.Core.ExternalDataSharesProviderServer.DeleteExternalDataShare = fakeDeleteExternalDataShareProvider()

	resource.Test(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		// Create and Read - recipient.type defaults to User when not specified
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": workspaceID,
					"item_id":      *entity.ItemID,
					"paths":        entity.Paths,
					"recipient": map[string]any{
						"user_principal_name": *entityRecipient.UserPrincipalName,
					},
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "workspace_id", *entity.WorkspaceID),
				resource.TestCheckResourceAttr(testResourceItemFQN, "item_id", *entity.ItemID),
				resource.TestCheckResourceAttr(testResourceItemFQN, "recipient.type", string(fabcore.ExternalDataShareRecipientTypeUser)),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "id"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "status"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "invitation_url"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "recipient.user_principal_name"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "principal_model.id"),
			),
		},
	}))
}

func TestAcc_ExternalDataShareResource_CRUD(t *testing.T) {
	workspace := testhelp.WellKnown()["WorkspaceRS"].(map[string]any)
	workspaceID := workspace["id"].(string)

	lakehouseResourceHCL, lakehouseResourceFQN := lakehouseResource(t, workspaceID)
	userPrincipalName := "test@example.com"

	resource.ParallelTest(t, testhelp.NewTestAccCase(t, nil, nil, []resource.TestStep{
		// read
		{
			Config: at.JoinConfigs(
				lakehouseResourceHCL,
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"workspace_id": workspaceID,
						"item_id":      testhelp.RefByFQN(lakehouseResourceFQN, "id"),
						"paths":        []string{"Tables/publicholidays"},
						"recipient": map[string]any{
							"type":                "User",
							"user_principal_name": userPrincipalName,
						},
					},
				)),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "workspace_id", workspaceID),
				resource.TestCheckResourceAttr(testResourceItemFQN, "recipient.user_principal_name", userPrincipalName),
				resource.TestCheckResourceAttrPair(testResourceItemFQN, "item_id", lakehouseResourceFQN, "id"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "id"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "status"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "invitation_url"),
			),
		},
	},
	))
}
