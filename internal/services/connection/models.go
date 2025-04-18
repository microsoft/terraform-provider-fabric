// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package connection

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
)

type baseConnectionModel[ConnectionDetails dsConnectionDetailsModel | rsConnectionDetailsModel, CredentialDetails dsCredentialDetailsModel | rsCredentialDetailsModel] struct {
	ID                customtypes.UUID                                        `tfsdk:"id"`
	DisplayName       types.String                                            `tfsdk:"display_name"`
	GatewayID         customtypes.UUID                                        `tfsdk:"gateway_id"`
	ConnectivityType  types.String                                            `tfsdk:"connectivity_type"`
	PrivacyLevel      types.String                                            `tfsdk:"privacy_level"`
	ConnectionDetails supertypes.SingleNestedObjectValueOf[ConnectionDetails] `tfsdk:"connection_details"`
	CredentialDetails supertypes.SingleNestedObjectValueOf[CredentialDetails] `tfsdk:"credential_details"`
}

func (to *baseConnectionModel[ConnectionDetails, CredentialDetails]) set(ctx context.Context, from fabcore.Connection) diag.Diagnostics {
	to.ID = customtypes.NewUUIDPointerValue(from.ID)
	to.DisplayName = types.StringPointerValue(from.DisplayName)
	to.GatewayID = customtypes.NewUUIDPointerValue(from.GatewayID)
	to.ConnectivityType = types.StringPointerValue((*string)(from.ConnectivityType))
	to.PrivacyLevel = types.StringPointerValue((*string)(from.PrivacyLevel))

	var diags diag.Diagnostics

	var connectionDetails supertypes.SingleNestedObjectValueOf[ConnectionDetails]
	var connectionDetailsModel *ConnectionDetails

	if to.ConnectionDetails.IsKnown() {
		connectionDetails = to.ConnectionDetails

		connectionDetailsModel, diags = to.ConnectionDetails.Get(ctx)
		if diags.HasError() {
			return diags
		}
	} else {
		connectionDetails = supertypes.NewSingleNestedObjectValueOfNull[ConnectionDetails](ctx)
	}

	if from.ConnectionDetails != nil {
		var connectionDetailsModelPtr *ConnectionDetails

		switch v := any(connectionDetailsModel).(type) {
		case *dsConnectionDetailsModel:
			setDSConnectionDetails(*from.ConnectionDetails, v)
			if convertedValue, ok := any(*v).(ConnectionDetails); ok {
				connectionDetailsModelPtr = &convertedValue
			}
		case *rsConnectionDetailsModel:
			setRSConnectionDetails(*from.ConnectionDetails, v)
			if convertedValue, ok := any(*v).(ConnectionDetails); ok {
				connectionDetailsModelPtr = &convertedValue
			}
		}

		if diags := connectionDetails.Set(ctx, connectionDetailsModelPtr); diags.HasError() {
			return diags
		}
	}

	to.ConnectionDetails = connectionDetails

	var credentialDetails supertypes.SingleNestedObjectValueOf[CredentialDetails]
	var credentialDetailsModel *CredentialDetails

	if to.CredentialDetails.IsKnown() {
		credentialDetails = to.CredentialDetails

		credentialDetailsModel, diags = to.CredentialDetails.Get(ctx)
		if diags.HasError() {
			return diags
		}

	} else {
		credentialDetails = supertypes.NewSingleNestedObjectValueOfNull[CredentialDetails](ctx)
	}

	if from.CredentialDetails != nil {
		var credentialDetailsModelPtr *CredentialDetails
		switch v := any(credentialDetailsModel).(type) {
		case *dsCredentialDetailsModel:
			setDSCredentialDetails(*from.CredentialDetails, v)
			if convertedValue, ok := any(*v).(CredentialDetails); ok {
				credentialDetailsModelPtr = &convertedValue
			}
		case *rsCredentialDetailsModel:
			setRSCredentialDetails(*from.CredentialDetails, v)
			if convertedValue, ok := any(*v).(CredentialDetails); ok {
				credentialDetailsModelPtr = &convertedValue
			}
		}

		if diags := credentialDetails.Set(ctx, credentialDetailsModelPtr); diags.HasError() {
			return diags
		}
	}

	to.CredentialDetails = credentialDetails

	return nil
}

type credentialsBasicModel struct {
	Username          types.String `tfsdk:"username"`
	Password          types.String `tfsdk:"password"`
	PasswordWO        types.String `tfsdk:"password_wo"`
	PasswordWOVersion types.Int32  `tfsdk:"password_wo_version"`
}

type credentialsKeyModel struct {
	Key          types.String `tfsdk:"key"`
	KeyWO        types.String `tfsdk:"key_wo"`
	KeyWOVersion types.Int32  `tfsdk:"key_wo_version"`
}

type credentialsServicePrincipalModel struct {
	TenantID              types.String `tfsdk:"tenant_id"`
	ClientID              types.String `tfsdk:"client_id"`
	ClientSecret          types.String `tfsdk:"client_secret"`
	ClientSecretWO        types.String `tfsdk:"client_secret_wo"`
	ClientSecretWOVersion types.Int32  `tfsdk:"client_secret_wo_version"`
}

type credentialsSharedAccessSignatureModel struct {
	Token          types.String `tfsdk:"token"`
	TokenWO        types.String `tfsdk:"token_wo"`
	TokenWOVersion types.Int32  `tfsdk:"token_wo_version"`
}

type credentialsWindowsModel struct {
	Username          types.String `tfsdk:"username"`
	Password          types.String `tfsdk:"password"`
	PasswordWO        types.String `tfsdk:"password_wo"`
	PasswordWOVersion types.Int32  `tfsdk:"password_wo_version"`
}
