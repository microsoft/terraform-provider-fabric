// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package domain

import (
	"context"

	timeoutsD "github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts" //revive:disable-line:import-alias-naming
	timeoutsR "github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"   //revive:disable-line:import-alias-naming
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	fabadmin "github.com/microsoft/fabric-sdk-go/fabric/admin"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
)

/*
BASE MODEL
*/

type baseDomainModel struct {
	ID                customtypes.UUID `tfsdk:"id"`
	DisplayName       types.String     `tfsdk:"display_name"`
	Description       types.String     `tfsdk:"description"`
	ParentDomainID    customtypes.UUID `tfsdk:"parent_domain_id"`
	ContributorsScope types.String     `tfsdk:"contributors_scope"`
}

func (to *baseDomainModel) set(from fabadmin.DomainPreview) {
	to.ID = customtypes.NewUUIDPointerValue(from.ID)
	to.DisplayName = types.StringPointerValue(from.DisplayName)
	to.Description = types.StringPointerValue(from.Description)
	to.ParentDomainID = customtypes.NewUUIDPointerValue(from.ParentDomainID)
	to.ContributorsScope = types.StringPointerValue((*string)(from.ContributorsScope))
}

/*
DATA-SOURCE
*/

type dataSourceDomainModel struct {
	baseDomainModel

	Timeouts timeoutsD.Value `tfsdk:"timeouts"`
}

/*
DATA-SOURCE (list)
*/

type dataSourceDomainsModel struct {
	Values   supertypes.SetNestedObjectValueOf[baseDomainModel] `tfsdk:"values"`
	Timeouts timeoutsD.Value                                    `tfsdk:"timeouts"`
}

func (to *dataSourceDomainsModel) setValues(ctx context.Context, from []fabadmin.DomainPreview) diag.Diagnostics {
	slice := make([]*baseDomainModel, 0, len(from))

	for _, entity := range from {
		var entityModel baseDomainModel
		entityModel.set(entity)
		slice = append(slice, &entityModel)
	}

	return to.Values.Set(ctx, slice)
}

/*
RESOURCE
*/

type resourceDomainModel struct {
	baseDomainModel

	Timeouts timeoutsR.Value `tfsdk:"timeouts"`
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
	fabadmin.UpdateDomainRequestPreview
}

func (to *requestUpdateDomain) set(from resourceDomainModel) {
	to.DisplayName = from.DisplayName.ValueStringPointer()
	to.Description = from.Description.ValueStringPointer()
	to.ContributorsScope = (*fabadmin.ContributorsScopeType)(from.ContributorsScope.ValueStringPointer())
}
