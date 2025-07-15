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

func TestUnit_OneLakeDataAccessSecurityResource_CRUD(t *testing.T) {
	entity := fakes.NewRandomOneLakeDataAccessesSecurityClient(testhelp.RandomUUID())

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
