// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package connection

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
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

	BasicCredentials                 supertypes.SingleNestedObjectValueOf[credentialsBasicModel]                 `tfsdk:"basic_credentials"`
	KeyCredentials                   supertypes.SingleNestedObjectValueOf[credentialsKeyModel]                   `tfsdk:"key_credentials"`
	ServicePrincipalCredentials      supertypes.SingleNestedObjectValueOf[credentialsServicePrincipalModel]      `tfsdk:"service_principal_credentials"`
	SharedAccessSignatureCredentials supertypes.SingleNestedObjectValueOf[credentialsSharedAccessSignatureModel] `tfsdk:"shared_access_signature_credentials"`
}

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
) diag.Diagnostics {
	connectivityType := (fabcore.ConnectivityType)(plan.ConnectivityType.ValueString())

	var requestCreateConnectionDetails requestCreateConnectionDetails
	if diags := requestCreateConnectionDetails.set(ctx, plan.ConnectionDetails); diags.HasError() {
		return diags
	}

	var requestCreateCredentialDetails requestCreateCredentialDetails
	if diags := requestCreateCredentialDetails.set(ctx, config.CredentialDetails); diags.HasError() {
		return diags
	}

	displayName := plan.DisplayName.ValueStringPointer()
	privacyLevel := (*fabcore.PrivacyLevel)(plan.PrivacyLevel.ValueStringPointer())

	switch connectivityType {
	case fabcore.ConnectivityTypeShareableCloud: // fabcore.ConnectivityTypePersonalCloud:
		to.CreateConnectionRequestClassification = &fabcore.CreateCloudConnectionRequest{
			DisplayName:                   displayName,
			PrivacyLevel:                  privacyLevel,
			ConnectivityType:              &connectivityType,
			ConnectionDetails:             &requestCreateConnectionDetails.CreateConnectionDetails,
			CredentialDetails:             &requestCreateCredentialDetails.CreateCredentialDetails,
			AllowConnectionUsageInGateway: plan.AllowConnectionUsageInGateway.ValueBoolPointer(),
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
	default:
		var diags diag.Diagnostics
		diags.AddError(
			"Unsupported connectivity type",
			fmt.Sprintf("Connectivity type '%s' is not supported", connectivityType),
		)

		return diags
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

					return diags
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
			default:
				diags.AddError(
					"Unsupported data type",
					fmt.Sprintf("Data type '%s' is not supported", dataType),
				)

				return diags
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

	case fabcore.CredentialTypeWorkspaceIdentity:
		requestCreateCredential = &fabcore.WorkspaceIdentityCredentials{
			CredentialType: &credentialType,
		}
	default:
		var diags diag.Diagnostics
		diags.AddError(
			"Unsupported credential type",
			fmt.Sprintf("Credential type '%s' is not supported", credentialType),
		)

		return diags
	}

	to.Credentials = requestCreateCredential

	return nil
}

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

	case fabcore.CredentialTypeWorkspaceIdentity:
		requestUpdateCredential = &fabcore.WorkspaceIdentityCredentials{
			CredentialType: &credentialType,
		}
	default:
		var diags diag.Diagnostics
		diags.AddError(
			"Unsupported credential type",
			fmt.Sprintf("Credential type '%s' is not supported", credentialType),
		)

		return diags
	}

	to.Credentials = requestUpdateCredential

	return nil
}

type requestUpdateConnection struct {
	fabcore.UpdateConnectionRequestClassification
}

func (to *requestUpdateConnection) set(
	ctx context.Context,
	plan, config resourceConnectionModel[rsConnectionDetailsModel, rsCredentialDetailsModel],
) diag.Diagnostics {
	connectivityType := (fabcore.ConnectivityType)(plan.ConnectivityType.ValueString())

	var requestUpdateCredentialDetails requestUpdateCredentialDetails
	if connectivityType == fabcore.ConnectivityTypeShareableCloud ||
		connectivityType == fabcore.ConnectivityTypeVirtualNetworkGateway {
		if diags := requestUpdateCredentialDetails.set(ctx, config.CredentialDetails); diags.HasError() {
			return diags
		}
	}

	displayName := plan.DisplayName.ValueStringPointer()
	privacyLevel := (*fabcore.PrivacyLevel)(plan.PrivacyLevel.ValueStringPointer())

	switch connectivityType {
	case fabcore.ConnectivityTypeShareableCloud:
		to.UpdateConnectionRequestClassification = &fabcore.UpdateShareableCloudConnectionRequest{
			DisplayName:                   displayName,
			ConnectivityType:              &connectivityType,
			PrivacyLevel:                  privacyLevel,
			CredentialDetails:             &requestUpdateCredentialDetails.UpdateCredentialDetails,
			AllowConnectionUsageInGateway: plan.AllowConnectionUsageInGateway.ValueBoolPointer(),
		}

	case fabcore.ConnectivityTypeVirtualNetworkGateway:
		to.UpdateConnectionRequestClassification = &fabcore.UpdateVirtualNetworkGatewayConnectionRequest{
			DisplayName:       displayName,
			ConnectivityType:  &connectivityType,
			PrivacyLevel:      privacyLevel,
			CredentialDetails: &requestUpdateCredentialDetails.UpdateCredentialDetails,
		}
	default:
		var diags diag.Diagnostics
		diags.AddError(
			"Unsupported connectivity type",
			fmt.Sprintf("Connectivity type '%s' is not supported", connectivityType),
		)

		return diags
	}

	return nil
}
