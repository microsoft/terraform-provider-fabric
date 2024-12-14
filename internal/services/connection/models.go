// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package connection

import (
	"context"

	supertypes "github.com/FrangipaneTeam/terraform-plugin-framework-supertypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
)

type baseConnectionModel struct {
	ID                customtypes.UUID                                             `tfsdk:"id"`
	DisplayName       types.String                                                 `tfsdk:"display_name"`
	GatewayID         customtypes.UUID                                             `tfsdk:"gateway_id"`
	ConnectivityType  types.String                                                 `tfsdk:"connectivity_type"`
	PrivacyLevel      types.String                                                 `tfsdk:"privacy_level"`
	ConnectionDetails supertypes.SingleNestedObjectValueOf[connectionDetailsModel] `tfsdk:"connection_details"`
	CredentialDetails supertypes.SingleNestedObjectValueOf[credentialDetailsModel] `tfsdk:"credential_details"`
}

func (to *baseConnectionModel) set(ctx context.Context, from fabcore.Connection) diag.Diagnostics {
	to.ID = customtypes.NewUUIDPointerValue(from.ID)
	to.DisplayName = types.StringPointerValue(from.DisplayName)
	to.GatewayID = customtypes.NewUUIDPointerValue(from.GatewayID)
	to.ConnectivityType = types.StringPointerValue((*string)(from.ConnectivityType))
	to.PrivacyLevel = types.StringPointerValue((*string)(from.PrivacyLevel))

	connectionDetails := supertypes.NewSingleNestedObjectValueOfNull[connectionDetailsModel](ctx)

	if from.ConnectionDetails != nil {
		connectionDetailsModel := &connectionDetailsModel{}
		connectionDetailsModel.set(*from.ConnectionDetails)

		diags := connectionDetails.Set(ctx, connectionDetailsModel)
		if diags.HasError() {
			return diags
		}
	}

	to.ConnectionDetails = connectionDetails

	credentialDetails := supertypes.NewSingleNestedObjectValueOfNull[credentialDetailsModel](ctx)

	if from.CredentialDetails != nil {
		credentialDetailsModel := &credentialDetailsModel{}
		credentialDetailsModel.set(*from.CredentialDetails)

		diags := credentialDetails.Set(ctx, credentialDetailsModel)
		if diags.HasError() {
			return diags
		}
	}

	to.CredentialDetails = credentialDetails

	return nil
}

type connectionDetailsModel struct {
	Path types.String `tfsdk:"path"`
	Type types.String `tfsdk:"type"`
}

func (to *connectionDetailsModel) set(from fabcore.ListConnectionDetails) {
	to.Path = types.StringPointerValue(from.Path)
	to.Type = types.StringPointerValue(from.Type)
}

type credentialDetailsModel struct {
	ConnectionEncryption types.String `tfsdk:"connection_encryption"`
	CredentialType       types.String `tfsdk:"credential_type"`
	SingleSignOnType     types.String `tfsdk:"single_sign_on_type"`
	SkipTestConnection   types.Bool   `tfsdk:"skip_test_connection"`
}

func (to *credentialDetailsModel) set(from fabcore.ListCredentialDetails) {
	to.ConnectionEncryption = types.StringPointerValue((*string)(from.ConnectionEncryption))
	to.CredentialType = types.StringPointerValue((*string)(from.CredentialType))
	to.SingleSignOnType = types.StringPointerValue((*string)(from.SingleSignOnType))
	to.SkipTestConnection = types.BoolPointerValue(from.SkipTestConnection)
}
