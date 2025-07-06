// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package fakes

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
)

// operationsOnelakeDataAccessSecurity implements SimpleIDOperations.
type operationsOnelakeDataAccessSecurity struct{}

func (o *operationsOnelakeDataAccessSecurity) Create(data fabcore.CreateOrUpdateDataAccessRolesRequest) fabcore.OneLakeDataAccessSecurityClientCreateOrUpdateDataAccessRolesResponse {
}

func (o *operationsOnelakeDataAccessSecurity) TransformCreate(entity fabcore.CreateOrUpdateDataAccessRolesRequest) fabcore.OneLakeDataAccessSecurityClientCreateOrUpdateDataAccessRolesResponse {
	return fabcore.OneLakeDataAccessSecurityClientCreateOrUpdateDataAccessRolesResponse{
		Etag: to.Ptr(testhelp.RandomName()),
	}
}

func (o *operationsOnelakeDataAccessSecurity) TransformGet(entity fabcore.dataaccess) fabcore.OneLakeDataAccessSecurityClientListDataAccessRolesResponse {
}

func (o *operationsOnelakeDataAccessSecurity) TransformList(entities []fabcore.DeploymentPipelineExtendedInfo) fabcore.DeploymentPipelinesClientListDeploymentPipelinesResponse {
}

func (o *operationsOnelakeDataAccessSecurity) Update(
	base fabcore.DeploymentPipelineExtendedInfo,
	data fabcore.CreateOrUpdateDataAccessRolesRequest,
) fabcore.OneLakeDataAccessSecurityClientCreateOrUpdateDataAccessRolesResponse {
	return nil
}
