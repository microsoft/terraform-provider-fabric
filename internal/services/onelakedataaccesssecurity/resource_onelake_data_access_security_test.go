// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package onelakedataaccesssecurity_test

import (
	"errors"
	"fmt"
	"regexp"
	"testing"

	at "github.com/dcarbone/terraform-plugin-framework-utils/v3/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp/fakes"
)

var testResourceItemFQN, testResourceItemHeader = testhelp.TFResource(common.ProviderTypeName, itemTypeInfo.Types, "test")

func TestUnit_OneLakeDataAccessSecurityResource_Attributes(t *testing.T) {
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
		// error - unexpected attribute
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"item_id":         testhelp.RandomUUID(),
					"workspace_id":    testhelp.RandomUUID(),
					"unexpected_attr": "test",
				},
			),
			ExpectError: regexp.MustCompile(`An argument named "unexpected_attr" is not expected here`),
		},
		// error - no required attribute - item_id
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
		// error - no required attribute - workspace_id
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"item_id": testhelp.RandomUUID(),
				},
			),
			ExpectError: regexp.MustCompile(`The argument "workspace_id" is required, but no definition was found.`),
		},
	}))
}

func TestUnit_OneLakeDataAccessSecurityResource_ImportState(t *testing.T) {
	workspaceID := testhelp.RandomUUID()
	itemID := testhelp.RandomUUID()
	entity := fakes.NewRandomOneLakeDataAccessesSecurityClient(itemID, workspaceID)

	fakes.FakeServer.Upsert(fakes.NewRandomOneLakeDataAccessesSecurityClient(testhelp.RandomUUID(), workspaceID))
	fakes.FakeServer.Upsert(entity)
	fakes.FakeServer.Upsert(fakes.NewRandomOneLakeDataAccessesSecurityClient(testhelp.RandomUUID(), workspaceID))

	testCase := at.CompileConfig(
		testResourceItemHeader,
		map[string]any{
			"workspace_id": workspaceID,
			"item_id":      itemID,
			"value": []map[string]any{
				{
					"name": *entity.Value[0].Name,
					"decision_rules": []map[string]any{
						{
							"effect": (string)(*entity.Value[0].DecisionRules[0].Effect),
							"permission": []map[string]any{
								{
									"attribute_name":              string(*entity.Value[0].DecisionRules[0].Permission[0].AttributeName),
									"attribute_value_included_in": []string{"*"},
								},
								{
									"attribute_name":              string(*entity.Value[0].DecisionRules[0].Permission[1].AttributeName),
									"attribute_value_included_in": []string{"Read"},
								},
							},
						},
					},
					"members": map[string]any{
						"fabric_item_members": []map[string]any{
							{
								"item_access": []string{"ReadAll"},
								"source_path": *entity.Value[0].Members.FabricItemMembers[0].SourcePath,
							},
						},
					},
				},
			},
		},
	)

	resource.Test(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		{
			ResourceName:  testResourceItemFQN,
			Config:        testCase,
			ImportStateId: "not-valid",
			ImportState:   true,
			ExpectError:   regexp.MustCompile(fmt.Sprintf(common.ErrorImportIdentifierDetails, "WorkspaceID/ItemID")),
		},
		{
			ResourceName:  testResourceItemFQN,
			Config:        testCase,
			ImportStateId: "test/id",
			ImportState:   true,
			ExpectError:   regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
		// Import state testing - onelake data access security
		{
			ResourceName:       testResourceItemFQN,
			Config:             testCase,
			ImportStateId:      fmt.Sprintf("%s/%s", workspaceID, itemID),
			ImportState:        true,
			ImportStatePersist: true,
			ImportStateCheck: func(is []*terraform.InstanceState) error {
				if len(is) != 1 {
					return errors.New("expected one instance state")
				}

				return nil
			},
		},
	}))
}

func TestUnit_OneLakeDataAccessSecurityResource_CRUD(t *testing.T) {
	entity := fakes.NewRandomOneLakeDataAccessesSecurityClient(testhelp.RandomUUID(), testhelp.RandomUUID())

	fakes.FakeServer.ServerFactory.Core.OneLakeDataAccessSecurityServer.CreateOrUpdateDataAccessRoles = fakeCreateOrUpdateOneLakeDataAccessSecurity()

	resource.Test(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		// Create and Read
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": testhelp.RandomUUID(),
					"item_id":      testhelp.RandomUUID(),
					"value": []map[string]any{
						{
							"name": *entity.Value[0].Name,
							"decision_rules": []map[string]any{
								{
									"effect": (string)(*entity.Value[0].DecisionRules[0].Effect),
									"permission": []map[string]any{
										{
											"attribute_name":              string(*entity.Value[0].DecisionRules[0].Permission[0].AttributeName),
											"attribute_value_included_in": []string{"*"},
										},
										{
											"attribute_name":              string(*entity.Value[0].DecisionRules[0].Permission[1].AttributeName),
											"attribute_value_included_in": []string{"Read"},
										},
									},
								},
							},
							"members": map[string]any{
								"fabric_item_members": []map[string]any{
									{
										"item_access": []string{"ReadAll"},
										"source_path": *entity.Value[0].Members.FabricItemMembers[0].SourcePath,
									},
								},
							},
						},
					},
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "etag"),
			),
		},
	}))
}

func TestAcc_OneLakeDataAccessSecurityResource_CRUD(t *testing.T) {
	workspace := testhelp.WellKnown()["WorkspaceRS"].(map[string]any)
	workspaceID := workspace["id"].(string)

	lakehouse := testhelp.WellKnown()["LakehouseRS"].(map[string]any)
	itemID := lakehouse["id"].(string)

	entity := fakes.NewRandomOneLakeDataAccessesSecurityClient(itemID, workspaceID)

	resource.ParallelTest(t, testhelp.NewTestAccCase(t, &testResourceItemFQN, nil, []resource.TestStep{
		{
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": workspaceID,
					"item_id":      itemID,
					"value": []map[string]any{
						{
							"name": *entity.Value[0].Name,
							"decision_rules": []map[string]any{
								{
									"effect": (string)(*entity.Value[0].DecisionRules[0].Effect),
									"permission": []map[string]any{
										{
											"attribute_name":              string(*entity.Value[0].DecisionRules[0].Permission[0].AttributeName),
											"attribute_value_included_in": []string{"*"},
										},
										{
											"attribute_name":              string(*entity.Value[0].DecisionRules[0].Permission[1].AttributeName),
											"attribute_value_included_in": []string{"Read"},
										},
									},
								},
							},
							"members": map[string]any{
								"fabric_item_members": []map[string]any{
									{
										"item_access": []string{"ReadAll"},
										"source_path": *entity.Value[0].Members.FabricItemMembers[0].SourcePath,
									},
								},
							},
						},
					},
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "etag"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "value.0.name", *entity.Value[0].Name),
			),
		},
	}))
}
