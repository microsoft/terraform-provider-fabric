// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package fakes

import (
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	fabfake "github.com/microsoft/fabric-sdk-go/fabric/fake"
	fabreport "github.com/microsoft/fabric-sdk-go/fabric/report"

	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
)

type operationsReport struct{}

// CreateDefinition implements concreteDefinitionOperations.
func (o *operationsReport) CreateDefinition(data fabreport.CreateReportRequest) *fabreport.Definition {
	return data.Definition
}

// TransformDefinition implements concreteDefinitionOperations.
func (o *operationsReport) TransformDefinition(entity *fabreport.Definition) fabreport.ItemsClientGetReportDefinitionResponse {
	return fabreport.ItemsClientGetReportDefinitionResponse{
		DefinitionResponse: fabreport.DefinitionResponse{
			Definition: entity,
		},
	}
}

// UpdateDefinition implements concreteDefinitionOperations.
func (o *operationsReport) UpdateDefinition(_ *fabreport.Definition, data fabreport.UpdateReportDefinitionRequest) *fabreport.Definition {
	return data.Definition
}

// CreateWithParentID implements concreteOperations.
func (o *operationsReport) CreateWithParentID(parentID string, data fabreport.CreateReportRequest) fabreport.Report {
	result := NewRandomReportWithWorkspace(parentID)
	result.DisplayName = data.DisplayName
	result.Description = data.Description
	result.FolderID = data.FolderID

	return result
}

// Filter implements concreteOperations.
func (o *operationsReport) Filter(entities []fabreport.Report, parentID string) []fabreport.Report {
	ret := make([]fabreport.Report, 0)

	for _, entity := range entities {
		if *entity.WorkspaceID == parentID {
			ret = append(ret, entity)
		}
	}

	return ret
}

// GetID implements concreteOperations.
func (o *operationsReport) GetID(entity fabreport.Report) string {
	return generateID(*entity.WorkspaceID, *entity.ID)
}

// TransformCreate implements concreteOperations.
func (o *operationsReport) TransformCreate(entity fabreport.Report) fabreport.ItemsClientCreateReportResponse {
	return fabreport.ItemsClientCreateReportResponse{
		Report: entity,
	}
}

// TransformGet implements concreteOperations.
func (o *operationsReport) TransformGet(entity fabreport.Report) fabreport.ItemsClientGetReportResponse {
	return fabreport.ItemsClientGetReportResponse{
		Report: entity,
	}
}

// TransformList implements concreteOperations.
func (o *operationsReport) TransformList(entities []fabreport.Report) fabreport.ItemsClientListReportsResponse {
	return fabreport.ItemsClientListReportsResponse{
		Reports: fabreport.Reports{
			Value: entities,
		},
	}
}

// TransformUpdate implements concreteOperations.
func (o *operationsReport) TransformUpdate(entity fabreport.Report) fabreport.ItemsClientUpdateReportResponse {
	return fabreport.ItemsClientUpdateReportResponse{
		Report: entity,
	}
}

// Update implements concreteOperations.
func (o *operationsReport) Update(base fabreport.Report, data fabreport.UpdateReportRequest) fabreport.Report {
	base.Description = data.Description
	base.DisplayName = data.DisplayName

	return base
}

// Validate implements concreteOperations.
func (o *operationsReport) Validate(newEntity fabreport.Report, existing []fabreport.Report) (int, error) {
	for _, entity := range existing {
		if *entity.DisplayName == *newEntity.DisplayName {
			return http.StatusConflict, fabfake.SetResponseError(http.StatusConflict, fabcore.ErrItem.ItemDisplayNameAlreadyInUse.Error(), fabcore.ErrItem.ItemDisplayNameAlreadyInUse.Error())
		}
	}

	return http.StatusCreated, nil
}

func configureReport(server *fakeServer) fabreport.Report {
	type concreteEntityOperations interface {
		parentIDOperations[
			fabreport.Report,
			fabreport.ItemsClientGetReportResponse,
			fabreport.ItemsClientUpdateReportResponse,
			fabreport.ItemsClientCreateReportResponse,
			fabreport.ItemsClientListReportsResponse,
			fabreport.CreateReportRequest,
			fabreport.UpdateReportRequest]
	}

	type concreteDefinitionOperations interface {
		definitionOperations[
			fabreport.Definition,
			fabreport.CreateReportRequest,
			fabreport.UpdateReportDefinitionRequest,
			fabreport.ItemsClientGetReportDefinitionResponse,
			fabreport.ItemsClientUpdateReportDefinitionResponse]
	}

	var entityOperations concreteEntityOperations = &operationsReport{}

	var definitionOperations concreteDefinitionOperations = &operationsReport{}

	handler := newTypedHandler(server, entityOperations)

	configureEntityWithParentID(
		handler,
		entityOperations,
		&server.ServerFactory.Report.ItemsServer.GetReport,
		&server.ServerFactory.Report.ItemsServer.UpdateReport,
		&server.ServerFactory.Report.ItemsServer.BeginCreateReport,
		&server.ServerFactory.Report.ItemsServer.NewListReportsPager,
		&server.ServerFactory.Report.ItemsServer.DeleteReport)

	configureDefinitions(
		handler,
		entityOperations,
		definitionOperations,
		&server.ServerFactory.Report.ItemsServer.BeginCreateReport,
		&server.ServerFactory.Report.ItemsServer.BeginGetReportDefinition,
		&server.ServerFactory.Report.ItemsServer.BeginUpdateReportDefinition)

	return fabreport.Report{}
}

func NewRandomReport() fabreport.Report {
	return fabreport.Report{
		ID:          to.Ptr(testhelp.RandomUUID()),
		DisplayName: to.Ptr(testhelp.RandomName()),
		Description: to.Ptr(testhelp.RandomName()),
		WorkspaceID: to.Ptr(testhelp.RandomUUID()),
		FolderID:    to.Ptr(testhelp.RandomUUID()),
		Type:        to.Ptr(fabreport.ItemTypeReport),
	}
}

func NewRandomReportWithWorkspace(workspaceID string) fabreport.Report {
	result := NewRandomReport()
	result.WorkspaceID = &workspaceID

	return result
}
