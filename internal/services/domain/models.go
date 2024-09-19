// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package domain

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	fabadmin "github.com/microsoft/fabric-sdk-go/fabric/admin"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
)

type baseDomainModel struct {
	ID                customtypes.UUID `tfsdk:"id"`
	DisplayName       types.String     `tfsdk:"display_name"`
	Description       types.String     `tfsdk:"description"`
	ParentDomainID    customtypes.UUID `tfsdk:"parent_domain_id"`
	ContributorsScope types.String     `tfsdk:"contributors_scope"`
}

func (to *baseDomainModel) set(from fabadmin.Domain) {
	to.ID = customtypes.NewUUIDPointerValue(from.ID)
	to.DisplayName = types.StringPointerValue(from.DisplayName)
	to.Description = types.StringPointerValue(from.Description)
	to.ParentDomainID = customtypes.NewUUIDPointerValue(from.ParentDomainID)
	to.ContributorsScope = types.StringPointerValue((*string)(from.ContributorsScope))
}

type workspaceModel struct {
	ID          customtypes.UUID `tfsdk:"id"`
	DisplayName types.String     `tfsdk:"display_name"`
}

func (to *workspaceModel) set(from fabadmin.DomainWorkspace) {
	to.ID = customtypes.NewUUIDPointerValue(from.ID)
	to.DisplayName = types.StringPointerValue(from.DisplayName)
}
