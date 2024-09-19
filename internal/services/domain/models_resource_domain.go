// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package domain

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	fabadmin "github.com/microsoft/fabric-sdk-go/fabric/admin"
)

type resourceDomainModel struct {
	baseDomainModel
	Timeouts timeouts.Value `tfsdk:"timeouts"`
}

type requestCreateDomain struct {
	fabadmin.CreateDomainRequest
}

func (to *requestCreateDomain) set(from resourceDomainModel) {
	to.DisplayName = from.DisplayName.ValueStringPointer()
	to.Description = from.Description.ValueStringPointer()
	to.ParentDomainID = from.ParentDomainID.ValueStringPointer()
}

type requestUpdateDomain struct {
	fabadmin.UpdateDomainRequest
}

func (to *requestUpdateDomain) set(from resourceDomainModel) {
	to.DisplayName = from.DisplayName.ValueStringPointer()
	to.Description = from.Description.ValueStringPointer()
	to.ContributorsScope = (*fabadmin.ContributorsScopeType)(from.ContributorsScope.ValueStringPointer())
}
