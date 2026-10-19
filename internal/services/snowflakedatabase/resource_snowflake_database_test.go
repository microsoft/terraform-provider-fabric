// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package snowflakedatabase_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"testing"

	at "github.com/dcarbone/terraform-plugin-framework-utils/v3/acctest"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/params"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/transforms"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp/fakes"
)

var testResourceItemFQN, testResourceItemHeader = testhelp.TFResource(common.ProviderTypeName, itemTypeInfo.Type, "test")

var testHelperLocals = at.CompileLocalsConfig(map[string]any{
	"path": testhelp.GetFixturesDirPath("snowflake_database"),
})

var testHelperDefinition = map[string]any{
	`"SnowflakeDatabaseProperties.json"`: map[string]any{
		"source": "${local.path}/SnowflakeDatabaseProperties.json.tmpl",
		"tokens": map[string]any{
			"DATABASE_NAME": "ExampleDatabase",
			"CONNECTION_ID": "11111111-1111-1111-1111-111111111111",
		},
	},
}

func TestUnit_SnowflakeDatabaseDefinitionTemplates(t *testing.T) {
	const connectionID = "11111111-1111-1111-1111-111111111111"

	tests := []struct {
		name            string
		source          string
		processingMode  string
		tokens          map[string]string
		parameters      []*params.ParametersModel
		tokensDelimiter string
		wantName        string
		wantConnection  string
	}{
		{
			name:            "Go template",
			source:          testhelp.GetFixturesDirPath("snowflake_database", "SnowflakeDatabaseProperties.json.tmpl"),
			processingMode:  transforms.ProcessingModeGoTemplate,
			tokens:          map[string]string{"DATABASE_NAME": "ExampleDatabase", "CONNECTION_ID": connectionID},
			tokensDelimiter: transforms.TokensDelimiterCurlyBraces,
			wantName:        "ExampleDatabase",
			wantConnection:  connectionID,
		},
		{
			name:            "custom Go template delimiter",
			source:          testhelp.GetFixturesDirPath("snowflake_database", "SnowflakeDatabasePropertiesCustomDelimiter.json.tmpl"),
			processingMode:  transforms.ProcessingModeGoTemplate,
			tokens:          map[string]string{"DATABASE_NAME": "ExampleDatabase", "CONNECTION_ID": connectionID},
			tokensDelimiter: "<<>>",
			wantName:        "ExampleDatabase",
			wantConnection:  connectionID,
		},
		{
			name:           "JSONPath parameters",
			source:         testhelp.GetFixturesDirPath("snowflake_database", "SnowflakeDatabasePropertiesParameters.json"),
			processingMode: transforms.ProcessingModeParameters,
			parameters: []*params.ParametersModel{
				{
					Type:  types.StringValue(transforms.ParameterTypeJSONPathReplace),
					Find:  types.StringValue("$.snowflakeDatabaseName"),
					Value: types.StringValue("UpdatedDatabase"),
				},
				{
					Type:  types.StringValue(transforms.ParameterTypeJSONPathReplace),
					Find:  types.StringValue("$.connectionId"),
					Value: types.StringValue(connectionID),
				},
			},
			tokensDelimiter: transforms.TokensDelimiterCurlyBraces,
			wantName:        "UpdatedDatabase",
			wantConnection:  connectionID,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			payload, _, diags := transforms.SourceFileToPayload(
				test.source,
				test.processingMode,
				test.tokens,
				test.parameters,
				test.tokensDelimiter,
			)
			if diags.HasError() {
				t.Fatalf("unexpected template processing diagnostics: %v", diags)
			}

			content, err := transforms.Base64Decode(payload)
			if err != nil {
				t.Fatalf("failed to decode definition payload: %v", err)
			}

			var definition struct {
				SnowflakeDatabaseName string `json:"snowflakeDatabaseName"`
				ConnectionID          string `json:"connectionId"`
			}
			if err := json.Unmarshal([]byte(content), &definition); err != nil {
				t.Fatalf("failed to decode Snowflake Database definition: %v", err)
			}

			if definition.SnowflakeDatabaseName != test.wantName {
				t.Errorf("unexpected snowflakeDatabaseName: got %q, want %q", definition.SnowflakeDatabaseName, test.wantName)
			}

			if definition.ConnectionID != test.wantConnection {
				t.Errorf("unexpected connectionId: got %q, want %q", definition.ConnectionID, test.wantConnection)
			}
		})
	}
}

func TestUnit_SnowflakeDatabaseResource_Attributes(t *testing.T) {
	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		// error - no attributes
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{},
			),
			ExpectError: regexp.MustCompile(`Missing required argument`),
		},
		// error - workspace_id - invalid UUID
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": "invalid uuid",
					"display_name": "test",
				},
			),
			ExpectError: regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
		// error - unexpected attribute
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id":    "00000000-0000-0000-0000-000000000000",
					"display_name":    "test",
					"unexpected_attr": "test",
				},
			),
			ExpectError: regexp.MustCompile(`An argument named "unexpected_attr" is not expected here`),
		},
		// error - no required attributes
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"display_name": "test",
				},
			),
			ExpectError: regexp.MustCompile(`The argument "workspace_id" is required, but no definition was found.`),
		},
		// error - no required attributes
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": "00000000-0000-0000-0000-000000000000",
				},
			),
			ExpectError: regexp.MustCompile(`The argument "display_name" is required, but no definition was found.`),
		},
		// error - no required attributes (configuration or definition)
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": "00000000-0000-0000-0000-000000000000",
					"display_name": "test",
				},
			),
			ExpectError: regexp.MustCompile(`Missing Attribute Configuration`),
		},
		// error - an empty definition does not provide the required definition part
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": "00000000-0000-0000-0000-000000000000",
					"display_name": "test",
					"format":       "Default",
					"definition":   map[string]any{},
				},
			),
			ExpectError: regexp.MustCompile(`Attribute definition map must contain at least 1 elements`),
		},
		// error - both configuration and definition are set
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": "00000000-0000-0000-0000-000000000000",
					"display_name": "test",
					"configuration": map[string]any{
						"connection_id":           testhelp.RandomUUID(),
						"snowflake_database_name": "ExampleDatabase",
					},
					"definition": map[string]any{},
				},
			),
			ExpectError: regexp.MustCompile(`Invalid Attribute Combination`),
		},
		// error - missing required configuration attribute
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id":  "00000000-0000-0000-0000-000000000000",
					"display_name":  "test",
					"configuration": map[string]any{},
				},
			),
			ExpectError: regexp.MustCompile(`Inappropriate value for attribute "configuration".`),
		},
	}))
}

func TestUnit_SnowflakeDatabaseResource_ImportState(t *testing.T) {
	workspaceID := testhelp.RandomUUID()
	entity := fakes.NewRandomSnowflakeDatabaseWithWorkspace(workspaceID)

	fakes.FakeServer.Upsert(fakes.NewRandomSnowflakeDatabaseWithWorkspace(workspaceID))
	fakes.FakeServer.Upsert(entity)
	fakes.FakeServer.Upsert(fakes.NewRandomSnowflakeDatabaseWithWorkspace(workspaceID))

	testCase := at.JoinConfigs(
		testHelperLocals,
		at.CompileConfig(
			testResourceItemHeader,
			map[string]any{
				"workspace_id": *entity.WorkspaceID,
				"display_name": *entity.DisplayName,
				"format":       "Default",
				"definition":   testHelperDefinition,
			},
		),
	)

	resource.Test(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		{
			ResourceName:  testResourceItemFQN,
			Config:        testCase,
			ImportStateId: "not-valid",
			ImportState:   true,
			ExpectError:   regexp.MustCompile(fmt.Sprintf(common.ErrorImportIdentifierDetails, fmt.Sprintf("WorkspaceID/%sID", string(fabricItemType)))),
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
			ImportStateId: fmt.Sprintf("%s/%s", "test", *entity.ID),
			ImportState:   true,
			ExpectError:   regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
		{
			ResourceName:  testResourceItemFQN,
			Config:        testCase,
			ImportStateId: fmt.Sprintf("%s/%s", *entity.WorkspaceID, "test"),
			ImportState:   true,
			ExpectError:   regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
		// Import state testing
		{
			ResourceName:       testResourceItemFQN,
			Config:             testCase,
			ImportStateId:      fmt.Sprintf("%s/%s", *entity.WorkspaceID, *entity.ID),
			ImportState:        true,
			ImportStatePersist: true,
			ImportStateCheck: func(is []*terraform.InstanceState) error {
				if len(is) != 1 {
					return errors.New("expected one instance state")
				}

				if is[0].ID != *entity.ID {
					return errors.New(testResourceItemFQN + ": unexpected ID")
				}

				return nil
			},
		},
	}))
}

func TestUnit_SnowflakeDatabaseResource_CRUD(t *testing.T) {
	workspaceID := testhelp.RandomUUID()
	entityExist := fakes.NewRandomSnowflakeDatabaseWithWorkspace(workspaceID)
	entityBefore := fakes.NewRandomSnowflakeDatabaseWithWorkspace(workspaceID)
	entityAfter := fakes.NewRandomSnowflakeDatabaseWithWorkspace(workspaceID)

	fakes.FakeServer.Upsert(fakes.NewRandomSnowflakeDatabaseWithWorkspace(workspaceID))
	fakes.FakeServer.Upsert(entityExist)
	fakes.FakeServer.Upsert(entityAfter)
	fakes.FakeServer.Upsert(fakes.NewRandomSnowflakeDatabaseWithWorkspace(workspaceID))

	resource.Test(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		// error - create - existing entity
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				testHelperLocals,
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"workspace_id": *entityExist.WorkspaceID,
						"display_name": *entityExist.DisplayName,
						"format":       "Default",
						"definition":   testHelperDefinition,
					},
				),
			),
			ExpectError: regexp.MustCompile(common.ErrorCreateHeader),
		},
		// Create and Read
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				testHelperLocals,
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"workspace_id": *entityBefore.WorkspaceID,
						"display_name": *entityBefore.DisplayName,
						"format":       "Default",
						"definition":   testHelperDefinition,
					},
				),
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "display_name", entityBefore.DisplayName),
				resource.TestCheckResourceAttr(testResourceItemFQN, "description", ""),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.connection_id"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.onelake_tables_path"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.snowflake_account_url"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.snowflake_database_name"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.snowflake_volume_path"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.sql_endpoint_properties.id"),
			),
		},
		// Update and Read
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				testHelperLocals,
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"workspace_id": *entityBefore.WorkspaceID,
						"display_name": *entityBefore.DisplayName,
						"description":  *entityBefore.Description,
						"format":       "Default",
						"definition":   testHelperDefinition,
					},
				),
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "display_name", entityBefore.DisplayName),
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "description", entityBefore.Description),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.connection_id"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.onelake_tables_path"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.snowflake_account_url"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.snowflake_database_name"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.snowflake_volume_path"),
			),
		},
		// Delete testing automatically occurs in TestCase
	}))
}

func TestUnit_SnowflakeDatabaseResource_CRUD_Configuration(t *testing.T) {
	workspaceID := testhelp.RandomUUID()
	entityCreate := fakes.NewRandomSnowflakeDatabaseWithWorkspace(workspaceID)

	resource.Test(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		// Create and Read with configuration
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": *entityCreate.WorkspaceID,
					"display_name": *entityCreate.DisplayName,
					"configuration": map[string]any{
						"connection_id":           testhelp.RandomUUID(),
						"snowflake_database_name": "ExampleDatabase",
					},
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "display_name", entityCreate.DisplayName),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "configuration.connection_id"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "configuration.snowflake_database_name", "ExampleDatabase"),
			),
		},
	}))
}
