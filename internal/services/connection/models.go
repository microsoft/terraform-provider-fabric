// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package connection

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
)

type connectionModel struct {
	ID               customtypes.UUID `tfsdk:"id"`
	DisplayName      types.String     `tfsdk:"display_name"`
	GatewayID        customtypes.UUID `tfsdk:"gateway_id"`
	ConnectivityType types.String     `tfsdk:"connectivity_type"`
	PrivacyLevel     types.String     `tfsdk:"privacy_level"`
}

func (to *connectionModel) set(from fabcore.Connection) {
	to.ID = customtypes.NewUUIDPointerValue(from.ID)
	to.DisplayName = types.StringPointerValue(from.DisplayName)
	to.GatewayID = customtypes.NewUUIDPointerValue(from.GatewayID)
	to.ConnectivityType = types.StringPointerValue((*string)(from.ConnectivityType))
	to.PrivacyLevel = types.StringPointerValue((*string)(from.PrivacyLevel))
}

type credentialsBasicModel struct {
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
}

type credentialsKeyModel struct {
	Key types.String `tfsdk:"key"`
}

type credentialsServicePrincipalModel struct {
	TenantID     types.String `tfsdk:"tenant_id"`
	ClientID     types.String `tfsdk:"client_id"`
	ClientSecret types.String `tfsdk:"client_secret"`
}

type credentialsSharedAccessSignatureModel struct {
	Token types.String `tfsdk:"token"`
}

type credentialsWindowsModel struct {
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
}

type credentialsEncryptedModel struct {
	Value types.String `tfsdk:"value"`
}
