// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package semanticmodelconnectionbinding_test

import (
	"encoding/base64"
	"errors"
	"fmt"
	"regexp"
	"testing"

	at "github.com/dcarbone/terraform-plugin-framework-utils/v3/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/services/semanticmodelconnectionbinding"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp/fakes"
)

var testResourceFQN, testResourceHeader = testhelp.TFResource(
	common.ProviderTypeName,
	semanticmodelconnectionbinding.ItemTypeInfo.Type,
	"test",
)

func TestUnit_SemanticModelConnectionBindingResource_Attributes(t *testing.T) {
	workspaceID := testhelp.RandomUUID()
	semanticModelID := testhelp.RandomUUID()

	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, &testResourceFQN, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		// error - no attributes
		{
			ResourceName: testResourceFQN,
			Config:       at.CompileConfig(testResourceHeader, map[string]any{}),
			ExpectError:  regexp.MustCompile(`Missing required argument`),
		},
		// error - invalid workspace UUID
		{
			ResourceName: testResourceFQN,
			Config: at.CompileConfig(testResourceHeader, map[string]any{
				"workspace_id":      "invalid uuid",
				"semantic_model_id": semanticModelID,
				"connectivity_type": "ShareableCloud",
				"connection_id":     testhelp.RandomUUID(),
				"connection_details": map[string]any{
					"path": "server;database",
					"type": "SQL",
				},
			}),
			ExpectError: regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
		// error - Automatic connectivity cannot use a connection ID
		{
			ResourceName: testResourceFQN,
			Config: at.CompileConfig(testResourceHeader, map[string]any{
				"workspace_id":      workspaceID,
				"semantic_model_id": semanticModelID,
				"connectivity_type": "Automatic",
				"connection_id":     testhelp.RandomUUID(),
				"connection_details": map[string]any{
					"path": "server;database",
					"type": "SQL",
				},
			}),
			ExpectError: regexp.MustCompile(
				`(?s)Invalid configuration for attribute connection_id.*If connectivity_type attribute is set and the value is "Automatic" this\s+attribute is NULL`,
			),
		},
		// error - missing connection_details.type
		{
			ResourceName: testResourceFQN,
			Config: at.CompileConfig(testResourceHeader, map[string]any{
				"workspace_id":      workspaceID,
				"semantic_model_id": semanticModelID,
				"connectivity_type": "ShareableCloud",
				"connection_details": map[string]any{
					"path": "server;database",
				},
			}),
			ExpectError: regexp.MustCompile(`(?s)attribute "type" is\s+required`),
		},
	}))
}

func TestUnit_SemanticModelConnectionBindingResource_ImportState(t *testing.T) {
	workspaceID := testhelp.RandomUUID()
	semanticModelID := testhelp.RandomUUID()
	notFoundWorkspaceID := testhelp.RandomUUID()
	notFoundSemanticModelID := testhelp.RandomUUID()
	unboundWorkspaceID := testhelp.RandomUUID()
	unboundSemanticModelID := testhelp.RandomUUID()
	state := newSemanticModelConnectionBindingFakeState()
	state.add(unboundWorkspaceID, unboundSemanticModelID, newRandomItemConnection("", "server;database", "None"))
	state.configure()

	resource.Test(t, testhelp.NewTestUnitCase(t, &testResourceFQN, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		// error - invalid identifier format
		{
			ResourceName:  testResourceFQN,
			Config:        at.CompileConfig(testResourceHeader, map[string]any{}),
			ImportStateId: "not-valid",
			ImportState:   true,
			ExpectError:   regexp.MustCompile("WorkspaceID/SemanticModelID/Base64URLPath/Base64URLType"),
		},
		// error - invalid base64url path
		{
			ResourceName:  testResourceFQN,
			Config:        at.CompileConfig(testResourceHeader, map[string]any{}),
			ImportStateId: fmt.Sprintf("%s/%s/%%/U1FM", workspaceID, semanticModelID),
			ImportState:   true,
			ExpectError:   regexp.MustCompile(`(?s)decode connection\s+path`),
		},
		// error - invalid workspace UUID
		{
			ResourceName:  testResourceFQN,
			Config:        at.CompileConfig(testResourceHeader, map[string]any{}),
			ImportStateId: bindingImportID("invalid", semanticModelID, "server;database", "SQL"),
			ImportState:   true,
			ExpectError:   regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
		// error - binding not found
		{
			ResourceName:  testResourceFQN,
			Config:        at.CompileConfig(testResourceHeader, map[string]any{}),
			ImportStateId: bindingImportID(notFoundWorkspaceID, notFoundSemanticModelID, "server;database", "SQL"),
			ImportState:   true,
			ExpectError:   regexp.MustCompile(`No binding exists`),
		},
		// error - binding is explicitly unbound
		{
			ResourceName:  testResourceFQN,
			Config:        at.CompileConfig(testResourceHeader, map[string]any{}),
			ImportStateId: bindingImportID(unboundWorkspaceID, unboundSemanticModelID, "server;database", "SQL"),
			ImportState:   true,
			ExpectError:   regexp.MustCompile(`No binding exists or is currently unbound`),
		},
	}))
}

func TestUnit_SemanticModelConnectionBindingResource_CRUD(t *testing.T) {
	workspaceID := testhelp.RandomUUID()
	semanticModelID := testhelp.RandomUUID()
	connectionID := testhelp.RandomUUID()
	updatedConnectionID := testhelp.RandomUUID()
	connectionPath := "server;database"
	connectionType := "SQL"

	state := newSemanticModelConnectionBindingFakeState()
	state.configure()

	config := func(connectivityType, connectionID string, includeConnectionID bool) string {
		attributes := map[string]any{
			"workspace_id":      workspaceID,
			"semantic_model_id": semanticModelID,
			"connectivity_type": connectivityType,
			"connection_details": map[string]any{
				"path": connectionPath,
				"type": connectionType,
			},
		}
		if includeConnectionID {
			attributes["connection_id"] = connectionID
		}

		return at.CompileConfig(testResourceHeader, attributes)
	}

	importID := bindingImportID(workspaceID, semanticModelID, connectionPath, connectionType)

	resource.Test(t, testhelp.NewTestUnitCase(t, &testResourceFQN, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		// create and read
		{
			ResourceName: testResourceFQN,
			Config:       config("ShareableCloud", connectionID, true),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceFQN, "connection_id", connectionID),
				resource.TestCheckResourceAttr(testResourceFQN, "connectivity_type", "ShareableCloud"),
				resource.TestCheckResourceAttr(testResourceFQN, "connection_details.path", connectionPath),
				resource.TestCheckResourceAttr(testResourceFQN, "connection_details.type", connectionType),
				resource.TestCheckResourceAttrSet(testResourceFQN, "id"),
			),
		},
		// Remove the remote binding and verify Read removes the managed resource.
		{
			PreConfig: func() {
				state.unbind(workspaceID, semanticModelID, connectionPath, connectionType)
			},
			Config: at.CompileConfig(testDataSourceHeader, map[string]any{
				"workspace_id":      workspaceID,
				"semantic_model_id": semanticModelID,
				"connection_details": map[string]any{
					"path": connectionPath,
					"type": connectionType,
				},
			}),
			Check: func(state *terraform.State) error {
				if _, exists := state.RootModule().Resources[testResourceFQN]; exists {
					return errors.New("expected the unbound resource to be removed from state")
				}

				checkDataSource := resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(testDataSourceFQN, "connectivity_type", "None"),
					resource.TestCheckNoResourceAttr(testDataSourceFQN, "connection_id"),
				)

				return checkDataSource(state)
			},
		},
		// Recreate the binding after refresh removed it from state.
		{
			ResourceName: testResourceFQN,
			Config:       config("ShareableCloud", connectionID, true),
			Check:        resource.TestCheckResourceAttr(testResourceFQN, "connectivity_type", "ShareableCloud"),
		},
		// update connection ID and connectivity type
		{
			ResourceName: testResourceFQN,
			Config:       config("VirtualNetworkGateway", updatedConnectionID, true),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceFQN, "connection_id", updatedConnectionID),
				resource.TestCheckResourceAttr(testResourceFQN, "connectivity_type", "VirtualNetworkGateway"),
			),
		},
		// update to automatic connectivity without a connection ID
		{
			ResourceName: testResourceFQN,
			Config:       config("Automatic", "", false),
			Check:        resource.TestCheckResourceAttr(testResourceFQN, "connectivity_type", "Automatic"),
		},
		// import the binding
		{
			ResourceName:            testResourceFQN,
			Config:                  config("Automatic", "", false),
			ImportStateId:           importID,
			ImportState:             true,
			ImportStateVerify:       true,
			ImportStateVerifyIgnore: []string{"timeouts", "display_name"},
		},
	}))

	if _, ok := state.bindings[bindingKey(workspaceID, semanticModelID, connectionPath, connectionType)]; ok {
		t.Fatal("expected delete to remove the binding from the fake server")
	}
}

func TestUnit_SemanticModelConnectionBindingResource_ErrorPaths(t *testing.T) {
	state := newSemanticModelConnectionBindingFakeState()
	state.configure()

	fakes.FakeServer.ServerFactory.SemanticModel.ItemsServer.BindSemanticModelConnection = fakeBindSemanticModelConnectionError()

	resource.Test(t, testhelp.NewTestUnitCase(t, &testResourceFQN, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		// error - None is reserved for Delete
		{
			ResourceName: testResourceFQN,
			Config: at.CompileConfig(testResourceHeader, map[string]any{
				"workspace_id":      testhelp.RandomUUID(),
				"semantic_model_id": testhelp.RandomUUID(),
				"connectivity_type": "None",
				"connection_details": map[string]any{
					"path": "server;database",
					"type": "SQL",
				},
			}),
			ExpectError: regexp.MustCompile(`value must be one of:`),
		},
		// error - BindSemanticModelConnection API error
		{
			ResourceName: testResourceFQN,
			Config: at.CompileConfig(testResourceHeader, map[string]any{
				"workspace_id":      testhelp.RandomUUID(),
				"semantic_model_id": testhelp.RandomUUID(),
				"connectivity_type": "ShareableCloud",
				"connection_id":     testhelp.RandomUUID(),
				"connection_details": map[string]any{
					"path": "server;database",
					"type": "SQL",
				},
			}),
			ExpectError: regexp.MustCompile(`Bind failed`),
		},
	}))
}

func TestAcc_SemanticModelConnectionBindingResource_CRUD(t *testing.T) {
	workspace := testhelp.WellKnown()["WorkspaceRS"].(map[string]any)
	semanticModel := testhelp.WellKnown()["SemanticModelRS"].(map[string]any)
	shareableCloudConnection := testhelp.WellKnown()["ShareableCloudConnection"].(map[string]any)
	virtualNetworkGatewayConnection := testhelp.WellKnown()["VirtualNetworkGatewayConnection"].(map[string]any)
	connectionDetails := semanticModel["connectionDetails"].(map[string]any)
	importID := bindingImportID(
		workspace["id"].(string),
		semanticModel["id"].(string),
		connectionDetails["path"].(string),
		connectionDetails["type"].(string),
	)

	config := func(connectionID, connectivityType string) string {
		return at.CompileConfig(testResourceHeader, map[string]any{
			"workspace_id":      workspace["id"].(string),
			"semantic_model_id": semanticModel["id"].(string),
			"connection_id":     connectionID,
			"connectivity_type": connectivityType,
			"connection_details": map[string]any{
				"path": connectionDetails["path"],
				"type": connectionDetails["type"],
			},
		})
	}

	resource.ParallelTest(t, testhelp.NewTestAccCase(t, &testResourceFQN, nil, []resource.TestStep{
		// Create and read the WorkspaceRS semantic model binding.
		{
			ResourceName: testResourceFQN,
			Config:       config(shareableCloudConnection["id"].(string), "ShareableCloud"),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceFQN, "connection_id", shareableCloudConnection["id"].(string)),
				resource.TestCheckResourceAttr(testResourceFQN, "connectivity_type", "ShareableCloud"),
				resource.TestCheckResourceAttr(testResourceFQN, "connection_details.path", connectionDetails["path"].(string)),
				resource.TestCheckResourceAttr(testResourceFQN, "connection_details.type", connectionDetails["type"].(string)),
			),
		},
		// Update the source binding to the existing Virtual Network Gateway connection.
		{
			ResourceName: testResourceFQN,
			Config:       config(virtualNetworkGatewayConnection["id"].(string), "VirtualNetworkGateway"),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceFQN, "connection_id", virtualNetworkGatewayConnection["id"].(string)),
				resource.TestCheckResourceAttr(testResourceFQN, "connectivity_type", "VirtualNetworkGateway"),
			),
		},
		// Import the binding using its natural key; final cleanup exercises unbind.
		{
			ResourceName:            testResourceFQN,
			Config:                  config(virtualNetworkGatewayConnection["id"].(string), "VirtualNetworkGateway"),
			ImportStateId:           importID,
			ImportState:             true,
			ImportStateVerify:       true,
			ImportStateVerifyIgnore: []string{"timeouts", "display_name", "gateway_id"},
		},
	}))
}

func bindingImportID(workspaceID, semanticModelID, connectionPath, connectionType string) string {
	encodedPath := base64.RawURLEncoding.EncodeToString([]byte(connectionPath))
	encodedType := base64.RawURLEncoding.EncodeToString([]byte(connectionType))

	return fmt.Sprintf("%s/%s/%s/%s", workspaceID, semanticModelID, encodedPath, encodedType)
}
