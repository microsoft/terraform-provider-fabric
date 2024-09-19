// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package domain_test

import (
	"regexp"
	"testing"

	at "github.com/dcarbone/terraform-plugin-framework-utils/v3/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/services/domain"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp/fakes"
)

var (
	testResourceDomainWorkspaceAssignments       = testhelp.ResourceFQN("fabric", domain.DomainWorkspaceAssignmentsTFName, "test")
	testResourceDomainWorkspaceAssignmentsHeader = at.ResourceHeader(testhelp.TypeName("fabric", domain.DomainWorkspaceAssignmentsTFName), "test")
)

func TestUnit_DomainWorkspaceAssignmentsResource_Attributes(t *testing.T) {
	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, &testResourceDomainWorkspaceAssignments, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		// error - no required attributes - domain_id
		{
			ResourceName: testResourceDomainWorkspaceAssignments,
			Config: at.CompileConfig(
				testResourceDomainWorkspaceAssignmentsHeader,
				map[string]any{
					"workspace_ids": []string{},
				},
			),
			ExpectError: regexp.MustCompile(`The argument "domain_id" is required, but no definition was found.`),
		},
		// error - no required attributes - workspace_ids
		{
			ResourceName: testResourceDomainWorkspaceAssignments,
			Config: at.CompileConfig(
				testResourceDomainWorkspaceAssignmentsHeader,
				map[string]any{
					"domain_id": "00000000-0000-0000-0000-000000000000",
				},
			),
			ExpectError: regexp.MustCompile(`The argument "workspace_ids" is required, but no definition was found.`),
		},
		// error - invalid UUID - domain_id
		{
			ResourceName: testResourceDomainWorkspaceAssignments,
			Config: at.CompileConfig(
				testResourceDomainWorkspaceAssignmentsHeader,
				map[string]any{
					"domain_id":     "invalid uuid",
					"workspace_ids": []string{},
				},
			),
			ExpectError: regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
		// error - invalid UUID - workspace_ids[0]
		{
			ResourceName: testResourceDomainWorkspaceAssignments,
			Config: at.CompileConfig(
				testResourceDomainWorkspaceAssignmentsHeader,
				map[string]any{
					"domain_id": "00000000-0000-0000-0000-000000000000",
					"workspace_ids": []string{
						"invalid uuid",
					},
				},
			),
			ExpectError: regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
	}))
}

func TestAcc_DomainWorkspaceAssignmentsResource_CRUD(t *testing.T) {
	if testhelp.ShouldSkipTest(t) {
		t.Skip("No SPN support")
	}

	domainResourceHCL := at.CompileConfig(
		at.ResourceHeader(testhelp.TypeName("fabric", itemTFName), "test"),
		map[string]any{
			"display_name": testhelp.RandomName(),
		},
	)

	domainResourceFQN := testhelp.ResourceFQN("fabric", itemTFName, "test")

	entity := testhelp.WellKnown().Workspace

	resource.Test(t, testhelp.NewTestAccCase(t, nil, nil, []resource.TestStep{
		// Create and Read
		{
			ResourceName: testResourceDomainWorkspaceAssignments,
			Config: at.JoinConfigs(
				domainResourceHCL,
				at.CompileConfig(
					testResourceDomainWorkspaceAssignmentsHeader,
					map[string]any{
						"domain_id": testhelp.RefByFQN(domainResourceFQN, "id"),
						"workspace_ids": []string{
							*entity.ID,
						},
					},
				),
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceDomainWorkspaceAssignments, "workspace_ids.#", "1"),
				resource.TestCheckResourceAttrPtr(testResourceDomainWorkspaceAssignments, "workspace_ids.0", entity.ID),
			),
		},
		// Update and Read
		{
			ResourceName: testResourceDomainWorkspaceAssignments,
			Config: at.JoinConfigs(
				domainResourceHCL,
				at.CompileConfig(
					testResourceDomainWorkspaceAssignmentsHeader,
					map[string]any{
						"domain_id":     testhelp.RefByFQN(domainResourceFQN, "id"),
						"workspace_ids": []string{},
					},
				),
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceDomainWorkspaceAssignments, "workspace_ids.#", "0"),
			),
		},
	}))
}
