// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package semanticmodelcb_test

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

var (
	testResourceItemFQN    = testhelp.ResourceFQN(common.ProviderTypeName, itemTypeInfo.Type, "test")
	testResourceItemHeader = at.ResourceHeader(testhelp.TypeName(common.ProviderTypeName, itemTypeInfo.Type), "test")
)

func TestUnit_SemanticModelConnectionBindingResource_Attributes(t *testing.T) {
	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		// error - missing workspace_id
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"semantic_model_id": "00000000-0000-0000-0000-000000000000",
					"connectivity_type": "ShareableCloud",
					"connection_id":     "00000000-0000-0000-0000-000000000000",
					"connection_details": map[string]any{
						"path": "https://example.com",
						"type": "Sql",
					},
				},
			),
			ExpectError: regexp.MustCompile(`The argument "workspace_id" is required, but no definition was found.`),
		},
		// error - missing semantic_model_id
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id":      "00000000-0000-0000-0000-000000000000",
					"connectivity_type": "ShareableCloud",
					"connection_id":     "00000000-0000-0000-0000-000000000000",
					"connection_details": map[string]any{
						"path": "https://example.com",
						"type": "Sql",
					},
				},
			),
			ExpectError: regexp.MustCompile(`The argument "semantic_model_id" is required, but no definition was found.`),
		},
		// error - missing connectivity_type
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id":      "00000000-0000-0000-0000-000000000000",
					"semantic_model_id": "00000000-0000-0000-0000-000000000000",
					"connection_id":     "00000000-0000-0000-0000-000000000000",
					"connection_details": map[string]any{
						"path": "https://example.com",
						"type": "Sql",
					},
				},
			),
			ExpectError: regexp.MustCompile(`The argument "connectivity_type" is required, but no definition was found.`),
		},
		// error - missing connection_details
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id":      "00000000-0000-0000-0000-000000000000",
					"semantic_model_id": "00000000-0000-0000-0000-000000000000",
					"connectivity_type": "ShareableCloud",
					"connection_id":     "00000000-0000-0000-0000-000000000000",
				},
			),
			ExpectError: regexp.MustCompile(`The argument "connection_details" is required, but no definition was found.`),
		},
		// error - invalid connectivity_type
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id":      "00000000-0000-0000-0000-000000000000",
					"semantic_model_id": "00000000-0000-0000-0000-000000000000",
					"connectivity_type": "Invalid",
					"connection_id":     "00000000-0000-0000-0000-000000000000",
					"connection_details": map[string]any{
						"path": "https://example.com",
						"type": "Sql",
					},
				},
			),
			ExpectError: regexp.MustCompile(`Attribute connectivity_type value must be one of`),
		},
		// error - invalid UUID - workspace_id
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id":      "invalid uuid",
					"semantic_model_id": "00000000-0000-0000-0000-000000000000",
					"connectivity_type": "ShareableCloud",
					"connection_id":     "00000000-0000-0000-0000-000000000000",
					"connection_details": map[string]any{
						"path": "https://example.com",
						"type": "Sql",
					},
				},
			),
			ExpectError: regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
	}))
}

func TestUnit_SemanticModelConnectionBindingResource_CRUD(t *testing.T) {
	workspaceID := testhelp.RandomUUID()
	semanticModelID := testhelp.RandomUUID()
	connectionID := testhelp.RandomUUID()
	updatedConnectionID := testhelp.RandomUUID()

	fakes.FakeServer.ServerFactory.SemanticModel.ItemsServer.BindSemanticModelConnection = fakeBindSemanticModelConnection()

	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		// Create and Read
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id":      workspaceID,
					"semantic_model_id": semanticModelID,
					"connectivity_type": "ShareableCloud",
					"connection_id":     connectionID,
					"connection_details": map[string]any{
						"path": "https://example.com;sales",
						"type": "Sql",
					},
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "workspace_id", workspaceID),
				resource.TestCheckResourceAttr(testResourceItemFQN, "semantic_model_id", semanticModelID),
				resource.TestCheckResourceAttr(testResourceItemFQN, "connectivity_type", "ShareableCloud"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "connection_id", connectionID),
				resource.TestCheckResourceAttr(testResourceItemFQN, "connection_details.path", "https://example.com;sales"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "connection_details.type", "Sql"),
			),
		},
		// Update connection_id (in-place) and Read
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id":      workspaceID,
					"semantic_model_id": semanticModelID,
					"connectivity_type": "ShareableCloud",
					"connection_id":     updatedConnectionID,
					"connection_details": map[string]any{
						"path": "https://example.com;sales",
						"type": "Sql",
					},
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "connection_id", updatedConnectionID),
			),
		},
	}))
}
