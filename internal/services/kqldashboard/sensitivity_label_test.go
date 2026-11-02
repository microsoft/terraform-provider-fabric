// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package kqldashboard_test

import (
	"testing"

	at "github.com/dcarbone/terraform-plugin-framework-utils/v3/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/stretchr/testify/require"

	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
)

func TestAcc_KQLDashboardDataSource_SensitivityLabel(t *testing.T) {
	wellKnown := testhelp.WellKnown()

	entity, ok := wellKnown["KQLDashboard"].(map[string]any)
	require.True(t, ok)

	fixtureLabel, exists := entity["sensitivityLabel"]
	if !exists || fixtureLabel == nil {
		t.Skip("Set FABRIC_TESTACC_WELLKNOWN_SENSITIVITY_LABEL_ID and rerun task testacc:setup to label the existing KQLDashboard fixture")
	}

	workspace, ok := wellKnown["WorkspaceDS"].(map[string]any)
	require.True(t, ok)

	label, ok := fixtureLabel.(map[string]any)
	require.True(t, ok)

	labelID, ok := label["id"].(string)
	require.True(t, ok)
	require.NotEmpty(t, labelID)

	entityID, ok := entity["id"].(string)
	require.True(t, ok)
	require.NotEmpty(t, entityID)

	resource.ParallelTest(t, testhelp.NewTestAccCase(t, nil, nil, []resource.TestStep{
		{
			Config: at.JoinConfigs(
				at.CompileConfig(testDataSourceItemHeader, map[string]any{
					"workspace_id": workspace["id"],
					"id":           entityID,
				}),
				at.CompileConfig(testDataSourceItemsHeader, map[string]any{
					"workspace_id": workspace["id"],
				}),
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "id", entityID),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "sensitivity_label.id", labelID),
				resource.TestCheckNoResourceAttr(testDataSourceItemFQN, "sensitivity_label.label_id"),
			),
			ConfigStateChecks: []statecheck.StateCheck{
				statecheck.ExpectKnownValue(
					testDataSourceItemsFQN,
					tfjsonpath.New("values"),
					knownvalue.SetPartial([]knownvalue.Check{
						knownvalue.ObjectPartial(map[string]knownvalue.Check{
							"id": knownvalue.StringExact(entityID),
							"sensitivity_label": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"id": knownvalue.StringExact(labelID),
							}),
						}),
					}),
				),
			},
		},
	}))
}
