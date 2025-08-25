package externaldatasharesprovider

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	fabadmin "github.com/microsoft/fabric-sdk-go/fabric/admin"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
)

type externalDataSharesModel struct {
	ID                 customtypes.UUID                                            `tfsdk:"id"`
	Paths              supertypes.SetValueOf[types.String]                         `tfsdk:"paths"`
	Status             types.String                                                `tfsdk:"status"`
	Recipient          supertypes.SingleNestedObjectValueOf[recipientModel]        `tfsdk:"recipient"`
	ExpirationTimeUtc  types.String                                                `tfsdk:"expiration_time_utc"`
	CreatorPrincipal   supertypes.SingleNestedObjectValueOf[creatorPrincipalModel] `tfsdk:"creator_principal"`
	WorkspaceID        customtypes.UUID                                            `tfsdk:"workspace_id"`
	ItemID             customtypes.UUID                                            `tfsdk:"item_id"`
	InvitationUrl      types.String                                                `tfsdk:"invitation_url"`
	AcceptedByTenantID customtypes.UUID                                            `tfsdk:"accepted_by_tenant_id"`
}

type creatorPrincipalModel struct {
	ID          customtypes.UUID                                       `tfsdk:"id"`
	DisplayName types.String                                           `tfsdk:"display_name"`
	Type        types.String                                           `tfsdk:"type"`
	UserDetails supertypes.SingleNestedObjectValueOf[userDetailsModel] `tfsdk:"user_details"`
}

type recipientModel struct {
	UserPrincipalName types.String `tfsdk:"user_principal_name"`
}

type userDetailsModel struct {
	UserPrincipalName types.String `tfsdk:"user_principal_name"`
}

type baseExternalDataSharesProviderModel struct {
	Value supertypes.SetNestedObjectValueOf[externalDataSharesModel] `tfsdk:"value"`
}

func (to *baseExternalDataSharesProviderModel) set(ctx context.Context, from []fabadmin.ExternalDataShare) diag.Diagnostics {
	slice := make([]*externalDataSharesModel, 0, len(from))
	for _, item := range from {
		externalDataShare := &externalDataSharesModel{
			ID:                 customtypes.NewUUIDPointerValue(item.ID),
			Status:             types.StringPointerValue((*string)(item.Status)),
			WorkspaceID:        customtypes.NewUUIDPointerValue(item.WorkspaceID),
			ItemID:             customtypes.NewUUIDPointerValue(item.ItemID),
			InvitationUrl:      types.StringPointerValue(item.InvitationURL),
			AcceptedByTenantID: customtypes.NewUUIDPointerValue(item.AcceptedByTenantID),
		}

		if item.ExpirationTimeUTC != nil {
			externalDataShare.ExpirationTimeUtc = types.StringValue(item.ExpirationTimeUTC.Format(time.RFC3339))
		}
		externalDataShare.Paths = supertypes.NewSetValueOfNull[types.String](ctx)
		externalDataShare.CreatorPrincipal = supertypes.NewSingleNestedObjectValueOfNull[creatorPrincipalModel](ctx)
		externalDataShare.Recipient = supertypes.NewSingleNestedObjectValueOfNull[recipientModel](ctx)

		if item.Paths != nil {
			elems := make([]types.String, 0, len(item.Paths))
			for _, p := range item.Paths {
				elems = append(elems, types.StringValue(p))
			}
			externalDataShare.Paths.Set(ctx, elems)
		}

		if item.Recipient != nil {
			recipient := &recipientModel{}
			recipient.set(*item.Recipient)
			externalDataShare.Recipient.Set(ctx, recipient)
		}

		if item.CreatorPrincipal != nil {
			creatorPrincipal := &creatorPrincipalModel{}
			creatorPrincipal.set(ctx, *item.CreatorPrincipal)
			externalDataShare.CreatorPrincipal.Set(ctx, creatorPrincipal)
		}

		slice = append(slice, externalDataShare)
	}

	if diags := to.Value.Set(ctx, slice); diags.HasError() {
		return diags
	}

	return nil
}

func (to *creatorPrincipalModel) set(ctx context.Context, from fabadmin.Principal) diag.Diagnostics {
	to.ID = customtypes.NewUUIDPointerValue(from.ID)
	to.DisplayName = types.StringPointerValue(from.DisplayName)
	to.Type = types.StringPointerValue((*string)(from.Type))
	to.UserDetails = supertypes.NewSingleNestedObjectValueOfNull[userDetailsModel](ctx)

	if from.UserDetails != nil {
		userDetails := &userDetailsModel{}
		userDetails.set(*from.UserDetails)
		to.UserDetails.Set(ctx, userDetails)
	}

	return nil
}

func (to *userDetailsModel) set(from fabadmin.PrincipalUserDetails) {
	to.UserPrincipalName = types.StringPointerValue(from.UserPrincipalName)
}

func (to *recipientModel) set(from fabadmin.ExternalDataShareRecipient) {
	to.UserPrincipalName = types.StringPointerValue(from.UserPrincipalName)
}

type resourceExternalDataSharesProviderModel struct {
	ExternalDataShareID customtypes.UUID `tfsdk:"external_data_share_id"`
	ItemID              customtypes.UUID `tfsdk:"item_id"`
	WorkspaceID         customtypes.UUID `tfsdk:"workspace_id"`
}
