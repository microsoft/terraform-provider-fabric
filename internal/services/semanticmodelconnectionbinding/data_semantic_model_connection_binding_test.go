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

var testDataSourceFQN, testDataSourceHeader = testhelp.TFDataSource(
	common.ProviderTypeName,
	semanticmodelconnectionbinding.ItemTypeInfo.Type,
	"test",
)

func TestUnit_SemanticModelConnectionBindingDataSource_Attributes(t *testing.T) {
	workspaceID := testhelp.RandomUUID()
	semanticModelID := testhelp.RandomUUID()

	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, nil, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		// error - no attributes
		{
			Config:      at.CompileConfig(testDataSourceHeader, map[string]any{}),
			ExpectError: regexp.MustCompile(`The argument "workspace_id" is required`),
		},
		// error - missing semantic_model_id
		{
			Config: at.CompileConfig(testDataSourceHeader, map[string]any{
				"workspace_id": workspaceID,
			}),
			ExpectError: regexp.MustCompile(`The argument "semantic_model_id" is required`),
		},
		// error - missing workspace_id
		{
			Config: at.CompileConfig(testDataSourceHeader, map[string]any{
				"semantic_model_id": semanticModelID,
			}),
			ExpectError: regexp.MustCompile(`The argument "workspace_id" is required`),
		},
		// error - invalid workspace_id
		{
			Config: at.CompileConfig(testDataSourceHeader, map[string]any{
				"workspace_id":      "invalid uuid",
				"semantic_model_id": semanticModelID,
				"connection_details": map[string]any{
					"path": "server;database",
					"type": "SQL",
				},
			}),
			ExpectError: regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
		// error - invalid semantic_model_id
		{
			Config: at.CompileConfig(testDataSourceHeader, map[string]any{
				"workspace_id":      workspaceID,
				"semantic_model_id": "invalid uuid",
				"connection_details": map[string]any{
					"path": "server;database",
					"type": "SQL",
				},
			}),
			ExpectError: regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
		// error - missing connection_details
		{
			Config: at.CompileConfig(testDataSourceHeader, map[string]any{
				"workspace_id":      workspaceID,
				"semantic_model_id": semanticModelID,
			}),
			ExpectError: regexp.MustCompile(`The argument "connection_details" is required`),
		},
		// error - missing connection_details.type
		{
			Config: at.CompileConfig(testDataSourceHeader, map[string]any{
				"workspace_id":      workspaceID,
				"semantic_model_id": semanticModelID,
				"connection_details": map[string]any{
					"path": "server;database",
				},
			}),
			ExpectError: regexp.MustCompile(`(?s)attribute "type" is\s+required`),
		},
		// error - missing connection_details.path
		{
			Config: at.CompileConfig(testDataSourceHeader, map[string]any{
				"workspace_id":      workspaceID,
				"semantic_model_id": semanticModelID,
				"connection_details": map[string]any{
					"type": "SQL",
				},
			}),
			ExpectError: regexp.MustCompile(`(?s)attribute "path" is\s+required`),
		},
		// error - unexpected attribute
		{
			Config: at.CompileConfig(testDataSourceHeader, map[string]any{
				"workspace_id":      workspaceID,
				"semantic_model_id": semanticModelID,
				"connection_details": map[string]any{
					"path": "server;database",
					"type": "SQL",
				},
				"unexpected_attr": "test",
			}),
			ExpectError: regexp.MustCompile(`An argument named "unexpected_attr" is not expected here`),
		},
	}))
}

func TestUnit_SemanticModelConnectionBindingDataSource(t *testing.T) {
	workspaceID := testhelp.RandomUUID()
	semanticModelID := testhelp.RandomUUID()
	connectionID := testhelp.RandomUUID()
	connectionPath := "server;database"
	connectionType := "SQL"
	unboundPath := "server;database-unbound"

	state := newSemanticModelConnectionBindingFakeState()
	state.add(workspaceID, semanticModelID, newRandomItemConnection(connectionID, connectionPath, "ShareableCloud"))
	state.add(workspaceID, semanticModelID, newUnboundItemConnection(unboundPath))
	state.configure()

	resource.Test(t, testhelp.NewTestUnitCase(t, nil, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		// read - finds binding by connection details
		{
			Config: at.CompileConfig(testDataSourceHeader, map[string]any{
				"workspace_id":      workspaceID,
				"semantic_model_id": semanticModelID,
				"connection_details": map[string]any{
					"path": connectionPath,
					"type": connectionType,
				},
			}),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testDataSourceFQN, "workspace_id", workspaceID),
				resource.TestCheckResourceAttr(testDataSourceFQN, "semantic_model_id", semanticModelID),
				resource.TestCheckResourceAttr(testDataSourceFQN, "connection_id", connectionID),
				resource.TestCheckResourceAttr(testDataSourceFQN, "connectivity_type", "ShareableCloud"),
				resource.TestCheckResourceAttr(testDataSourceFQN, "connection_details.path", connectionPath),
				resource.TestCheckResourceAttr(testDataSourceFQN, "connection_details.type", connectionType),
			),
		},
		// error - no matching connection details
		{
			Config: at.CompileConfig(testDataSourceHeader, map[string]any{
				"workspace_id":      workspaceID,
				"semantic_model_id": semanticModelID,
				"connection_details": map[string]any{
					"path": "missing",
					"type": connectionType,
				},
			}),
			ExpectError: regexp.MustCompile(`No binding exists`),
		},
		// read - the data source reference exists but is unbound
		{
			Config: at.CompileConfig(testDataSourceHeader, map[string]any{
				"workspace_id":      workspaceID,
				"semantic_model_id": semanticModelID,
				"connection_details": map[string]any{
					"path": unboundPath,
					"type": connectionType,
				},
			}),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testDataSourceFQN, "connectivity_type", "None"),
				resource.TestCheckNoResourceAttr(testDataSourceFQN, "connection_id"),
			),
		},
		// error - the API fails while listing item connections
		{
			PreConfig: func() {
				fakes.FakeServer.ServerFactory.Core.ItemsServer.NewListItemConnectionsPager = fakeListItemConnectionsError()
			},
			Config: at.CompileConfig(testDataSourceHeader, map[string]any{
				"workspace_id":      workspaceID,
				"semantic_model_id": semanticModelID,
				"connection_details": map[string]any{
					"path": connectionPath,
					"type": connectionType,
				},
			}),
			ExpectError: regexp.MustCompile(`list failed`),
		},
	}))
}

func TestAcc_SemanticModelConnectionBindingDataSource(t *testing.T) {
	workspace := testhelp.WellKnown()["WorkspaceDS"].(map[string]any)
	semanticModel := testhelp.WellKnown()["SemanticModelConnectionBinding"].(map[string]any)
	connectionDetails := semanticModel["connectionDetails"].(map[string]any)

	resource.ParallelTest(t, testhelp.NewTestAccCase(t, nil, nil, []resource.TestStep{
		{
			Config: at.CompileConfig(testDataSourceHeader, map[string]any{
				"workspace_id":      workspace["id"].(string),
				"semantic_model_id": semanticModel["id"].(string),
				"connection_details": map[string]any{
					"path": connectionDetails["path"],
					"type": connectionDetails["type"],
				},
			}),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testDataSourceFQN, "connection_id", semanticModel["connectionId"].(string)),
				resource.TestCheckResourceAttr(testDataSourceFQN, "connectivity_type", "ShareableCloud"),
				resource.TestCheckResourceAttr(testDataSourceFQN, "connection_details.path", connectionDetails["path"].(string)),
				resource.TestCheckResourceAttr(testDataSourceFQN, "connection_details.type", connectionDetails["type"].(string)),
			),
		},
	}))
}
