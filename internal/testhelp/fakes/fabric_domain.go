// Copyright Microsoft Corporation 2026
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
func (o *operationsDomain) GetID(entity fabadmin.DomainPreview) string {
	return *entity.ID
}

// TransformCreate implements concreteOperations.
func (o *operationsDomain) TransformCreate(entity fabadmin.DomainPreview) fabadmin.DomainsClientCreateDomainPreviewResponse {
	return fabadmin.DomainsClientCreateDomainPreviewResponse{
		DomainPreview: entity,
	}
}

func (o *operationsDomain) Create(data fabadmin.CreateDomainRequest) fabadmin.DomainPreview {
	entity := NewRandomDomain()
	entity.DisplayName = data.DisplayName
	entity.Description = data.Description
	entity.ParentDomainID = data.ParentDomainID

	return entity
}

// TransformGet implements concreteOperations.
func (o *operationsDomain) TransformGet(entity fabadmin.DomainPreview) fabadmin.DomainsClientGetDomainPreviewResponse {
	return fabadmin.DomainsClientGetDomainPreviewResponse{
		DomainPreview: entity,
	}
}

// TransformList implements concreteOperations.
func (o *operationsDomain) TransformList(entities []fabadmin.DomainPreview) fabadmin.DomainsClientListDomainsPreviewResponse {
	return fabadmin.DomainsClientListDomainsPreviewResponse{
		DomainsResponsePreview: fabadmin.DomainsResponsePreview{
			Domains: entities,
		},
	}
}

// TransformUpdate implements concreteOperations.
func (o *operationsDomain) TransformUpdate(entity fabadmin.DomainPreview) fabadmin.DomainsClientUpdateDomainPreviewResponse {
	return fabadmin.DomainsClientUpdateDomainPreviewResponse{
		DomainPreview: entity,
	}
}

// Update implements concreteOperations.
func (o *operationsDomain) Update(base fabadmin.DomainPreview, data fabadmin.UpdateDomainRequestPreview) fabadmin.DomainPreview {
	base.Description = data.Description
	base.DisplayName = data.DisplayName
	base.ContributorsScope = data.ContributorsScope

	return base
}

// Validate implements concreteOperations.
func (o *operationsDomain) Validate(newEntity fabadmin.DomainPreview, existing []fabadmin.DomainPreview) (int, error) {
	for _, entity := range existing {
		if *entity.DisplayName == *newEntity.DisplayName {
			return http.StatusConflict, fabfake.SetResponseError(http.StatusConflict, fabcore.ErrItem.ItemDisplayNameAlreadyInUse.Error(), fabcore.ErrItem.ItemDisplayNameAlreadyInUse.Error())
		}
	}

	return http.StatusCreated, nil
}

func configureDomain(server *fakeServer) fabadmin.DomainPreview {
	type concreteEntityOperations interface {
		simpleIDOperations[
			fabadmin.DomainPreview,
			fabadmin.DomainsClientGetDomainPreviewResponse,
			fabadmin.DomainsClientUpdateDomainPreviewResponse,
			fabadmin.DomainsClientCreateDomainPreviewResponse,
			fabadmin.DomainsClientListDomainsPreviewResponse,
			fabadmin.CreateDomainRequest,
			fabadmin.UpdateDomainRequestPreview,
		]
	}

	var entityOperations concreteEntityOperations = &operationsDomain{}

	handler := newTypedHandler(server, entityOperations)

	getDomainWrapper := func(_ context.Context, id string, _ bool, _ *fabadmin.DomainsClientGetDomainPreviewOptions) (_ azfake.Responder[fabadmin.DomainsClientGetDomainPreviewResponse], _ azfake.ErrorResponder) {
		return getByID(handler, id, entityOperations)
	}

	updateDomainWrapper := func(_ context.Context, id string, _ bool, updateRequest fabadmin.UpdateDomainRequestPreview, _ *fabadmin.DomainsClientUpdateDomainPreviewOptions) (_ azfake.Responder[fabadmin.DomainsClientUpdateDomainPreviewResponse], _ azfake.ErrorResponder) {
		return updateByID(handler, id, updateRequest, entityOperations, entityOperations)
	}

	var createDomainWrapper func(context.Context, fabadmin.CreateDomainRequest, *fabadmin.DomainsClientCreateDomainPreviewOptions) (azfake.Responder[fabadmin.DomainsClientCreateDomainPreviewResponse], azfake.ErrorResponder)
	handleCreateWithoutWorkspace(handler, entityOperations, entityOperations, entityOperations, &createDomainWrapper)

	createDomainWrapperWithPreview := func(ctx context.Context, _ bool, createRequest fabadmin.CreateDomainRequest, options *fabadmin.DomainsClientCreateDomainPreviewOptions) (azfake.Responder[fabadmin.DomainsClientCreateDomainPreviewResponse], azfake.ErrorResponder) {
		return createDomainWrapper(ctx, createRequest, options)
	}

	listDomainWrapper := func(_ context.Context, _ bool, options *fabadmin.DomainsClientListDomainsPreviewOptions) (_ azfake.Responder[fabadmin.DomainsClientListDomainsPreviewResponse], _ azfake.ErrorResponder) {
		return listWithFilter[fabadmin.DomainPreview, fabadmin.DomainsClientListDomainsPreviewOptions](handler, nil, entityOperations)("", options)
	}

	deleteDomainWrapper := func(_ context.Context, id string, _ *fabadmin.DomainsClientDeleteDomainOptions) (_ azfake.Responder[fabadmin.DomainsClientDeleteDomainResponse], _ azfake.ErrorResponder) {
		return deleteByID[fabadmin.DomainPreview, fabadmin.DomainsClientDeleteDomainResponse](handler, id)
	}

	server.ServerFactory.Admin.DomainsServer.GetDomainPreview = getDomainWrapper
	server.ServerFactory.Admin.DomainsServer.UpdateDomainPreview = updateDomainWrapper
	server.ServerFactory.Admin.DomainsServer.CreateDomainPreview = createDomainWrapperWithPreview
	server.ServerFactory.Admin.DomainsServer.ListDomainsPreview = listDomainWrapper
	server.ServerFactory.Admin.DomainsServer.DeleteDomain = deleteDomainWrapper

	return fabadmin.DomainPreview{}
}

func NewRandomDomain() fabadmin.DomainPreview {
	return fabadmin.DomainPreview{
		ID:                to.Ptr(testhelp.RandomUUID()),
		DisplayName:       to.Ptr(testhelp.RandomName()),
		Description:       to.Ptr(testhelp.RandomName()),
		ContributorsScope: to.Ptr(fabadmin.ContributorsScopeTypeAllTenant),
	}
}

func NewRandomDomainWithParentDomain(parentDomainID string) fabadmin.DomainPreview {
	entity := NewRandomDomain()
	entity.ParentDomainID = to.Ptr(parentDomainID)

	return entity
}

func NewRandomDomainWithContributorsScope(contributorsScope fabadmin.ContributorsScopeType) fabadmin.DomainPreview {
	entity := NewRandomDomain()
	entity.ContributorsScope = to.Ptr(contributorsScope)

	return entity
}
