// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package onelakedataaccesssecurity_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp/fakes"
)

func TestUnit_OneLakeDataAccessSecurityResource_FabricItemMembersDefaults(t *testing.T) {
	for name, attribute := range map[string]string{
		"omitted": "",
		"null":    "fabric_item_members = null",
		"empty":   "fabric_item_members = []",
	} {
		t.Run(name, func(t *testing.T) {
			workspaceID := testhelp.RandomUUID()
			itemID := testhelp.RandomUUID()
			objectID := testhelp.RandomUUID()
			tenantID := testhelp.RandomUUID()
			roleID := GenerateOneLakeDataAccessRoleID(workspaceID, itemID, "member-defaults")

			config := func(fabricMembers string) string {
				return fmt.Sprintf(`
					%s {
						workspace_id = "%s"
						item_id      = "%s"
						role_name    = "member-defaults"
						kind         = "Policy"
						decision_rules = [{
							effect = "Permit"
							permission = [
								{ attribute_name = "Path", attribute_value_included_in = ["*"] },
								{ attribute_name = "Action", attribute_value_included_in = ["Read"] }
							]
						}]
						members = {
							%s
							microsoft_entra_members = [{
								object_id = "%s"
								tenant_id = "%s"
							}]
						}
					}
				`, testResourceItemHeader, workspaceID, itemID, fabricMembers, objectID, tenantID)
			}
			emptyConfig := config(attribute)
			populatedConfig := config(fmt.Sprintf(`
				fabric_item_members = [{
					item_access = ["ReadAll"]
					source_path = "%s/%s"
				}]
			`, workspaceID, itemID))

			fakeServer := fakes.NewFakeServer()
			fakeServer.ServerFactory.Core.OneLakeDataAccessSecurityServer.CreateOrUpdateSingleDataAccessRole = fakeCreateOrUpdateSingleDataAccessRoleFunc()
			fakeServer.ServerFactory.Core.OneLakeDataAccessSecurityServer.GetDataAccessRole = fakeGetDataAccessRoleFunc()
			fakeServer.ServerFactory.Core.OneLakeDataAccessSecurityServer.DeleteDataAccessRole = fakeDeleteDataAccessRoleFunc()

			resource.Test(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, fakeServer.ServerFactory, nil, []resource.TestStep{
				{
					Config: emptyConfig,
					Check:  resource.TestCheckResourceAttr(testResourceItemFQN, "members.fabric_item_members.#", "0"),
				},
				{
					Config:   emptyConfig,
					PlanOnly: true,
				},
				{
					Config: populatedConfig,
					Check:  resource.TestCheckResourceAttr(testResourceItemFQN, "members.fabric_item_members.#", "1"),
				},
				{
					Config: emptyConfig,
					Check:  resource.TestCheckResourceAttr(testResourceItemFQN, "members.fabric_item_members.#", "0"),
				},
				{
					Config:   emptyConfig,
					PlanOnly: true,
				},
				{
					ResourceName:                         testResourceItemFQN,
					ImportState:                          true,
					ImportStateId:                        roleID,
					ImportStateVerify:                    true,
					ImportStateVerifyIdentifierAttribute: "role_name",
				},
			}))
		})
	}
}
