// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package externaldatasharesprovider

import (
	"context"
	"time"

	timeoutsD "github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts" //revive:disable-line:import-alias-naming
	timeoutsR "github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"   //revive:disable-line:import-alias-naming
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
)

/*
BASE MODEL
*/

type externalDataSharesModel struct {
	Paths              supertypes.SetValueOf[types.String]                         `tfsdk:"paths"`
	Status             types.String                                                `tfsdk:"status"`
	Recipient          supertypes.SingleNestedObjectValueOf[recipientModel]        `tfsdk:"recipient"`
	ExpirationTimeUtc  types.String                                                `tfsdk:"expiration_time_utc"`
	CreatorPrincipal   supertypes.SingleNestedObjectValueOf[creatorPrincipalModel] `tfsdk:"creator_principal"`
	WorkspaceID        customtypes.UUID                                            `tfsdk:"workspace_id"`
	ItemID             customtypes.UUID                                            `tfsdk:"item_id"`
	ID                 customtypes.UUID                                            `tfsdk:"id"`
	InvitationURL      types.String                                                `tfsdk:"invitation_url"`
	AcceptedByTenantID customtypes.UUID                                            `tfsdk:"accepted_by_tenant_id"`
}

type creatorPrincipalModel struct {
	ID          customtypes.UUID                                       `tfsdk:"id"`
	DisplayName types.String                                           `tfsdk:"display_name"`
	Type        types.String                                           `tfsdk:"type"`
	UserDetails supertypes.SingleNestedObjectValueOf[userDetailsModel] `tfsdk:"user_details"`
}

type recipientModel struct {
	UserPrincipalName types.String     `tfsdk:"user_principal_name"`
	TenantID          customtypes.UUID `tfsdk:"tenant_id"`
}

type userDetailsModel struct {
	UserPrincipalName types.String `tfsdk:"user_principal_name"`
}

/*
DATA-SOURCE
*/

type dataSourceExternalDataShareProviderModel struct {
	externalDataSharesModel

	Timeouts timeoutsD.Value `tfsdk:"timeouts"`
}

/*
DATA-SOURCE (list)
*/

type dataSourceExternalDataSharesProviderModel struct {
	WorkspaceID customtypes.UUID                                           `tfsdk:"workspace_id"`
	ItemID      customtypes.UUID                                           `tfsdk:"item_id"`
	Values      supertypes.SetNestedObjectValueOf[externalDataSharesModel] `tfsdk:"values"`
	Timeouts    timeoutsD.Value                                            `tfsdk:"timeouts"`
}

/*
RESOURCE
*/

type resourceExternalDataSharesProviderModel struct {
	externalDataSharesModel

	Timeouts timeoutsR.Value `tfsdk:"timeouts"`
}

type requestCreateExternalDataShare struct {
	fabcore.CreateExternalDataShareRequest
}

func (to *externalDataSharesModel) set(ctx context.Context, from *fabcore.ExternalDataShare) diag.Diagnostics {
	to.ID = customtypes.NewUUIDPointerValue(from.ID)
	to.Status = types.StringPointerValue((*string)(from.Status))
	to.WorkspaceID = customtypes.NewUUIDPointerValue(from.WorkspaceID)
	to.ItemID = customtypes.NewUUIDPointerValue(from.ItemID)
	to.InvitationURL = types.StringPointerValue(from.InvitationURL)
	to.AcceptedByTenantID = customtypes.NewUUIDPointerValue(from.AcceptedByTenantID)

	if from.ExpirationTimeUTC != nil {
		to.ExpirationTimeUtc = types.StringValue(from.ExpirationTimeUTC.Format(time.RFC3339))
	}

	to.Paths.SetNull(ctx)

	if from.Paths != nil {
		values := make([]types.String, 0, len(from.Paths))
		for _, value := range from.Paths {
			values = append(values, types.StringValue(value))
		}

		if diags := to.Paths.Set(ctx, values); diags.HasError() {
			return diags
		}
	}

	to.CreatorPrincipal = supertypes.NewSingleNestedObjectValueOfNull[creatorPrincipalModel](ctx)
	to.Recipient = supertypes.NewSingleNestedObjectValueOfNull[recipientModel](ctx)

	if from.Recipient != nil {
		recipient := &recipientModel{}
		recipient.set(*from.Recipient)
		if diags := to.Recipient.Set(ctx, recipient); diags.HasError() {
			return diags
		}
	}

	if from.CreatorPrincipal != nil {
		creatorPrincipal := &creatorPrincipalModel{}
		creatorPrincipal.set(ctx, *from.CreatorPrincipal)
		if diags := to.CreatorPrincipal.Set(ctx, creatorPrincipal); diags.HasError() {
			return diags
		}
	}

	return nil
}

func (to *dataSourceExternalDataSharesProviderModel) set(ctx context.Context, from []fabcore.ExternalDataShare) diag.Diagnostics {
	to.WorkspaceID = customtypes.NewUUIDPointerValue(from[0].WorkspaceID)
	to.ItemID = customtypes.NewUUIDPointerValue(from[0].ItemID)

	slice := make([]*externalDataSharesModel, 0, len(from))
	for _, item := range from {
		externalDataShare := &externalDataSharesModel{}
		if diags := externalDataShare.set(ctx, &item); diags.HasError() {
			return diags
		}

		slice = append(slice, externalDataShare)
	}

	if diags := to.Values.Set(ctx, slice); diags.HasError() {
		return diags
	}

	return nil
}

func (to *creatorPrincipalModel) set(ctx context.Context, from fabcore.Principal) diag.Diagnostics {
	to.ID = customtypes.NewUUIDPointerValue(from.ID)
	to.DisplayName = types.StringPointerValue(from.DisplayName)
	to.Type = types.StringPointerValue((*string)(from.Type))
	to.UserDetails = supertypes.NewSingleNestedObjectValueOfNull[userDetailsModel](ctx)

	if from.UserDetails != nil {
		userDetails := &userDetailsModel{}
		userDetails.set(*from.UserDetails)
		if diags := to.UserDetails.Set(ctx, userDetails); diags.HasError() {
			return diags
		}
	}

	return nil
}

func (to *userDetailsModel) set(from fabcore.PrincipalUserDetails) {
	to.UserPrincipalName = types.StringPointerValue(from.UserPrincipalName)
}

func (to *recipientModel) set(from fabcore.ExternalDataShareRecipient) {
	to.UserPrincipalName = types.StringPointerValue(from.UserPrincipalName)
	to.TenantID = customtypes.NewUUIDPointerValue(from.TenantID)
}

func (to *requestCreateExternalDataShare) set(ctx context.Context, from resourceExternalDataSharesProviderModel) {
	paths, _ := from.Paths.Get(ctx)

	to.Paths = make([]string, 0, len(paths))
	for _, path := range paths {
		to.Paths = append(to.Paths, path.ValueString())
	}

	to.Recipient = &fabcore.ExternalDataShareRecipient{}

	recipient, _ := from.Recipient.Get(ctx)

	to.Recipient.UserPrincipalName = recipient.UserPrincipalName.ValueStringPointer()
	to.Recipient.TenantID = recipient.TenantID.ValueStringPointer()
}
