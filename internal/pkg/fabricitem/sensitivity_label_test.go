// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package fabricitem_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"go/parser"
	"go/token"
	"net/http"
	"path/filepath"
	"regexp"
	"testing"
	"text/template"

	azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	fabfake "github.com/microsoft/fabric-sdk-go/fabric/fake"
	"github.com/stretchr/testify/require"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp/fakes"
)

func TestUnit_FabricItemSensitivityLabelTemplates(t *testing.T) {
	for _, archetype := range []string{"properties", "definition-properties", "config-properties", "config-definition-properties"} {
		t.Run(archetype, func(t *testing.T) {
			path := filepath.Join("../../../tools/itemgen/templates", archetype, "data_items.go.tmpl")
			tmpl, err := template.ParseFiles(path)
			require.NoError(t, err)

			var output bytes.Buffer

			err = tmpl.Execute(&output, map[string]string{
				"Package":   "lakehouse",
				"ModelName": "lakehouse",
				"TypeInfo":  "Lakehouse",
				"TypesInfo": "Lakehouses",
			})
			require.NoError(t, err)
			require.Contains(t, output.String(), "*fabricitem.FabricItemPropertiesModel[lakehousePropertiesModel, fablakehouse.Properties]")

			_, err = parser.ParseFile(token.NewFileSet(), path, output.Bytes(), parser.AllErrors)
			require.NoError(t, err)
		})
	}
}

func TestUnit_FabricItemSensitivityLabels(t *testing.T) {
	testCases := []struct {
		name        string
		plural      string
		extraConfig string
	}{
		{name: "kql_dashboard", plural: "kql_dashboards"},
		{name: "notebook", plural: "notebooks"},
		{name: "environment", plural: "environments"},
		{name: "spark_job_definition", plural: "spark_job_definitions"},
		{name: "warehouse", plural: "warehouses"},
		{name: "lakehouse", plural: "lakehouses"},
		{name: "operations_agent", plural: "operations_agents"},
		{name: "org_app", plural: "org_apps"},
		{
			name:   "mirrored_catalog",
			plural: "mirrored_catalogs",
			extraConfig: fmt.Sprintf(`
				format = "Default"
				definition = {
					"MirroredCatalogDefinition.json" = {
						source = %q
						tokens = {
							CONNECTION_ID = "11111111-1111-1111-1111-111111111111"
							SCOPE = "default"
						}
					}
				}`, filepath.Join(testhelp.GetFixturesDirPath("mirrored_catalog"), "MirroredCatalogDefinition.json.tmpl")),
		},
		{
			name:   "mirrored_azure_databricks_catalog",
			plural: "mirrored_azure_databricks_catalogs",
			extraConfig: `
				configuration = {
					catalog_name = "test"
					databricks_workspace_connection_id = "11111111-1111-1111-1111-111111111111"
					mirroring_mode = "Partial"
				}`,
		},
		{
			name:   "snowflake_database",
			plural: "snowflake_databases",
			extraConfig: `
				configuration = {
					connection_id = "11111111-1111-1111-1111-111111111111"
					snowflake_database_name = "test"
				}`,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			server := fakes.NewFakeServer()
			requests := make(chan fabcore.CreateItemRequest, 1)
			createItem := server.ServerFactory.Core.ItemsServer.BeginCreateItem
			server.ServerFactory.Core.ItemsServer.BeginCreateItem = func(ctx context.Context, workspaceID string, request fabcore.CreateItemRequest, options *fabcore.ItemsClientBeginCreateItemOptions) (azfake.PollerResponder[fabcore.ItemsClientCreateItemResponse], azfake.ErrorResponder) {
				requests <- request

				return createItem(ctx, workspaceID, request, options)
			}
			workspaceID := testhelp.RandomUUID()
			labelID := testhelp.RandomUUID()
			tagID := testhelp.RandomUUID()
			resourceFQN := "fabric_" + testCase.name + ".test"
			dataSourceFQN := "data." + resourceFQN
			dataSourcesFQN := "data.fabric_" + testCase.plural + ".test"

			config := fmt.Sprintf(`
				resource "fabric_%[1]s" "test" {
					workspace_id = %[3]q
					display_name = "sensitivity-label-test"
					sensitivity_label = {
						id = %[4]q
					}
					tags = [%[5]q]
					%[6]s
				}
				data "fabric_%[1]s" "test" {
					workspace_id = fabric_%[1]s.test.workspace_id
					id = fabric_%[1]s.test.id
				}
				data "fabric_%[2]s" "test" {
					workspace_id = fabric_%[1]s.test.workspace_id
					depends_on = [fabric_%[1]s.test]
				}`,
				testCase.name, testCase.plural, workspaceID, labelID, tagID, testCase.extraConfig,
			)

			resource.ParallelTest(t, testhelp.NewTestUnitCase(t, &resourceFQN, server.ServerFactory, nil, []resource.TestStep{
				{
					Config: config,
					Check: resource.ComposeAggregateTestCheckFunc(
						checkSensitivityLabelRequest(requests, &labelID),
						resource.TestCheckResourceAttr(resourceFQN, "sensitivity_label.id", labelID),
						resource.TestCheckTypeSetElemAttr(resourceFQN, "tags.*", tagID),
						resource.TestCheckResourceAttr(dataSourceFQN, "sensitivity_label.id", labelID),
						resource.TestCheckTypeSetElemAttr(dataSourceFQN, "tags.*", tagID),
					),
					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectKnownValue(
							dataSourcesFQN,
							tfjsonpath.New("values"),
							knownvalue.SetExact([]knownvalue.Check{
								knownvalue.ObjectPartial(map[string]knownvalue.Check{
									"sensitivity_label": knownvalue.ObjectExact(map[string]knownvalue.Check{
										"id": knownvalue.StringExact(labelID),
									}),
									"tags": knownvalue.SetExact([]knownvalue.Check{knownvalue.StringExact(tagID)}),
								}),
							}),
						),
					},
				},
			}))
		})
	}
}

func checkSensitivityLabelRequest(requests <-chan fabcore.CreateItemRequest, labelID *string) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		select {
		case request := <-requests:
			if labelID == nil {
				if request.SensitivityLabelSettings != nil {
					return errors.New("unconfigured sensitivity label must not be sent in the create request")
				}

				return nil
			}

			settings := request.SensitivityLabelSettings
			if settings == nil || settings.LabelID == nil || *settings.LabelID != *labelID {
				return errors.New("create request did not contain the configured sensitivity label ID")
			}

			if settings.SensitivityLabelApplyStrategy == nil || *settings.SensitivityLabelApplyStrategy != fabcore.SensitivityLabelApplyStrategyApplyOrFail {
				return errors.New("create request must explicitly use ApplyOrFail")
			}

			return nil
		default:
			return errors.New("expected a Fabric create request")
		}
	}
}

func TestUnit_FabricItemSensitivityLabel_Import(t *testing.T) {
	for _, labelled := range []bool{false, true} {
		t.Run(fmt.Sprintf("labelled=%t", labelled), func(t *testing.T) {
			server := fakes.NewFakeServer()
			entity := fakes.NewRandomItem(fabcore.ItemTypeMLModel)
			entity.FolderID = nil
			entity.SensitivityLabel = nil

			labelID := testhelp.RandomUUID()
			if labelled {
				entity.SensitivityLabel = &fabcore.SensitivityLabel{ID: &labelID}
			}

			server.Upsert(entity)

			fqn := "fabric_ml_model.test"
			config := fmt.Sprintf(`
				resource "fabric_ml_model" "test" {
					workspace_id = %q
					display_name = %q
					description = %q
				}`, *entity.WorkspaceID, *entity.DisplayName, *entity.Description)

			labelConfig := config
			if labelled {
				labelConfig = fmt.Sprintf(`
					resource "fabric_ml_model" "test" {
						workspace_id = %q
						display_name = %q
						description = %q
						sensitivity_label = { id = %q }
					}`, *entity.WorkspaceID, *entity.DisplayName, *entity.Description, labelID)
			}

			check := resource.TestCheckNoResourceAttr(fqn, "sensitivity_label.id")
			if labelled {
				check = resource.TestCheckResourceAttr(fqn, "sensitivity_label.id", labelID)
			}

			steps := []resource.TestStep{
				{
					Config:             labelConfig,
					ResourceName:       fqn,
					ImportState:        true,
					ImportStateId:      *entity.WorkspaceID + "/" + *entity.ID,
					ImportStatePersist: true,
				},
				{
					Config: labelConfig,
					Check:  check,
					ConfigPlanChecks: resource.ConfigPlanChecks{
						PreApply: []plancheck.PlanCheck{
							plancheck.ExpectResourceAction(fqn, plancheck.ResourceActionNoop),
						},
					},
				},
				{
					Config: config,
					Check:  check,
					ConfigPlanChecks: resource.ConfigPlanChecks{
						PreApply: []plancheck.PlanCheck{
							plancheck.ExpectResourceAction(fqn, plancheck.ResourceActionNoop),
						},
					},
				},
			}

			if labelled {
				replacementLabelID := testhelp.RandomUUID()
				steps = append(steps, resource.TestStep{
					Config: fmt.Sprintf(`
						resource "fabric_ml_model" "test" {
							workspace_id = %q
							display_name = %q
							description = %q
							sensitivity_label = { id = %q }
						}`, *entity.WorkspaceID, *entity.DisplayName, *entity.Description, replacementLabelID),
					Check: resource.TestCheckResourceAttr(fqn, "sensitivity_label.id", replacementLabelID),
					ConfigPlanChecks: resource.ConfigPlanChecks{
						PreApply: []plancheck.PlanCheck{
							plancheck.ExpectResourceAction(fqn, plancheck.ResourceActionDestroyBeforeCreate),
						},
					},
				})
			}

			resource.ParallelTest(t, testhelp.NewTestUnitCase(t, &fqn, server.ServerFactory, nil, steps))
		})
	}
}

func TestUnit_FabricItemSensitivityLabel_OmittedOnCreation(t *testing.T) {
	server := fakes.NewFakeServer()
	requests := make(chan fabcore.CreateItemRequest, 1)
	createItem := server.ServerFactory.Core.ItemsServer.BeginCreateItem
	server.ServerFactory.Core.ItemsServer.BeginCreateItem = func(ctx context.Context, workspaceID string, request fabcore.CreateItemRequest, options *fabcore.ItemsClientBeginCreateItemOptions) (azfake.PollerResponder[fabcore.ItemsClientCreateItemResponse], azfake.ErrorResponder) {
		requests <- request

		return createItem(ctx, workspaceID, request, options)
	}
	fqn := "fabric_ml_model.test"
	config := fmt.Sprintf(`
		resource "fabric_ml_model" "test" {
			workspace_id = %q
			display_name = "inherited-label"
		}`, testhelp.RandomUUID())

	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, &fqn, server.ServerFactory, nil, []resource.TestStep{
		{
			Config: config,
			Check: resource.ComposeAggregateTestCheckFunc(
				checkSensitivityLabelRequest(requests, nil),
				resource.TestCheckResourceAttrSet(fqn, "sensitivity_label.id"),
			),
		},
		{
			Config: config,
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PreApply: []plancheck.PlanCheck{
					plancheck.ExpectResourceAction(fqn, plancheck.ResourceActionNoop),
				},
			},
		},
	}))
}

func TestUnit_FabricItemSensitivityLabel_ApplyFailure(t *testing.T) {
	server := fakes.NewFakeServer()
	server.ServerFactory.Core.ItemsServer.BeginCreateItem = func(_ context.Context, _ string, _ fabcore.CreateItemRequest, _ *fabcore.ItemsClientBeginCreateItemOptions) (resp azfake.PollerResponder[fabcore.ItemsClientCreateItemResponse], errResp azfake.ErrorResponder) {
		errResp.SetError(fabfake.SetResponseError(http.StatusBadRequest, "SensitivityLabelNotFound", "The requested sensitivity label could not be applied."))

		return resp, errResp
	}

	config := fmt.Sprintf(`
		resource "fabric_ml_model" "test" {
			workspace_id = %q
			display_name = "label-application-fails"
			sensitivity_label = { id = %q }
		}`, testhelp.RandomUUID(), testhelp.RandomUUID())

	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, nil, server.ServerFactory, nil, []resource.TestStep{
		{
			Config:      config,
			ExpectError: regexp.MustCompile("SensitivityLabelNotFound"),
		},
	}))
}

func TestUnit_FabricItemSensitivityLabel_Attributes(t *testing.T) {
	server := fakes.NewFakeServer()
	config := fmt.Sprintf(`
		resource "fabric_ml_model" "test" {
			workspace_id = %q
			display_name = "label-validation"
			sensitivity_label = %%s
		}`, testhelp.RandomUUID())

	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, nil, server.ServerFactory, nil, []resource.TestStep{
		{
			Config:      fmt.Sprintf(config, "{}"),
			ExpectError: regexp.MustCompile(`sensitivity_label\.id`),
		},
		{
			Config:      fmt.Sprintf(config, `{ id = "not-a-uuid" }`),
			ExpectError: regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
	}))
}

func TestUnit_FabricItemSensitivityLabel_UnknownDuringPlan(t *testing.T) {
	server := fakes.NewFakeServer()
	workspaceID := testhelp.RandomUUID()
	fqn := "fabric_ml_model.target"
	config := fmt.Sprintf(`
		resource "fabric_ml_model" "source" {
			workspace_id = %[1]q
			display_name = "label-source"
		}
		resource "fabric_ml_model" "target" {
			workspace_id = %[1]q
			display_name = "label-target"
			sensitivity_label = {
				id = fabric_ml_model.source.sensitivity_label.id
			}
		}`, workspaceID)

	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, &fqn, server.ServerFactory, nil, []resource.TestStep{
		{
			Config: config,
			Check: resource.TestCheckResourceAttrPair(
				fqn, "sensitivity_label.id",
				"fabric_ml_model.source", "sensitivity_label.id",
			),
		},
	}))
}
