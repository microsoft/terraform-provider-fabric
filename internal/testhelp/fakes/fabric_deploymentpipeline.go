// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package fakes

import (
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	fabfake "github.com/microsoft/fabric-sdk-go/fabric/fake"

	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
)

// operationsDeploymentPipeline implements SimpleIDOperations.
type operationsDeploymentPipeline struct{}

func (o *operationsDeploymentPipeline) Create(data fabcore.CreateDeploymentPipelineRequest) fabcore.DeploymentPipelineExtendedInfo {
	entity := NewRandomDeploymentPipelineWithStages()
	entity.DisplayName = data.DisplayName
	entity.Description = data.Description
	entity.Stages = make([]fabcore.DeploymentPipelineStage, len(data.Stages))

	for i := range data.Stages {
		entity.Stages[i].DisplayName = data.Stages[i].DisplayName
		entity.Stages[i].Description = data.Stages[i].Description
		entity.Stages[i].IsPublic = data.Stages[i].IsPublic
	}

	return entity
}

func (o *operationsDeploymentPipeline) TransformCreate(entity fabcore.DeploymentPipelineExtendedInfo) fabcore.DeploymentPipelinesClientCreateDeploymentPipelineResponse {
	return fabcore.DeploymentPipelinesClientCreateDeploymentPipelineResponse{
		DeploymentPipelineExtendedInfo: entity,
	}
}

func (o *operationsDeploymentPipeline) TransformGet(entity fabcore.DeploymentPipelineExtendedInfo) fabcore.DeploymentPipelinesClientGetDeploymentPipelineResponse {
	return fabcore.DeploymentPipelinesClientGetDeploymentPipelineResponse{
		DeploymentPipelineExtendedInfo: entity,
	}
}

func (o *operationsDeploymentPipeline) TransformList(entities []fabcore.DeploymentPipelineExtendedInfo) fabcore.DeploymentPipelinesClientListDeploymentPipelinesResponse {
	list := make([]fabcore.DeploymentPipeline, len(entities))
	for i, entity := range entities {
		list[i] = transformDeploymentPipeline(entity)
	}

	return fabcore.DeploymentPipelinesClientListDeploymentPipelinesResponse{
		DeploymentPipelines: fabcore.DeploymentPipelines{
			Value: list,
		},
	}
}

func (o *operationsDeploymentPipeline) TransformUpdate(entity fabcore.DeploymentPipelineExtendedInfo) fabcore.DeploymentPipelinesClientUpdateDeploymentPipelineResponse {
	return fabcore.DeploymentPipelinesClientUpdateDeploymentPipelineResponse{
		DeploymentPipelineExtendedInfo: entity,
	}
}

func (o *operationsDeploymentPipeline) Update(base fabcore.DeploymentPipelineExtendedInfo, data fabcore.UpdateDeploymentPipelineRequest) fabcore.DeploymentPipelineExtendedInfo {
	base.DisplayName = data.DisplayName
	base.Description = data.Description

	return base
}

func (o *operationsDeploymentPipeline) Validate(newEntity fabcore.DeploymentPipelineExtendedInfo, existing []fabcore.DeploymentPipelineExtendedInfo) (int, error) {
	for _, entity := range existing {
		if *entity.DisplayName == *newEntity.DisplayName {
			return http.StatusConflict, fabfake.SetResponseError(http.StatusConflict, fabcore.ErrItem.ItemDisplayNameAlreadyInUse.Error(), fabcore.ErrItem.ItemDisplayNameAlreadyInUse.Error())
		}
	}

	return http.StatusCreated, nil
}

func (o *operationsDeploymentPipeline) GetID(entity fabcore.DeploymentPipelineExtendedInfo) string {
	return *entity.ID
}

func transformDeploymentPipeline(entity fabcore.DeploymentPipelineExtendedInfo) fabcore.DeploymentPipeline {
	return fabcore.DeploymentPipeline{
		ID:          entity.ID,
		DisplayName: entity.DisplayName,
		Description: entity.Description,
	}
}

func configureDeploymentPipeline(server *fakeServer) fabcore.DeploymentPipelineExtendedInfo {
	type concreteEntityOperations interface {
		simpleIDOperations[
			fabcore.DeploymentPipelineExtendedInfo,
			fabcore.DeploymentPipelinesClientGetDeploymentPipelineResponse,
			fabcore.DeploymentPipelinesClientUpdateDeploymentPipelineResponse,
			fabcore.DeploymentPipelinesClientCreateDeploymentPipelineResponse,
			fabcore.DeploymentPipelinesClientListDeploymentPipelinesResponse,
			fabcore.CreateDeploymentPipelineRequest,
			fabcore.UpdateDeploymentPipelineRequest]
	}

	var entityOperations concreteEntityOperations = &operationsDeploymentPipeline{}

	handler := newTypedHandler(server, entityOperations)

	configureEntityPagerWithSimpleID(
		handler,
		entityOperations,
		&handler.ServerFactory.Core.DeploymentPipelinesServer.GetDeploymentPipeline,
		&handler.ServerFactory.Core.DeploymentPipelinesServer.UpdateDeploymentPipeline,
		&handler.ServerFactory.Core.DeploymentPipelinesServer.CreateDeploymentPipeline,
		&handler.ServerFactory.Core.DeploymentPipelinesServer.NewListDeploymentPipelinesPager,
		&handler.ServerFactory.Core.DeploymentPipelinesServer.DeleteDeploymentPipeline)

	return fabcore.DeploymentPipelineExtendedInfo{}
}

func NewRandomDeploymentPipeline() fabcore.DeploymentPipelineExtendedInfo {
	return fabcore.DeploymentPipelineExtendedInfo{
		ID:          to.Ptr(testhelp.RandomUUID()),
		DisplayName: to.Ptr(testhelp.RandomName()),
		Description: to.Ptr(testhelp.RandomName()),
	}
}

func NewRandomDeploymentPipelineWithStages() fabcore.DeploymentPipelineExtendedInfo {
	entity := NewRandomDeploymentPipeline()
	entity.Stages = []fabcore.DeploymentPipelineStage{
		{
			DisplayName: to.Ptr(testhelp.RandomName()),
			Description: to.Ptr(testhelp.RandomName()),
			IsPublic:    to.Ptr(testhelp.RandomBool()),
		},
		{
			DisplayName: to.Ptr(testhelp.RandomName()),
			Description: to.Ptr(testhelp.RandomName()),
			IsPublic:    to.Ptr(testhelp.RandomBool()),
		},
	}

	return entity
}
