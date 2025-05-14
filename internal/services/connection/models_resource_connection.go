// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package connection

import (
	"context"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/credentials"
)

type resourceConnectionModel[ConnectionDetails dsConnectionDetailsModel | rsConnectionDetailsModel, CredentialDetails dsCredentialDetailsModel | rsCredentialDetailsModel] struct {
	baseConnectionModel[ConnectionDetails, CredentialDetails]
	Timeouts timeouts.Value `tfsdk:"timeouts"`
}

func (m resourceConnectionModel[ConnectionDetails, CredentialDetails]) getConnectionDetails(ctx context.Context) (*ConnectionDetails, diag.Diagnostics) {
	if !m.ConnectionDetails.IsNull() && !m.ConnectionDetails.IsUnknown() {
		return m.ConnectionDetails.Get(ctx)
	}

	return nil, nil
}

func (m resourceConnectionModel[ConnectionDetails, CredentialDetails]) getCredentialDetails(ctx context.Context) (*CredentialDetails, diag.Diagnostics) {
	if !m.CredentialDetails.IsNull() && !m.CredentialDetails.IsUnknown() {
		return m.CredentialDetails.Get(ctx)
	}

	return nil, nil
}

type rsConnectionDetailsModel struct {
	Path           types.String                                                 `tfsdk:"path"` // computed
	Type           types.String                                                 `tfsdk:"type"`
	CreationMethod types.String                                                 `tfsdk:"creation_method"`
	Parameters     supertypes.SetNestedObjectValueOf[connectionParametersModel] `tfsdk:"parameters"`
}

type connectionParametersModel struct {
	Name     types.String `tfsdk:"name"`
	Value    types.String `tfsdk:"value"`
	DataType types.String `tfsdk:"data_type"`
}

func (m rsConnectionDetailsModel) getParameters(ctx context.Context) (map[string]string, diag.Diagnostics) {
	if !m.Parameters.IsNull() && !m.Parameters.IsUnknown() {
		parametersModel, diags := m.Parameters.Get(ctx)
		if diags.HasError() {
			return nil, diags
		}

		parameters := make(map[string]string, len(parametersModel))
		for _, parameterModel := range parametersModel {
			parameters[parameterModel.Name.ValueString()] = parameterModel.Value.ValueString()
		}

		return parameters, nil
	}

	return nil, nil
}

func setRSConnectionDetails(from fabcore.ListConnectionDetails, to *rsConnectionDetailsModel) {
	to.Path = types.StringPointerValue(from.Path)
	to.Type = types.StringPointerValue(from.Type)
}

type rsCredentialDetailsModel struct {
	ConnectionEncryption types.String `tfsdk:"connection_encryption"`
	SingleSignOnType     types.String `tfsdk:"single_sign_on_type"`
	SkipTestConnection   types.Bool   `tfsdk:"skip_test_connection"`
	CredentialType       types.String `tfsdk:"credential_type"`

	// AnonymousCredentials                   supertypes.SingleNestedObjectValueOf[anonymousCredentialsModel]                   `tfsdk:"anonymous_credentials"`
	BasicCredentials                 supertypes.SingleNestedObjectValueOf[credentialsBasicModel]                 `tfsdk:"basic_credentials"`
	KeyCredentials                   supertypes.SingleNestedObjectValueOf[credentialsKeyModel]                   `tfsdk:"key_credentials"`
	ServicePrincipalCredentials      supertypes.SingleNestedObjectValueOf[credentialsServicePrincipalModel]      `tfsdk:"service_principal_credentials"`
	SharedAccessSignatureCredentials supertypes.SingleNestedObjectValueOf[credentialsSharedAccessSignatureModel] `tfsdk:"shared_access_signature_credentials"`
	WindowsCredentials               supertypes.SingleNestedObjectValueOf[credentialsWindowsModel]               `tfsdk:"windows_credentials"`
	OAuth2Credentials                supertypes.SingleNestedObjectValueOf[credentialsOAuth2Model]                `tfsdk:"oauth2_credentials"`
	// EncryptedCredentials             supertypes.SingleNestedObjectValueOf[credentialsEncryptedModel]             `tfsdk:"encrypted_credentials"`
	// WindowsWithoutImpersonationCredentials supertypes.SingleNestedObjectValueOf[credentialsWindowsWithoutImpersonationModel] `tfsdk:"windows_without_impersonation_credentials"`
	// WorkspaceIdentityCredentials           supertypes.SingleNestedObjectValueOf[credentialsWorkspaceIdentityModel]           `tfsdk:"workspace_identity_credentials"`
	// OnPremisesGatewayCredentials supertypes.ListNestedObjectValueOf[credentialsOnPremisesGatewayModel] `tfsdk:"on_premises_gateway_credentials"`
}

// func (to *rsCredentialDetailsModel) set(from fabcore.ListCredentialDetails) {
// 	to.ConnectionEncryption = types.StringPointerValue((*string)(from.ConnectionEncryption))
// 	to.SingleSignOnType = types.StringPointerValue((*string)(from.SingleSignOnType))
// 	to.SkipTestConnection = types.BoolPointerValue(from.SkipTestConnection)
// 	to.CredentialType = types.StringPointerValue((*string)(from.CredentialType))
// }

func setRSCredentialDetails(from fabcore.ListCredentialDetails, to *rsCredentialDetailsModel) {
	to.ConnectionEncryption = types.StringPointerValue((*string)(from.ConnectionEncryption))
	to.SingleSignOnType = types.StringPointerValue((*string)(from.SingleSignOnType))
	to.SkipTestConnection = types.BoolPointerValue(from.SkipTestConnection)
	to.CredentialType = types.StringPointerValue((*string)(from.CredentialType))
}

type requestCreateConnection struct {
	fabcore.CreateConnectionRequestClassification
}

func (to *requestCreateConnection) set(
	ctx context.Context,
	plan, config resourceConnectionModel[rsConnectionDetailsModel, rsCredentialDetailsModel],
	gwMembers []fabcore.OnPremisesGatewayMember,
) diag.Diagnostics {
	connectivityType := (fabcore.ConnectivityType)(plan.ConnectivityType.ValueString())

	var requestCreateConnectionDetails requestCreateConnectionDetails
	if diags := requestCreateConnectionDetails.set(ctx, plan.ConnectionDetails); diags.HasError() {
		return diags
	}

	var requestCreateCredentialDetails requestCreateCredentialDetails
	if connectivityType == fabcore.ConnectivityTypeShareableCloud ||
		connectivityType == fabcore.ConnectivityTypeVirtualNetworkGateway || connectivityType == fabcore.ConnectivityTypeOnPremisesGateway {
		if diags := requestCreateCredentialDetails.set(ctx, config.CredentialDetails); diags.HasError() {
			return diags
		}
	}

	displayName := plan.DisplayName.ValueStringPointer()
	privacyLevel := (*fabcore.PrivacyLevel)(plan.PrivacyLevel.ValueStringPointer())

	switch connectivityType {
	case fabcore.ConnectivityTypeShareableCloud: // fabcore.ConnectivityTypePersonalCloud:
		to.CreateConnectionRequestClassification = &fabcore.CreateCloudConnectionRequest{
			DisplayName:       displayName,
			PrivacyLevel:      privacyLevel,
			ConnectivityType:  &connectivityType,
			ConnectionDetails: &requestCreateConnectionDetails.CreateConnectionDetails,
			CredentialDetails: &requestCreateCredentialDetails.CreateCredentialDetails,
		}

	case fabcore.ConnectivityTypeVirtualNetworkGateway:
		to.CreateConnectionRequestClassification = &fabcore.CreateVirtualNetworkGatewayConnectionRequest{
			DisplayName:       displayName,
			PrivacyLevel:      privacyLevel,
			ConnectivityType:  &connectivityType,
			ConnectionDetails: &requestCreateConnectionDetails.CreateConnectionDetails,
			CredentialDetails: &requestCreateCredentialDetails.CreateCredentialDetails,
			GatewayID:         plan.GatewayID.ValueStringPointer(),
		}

	case fabcore.ConnectivityTypeOnPremisesGateway: // fabcore.ConnectivityTypeOnPremisesGatewayPersonal:
		var credentialDetails requestCreateOnPremisesCredentialDetails

		if diags := credentialDetails.set(ctx, config.CredentialDetails, requestCreateCredentialDetails.Credentials, gwMembers); diags.HasError() {
			return diags
		}

		to.CreateConnectionRequestClassification = &fabcore.CreateOnPremisesConnectionRequest{
			DisplayName:       displayName,
			PrivacyLevel:      privacyLevel,
			ConnectivityType:  &connectivityType,
			ConnectionDetails: &requestCreateConnectionDetails.CreateConnectionDetails,
			CredentialDetails: &credentialDetails.CreateOnPremisesCredentialDetails,
			GatewayID:         plan.GatewayID.ValueStringPointer(),
		}

	case fabcore.ConnectivityTypeAutomatic: // fabcore.ConnectivityTypeNone:
		to.CreateConnectionRequestClassification = &fabcore.CreateConnectionRequest{
			DisplayName:       displayName,
			PrivacyLevel:      privacyLevel,
			ConnectivityType:  &connectivityType,
			ConnectionDetails: &requestCreateConnectionDetails.CreateConnectionDetails,
		}
	}

	return nil
}

type requestCreateConnectionDetails struct {
	fabcore.CreateConnectionDetails
}

func (to *requestCreateConnectionDetails) set(ctx context.Context, from supertypes.SingleNestedObjectValueOf[rsConnectionDetailsModel]) diag.Diagnostics {
	var diags diag.Diagnostics

	connectionDetails, diags := from.Get(ctx)
	if diags.HasError() {
		return diags
	}

	var params []fabcore.ConnectionDetailsParameterClassification

	if !connectionDetails.Parameters.IsNull() && !connectionDetails.Parameters.IsUnknown() {
		parameters, diags := connectionDetails.Parameters.Get(ctx)
		if diags.HasError() {
			return diags
		}

		for _, parameter := range parameters {
			var requestParameter fabcore.ConnectionDetailsParameterClassification

			dataType := (fabcore.DataType)(parameter.DataType.ValueString())
			name := parameter.Name.ValueString()
			value := parameter.Value.ValueString()

			switch dataType {
			case fabcore.DataTypeBoolean:
				boolValue, err := strconv.ParseBool(parameter.Value.ValueString())
				if err != nil {
					diags.AddError(
						"Boolean parameter",
						err.Error(),
					)
				}

				requestParameter = &fabcore.ConnectionDetailsBooleanParameter{
					DataType: &dataType,
					Name:     &name,
					Value:    &boolValue,
				}

			case fabcore.DataTypeDate:
				dateValue, err := time.Parse("2006-01-02", value)
				if err != nil {
					diags.AddError(
						"Date parameter",
						err.Error(),
					)
					return diags
				}

				requestParameter = &fabcore.ConnectionDetailsDateParameter{
					DataType: &dataType,
					Name:     &name,
					Value:    &dateValue,
				}

			case fabcore.DataTypeDateTime:
				dateTimeValue, err := time.Parse("2006-01-02T15:04:05.000Z07:00", value)
				if err != nil {
					diags.AddError(
						"DateTime parameter",
						err.Error(),
					)
					return diags
				}

				requestParameter = &fabcore.ConnectionDetailsDateTimeParameter{
					DataType: &dataType,
					Name:     &name,
					Value:    &dateTimeValue,
				}

			case fabcore.DataTypeDateTimeZone:
				requestParameter = &fabcore.ConnectionDetailsDateTimeZoneParameter{
					DataType: &dataType,
					Name:     &name,
					Value:    &value,
				}

			case fabcore.DataTypeDuration:
				requestParameter = &fabcore.ConnectionDetailsDurationParameter{
					DataType: &dataType,
					Name:     &name,
					Value:    &value,
				}

			case fabcore.DataTypeNumber:
				float64Value, err := strconv.ParseFloat(value, 32)
				if err != nil {
					diags.AddError(
						"Number parameter",
						err.Error(),
					)
					return diags
				}

				float32Value := float32(float64Value)

				requestParameter = &fabcore.ConnectionDetailsNumberParameter{
					DataType: &dataType,
					Name:     &name,
					Value:    &float32Value,
				}

			case fabcore.DataTypeText:
				requestParameter = &fabcore.ConnectionDetailsTextParameter{
					DataType: &dataType,
					Name:     &name,
					Value:    &value,
				}

			case fabcore.DataTypeTime:
				timeValue, err := time.Parse("15:04:05.000Z07:00", value)
				if err != nil {
					diags.AddError(
						"DateTime parameter",
						err.Error(),
					)

					return diags
				}

				requestParameter = &fabcore.ConnectionDetailsTimeParameter{
					DataType: &dataType,
					Name:     &name,
					Value:    &timeValue,
				}
			}

			params = append(params, requestParameter)
		}
	}

	to.Parameters = params
	to.CreationMethod = connectionDetails.CreationMethod.ValueStringPointer()
	to.Type = connectionDetails.Type.ValueStringPointer()

	return nil
}

type requestCreateCredentialDetails struct {
	fabcore.CreateCredentialDetails
}

func (to *requestCreateCredentialDetails) set(ctx context.Context, from supertypes.SingleNestedObjectValueOf[rsCredentialDetailsModel]) diag.Diagnostics {
	credentialDetails, diags := from.Get(ctx)
	if diags.HasError() {
		return diags
	}

	to.ConnectionEncryption = (*fabcore.ConnectionEncryption)(credentialDetails.ConnectionEncryption.ValueStringPointer())
	to.SingleSignOnType = (*fabcore.SingleSignOnType)(credentialDetails.SingleSignOnType.ValueStringPointer())
	to.SkipTestConnection = credentialDetails.SkipTestConnection.ValueBoolPointer()

	credentialType := (fabcore.CredentialType)(credentialDetails.CredentialType.ValueString())

	var requestCreateCredential fabcore.CredentialsClassification

	switch credentialType {
	case fabcore.CredentialTypeAnonymous:
		requestCreateCredential = &fabcore.AnonymousCredentials{
			CredentialType: &credentialType,
		}

	case fabcore.CredentialTypeBasic:
		cred, diags := credentialDetails.BasicCredentials.Get(ctx)
		if diags.HasError() {
			return diags
		}

		requestCreateCredential = &fabcore.BasicCredentials{
			CredentialType: &credentialType,
			Username:       cred.Username.ValueStringPointer(),
			Password:       cred.PasswordWO.ValueStringPointer(),
		}

	case fabcore.CredentialTypeKey:
		cred, diags := credentialDetails.KeyCredentials.Get(ctx)
		if diags.HasError() {
			return diags
		}

		requestCreateCredential = &fabcore.KeyCredentials{
			CredentialType: &credentialType,
			Key:            cred.KeyWO.ValueStringPointer(),
		}

	case fabcore.CredentialTypeOAuth2:
		requestCreateCredential = &fabcore.Credentials{
			CredentialType: &credentialType,
		}

	case fabcore.CredentialTypeServicePrincipal:
		cred, diags := credentialDetails.ServicePrincipalCredentials.Get(ctx)
		if diags.HasError() {
			return diags
		}

		requestCreateCredential = &fabcore.ServicePrincipalCredentials{
			CredentialType:           &credentialType,
			TenantID:                 cred.TenantID.ValueStringPointer(),
			ServicePrincipalClientID: cred.ClientID.ValueStringPointer(),
			ServicePrincipalSecret:   cred.ClientSecretWO.ValueStringPointer(),
		}

	case fabcore.CredentialTypeSharedAccessSignature:
		cred, diags := credentialDetails.SharedAccessSignatureCredentials.Get(ctx)
		if diags.HasError() {
			return diags
		}

		requestCreateCredential = &fabcore.SharedAccessSignatureCredentials{
			CredentialType: &credentialType,
			Token:          cred.TokenWO.ValueStringPointer(),
		}

	case fabcore.CredentialTypeWindows:
		cred, diags := credentialDetails.WindowsCredentials.Get(ctx)
		if diags.HasError() {
			return diags
		}

		requestCreateCredential = &fabcore.WindowsCredentials{
			CredentialType: &credentialType,
			Username:       cred.Username.ValueStringPointer(),
			Password:       cred.PasswordWO.ValueStringPointer(),
		}

	case fabcore.CredentialTypeWindowsWithoutImpersonation:
		requestCreateCredential = &fabcore.WindowsWithoutImpersonationCredentials{
			CredentialType: &credentialType,
		}
	case fabcore.CredentialTypeWorkspaceIdentity:
		requestCreateCredential = &fabcore.WorkspaceIdentityCredentials{
			CredentialType: &credentialType,
		}
	}

	to.Credentials = requestCreateCredential

	return nil
}

type requestCreateOnPremisesCredentialDetails struct {
	fabcore.CreateOnPremisesCredentialDetails
}

func (to *requestCreateOnPremisesCredentialDetails) set(
	ctx context.Context,
	from supertypes.SingleNestedObjectValueOf[rsCredentialDetailsModel],
	creds fabcore.CredentialsClassification,
	gwMembers []fabcore.OnPremisesGatewayMember,
) diag.Diagnostics {
	var err error

	credentialDetails, diags := from.Get(ctx)
	if diags.HasError() {
		return diags
	}

	to.ConnectionEncryption = (*fabcore.ConnectionEncryption)(credentialDetails.ConnectionEncryption.ValueStringPointer())
	to.SingleSignOnType = (*fabcore.SingleSignOnType)(credentialDetails.SingleSignOnType.ValueStringPointer())
	to.SkipTestConnection = credentialDetails.SkipTestConnection.ValueBoolPointer()

	c, err := credentials.NewCredentials(creds)
	if err != nil {
		diags.AddError("summary string", err.Error())

		return diags
	}

	var values []fabcore.OnPremisesCredentialEntry

	for _, gwMember := range gwMembers {
		encryptedCredentials, err := credentials.EncryptCredentials(*c, *gwMember.PublicKey)
		if err != nil {
			diags.AddError("summary string", err.Error())
			return diags
		}

		values = append(values, fabcore.OnPremisesCredentialEntry{
			GatewayID:            gwMember.ID,
			EncryptedCredentials: &encryptedCredentials,
		})
	}

	to.Credentials = &fabcore.OnPremisesGatewayCredentials{
		CredentialType: &c.Type,
		Values:         values,
	}

	return nil
}

type (
	gatewayCredential struct {
		CredentialData []gatewayCredentialEntity `json:"credentialData"`
	}
	gatewayCredentialEntity struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	}
)

// func newOnPremisesCredentialEntry(ctx context.Context, entity credentialsOnPremisesGatewayModel) (*gatewayCredential, diag.Diagnostics) {
// 	var diags diag.Diagnostics
// 	var gwCredential *gatewayCredential

// 	switch fabcore.CredentialType(entity.CredentialType.ValueString()) {
// 	case fabcore.CredentialTypeWindows:
// 		credentials, diags := entity.WindowsCredentials.Get(ctx)
// 		if diags.HasError() {
// 			return nil, diags
// 		}

// 		username := gatewayCredentialEntity{
// 			Name:  "username",
// 			Value: credentials.Username.ValueString(),
// 		}

// 		password := gatewayCredentialEntity{
// 			Name:  "password",
// 			Value: credentials.PasswordWO.ValueString(),
// 		}

// 		gwCredential = &gatewayCredential{
// 			CredentialData: []gatewayCredentialEntity{username, password},
// 		}

// 	case fabcore.CredentialTypeBasic:
// 		credentials, diags := entity.BasicCredentials.Get(ctx)
// 		if diags.HasError() {
// 			return nil, diags
// 		}

// 		username := gatewayCredentialEntity{
// 			Name:  "username",
// 			Value: credentials.Username.ValueString(),
// 		}

// 		password := gatewayCredentialEntity{
// 			Name:  "password",
// 			Value: credentials.PasswordWO.ValueString(),
// 		}

// 		gwCredential = &gatewayCredential{
// 			CredentialData: []gatewayCredentialEntity{username, password},
// 		}

// 	case fabcore.CredentialTypeKey:
// 		credentials, diags := entity.KeyCredentials.Get(ctx)
// 		if diags.HasError() {
// 			return nil, diags
// 		}

// 		key := gatewayCredentialEntity{
// 			Name:  "key",
// 			Value: credentials.KeyWO.ValueString(),
// 		}

// 		gwCredential = &gatewayCredential{
// 			CredentialData: []gatewayCredentialEntity{key},
// 		}

// 	case fabcore.CredentialTypeOAuth2:
// 		credentials, diags := entity.OAuth2Credentials.Get(ctx)
// 		if diags.HasError() {
// 			return nil, diags
// 		}

// 		accessToken := gatewayCredentialEntity{
// 			Name:  "accessToken",
// 			Value: credentials.AccessTokenWO.ValueString(),
// 		}

// 		gwCredential = &gatewayCredential{
// 			CredentialData: []gatewayCredentialEntity{accessToken},
// 		}
// 	}

// 	return gwCredential, diags
// }

type requestUpdateCredentialDetails struct {
	fabcore.UpdateCredentialDetails
}

func (to *requestUpdateCredentialDetails) set(ctx context.Context, from supertypes.SingleNestedObjectValueOf[rsCredentialDetailsModel]) diag.Diagnostics {
	credentialDetails, diags := from.Get(ctx)
	if diags.HasError() {
		return diags
	}

	to.ConnectionEncryption = (*fabcore.ConnectionEncryption)(credentialDetails.ConnectionEncryption.ValueStringPointer())
	to.SingleSignOnType = (*fabcore.SingleSignOnType)(credentialDetails.SingleSignOnType.ValueStringPointer())
	to.SkipTestConnection = credentialDetails.SkipTestConnection.ValueBoolPointer()

	credentialType := (fabcore.CredentialType)(credentialDetails.CredentialType.ValueString())

	var requestUpdateCredential fabcore.CredentialsClassification

	switch credentialType {
	case fabcore.CredentialTypeAnonymous:
		requestUpdateCredential = &fabcore.AnonymousCredentials{
			CredentialType: &credentialType,
		}

	case fabcore.CredentialTypeBasic:
		cred, diags := credentialDetails.BasicCredentials.Get(ctx)
		if diags.HasError() {
			return diags
		}

		requestUpdateCredential = &fabcore.BasicCredentials{
			CredentialType: &credentialType,
			Username:       cred.Username.ValueStringPointer(),
			Password:       cred.PasswordWO.ValueStringPointer(),
		}

	case fabcore.CredentialTypeKey:
		cred, diags := credentialDetails.KeyCredentials.Get(ctx)
		if diags.HasError() {
			return diags
		}

		requestUpdateCredential = &fabcore.KeyCredentials{
			CredentialType: &credentialType,
			Key:            cred.KeyWO.ValueStringPointer(),
		}

	case fabcore.CredentialTypeOAuth2:
		requestUpdateCredential = &fabcore.Credentials{
			CredentialType: &credentialType,
		}

	case fabcore.CredentialTypeServicePrincipal:
		cred, diags := credentialDetails.ServicePrincipalCredentials.Get(ctx)
		if diags.HasError() {
			return diags
		}

		requestUpdateCredential = &fabcore.ServicePrincipalCredentials{
			CredentialType:           &credentialType,
			TenantID:                 cred.TenantID.ValueStringPointer(),
			ServicePrincipalClientID: cred.ClientID.ValueStringPointer(),
			ServicePrincipalSecret:   cred.ClientSecretWO.ValueStringPointer(),
		}

	case fabcore.CredentialTypeSharedAccessSignature:
		cred, diags := credentialDetails.SharedAccessSignatureCredentials.Get(ctx)
		if diags.HasError() {
			return diags
		}

		requestUpdateCredential = &fabcore.SharedAccessSignatureCredentials{
			CredentialType: &credentialType,
			Token:          cred.TokenWO.ValueStringPointer(),
		}

	case fabcore.CredentialTypeWindows:
		cred, diags := credentialDetails.WindowsCredentials.Get(ctx)
		if diags.HasError() {
			return diags
		}

		requestUpdateCredential = &fabcore.WindowsCredentials{
			CredentialType: &credentialType,
			Username:       cred.Username.ValueStringPointer(),
			Password:       cred.PasswordWO.ValueStringPointer(),
		}

	case fabcore.CredentialTypeWindowsWithoutImpersonation:
		requestUpdateCredential = &fabcore.WindowsWithoutImpersonationCredentials{
			CredentialType: &credentialType,
		}

	case fabcore.CredentialTypeWorkspaceIdentity:
		requestUpdateCredential = &fabcore.WorkspaceIdentityCredentials{
			CredentialType: &credentialType,
		}
	}

	to.Credentials = requestUpdateCredential

	return nil
}

type requestUpdateOnPremisesCredentialDetails struct {
	fabcore.UpdateOnPremisesGatewayCredentialDetails
}

func (to *requestUpdateOnPremisesCredentialDetails) set(
	ctx context.Context,
	from supertypes.SingleNestedObjectValueOf[rsCredentialDetailsModel],
	creds fabcore.CredentialsClassification,
	gwMembers []fabcore.OnPremisesGatewayMember,
) diag.Diagnostics {
	credentialDetails, diags := from.Get(ctx)
	if diags.HasError() {
		return diags
	}

	to.ConnectionEncryption = (*fabcore.ConnectionEncryption)(credentialDetails.ConnectionEncryption.ValueStringPointer())
	to.SingleSignOnType = (*fabcore.SingleSignOnType)(credentialDetails.SingleSignOnType.ValueStringPointer())
	to.SkipTestConnection = credentialDetails.SkipTestConnection.ValueBoolPointer()

	c, err := credentials.NewCredentials(creds)
	if err != nil {
		diags.AddError("summary string", err.Error())

		return diags
	}

	var values []fabcore.OnPremisesCredentialEntry
	for _, gwMember := range gwMembers {
		encryptedCredentials, err := credentials.EncryptCredentials(*c, *gwMember.PublicKey)
		if err != nil {
			diags.AddError("summary string", err.Error())
			return diags
		}

		values = append(values, fabcore.OnPremisesCredentialEntry{
			GatewayID:            gwMember.ID,
			EncryptedCredentials: &encryptedCredentials,
		})
	}

	to.Credentials = &fabcore.OnPremisesGatewayCredentials{
		CredentialType: &c.Type,
		Values:         values,
	}

	return nil
}

type requestUpdateConnection struct {
	fabcore.UpdateConnectionRequestClassification
}

func (to *requestUpdateConnection) set(
	ctx context.Context,
	plan, config resourceConnectionModel[rsConnectionDetailsModel, rsCredentialDetailsModel],
	gwMembers []fabcore.OnPremisesGatewayMember,
) diag.Diagnostics {
	connectivityType := (fabcore.ConnectivityType)(plan.ConnectivityType.ValueString())

	var requestUpdateCredentialDetails requestUpdateCredentialDetails
	if connectivityType == fabcore.ConnectivityTypeShareableCloud ||
		connectivityType == fabcore.ConnectivityTypeVirtualNetworkGateway || connectivityType == fabcore.ConnectivityTypeOnPremisesGateway {
		if diags := requestUpdateCredentialDetails.set(ctx, config.CredentialDetails); diags.HasError() {
			return diags
		}
	}

	displayName := plan.DisplayName.ValueStringPointer()
	privacyLevel := (*fabcore.PrivacyLevel)(plan.PrivacyLevel.ValueStringPointer())

	// var requestUpdateConnection fabcore.UpdateConnectionRequestClassification

	switch connectivityType {
	case fabcore.ConnectivityTypeShareableCloud:
		to.UpdateConnectionRequestClassification = &fabcore.UpdateShareableCloudConnectionRequest{
			DisplayName:       displayName,
			ConnectivityType:  &connectivityType,
			PrivacyLevel:      privacyLevel,
			CredentialDetails: &requestUpdateCredentialDetails.UpdateCredentialDetails,
		}
	// case fabcore.ConnectivityTypePersonalCloud:
	// 	// todo

	// 	requestUpdateConnection = &fabcore.UpdatePersonalCloudConnectionRequest{
	// 		ConnectivityType: &connectivityType,
	// 		PrivacyLevel:     privacyLevel,
	// 	}
	case fabcore.ConnectivityTypeVirtualNetworkGateway:
		to.UpdateConnectionRequestClassification = &fabcore.UpdateVirtualNetworkGatewayConnectionRequest{
			DisplayName:       displayName,
			ConnectivityType:  &connectivityType,
			PrivacyLevel:      privacyLevel,
			CredentialDetails: &requestUpdateCredentialDetails.UpdateCredentialDetails,
		}

	case fabcore.ConnectivityTypeOnPremisesGateway:
		var credentialDetails requestUpdateOnPremisesCredentialDetails

		if diags := credentialDetails.set(ctx, config.CredentialDetails, requestUpdateCredentialDetails.Credentials, gwMembers); diags.HasError() {
			return diags
		}

		to.UpdateConnectionRequestClassification = &fabcore.UpdateOnPremisesGatewayConnectionRequest{
			DisplayName:       displayName,
			ConnectivityType:  &connectivityType,
			PrivacyLevel:      privacyLevel,
			CredentialDetails: &credentialDetails.UpdateOnPremisesGatewayCredentialDetails,
		}

	// case fabcore.ConnectivityTypeOnPremisesGatewayPersonal:

	// 	requestUpdateConnection = &fabcore.UpdateOnPremisesGatewayPersonalConnectionRequest{
	// 		ConnectivityType: &connectivityType,
	// 		PrivacyLevel:     privacyLevel,
	// 	}
	case fabcore.ConnectivityTypeAutomatic: // fabcore.ConnectivityTypeNone:
		to.UpdateConnectionRequestClassification = &fabcore.UpdateConnectionRequest{
			ConnectivityType: &connectivityType,
			PrivacyLevel:     privacyLevel,
		}
	}

	// to.UpdateConnectionRequestClassification = requestUpdateConnection

	return nil
}
