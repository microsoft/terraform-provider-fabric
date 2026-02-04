// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package domainra_test

import (
	"regexp"
	"testing"

	at "github.com/dcarbone/terraform-plugin-framework-utils/v3/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	fabadmin "github.com/microsoft/fabric-sdk-go/fabric/admin"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/services/domainra"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp/fakes"
)

var testResourceItemsFQN, testResourceItemsHeader = testhelp.TFResource(common.ProviderTypeName, itemTypeInfo.Types, "test")

func TestUnit_DomainRoleAssignmentsResource_Attributes(t *testing.T) {
	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, &testResourceItemsFQN, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		// error - no required attributes - domain_id
		{
			ResourceName: testResourceItemsFQN,
			Config: at.CompileConfig(
				testResourceItemsHeader,
				map[string]any{
					"principals": []map[string]any{},
					"role":       string(domainra.DomainRoleContributors),
				},
			),
			ExpectError: regexp.MustCompile(`The argument "domain_id" is required, but no definition was found.`),
		},
		// error - no required attributes - principals
		{
			ResourceName: testResourceItemsFQN,
			Config: at.CompileConfig(
				testResourceItemsHeader,
				map[string]any{
					"domain_id": "00000000-0000-0000-0000-000000000000",
					"role":      string(domainra.DomainRoleContributors),
				},
			),
			ExpectError: regexp.MustCompile(`The argument "principals" is required, but no definition was found.`),
		},
		// error - no required attributes - role
		{
			ResourceName: testResourceItemsFQN,
			Config: at.CompileConfig(
				testResourceItemsHeader,
				map[string]any{
					"domain_id":  "00000000-0000-0000-0000-000000000000",
					"principals": []map[string]any{},
				},
			),
			ExpectError: regexp.MustCompile(`The argument "role" is required, but no definition was found.`),
		},
		// error - invalid UUID - domain_id
		{
			ResourceName: testResourceItemsFQN,
			Config: at.CompileConfig(
				testResourceItemsHeader,
				map[string]any{
					"domain_id":  "invalid uuid",
					"principals": []map[string]any{},
					"role":       string(domainra.DomainRoleContributors),
				},
			),
			ExpectError: regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
		// error - invalid value - role
		{
			ResourceName: testResourceItemsFQN,
			Config: at.CompileConfig(
				testResourceItemsHeader,
				map[string]any{
					"domain_id":  "00000000-0000-0000-0000-000000000000",
					"principals": []map[string]any{},
					"role":       "invalid role",
				},
			),
			ExpectError: regexp.MustCompile(common.ErrorAttValueMatch),
		},
		// error - no required attributes - principals[0].id
		{
			ResourceName: testResourceItemsFQN,
			Config: at.CompileConfig(
				testResourceItemsHeader,
				map[string]any{
					"domain_id": "00000000-0000-0000-0000-000000000000",
					"role":      string(domainra.DomainRoleContributors),
					"principals": []map[string]any{
						{
							"type": string(fabadmin.PrincipalTypeUser),
						},
					},
				},
			),
			ExpectError: regexp.MustCompile("Inappropriate value for attribute \"principals\": element 0: attribute \"id\" is\nrequired."),
		},
		// error - no required attributes - principals[0].type
		{
			ResourceName: testResourceItemsFQN,
			Config: at.CompileConfig(
				testResourceItemsHeader,
				map[string]any{
					"domain_id": "00000000-0000-0000-0000-000000000000",
					"role":      string(domainra.DomainRoleContributors),
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
			ResourceName: testResourceItemsFQN,
			Config: at.CompileConfig(
				testResourceItemsHeader,
				map[string]any{
					"domain_id": "00000000-0000-0000-0000-000000000000",
					"role":      string(domainra.DomainRoleContributors),
					"principals": []map[string]any{
						{
							"id":   "invalid uuid",
							"type": string(fabadmin.PrincipalTypeUser),
						},
					},
				},
			),
			ExpectError: regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
	}))
}

func TestAcc_DomainRoleAssignmentsResource_CRUD(t *testing.T) {
	domainResourceHCL := at.CompileConfig(
		at.ResourceHeader(testhelp.TypeName(common.ProviderTypeName, "domain"), "test"),
		map[string]any{
			"display_name": testhelp.RandomName(),
		},
	)

	domainResourceFQN := testhelp.ResourceFQN(common.ProviderTypeName, "domain", "test")

	entity := testhelp.WellKnown()["Group"].(map[string]any)
	entityID := entity["id"].(string)
	entityType := entity["type"].(string)

	resource.Test(t, testhelp.NewTestAccCase(t, nil, nil, []resource.TestStep{
		// Create and Read
		{
			ResourceName: testResourceItemsFQN,
			Config: at.JoinConfigs(
				domainResourceHCL,
				at.CompileConfig(
					testResourceItemsHeader,
					map[string]any{
						"domain_id": testhelp.RefByFQN(domainResourceFQN, "id"),
						"role":      string(domainra.DomainRoleContributors),
						"principals": []map[string]any{
							{
								"id":   entityID,
								"type": entityType,
							},
						},
					},
				),
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemsFQN, "principals.0.id", entityID),
				resource.TestCheckResourceAttr(testResourceItemsFQN, "principals.0.type", entityType),
				resource.TestCheckResourceAttr(testResourceItemsFQN, "role", string(domainra.DomainRoleContributors)),
			),
		},
		// Update and Read
		{
			ResourceName: testResourceItemsFQN,
			Config: at.JoinConfigs(
				domainResourceHCL,
				at.CompileConfig(
					testResourceItemsHeader,
					map[string]any{
						"domain_id": testhelp.RefByFQN(domainResourceFQN, "id"),
						"role":      string(domainra.DomainRoleAdmins),
						"principals": []map[string]any{
							{
								"id":   entityID,
								"type": entityType,
							},
						},
					},
				),
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemsFQN, "principals.0.id", entityID),
				resource.TestCheckResourceAttr(testResourceItemsFQN, "principals.0.type", entityType),
				resource.TestCheckResourceAttr(testResourceItemsFQN, "role", string(domainra.DomainRoleAdmins)),
			),
		},
	}))
}
