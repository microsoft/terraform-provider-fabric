// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package domainwa_test

import (
	"regexp"
	"testing"

	at "github.com/dcarbone/terraform-plugin-framework-utils/v3/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp/fakes"
)

var testResourceItemsFQN, testResourceItemsHeader = testhelp.TFResource(common.ProviderTypeName, itemTypeInfo.Types, "test")

func TestUnit_DomainWorkspaceAssignmentsResource_Attributes(t *testing.T) {
	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, &testResourceItemsFQN, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		// error - no required attributes - domain_id
		{
			ResourceName: testResourceItemsFQN,
			Config: at.CompileConfig(
				testResourceItemsHeader,
				map[string]any{
					"workspace_ids": []string{},
				},
			),
			ExpectError: regexp.MustCompile(`The argument "domain_id" is required, but no definition was found.`),
		},
		// error - no required attributes - workspace_ids
		{
			ResourceName: testResourceItemsFQN,
			Config: at.CompileConfig(
				testResourceItemsHeader,
				map[string]any{
					"domain_id": "00000000-0000-0000-0000-000000000000",
				},
			),
			ExpectError: regexp.MustCompile(`The argument "workspace_ids" is required, but no definition was found.`),
		},
		// error - invalid UUID - domain_id
		{
			ResourceName: testResourceItemsFQN,
			Config: at.CompileConfig(
				testResourceItemsHeader,
				map[string]any{
					"domain_id":     "invalid uuid",
					"workspace_ids": []string{},
				},
			),
			ExpectError: regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
		// error - invalid UUID - workspace_ids[0]
		{
			ResourceName: testResourceItemsFQN,
			Config: at.CompileConfig(
				testResourceItemsHeader,
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
	domainResourceHCL := at.CompileConfig(
		at.ResourceHeader(testhelp.TypeName("fabric", "domain"), "test"),
		map[string]any{
			"display_name": testhelp.RandomName(),
		},
	)
	domainResourceFQN := testhelp.ResourceFQN("fabric", "domain", "test")

	workspace1ResourceHCL := at.CompileConfig(
		at.ResourceHeader(testhelp.TypeName("fabric", "workspace"), "test1"),
		map[string]any{
			"display_name": testhelp.RandomName(),
		},
	)
	workspace1ResourceFQN := testhelp.ResourceFQN("fabric", "workspace", "test1")

	workspace2ResourceHCL := at.CompileConfig(
		at.ResourceHeader(testhelp.TypeName("fabric", "workspace"), "test2"),
		map[string]any{
			"display_name": testhelp.RandomName(),
		},
	)
	workspace2ResourceFQN := testhelp.ResourceFQN("fabric", "workspace", "test2")

	// entity1 := testhelp.WellKnown()["WorkspaceRS"].(map[string]any)
	// entity1ID := entity1["id"].(string)

	// entity2 := testhelp.WellKnown()["WorkspaceDS"].(map[string]any)
	// entity2ID := entity2["id"].(string)

	resource.Test(t, testhelp.NewTestAccCase(t, nil, nil, []resource.TestStep{
		// Create and Read
		{
			ResourceName: testResourceItemsFQN,
			Config: at.JoinConfigs(
				domainResourceHCL,
				workspace1ResourceHCL,
				workspace2ResourceHCL,
				at.CompileConfig(
					testResourceItemsHeader,
					map[string]any{
						"domain_id": testhelp.RefByFQN(domainResourceFQN, "id"),
						"workspace_ids": []string{
							testhelp.RefByFQN(workspace1ResourceFQN, "id"),
						},
					},
				),
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemsFQN, "workspace_ids.#", "1"),
				// resource.TestCheckResourceAttr(testResourceItemsFQN, "workspace_ids.0", entity1ID),
			),
		},
		// Update and Read
		{
			ResourceName: testResourceItemsFQN,
			Config: at.JoinConfigs(
				domainResourceHCL,
				workspace1ResourceHCL,
				workspace2ResourceHCL,
				at.CompileConfig(
					testResourceItemsHeader,
					map[string]any{
						"domain_id": testhelp.RefByFQN(domainResourceFQN, "id"),
						"workspace_ids": []string{
							testhelp.RefByFQN(workspace1ResourceFQN, "id"),
							testhelp.RefByFQN(workspace2ResourceFQN, "id"),
						},
					},
				),
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemsFQN, "workspace_ids.#", "2"),
			),
		},
		// Update and Read
		{
			ResourceName: testResourceItemsFQN,
			Config: at.JoinConfigs(
				domainResourceHCL,
				workspace1ResourceHCL,
				workspace2ResourceHCL,
				at.CompileConfig(
					testResourceItemsHeader,
					map[string]any{
						"domain_id": testhelp.RefByFQN(domainResourceFQN, "id"),
						"workspace_ids": []string{
							testhelp.RefByFQN(workspace2ResourceFQN, "id"),
						},
					},
				),
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemsFQN, "workspace_ids.#", "1"),
				// resource.TestCheckResourceAttr(testResourceItemsFQN, "workspace_ids.0", entity2ID),
			),
		},
	}))
}
