// Copyright Microsoft Corporation 2026
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
	ID                            customtypes.UUID                                        `tfsdk:"id"`
	DisplayName                   types.String                                            `tfsdk:"display_name"`
	GatewayID                     customtypes.UUID                                        `tfsdk:"gateway_id"`
	ConnectivityType              types.String                                            `tfsdk:"connectivity_type"`
	PrivacyLevel                  types.String                                            `tfsdk:"privacy_level"`
	AllowConnectionUsageInGateway types.Bool                                              `tfsdk:"allow_connection_usage_in_gateway"`
	ConnectionDetails             supertypes.SingleNestedObjectValueOf[ConnectionDetails] `tfsdk:"connection_details"`
	CredentialDetails             supertypes.SingleNestedObjectValueOf[CredentialDetails] `tfsdk:"credential_details"`
}

func (to *baseConnectionModel[ConnectionDetails, CredentialDetails]) set(ctx context.Context, from fabcore.ConnectionClassification) diag.Diagnostics { //nolint:gocognit
	fromConnection := from.GetConnection()
	to.ID = customtypes.NewUUIDPointerValue(fromConnection.ID)
	to.DisplayName = types.StringPointerValue(fromConnection.DisplayName)
	to.ConnectivityType = types.StringPointerValue((*string)(fromConnection.ConnectivityType))
	to.PrivacyLevel = types.StringPointerValue((*string)(fromConnection.PrivacyLevel))

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
		var defaultModel ConnectionDetails
		connectionDetailsModel = &defaultModel
	}

	if fromConnection.ConnectionDetails != nil {
		var connectionDetailsModelPtr *ConnectionDetails

		switch v := any(connectionDetailsModel).(type) {
		case *dsConnectionDetailsModel:
			setDSConnectionDetails(*fromConnection.ConnectionDetails, v)

			if convertedValue, ok := any(*v).(ConnectionDetails); ok {
				connectionDetailsModelPtr = &convertedValue
			}
		case *rsConnectionDetailsModel:
			setRSConnectionDetails(*fromConnection.ConnectionDetails, v)

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
		var defaultModel CredentialDetails
		credentialDetailsModel = &defaultModel
	}

	if fromConnection.CredentialDetails != nil {
		var credentialDetailsModelPtr *CredentialDetails

		switch v := any(credentialDetailsModel).(type) {
		case *dsCredentialDetailsModel:
			setDSCredentialDetails(*fromConnection.CredentialDetails, v)

			if convertedValue, ok := any(*v).(CredentialDetails); ok {
				credentialDetailsModelPtr = &convertedValue
			}
		case *rsCredentialDetailsModel:
			setRSCredentialDetails(*fromConnection.CredentialDetails, v)

			if convertedValue, ok := any(*v).(CredentialDetails); ok {
				credentialDetailsModelPtr = &convertedValue
			}
		}

		if diags := credentialDetails.Set(ctx, credentialDetailsModelPtr); diags.HasError() {
			return diags
		}
	}

	to.CredentialDetails = credentialDetails

	// connectivity type specific information
	switch v := from.(type) {
	case *fabcore.ShareableCloudConnection:
		if v.AllowConnectionUsageInGateway != nil {
			to.AllowConnectionUsageInGateway = types.BoolValue(*v.AllowConnectionUsageInGateway)
		} else {
			to.AllowConnectionUsageInGateway = types.BoolValue(false) // default is false
		}
	case *fabcore.VirtualNetworkGatewayConnection:
		to.GatewayID = customtypes.NewUUIDPointerValue(v.GatewayID)
		// keep it here due to default being "false"
		to.AllowConnectionUsageInGateway = types.BoolNull()
	}

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
