// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package semanticmodelconnectionbinding_test

import (
	"regexp"
	"testing"

	at "github.com/dcarbone/terraform-plugin-framework-utils/v3/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/services/semanticmodelconnectionbinding"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp/fakes"
)

var testDataSourcesFQN, testDataSourcesHeader = testhelp.TFDataSource(
	common.ProviderTypeName,
	semanticmodelconnectionbinding.ItemTypeInfo.Types,
	"test",
)

func TestUnit_SemanticModelConnectionBindingsDataSource_Attributes(t *testing.T) {
	workspaceID := testhelp.RandomUUID()
	semanticModelID := testhelp.RandomUUID()

	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, nil, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		// error - no attributes
		{
			Config:      at.CompileConfig(testDataSourcesHeader, map[string]any{}),
			ExpectError: regexp.MustCompile(`The argument "workspace_id" is required`),
		},
		// error - missing semantic_model_id
		{
			Config: at.CompileConfig(testDataSourcesHeader, map[string]any{
				"workspace_id": workspaceID,
			}),
			ExpectError: regexp.MustCompile(`The argument "semantic_model_id" is required`),
		},
		// error - missing workspace_id
		{
			Config: at.CompileConfig(testDataSourcesHeader, map[string]any{
				"semantic_model_id": semanticModelID,
			}),
			ExpectError: regexp.MustCompile(`The argument "workspace_id" is required`),
		},
		// error - invalid workspace_id
		{
			Config: at.CompileConfig(testDataSourcesHeader, map[string]any{
				"workspace_id":      "invalid uuid",
				"semantic_model_id": semanticModelID,
			}),
			ExpectError: regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
		// error - invalid semantic_model_id
		{
			Config: at.CompileConfig(testDataSourcesHeader, map[string]any{
				"workspace_id":      workspaceID,
				"semantic_model_id": "invalid uuid",
			}),
			ExpectError: regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
		// error - unexpected attribute
		{
			Config: at.CompileConfig(testDataSourcesHeader, map[string]any{
				"workspace_id":      workspaceID,
				"semantic_model_id": semanticModelID,
				"unexpected_attr":   "test",
			}),
			ExpectError: regexp.MustCompile(`An argument named "unexpected_attr" is not expected here`),
		},
	}))
}

func TestUnit_SemanticModelConnectionBindingsDataSource(t *testing.T) {
	workspaceID := testhelp.RandomUUID()
	semanticModelID := testhelp.RandomUUID()
	firstConnectionID := testhelp.RandomUUID()
	secondConnectionID := testhelp.RandomUUID()
	unboundPath := "server;database-unbound"
	emptyWorkspaceID := testhelp.RandomUUID()
	emptySemanticModelID := testhelp.RandomUUID()

	state := newSemanticModelConnectionBindingFakeState()
	state.add(workspaceID, semanticModelID, newRandomItemConnection(firstConnectionID, "server;database-a", "ShareableCloud"))
	state.add(workspaceID, semanticModelID, newRandomItemConnection(secondConnectionID, "server;database-b", "VirtualNetworkGateway"))
	state.add(workspaceID, semanticModelID, newUnboundItemConnection(unboundPath))
	state.configure()

	resource.Test(t, testhelp.NewTestUnitCase(t, nil, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		// read - returns all bindings for the semantic model
		{
			Config: at.CompileConfig(testDataSourcesHeader, map[string]any{
				"workspace_id":      workspaceID,
				"semantic_model_id": semanticModelID,
			}),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testDataSourcesFQN, "workspace_id", workspaceID),
				resource.TestCheckResourceAttr(testDataSourcesFQN, "semantic_model_id", semanticModelID),
				resource.TestCheckResourceAttr(testDataSourcesFQN, "values.#", "3"),
				resource.TestCheckTypeSetElemNestedAttrs(testDataSourcesFQN, "values.*", map[string]string{
					"connection_id":     firstConnectionID,
					"connectivity_type": "ShareableCloud",
				}),
				resource.TestCheckTypeSetElemNestedAttrs(testDataSourcesFQN, "values.*", map[string]string{
					"connectivity_type":       "None",
					"connection_details.path": unboundPath,
					"connection_details.type": "SQL",
				}),
			),
		},
		// read - no connections
		{
			Config: at.CompileConfig(testDataSourcesHeader, map[string]any{
				"workspace_id":      emptyWorkspaceID,
				"semantic_model_id": emptySemanticModelID,
			}),
			Check: resource.TestCheckResourceAttr(testDataSourcesFQN, "values.#", "0"),
		},
		// error - the API fails while listing item connections
		{
			PreConfig: func() {
				fakes.FakeServer.ServerFactory.Core.ItemsServer.NewListItemConnectionsPager = fakeListItemConnectionsError()
			},
			Config: at.CompileConfig(testDataSourcesHeader, map[string]any{
				"workspace_id":      workspaceID,
				"semantic_model_id": semanticModelID,
			}),
			ExpectError: regexp.MustCompile(`list failed`),
		},
	}))
}

func TestAcc_SemanticModelConnectionBindingsDataSource(t *testing.T) {
	workspace := testhelp.WellKnown()["WorkspaceDS"].(map[string]any)
	semanticModel := testhelp.WellKnown()["SemanticModelConnectionBinding"].(map[string]any)
	connectionDetails := semanticModel["connectionDetails"].(map[string]any)
	connectionID := semanticModel["connectionId"].(string)

	resource.ParallelTest(t, testhelp.NewTestAccCase(t, nil, nil, []resource.TestStep{
		{
			Config: at.CompileConfig(testDataSourcesHeader, map[string]any{
				"workspace_id":      workspace["id"].(string),
				"semantic_model_id": semanticModel["id"].(string),
			}),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testDataSourcesFQN, "workspace_id", workspace["id"].(string)),
				resource.TestCheckResourceAttr(testDataSourcesFQN, "semantic_model_id", semanticModel["id"].(string)),
				resource.TestCheckResourceAttr(testDataSourcesFQN, "values.#", "1"),
				resource.TestCheckTypeSetElemNestedAttrs(testDataSourcesFQN, "values.*", map[string]string{
					"connection_id":           connectionID,
					"connectivity_type":       "ShareableCloud",
					"connection_details.path": connectionDetails["path"].(string),
					"connection_details.type": connectionDetails["type"].(string),
				}),
			),
		},
	}))
}
