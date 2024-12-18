// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package domain_test

import (
	"net/http"
	"regexp"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	at "github.com/dcarbone/terraform-plugin-framework-utils/v3/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/microsoft/fabric-sdk-go/fabric/admin"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/services/domain"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp/fakes"
)

var (
	testDataSourceDomainWorkspaceAssignments       = testhelp.DataSourceFQN("fabric", domain.DomainWorkspaceAssignmentsTFName, "test")
	testDataSourceDomainWorkspaceAssignmentsHeader = at.DataSourceHeader(testhelp.TypeName("fabric", domain.DomainWorkspaceAssignmentsTFName), "test")
)

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

		return
	}

	entity := fakes.NewRandomDomain()

	fakes.FakeServer.Upsert(fakes.NewRandomDomain())
	fakes.FakeServer.Upsert(entity)
	fakes.FakeServer.Upsert(fakes.NewRandomDomain())

	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, nil, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		// error - unexpected_attr
		{
			Config: at.CompileConfig(
				testDataSourceDomainWorkspaceAssignmentsHeader,
				map[string]any{
					"unexpected_attr": "test",
				},
			),
			ExpectError: regexp.MustCompile(`An argument named "unexpected_attr" is not expected here`),
		},
		// error - no required attributes - domain_id
		{
			Config: at.CompileConfig(
				testDataSourceDomainWorkspaceAssignmentsHeader,
				map[string]any{},
			),
			ExpectError: regexp.MustCompile(`The argument "domain_id" is required, but no definition was found.`),
		},
		// error - invalid UUID - domain_id
		{
			Config: at.CompileConfig(
				testDataSourceDomainWorkspaceAssignmentsHeader,
				map[string]any{
					"domain_id": "invalid uuid",
				},
			),
			ExpectError: regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
		// read
		{
			Config: at.CompileConfig(
				testDataSourceDomainWorkspaceAssignmentsHeader,
				map[string]any{
					"domain_id": *entity.ID,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrPtr(testDataSourceDomainWorkspaceAssignments, "domain_id", entity.ID),
				resource.TestCheckResourceAttr(testDataSourceDomainWorkspaceAssignments, "values.#", "3"),
				resource.TestCheckResourceAttrPtr(testDataSourceDomainWorkspaceAssignments, "values.0.id", dw1.ID),
				resource.TestCheckResourceAttrPtr(testDataSourceDomainWorkspaceAssignments, "values.0.display_name", dw1.DisplayName),
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
				testDataSourceDomainWorkspaceAssignmentsHeader,
				map[string]any{
					"domain_id": entityID,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testDataSourceDomainWorkspaceAssignments, "domain_id", entityID),
			),
		},
	},
	))
}
