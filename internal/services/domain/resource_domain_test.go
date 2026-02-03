// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package domain_test

import (
	"errors"
	"regexp"
	"testing"

	at "github.com/dcarbone/terraform-plugin-framework-utils/v3/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/microsoft/fabric-sdk-go/fabric/admin"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp/fakes"
)

var testResourceItemFQN, testResourceItemHeader = testhelp.TFResource(common.ProviderTypeName, itemTypeInfo.Type, "test")

func TestUnit_DomainResource_Attributes(t *testing.T) {
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
					"display_name":    "test",
					"unexpected_attr": "test",
				},
			),
			ExpectError: regexp.MustCompile(`An argument named "unexpected_attr" is not expected here`),
		},
		// // error - no required attributes
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"description": "test",
				},
			),
			ExpectError: regexp.MustCompile(`The argument "display_name" is required, but no definition was found.`),
		},
		// // error - invalid uuid - capacity_id
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"display_name":     "test",
					"parent_domain_id": "invalid uuid",
				},
			),
			ExpectError: regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
		// // error - invalid contributors_scope
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"display_name":       "test",
					"contributors_scope": "invalid value",
				},
			),
			ExpectError: regexp.MustCompile(`Attribute contributors_scope value must be one of`),
		},
	}))
}

func TestUnit_DomainResource_ImportState(t *testing.T) {
	entity := fakes.NewRandomDomain()

	fakes.FakeServer.Upsert(fakes.NewRandomDomain())
	fakes.FakeServer.Upsert(entity)
	fakes.FakeServer.Upsert(fakes.NewRandomDomain())

	testCase := at.CompileConfig(
		testResourceItemHeader,
		map[string]any{
			"display_name": *entity.DisplayName,
			"description":  *entity.Description,
		},
	)

	resource.Test(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		{
			ResourceName:  testResourceItemFQN,
			Config:        testCase,
			ImportStateId: "not-valid",
			ImportState:   true,
			ExpectError:   regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
		// Import state testing
		{
			ResourceName:       testResourceItemFQN,
			Config:             testCase,
			ImportStateId:      *entity.ID,
			ImportState:        true,
			ImportStatePersist: true,
			ImportStateCheck: func(is []*terraform.InstanceState) error {
				if len(is) != 1 {
					return errors.New("expected one instance state")
				}

				if is[0].ID != *entity.ID {
					return errors.New(testResourceItemFQN + ": unexpected ID")
				}

				if is[0].Attributes["display_name"] != *entity.DisplayName {
					return errors.New(testResourceItemFQN + ": unexpected display_name")
				}

				if is[0].Attributes["description"] != *entity.Description {
					return errors.New(testResourceItemFQN + ": unexpected description")
				}

				if is[0].Attributes["contributors_scope"] != string(*entity.ContributorsScope) {
					return errors.New(testResourceItemFQN + ": unexpected contributors_scope")
				}

				return nil
			},
		},
	}))
}

func TestUnit_DomainResource_CRUD(t *testing.T) {
	entityExist := fakes.NewRandomDomain()
	entityBefore := fakes.NewRandomDomain()
	entityAfter := fakes.NewRandomDomainWithContributorsScope(admin.ContributorsScopeTypeAdminsOnly)

	fakes.FakeServer.Upsert(fakes.NewRandomDomain())
	fakes.FakeServer.Upsert(entityExist)
	fakes.FakeServer.Upsert(fakes.NewRandomDomain())

	defaultContributorsScope := string(admin.ContributorsScopeTypeAllTenant)

	resource.Test(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		// error - create - existing entity
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"display_name": *entityExist.DisplayName,
				},
			),
			ExpectError: regexp.MustCompile(common.ErrorCreateHeader),
		},
		// Create and Read
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"display_name": *entityBefore.DisplayName,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "display_name", entityBefore.DisplayName),
				resource.TestCheckResourceAttr(testResourceItemFQN, "description", ""),
				resource.TestCheckResourceAttr(testResourceItemFQN, "contributors_scope", defaultContributorsScope),
			),
		},
		// Update and Read
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"display_name":       *entityBefore.DisplayName,
					"description":        *entityAfter.Description,
					"contributors_scope": string(*entityAfter.ContributorsScope),
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "display_name", entityBefore.DisplayName),
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "description", entityAfter.Description),
				resource.TestCheckResourceAttr(testResourceItemFQN, "contributors_scope", string(*entityAfter.ContributorsScope)),
			),
		},
	}))
}

func TestAcc_DomainResource_CRUD(t *testing.T) {
	entityCreateDisplayName := testhelp.RandomName()
	entityUpdateDisplayName := testhelp.RandomName()
	entityUpdateDescription := testhelp.RandomName()
	defaultContributorsScope := string(admin.ContributorsScopeTypeAllTenant)
	updatedContributorsScope := string(admin.ContributorsScopeTypeAdminsOnly)

	resource.Test(t, testhelp.NewTestAccCase(t, &testResourceItemFQN, nil, []resource.TestStep{
		// Create and Read
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"display_name": entityCreateDisplayName,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "display_name", entityCreateDisplayName),
				resource.TestCheckResourceAttr(testResourceItemFQN, "description", ""),
				resource.TestCheckResourceAttr(testResourceItemFQN, "contributors_scope", defaultContributorsScope),
			),
		},
		// Update and Read
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"display_name":       entityUpdateDisplayName,
					"description":        entityUpdateDescription,
					"contributors_scope": updatedContributorsScope,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "display_name", entityUpdateDisplayName),
				resource.TestCheckResourceAttr(testResourceItemFQN, "description", entityUpdateDescription),
				resource.TestCheckResourceAttr(testResourceItemFQN, "contributors_scope", updatedContributorsScope),
			),
		},
	},
	))
}
