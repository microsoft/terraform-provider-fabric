// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package fakes

import (
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	fabdeploymentplan "github.com/microsoft/fabric-sdk-go/fabric/deploymentplan"
	fabfake "github.com/microsoft/fabric-sdk-go/fabric/fake"

	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
)

type operationsDeploymentPlan struct{}

// CreateDefinition implements concreteDefinitionOperations.
func (o *operationsDeploymentPlan) CreateDefinition(data fabdeploymentplan.CreateDeploymentPlanRequest) *fabdeploymentplan.PublicDefinition {
	return data.Definition
}

// TransformDefinition implements concreteDefinitionOperations.
func (o *operationsDeploymentPlan) TransformDefinition(entity *fabdeploymentplan.PublicDefinition) fabdeploymentplan.ItemsClientGetDeploymentPlanDefinitionResponse {
	return fabdeploymentplan.ItemsClientGetDeploymentPlanDefinitionResponse{
		Definition: entity,
	}
}

// UpdateDefinition implements concreteDefinitionOperations.
func (o *operationsDeploymentPlan) UpdateDefinition(_ *fabdeploymentplan.PublicDefinition, data fabdeploymentplan.UpdateDeploymentPlanDefinitionRequest) *fabdeploymentplan.PublicDefinition {
	return data.Definition
}

// CreateWithParentID implements concreteOperations.
func (o *operationsDeploymentPlan) CreateWithParentID(parentID string, data fabdeploymentplan.CreateDeploymentPlanRequest) fabdeploymentplan.DeploymentPlan {
	entity := NewRandomDeploymentPlanWithWorkspace(parentID)
	entity.DisplayName = data.DisplayName
	entity.Description = data.Description
	entity.FolderID = data.FolderID

	return entity
}

// Filter implements concreteOperations.
func (o *operationsDeploymentPlan) Filter(entities []fabdeploymentplan.DeploymentPlan, parentID string) []fabdeploymentplan.DeploymentPlan {
	ret := make([]fabdeploymentplan.DeploymentPlan, 0)

	for _, entity := range entities {
		if *entity.WorkspaceID == parentID {
			ret = append(ret, entity)
		}
	}

	return ret
}

// GetID implements concreteOperations.
func (o *operationsDeploymentPlan) GetID(entity fabdeploymentplan.DeploymentPlan) string {
	return generateID(*entity.WorkspaceID, *entity.ID)
}

// TransformCreate implements concreteOperations.
func (o *operationsDeploymentPlan) TransformCreate(entity fabdeploymentplan.DeploymentPlan) fabdeploymentplan.ItemsClientCreateDeploymentPlanResponse {
	return fabdeploymentplan.ItemsClientCreateDeploymentPlanResponse{
		DeploymentPlan: entity,
	}
}

// TransformGet implements concreteOperations.
func (o *operationsDeploymentPlan) TransformGet(entity fabdeploymentplan.DeploymentPlan) fabdeploymentplan.ItemsClientGetDeploymentPlanResponse {
	return fabdeploymentplan.ItemsClientGetDeploymentPlanResponse{
		DeploymentPlan: entity,
	}
}

// TransformList implements concreteOperations.
func (o *operationsDeploymentPlan) TransformList(entities []fabdeploymentplan.DeploymentPlan) fabdeploymentplan.ItemsClientListDeploymentPlansResponse {
	return fabdeploymentplan.ItemsClientListDeploymentPlansResponse{
		Value: entities,
	}
}

// TransformUpdate implements concreteOperations.
func (o *operationsDeploymentPlan) TransformUpdate(entity fabdeploymentplan.DeploymentPlan) fabdeploymentplan.ItemsClientUpdateDeploymentPlanResponse {
	return fabdeploymentplan.ItemsClientUpdateDeploymentPlanResponse{
		DeploymentPlan: entity,
	}
}

// Update implements concreteOperations.
func (o *operationsDeploymentPlan) Update(base fabdeploymentplan.DeploymentPlan, data fabdeploymentplan.UpdateDeploymentPlanRequest) fabdeploymentplan.DeploymentPlan {
	base.Description = data.Description
	base.DisplayName = data.DisplayName

	return base
}

// Validate implements concreteOperations.
func (o *operationsDeploymentPlan) Validate(newEntity fabdeploymentplan.DeploymentPlan, existing []fabdeploymentplan.DeploymentPlan) (int, error) {
	for _, entity := range existing {
		if *entity.WorkspaceID == *newEntity.WorkspaceID && *entity.DisplayName == *newEntity.DisplayName {
			return http.StatusConflict, fabfake.SetResponseError(http.StatusConflict, fabcore.ErrItem.ItemDisplayNameAlreadyInUse.Error(), fabcore.ErrItem.ItemDisplayNameAlreadyInUse.Error())
		}
	}

	return http.StatusCreated, nil
}

func configureDeploymentPlan(server *fakeServer) fabdeploymentplan.DeploymentPlan {
	type concreteEntityOperations interface {
		parentIDOperations[
			fabdeploymentplan.DeploymentPlan,
			fabdeploymentplan.ItemsClientGetDeploymentPlanResponse,
			fabdeploymentplan.ItemsClientUpdateDeploymentPlanResponse,
			fabdeploymentplan.ItemsClientCreateDeploymentPlanResponse,
			fabdeploymentplan.ItemsClientListDeploymentPlansResponse,
			fabdeploymentplan.CreateDeploymentPlanRequest,
			fabdeploymentplan.UpdateDeploymentPlanRequest]
	}

	type concreteDefinitionOperations interface {
		definitionOperations[
			fabdeploymentplan.PublicDefinition,
			fabdeploymentplan.CreateDeploymentPlanRequest,
			fabdeploymentplan.UpdateDeploymentPlanDefinitionRequest,
			fabdeploymentplan.ItemsClientGetDeploymentPlanDefinitionResponse,
			fabdeploymentplan.ItemsClientUpdateDeploymentPlanDefinitionResponse]
	}

	var entityOperations concreteEntityOperations = &operationsDeploymentPlan{}

	var definitionOperations concreteDefinitionOperations = &operationsDeploymentPlan{}

	handler := newTypedHandler(server, entityOperations)

	// The fabric-sdk-go fake ServerFactory exposes the Deployment Plan servers under the "Deployment" field.
	configureEntityWithParentID(
		handler,
		entityOperations,
		&server.ServerFactory.Deployment.ItemsServer.GetDeploymentPlan,
		&server.ServerFactory.Deployment.ItemsServer.UpdateDeploymentPlan,
		&server.ServerFactory.Deployment.ItemsServer.BeginCreateDeploymentPlan,
		&server.ServerFactory.Deployment.ItemsServer.NewListDeploymentPlansPager,
		&server.ServerFactory.Deployment.ItemsServer.DeleteDeploymentPlan,
	)

	configureDefinitions(
		handler,
		entityOperations,
		definitionOperations,
		&server.ServerFactory.Deployment.ItemsServer.BeginCreateDeploymentPlan,
		&server.ServerFactory.Deployment.ItemsServer.BeginGetDeploymentPlanDefinition,
		&server.ServerFactory.Deployment.ItemsServer.BeginUpdateDeploymentPlanDefinition,
	)

	return fabdeploymentplan.DeploymentPlan{}
}

func NewRandomDeploymentPlan() fabdeploymentplan.DeploymentPlan {
	return fabdeploymentplan.DeploymentPlan{
		ID:          new(testhelp.RandomUUID()),
		DisplayName: new(testhelp.RandomName()),
		Description: new(testhelp.RandomName()),
		WorkspaceID: new(testhelp.RandomUUID()),
		FolderID:    new(testhelp.RandomUUID()),
		Type:        to.Ptr(fabdeploymentplan.ItemTypeDeploymentPlan),
	}
}

func NewRandomDeploymentPlanWithWorkspace(workspaceID string) fabdeploymentplan.DeploymentPlan {
	result := NewRandomDeploymentPlan()
	result.WorkspaceID = &workspaceID

	return result
}

func NewRandomDeploymentPlanDefinition() fabdeploymentplan.PublicDefinition {
	defPart := fabdeploymentplan.PublicDefinitionPart{
		PayloadType: to.Ptr(fabdeploymentplan.PayloadTypeInlineBase64),
		Path:        new("plan.yml"),
		Payload: new(
			"JHNjaGVtYTogaHR0cHM6Ly9kZXZlbG9wZXIubWljcm9zb2Z0LmNvbS9qc29uLXNjaGVtYXMvZmFicmljL2l0ZW0vZGVwbG95bWVudFBsYW4vZGVmaW5pdGlvbi9wbGFuLzEuMC4wL3NjaGVtYS5qc29uCnZlcnNpb246IDEuMC4wCmdyb3VwczogW10K", // plan.yml with empty groups in base64
		),
	}

	return fabdeploymentplan.PublicDefinition{
		Format: to.Ptr(fabdeploymentplan.DefinitionFormatDeploymentPlanFormat),
		Parts:  []fabdeploymentplan.PublicDefinitionPart{defPart},
	}
}
