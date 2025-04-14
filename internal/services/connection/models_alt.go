// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package connection

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoft/fabric-sdk-go/fabric/core"
)

// Helper function to try and convert a credential_details object to this model
func tryConvertCredentialDetails(ctx context.Context, obj types.Object) (connectionCredentialDetailsModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	var result connectionCredentialDetailsModel

	diags = obj.As(ctx, &result, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		diags.AddError(
			"Unable to parse credential_details",
			"Failed to parse credential_details",
		)
	}

	return result, diags
}

// Helper function to try and convert a connection_details object to this model
func tryConvertConnectionDetails(ctx context.Context, obj types.Object) (connectionDetailsModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	var result connectionDetailsModel

	diags = obj.As(ctx, &result, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		diags.AddError(
			"Unable to parse connection_details",
			"Failed to parse connection_details",
		)
	}

	return result, diags
}

// Helper function to build a connection create request from either object or block syntax
func buildConnectionRequest(ctx context.Context,
	connectionDetails types.Object,
	credentialDetails types.Object,
	displayName string,
	connectivityType string,
	privacyLevel string,
	gatewayId string,
) (*requestCreateConnection, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Parse connection details
	connectionDetailsModel, parseDiags := tryConvertConnectionDetails(ctx, connectionDetails)
	diags.Append(parseDiags...)
	if diags.HasError() {
		return nil, diags
	}

	// Parse credential details
	credentialDetailsModel, parseDiags := tryConvertCredentialDetails(ctx, credentialDetails)
	diags.Append(parseDiags...)
	if diags.HasError() {
		return nil, diags
	}

	// Create connection parameters
	parameters := make([]core.ConnectionDetailsParameterClassification, 0, len(connectionDetailsModel.Parameters))

	for _, param := range connectionDetailsModel.Parameters {
		name := param.Name.ValueString()
		dataType := param.DataType.ValueString()
		value := param.Value.ValueString()

		// For now, we only support text parameters
		textParam := core.ConnectionDetailsTextParameter{
			Name:     &name,
			DataType: (*core.DataType)(&dataType),
			Value:    &value,
		}
		parameters = append(parameters, &textParam)
	}

	// Prepare credential details
	credType := core.CredentialType(credentialDetailsModel.Credentials.CredentialType.ValueString())
	ssoType := core.SingleSignOnType(credentialDetailsModel.SingleSignOnType.ValueString())
	connEncryption := core.ConnectionEncryption(credentialDetailsModel.ConnectionEncryption.ValueString())
	skipTest := credentialDetailsModel.SkipTestConnection.ValueBool()

	// Create the appropriate credentials object based on the credential type
	var credentials core.CredentialsClassification

	// Set the appropriate credential fields based on the credential type
	switch credType {
	case core.CredentialTypeBasic:
		// Validate required fields for BasicCredentials
		if credentialDetailsModel.Credentials.Username.IsNull() {
			diags.AddError(
				"Missing required field for BasicCredentials",
				"username is required when credential_type is Basic",
			)
		}
		if credentialDetailsModel.Credentials.Password.IsNull() {
			diags.AddError(
				"Missing required field for BasicCredentials",
				"password is required when credential_type is Basic",
			)
		}
		if diags.HasError() {
			return nil, diags
		}

		username := credentialDetailsModel.Credentials.Username.ValueString()
		password := credentialDetailsModel.Credentials.Password.ValueString()
		credentials = &core.BasicCredentials{
			CredentialType: &credType,
			Username:       &username,
			Password:       &password,
		}
	case core.CredentialTypeWindows:
		// Validate required fields for WindowsCredentials
		if credentialDetailsModel.Credentials.Username.IsNull() {
			diags.AddError(
				"Missing required field for WindowsCredentials",
				"username is required when credential_type is Windows",
			)
		}
		if credentialDetailsModel.Credentials.Password.IsNull() {
			diags.AddError(
				"Missing required field for WindowsCredentials",
				"password is required when credential_type is Windows",
			)
		}
		if diags.HasError() {
			return nil, diags
		}

		username := credentialDetailsModel.Credentials.Username.ValueString()
		password := credentialDetailsModel.Credentials.Password.ValueString()

		credentials = &core.WindowsCredentials{
			CredentialType: &credType,
			Username:       &username,
			Password:       &password,
		}
	case core.CredentialTypeWindowsWithoutImpersonation:
		credentials = &core.WindowsWithoutImpersonationCredentials{
			CredentialType: &credType,
		}
	case core.CredentialTypeKey:
		// Validate required fields for KeyCredentials
		if credentialDetailsModel.Credentials.Key.IsNull() {
			diags.AddError(
				"Missing required field for KeyCredentials",
				"key is required when credential_type is Key",
			)
			return nil, diags
		}

		key := credentialDetailsModel.Credentials.Key.ValueString()
		credentials = &core.KeyCredentials{
			CredentialType: &credType,
			Key:            &key,
		}
	case core.CredentialTypeServicePrincipal:
		// Validate required fields for ServicePrincipalCredentials
		if credentialDetailsModel.Credentials.ApplicationID.IsNull() {
			diags.AddError(
				"Missing required field for ServicePrincipalCredentials",
				"application_id is required when credential_type is ServicePrincipal",
			)
		}
		if credentialDetailsModel.Credentials.ApplicationSecret.IsNull() {
			diags.AddError(
				"Missing required field for ServicePrincipalCredentials",
				"application_secret is required when credential_type is ServicePrincipal",
			)
		}
		if credentialDetailsModel.Credentials.TenantID.IsNull() {
			diags.AddError(
				"Missing required field for ServicePrincipalCredentials",
				"tenant_id is required when credential_type is ServicePrincipal",
			)
		}
		if diags.HasError() {
			return nil, diags
		}

		// Use application_id and application_secret directly
		appID := credentialDetailsModel.Credentials.ApplicationID.ValueString()
		appSecret := credentialDetailsModel.Credentials.ApplicationSecret.ValueString()
		tenantID := credentialDetailsModel.Credentials.TenantID.ValueString()
		credentials = &core.ServicePrincipalCredentials{
			CredentialType:           &credType,
			ServicePrincipalClientID: &appID,
			ServicePrincipalSecret:   &appSecret,
			TenantID:                 &tenantID,
		}
	case core.CredentialTypeSharedAccessSignature:
		// Validate required fields for SharedAccessSignatureCredentials
		if credentialDetailsModel.Credentials.SasToken.IsNull() {
			diags.AddError(
				"Missing required field for SharedAccessSignatureCredentials",
				"sas_token is required when credential_type is SharedAccessSignature",
			)
			return nil, diags
		}

		sasToken := credentialDetailsModel.Credentials.SasToken.ValueString()
		credentials = &core.SharedAccessSignatureCredentials{
			CredentialType: &credType,
			Token:          &sasToken,
		}
	case core.CredentialTypeWorkspaceIdentity:
		credentials = &core.WorkspaceIdentityCredentials{
			CredentialType: &credType,
		}
	case core.CredentialTypeAnonymous:
		credentials = &core.AnonymousCredentials{
			CredentialType: &credType,
		}
	default:
		// Default to base credentials type with just the credential type
		credentials = &core.Credentials{
			CredentialType: &credType,
		}
	}

	// Create the request
	connType := core.ConnectivityType(connectivityType)
	connDetailsType := connectionDetailsModel.Type.ValueString()
	creationMethod := connectionDetailsModel.CreationMethod.ValueString()
	privacyLevelVal := core.PrivacyLevel(privacyLevel)

	// Make sure gateway_id is provided if required
	if (connType == core.ConnectivityTypeOnPremisesGateway || connType == core.ConnectivityTypeVirtualNetworkGateway) && gatewayId == "" {
		diags.AddError(
			"Missing required field for OnPremisesGateway/VirtualNetworkGateway",
			"gateway_id is required when connectivity_type is OnPremisesGateway or VirtualNetworkGateway",
		)
		return nil, diags
	}

	// Create the connection details
	createConnDetails := &core.CreateConnectionDetails{
		Type:           &connDetailsType,
		CreationMethod: &creationMethod,
		Parameters:     parameters,
	}

	// Create standard credential details
	credDetails := &core.CreateCredentialDetails{
		Credentials:          credentials,
		ConnectionEncryption: &connEncryption,
		SingleSignOnType:     &ssoType,
		SkipTestConnection:   &skipTest,
	}

	// Prepare the result
	var result requestCreateConnection

	// Create the appropriate request type based on the connectivity type
	switch connType {
	case core.ConnectivityTypeOnPremisesGateway:
		// For OnPremisesGateway, we need to use CreateOnPremisesConnectionRequest with OnPremisesGatewayCredentials
		// Special handling for on-premises gateway credentials
		onPremCredentials := &core.OnPremisesGatewayCredentials{
			CredentialType: &credType,
			Values: []core.OnPremisesCredentialEntry{},
		}

		// Create on-premises credential details
		onPremCredDetails := &core.CreateOnPremisesCredentialDetails{
			Credentials:          onPremCredentials,
			ConnectionEncryption: &connEncryption,
			SingleSignOnType:     &ssoType,
			SkipTestConnection:   &skipTest,
		}

		// Create the on-premises connection request
		onPremRequest := &core.CreateOnPremisesConnectionRequest{
			ConnectionDetails: createConnDetails,
			ConnectivityType:  &connType,
			CredentialDetails: onPremCredDetails,
			DisplayName:       &displayName,
			GatewayID:         &gatewayId,
			PrivacyLevel:      &privacyLevelVal,
		}

		result.CreateConnectionRequestClassification = onPremRequest

	case core.ConnectivityTypeVirtualNetworkGateway:
		// For VirtualNetworkGateway, use CreateVirtualNetworkGatewayConnectionRequest
		vnGatewayRequest := &core.CreateVirtualNetworkGatewayConnectionRequest{
			ConnectionDetails: createConnDetails,
			ConnectivityType:  &connType,
			CredentialDetails: credDetails,
			DisplayName:       &displayName,
			GatewayID:         &gatewayId,
			PrivacyLevel:      &privacyLevelVal,
		}

		result.CreateConnectionRequestClassification = vnGatewayRequest

	default: // ShareableCloud or other types
		// Use the standard CreateCloudConnectionRequest for other types (including ShareableCloud)
		cloudRequest := &core.CreateCloudConnectionRequest{
			ConnectionDetails: createConnDetails,
			ConnectivityType:  &connType,
			CredentialDetails: credDetails,
			DisplayName:       &displayName,
			PrivacyLevel:      &privacyLevelVal,
		}

		result.CreateConnectionRequestClassification = cloudRequest
	}

	// Log the request type we're creating
	tflog.Info(ctx, "Creating connection request", map[string]interface{}{
		"connectivity_type": connectivityType,
		"request_type":      fmt.Sprintf("%T", result.CreateConnectionRequestClassification),
	})

	return &result, diags
}
