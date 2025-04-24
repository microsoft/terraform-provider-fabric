package onelakeshortcut_test

import (
	"testing"

	at "github.com/dcarbone/terraform-plugin-framework-utils/v3/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
)

var testResourceItemFQN, testResourceItemHeader = testhelp.TFResource(common.ProviderTypeName, itemTypeInfo.Type, "test")

func TestAcc_DomainResource_CRUD(t *testing.T) {
	// entityCreateDisplayName := testhelp.RandomName()
	// entityUpdateDisplayName := testhelp.RandomName()
	// entityUpdateDescription := testhelp.RandomName()
	// defaultContributorsScope := string(admin.ContributorsScopeTypeAllTenant)
	// updatedContributorsScope := string(admin.ContributorsScopeTypeAdminsOnly)

	resource.Test(t, testhelp.NewTestAccCase(t, &testResourceItemFQN, nil, []resource.TestStep{
		// Create and Read
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"item_id":                  testhelp.WellKnown()["Lakehouse"].(map[string]any)["id"].(string),
					"workspace_id":             testhelp.WellKnown()["WorkspaceDS"].(map[string]any)["id"].(string),
					"shortcut_conflict_policy": "CreateOrOverwrite",
					"name":                     "acc_test",
					"path":                     "/Tables",
					"target": map[string]any{
						"type": "OneLake",
						"onelake": map[string]any{
							"workspace_id": testhelp.WellKnown()["WorkspaceDS"].(map[string]any)["id"].(string),
							"item_id":      testhelp.WellKnown()["Lakehouse"].(map[string]any)["id"].(string),
							"path":         "/Tables/publicholidays_1",
						},
					},
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "name", "acc_test"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "path", "/Tables"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "target.0.type", "onelake"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "target.0.onelake.workspace_id", testhelp.WellKnown()["WorkspaceDS"].(map[string]any)["id"].(string)),
				resource.TestCheckResourceAttr(testResourceItemFQN, "target.0.onelake.item_id", testhelp.WellKnown()["Lakehouse"].(map[string]any)["id"].(string)),
				resource.TestCheckResourceAttr(testResourceItemFQN, "target.0.onelake.path", "/Tables/publicholidays"),
			),
		},
		// // Update and Read
		// {
		// 	ResourceName: testResourceItemFQN,
		// 	Config: at.CompileConfig(
		// 		testResourceItemHeader,
		// 		map[string]any{
		// 			"display_name":       entityUpdateDisplayName,
		// 			"description":        entityUpdateDescription,
		// 			"contributors_scope": updatedContributorsScope,
		// 		},
		// 	),
		// 	Check: resource.ComposeAggregateTestCheckFunc(
		// 		resource.TestCheckResourceAttr(testResourceItemFQN, "display_name", entityUpdateDisplayName),
		// 		resource.TestCheckResourceAttr(testResourceItemFQN, "description", entityUpdateDescription),
		// 		resource.TestCheckResourceAttr(testResourceItemFQN, "contributors_scope", updatedContributorsScope),
		// 	),
		// },
	},
	))
}
