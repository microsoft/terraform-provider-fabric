// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package domainwa_test

import (
	"net/http"
	"regexp"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	at "github.com/dcarbone/terraform-plugin-framework-utils/v3/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/microsoft/fabric-sdk-go/fabric/admin"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp/fakes"
)

var testDataSourceItemsFQN, testDataSourceItemsHeader = testhelp.TFDataSource(common.ProviderTypeName, itemTypeInfo.Types, "test")

func newRandomDomainWorkspace() admin.DomainWorkspace {
	randomName := testhelp.RandomName()
	randomID := testhelp.RandomUUID()

	return admin.DomainWorkspace{
		DisplayName: &randomName,
		ID:          &randomID,
	}
}

func TestUnit_DomainWorkspaceAssignmentsDataSource_Attributes(t *testing.T) {
	dw1 := newRandomDomainWorkspace()
	dw2 := newRandomDomainWorkspace()
	dw3 := newRandomDomainWorkspace()

	fakes.FakeServer.ServerFactory.Admin.DomainsServer.NewListDomainWorkspacesPager = func(_ string, _ *admin.DomainsClientListDomainWorkspacesOptions) (resp fake.PagerResponder[admin.DomainsClientListDomainWorkspacesResponse]) {
		resp.AddPage(http.StatusOK, admin.DomainsClientListDomainWorkspacesResponse{
			DomainWorkspaces: admin.DomainWorkspaces{
				Value: []admin.DomainWorkspace{dw1, dw2, dw3},
			},
		}, nil)

		return resp
	}

	entity := fakes.NewRandomDomain()

	fakes.FakeServer.Upsert(fakes.NewRandomDomain())
	fakes.FakeServer.Upsert(entity)
	fakes.FakeServer.Upsert(fakes.NewRandomDomain())

	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, nil, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		// error - unexpected_attr
		{
			Config: at.CompileConfig(
				testDataSourceItemsHeader,
				map[string]any{
					"unexpected_attr": "test",
				},
			),
			ExpectError: regexp.MustCompile(`An argument named "unexpected_attr" is not expected here`),
		},
		// error - no required attributes - domain_id
		{
			Config: at.CompileConfig(
				testDataSourceItemsHeader,
				map[string]any{},
			),
			ExpectError: regexp.MustCompile(`The argument "domain_id" is required, but no definition was found.`),
		},
		// error - invalid UUID - domain_id
		{
			Config: at.CompileConfig(
				testDataSourceItemsHeader,
				map[string]any{
					"domain_id": "invalid uuid",
				},
			),
			ExpectError: regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
		// read
		{
			Config: at.CompileConfig(
				testDataSourceItemsHeader,
				map[string]any{
					"domain_id": *entity.ID,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrPtr(testDataSourceItemsFQN, "domain_id", entity.ID),
				resource.TestCheckResourceAttr(testDataSourceItemsFQN, "workspace_ids.#", "3"),
			),
		},
	}))
}

func TestAcc_DomainWorkspaceAssignmentsDataSource(t *testing.T) {
	entity := testhelp.WellKnown()["DomainParent"].(map[string]any)
	entityID := entity["id"].(string)

	resource.ParallelTest(t, testhelp.NewTestAccCase(t, nil, nil, []resource.TestStep{
		// read
		{
			Config: at.CompileConfig(
				testDataSourceItemsHeader,
				map[string]any{
					"domain_id": entityID,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testDataSourceItemsFQN, "domain_id", entityID),
			),
		},
	},
	))
}
