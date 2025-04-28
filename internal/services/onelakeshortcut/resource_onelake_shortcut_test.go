package onelakeshortcut_test

import (
	"testing"

	at "github.com/dcarbone/terraform-plugin-framework-utils/v3/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
)

var testResourceItemFQN, testResourceItemHeader = testhelp.TFResource(common.ProviderTypeName, itemTypeInfo.Type, "test")

func TestAcc_OneLakeShortcutResource_CRUD(t *testing.T) {
	entityCreateDisplayName := testhelp.RandomName()
	entityTargetPath := "Files/images"
	entityType := "OneLake"
	entityUpdatedTargetPath := "Files/sample_datasets"
	workspaceId := testhelp.WellKnown()["WorkspaceDS"].(map[string]any)["id"].(string)
	lakehouseId := testhelp.WellKnown()["Lakehouse"].(map[string]any)["id"].(string)
	resource.Test(t, testhelp.NewTestAccCase(t, &testResourceItemFQN, nil, []resource.TestStep{
		// Create and Read
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"item_id":      lakehouseId,
					"workspace_id": workspaceId,
					"name":         entityCreateDisplayName,
					"path":         "Files",
					"target": map[string]any{
						"type": entityType,
						"onelake": map[string]any{
							"workspace_id": workspaceId,
							"item_id":      lakehouseId,
							"path":         entityTargetPath,
						},
					},
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "name", entityCreateDisplayName),
				resource.TestCheckResourceAttr(testResourceItemFQN, "target.onelake.item_id", lakehouseId),
				resource.TestCheckResourceAttr(testResourceItemFQN, "target.onelake.workspace_id", workspaceId),
				resource.TestCheckResourceAttr(testResourceItemFQN, "target.onelake.path", entityTargetPath),
				resource.TestCheckResourceAttr(testResourceItemFQN, "target.type", entityType),
			),
		},
		// Update - Create with OverwriteOnly and Read
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"item_id":                  lakehouseId,
					"workspace_id":             workspaceId,
					"shortcut_conflict_policy": "OverwriteOnly",
					"name":                     entityCreateDisplayName,
					"path":                     "Files",
					"target": map[string]any{
						"type": entityType,
						"onelake": map[string]any{
							"workspace_id": workspaceId,
							"item_id":      lakehouseId,
							"path":         entityUpdatedTargetPath,
						},
					},
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "name", entityCreateDisplayName),
				resource.TestCheckResourceAttr(testResourceItemFQN, "target.onelake.item_id", lakehouseId),
				resource.TestCheckResourceAttr(testResourceItemFQN, "target.onelake.workspace_id", workspaceId),
				resource.TestCheckResourceAttr(testResourceItemFQN, "target.onelake.path", entityUpdatedTargetPath),
				resource.TestCheckResourceAttr(testResourceItemFQN, "target.type", entityType),
			),
		},
	},
	))
}
