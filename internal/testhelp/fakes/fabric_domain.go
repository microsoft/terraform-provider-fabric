// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package fakes

import (
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	fabadmin "github.com/microsoft/fabric-sdk-go/fabric/admin"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	fabfake "github.com/microsoft/fabric-sdk-go/fabric/fake"

	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
)

// operationsDomain implements SimpleIDOperations.
type operationsDomain struct{}

// GetID implements concreteOperations.
func (o *operationsDomain) GetID(entity fabadmin.Domain) string {
	return *entity.ID
}

// TransformCreate implements concreteOperations.
func (o *operationsDomain) TransformCreate(entity fabadmin.Domain) fabadmin.DomainsClientCreateDomainResponse {
	return fabadmin.DomainsClientCreateDomainResponse{
		Domain: entity,
	}
}

func (o *operationsDomain) Create(data fabadmin.CreateDomainRequest) fabadmin.Domain {
	entity := NewRandomDomain()
	entity.DisplayName = data.DisplayName
	entity.Description = data.Description
	entity.ParentDomainID = data.ParentDomainID

	return entity
}

// TransformGet implements concreteOperations.
func (o *operationsDomain) TransformGet(entity fabadmin.Domain) fabadmin.DomainsClientGetDomainResponse {
	return fabadmin.DomainsClientGetDomainResponse{
		Domain: entity,
	}
}

// TransformList implements concreteOperations.
func (o *operationsDomain) TransformList(entities []fabadmin.Domain) fabadmin.DomainsClientListDomainsResponse {
	return fabadmin.DomainsClientListDomainsResponse{
		DomainsResponse: fabadmin.DomainsResponse{
			Domains: entities,
		},
	}
}

// TransformUpdate implements concreteOperations.
func (o *operationsDomain) TransformUpdate(entity fabadmin.Domain) fabadmin.DomainsClientUpdateDomainResponse {
	return fabadmin.DomainsClientUpdateDomainResponse{
		Domain: entity,
	}
}

// Update implements concreteOperations.
func (o *operationsDomain) Update(base fabadmin.Domain, data fabadmin.UpdateDomainRequest) fabadmin.Domain {
	base.Description = data.Description
	base.DisplayName = data.DisplayName
	base.ContributorsScope = data.ContributorsScope

	return base
}

// Validate implements concreteOperations.
func (o *operationsDomain) Validate(newEntity fabadmin.Domain, existing []fabadmin.Domain) (int, error) {
	for _, entity := range existing {
		if *entity.DisplayName == *newEntity.DisplayName {
			return http.StatusConflict, fabfake.SetResponseError(http.StatusConflict, fabcore.ErrItem.ItemDisplayNameAlreadyInUse.Error(), fabcore.ErrItem.ItemDisplayNameAlreadyInUse.Error())
		}
	}

	return http.StatusCreated, nil
}

func configureDomain(server *fakeServer) fabadmin.Domain {
	type concreteEntityOperations interface {
		simpleIDOperations[
			fabadmin.Domain,
			fabadmin.DomainsClientGetDomainResponse,
			fabadmin.DomainsClientUpdateDomainResponse,
			fabadmin.DomainsClientCreateDomainResponse,
			fabadmin.DomainsClientListDomainsResponse,
			fabadmin.CreateDomainRequest,
			fabadmin.UpdateDomainRequest,
		]
	}

	var entityOperations concreteEntityOperations = &operationsDomain{}

	handler := newTypedHandler(server, entityOperations)

	configureEntityWithSimpleID(
		handler,
		entityOperations,
		&handler.ServerFactory.Admin.DomainsServer.GetDomain,
		&handler.ServerFactory.Admin.DomainsServer.UpdateDomain,
		&handler.ServerFactory.Admin.DomainsServer.CreateDomain,
		&handler.ServerFactory.Admin.DomainsServer.ListDomains,
		&handler.ServerFactory.Admin.DomainsServer.DeleteDomain,
	)

	return fabadmin.Domain{}
}

func NewRandomDomain() fabadmin.Domain {
	return fabadmin.Domain{
		ID:                to.Ptr(testhelp.RandomUUID()),
		DisplayName:       to.Ptr(testhelp.RandomName()),
		Description:       to.Ptr(testhelp.RandomName()),
		ContributorsScope: to.Ptr(fabadmin.ContributorsScopeTypeAllTenant),
	}
}

func NewRandomDomainWithParentDomain(parentDomainID string) fabadmin.Domain {
	entity := NewRandomDomain()
	entity.ParentDomainID = to.Ptr(parentDomainID)

	return entity
}

func NewRandomDomainWithContributorsScope(contributorsScope fabadmin.ContributorsScopeType) fabadmin.Domain {
	entity := NewRandomDomain()
	entity.ContributorsScope = to.Ptr(contributorsScope)

	return entity
}
