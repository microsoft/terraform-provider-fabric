// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package connection

import (
	"context"

	timeoutsD "github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts" //revive:disable-line:import-alias-naming
	timeoutsR "github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"   //revive:disable-line:import-alias-naming
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/microsoft/fabric-sdk-go/fabric/core"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
)

// Terraform model definitions below

type connectionParameterModel struct {
	Name     types.String `tfsdk:"name"`
	DataType types.String `tfsdk:"data_type"`
	Value    types.String `tfsdk:"value"`
}

type connectionDetailsModel struct {
	Type           types.String               `tfsdk:"type"`
	CreationMethod types.String               `tfsdk:"creation_method"`
	Parameters     []connectionParameterModel `tfsdk:"parameters"`
}

type connectionCredentialsModel struct {
	CredentialType types.String `tfsdk:"credential_type"`
	// For BasicCredentials
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
	// For KeyCredentials
	Key types.String `tfsdk:"key"`
	// For ServicePrincipalCredentials
	ApplicationID     types.String `tfsdk:"application_id"`
	ApplicationSecret types.String `tfsdk:"application_secret"`
	TenantID          types.String `tfsdk:"tenant_id"`
	// For SharedAccessSignatureCredentials
	SasToken types.String `tfsdk:"sas_token"`
	// For WindowsCredentials and WindowsWithoutImpersonationCredentials
	Domain types.String `tfsdk:"domain"`
	// For WorkspaceIdentityCredentials - no additional fields required
	// For AnonymousCredentials - no additional fields required
}

type connectionCredentialDetailsModel struct {
	SingleSignOnType     types.String               `tfsdk:"single_sign_on_type"`
	ConnectionEncryption types.String               `tfsdk:"connection_encryption"`
	SkipTestConnection   types.Bool                 `tfsdk:"skip_test_connection"`
	Credentials          connectionCredentialsModel `tfsdk:"credentials"`
}

type connectionPropertiesModel struct {
	ConnectivityType  types.String     `tfsdk:"connectivity_type"`
	GatewayID         customtypes.UUID `tfsdk:"gateway_id"`
	PrivacyLevel      types.String     `tfsdk:"privacy_level"`
	ConnectionDetails types.Object     `tfsdk:"connection_details"`
	CredentialDetails types.Object     `tfsdk:"credential_details"`
}

// Single data source model.
type baseConnectionModel struct {
	ID                customtypes.UUID `tfsdk:"id"`
	DisplayName       types.String     `tfsdk:"display_name"`
	ConnectivityType  types.String     `tfsdk:"connectivity_type"`
	GatewayID         customtypes.UUID `tfsdk:"gateway_id"`
	PrivacyLevel      types.String     `tfsdk:"privacy_level"`
	ConnectionDetails types.Object     `tfsdk:"connection_details"`
	CredentialDetails types.Object     `tfsdk:"credential_details"`
}

type dataSourceConnectionModel struct {
	baseConnectionModel
	Timeouts timeoutsD.Value `tfsdk:"timeouts"`
}

type resourceConnectionModel struct {
	baseConnectionModel
	Timeouts timeoutsR.Value `tfsdk:"timeouts"`
}

type requestCreateConnection struct {
	core.CreateConnectionRequestClassification
}

type requestUpdateConnection struct {
	core.UpdateConnectionRequestClassification
}

func (to *requestUpdateConnection) set(plan resourceConnectionModel) diag.Diagnostics {
	// Create the update request
	updateRequest := &core.UpdateConnectionRequest{}

	// Set connectivity type if provided
	if !plan.ConnectivityType.IsNull() && !plan.ConnectivityType.IsUnknown() {
		connType := core.ConnectivityType(plan.ConnectivityType.ValueString())
		updateRequest.ConnectivityType = &connType
	}

	// Set privacy level if provided
	if !plan.PrivacyLevel.IsNull() && !plan.PrivacyLevel.IsUnknown() {
		privacyLevel := core.PrivacyLevel(plan.PrivacyLevel.ValueString())
		updateRequest.PrivacyLevel = &privacyLevel
	}

	to.UpdateConnectionRequestClassification = updateRequest

	return nil
}

func (to *requestCreateConnection) set(ctx context.Context, from resourceConnectionModel) diag.Diagnostics {
	// Unmarshal connection details from types.Object to our struct
	var connectionDetailsModel connectionDetailsModel
	diags := from.ConnectionDetails.As(ctx, &connectionDetailsModel, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		return diags
	}

	// Unmarshal credential details from types.Object to our struct
	var credentialDetailsModel connectionCredentialDetailsModel
	diags = from.CredentialDetails.As(ctx, &credentialDetailsModel, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		return diags
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
		username := credentialDetailsModel.Credentials.Username.ValueString()
		password := credentialDetailsModel.Credentials.Password.ValueString()
		credentials = &core.BasicCredentials{
			CredentialType: &credType,
			Username:       &username,
			Password:       &password,
		}
	case core.CredentialTypeWindows, core.CredentialTypeWindowsWithoutImpersonation:
		username := credentialDetailsModel.Credentials.Username.ValueString()
		password := credentialDetailsModel.Credentials.Password.ValueString()

		if credType == core.CredentialTypeWindows {
			credentials = &core.WindowsCredentials{
				CredentialType: &credType,
				Username:       &username,
				Password:       &password,
			}
		} else {
			credentials = &core.WindowsWithoutImpersonationCredentials{
				CredentialType: &credType,
			}
		}
	case core.CredentialTypeKey:
		key := credentialDetailsModel.Credentials.Key.ValueString()
		credentials = &core.KeyCredentials{
			CredentialType: &credType,
			Key:            &key,
		}
	case core.CredentialTypeServicePrincipal:
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
	displayName := from.DisplayName.ValueString()
	connType := core.ConnectivityType(from.ConnectivityType.ValueString())
	connDetailsType := connectionDetailsModel.Type.ValueString()
	creationMethod := connectionDetailsModel.CreationMethod.ValueString()
	privacyLevel := core.PrivacyLevel(from.PrivacyLevel.ValueString())

	createConnRequest := &core.CreateConnectionRequest{
		DisplayName:      &displayName,
		ConnectivityType: &connType,
		ConnectionDetails: &core.CreateConnectionDetails{
			Type:           &connDetailsType,
			CreationMethod: &creationMethod,
			Parameters:     parameters,
		},
		PrivacyLevel: &privacyLevel,
	}

	// Create credential details for the connection
	credDetails := &core.CreateCredentialDetails{
		Credentials:          credentials,
		ConnectionEncryption: &connEncryption,
		SingleSignOnType:     &ssoType,
		SkipTestConnection:   &skipTest,
	}

	// Since GatewayID isn't directly available on CreateConnectionRequest,
	// we'll need to determine the appropriate type of CreateConnectionRequestClassification to use
	// For now, we'll use the base CreateConnectionRequest type
	if !from.GatewayID.IsNull() && !from.GatewayID.IsUnknown() {
		// If we have a gateway ID, we need to create an appropriate connection request with the gateway ID
		// For now we're using a Cloud connection request
		cloudRequest := &core.CreateCloudConnectionRequest{
			ConnectionDetails: createConnRequest.ConnectionDetails,
			ConnectivityType:  createConnRequest.ConnectivityType,
			CredentialDetails: credDetails,
			DisplayName:       createConnRequest.DisplayName,
			PrivacyLevel:      createConnRequest.PrivacyLevel,
		}

		to.CreateConnectionRequestClassification = cloudRequest
	} else {
		// Create a cloud connection by default
		cloudRequest := &core.CreateCloudConnectionRequest{
			ConnectionDetails: createConnRequest.ConnectionDetails,
			ConnectivityType:  createConnRequest.ConnectivityType,
			CredentialDetails: credDetails,
			DisplayName:       createConnRequest.DisplayName,
			PrivacyLevel:      createConnRequest.PrivacyLevel,
		}

		to.CreateConnectionRequestClassification = cloudRequest
	}

	return nil
}

func (to *baseConnectionModel) set(ctx context.Context, from core.Connection) diag.Diagnostics {
	to.ID = customtypes.NewUUIDPointerValue(from.ID)
	to.DisplayName = types.StringPointerValue(from.DisplayName)
	to.ConnectivityType = types.StringValue(string(*from.ConnectivityType))
	to.GatewayID = customtypes.NewUUIDPointerValue(from.GatewayID)
	to.PrivacyLevel = types.StringValue(string(*from.PrivacyLevel))

	// Important: Do NOT overwrite the connection_details and credential_details objects
	// from the API as they don't contain sensitive values
	// Instead, keep the values from the config/plan to maintain any sensitive values

	// We only set these if they're empty (like during import)
	if to.ConnectionDetails.IsNull() || to.ConnectionDetails.IsUnknown() {
		// Create a minimal connection details object
		if from.ConnectionDetails != nil {
			connDetailsModel := connectionDetailsModel{
				Type:           types.StringValue(string(*from.ConnectionDetails.Type)),
				CreationMethod: types.StringValue(""),        // Not available in API response
				Parameters:     []connectionParameterModel{}, // Not available in API response
			}

			// Convert the model to a types.Object
			connDetailsObj, diags := types.ObjectValueFrom(ctx,
				map[string]attr.Type{
					"type":            types.StringType,
					"creation_method": types.StringType,
					"parameters": types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{
						"name":      types.StringType,
						"data_type": types.StringType,
						"value":     types.StringType,
					}}},
				},
				connDetailsModel)
			if diags.HasError() {
				return diags
			}

			to.ConnectionDetails = connDetailsObj
		}
	} else {
		// Do not update ConnectionDetails at all - preserve all values from the current state, including the 
		// type, creation_method, and parameters that the user has defined
		// This prevents inconsistencies after apply where values might vanish
	}

	// We only set credential details if they're empty (like during import)
	if to.CredentialDetails.IsNull() || to.CredentialDetails.IsUnknown() {
		// Create a minimal credential details object
		if from.CredentialDetails != nil {
			credDetailsModel := connectionCredentialDetailsModel{
				SingleSignOnType:     types.StringValue(string(*from.CredentialDetails.SingleSignOnType)),
				ConnectionEncryption: types.StringValue(string(*from.CredentialDetails.ConnectionEncryption)),
				SkipTestConnection:   types.BoolValue(*from.CredentialDetails.SkipTestConnection),
				Credentials: connectionCredentialsModel{
					CredentialType: types.StringValue(string(*from.CredentialDetails.CredentialType)),
					// Note: We don't get sensitive credential values from the API
				},
			}

			// Convert the model to a types.Object
			credDetailsObj, diags := types.ObjectValueFrom(ctx,
				map[string]attr.Type{
					"single_sign_on_type":   types.StringType,
					"connection_encryption": types.StringType,
					"skip_test_connection":  types.BoolType,
					"credentials": types.ObjectType{AttrTypes: map[string]attr.Type{
						"credential_type":    types.StringType,
						"username":           types.StringType,
						"password":           types.StringType,
						"key":                types.StringType,
						"application_id":     types.StringType,
						"application_secret": types.StringType,
						"tenant_id":          types.StringType,
						"sas_token":          types.StringType,
						"domain":             types.StringType,
					}},
				},
				credDetailsModel)
			if diags.HasError() {
				return diags
			}

			to.CredentialDetails = credDetailsObj
		}
	} else {
		// Do not update CredentialDetails at all - preserve all values from the current state
		// This prevents inconsistencies with sensitive values after apply
	}

	return nil
}

// List data source model.
type dataSourceConnectionsModel struct {
	WorkspaceID customtypes.UUID                                       `tfsdk:"workspace_id"`
	Connections []core.Item                                            `tfsdk:"-"` // Internal use only
	Values      supertypes.SetNestedObjectValueOf[baseConnectionModel] `tfsdk:"values"`
}

func (to *dataSourceConnectionsModel) setValues(ctx context.Context, from []core.Connection) diag.Diagnostics {
	slice := make([]*baseConnectionModel, 0, len(from))

	for _, entity := range from {
		var entityModel baseConnectionModel

		if diags := entityModel.set(ctx, entity); diags.HasError() {
			return diags
		}

		slice = append(slice, &entityModel)
	}

	return to.Values.Set(ctx, slice)
}
