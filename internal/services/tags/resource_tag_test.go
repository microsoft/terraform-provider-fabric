// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package tags_test

import (
	"regexp"
	"testing"

	at "github.com/dcarbone/terraform-plugin-framework-utils/v3/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	fabadmin "github.com/microsoft/fabric-sdk-go/fabric/admin"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp/fakes"
)

var testResourceItemFQN, testResourceItemHeader = testhelp.TFResource(common.ProviderTypeName, itemTypeInfo.Type, "test")

func TestUnit_TagResource_Attributes(t *testing.T) {
	fakes.FakeServer.ServerFactory.Admin.TagsServer.NewListTagsPager = fakeTagsFunc()
	fakes.FakeServer.ServerFactory.Admin.TagsServer.BulkCreateTags = fakeBulkCreateTagsFunc()
	fakes.FakeServer.ServerFactory.Admin.TagsServer.DeleteTag = fakeDeleteTagFunc()

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
		// error - unexpected attribute
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"id":              testhelp.RandomUUID(),
					"unexpected_attr": "test",
				},
			),
			ExpectError: regexp.MustCompile(`An argument named "unexpected_attr" is not expected here`),
		},
		// error - no required attributes
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"scope": map[string]any{
						"type": string(fabadmin.TagScopeTypeTenant),
					},
				},
			),
			ExpectError: regexp.MustCompile(`The argument "display_name" is required, but no definition was found.`),
		},
	}))
}

func TestUnit_TagResource_CRUD(t *testing.T) {
	entity1DisplayName := testhelp.RandomName()
	entity2DisplayName := testhelp.RandomName()
	entity3DisplayName := testhelp.RandomName()
	entity4DisplayName := testhelp.RandomName()
	domainID := testhelp.RandomUUID()

	fakes.FakeServer.ServerFactory.Admin.TagsServer.NewListTagsPager = fakeTagsFunc()
	fakes.FakeServer.ServerFactory.Admin.TagsServer.BulkCreateTags = fakeBulkCreateTagsFunc()
	fakes.FakeServer.ServerFactory.Admin.TagsServer.DeleteTag = fakeDeleteTagFunc()
	fakes.FakeServer.ServerFactory.Admin.TagsServer.UpdateTag = fakeUpdateTagFunc()

	resource.Test(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		// Create and Read
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"display_name": entity1DisplayName,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "id"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "display_name", entity1DisplayName),
				resource.TestCheckResourceAttr(testResourceItemFQN, "scope.type", string(fabadmin.TagScopeTypeTenant)),
			),
		},
		// Create with Scope
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"display_name": entity2DisplayName,
					"scope": map[string]any{
						"type": string(fabadmin.TagScopeTypeTenant),
					},
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "id"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "display_name", entity2DisplayName),
				resource.TestCheckResourceAttr(testResourceItemFQN, "scope.type", string(fabadmin.TagScopeTypeTenant)),
			),
		},
		// Update display_name
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"display_name": entity3DisplayName,
					"scope": map[string]any{
						"type": string(fabadmin.TagScopeTypeTenant),
					},
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "id"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "display_name", entity3DisplayName),
				resource.TestCheckResourceAttr(testResourceItemFQN, "scope.type", string(fabadmin.TagScopeTypeTenant)),
			),
		},
		// Create with Domain Scope
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"display_name": entity4DisplayName,
					"scope": map[string]any{
						"type":      string(fabadmin.TagScopeTypeDomain),
						"domain_id": domainID,
					},
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "id"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "display_name", entity4DisplayName),
				resource.TestCheckResourceAttr(testResourceItemFQN, "scope.type", string(fabadmin.TagScopeTypeDomain)),
				resource.TestCheckResourceAttr(testResourceItemFQN, "scope.domain_id", domainID),
			),
		},
		// Delete testing automatically occurs in TestCase
	}))
}

func TestAcc_TagResource_CRUD(t *testing.T) {
	entity1DisplayName := testhelp.RandomName()
	entity2DisplayName := testhelp.RandomName()

	domain := testhelp.WellKnown()["DomainParent"].(map[string]any)
	domainID := domain["id"].(string)

	resource.Test(t, testhelp.NewTestAccCase(t, &testResourceItemFQN, nil, []resource.TestStep{
		// Create and Read
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"display_name": entity1DisplayName,
					"scope": map[string]any{
						"domain_id": domainID,
						"type":      string(fabadmin.TagScopeTypeDomain),
					},
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "id"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "display_name", entity1DisplayName),
				resource.TestCheckResourceAttr(testResourceItemFQN, "scope.type", string(fabadmin.TagScopeTypeDomain)),
				resource.TestCheckResourceAttr(testResourceItemFQN, "scope.domain_id", domainID),
			),
		},
		// Update display_name
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"display_name": entity2DisplayName,
					"scope": map[string]any{
						"domain_id": domainID,
						"type":      string(fabadmin.TagScopeTypeDomain),
					},
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "id"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "display_name", entity2DisplayName),
				resource.TestCheckResourceAttr(testResourceItemFQN, "scope.type", string(fabadmin.TagScopeTypeDomain)),
				resource.TestCheckResourceAttr(testResourceItemFQN, "scope.domain_id", domainID),
			),
		},
	}))
}
