// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package fakes

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
)

func NewRandomOneLakeDataAccessRole() fabcore.DataAccessRoleBase {
	return fabcore.DataAccessRoleBase{
		Name: to.Ptr(testhelp.RandomName()),
		Kind: to.Ptr(fabcore.DataAccessRoleKindPolicy),
		DecisionRules: []fabcore.DecisionRule{
			{
				Effect: to.Ptr(fabcore.EffectPermit),
				Permission: []fabcore.PermissionScope{
					{
						AttributeName:            to.Ptr(fabcore.AttributeNamePath),
						AttributeValueIncludedIn: []string{"*"},
					},
					{
						AttributeName:            to.Ptr(fabcore.AttributeNameAction),
						AttributeValueIncludedIn: []string{"Read"},
					},
				},
			},
		},
		Members: &fabcore.Members{
			FabricItemMembers: []fabcore.FabricItemMember{
				{
					ItemAccess: []fabcore.ItemAccess{fabcore.ItemAccessReadAll},
					SourcePath: to.Ptr(testhelp.RandomUUID() + "/" + testhelp.RandomUUID()),
				},
			},
		},
	}
}

func NewRandomOneLakeDataAccessRoleListItem() fabcore.DataAccessRoleListItem {
	role := NewRandomOneLakeDataAccessRole()

	return fabcore.DataAccessRoleListItem{
		ID:            to.Ptr(testhelp.RandomUUID()),
		Name:          role.Name,
		Kind:          role.Kind,
		DecisionRules: role.DecisionRules,
		Members:       role.Members,
	}
}
