// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package connection

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
)

type dataSourceConnectionModel struct {
	baseDataSourceConnectionModel
	Timeouts timeouts.Value `tfsdk:"timeouts"`
}

type baseDataSourceConnectionModel struct {
	ID                customtypes.UUID                                               `tfsdk:"id"`
	DisplayName       types.String                                                   `tfsdk:"display_name"`
	GatewayID         customtypes.UUID                                               `tfsdk:"gateway_id"`
	ConnectivityType  types.String                                                   `tfsdk:"connectivity_type"`
	PrivacyLevel      types.String                                                   `tfsdk:"privacy_level"`
	ConnectionDetails supertypes.SingleNestedObjectValueOf[dsConnectionDetailsModel] `tfsdk:"connection_details"`
	CredentialDetails supertypes.SingleNestedObjectValueOf[dsCredentialDetailsModel] `tfsdk:"credential_details"`
}

func (to *baseDataSourceConnectionModel) set(ctx context.Context, from fabcore.Connection) diag.Diagnostics {
	to.ID = customtypes.NewUUIDPointerValue(from.ID)
	to.DisplayName = types.StringPointerValue(from.DisplayName)
	to.GatewayID = customtypes.NewUUIDPointerValue(from.GatewayID)
	to.ConnectivityType = types.StringPointerValue((*string)(from.ConnectivityType))
	to.PrivacyLevel = types.StringPointerValue((*string)(from.PrivacyLevel))

	connectionDetails := supertypes.NewSingleNestedObjectValueOfNull[dsConnectionDetailsModel](ctx)

	if from.ConnectionDetails != nil {
		connectionDetailsModel := &dsConnectionDetailsModel{}
		connectionDetailsModel.set(*from.ConnectionDetails)

		if diags := connectionDetails.Set(ctx, connectionDetailsModel); diags.HasError() {
			return diags
		}
	}

	to.ConnectionDetails = connectionDetails

	credentialDetails := supertypes.NewSingleNestedObjectValueOfNull[dsCredentialDetailsModel](ctx)

	if from.CredentialDetails != nil {
		credentialDetailsModel := &dsCredentialDetailsModel{}
		credentialDetailsModel.set(*from.CredentialDetails)

		if diags := credentialDetails.Set(ctx, credentialDetailsModel); diags.HasError() {
			return diags
		}
	}

	to.CredentialDetails = credentialDetails

	return nil
}

// func (to *baseDataSourceConnectionModel) setConnectionDetails(ctx context.Context, from *fabcore.ListConnectionDetails) diag.Diagnostics {
// 	connectionDetails := supertypes.NewSingleNestedObjectValueOfNull[dsConnectionDetailsModel](ctx)

// 	if from != nil {
// 		connectionDetailsModel := &dsConnectionDetailsModel{}
// 		connectionDetailsModel.set(*from)

// 		if diags := connectionDetails.Set(ctx, connectionDetailsModel); diags.HasError() {
// 			return diags
// 		}
// 	}

// 	to.ConnectionDetails = connectionDetails

// 	return nil
// }

// func (to *baseDataSourceConnectionModel) setCredentialDetails(ctx context.Context, from *fabcore.ListCredentialDetails) diag.Diagnostics {
// 	credentialDetails := supertypes.NewSingleNestedObjectValueOfNull[dsCredentialDetailsModel](ctx)

// 	if from != nil {
// 		credentialDetailsModel := &dsCredentialDetailsModel{}
// 		credentialDetailsModel.set(*from)

// 		if diags := credentialDetails.Set(ctx, credentialDetailsModel); diags.HasError() {
// 			return diags
// 		}
// 	}

// 	to.CredentialDetails = credentialDetails

// 	return nil
// }

type dsConnectionDetailsModel struct {
	Path types.String `tfsdk:"path"`
	Type types.String `tfsdk:"type"`
}

func (to *dsConnectionDetailsModel) set(from fabcore.ListConnectionDetails) {
	to.Path = types.StringPointerValue(from.Path)
	to.Type = types.StringPointerValue(from.Type)
}

type dsCredentialDetailsModel struct {
	ConnectionEncryption types.String `tfsdk:"connection_encryption"`
	CredentialType       types.String `tfsdk:"credential_type"`
	SingleSignOnType     types.String `tfsdk:"single_sign_on_type"`
	SkipTestConnection   types.Bool   `tfsdk:"skip_test_connection"`
}

func (to *dsCredentialDetailsModel) set(from fabcore.ListCredentialDetails) {
	to.CredentialType = types.StringPointerValue((*string)(from.CredentialType))
	to.ConnectionEncryption = types.StringPointerValue((*string)(from.ConnectionEncryption))
	to.SingleSignOnType = types.StringPointerValue((*string)(from.SingleSignOnType))
	to.SkipTestConnection = types.BoolPointerValue(from.SkipTestConnection)
}
