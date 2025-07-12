// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package fakes

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
)

// operationsOneLakeDataAccessSecurity implements SimpleIDOperations.
type operationsOneLakeDataAccessSecurity struct{}

func (o *operationsOneLakeDataAccessSecurity) TransformGet(entity fabcore.DataAccessRoles) fabcore.OneLakeDataAccessSecurityClientListDataAccessRolesResponse {
	return fabcore.OneLakeDataAccessSecurityClientListDataAccessRolesResponse{
		Etag: to.Ptr("123"),
	}
}

func transformDataAccessRole(entity fabcore.DataAccessRole) fabcore.DataAccessRole {
	return fabcore.DataAccessRole{
		ID:            entity.ID,
		Members:       entity.Members,
		Name:          entity.Name,
		DecisionRules: entity.DecisionRules,
	}
}

func (o *operationsOneLakeDataAccessSecurity) TransformList(entities []fabcore.DataAccessRole) fabcore.OneLakeDataAccessSecurityClientListDataAccessRolesResponse {
	list := make([]fabcore.DataAccessRole, len(entities))
	for i, entity := range entities {
		list[i] = transformDataAccessRole(entity)
	}

	return fabcore.OneLakeDataAccessSecurityClientListDataAccessRolesResponse{
		DataAccessRoles: fabcore.DataAccessRoles{
			Value: list,
		},
	}
}

func (o *operationsOneLakeDataAccessSecurity) GetID(entity fabcore.DataAccessRoles) string {
	return *entity.Value[0].ID
}

func configureOneLakeDataAccessSecurity(server *fakeServer) fabcore.DataAccessRoles {
	type concreteEntityOperations interface {
		identifier[fabcore.DataAccessRoles]
		getTransformer[fabcore.DataAccessRoles, fabcore.OneLakeDataAccessSecurityClientListDataAccessRolesResponse]
	}

	var entityOperations concreteEntityOperations = &operationsOneLakeDataAccessSecurity{}

	handler := newTypedHandler(server, entityOperations)

	handleGetWithParentID(handler,
		entityOperations,
		&handler.ServerFactory.Core.OneLakeDataAccessSecurityServer.ListDataAccessRoles)

	return fabcore.DataAccessRoles{}
}

func NewRandomOneLakeDataAccessSecurityClient() fabcore.DataAccessRole {
	return fabcore.DataAccessRole{
		ID:   to.Ptr(testhelp.RandomUUID()),
		Name: to.Ptr(testhelp.RandomName()),
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
					SourcePath: to.Ptr("cfafbeb1-8037-4d0c-896e-a46fb27ff222/25bac802-080d-4f73-8a42-1b406eb1fceb"),
				},
			},
		},
	}
}

func NewRandomOneLakeDataAccessesSecurityClient(itemID string) fabcore.DataAccessRoles {
	return fabcore.DataAccessRoles{
		Value: []fabcore.DataAccessRole{
			{
				ID:   to.Ptr(itemID),
				Name: to.Ptr(testhelp.RandomName()),
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
							SourcePath: to.Ptr("cfafbeb1-8037-4d0c-896e-a46fb27ff222/25bac802-080d-4f73-8a42-1b406eb1fceb"),
						},
					},
				},
			},
		},
	}
}
