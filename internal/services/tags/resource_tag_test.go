// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package tags_test

import (
	"testing"

	at "github.com/dcarbone/terraform-plugin-framework-utils/v3/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	fabadmin "github.com/microsoft/fabric-sdk-go/fabric/admin"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp/fakes"
)

var testResourceItemFQN, testResourceItemHeader = testhelp.TFResource(common.ProviderTypeName, itemTypeInfo.Type, "test")

func TestUnit_TagResource_CRUD(t *testing.T) {
	entity1DisplayName := testhelp.RandomName()
	entity2DisplayName := testhelp.RandomName()

	tag := NewRandomTag()
	fakeTestUpsert(tag)

	updateTag := NewRandomTag()

	fakes.FakeServer.ServerFactory.Admin.TagsServer.NewListTagsPager = fakeTagsFunc()
	fakes.FakeServer.ServerFactory.Admin.TagsServer.BulkCreateTags = fakeBulkCreateTagsFunc()
	fakes.FakeServer.ServerFactory.Admin.TagsServer.DeleteTag = fakeDeleteTagFunc()
	fakes.FakeServer.ServerFactory.Admin.TagsServer.UpdateTag = fakeUpdateTagFunc(updateTag)

	resource.Test(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		// Create and Read
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"tags": []map[string]any{
						{
							"display_name": entity1DisplayName,
						},
						{
							"display_name": entity2DisplayName,
						},
					},
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "tags.0.display_name"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "tags.1.display_name"),
			),
		},
		// Update and Read
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"id":           *tag.ID,
					"display_name": *updateTag.DisplayName,
					"scope": map[string]any{
						"type": string(fabadmin.TagScopeTypeTenant),
					},
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "id", tag.ID),
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "display_name", updateTag.DisplayName),
			),
		},
		// Delete testing automatically occurs in TestCase
	}))
}

func TestAcc_TagResource_CRUD(t *testing.T) {
	entity := testhelp.WellKnown()["Tags"].(map[string]any)
	entityID := entity["id"].(string)
	entityScope := entity["scopeType"].(string)
	entityUpdateDisplayName := testhelp.RandomName()
	entity1DisplayName := testhelp.RandomName()
	entity2DisplayName := testhelp.RandomName()

	resource.Test(t, testhelp.NewTestAccCase(t, &testResourceItemFQN, nil, []resource.TestStep{
		// Create and Read
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"tags": []map[string]any{
						{
							"display_name": entity1DisplayName,
						},
						{
							"display_name": entity2DisplayName,
						},
					},
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "tags.0.display_name", entity1DisplayName),
				resource.TestCheckResourceAttr(testResourceItemFQN, "tags.1.display_name", entity2DisplayName),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "tags.0.id"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "tags.1.id"),
			),
		},
		// Update and Read
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"id":           entityID,
					"display_name": entityUpdateDisplayName,
					"scope": map[string]any{
						"type": entityScope,
					},
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "id", entityID),
				resource.TestCheckResourceAttr(testResourceItemFQN, "display_name", entityUpdateDisplayName),
			),
		},
	}))
}
