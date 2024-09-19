// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package domain_test

import (
	"regexp"
	"testing"

	at "github.com/dcarbone/terraform-plugin-framework-utils/v3/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/microsoft/fabric-sdk-go/fabric/admin"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/services/domain"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp/fakes"
)

var (
	testResourceDomainRoleAssignments       = testhelp.ResourceFQN("fabric", domain.DomainRoleAssignmentsTFName, "test")
	testResourceDomainRoleAssignmentsHeader = at.ResourceHeader(testhelp.TypeName("fabric", domain.DomainRoleAssignmentsTFName), "test")
)

func TestUnit_DomainRoleAssignmentsResource_Attributes(t *testing.T) {
	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, &testResourceDomainRoleAssignments, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		// error - no required attributes - domain_id
		{
			ResourceName: testResourceDomainRoleAssignments,
			Config: at.CompileConfig(
				testResourceDomainRoleAssignmentsHeader,
				map[string]any{
					"principals": []map[string]any{},
					"role":       string(admin.DomainRoleContributors),
				},
			),
			ExpectError: regexp.MustCompile(`The argument "domain_id" is required, but no definition was found.`),
		},
		// error - no required attributes - principals
		{
			ResourceName: testResourceDomainRoleAssignments,
			Config: at.CompileConfig(
				testResourceDomainRoleAssignmentsHeader,
				map[string]any{
					"domain_id": "00000000-0000-0000-0000-000000000000",
					"role":      string(admin.DomainRoleContributors),
				},
			),
			ExpectError: regexp.MustCompile(`The argument "principals" is required, but no definition was found.`),
		},
		// error - no required attributes - role
		{
			ResourceName: testResourceDomainRoleAssignments,
			Config: at.CompileConfig(
				testResourceDomainRoleAssignmentsHeader,
				map[string]any{
					"domain_id":  "00000000-0000-0000-0000-000000000000",
					"principals": []map[string]any{},
				},
			),
			ExpectError: regexp.MustCompile(`The argument "role" is required, but no definition was found.`),
		},
		// error - invalid UUID - domain_id
		{
			ResourceName: testResourceDomainRoleAssignments,
			Config: at.CompileConfig(
				testResourceDomainRoleAssignmentsHeader,
				map[string]any{
					"domain_id":  "invalid uuid",
					"principals": []map[string]any{},
					"role":       string(admin.DomainRoleContributors),
				},
			),
			ExpectError: regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
		// error - invalid value - role
		{
			ResourceName: testResourceDomainRoleAssignments,
			Config: at.CompileConfig(
				testResourceDomainRoleAssignmentsHeader,
				map[string]any{
					"domain_id":  "00000000-0000-0000-0000-000000000000",
					"principals": []map[string]any{},
					"role":       "invalid role",
				},
			),
			ExpectError: regexp.MustCompile("Invalid Attribute Value Match"),
		},
		// error - no required attributes - principals[0].id
		{
			ResourceName: testResourceDomainRoleAssignments,
			Config: at.CompileConfig(
				testResourceDomainRoleAssignmentsHeader,
				map[string]any{
					"domain_id": "00000000-0000-0000-0000-000000000000",
					"role":      string(admin.DomainRoleContributors),
					"principals": []map[string]any{
						{
							"type": string(admin.PrincipalTypeUser),
						},
					},
				},
			),
			ExpectError: regexp.MustCompile("Inappropriate value for attribute \"principals\": element 0: attribute \"id\" is\nrequired."),
		},
		// error - no required attributes - principals[0].type
		{
			ResourceName: testResourceDomainRoleAssignments,
			Config: at.CompileConfig(
				testResourceDomainRoleAssignmentsHeader,
				map[string]any{
					"domain_id": "00000000-0000-0000-0000-000000000000",
					"role":      string(admin.DomainRoleContributors),
					"principals": []map[string]any{
						{
							"id": "00000000-0000-0000-0000-000000000000",
						},
					},
				},
			),
			ExpectError: regexp.MustCompile("Inappropriate value for attribute \"principals\": element 0: attribute \"type\"\nis required."),
		},
		// error - invalid UUID - principals[0].id
		{
			ResourceName: testResourceDomainRoleAssignments,
			Config: at.CompileConfig(
				testResourceDomainRoleAssignmentsHeader,
				map[string]any{
					"domain_id": "00000000-0000-0000-0000-000000000000",
					"role":      string(admin.DomainRoleContributors),
					"principals": []map[string]any{
						{
							"id":   "invalid uuid",
							"type": string(admin.PrincipalTypeUser),
						},
					},
				},
			),
			ExpectError: regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
	}))
}

// TEST COMMENTED OUT FOR NOW. WHEN A DOMAIN IS CREATED, IT HAS NO ROLE ASSIGNMENTS.
// HOWEVER, AFTER CREATING A ROLE ASSIGNMENT, TRYING TO DELETE IT WILL RESULT IN AN ERROR.
// IN ITS CURRENT STATE, THE TEST WILL FAIL WHEN RUNNING THE DESTROY OPERATION AFTER THE PLAN.

// func TestAcc_DomainRoleAssignmentsResource_CRUD(t *testing.T) {
//  if testhelp.ShouldSkipTest(t) {
// 	 t.Skip("No SPN support")
//  }
// 	domainResourceHCL := at.CompileConfig(
// 		at.ResourceHeader(testhelp.TypeName("fabric", itemTFName), "test"),
// 		map[string]any{
// 			"display_name":       testhelp.RandomName(),
// 			"contributors_scope": string(admin.ContributorsScopeTypeSpecificUsersAndGroups),
// 		},
// 	)

// 	domainResourceFQN := testhelp.ResourceFQN("fabric", itemTFName, "test")

// 	entity := testhelp.WellKnown().Group

//nolint:dupword
// 	resource.Test(t, testhelp.NewTestAccCase(t, nil, nil, []resource.TestStep{
// 		// Create and Read
// 		{
// 			ResourceName: testResourceDomainRoleAssignments,
// 			Config: at.JoinConfigs(
// 				domainResourceHCL,
// 				at.CompileConfig(
// 					testResourceDomainRoleAssignmentsHeader,
// 					map[string]any{
// 						"domain_id": testhelp.RefByFQN(domainResourceFQN, "id"),
// 						"role":      string(admin.DomainRoleContributors),
// 						"principals": []map[string]any{
// 							{
// 								"id":   *entity.ID,
// 								"type": *entity.Type,
// 							},
// 						},
// 					},
// 				),
// 			),
// 			Check: resource.ComposeAggregateTestCheckFunc(
// 				resource.TestCheckResourceAttrPtr(testResourceDomainRoleAssignments, "principals.0.id", entity.ID),
// 				resource.TestCheckResourceAttrPtr(testResourceDomainRoleAssignments, "principals.0.type", entity.Type),
// 				resource.TestCheckResourceAttr(testResourceDomainRoleAssignments, "role", string(admin.DomainRoleContributors)),
// 			),
// 		},
// 		// Update and Read
// 		{
// 			ResourceName: testResourceDomainRoleAssignments,
// 			Config: at.JoinConfigs(
// 				domainResourceHCL,
// 				at.CompileConfig(
// 					testResourceDomainRoleAssignmentsHeader,
// 					map[string]any{
// 						"domain_id": testhelp.RefByFQN(domainResourceFQN, "id"),
// 						"role":      string(admin.DomainRoleAdmins),
// 						"principals": []map[string]any{
// 							{
// 								"id":   *entity.ID,
// 								"type": *entity.Type,
// 							},
// 						},
// 					},
// 				),
// 			),
// 			Check: resource.ComposeAggregateTestCheckFunc(
// 				resource.TestCheckResourceAttrPtr(testResourceDomainRoleAssignments, "principals.0.id", entity.ID),
// 				resource.TestCheckResourceAttrPtr(testResourceDomainRoleAssignments, "principals.0.type", entity.Type),
// 				resource.TestCheckResourceAttr(testResourceDomainRoleAssignments, "role", string(admin.DomainRoleAdmins)),
// 			),
// 		},
// 	}))
// }
