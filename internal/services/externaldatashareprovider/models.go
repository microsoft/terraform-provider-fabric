// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package externaldatashareprovider

import (
	"context"

	timeoutsD "github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts" //revive:disable-line:import-alias-naming
	timeoutsR "github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"   //revive:disable-line:import-alias-naming
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
)

/*
BASE MODEL
*/

type baseExternalDataShareModel struct {
	Paths              supertypes.SetValueOf[types.String]                         `tfsdk:"paths"`
	Status             types.String                                                `tfsdk:"status"`
	Recipient          supertypes.SingleNestedObjectValueOf[recipientModel]        `tfsdk:"recipient"`
	ExpirationTime     timetypes.RFC3339                                           `tfsdk:"expiration_time"`
	PrincipalModel     supertypes.SingleNestedObjectValueOf[common.PrincipalModel] `tfsdk:"principal_model"`
	WorkspaceID        customtypes.UUID                                            `tfsdk:"workspace_id"`
	ItemID             customtypes.UUID                                            `tfsdk:"item_id"`
	ID                 customtypes.UUID                                            `tfsdk:"id"`
	InvitationURL      customtypes.URL                                             `tfsdk:"invitation_url"`
	AcceptedByTenantID customtypes.UUID                                            `tfsdk:"accepted_by_tenant_id"`
}

type recipientModel struct {
	UserPrincipalName types.String     `tfsdk:"user_principal_name"`
	TenantID          customtypes.UUID `tfsdk:"tenant_id"`
}

/*
DATA-SOURCE
*/

type dataSourceExternalDataShareProviderModel struct {
	baseExternalDataShareModel

	Timeouts timeoutsD.Value `tfsdk:"timeouts"`
}

/*
DATA-SOURCE (list)
*/

type dataSourceExternalDataSharesProviderModel struct {
	WorkspaceID customtypes.UUID                                              `tfsdk:"workspace_id"`
	ItemID      customtypes.UUID                                              `tfsdk:"item_id"`
	Values      supertypes.SetNestedObjectValueOf[baseExternalDataShareModel] `tfsdk:"values"`
	Timeouts    timeoutsD.Value                                               `tfsdk:"timeouts"`
}

/*
RESOURCE
*/

type resourceExternalDataSharesProviderModel struct {
	baseExternalDataShareModel

	Timeouts timeoutsR.Value `tfsdk:"timeouts"`
}

type requestCreateExternalDataShare struct {
	fabcore.CreateExternalDataShareRequest
}

func (to *baseExternalDataShareModel) set(ctx context.Context, workspaceID, itemID *string, from *fabcore.ExternalDataShare) diag.Diagnostics {
	to.ID = customtypes.NewUUIDPointerValue(from.ID)
	to.Status = types.StringPointerValue((*string)(from.Status))
	to.WorkspaceID = customtypes.NewUUIDPointerValue(workspaceID)
	to.ItemID = customtypes.NewUUIDPointerValue(itemID)
	to.InvitationURL = customtypes.NewURLPointerValue(from.InvitationURL)
	to.AcceptedByTenantID = customtypes.NewUUIDPointerValue(from.AcceptedByTenantID)
	to.ExpirationTime = timetypes.NewRFC3339TimePointerValue(from.ExpirationTimeUTC)
	to.PrincipalModel = supertypes.NewSingleNestedObjectValueOfNull[common.PrincipalModel](ctx)
	to.Recipient = supertypes.NewSingleNestedObjectValueOfNull[recipientModel](ctx)
	to.Paths = supertypes.NewSetValueOfNull[types.String](ctx)

	if from.Paths != nil {
		values := make([]types.String, 0, len(from.Paths))
		for _, value := range from.Paths {
			values = append(values, types.StringValue(value))
		}

		if diags := to.Paths.Set(ctx, values); diags.HasError() {
			return diags
		}
	}

	if from.Recipient != nil {
		recipient := &recipientModel{}
		recipient.set(*from.Recipient)
		if diags := to.Recipient.Set(ctx, recipient); diags.HasError() {
			return diags
		}
	}

	if from.CreatorPrincipal != nil {
		creatorPrincipalModel := &common.PrincipalModel{}
		creatorPrincipalModel.Set(*from.CreatorPrincipal)
		if diags := to.PrincipalModel.Set(ctx, creatorPrincipalModel); diags.HasError() {
			return diags
		}
	}

	return nil
}

func (to *dataSourceExternalDataSharesProviderModel) set(ctx context.Context, workspaceID, itemID *string, from []fabcore.ExternalDataShare) diag.Diagnostics {
	to.WorkspaceID = customtypes.NewUUIDPointerValue(workspaceID)
	to.ItemID = customtypes.NewUUIDPointerValue(itemID)

	slice := make([]*baseExternalDataShareModel, 0, len(from))
	for _, entity := range from {
		var entityModel baseExternalDataShareModel
		if diags := entityModel.set(ctx, workspaceID, itemID, &entity); diags.HasError() {
			return diags
		}

		slice = append(slice, &entityModel)
	}

	return to.Values.Set(ctx, slice)
}

func (to *recipientModel) set(from fabcore.ExternalDataShareRecipient) {
	to.UserPrincipalName = types.StringPointerValue(from.UserPrincipalName)
	to.TenantID = customtypes.NewUUIDPointerValue(from.TenantID)
}

func (to *requestCreateExternalDataShare) set(ctx context.Context, from resourceExternalDataSharesProviderModel) diag.Diagnostics {
	paths, diags := from.Paths.Get(ctx)

	if diags.HasError() {
		return diags
	}

	values := make([]string, 0, len(paths))

	for _, path := range paths {
		values = append(values, path.ValueString())
	}

	to.Paths = values

	recipientModel, diags := from.Recipient.Get(ctx)

	if diags.HasError() {
		return diags
	}

	to.Recipient = &fabcore.ExternalDataShareRecipient{
		UserPrincipalName: recipientModel.UserPrincipalName.ValueStringPointer(),
		TenantID:          recipientModel.TenantID.ValueStringPointer(),
	}

	return nil
}
