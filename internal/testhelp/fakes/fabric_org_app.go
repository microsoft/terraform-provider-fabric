// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package fakes

import (
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	fabfake "github.com/microsoft/fabric-sdk-go/fabric/fake"
	faborgapp "github.com/microsoft/fabric-sdk-go/fabric/orgapp"

	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
)

type operationsOrgApp struct{}

// ConvertItemToEntity implements itemConverter.
func (o *operationsOrgApp) ConvertItemToEntity(item fabcore.Item) faborgapp.OrgApp {
	return faborgapp.OrgApp{
		ID:          item.ID,
		DisplayName: item.DisplayName,
		Description: item.Description,
		WorkspaceID: item.WorkspaceID,
		FolderID:    item.FolderID,
		Type:        to.Ptr(faborgapp.ItemTypeOrgApp),
		Properties:  NewRandomOrgApp().Properties,
		Tags:        convertItemTags[faborgapp.ItemTag](item.Tags),
	}
}

// CreateDefinition implements concreteDefinitionOperations.
func (o *operationsOrgApp) CreateDefinition(data faborgapp.CreateOrgAppRequest) *faborgapp.PublicDefinition {
	return data.Definition
}

// TransformDefinition implements concreteDefinitionOperations.
func (o *operationsOrgApp) TransformDefinition(entity *faborgapp.PublicDefinition) faborgapp.ItemsClientGetOrgAppDefinitionResponse {
	return faborgapp.ItemsClientGetOrgAppDefinitionResponse{
		Definition: entity,
	}
}

// UpdateDefinition implements concreteDefinitionOperations.
func (o *operationsOrgApp) UpdateDefinition(_ *faborgapp.PublicDefinition, data faborgapp.UpdateOrgAppDefinitionRequest) *faborgapp.PublicDefinition {
	return data.Definition
}

// CreateWithParentID implements concreteOperations.
func (o *operationsOrgApp) CreateWithParentID(parentID string, data faborgapp.CreateOrgAppRequest) faborgapp.OrgApp {
	entity := NewRandomOrgAppWithWorkspace(parentID)
	entity.DisplayName = data.DisplayName
	entity.Description = data.Description
	entity.FolderID = data.FolderID

	return entity
}

// Filter implements concreteOperations.
func (o *operationsOrgApp) Filter(entities []faborgapp.OrgApp, parentID string) []faborgapp.OrgApp {
	ret := make([]faborgapp.OrgApp, 0)

	for _, entity := range entities {
		if *entity.WorkspaceID == parentID {
			ret = append(ret, entity)
		}
	}

	return ret
}

// GetID implements concreteOperations.
func (o *operationsOrgApp) GetID(entity faborgapp.OrgApp) string {
	return generateID(*entity.WorkspaceID, *entity.ID)
}

// TransformCreate implements concreteOperations.
func (o *operationsOrgApp) TransformCreate(entity faborgapp.OrgApp) faborgapp.ItemsClientCreateOrgAppResponse {
	return faborgapp.ItemsClientCreateOrgAppResponse{
		OrgApp: entity,
	}
}

// TransformGet implements concreteOperations.
func (o *operationsOrgApp) TransformGet(entity faborgapp.OrgApp) faborgapp.ItemsClientGetOrgAppResponse {
	return faborgapp.ItemsClientGetOrgAppResponse{
		OrgApp: entity,
	}
}

// TransformList implements concreteOperations.
func (o *operationsOrgApp) TransformList(entities []faborgapp.OrgApp) faborgapp.ItemsClientListOrgAppsResponse {
	return faborgapp.ItemsClientListOrgAppsResponse{
		Value: entities,
	}
}

// TransformUpdate implements concreteOperations.
func (o *operationsOrgApp) TransformUpdate(entity faborgapp.OrgApp) faborgapp.ItemsClientUpdateOrgAppResponse {
	return faborgapp.ItemsClientUpdateOrgAppResponse{
		OrgApp: entity,
	}
}

// Update implements concreteOperations.
func (o *operationsOrgApp) Update(base faborgapp.OrgApp, data faborgapp.UpdateOrgAppRequest) faborgapp.OrgApp {
	base.Description = data.Description
	base.DisplayName = data.DisplayName

	return base
}

// Validate implements concreteOperations.
func (o *operationsOrgApp) Validate(newEntity faborgapp.OrgApp, existing []faborgapp.OrgApp) (int, error) {
	for _, entity := range existing {
		if *entity.DisplayName == *newEntity.DisplayName {
			return http.StatusConflict, fabfake.SetResponseError(http.StatusConflict, fabcore.ErrItem.ItemDisplayNameAlreadyInUse.Error(), fabcore.ErrItem.ItemDisplayNameAlreadyInUse.Error())
		}
	}

	return http.StatusCreated, nil
}

func configureOrgApp(server *fakeServer) faborgapp.OrgApp {
	type concreteEntityOperations interface {
		parentIDOperations[
			faborgapp.OrgApp,
			faborgapp.ItemsClientGetOrgAppResponse,
			faborgapp.ItemsClientUpdateOrgAppResponse,
			faborgapp.ItemsClientCreateOrgAppResponse,
			faborgapp.ItemsClientListOrgAppsResponse,
			faborgapp.CreateOrgAppRequest,
			faborgapp.UpdateOrgAppRequest]
	}
	type concreteDefinitionOperations interface {
		definitionOperations[
			faborgapp.PublicDefinition,
			faborgapp.CreateOrgAppRequest,
			faborgapp.UpdateOrgAppDefinitionRequest,
			faborgapp.ItemsClientGetOrgAppDefinitionResponse,
			faborgapp.ItemsClientUpdateOrgAppDefinitionResponse]
	}

	var entityOperations concreteEntityOperations = &operationsOrgApp{}
	var definitionOperations concreteDefinitionOperations = &operationsOrgApp{}
	var converter itemConverter[faborgapp.OrgApp] = &operationsOrgApp{}

	handler := newTypedHandlerWithConverter(server, entityOperations, converter)

	configureEntityWithParentID(
		handler,
		entityOperations,
		&server.ServerFactory.OrgApp.ItemsServer.GetOrgApp,
		&server.ServerFactory.OrgApp.ItemsServer.UpdateOrgApp,
		&server.ServerFactory.OrgApp.ItemsServer.BeginCreateOrgApp,
		&server.ServerFactory.OrgApp.ItemsServer.NewListOrgAppsPager,
		&server.ServerFactory.OrgApp.ItemsServer.DeleteOrgApp,
	)
	configureDefinitions(
		handler,
		entityOperations,
		definitionOperations,
		&server.ServerFactory.OrgApp.ItemsServer.BeginCreateOrgApp,
		&server.ServerFactory.OrgApp.ItemsServer.BeginGetOrgAppDefinition,
		&server.ServerFactory.OrgApp.ItemsServer.BeginUpdateOrgAppDefinition,
	)

	return faborgapp.OrgApp{}
}

func NewRandomOrgApp() faborgapp.OrgApp {
	return faborgapp.OrgApp{
		ID:          new(testhelp.RandomUUID()),
		DisplayName: new(testhelp.RandomName()),
		Description: new(testhelp.RandomName()),
		WorkspaceID: new(testhelp.RandomUUID()),
		FolderID:    new(testhelp.RandomUUID()),
		Type:        to.Ptr(faborgapp.ItemTypeOrgApp),
		Properties: &faborgapp.Properties{
			OneLakeRootPath: new(testhelp.RandomURI()),
		},
	}
}

func NewRandomOrgAppWithWorkspace(workspaceID string) faborgapp.OrgApp {
	result := NewRandomOrgApp()
	result.WorkspaceID = &workspaceID

	return result
}

func NewRandomOrgAppDefinition() faborgapp.PublicDefinition {
	defPart := faborgapp.PublicDefinitionPart{
		PayloadType: to.Ptr(faborgapp.PayloadTypeInlineBase64),
		Path:        new("definition.json"),
		Payload: new(
			"eyJjb250ZW50IjoiSGVsbG8gV29ybGQifQ==", // {"content":"Hello World"} in base64
		),
	}

	defParts := make([]faborgapp.PublicDefinitionPart, 0, 1)

	defParts = append(defParts, defPart)

	return faborgapp.PublicDefinition{
		Format: new("JSON"),
		Parts:  defParts,
	}
}
