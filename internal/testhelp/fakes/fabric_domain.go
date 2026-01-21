// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package fakes

import (
	"context"
	"net/http"

	azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
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

	// Handle default_label_id clearing logic (mimics real API behavior):
	// - If empty GUID is sent AND base had a label set, clear it
	// - If a real UUID is sent, set it
	if data.DefaultLabelID != nil && *data.DefaultLabelID == "00000000-0000-0000-0000-000000000000" {
		if base.DefaultLabelID != nil {
			base.DefaultLabelID = nil
		}
	} else if data.DefaultLabelID != nil {
		base.DefaultLabelID = data.DefaultLabelID
	}

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

	getDomainWrapper := func(_ context.Context, id string, _ bool, _ *fabadmin.DomainsClientGetDomainOptions) (_ azfake.Responder[fabadmin.DomainsClientGetDomainResponse], _ azfake.ErrorResponder) {
		return getByID(handler, id, entityOperations)
	}

	updateDomainWrapper := func(_ context.Context, id string, _ bool, updateRequest fabadmin.UpdateDomainRequest, _ *fabadmin.DomainsClientUpdateDomainOptions) (_ azfake.Responder[fabadmin.DomainsClientUpdateDomainResponse], _ azfake.ErrorResponder) {
		return updateByID(handler, id, updateRequest, entityOperations, entityOperations)
	}

	var createDomainWrapper func(context.Context, fabadmin.CreateDomainRequest, *fabadmin.DomainsClientCreateDomainOptions) (azfake.Responder[fabadmin.DomainsClientCreateDomainResponse], azfake.ErrorResponder)
	handleCreateWithoutWorkspace(handler, entityOperations, entityOperations, entityOperations, &createDomainWrapper)

	createDomainWrapperWith := func(ctx context.Context, _ bool, createRequest fabadmin.CreateDomainRequest, options *fabadmin.DomainsClientCreateDomainOptions) (azfake.Responder[fabadmin.DomainsClientCreateDomainResponse], azfake.ErrorResponder) {
		return createDomainWrapper(ctx, createRequest, options)
	}

	listDomainWrapper := func(_ context.Context, _ bool, options *fabadmin.DomainsClientListDomainsOptions) (_ azfake.Responder[fabadmin.DomainsClientListDomainsResponse], _ azfake.ErrorResponder) {
		return listWithFilter[fabadmin.Domain, fabadmin.DomainsClientListDomainsOptions](handler, nil, entityOperations)("", options)
	}

	deleteDomainWrapper := func(_ context.Context, id string, _ *fabadmin.DomainsClientDeleteDomainOptions) (_ azfake.Responder[fabadmin.DomainsClientDeleteDomainResponse], _ azfake.ErrorResponder) {
		return deleteByID[fabadmin.Domain, fabadmin.DomainsClientDeleteDomainResponse](handler, id)
	}

	server.ServerFactory.Admin.DomainsServer.GetDomain = getDomainWrapper
	server.ServerFactory.Admin.DomainsServer.UpdateDomain = updateDomainWrapper
	server.ServerFactory.Admin.DomainsServer.CreateDomain = createDomainWrapperWith
	server.ServerFactory.Admin.DomainsServer.ListDomains = listDomainWrapper
	server.ServerFactory.Admin.DomainsServer.DeleteDomain = deleteDomainWrapper

	return fabadmin.Domain{}
}

func NewRandomDomain() fabadmin.Domain {
	return fabadmin.Domain{
		ID:          to.Ptr(testhelp.RandomUUID()),
		DisplayName: to.Ptr(testhelp.RandomName()),
		Description: to.Ptr(testhelp.RandomName()),
	}
}

func NewRandomDomainWithDefaultLabelID() fabadmin.Domain {
	entity := NewRandomDomain()
	entity.DefaultLabelID = to.Ptr(testhelp.RandomUUID())

	return entity
}

func NewRandomDomainWithParentDomain(parentDomainID string) fabadmin.Domain {
	entity := NewRandomDomain()
	entity.ParentDomainID = to.Ptr(parentDomainID)

	return entity
}
