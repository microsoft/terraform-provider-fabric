// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package deploymentpipelinera_test

import (
	"fmt"
	"regexp"
	"testing"

	at "github.com/dcarbone/terraform-plugin-framework-utils/v3/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp/fakes"
)

var testResourceItemFQN, testResourceItemHeader = testhelp.TFResource(common.ProviderTypeName, itemTypeInfo.Type, "test")

func TestUnit_DeploymentPipelineRoleAssignmentResource_Attributes(t *testing.T) {
	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		// error - no required attributes - deployment_pipeline_id
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"principal": map[string]any{
						"id":   "00000000-0000-0000-0000-000000000000",
						"type": "User",
					},
					"role": "Admin",
				},
			),
			ExpectError: regexp.MustCompile(`The argument "deployment_pipeline_id" is required, but no definition was found.`),
		},
		// error - no required attributes - principal.id
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"deployment_pipeline_id": "00000000-0000-0000-0000-000000000000",
					"principal": map[string]any{
						"type": "User",
					},
					"role": "Admin",
				},
			),
			ExpectError: regexp.MustCompile(`Inappropriate value for attribute "principal": attribute "id" is required.`),
		},
		// error - no required attributes - principal.type
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"deployment_pipeline_id": "00000000-0000-0000-0000-000000000000",
					"principal": map[string]any{
						"id": "00000000-0000-0000-0000-000000000000",
					},
					"role": "Admin",
				},
			),
			ExpectError: regexp.MustCompile(`Inappropriate value for attribute "principal": attribute "type" is required.`),
		},
		// error - no required attributes - role
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"deployment_pipeline_id": "00000000-0000-0000-0000-000000000000",
					"principal": map[string]any{
						"id":   "00000000-0000-0000-0000-000000000000",
						"type": "User",
					},
				},
			),
			ExpectError: regexp.MustCompile(`The argument "role" is required, but no definition was found.`),
		},
		// error - invalid UUID - deployment_pipeline_id
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"deployment_pipeline_id": "invalid uuid",
					"principal": map[string]any{
						"id":   "00000000-0000-0000-0000-000000000000",
						"type": "User",
					},
					"role": "Admin",
				},
			),
			ExpectError: regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
		// error - invalid UUID - principal.id
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"deployment_pipeline_id": "00000000-0000-0000-0000-000000000000",
					"principal": map[string]any{
						"id":   "invalid uuid",
						"type": "User",
					},
					"role": "Admin",
				},
			),
			ExpectError: regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
	}))
}

func TestUnit_DeploymentPipelineRoleAssignmentResource_ImportState(t *testing.T) {
	testCase := at.CompileConfig(
		testResourceItemHeader,
		map[string]any{},
	)

	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		{
			ResourceName:  testResourceItemFQN,
			Config:        testCase,
			ImportStateId: "not-valid",
			ImportState:   true,
			ExpectError:   regexp.MustCompile("DeploymentPipelineID/DeploymentPipelineRoleAssignmentID"),
		},
		{
			ResourceName:  testResourceItemFQN,
			Config:        testCase,
			ImportStateId: "test/id",
			ImportState:   true,
			ExpectError:   regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
		{
			ResourceName:  testResourceItemFQN,
			Config:        testCase,
			ImportStateId: fmt.Sprintf("%s/%s", "test", "00000000-0000-0000-0000-000000000000"),
			ImportState:   true,
			ExpectError:   regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
		{
			ResourceName:  testResourceItemFQN,
			Config:        testCase,
			ImportStateId: fmt.Sprintf("%s/%s", "00000000-0000-0000-0000-000000000000", "test"),
			ImportState:   true,
			ExpectError:   regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
	}))
}
