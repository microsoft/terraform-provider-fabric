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
	PasswordWO        types.String `tfsdk:"password_wo"`
	PasswordWOVersion types.Int32  `tfsdk:"password_wo_version"`
}

type credentialsKeyModel struct {
	KeyWO        types.String `tfsdk:"key_wo"`
	KeyWOVersion types.Int32  `tfsdk:"key_wo_version"`
}

type credentialsServicePrincipalModel struct {
	TenantID              types.String `tfsdk:"tenant_id"`
	ClientID              types.String `tfsdk:"client_id"`
	ClientSecretWO        types.String `tfsdk:"client_secret_wo"`
	ClientSecretWOVersion types.Int32  `tfsdk:"client_secret_wo_version"`
}

type credentialsSharedAccessSignatureModel struct {
	TokenWO        types.String `tfsdk:"token_wo"`
	TokenWOVersion types.Int32  `tfsdk:"token_wo_version"`
}

type credentialsWindowsModel struct {
	Username          types.String `tfsdk:"username"`
	PasswordWO        types.String `tfsdk:"password_wo"`
	PasswordWOVersion types.Int32  `tfsdk:"password_wo_version"`
}

type credentialsOAuth2Model struct {
	AccessTokenWO        types.String `tfsdk:"access_token_wo"`
	AccessTokenWOVersion types.Int32  `tfsdk:"access_token_wo_version"`
}

// type credentialsEncryptedModel struct {
// 	ValueWO        types.String `tfsdk:"value_wo"`
// 	ValueWOVersion types.Int32  `tfsdk:"value_wo_version"`
// }

// type credentialsOnPremisesGatewayModel struct {
// 	GatewayID customtypes.UUID `tfsdk:"gateway_id"`
// 	// EncryptedCredentialsWO        types.String     `tfsdk:"encrypted_credentials_wo"`
// 	// EncryptedCredentialsWOVersion types.Int32      `tfsdk:"encrypted_credentials_wo_version"`
// 	EncryptedCredentials supertypes.SingleNestedObjectValueOf[credentialsEncryptedModel] `tfsdk:"encrypted_credentials"`

// 	CredentialType     types.String                                                  `tfsdk:"credential_type"`
// 	BasicCredentials   supertypes.SingleNestedObjectValueOf[credentialsBasicModel]   `tfsdk:"basic_credentials"`
// 	WindowsCredentials supertypes.SingleNestedObjectValueOf[credentialsWindowsModel] `tfsdk:"windows_credentials"`
// 	KeyCredentials     supertypes.SingleNestedObjectValueOf[credentialsKeyModel]     `tfsdk:"key_credentials"`
// 	OAuth2Credentials  supertypes.SingleNestedObjectValueOf[credentialsOAuth2Model]  `tfsdk:"oauth2_credentials"`
// 	PublicKey          supertypes.SingleNestedObjectValueOf[publicKeyModel]          `tfsdk:"public_key"`
// }

// type publicKeyModel struct {
// 	Exponent types.String `tfsdk:"exponent"`
// 	Modulus  types.String `tfsdk:"modulus"`
// }

// func (to *publicKeyModel) set(from fabcore.PublicKey) {
// 	to.Exponent = types.StringPointerValue(from.Exponent)
// 	to.Modulus = types.StringPointerValue(from.Modulus)
// }

// type publicKey struct {
// 	Exponent string
// 	Modulus  string
// }
