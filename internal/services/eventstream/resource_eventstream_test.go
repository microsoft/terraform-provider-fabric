// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package eventstream_test

import (
	"errors"
	"fmt"
	"regexp"
	"testing"

	at "github.com/dcarbone/terraform-plugin-framework-utils/v3/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp/fakes"
)

var testResourceItemFQN, testResourceItemHeader = testhelp.TFResource(common.ProviderTypeName, itemTypeInfo.Type, "test")

var testHelperLocals = at.CompileLocalsConfig(map[string]any{
	"path": testhelp.GetFixturesDirPath("eventstream"),
})

func TestUnit_EventstreamResource_Attributes(t *testing.T) {
	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		// error - no attributes
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				testHelperLocals,
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{},
				),
			),
			ExpectError: regexp.MustCompile(`Missing required argument`),
		},
		// error - workspace_id - invalid UUID
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				testHelperLocals,
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"workspace_id": "invalid uuid",
						"display_name": "test",
					},
				)),
			ExpectError: regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
		// error - unexpected attribute
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				testHelperLocals,
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"workspace_id":    "00000000-0000-0000-0000-000000000000",
						"unexpected_attr": "test",
					},
				)),
			ExpectError: regexp.MustCompile(`An argument named "unexpected_attr" is not expected here`),
		},
		// error - no required attributes
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				testHelperLocals,
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"display_name": "test",
					},
				)),
			ExpectError: regexp.MustCompile(`The argument "workspace_id" is required, but no definition was found.`),
		},
		// error - no required attributes
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				testHelperLocals,
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"workspace_id": "00000000-0000-0000-0000-000000000000",
					},
				)),
			ExpectError: regexp.MustCompile(`The argument "display_name" is required, but no definition was found.`),
		},
	}))
}

func TestUnit_EventstreamResource_ImportState(t *testing.T) {
	workspaceID := testhelp.RandomUUID()
	entity := fakes.NewRandomItemWithWorkspace(fabricItemType, workspaceID)

	fakes.FakeServer.Upsert(fakes.NewRandomItemWithWorkspace(fabricItemType, workspaceID))
	fakes.FakeServer.Upsert(entity)
	fakes.FakeServer.Upsert(fakes.NewRandomItemWithWorkspace(fabricItemType, workspaceID))

	testCase := at.JoinConfigs(
		testHelperLocals,
		at.CompileConfig(
			testResourceItemHeader,
			map[string]any{
				"workspace_id": *entity.WorkspaceID,
				"display_name": *entity.DisplayName,
			},
		))

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

func TestUnit_EventstreamResource_CRUD(t *testing.T) {
	workspaceID := testhelp.RandomUUID()
	entityExist := fakes.NewRandomItemWithWorkspace(fabricItemType, workspaceID)
	entityBefore := fakes.NewRandomItemWithWorkspace(fabricItemType, workspaceID)
	entityAfter := fakes.NewRandomItemWithWorkspace(fabricItemType, workspaceID)

	fakes.FakeServer.Upsert(fakes.NewRandomItemWithWorkspace(fabricItemType, workspaceID))
	fakes.FakeServer.Upsert(entityExist)
	fakes.FakeServer.Upsert(entityAfter)
	fakes.FakeServer.Upsert(fakes.NewRandomItemWithWorkspace(fabricItemType, workspaceID))

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
					},
				)),
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
					},
				)),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "display_name", entityBefore.DisplayName),
				resource.TestCheckResourceAttr(testResourceItemFQN, "description", ""),
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
						"display_name": *entityAfter.DisplayName,
						"description":  *entityAfter.Description,
					},
				)),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "display_name", entityAfter.DisplayName),
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "description", entityAfter.Description),
			),
		},
		// Delete testing automatically occurs in TestCase
	}))
}

func TestAcc_EventstreamResource_CRUD(t *testing.T) {
	workspace := testhelp.WellKnown()["WorkspaceRS"].(map[string]any)
	workspaceID := workspace["id"].(string)

	entityCreateDisplayName := testhelp.RandomName()
	entityUpdateDisplayName := testhelp.RandomName()
	entityUpdateDescription := testhelp.RandomName()

	resource.Test(t, testhelp.NewTestAccCase(t, &testResourceItemFQN, nil, []resource.TestStep{
		// Create and Read
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				testHelperLocals,
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"workspace_id": workspaceID,
						"display_name": entityCreateDisplayName,
					},
				)),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "display_name", entityCreateDisplayName),
				resource.TestCheckResourceAttr(testResourceItemFQN, "description", ""),
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
						"workspace_id": workspaceID,
						"display_name": entityUpdateDisplayName,
						"description":  entityUpdateDescription,
					},
				)),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "display_name", entityUpdateDisplayName),
				resource.TestCheckResourceAttr(testResourceItemFQN, "description", entityUpdateDescription),
			),
		},
	},
	))
}

func TestAcc_EventstreamDefinitionResource_CRUD(t *testing.T) {
	workspace := testhelp.WellKnown()["WorkspaceRS"].(map[string]any)
	workspaceID := workspace["id"].(string)

	entityCreateDisplayName := testhelp.RandomName()
	entityUpdateDisplayName := testhelp.RandomName()
	entityUpdateDescription := testhelp.RandomName()

	lakehouseResourceHCL, lakehouseResourceFQN := lakehouseResource(t, workspaceID)

	testHelperDefinition := map[string]any{
		`"eventstream.json"`: map[string]any{
			"source": "${local.path}/eventstream.json.tmpl",
			"tokens": map[string]any{
				"LakehouseWorkspaceID": testhelp.RefByFQN(lakehouseResourceFQN, "workspace_id"),
				"LakehouseID":          testhelp.RefByFQN(lakehouseResourceFQN, "id"),
			},
		},
	}

	resource.Test(t, testhelp.NewTestAccCase(t, &testResourceItemFQN, nil, []resource.TestStep{
		// Create and Read
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				lakehouseResourceHCL,
				testHelperLocals,
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"workspace_id": workspaceID,
						"display_name": entityCreateDisplayName,
						"format":       "Default",
						"definition":   testHelperDefinition,
					},
				)),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "display_name", entityCreateDisplayName),
				resource.TestCheckResourceAttr(testResourceItemFQN, "description", ""),
				resource.TestCheckResourceAttr(testResourceItemFQN, "definition_update_enabled", "true"),
			),
		},
		// Update and Read
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				lakehouseResourceHCL,
				testHelperLocals,
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"workspace_id": workspaceID,
						"display_name": entityUpdateDisplayName,
						"description":  entityUpdateDescription,
						"format":       "Default",
						"definition":   testHelperDefinition,
					},
				)),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "display_name", entityUpdateDisplayName),
				resource.TestCheckResourceAttr(testResourceItemFQN, "description", entityUpdateDescription),
				resource.TestCheckResourceAttr(testResourceItemFQN, "definition_update_enabled", "true"),
			),
		},
	},
	))
}

func TestAcc_EventstreamDefinitionResource_CRUD_Parameters(t *testing.T) {
	workspace := testhelp.WellKnown()["WorkspaceRS"].(map[string]any)
	workspaceID := workspace["id"].(string)

	entityCreateDisplayName := testhelp.RandomName()
	entityUpdateDisplayName := testhelp.RandomName()
	entityUpdateDescription := testhelp.RandomName()

	lakehouseResourceHCL, lakehouseResourceFQN := lakehouseResource(t, workspaceID)

	testHelperDefinition := map[string]any{
		`"eventstream.json"`: map[string]any{
			"source":          "${local.path}/eventstream.json",
			"processing_mode": "Parameters",
			"parameters": []map[string]any{
				{
					"type":  "TextReplace",
					"find":  "00000000-0000-0000-0000-000000000000",
					"value": testhelp.RefByFQN(lakehouseResourceFQN, "workspace_id"),
				},
				{
					"type":  "JsonPathReplace",
					"find":  `$.destinations[?(@.name=='Lakehouse')].properties.itemId`,
					"value": testhelp.RefByFQN(lakehouseResourceFQN, "id"),
				},
			},
		},
	}

	resource.Test(t, testhelp.NewTestAccCase(t, &testResourceItemFQN, nil, []resource.TestStep{
		// Create and Read
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				lakehouseResourceHCL,
				testHelperLocals,
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"workspace_id": workspaceID,
						"display_name": entityCreateDisplayName,
						"format":       "Default",
						"definition":   testHelperDefinition,
					},
				)),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "display_name", entityCreateDisplayName),
				resource.TestCheckResourceAttr(testResourceItemFQN, "description", ""),
				resource.TestCheckResourceAttr(testResourceItemFQN, "definition_update_enabled", "true"),
			),
		},
		// Update and Read
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				lakehouseResourceHCL,
				testHelperLocals,
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"workspace_id": workspaceID,
						"display_name": entityUpdateDisplayName,
						"description":  entityUpdateDescription,
						"format":       "Default",
						"definition":   testHelperDefinition,
					},
				)),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "display_name", entityUpdateDisplayName),
				resource.TestCheckResourceAttr(testResourceItemFQN, "description", entityUpdateDescription),
				resource.TestCheckResourceAttr(testResourceItemFQN, "definition_update_enabled", "true"),
			),
		},
	},
	))
}

func TestUnit_EventstreamDefinitionResource_DefinitionPartProcessing_Attributes(t *testing.T) { //nolint:go-golangci-lint
	workspaceID := testhelp.RandomUUID()

	testHelperDefinitionBase := map[string]any{
		`"eventstream.json"`: map[string]any{
			"source": "${local.path}/eventstream.json",
		},
	}

	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		// error - invalid processing_mode value
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				testHelperLocals,
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"workspace_id": workspaceID,
						"display_name": "test",
						"format":       "Default",
						"definition": map[string]any{
							`"eventstream.json"`: map[string]any{
								"source":          "${local.path}/eventstream.json",
								"processing_mode": "InvalidMode",
							},
						},
					},
				)),
			ExpectError: regexp.MustCompile(common.ErrorAttValueMatch),
		},
		// error - invalid tokens_delimiter value
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				testHelperLocals,
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"workspace_id": workspaceID,
						"display_name": "test",
						"format":       "Default",
						"definition": map[string]any{
							`"eventstream.json"`: map[string]any{
								"source":           "${local.path}/eventstream.json",
								"tokens":           map[string]any{"TestToken": "value"},
								"tokens_delimiter": "%%",
								"processing_mode":  "GoTemplate",
							},
						},
					},
				)),
			ExpectError: regexp.MustCompile(common.ErrorAttValueMatch),
		},
		// error - tokens_delimiter without tokens
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				testHelperLocals,
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"workspace_id": workspaceID,
						"display_name": "test",
						"format":       "Default",
						"definition": map[string]any{
							`"eventstream.json"`: map[string]any{
								"source":           "${local.path}/eventstream.json",
								"tokens_delimiter": "{{}}",
							},
						},
					},
				)),
			ExpectError: regexp.MustCompile(common.ErrorAttComboInvalid),
		},
		// error - tokens_delimiter conflicts with parameters
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				testHelperLocals,
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"workspace_id": workspaceID,
						"display_name": "test",
						"format":       "Default",
						"definition": map[string]any{
							`"eventstream.json"`: map[string]any{
								"source":           "${local.path}/eventstream.json",
								"processing_mode":  "Parameters",
								"tokens_delimiter": "{{}}",
								"parameters": []map[string]any{
									{"type": "TextReplace", "find": "test", "value": "value"},
								},
							},
						},
					},
				)),
			ExpectError: regexp.MustCompile(common.ErrorAttComboInvalid),
		},
		// error - tokens with processing_mode=Parameters (should use parameters instead)
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				testHelperLocals,
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"workspace_id": workspaceID,
						"display_name": "test",
						"format":       "Default",
						"definition": map[string]any{
							`"eventstream.json"`: map[string]any{
								"source":          "${local.path}/eventstream.json",
								"processing_mode": "Parameters",
								"tokens": map[string]any{
									"TestToken": "value",
								},
							},
						},
					},
				)),
			ExpectError: regexp.MustCompile("Invalid configuration for attribute"),
		},
		// error - tokens with processing_mode=None
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				testHelperLocals,
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"workspace_id": workspaceID,
						"display_name": "test",
						"format":       "Default",
						"definition": map[string]any{
							`"eventstream.json"`: map[string]any{
								"source":          "${local.path}/eventstream.json",
								"processing_mode": "None",
								"tokens": map[string]any{
									"TestToken": "value",
								},
							},
						},
					},
				)),
			ExpectError: regexp.MustCompile("Invalid configuration for attribute"),
		},
		// error - parameters with processing_mode=GoTemplate (should use tokens instead)
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				testHelperLocals,
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"workspace_id": workspaceID,
						"display_name": "test",
						"format":       "Default",
						"definition": map[string]any{
							`"eventstream.json"`: map[string]any{
								"source":          "${local.path}/eventstream.json",
								"processing_mode": "GoTemplate",
								"parameters": []map[string]any{
									{"type": "TextReplace", "find": "test", "value": "value"},
								},
							},
						},
					},
				)),
			ExpectError: regexp.MustCompile("Invalid configuration for attribute"),
		},
		// error - parameters with processing_mode=None
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				testHelperLocals,
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"workspace_id": workspaceID,
						"display_name": "test",
						"format":       "Default",
						"definition": map[string]any{
							`"eventstream.json"`: map[string]any{
								"source":          "${local.path}/eventstream.json",
								"processing_mode": "None",
								"parameters": []map[string]any{
									{"type": "TextReplace", "find": "test", "value": "value"},
								},
							},
						},
					},
				)),
			ExpectError: regexp.MustCompile("Invalid configuration for attribute"),
		},
		// error - parameters without processing_mode=Parameters
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				testHelperLocals,
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"workspace_id": workspaceID,
						"display_name": "test",
						"format":       "Default",
						"definition": map[string]any{
							`"eventstream.json"`: map[string]any{
								"source": "${local.path}/eventstream.json",
								"parameters": []map[string]any{
									{"type": "TextReplace", "find": "test", "value": "value"},
								},
							},
						},
					},
				)),
			ExpectError: regexp.MustCompile(common.ErrorAttComboInvalid),
		},
		// error - invalid parameter type
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				testHelperLocals,
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"workspace_id": workspaceID,
						"display_name": "test",
						"format":       "Default",
						"definition": map[string]any{
							`"eventstream.json"`: map[string]any{
								"source":          "${local.path}/eventstream.json",
								"processing_mode": "Parameters",
								"parameters": []map[string]any{
									{"type": "InvalidType", "find": "test", "value": "value"},
								},
							},
						},
					},
				)),
			ExpectError: regexp.MustCompile(common.ErrorAttValueMatch),
		},
		// success - valid tokens with GoTemplate mode (default delimiter)
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				testHelperLocals,
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"workspace_id": workspaceID,
						"display_name": "test_tokens_default",
						"format":       "Default",
						"definition": map[string]any{
							`"eventstream.json"`: map[string]any{
								"source":          "${local.path}/eventstream.json",
								"processing_mode": "GoTemplate",
								"tokens": map[string]any{
									"ValidToken":    "value1",
									"Another_Token": "value2",
								},
							},
						},
					},
				)),
		},
		// success - valid tokens with default processing_mode (default)
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				testHelperLocals,
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"workspace_id": workspaceID,
						"display_name": "test_tokens_default",
						"format":       "Default",
						"definition": map[string]any{
							`"eventstream.json"`: map[string]any{
								"source":          "${local.path}/eventstream.json",
								"processing_mode": "GoTemplate",
								"tokens": map[string]any{
									"ValidToken":    "value1",
									"Another_Token": "value2",
								},
							},
						},
					},
				)),
		},
		// success - valid tokens with custom delimiter
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				testHelperLocals,
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"workspace_id": workspaceID,
						"display_name": "test_tokens_custom_delim",
						"format":       "Default",
						"definition": map[string]any{
							`"eventstream.json"`: map[string]any{
								"source":           "${local.path}/eventstream.json",
								"processing_mode":  "GoTemplate",
								"tokens_delimiter": "<<>>",
								"tokens": map[string]any{
									"ValidToken": "value",
								},
							},
						},
					},
				)),
		},
		// success - valid parameters with Parameters mode
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				testHelperLocals,
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"workspace_id": workspaceID,
						"display_name": "test_parameters",
						"format":       "Default",
						"definition": map[string]any{
							`"eventstream.json"`: map[string]any{
								"source":          "${local.path}/eventstream.json",
								"processing_mode": "Parameters",
								"parameters": []map[string]any{
									{"type": "TextReplace", "find": "old", "value": "new"},
									{"type": "JsonPathReplace", "find": "$.path", "value": "newvalue"},
								},
							},
						},
					},
				)),
		},
		// success - None processing mode (no tokens or parameters)
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				testHelperLocals,
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"workspace_id": workspaceID,
						"display_name": "test_none_mode",
						"format":       "Default",
						"definition":   testHelperDefinitionBase,
					},
				)),
		},
	}))
}
